package service

import (
	"database/sql"
	"ecommerce-app/internal/api/apierr"
	"errors"
)

func NoRowsErr(entity string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return apierr.DataNotFound(entity, err)
	}

	return err
}
