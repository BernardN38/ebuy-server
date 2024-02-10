CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    owner_id int NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    price int NOT NULL,
    product_type string NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    media_id UUID,
    FOREIGN KEY (product_types) REFERENCES product_types(type_name)
);
CREATE TABLE product_types (
    id SERIAL PRIMARY KEY,
    type_name VARCHAR(30) NOT NULL UNIQUE
);
