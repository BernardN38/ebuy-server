package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/BernardN38/ebuy-server/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	productService *service.ProductService
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
	// decode json body
	var createProductBody ProductPayload
	err := json.NewDecoder(r.Body).Decode(&createProductBody)
	if err != nil {
		log.Println(err)
		http.Error(w, "error decoding body", http.StatusBadRequest)
		return
	}
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
	json.NewEncoder(w).Encode(
		ProductCreationResponse{
			PoductID:     1,
			ErrorMessage: "",
		},
	)
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
		ProductResponse{
			PoductID:    int(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Price:       int(product.Price),
		},
	)
}

func (h *Handler) PatchProduct(w http.ResponseWriter, r *http.Request) {
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
	productId := chi.URLParam(r, "productId")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		log.Println(err)
		http.Error(w, "product id invalid", http.StatusBadRequest)
		return
	}
	err = h.productService.DeleteProduct(r.Context(), productIdInt)
	if err != nil {
		http.Error(w, "unabble to delete product", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
