package model

import "database/sql"

type Audit struct {
	CreatedAt sql.NullTime
	CreatedBy sql.NullInt64
	UpdatedAt sql.NullTime
	UpdatedBy sql.NullInt64
	DeletedAt sql.NullTime
}
