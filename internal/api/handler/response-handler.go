package handler

import (
	"ecommerce-app/internal/api/apierr"
	"errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool
	Code    string
	Message string
	Details interface{}
}

func SuccessJSON(ctx *gin.Context, message string) {
	ctx.JSON(200, Response{
		Success: true,
		Code:    "SUCCESS",
		Message: message,
	})
}

func ErrorJSON(ctx *gin.Context, err error) {
	statusCode := 500
	response := Response{
		Success: false,
		Code:    "FAILED",
		Message: err.Error(),
	}

	apiErr := &apierr.Error{}
	if errors.As(err, &apiErr) {
		statusCode = apiErr.StatusCode
		response.Code = apiErr.Code
		response.Message = apiErr.Message
		response.Details = apiErr.Details
	}

	ctx.JSON(statusCode, response)
}
