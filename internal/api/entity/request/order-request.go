package request

import "ecommerce-app/internal/api/apierr"

type ProductVariant struct {
	ProductVariantId int64 `json:"product_variant_id"`
	Qty              int   `json:"qty"`
}

type CalculateSummaryRequest struct {
	Products              []ProductVariant `json:"products"`
	ShipperDestinationId  int64            `json:"shipper_destination_id"`
	ReceiverDestinationId int64            `json:"receiver_destination_id"`
	Courier               string           `json:"courier"`
	CourierType           string           `json:"courier_type"`
	CourierService        string           `json:"courier_service"`
	IsInitial             bool             `json:"is_initial"`
}

func (r CalculateSummaryRequest) ValidateCalculateSummary() *apierr.Error {
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

	if r.ShipperDestinationId <= 0 {
		validationErrors.Add("shipper_destination_id", apierr.EmptyFieldMessage())
	}

	if r.ReceiverDestinationId <= 0 {
		validationErrors.Add("receiver_destination_id", apierr.EmptyFieldMessage())
	}

	if !r.IsInitial {
		if r.Courier == "" {
			validationErrors.Add("courier", apierr.EmptyFieldMessage())
		}

		if r.CourierType == "" {
			validationErrors.Add("courier_type", apierr.EmptyFieldMessage())
		}

		if r.CourierService == "" {
			validationErrors.Add("courier_service", apierr.EmptyFieldMessage())
		}
	}

	return validationErrors.GetError()
}
