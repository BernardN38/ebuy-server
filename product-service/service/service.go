package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	products_sql "github.com/BernardN38/ebuy-server/product-service/sqlc/products"
)

type ProductService struct {
	productsDbQueries queries
}

// for testing purposes
type queries interface {
	CreateProduct(context.Context, products_sql.CreateProductParams) error
	GetProduct(context.Context, int32) (products_sql.Product, error)
	PatchProduct(context.Context, products_sql.PatchProductParams) (int64, error)
	DeleteProduct(context.Context, products_sql.DeleteProductParams) (int64, error)
	GetRecentProducts(ctx context.Context, limit int32) ([]products_sql.Product, error)
}

func New(db *sql.DB) *ProductService {
	productQueries := products_sql.New(db)
	// productQueries.GetRecentProducts()
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
			OwnerID:     int32(product.OwnerId),
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

func (p *ProductService) GetRecentProducts(ctx context.Context) ([]products_sql.Product, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	respCh := make(chan []products_sql.Product)
	errCh := make(chan error)
	go func() {
		recentProducts, err := p.productsDbQueries.GetRecentProducts(timeoutCtx, 20)
		if err != nil {
			errCh <- err
			return
		}

		respCh <- recentProducts
	}()
	select {
	case product := <-respCh:
		return product, nil
	case err := <-errCh:
		return nil, err
	case <-timeoutCtx.Done():
		return nil, timeoutCtx.Err()
	}
}

func (p *ProductService) PatchProduct(ctx context.Context, productId int, productUpdate ProductParams) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()
	successCh := make(chan bool)
	errCh := make(chan error)
	go func() {
		count, err := p.productsDbQueries.PatchProduct(timeoutCtx, products_sql.PatchProductParams{
			ID:      int32(productId),
			OwnerID: int32(productUpdate.OwnerId),
			Column3: productUpdate.Name,
			Column4: productUpdate.Description,
			Column5: int32(productUpdate.Price),
		})
		if count == 0 {
			errCh <- errors.New("no rows found")
			return
		}
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

func (p *ProductService) DeleteProduct(ctx context.Context, productid int, ownerId int) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()
	successCh := make(chan bool)
	errCh := make(chan error)
	go func() {
		count, err := p.productsDbQueries.DeleteProduct(timeoutCtx, products_sql.DeleteProductParams{
			ID:      int32(productid),
			OwnerID: int32(ownerId),
		})
		if count == 0 {
			errCh <- errors.New("no rows found")
			return
		}
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
