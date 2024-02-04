package service

import (
	"context"
	"errors"
	"testing"
	"time"

	products_sql "github.com/BernardN38/ebuy-server/product-service/sqlc/products"
)

type MockQueries struct {
	productStore []products_sql.Product
	sleepTime    int
}

func (m *MockQueries) CreateProduct(ctx context.Context, product products_sql.CreateProductParams) error {
	time.Sleep(time.Millisecond * time.Duration(m.sleepTime))
	m.productStore = append(m.productStore, products_sql.Product{
		ID:          int32(len(m.productStore) + 1),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   time.Now(),
	})
	return nil
}
func (m *MockQueries) GetProduct(ctx context.Context, id int32) (products_sql.Product, error) {
	time.Sleep(time.Millisecond * time.Duration(m.sleepTime))
	for _, product := range m.productStore {
		if product.ID == id {
			return product, nil
		}
	}
	return products_sql.Product{}, errors.New("product id not found in db")
}
func (m *MockQueries) GetRecentProducts(ctx context.Context, limit int32) ([]products_sql.Product, error) {
	time.Sleep(time.Millisecond * time.Duration(m.sleepTime))
	return m.productStore[:10], nil
}
func (m *MockQueries) PatchProduct(context.Context, products_sql.PatchProductParams) (int64, error) {
	return 1, nil
}
func (m *MockQueries) DeleteProduct(context.Context, products_sql.DeleteProductParams) (int64, error) {
	return 1, nil
}

// happy path
func TestCreateProductHappyPath(t *testing.T) {
	testCases := []ProductParams{
		{Name: "test", Description: "test description", Price: 100},
		{Name: "test2", Description: "second test description", Price: 200},
	}
	mockQueries := MockQueries{
		sleepTime: 100,
	}
	service := ProductService{
		productsDbQueries: &mockQueries,
	}
	for _, test := range testCases {
		ctx := context.Background()
		err := service.CreateProduct(ctx, ProductParams{
			Name:        test.Name,
			Description: test.Description,
			Price:       test.Price,
		})
		if err != nil {
			t.Error(err)
		}
	}
}

// bad path timeout is hit
func TestCreateProductUnHappyPath(t *testing.T) {
	testCases := []ProductParams{
		{Name: "test", Description: "test description", Price: 100},
		{Name: "test2", Description: "second test description", Price: 200},
	}
	mockQueries := MockQueries{
		sleepTime: 300,
	}
	service := ProductService{
		productsDbQueries: &mockQueries,
	}
	for _, test := range testCases {
		ctx := context.Background()
		err := service.CreateProduct(ctx, ProductParams{
			Name:        test.Name,
			Description: test.Description,
			Price:       test.Price,
		})
		//error should be returned
		if err == nil {
			t.Error("timeout reached and error is nil")
		}
	}
}

// happy path
func TestGetProductHappyPath(t *testing.T) {
	testCases := []products_sql.Product{
		{ID: 1, Name: "test", Description: "test description", Price: 100},
		{ID: 2, Name: "test2", Description: "second test description", Price: 200},
	}
	mockQueries := MockQueries{
		sleepTime: 100,
	}
	service := ProductService{
		productsDbQueries: &mockQueries,
	}
	for _, test := range testCases {
		ctx := context.Background()
		service.CreateProduct(ctx, ProductParams{Name: test.Name, Description: test.Description, Price: int(test.Price)})
		product, err := service.GetProduct(ctx, int(test.ID))
		if err != nil {
			t.Error(err)
		}
		if product.ID != test.ID || product.Name != test.Name || product.Description != test.Description || product.Price != test.Price {
			t.Errorf("expected id: %v, got id:%v", test, product)
		}
	}
}

// tests when timout is hit, expect error returned
func TestGetProductUnHappyPath(t *testing.T) {
	testCases := []products_sql.Product{
		{ID: 1, Name: "test", Description: "test description", Price: 100},
		{ID: 2, Name: "test2", Description: "second test description", Price: 200},
	}
	mockQueries := MockQueries{
		sleepTime: 300,
	}
	service := ProductService{
		productsDbQueries: &mockQueries,
	}
	for _, test := range testCases {
		ctx := context.Background()

		//populate mock db
		mockQueries.sleepTime = 0
		service.CreateProduct(ctx, ProductParams{Name: test.Name, Description: test.Description, Price: int(test.Price)})

		//exceed dealine timeout
		mockQueries.sleepTime = 300
		_, err := service.GetProduct(ctx, int(test.ID))
		if err.Error() != "context deadline exceeded" {
			t.Error("timeout reached and error is not conext deadline exceeded")
		}
	}
}
