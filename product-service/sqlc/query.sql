-- name: GetAll :many
SELECT * FROM products;
-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;
-- name: GetRecentProducts :many
SELECT * 
FROM products
ORDER BY created_at DESC
LIMIT $1;

-- name: PatchProduct :one
with updated as (
UPDATE products
SET
    name = COALESCE(NULLIF($3,''), name),                   
    description = COALESCE(NULLIF($4,''), description), 
    price = COALESCE(NULLIF($5,0), price)                  
WHERE
    id = $1 and owner_id = $2 returning id) 
    select count(*)
from updated;                            

-- name: DeleteProduct :one
with deleted as (
   DELETE FROM products WHERE id  = $1 and owner_id = $2
   returning id
)
select count(*)
from deleted;

-- name: CreateProduct :one
INSERT INTO products (owner_id, name, description, price, media_id, product_type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;

-- name: GetProductTypes :many
SELECT * FROM product_types;