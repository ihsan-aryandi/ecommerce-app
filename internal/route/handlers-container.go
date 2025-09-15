package route

import "ecommerce-app/internal/api/handler"

type HandlersContainer struct {
	CheckoutHandler *handler.CheckoutHandler
	CalculateHandler *handler.CalculateHandler
	CartHandler *handler.CartHandler
	OrderHandler *handler.OrderHandler
}

func NewHandlersContainer(
	orderHandler *handler.OrderHandler,

	cartHandler *handler.CartHandler,
	calculateHandler *handler.CalculateHandler,
	checkoutHandler *handler.CheckoutHandler,) *HandlersContainer {
	return &HandlersContainer{
		CheckoutHandler: checkoutHandler,
		CalculateHandler: calculateHandler,
		CartHandler: cartHandler,
		OrderHandler: orderHandler,
	}
}
