package request

import "ecommerce-app/internal/api/apierr"

type CartRequest struct {
	ProductVariantId int64 `json:"product_variant_id"`
	Qty              int   `json:"qty"`
}

func (c CartRequest) ValidateAddToCart() *apierr.Error {
	validationError := make(apierr.ValidationError)

	if c.ProductVariantId <= 0 {
		validationError.Add("product_variant_id", apierr.EmptyFieldMessage())
	}

	if c.Qty <= 0 {
		validationError.Add("qty", apierr.EmptyFieldMessage())
	}

	return validationError.Error()
}
