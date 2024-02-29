use std::sync::Arc;

use crate::{
    models::User,
    service::{self, CreateUserParams},
};
use axum::{
    async_trait,
    extract::{FromRef, FromRequestParts, Path, State},
    http::{request::Parts, StatusCode},
    response::IntoResponse,
    Json,
};

use sqlx::PgPool;

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

pub async fn get_user(
    State(pool): State<Arc<PgPool>>,
    Path(user_id): Path<i32>,
) -> Result<Json<User>, (StatusCode, String)> {
    match service::get_user(&pool, user_id).await {
        Ok(user) => Ok(Json(user)),
        Err(_) => Err((StatusCode::NOT_FOUND, "User not found".to_string())),
    }
}

pub async fn create_user(
    State(pool): State<Arc<PgPool>>,
    Json(user): Json<User>,
) -> Result<impl IntoResponse, (StatusCode, String)> {
    let new_user = CreateUserParams {
        first_name: user.first_name,
        last_name: user.last_name,
        username: user.username,
        dob: user.dob,
        email: user.email,
    };
    match service::create_user(&pool, new_user).await {
        Ok(user_id) => Ok(user_id.to_string()),
        Err(_) => Err((
            StatusCode::INTERNAL_SERVER_ERROR,
            "Failed to create user".to_string(),
        )),
    }
}

pub fn internal_error<E>(err: E) -> (StatusCode, String)
where
    E: std::error::Error,
{
    (StatusCode::INTERNAL_SERVER_ERROR, err.to_string())
}
