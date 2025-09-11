package route

import "ecommerce-app/internal/api/handler"

type HandlersContainer struct {
	CalculateHandler *handler.CalculateHandler
	CartHandler *handler.CartHandler
	OrderHandler *handler.OrderHandler
}

func NewHandlersContainer(
	orderHandler *handler.OrderHandler,

	cartHandler *handler.CartHandler,
	calculateHandler *handler.CalculateHandler,) *HandlersContainer {
	return &HandlersContainer{
		CalculateHandler: calculateHandler,
		CartHandler: cartHandler,
		OrderHandler: orderHandler,
	}
}
