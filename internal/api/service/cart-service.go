package service

import (
	"database/sql"
	"ecommerce-app/internal/api/apierr"
	"ecommerce-app/internal/api/entity/request"
	"ecommerce-app/internal/api/repository"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type CartService struct {
	db             *goqu.Database
	cartRepository *repository.CartRepository
}

func NewCartService(
	db *goqu.Database,
	cartRepository *repository.CartRepository,
) *CartService {
	return &CartService{
		db:             db,
		cartRepository: cartRepository,
	}
}

func (svc CartService) CreateCart(body *request.CartRequest) error {
	var (
		userId int64 = 1
		_            = time.Now()
	)

	err := svc.db.WithTx(func(tx *goqu.TxDatabase) error {
		hasCart := true

		// Find the user's cart
		_, err := svc.cartRepository.FindByUserIdTx(tx, userId)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return apierr.InternalServer(err)
		}

		// Check if the user has cart
		if errors.Is(err, sql.ErrNoRows) {
			hasCart = false
		}

		// Create cart if the user does not have one
		svc.createNewCart(tx, userId, hasCart)

		return nil
	})

	if !apierr.IsAPIError(err) {
		return apierr.InternalServer(err)
	}

	return err
}

func (svc CartService) createNewCart(tx *goqu.TxDatabase, userId int64, hasCart bool) error {
	if hasCart {
		return nil
	}

	return nil
}
