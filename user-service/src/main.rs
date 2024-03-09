use std::sync::Arc;

use amqprs::{
    channel::BasicQosArguments,
    connection::{Connection, OpenConnectionArguments},
};
use axum::{
    routing::{get, post},
    Router,
};
use tokio::net::TcpListener;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

mod handlers;
use crate::{
    database::initialize_db,
    handlers::{check_health, create_user, get_user},
    messaging::run_rabbitmq_consumer,
};
mod database;
mod messaging;
mod models;
mod service;
#[tokio::main]
async fn main() {
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "example_tokio_postgres=debug".into()),
        )
        .with(tracing_subscriber::fmt::layer())
        .init();

    // set up postgres connection pool
    let pool = Arc::new(initialize_db().await.expect("error connecting to database"));

    // Clone pool for RabbitMQ consumer
    let pool_clone = Arc::clone(&pool);

    //initialize rabbitmq connection, channel and consumer
    let rabbitmq_connection = Connection::open(&OpenConnectionArguments::new(
        "rabbitmq", 5672, "guest", "guest",
    ))
    .await
    .unwrap();
    // open a channel on the connection
    let channel = rabbitmq_connection.open_channel(None).await.unwrap();
    // let channel2 = rabbitmq_connection.open_channel(None).await.unwrap();
    // Set the desired prefetch count
    channel
        .basic_qos(BasicQosArguments::new(0, 5, false))
        .await
        .expect("Failed to set prefetch count");
    // channel2
    //     .basic_qos(BasicQosArguments::new(0, 5, false))
    //     .await
    //     .expect("Failed to set prefetch count");
    run_rabbitmq_consumer(pool_clone.clone(), &rabbitmq_connection, &channel).await;
    // run_rabbitmq_consumer(pool_clone, &rabbitmq_connection, &channel2).await;
    let app = Router::new()
        .route("/api/v1/users/health", get(check_health))
        .route("/api/v1/users/:user_id", get(get_user))
        .route("/api/v1/users", post(create_user))
        .with_state(pool);

    // run it with hyper
    let listener = TcpListener::bind("0.0.0.0:8080").await.unwrap();
    println!("{:?}", "listening on port 8080");
    tracing::debug!("listening on {}", listener.local_addr().unwrap());
    axum::serve(listener, app).await.unwrap();
}
