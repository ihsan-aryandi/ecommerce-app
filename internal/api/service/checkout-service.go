package service

import (
	"ecommerce-app/internal/api/apierr"
	"ecommerce-app/internal/api/entity/request"
	"ecommerce-app/internal/api/model"
	"ecommerce-app/internal/api/repository"
	"time"
)

type CheckoutService struct {
	checkoutSessionRepo   *repository.CheckoutSessionRepository
	productVariantService *ProductVariantService
}

func NewCheckoutService(
	checkoutSessionRepo *repository.CheckoutSessionRepository,
	productVariantService *ProductVariantService,
) *CheckoutService {
	return &CheckoutService{
		checkoutSessionRepo,
		productVariantService,
	}
}

func (svc CheckoutService) CreateCheckoutSession(body *request.CreateCheckoutSessionRequest) error {
	variantMap, err := svc.productVariantService.GetProductVariantMap(body.Products)
	if err != nil {
		return apierr.InternalServer(err)
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now()

	expirationDuration := 15 * time.Minute

	createdAt := now.In(loc)
	expiresAt := createdAt.Add(expirationDuration)

	createdAtUTC := now.UTC()
	expiresAtUTC := createdAtUTC.Add(expirationDuration)

	checkoutSession := &model.CheckoutSession{
		ProductVariants: variantMap,
		UserId:          1, // temp hardcoded
		CreatedAt:       now,
		CreatedAtUTC:    createdAtUTC,
		ExpiresAt:       expiresAt,
		ExpiresAtUTC:    expiresAtUTC,
	}

	_, err = svc.checkoutSessionRepo.Save(checkoutSession, expirationDuration)
	if err != nil {
		return apierr.InternalServer(err)
	}

	return nil
}
