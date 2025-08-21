package service

import (
	"database/sql"
	"ecommerce-app/internal/api/model"
	"time"
)

func AuditInsert(createdAt time.Time, createdBy int64) model.Audit {
	return model.Audit{
		CreatedAt: NullTime(createdAt),
		CreatedBy: NullInt64(createdBy),
	}
}

func AuditUpdate(updatedAt time.Time, updatedBy int64) model.Audit {
	return model.Audit{
		UpdatedAt: NullTime(updatedAt),
		UpdatedBy: NullInt64(updatedBy),
	}
}

func AuditDelete(deletedAt time.Time, deletedBy int64) model.Audit {
	return model.Audit{
		DeletedAt: NullTime(deletedAt),
		UpdatedBy: NullInt64(deletedBy),
	}
}

func NullString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  value != "",
	}
}

func NullInt64(value int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: value,
		Valid: value != 0,
	}
}

func NullInt32(value int32) sql.NullInt32 {
	return sql.NullInt32{
		Int32: value,
		Valid: value != 0,
	}
}

func NullTime(value time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  value,
		Valid: !value.IsZero(),
	}
}
