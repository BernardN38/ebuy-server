-- +goose Up
CREATE TABLE product_types (
    id SERIAL PRIMARY KEY,
    type_name VARCHAR(30)
);
-- +goose Down
DROP TABLE product_types;