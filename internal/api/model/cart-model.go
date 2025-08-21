package model

import (
	"database/sql"
)

type CartModel struct {
	ID     sql.NullInt64 `db:"id"`
	UserId sql.NullInt64 `db:"user_id"`
	Audit
}
