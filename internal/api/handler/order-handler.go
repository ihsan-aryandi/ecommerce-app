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

func (h OrderHandler) CalculateSummaries(ctx *gin.Context) {
	summaryRequest := new(request.CalculateSummaryRequest)

	if err := ctx.ShouldBindJSON(summaryRequest); err != nil {
		ErrorJSON(ctx, err)
		return
	}

	if err := summaryRequest.ValidateCalculateSummary(); err != nil {
		ErrorJSON(ctx, err)
		return
	}

	summaries, err := h.orderService.CalculateSummaries(summaryRequest)
	if err != nil {
		ErrorJSON(ctx, err)
		return
	}

	SuccessJSON(ctx, "Summaries calculated successfully", summaries, nil)
}
