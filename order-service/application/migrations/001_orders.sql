-- +goose Up
CREATE TABLE orders(
    id SERIAL PRIMARY KEY
);
-- +goose Down
DROP TABLE orders;