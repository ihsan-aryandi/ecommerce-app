package service

import (
	"ecommerce-app/internal/api/repository"

	"github.com/doug-martin/goqu/v9"
)

type OrderService struct {
	db                       *goqu.Database
	orderRepository          *repository.OrderRepository
	productVariantRepository *repository.ProductVariantRepository
	rajaOngkirRepository     *repository.RajaOngkirRepository
}

func NewOrderService(
	db *goqu.Database,
	orderRepository *repository.OrderRepository,
	productVariantRepository *repository.ProductVariantRepository,
	rajaOngkirRepository *repository.RajaOngkirRepository,
) *OrderService {
	return &OrderService{
		db:                       db,
		orderRepository:          orderRepository,
		productVariantRepository: productVariantRepository,
		rajaOngkirRepository:     rajaOngkirRepository,
	}
}
