package response

import "github.com/shopspring/decimal"

type CalculateSummary struct {
	AvailableShipments AvailableShipments `json:"available_shipments"`
	Products           []ProductSummary   `json:"products"`
	WeightTotal        int                `json:"weight_total"`
	Courier            string             `json:"courier"`
	CourierType        string             `json:"courier_type"`
	CourierService     string             `json:"courier_service"`
	SubTotal           decimal.Decimal    `json:"sub_total"`
	ShippingCost       decimal.Decimal    `json:"shipping_cost"`
	Total              decimal.Decimal    `json:"total"`
}

type ProductSummary struct {
	ProductVariantId int64           `json:"product_variant_id"`
	Name             string          `json:"name"`
	Price            decimal.Decimal `json:"price"`
	Qty              int             `json:"qty"`
	Total            decimal.Decimal `json:"total"`
	Weight           int             `json:"weight"`
}

type AvailableShipments struct {
	Regular []Shipment `json:"regular"`
	Instant []Shipment `json:"instant"`
	Cargo   []Shipment `json:"cargo"`
}

type Shipment struct {
	CourierName  string          `json:"shipping_name"`
	ServiceName  string          `json:"service_name"`
	ShippingCost decimal.Decimal `json:"shipping_cost"`
	Etd          string          `json:"etd"`
}
