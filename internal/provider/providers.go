package provider

import (
	"ecommerce-app/internal/api/handler"
	"ecommerce-app/internal/api/repository"
	"ecommerce-app/internal/api/service"

	"github.com/google/wire"
)

var HandlersSet = wire.NewSet(
	handler.NewOrderHandler,
)

var ServicesSet = wire.NewSet(
	service.NewOrderService,
)

var RepositoriesSet = wire.NewSet(
	repository.NewOrderRepository,
)
