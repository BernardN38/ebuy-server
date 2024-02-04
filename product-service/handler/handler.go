package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BernardN38/ebuy-server/product-service/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
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
	_, claims, _ := jwtauth.FromContext(r.Context())
	fmt.Println(claims["user_id"])
	userId, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "bad token", http.StatusBadRequest)
		return
	}
	// userIdInt, err := strconv.Atoi(userId)
	// if !ok {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// decode json body
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

	//change this to correct porduct id in future
	json.NewEncoder(w).Encode(
		ProductCreationResponse{
			PoductID:     0,
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
		product,
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
