package route

import "ecommerce-app/internal/api/handler"

type HandlersContainer struct {
	OrderHandler *handler.OrderHandler
}

func NewHandlersContainer(
	orderHandler *handler.OrderHandler,
) *HandlersContainer {
	return &HandlersContainer{
		OrderHandler: orderHandler,
	}
}
