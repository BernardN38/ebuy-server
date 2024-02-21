use axum::{routing::get, Router};

use sqlx::postgres::PgPoolOptions;
use tokio::net::TcpListener;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

mod handlers;

// use futures_util::stream::stream::StreamExt;
use handlers::using_connection_extractor;
use std::time::Duration;

use crate::handlers::{create_user, get_user};

#[tokio::main]
async fn main() {
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "example_tokio_postgres=debug".into()),
        )
        .with(tracing_subscriber::fmt::layer())
        .init();

    let db_connection_str = std::env::var("DATABASE_URL")
        .unwrap_or_else(|_| "postgres://bernardn:password@localhost:5438/user_service".to_string());

    // set up connection pool
    let pool = PgPoolOptions::new()
        .max_connections(5)
        .acquire_timeout(Duration::from_secs(3))
        .connect(&db_connection_str)
        .await
        .expect("can't connect to database");

    sqlx::migrate!("./migrations")
        .run(&pool)
        .await
        .expect("error running migrations");

    // let user = User {
    //     id: 1,
    //     first_name: "firstName".to_string(),
    //     last_name: "lastName".to_string(),
    //     username: "username".to_string(),
    //     email: "email@gmail.com".to_string(),
    //     dob: Utc::now(),
    //     created_at: Utc::now(),
    //     last_updated_at: Utc::now(),
    // };
    // create_user(&user, &pool)
    //     .await
    //     .expect("error creating user");
    // let user = get_user(1, &pool).await.expect("error getting user");
    // let changed_user = User {
    //     id: 1,
    //     first_name: "chngfirstName".to_string(),
    //     last_name: "chnglastName".to_string(),
    //     username: "chngusername".to_string(),
    //     email: "email@gmail.com".to_string(),
    //     dob: Utc::now(),
    //     created_at: Utc::now(),
    //     last_updated_at: Utc::now(),
    // };
    // update_user(&changed_user, &pool)
    //     .await
    //     .expect("error updating users");

    // println!("{:?}", user);
    // build our application with some routes
    let app = Router::new()
        .route("/", get(get_user).post(create_user))
        .with_state(pool);

    // run it with hyper
    let listener = TcpListener::bind("127.0.0.1:3000").await.unwrap();
    tracing::debug!("listening on {}", listener.local_addr().unwrap());
    axum::serve(listener, app).await.unwrap();
}
