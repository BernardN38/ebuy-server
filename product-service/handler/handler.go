package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BernardN38/ebuy-server/product-service/service"
	products_sql "github.com/BernardN38/ebuy-server/product-service/sqlc/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	CreateProduct(context.Context, service.ProductParams) error
	GetProduct(context.Context, int) (products_sql.Product, error)
	GetRecentProducts(context.Context, int) ([]products_sql.Product, error)
	PatchProduct(context.Context, int, service.ProductParams) error
	DeleteProduct(context.Context, int, int) error
	GetProductTypes(context.Context) ([]products_sql.ProductType, error)
}

type Handler struct {
	productService Service
	validator      *validator.Validate
}

func New(service *service.ProductService) *Handler {
	v := validator.New()
	return &Handler{
		productService: service,
		validator:      v,
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
	_, claims, _ := jwtauth.FromContext(r.Context())
	userId, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "bad token", http.StatusBadRequest)
		return
	}
	var createProductBody ProductPayload
	err := json.NewDecoder(r.Body).Decode(&createProductBody)
	if err != nil {
		log.Println(err)
		http.Error(w, "error decoding body", http.StatusBadRequest)
		return
	}
	createProductBody.OwnerId = int(userId)
	//validate json body
	err = h.validator.Struct(createProductBody)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// send request to service layer
	err = h.productService.CreateProduct(r.Context(), service.ProductParams(createProductBody))
	if err != nil {
		//create resposes for specific error reasons
		log.Println(err)
		http.Error(w, "error creating product", http.StatusBadRequest)
		return
	}
	//write success response
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "productId")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		log.Println(err)
		http.Error(w, "product id invalid", http.StatusBadRequest)
		return
	}
	product, err := h.productService.GetProduct(r.Context(), productIdInt)
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to get product", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		product,
	)
}

func (h *Handler) GetProductTypes(w http.ResponseWriter, r *http.Request) {
	productsTypes, err := h.productService.GetProductTypes(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to get product", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	types := make([]string, 0)
	for _, v := range productsTypes {
		types = append(types, v.TypeName)
	}
	json.NewEncoder(w).Encode(
		map[string]any{
			"productTypes": types,
		},
	)
}
func (h *Handler) GetRecentProducts(w http.ResponseWriter, r *http.Request) {
	// Get query parameters using chi
	queryParams := r.URL.Query()

	// Extract specific parameters
	limit := queryParams.Get("limit")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Println(err)
		limitInt = 10
	}
	if limitInt > 100 {
		http.Error(w, "too many entries requested", http.StatusBadRequest)
		return
	}
	products, err := h.productService.GetRecentProducts(r.Context(), limitInt)
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to get product", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		map[string]any{
			"products": products,
		},
	)
}
func (h *Handler) PatchProduct(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	fmt.Println(claims["user_id"])
	userId, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "bad token", http.StatusBadRequest)
		return
	}
	productId := chi.URLParam(r, "productId")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		log.Println(err)
		http.Error(w, "product id invalid", http.StatusBadRequest)
		return
	}
	var productUpdate ProductPayload
	json.NewDecoder(r.Body).Decode(&productUpdate)
	err = h.productService.PatchProduct(r.Context(), productIdInt, service.ProductParams{
		OwnerId:     int(userId),
		Name:        productUpdate.Name,
		Description: productUpdate.Description,
		Price:       productUpdate.Price,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "unable to update product", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	fmt.Println(claims["user_id"])
	userId, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "bad token", http.StatusBadRequest)
		return
	}
	productId := chi.URLParam(r, "productId")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		log.Println(err)
		http.Error(w, "product id invalid", http.StatusBadRequest)
		return
	}
	err = h.productService.DeleteProduct(r.Context(), productIdInt, int(userId))
	if err != nil {
		http.Error(w, "unabble to delete product", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
