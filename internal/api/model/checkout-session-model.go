package model

type CheckoutSession struct {
	CheckoutID      string            `json:"checkout_id"`
	ProductVariants ProductVariantMap `json:"product_variants"`
	UserId          int64             `json:"user_id"`
}
