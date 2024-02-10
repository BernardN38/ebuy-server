-- +goose Up
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    owner_id int NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    price int NOT NULL,
    product_type VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    media_id UUID
);
-- +goose Down
DROP TABLE products;