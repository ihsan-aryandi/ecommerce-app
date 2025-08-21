package model

import (
	"database/sql"
)

type Audit struct {
	CreatedAt sql.NullTime  `db:"created_at" goqu:"skipupdate"`
	CreatedBy sql.NullInt64 `db:"created_by" goqu:"skipupdate"`
	UpdatedAt sql.NullTime  `db:"updated_at" goqu:"skipinsert"`
	UpdatedBy sql.NullInt64 `db:"updated_by" goqu:"skipinsert"`
	DeletedAt sql.NullTime  `db:"deleted_at" goqu:"skipinsert"`
}
