use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::prelude::FromRow;

#[derive(FromRow, Serialize, Deserialize, Debug)]
pub struct User {
    #[serde(skip_deserializing)]
    pub id: i32,
    pub first_name: String,
    pub last_name: String,
    pub username: String,
    pub email: String,
    pub dob: DateTime<Utc>,
    #[serde(skip_deserializing)]
    pub created_at: DateTime<Utc>,
    #[serde(skip_deserializing)]
    pub last_updated_at: DateTime<Utc>,
}
