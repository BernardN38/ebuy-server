package service

import "github.com/google/uuid"

type ProductParams struct {
	OwnerId     int
	Name        string
	Description string
	Price       int
	ProductType string
	MediaId     uuid.NullUUID
}
type ProductCreatedMsg struct {
	productId int
}

var (
	productTypeIDs = map[string]int32{
		"electronics": 1,
	}
)
