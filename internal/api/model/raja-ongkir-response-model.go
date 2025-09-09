package model

import "github.com/shopspring/decimal"

type RajaOngkirResponse[T any] struct {
	Meta RajaOngkirMeta `json:"meta"`
	Data T              `json:"data"`
}

type RajaOngkirMeta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type RajaOngkirCalculateResponse struct {
	CalculateRegular []RajaOngkirShipment `json:"calculate_reguler"`
	CalculateCargo   []RajaOngkirShipment `json:"calculate_cargo"`
	CalculateInstant []RajaOngkirShipment `json:"calculate_instant"`
}

type RajaOngkirShipment struct {
	ShippingName     string          `json:"shipping_name"`
	ServiceName      string          `json:"service_name"`
	Weight           float64         `json:"weight"` // In kilogram
	IsCod            bool            `json:"is_cod"`
	ShippingCost     decimal.Decimal `json:"shipping_cost"`
	ShippingCashback decimal.Decimal `json:"shipping_cashback"`
	ShippingCostNet  decimal.Decimal `json:"shipping_cost_net"`
	GrandTotal       decimal.Decimal `json:"grandtotal"`
	ServiceFee       decimal.Decimal `json:"service_fee"`
	NetIncome        decimal.Decimal `json:"net_income"`
	Etd              string          `json:"etd"`
}
