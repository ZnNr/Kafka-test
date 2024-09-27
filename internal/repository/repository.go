package repository

import "github.com/ZnNr/WB-test-L0/internal/models"

type Orders interface {
	AddOrder(order models.Order) error
	GetOrder(OrderUID string) (*models.Order, error)
	GetOrders() ([]models.Order, error)
}
