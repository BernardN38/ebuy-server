use axum::{
    async_trait,
    extract::{FromRef, FromRequestParts, State},
    http::{request::Parts, StatusCode},
    Json,
};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::{prelude::FromRow, PgPool, Row};

#[derive(FromRow, Serialize, Deserialize, Debug)]
pub struct User {
    pub id: i32,
    pub first_name: String,
    pub last_name: String,
    pub username: String,
    pub email: String,
    pub dob: DateTime<Utc>,
    pub created_at: DateTime<Utc>,
    pub last_updated_at: DateTime<Utc>,
}

pub struct DatabaseConnection(sqlx::pool::PoolConnection<sqlx::Postgres>);

#[async_trait]
impl<S> FromRequestParts<S> for DatabaseConnection
where
    PgPool: FromRef<S>,
    S: Send + Sync,
{
    type Rejection = (StatusCode, String);

    async fn from_request_parts(_parts: &mut Parts, state: &S) -> Result<Self, Self::Rejection> {
        let pool = PgPool::from_ref(state);

        let conn = pool.acquire().await.map_err(internal_error)?;

        Ok(Self(conn))
    }
}

// we can extract the connection pool with `State`
pub async fn get_user(State(pool): State<PgPool>) -> Result<Json<User>, (StatusCode, String)> {
    let q = "SELECT * FROM users WHERE id = $1";
    let query = sqlx::query(q).bind(1);
    let row = query.fetch_one(&pool).await.expect("error fetching users");
    let user = User {
        id: row.get("id"),
        first_name: row.get("first_name"),
        last_name: row.get("last_name"),
        username: row.get("username"),
        email: row.get("email"),
        dob: row.get("dob"),
        created_at: row.get("created_at"),
        last_updated_at: row.get("last_updated_at"),
    };
    Ok(Json(user))
}
pub async fn create_user(
    State(pool): State<PgPool>,
    Json(user): Json<User>,
) -> Result<(), (StatusCode, String)> {
    let query = "INSERT INTO users ( first_name,last_name,username,email,dob,created_at,last_updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)";
    sqlx::query(query)
        .bind(&user.first_name)
        .bind(&user.last_name)
        .bind(&user.username)
        .bind(&user.email)
        .bind(&user.dob)
        .bind(&user.created_at)
        .bind(&user.last_updated_at)
        .execute(&pool)
        .await
        .expect("error creating user");
    Ok(())
}

pub fn internal_error<E>(err: E) -> (StatusCode, String)
where
    E: std::error::Error,
{
    (StatusCode::INTERNAL_SERVER_ERROR, err.to_string())
}
