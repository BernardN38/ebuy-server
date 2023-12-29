package service

import (
	"database/sql"

	products_sql "github.com/BernardN38/ebuy-server/sqlc/products"
)

type ProductService struct {
	productsDbQueries *products_sql.Queries
}

func New(db *sql.DB) *ProductService {
	productQueries := products_sql.New(db)
	return &ProductService{
		productsDbQueries: productQueries,
	}
}
