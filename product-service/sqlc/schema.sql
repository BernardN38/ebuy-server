CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    owner_id int NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    price int NOT NULL,
    product_type_id int REFERENCES product_types(id), 
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE product_types (
    id SERIAL PRIMARY KEY,
    type_name VARCHAR(30)
);