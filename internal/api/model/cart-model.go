package model

import "database/sql"

type CartModel struct {
	ID     sql.NullInt64
	UserId sql.NullInt64
	Audit
}
