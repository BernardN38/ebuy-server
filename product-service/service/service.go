package service

import (
	"context"
	"database/sql"
	"time"

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

func (p *ProductService) CreateProduct(ctx context.Context, product ProductParams) error {
	// timeout for creating product set
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()
	//create response channels
	successCh := make(chan bool)
	errCh := make(chan error)
	// attempt create product
	go func() {
		err := p.productsDbQueries.CreateProduct(timeoutCtx, products_sql.CreateProductParams{
			Name:        product.Name,
			Description: product.Description,
			Price:       int32(product.Price),
		})
		if err != nil {
			errCh <- err
			return
		}
		successCh <- true
	}()

	// check if creating product takes too long or returns product
	select {
	case <-successCh:
		return nil
	case err := <-errCh:
		return err
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}

func (p *ProductService) GetProduct(ctx context.Context, productId int) (products_sql.Product, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	respCh := make(chan products_sql.Product)
	errCh := make(chan error)
	go func() {
		product, err := p.productsDbQueries.GetProduct(timeoutCtx, int32(productId))
		if err != nil {
			errCh <- err
			return
		}
		respCh <- product
	}()
	select {
	case product := <-respCh:
		return product, nil
	case err := <-errCh:
		return products_sql.Product{}, err
	case <-timeoutCtx.Done():
		return products_sql.Product{}, timeoutCtx.Err()
	}
}

func (p *ProductService) PatchProduct(ctx context.Context, productId int, productUpdate ProductParams) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()
	successCh := make(chan bool)
	errCh := make(chan error)
	go func() {
		err := p.productsDbQueries.PatchProduct(timeoutCtx, products_sql.PatchProductParams{
			ID:      int32(productId),
			Column2: productUpdate.Name,
			Column3: productUpdate.Description,
			Column4: int32(productUpdate.Price),
		})
		if err != nil {
			errCh <- err
			return
		}
		successCh <- true
	}()
	select {
	case <-successCh:
		return nil
	case err := <-errCh:
		return err
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}

func (p *ProductService) DeleteProduct(ctx context.Context, productid int) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()
	successCh := make(chan bool)
	errCh := make(chan error)
	go func() {
		err := p.productsDbQueries.DeleteProduct(timeoutCtx, int32(productid))
		if err != nil {
			errCh <- err
			return
		}
		successCh <- true
	}()
	select {
	case <-successCh:
		return nil
	case err := <-errCh:
		return err
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}
