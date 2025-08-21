package handler

import (
	"ecommerce-app/internal/api/apierr"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
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

	log.Printf("[INFO] Timestamp = %s | Method = %s | Host = %s | Path = %s | StatusCode = %d | Err = %s\n",
		time.Now().Format("2006-01-02 15-04-05.999999"),
		ctx.Request.Method,
		ctx.Request.Host,
		ctx.Request.RequestURI,
		statusCode,
		err)

	ctx.JSON(statusCode, response)
}
