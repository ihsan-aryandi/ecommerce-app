package route

import "ecommerce-app/internal/api/handler"

type HandlersContainer struct {
	CartHandler *handler.CartHandler
	OrderHandler *handler.OrderHandler
}

func NewHandlersContainer(
	orderHandler *handler.OrderHandler,

	cartHandler *handler.CartHandler,) *HandlersContainer {
	return &HandlersContainer{
		CartHandler: cartHandler,
		OrderHandler: orderHandler,
	}
}
