use crate::models::User;
use chrono::{DateTime, Utc};
use sqlx::{Error, PgPool, Row}; // Import Error type from sqlx

pub async fn get_user(pool: &PgPool, user_id: i32) -> Result<User, Error> {
    let q = "SELECT * FROM users WHERE id = $1";
    let row = sqlx::query(q).bind(user_id).fetch_one(pool).await?;

    let user = User {
        id: row.try_get("id")?, // Use try_get to handle potential errors
        first_name: row.try_get("first_name")?,
        last_name: row.try_get("last_name")?,
        username: row.try_get("username")?,
        email: row.try_get("email")?,
        dob: row.try_get("dob")?,
        created_at: row.try_get("created_at")?,
        last_updated_at: row.try_get("last_updated_at")?,
    };
    Ok(user)
}

pub struct CreateUserParams {
    pub first_name: String,
    pub last_name: String,
    pub username: String,
    pub email: String,
    pub dob: DateTime<Utc>,
}
pub async fn create_user(pool: &PgPool, user: CreateUserParams) -> Result<i32, Error> {
    let query = "INSERT INTO users ( first_name,last_name,username,email,dob,created_at,last_updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id";
    let result = sqlx::query(query)
        .bind(&user.first_name)
        .bind(&user.last_name)
        .bind(&user.username)
        .bind(&user.email)
        .bind(&user.dob)
        .bind(chrono::offset::Local::now())
        .bind(chrono::offset::Local::now())
        .fetch_one(pool)
        .await
        .expect("error creating user");
    let user_id = result.try_get("id")?;
    Ok(user_id)
}
