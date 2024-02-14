package handler

import (
	"encoding/json"
	"net/http"

	"github.com/BernardN38/ebuy-server/order-service/service"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	// CreateProduct(context.Context, service.ProductParams) error
	// GetProduct(context.Context, int) (products_sql.Product, error)
	// GetRecentProducts(context.Context, int) ([]products_sql.Product, error)
	// PatchProduct(context.Context, int, service.ProductParams) error
	// DeleteProduct(context.Context, int, int) error
	// GetProductTypes(context.Context) ([]products_sql.ProductType, error)
}

type Handler struct {
	orderService Service
	validator    *validator.Validate
}

func New(service *service.OrderService) *Handler {
	v := validator.New()
	return &Handler{
		orderService: service,
		validator:    v,
	}
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}{
		Name:   "product-service",
		Status: "up",
	})
}
