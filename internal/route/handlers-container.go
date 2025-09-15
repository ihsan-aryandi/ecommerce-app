package route

import "ecommerce-app/internal/api/handler"

type HandlersContainer struct {
	CheckoutHandler *handler.CheckoutHandler
	CartHandler     *handler.CartHandler
	OrderHandler    *handler.OrderHandler
}

func NewHandlersContainer(
	orderHandler *handler.OrderHandler,
	cartHandler *handler.CartHandler,
	checkoutHandler *handler.CheckoutHandler,
) *HandlersContainer {
	return &HandlersContainer{
		CheckoutHandler: checkoutHandler,
		CartHandler:     cartHandler,
		OrderHandler:    orderHandler,
	}
}
