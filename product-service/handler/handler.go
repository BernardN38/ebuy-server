package handler

import (
	"encoding/json"
	"net/http"

	"github.com/BernardN38/ebuy-server/service"
)

type Handler struct {
	productService *service.ProductService
}

func New(service *service.ProductService) *Handler {
	return &Handler{
		productService: service,
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

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {

}
