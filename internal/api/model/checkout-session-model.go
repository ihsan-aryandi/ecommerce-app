package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type CheckoutSession struct {
	CheckoutID      string            `json:"checkout_id"`
	ProductVariants ProductVariantMap `json:"product_variants"`
	UserId          int64             `json:"user_id"`
	SubTotal        decimal.Decimal   `json:"sub_total"`
	ShippingCost    decimal.Decimal   `json:"shipping_cost"`
	Total           decimal.Decimal   `json:"total"`
	CreatedAt       time.Time         `json:"created_at"`
	CreatedAtUTC    time.Time         `json:"created_at_utc"`
	ExpiresAt       time.Time         `json:"expires_at"`
	ExpiresAtUTC    time.Time         `json:"expires_at_utc"`
}
