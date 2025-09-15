package handler

import (
	"ecommerce-app/internal/api/entity/request"
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

func (h OrderHandler) CreateOrder(ctx *gin.Context) {
	body := new(request.OrderRequest)

	if err := ctx.ShouldBindJSON(body); err != nil {
		ErrorJSON(ctx, err)
		return
	}

	if err := body.ValidateCreateOrder(); err != nil {
		ErrorJSON(ctx, err)
		return
	}

	response, err := h.orderService.CreateOrder(body)
	if err != nil {
		ErrorJSON(ctx, err)
		return
	}

	SuccessJSON(ctx, "Success", response, nil)
}
