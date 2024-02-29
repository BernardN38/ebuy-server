use std::time::Duration;

use sqlx::{postgres::PgPoolOptions, Error, PgPool};

pub async fn initialize_db() -> Result<PgPool, Error> {
    let db_connection_str = std::env::var("DATABASE_URL")
        .unwrap_or_else(|_| "postgres://bernardn:password@postgres:5432/user_service".to_string());
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
    Ok(pool)
}
