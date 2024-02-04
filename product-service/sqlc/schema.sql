CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    owner_id int NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    price int NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);