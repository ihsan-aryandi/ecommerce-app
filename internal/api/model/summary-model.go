package model

import "github.com/shopspring/decimal"

type SummaryModel struct {
	AvailableShipments *RajaOngkirCalculateResponse
	Products           []ProductSummary
	Weight             int
	Courier            string
	CourierType        string
	CourierService     string
	SubTotal           decimal.Decimal
	ShippingCost       decimal.Decimal
	Total              decimal.Decimal
}

type ProductSummary struct {
	ProductVariant *ProductVariant
	Qty            int
	Total          decimal.Decimal
	Weight         int
}

type CalculateSummary struct {
	Variants              ProductVariantMap
	ShipperDestinationId  int64
	ReceiverDestinationId int64
	Courier               string
	CourierType           string
	CourierService        string
}

type ProductVariantMap map[int64]*ProductVariant
