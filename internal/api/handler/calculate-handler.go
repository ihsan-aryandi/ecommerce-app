package handler

import (
	"ecommerce-app/internal/api/entity/request"
	"ecommerce-app/internal/api/service"

	"github.com/gin-gonic/gin"
)

type CalculateHandler struct {
	calculateService *service.CalculateService
}

func NewCalculateHandler(
	calculateService *service.CalculateService,
) *CalculateHandler {
	return &CalculateHandler{
		calculateService: calculateService,
	}
}

func (h CalculateHandler) CalculateSummaries(ctx *gin.Context) {
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
