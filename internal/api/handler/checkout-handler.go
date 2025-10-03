package handler

import (
	"ecommerce-app/internal/api/entity/request"
	"ecommerce-app/internal/api/service"

	"github.com/gin-gonic/gin"
)

type CheckoutHandler struct {
	calculateService *service.CalculateService
	checkoutService  *service.CheckoutService
}

func NewCheckoutHandler(
	calculateService *service.CalculateService,
	checkoutService *service.CheckoutService,
) *CheckoutHandler {
	return &CheckoutHandler{
		calculateService: calculateService,
		checkoutService:  checkoutService,
	}
}

func (h CheckoutHandler) Checkout(ctx *gin.Context) {
	body := new(request.CreateCheckoutSessionRequest)

	if err := ctx.ShouldBindJSON(body); err != nil {
		ErrorJSON(ctx, err)
		return
	}

	if err := body.Validate(); err != nil {
		ErrorJSON(ctx, err)
		return
	}

	err := h.checkoutService.CreateCheckoutSession(body)
	if err != nil {
		ErrorJSON(ctx, err)
		return
	}

	SuccessJSON(ctx, "Checkout success", nil, nil)
}

func (h CheckoutHandler) CalculateSummaries(ctx *gin.Context) {
	summaryRequest := new(request.CalculateSummaryRequest)

	if err := ctx.ShouldBindJSON(summaryRequest); err != nil {
		ErrorJSON(ctx, err)
		return
	}

	if err := summaryRequest.ValidateCalculateSummary(); err != nil {
		ErrorJSON(ctx, err)
		return
	}

	summaries, err := h.calculateService.CalculateSummaries(summaryRequest)
	if err != nil {
		ErrorJSON(ctx, err)
		return
	}

	SuccessJSON(ctx, "Summaries calculated successfully", summaries, nil)
}
