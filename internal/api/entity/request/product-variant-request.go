package request

type ProductVariant struct {
	ProductVariantId int64 `json:"product_variant_id"`
	Qty              int   `json:"qty"`
}
