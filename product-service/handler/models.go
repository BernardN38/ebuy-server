package handler

import "github.com/google/uuid"

type ProductPayload struct {
	OwnerId     int           `json:"ownerId" validate:"required"`
	Name        string        `json:"name" validate:"required,min=5,max=20"`
	Description string        `json:"description" validate:"required,min=5,max=50"`
	Price       int           `json:"price" validate:"required"`
	ProductType string        `json:"productType" validate:"required"`
	MediaId     uuid.NullUUID `json:"mediaId"`
}

type ProductCreationResponse struct {
	PoductID     int    `json:"productId"`
	ErrorMessage string `json:"errorMessage"`
}
type ProductResponse struct {
	PoductID    int    `json:"productId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}
