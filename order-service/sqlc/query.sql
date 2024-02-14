-- name: GetAll :many
SELECT * FROM orders;
-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- -- name: PatchProduct :one
-- with updated as (
-- UPDATE products
-- SET
--     name = COALESCE(NULLIF($3,''), name),                   
--     description = COALESCE(NULLIF($4,''), description), 
--     price = COALESCE(NULLIF($5,0), price)                  
-- WHERE
--     id = $1 and owner_id = $2 returning id) 
--     select count(*)
-- from updated;                            



-- name: CreateProduct :one
INSERT INTO orders (id) VALUES ($1) RETURNING id;
