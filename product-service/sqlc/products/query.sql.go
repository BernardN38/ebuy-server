// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package products_sql

import (
	"context"

	"github.com/google/uuid"
)

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (owner_id, name, description, price, media_id, product_type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
`

type CreateProductParams struct {
	OwnerID     int32         `json:"ownerId"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       int32         `json:"price"`
	MediaID     uuid.NullUUID `json:"mediaId"`
	ProductType string        `json:"productType"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, createProduct,
		arg.OwnerID,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.MediaID,
		arg.ProductType,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const deleteProduct = `-- name: DeleteProduct :one
with deleted as (
   DELETE FROM products WHERE id  = $1 and owner_id = $2
   returning id
)
select count(*)
from deleted
`

type DeleteProductParams struct {
	ID      int32 `json:"id"`
	OwnerID int32 `json:"ownerId"`
}

func (q *Queries) DeleteProduct(ctx context.Context, arg DeleteProductParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, deleteProduct, arg.ID, arg.OwnerID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getAll = `-- name: GetAll :many
SELECT id, owner_id, name, description, price, product_type, created_at, media_id FROM products
`

func (q *Queries) GetAll(ctx context.Context) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.OwnerID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.ProductType,
			&i.CreatedAt,
			&i.MediaID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProduct = `-- name: GetProduct :one
SELECT id, owner_id, name, description, price, product_type, created_at, media_id FROM products WHERE id = $1
`

func (q *Queries) GetProduct(ctx context.Context, id int32) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProduct, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.OwnerID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.ProductType,
		&i.CreatedAt,
		&i.MediaID,
	)
	return i, err
}

const getProductTypes = `-- name: GetProductTypes :many
SELECT id, type_name FROM product_types
`

func (q *Queries) GetProductTypes(ctx context.Context) ([]ProductType, error) {
	rows, err := q.db.QueryContext(ctx, getProductTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ProductType
	for rows.Next() {
		var i ProductType
		if err := rows.Scan(&i.ID, &i.TypeName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRecentProducts = `-- name: GetRecentProducts :many
SELECT id, owner_id, name, description, price, product_type, created_at, media_id 
FROM products
ORDER BY created_at DESC
LIMIT $1
`

func (q *Queries) GetRecentProducts(ctx context.Context, limit int32) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getRecentProducts, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.OwnerID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.ProductType,
			&i.CreatedAt,
			&i.MediaID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const patchProduct = `-- name: PatchProduct :one
with updated as (
UPDATE products
SET
    name = COALESCE(NULLIF($3,''), name),                   
    description = COALESCE(NULLIF($4,''), description), 
    price = COALESCE(NULLIF($5,0), price)                  
WHERE
    id = $1 and owner_id = $2 returning id) 
    select count(*)
from updated
`

type PatchProductParams struct {
	ID      int32       `json:"id"`
	OwnerID int32       `json:"ownerId"`
	Column3 interface{} `json:"column3"`
	Column4 interface{} `json:"column4"`
	Column5 interface{} `json:"column5"`
}

func (q *Queries) PatchProduct(ctx context.Context, arg PatchProductParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, patchProduct,
		arg.ID,
		arg.OwnerID,
		arg.Column3,
		arg.Column4,
		arg.Column5,
	)
	var count int64
	err := row.Scan(&count)
	return count, err
}
