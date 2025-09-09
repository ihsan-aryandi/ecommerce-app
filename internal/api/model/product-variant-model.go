package model

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

type ProductVariant struct {
	ID          sql.NullInt64       `db:"id" goqu:"skipinsert,skipupdate"`
	ProductId   sql.NullInt64       `db:"product_id"`
	ProductName sql.NullString      `db:"product_name"`
	Price       decimal.NullDecimal `db:"price"`
	Stock       sql.NullInt32       `db:"stock"`
	Weight      sql.NullInt32       `db:"weight"`
	Audit
}
