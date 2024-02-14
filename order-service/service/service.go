package service

import (
	"context"
	"database/sql"

	"github.com/BernardN38/ebuy-server/order-service/messaging"
	orders_sql "github.com/BernardN38/ebuy-server/order-service/sqlc/orders"
)

type OrderService struct {
	config          config
	orderDbQueries  queries
	rabbitmqEmitter messaging.MessageEmitter
}

// for testing purposes
type queries interface {
	GetAll(context.Context) ([]int32, error)
}
type config struct {
	rabbitmqExchange string
}

func New(db *sql.DB, rabbitmqEmitter messaging.MessageEmitter) (*OrderService, error) {
	orderQueries := orders_sql.New(db)
	return &OrderService{
		orderDbQueries:  orderQueries,
		rabbitmqEmitter: rabbitmqEmitter,
	}, nil
}
