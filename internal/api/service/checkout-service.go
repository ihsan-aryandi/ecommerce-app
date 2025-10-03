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

	checkoutSession := &model.CheckoutSession{
		ProductVariants: variantMap,
		UserId:          1, // temp hardcoded
	}

	_, err = svc.checkoutSessionRepo.Save(checkoutSession, time.Minute*15)
	if err != nil {
		return apierr.InternalServer(err)
	}

	return nil
}
