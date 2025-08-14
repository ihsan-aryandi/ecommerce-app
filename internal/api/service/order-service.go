package service

import (
	"database/sql"
	"ecommerce-app/internal/api/repository"
)

type OrderService struct {
	db              *sql.DB
	orderRepository *repository.OrderRepository
}

func NewOrderService(db *sql.DB, orderRepository *repository.OrderRepository) *OrderService {
	return &OrderService{
		db:              db,
		orderRepository: orderRepository,
	}
}

func (svc OrderService) GetTotal() []string {
	return svc.orderRepository.GetTotal()
}
