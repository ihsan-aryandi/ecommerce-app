package service

import "ecommerce-app/internal/api/repository"

type CheckoutService struct {
	checkoutSessionRepo *repository.CheckoutSessionRepository
}

func NewCheckoutService(
	checkoutSessionRepo *repository.CheckoutSessionRepository,
) *CheckoutService {
	return &CheckoutService{
		checkoutSessionRepo,
	}
}

func (svc CheckoutService) CreateCheckoutSession() {

}
