package handler

import (
	"ecommerce-app/internal/api/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h OrderHandler) GetTotal(ctx *gin.Context) {
	totals := h.orderService.GetTotal()
	ctx.JSON(200, totals)
}
