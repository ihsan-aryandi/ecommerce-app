package service

import (
	"ecommerce-app/internal/api/entity/request"
	"ecommerce-app/internal/api/entity/response"
	"ecommerce-app/internal/api/repository"

	"github.com/doug-martin/goqu/v9"
)

type OrderService struct {
	db                       *goqu.Database
	orderRepository          *repository.OrderRepository
	productVariantRepository *repository.ProductVariantRepository
	rajaOngkirRepository     *repository.RajaOngkirRepository
	productVariantService    *ProductVariantService
}

func NewOrderService(
	db *goqu.Database,
	orderRepository *repository.OrderRepository,
	productVariantRepository *repository.ProductVariantRepository,
	rajaOngkirRepository *repository.RajaOngkirRepository,
	productVariantService *ProductVariantService,
) *OrderService {
	return &OrderService{
		db:                       db,
		orderRepository:          orderRepository,
		productVariantRepository: productVariantRepository,
		rajaOngkirRepository:     rajaOngkirRepository,
		productVariantService:    productVariantService,
	}
}

func (svc OrderService) CreateOrder(body *request.OrderRequest) (*response.CreateOrderResponse, error) {
	//productVariantMap, err := svc.productVariantService.GetProductVariantMap(body.Products)
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil
}
