package service

import (
	"ecommerce-app/internal/api/apierr"
	"ecommerce-app/internal/api/entity/request"
	"ecommerce-app/internal/api/model"
	"ecommerce-app/internal/api/repository"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type CartService struct {
	db                       *goqu.Database
	cartRepository           *repository.CartRepository
	cartItemRepository       *repository.CartItemRepository
	productVariantRepository *repository.ProductVariantRepository
}

func NewCartService(
	db *goqu.Database,
	cartRepository *repository.CartRepository,
	cartItemRepository *repository.CartItemRepository,
	productVariantRepository *repository.ProductVariantRepository,
) *CartService {
	return &CartService{
		db:                       db,
		cartRepository:           cartRepository,
		cartItemRepository:       cartItemRepository,
		productVariantRepository: productVariantRepository,
	}
}

func (svc CartService) CreateCart(body *request.CartRequest) error {
	var (
		userId int64 = 1
		now          = time.Now()
	)

	return svc.db.WithTx(func(tx *goqu.TxDatabase) error {
		var (
			cartId int64
		)

		// Validate product variant
		if err := svc.validateProductVariant(tx, body.ProductVariantId); err != nil {
			return err
		}

		// Find the user's cart
		cart, err := svc.cartRepository.FindByUserIdTx(tx, userId)
		if err != nil {
			return apierr.InternalServer(err)
		}

		// Create cart if the user does not have one
		if cart == nil {
			cartId, err = svc.createNewCart(tx, userId, now)
			if err != nil {
				return err
			}
		} else {
			cartId = cart.ID.Int64
		}

		// Add items
		if err = svc.addCartItem(tx, cartId, body, now); err != nil {
			return err
		}

		return nil
	})
}

func (svc CartService) createNewCart(tx *goqu.TxDatabase, userId int64, now time.Time) (cartId int64, err error) {
	cartId, err = svc.cartRepository.InsertTx(tx, &model.CartModel{
		UserId: NullInt64(userId),
		Audit:  AuditInsert(now, userId),
	})
	if err != nil {
		return 0, apierr.InternalServer(err)
	}

	return cartId, nil
}

func (svc CartService) addCartItem(tx *goqu.TxDatabase, cartId int64, body *request.CartRequest, now time.Time) (err error) {
	cartItem, err := svc.cartItemRepository.FindByCartIdAndProductVariantIdTx(tx, cartId, body.ProductVariantId)
	if err != nil {
		return err
	}

	// User already has the item in cart. Increase quantity
	if cartItem != nil {
		if err = svc.increaseQuantity(tx, cartItem, body.Qty, now); err != nil {
			return err
		}

		return nil
	}

	// User does not have the item
	if err = svc.createCartItem(tx, cartId, body, now); err != nil {
		return err
	}

	return nil
}

func (svc CartService) createCartItem(tx *goqu.TxDatabase, cartId int64, body *request.CartRequest, now time.Time) (err error) {
	_, err = svc.cartItemRepository.InsertTx(tx, &model.CartItemModel{
		CartId:           NullInt64(cartId),
		ProductVariantId: NullInt64(body.ProductVariantId),
		Qty:              NullInt32(body.Qty),
		Audit:            AuditInsert(now, 1),
	})
	if err != nil {
		return apierr.InternalServer(err)
	}

	return nil
}

func (svc CartService) increaseQuantity(tx *goqu.TxDatabase, cartItem *model.CartItemModel, qty int32, now time.Time) (err error) {
	err = svc.cartItemRepository.UpdateQuantityByIdTx(tx, &model.CartItemModel{
		ID:    NullInt64(cartItem.ID.Int64),
		Qty:   NullInt32(cartItem.Qty.Int32 + qty),
		Audit: AuditUpdate(now, 1),
	})
	if err != nil {
		return apierr.InternalServer(err)
	}

	return nil
}

func (svc CartService) validateProductVariant(tx *goqu.TxDatabase, productVariantId int64) error {
	productVariant, err := svc.productVariantRepository.FindByIdTx(tx, productVariantId)
	if err != nil {
		return apierr.InternalServer(err)
	}

	if productVariant == nil {
		return apierr.DataNotFound("Product Variant")
	}

	return nil
}
