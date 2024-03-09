use std::{fmt::Display, sync::Arc};

use crate::{
    models::{self, User, UserDto},
    service::{self, CreateUserParams},
};
use axum::{
    async_trait,
    extract::{FromRef, FromRequestParts, Path, State},
    http::{request::Parts, StatusCode},
    response::{IntoResponse, Response},
    Json,
};

use jsonwebtoken::{decode, DecodingKey, Validation};
use once_cell::sync::Lazy;
use serde::{Deserialize, Serialize};
use serde_json::json;
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

pub async fn check_health() -> Result<Json<serde_json::Value>, (StatusCode, String)> {
    let service_name = "user-service";
    let status = "up";

    let response_json = serde_json::json!({
        "name": service_name,
        "status": status,
    });

    Ok(Json(response_json))
}

pub async fn get_user(
    State(pool): State<Arc<PgPool>>,
    Path(user_id): Path<i32>,
    claims: JwtClaims,
) -> Result<Json<UserDto>, (StatusCode, String)> {
    println!("{:?}", claims);
    match service::get_user(&pool, user_id).await {
        Ok(user) => {
            let user_response = models::UserDto {
                first_name: user.first_name,
                last_name: user.last_name,
                email: user.email,
                username: user.username,
                dob: user.dob,
            };
            Ok(Json(user_response))
        }
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

static KEYS: Lazy<Keys> = Lazy::new(|| {
    // let secret = std::env::var("JWT_SECRET").expect("JWT_SECRET must be set");
    Keys::new("qwertyuiopasdfghjklzxcvbnm123456qwertyuiopasdfghjklzxcvbnm123456".as_bytes())
});

impl Display for Claims {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Email: {}\nCompany: {}", self.sub, self.company)
    }
}

#[async_trait]
impl<S> FromRequestParts<S> for JwtClaims
where
    S: Send + Sync,
{
    type Rejection = AuthError;

    async fn from_request_parts(
        parts: &mut Parts,
        _state: &S,
    ) -> Result<JwtClaims, Self::Rejection> {
        // Extract the cookies from the headers
        let cookie = parts
            .headers
            .get("Cookie")
            .ok_or(AuthError::MissingCredentials)?;

        let jwt_cookie = cookie
            .to_str()
            .map_err(|_| AuthError::MissingCredentials)?
            .split(';')
            .find(|cookie| cookie.trim().starts_with("jwt="))
            .ok_or(AuthError::MissingCredentials)?;
        // Extract the token from the jwt cookie
        let token = jwt_cookie
            .split('=')
            .nth(1)
            .ok_or(AuthError::MissingCredentials)?;
        // Decode the user data
        let token_data = decode::<JwtClaims>(
            token,
            &KEYS.decoding,
            &Validation::new(jsonwebtoken::Algorithm::HS512),
        );
        // Match on the result to handle errors
        match token_data {
            Ok(token_data) => Ok(token_data.claims),
            Err(err) => {
                println!("Error decoding token: {:?}", err);
                Err(AuthError::InvalidToken)
            }
        }

        // Ok(token_data.claims)
    }
}
#[derive(Debug, Serialize, Deserialize)]
pub struct JwtClaims {
    exp: usize,
    user_id: i32,
}
impl IntoResponse for AuthError {
    fn into_response(self) -> Response {
        let (status, error_message) = match self {
            AuthError::WrongCredentials => (StatusCode::UNAUTHORIZED, "Wrong credentials"),
            AuthError::MissingCredentials => (StatusCode::BAD_REQUEST, "Missing credentials"),
            AuthError::TokenCreation => (StatusCode::INTERNAL_SERVER_ERROR, "Token creation error"),
            AuthError::InvalidToken => (StatusCode::BAD_REQUEST, "Invalid token"),
        };
        let body = Json(json!({
            "error": error_message,
        }));
        (status, body).into_response()
    }
}

struct Keys {
    decoding: DecodingKey,
}

impl Keys {
    fn new(secret: &[u8]) -> Self {
        Self {
            decoding: DecodingKey::from_secret(secret),
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Claims {
    sub: String,
    company: String,
    exp: usize,
}

#[derive(Debug, Serialize)]
struct AuthBody {
    access_token: String,
    token_type: String,
}

#[derive(Debug)]
pub enum AuthError {
    WrongCredentials,
    MissingCredentials,
    TokenCreation,
    InvalidToken,
}
