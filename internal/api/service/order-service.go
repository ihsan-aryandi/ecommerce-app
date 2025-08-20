package service

import (
	"ecommerce-app/internal/api/repository"

	"github.com/doug-martin/goqu/v9"
)

type OrderService struct {
	db              *goqu.Database
	orderRepository *repository.OrderRepository
}

func NewOrderService(db *goqu.Database, orderRepository *repository.OrderRepository) *OrderService {
	return &OrderService{
		db:              db,
		orderRepository: orderRepository,
	}
}

func (svc OrderService) GetTotal() []string {
	return svc.orderRepository.GetTotal()
}
