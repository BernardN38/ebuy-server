package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/BernardN38/ebuy-server/product-service/messaging"
	products_sql "github.com/BernardN38/ebuy-server/product-service/sqlc/products"
)

type ProductService struct {
	config            config
	productsDbQueries queries
	productTypeIds    map[string]int32
	rabbitmqEmitter   messaging.MessageEmitter
}

// for testing purposes
type queries interface {
	CreateProduct(context.Context, products_sql.CreateProductParams) (int32, error)
	GetProduct(context.Context, int32) (products_sql.Product, error)
	PatchProduct(context.Context, products_sql.PatchProductParams) (int64, error)
	DeleteProduct(context.Context, products_sql.DeleteProductParams) (int64, error)
	GetRecentProducts(context.Context, int32) ([]products_sql.Product, error)
	GetProductTypes(context.Context) ([]products_sql.ProductType, error)
}
type config struct {
	rabbitmqExchange string
}

func New(db *sql.DB, rabbitmqEmitter messaging.MessageEmitter) (*ProductService, error) {
	productQueries := products_sql.New(db)
	// productQueries.GetProductTypes()
	productTypes, err := productQueries.GetProductTypes(context.Background())
	if err != nil {
		return nil, err
	}
	productTypeIDs := make(map[string]int32)
	for _, productType := range productTypes {
		productTypeIDs[productType.TypeName] = productType.ID
	}
	return &ProductService{
		productsDbQueries: productQueries,
		productTypeIds:    productTypeIDs,
		rabbitmqEmitter:   rabbitmqEmitter,
	}, nil
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
		productId, err := p.productsDbQueries.CreateProduct(timeoutCtx, products_sql.CreateProductParams{
			OwnerID:     int32(product.OwnerId),
			Name:        product.Name,
			Description: product.Description,
			Price:       int32(product.Price),
			MediaID:     product.MediaId,
			ProductType: product.ProductType,
		})
		if err != nil {
			errCh <- err
			return
		}
		msg := ProductCreatedMsg{
			productId: int(productId),
		}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
		}
		err = p.rabbitmqEmitter.SendMessage(timeoutCtx, msgBytes, p.config.rabbitmqExchange, "product_events", "product.created")
		if err != nil {
			log.Println(err)
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
func (p *ProductService) GetProductTypes(ctx context.Context) ([]products_sql.ProductType, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	respCh := make(chan []products_sql.ProductType)
	errCh := make(chan error)
	go func() {
		productTypes, err := p.productsDbQueries.GetProductTypes(timeoutCtx)
		if err != nil {
			errCh <- err
			return
		}
		respCh <- productTypes
	}()
	select {
	case productTypes := <-respCh:
		return productTypes, nil
	case err := <-errCh:
		return nil, err
	case <-timeoutCtx.Done():
		return nil, timeoutCtx.Err()
	}
}

func (p *ProductService) GetRecentProducts(ctx context.Context, limit int) ([]products_sql.Product, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	respCh := make(chan []products_sql.Product)
	errCh := make(chan error)
	go func() {
		recentProducts, err := p.productsDbQueries.GetRecentProducts(timeoutCtx, int32(limit))
		if err != nil {
			errCh <- err
			return
		}

		respCh <- recentProducts
	}()
	select {
	case products := <-respCh:
		return products, nil
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
