-- name: GetAll :many
SELECT * FROM products;
-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;
-- name: PatchProduct :exec
UPDATE products
SET
    name = COALESCE(NULLIF($2,''), name),                   
    description = COALESCE(NULLIF($3,''), description), 
    price = COALESCE(NULLIF($4,0), price)                  
WHERE
    id = $1;                                  

-- name: DeleteProduct :exec 
DELETE from products WHERE id = $1;
-- name: CreateProduct :exec
INSERT INTO products (name, description, price) VALUES ($1, $2, $3);