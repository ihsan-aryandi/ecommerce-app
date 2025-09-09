package handler

import (
	"ecommerce-app/internal/api/apierr"
	"ecommerce-app/internal/api/entity/request"
	"ecommerce-app/internal/api/service"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService *service.CartService
}

func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

func (h *CartHandler) CreateCart(ctx *gin.Context) {
	body := new(request.CartRequest)

	if err := ctx.ShouldBindJSON(body); err != nil {
		ErrorJSON(ctx, apierr.InvalidRequest(err))
		return
	}

	if err := body.ValidateAddToCart(); err != nil {
		ErrorJSON(ctx, err)
		return
	}

	if err := h.cartService.CreateCart(body); err != nil {
		ErrorJSON(ctx, err)
		return
	}

	SuccessJSON(ctx, "Item has been added successfully", nil, nil)
}
