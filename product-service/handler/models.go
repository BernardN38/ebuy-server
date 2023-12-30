package handler

type ProductPayload struct {
	Name        string `json:"name" validate:"required,min=5,max=20"`
	Description string `json:"description" validate:"required,min=5,max=50"`
	Price       int    `json:"price" validate:"required"`
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
