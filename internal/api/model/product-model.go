package model

import (
	"database/sql"
)

type Product struct {
	ID          sql.NullInt64  `db:"id" goqu:"skipinsert,skipupdate"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	Audit
}
