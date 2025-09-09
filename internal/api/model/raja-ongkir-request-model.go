package model

import "github.com/shopspring/decimal"

type RajaOngkirCalculateShippingCost struct {
	ShipperDestinationId  int64
	ReceiverDestinationId int64
	Weight                float64 // Kilogram
	ItemValue             decimal.Decimal
}
