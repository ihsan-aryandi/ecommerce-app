package response

import "github.com/shopspring/decimal"

type CreateOrderResponse struct {
	OrderId     int64           `json:"order_id"`
	PaymentType string          `json:"payment_type"`
	VANumber    string          `json:"va_number"`
	QrisUrl     string          `json:"qris_url"`
	GrandTotal  decimal.Decimal `json:"grand_total"`
}
