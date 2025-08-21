package model

import "database/sql"

type CartItemModel struct {
	ID               sql.NullInt64 `db:"id" goqu:"skipinsert,skipupdate"`
	CartId           sql.NullInt64 `db:"cart_id"`
	ProductVariantId sql.NullInt64 `db:"product_variant_id"`
	Qty              sql.NullInt32 `db:"qty"`
	Audit
}
