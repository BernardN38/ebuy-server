-- +goose Up
ALTER TABLE products ADD CONSTRAINT product_type FOREIGN KEY (product_type) REFERENCES product_types(type_name);

-- +goose Down
ALTER TABLE products DROP CONSTRAINT product_type;