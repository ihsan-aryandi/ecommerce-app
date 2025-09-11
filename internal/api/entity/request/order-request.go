package request

type ProductVariant struct {
	ProductVariantId int64 `json:"product_variant_id"`
	Qty              int   `json:"qty"`
}

type OrderRequest struct {
	Products              []ProductVariant `json:"products"`
	ShipperDestinationId  int64            `json:"shipper_destination_id"`
	ReceiverDestinationId int64            `json:"receiver_destination_id"`
	Courier               string           `json:"courier"`
	CourierType           string           `json:"courier_type"`
	CourierService        string           `json:"courier_service"`
	PaymentType           string           `json:"payment_type"`
}
