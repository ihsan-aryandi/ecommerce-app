package provider

import (
	"ecommerce-app/internal/api/handler"
	"ecommerce-app/internal/api/repository"
	"ecommerce-app/internal/api/service"

	"github.com/google/wire"
)

var HandlersSet = wire.NewSet(
	handler.NewCalculateHandler,
	handler.NewCartHandler,
	handler.NewOrderHandler,
)

var ServicesSet = wire.NewSet(
	service.NewCalculateService,
	service.NewCartService,
	service.NewOrderService,
)

var RepositoriesSet = wire.NewSet(
	repository.NewRajaOngkirRepository,
	repository.NewProductVariantRepository,
	repository.NewCartItemRepository,
	repository.NewCartRepository,
	repository.NewOrderRepository,
)
