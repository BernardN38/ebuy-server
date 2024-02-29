use crate::service::{self, CreateUserParams};
use amqprs::{
    callbacks::{DefaultChannelCallback, DefaultConnectionCallback},
    channel::{
        BasicConsumeArguments, Channel, ExchangeDeclareArguments, QueueBindArguments,
        QueueDeclareArguments,
    },
    connection::Connection,
    consumer::AsyncConsumer,
    BasicProperties, Deliver,
};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::PgPool;
use std::sync::Arc;

pub async fn run_rabbitmq_consumer(pool: Arc<PgPool>, connection: &Connection, channel: &Channel) {
    let exchange_name = "user_events";
    let exchange_args = ExchangeDeclareArguments::new(exchange_name, "topic")
        .durable(true)
        .finish();
    connection
        .register_callback(DefaultConnectionCallback)
        .await
        .unwrap();

    channel
        .register_callback(DefaultChannelCallback)
        .await
        .unwrap();
    channel.exchange_declare(exchange_args).await.unwrap();
    // declare a queue
    let (queue_name, _, _) = channel
        .queue_declare(QueueDeclareArguments::new("user-service"))
        .await
        .unwrap()
        .unwrap();

    // bind the queue to exchange
    let rounting_key = "user.#";
    channel
        .queue_bind(QueueBindArguments::new(
            &queue_name,
            exchange_name,
            rounting_key,
        ))
        .await
        .unwrap();

    //////////////////////////////////////////////////////////////////
    // start consumer with given name
    let args = BasicConsumeArguments::new(&queue_name, "example_basic_pub_sub")
        .auto_ack(true)
        .finish();
    let custom_consumer = CustomConsumer::new(pool);
    // res type is : impl Future<Output = Result<String, Error>>
    let _ = channel.basic_consume(custom_consumer, args).await;
}

pub(crate) struct CustomConsumer {
    pool: Arc<PgPool>,
}

impl CustomConsumer {
    pub fn new(pool: Arc<PgPool>) -> Self {
        CustomConsumer { pool }
    }
}

impl AsyncConsumer for CustomConsumer {
    fn consume<'life0, 'life1, 'async_trait>(
        &'life0 mut self,
        _channel: &'life1 Channel,
        deliver: Deliver,
        _basic_properties: BasicProperties,
        content: Vec<u8>,
    ) -> std::pin::Pin<Box<dyn std::future::Future<Output = ()> + std::marker::Send + 'async_trait>>
    where
        'life0: 'async_trait,
        'life1: 'async_trait,
        Self: 'async_trait,
    {
        Box::pin(async move {
            let routing_key = deliver.routing_key();
            let message = std::str::from_utf8(&content).unwrap();
            match routing_key.as_str() {
                "user.created" => {
                    let pool = self.pool.clone();
                    let message = message.to_owned();
                    tokio::spawn(async move {
                        handle_create_user(&pool, &message).await;
                    });
                    // acknowledge_message(channel, deliver)
                    //     .await
                    //     .expect("error acknowledging")
                }
                _ => println!("Routing key not matched"),
            }
        })
    }
}

#[derive(Serialize, Deserialize, Debug)]
struct UserCreatedMsg {
    #[serde(rename = "firstName")]
    first_name: String,

    #[serde(rename = "lastName")]
    last_name: String,

    #[serde(rename = "username")]
    username: String,

    #[serde(rename = "email")]
    email: String,

    #[serde(rename = "dob")]
    dob: DateTime<Utc>,
}
async fn handle_create_user(pool: &PgPool, message: &str) {
    println!("{:?}", message);
    let user_msg: Result<UserCreatedMsg, serde_json::Error> = serde_json::from_str(message);
    let msg = match user_msg {
        Ok(user) => user,
        Err(err) => {
            eprintln!("Failed to deserialize JSON: {}", err);
            return;
        }
    };
    let new_user = CreateUserParams {
        first_name: msg.first_name,
        last_name: msg.last_name,
        username: msg.username,
        dob: msg.dob,
        email: msg.email,
    };
    let result = service::create_user(pool, new_user).await;
    let _result = match result {
        Ok(user_id) => println!("user created successfully, id: {:?}", user_id),
        Err(err) => {
            println!("{:?}", err)
        }
    };
}

// async fn acknowledge_message(
//     channel: &Channel,
//     deliver: Deliver,
// ) -> Result<(), amqprs::error::Error> {
//     channel
//         .basic_ack(BasicAckArguments::new(deliver.delivery_tag(), false))
//         .await?;
//     Ok(())
// }
