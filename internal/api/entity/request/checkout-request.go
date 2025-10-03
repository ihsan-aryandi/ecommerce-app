package request

import "ecommerce-app/internal/api/apierr"

type CreateCheckoutSessionRequest struct {
	Products []ProductVariant `json:"products"`
}

func (r CreateCheckoutSessionRequest) Validate() *apierr.Error {
	validationErrors := apierr.NewValidationError()

	if len(r.Products) == 0 {
		validationErrors.Add("products", apierr.EmptyFieldMessage())
	}

	for _, product := range r.Products {
		if product.ProductVariantId <= 0 {
			return apierr.EmptyField("product_variant_id")
		}

		if product.Qty <= 0 {
			return apierr.EmptyField("qty")
		}
	}

	return validationErrors.GetError()
}
