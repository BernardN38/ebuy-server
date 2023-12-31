// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package products_sql

import (
	"context"
)

const createProduct = `-- name: CreateProduct :exec
INSERT INTO products (name, description, price) VALUES ($1, $2, $3)
`

type CreateProductParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int32  `json:"price"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) error {
	_, err := q.db.ExecContext(ctx, createProduct, arg.Name, arg.Description, arg.Price)
	return err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE from products WHERE id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteProduct, id)
	return err
}

const getAll = `-- name: GetAll :many
SELECT id, name, description, price, created_at FROM products
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
			&i.Name,
			&i.Description,
			&i.Price,
			&i.CreatedAt,
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
SELECT id, name, description, price, created_at FROM products WHERE id = $1
`

func (q *Queries) GetProduct(ctx context.Context, id int32) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProduct, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.CreatedAt,
	)
	return i, err
}

const patchProduct = `-- name: PatchProduct :exec
UPDATE products
SET
    name = COALESCE(NULLIF($2,''), name),                   
    description = COALESCE(NULLIF($3,''), description), 
    price = COALESCE(NULLIF($4,0), price)                  
WHERE
    id = $1
`

type PatchProductParams struct {
	ID      int32       `json:"id"`
	Column2 interface{} `json:"column2"`
	Column3 interface{} `json:"column3"`
	Column4 interface{} `json:"column4"`
}

func (q *Queries) PatchProduct(ctx context.Context, arg PatchProductParams) error {
	_, err := q.db.ExecContext(ctx, patchProduct,
		arg.ID,
		arg.Column2,
		arg.Column3,
		arg.Column4,
	)
	return err
}
