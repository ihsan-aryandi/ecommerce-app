package service

import (
	"ecommerce-app/internal/api/apierr"
	"ecommerce-app/internal/api/entity/request"
	"ecommerce-app/internal/api/entity/response"
	"ecommerce-app/internal/api/model"
	"ecommerce-app/internal/api/repository"

	"github.com/doug-martin/goqu/v9"
)

type OrderService struct {
	db                        *goqu.Database
	orderRepository           *repository.OrderRepository
	productVariantRepository  *repository.ProductVariantRepository
	rajaOngkirRepository      *repository.RajaOngkirRepository
	productVariantService     *ProductVariantService
	checkoutSessionRepository *repository.CheckoutSessionRepository
	calculateService          *CalculateService
}

func NewOrderService(
	db *goqu.Database,
	orderRepository *repository.OrderRepository,
	productVariantRepository *repository.ProductVariantRepository,
	rajaOngkirRepository *repository.RajaOngkirRepository,
	productVariantService *ProductVariantService,
	checkoutSessionRepository *repository.CheckoutSessionRepository,
	calculateService *CalculateService,
) *OrderService {
	return &OrderService{
		db:                        db,
		orderRepository:           orderRepository,
		productVariantRepository:  productVariantRepository,
		rajaOngkirRepository:      rajaOngkirRepository,
		productVariantService:     productVariantService,
		checkoutSessionRepository: checkoutSessionRepository,
		calculateService:          calculateService,
	}
}

func (svc OrderService) CreateOrder(body *request.OrderRequest) (*response.CreateOrderResponse, error) {
	// Retrieve checkout session by user id
	checkoutSession, err := svc.checkoutSessionRepository.GetByUserId(1)
	if err != nil {
		return nil, err
	}

	// Check if the user has checkout session
	if checkoutSession == nil {
		return nil, apierr.DataNotFound("Checkout session")
	}

	// Convert variant sessions to variant requests
	productVariantsReq := svc.variantSessionsToVariantRequests(checkoutSession.ProductVariants)

	// Fresh products from database
	productVariantMap, err := svc.productVariantService.GetProductVariantMap(productVariantsReq)
	if err != nil {
		return nil, err
	}

	// Validate prices & weights diffs
	if err = svc.compareProductPricesAndWeights(checkoutSession.ProductVariants, productVariantMap); err != nil {
		return nil, err
	}

	// Recalculate summary
	summary, err := svc.recalculateSummary(checkoutSession, productVariantMap)
	if err != nil {
		return nil, err
	}

	// Validate shipping cost diff
	if !checkoutSession.ShippingCost.Equal(summary.ShippingCost) {
		return nil, apierr.ShippingCostChanged(summary.Courier, summary.CourierService, checkoutSession.ShippingCost, summary.ShippingCost)
	}

	return nil, nil
}

func (svc OrderService) variantSessionsToVariantRequests(variantSessions model.ProductVariantMap) (result []request.ProductVariant) {
	for _, productVariant := range variantSessions {
		result = append(result, request.ProductVariant{
			ProductVariantId: productVariant.ID.Int64,
			Qty:              int(productVariant.Qty.Int32),
		})
	}

	return
}

func (svc OrderService) recalculateSummary(checkoutSession *model.CheckoutSession, productVariantMap model.ProductVariantMap) (*model.SummaryModel, error) {
	calcSummaryModel := &model.CalculateSummary{
		IsInitial:             false,
		Variants:              productVariantMap,
		ShipperDestinationId:  checkoutSession.ShipperDestinationId,
		ReceiverDestinationId: checkoutSession.ReceiverDestinationId,
		Courier:               checkoutSession.Courier,
		CourierType:           checkoutSession.CourierType,
		CourierService:        checkoutSession.CourierService,
	}

	return svc.calculateService.CalculateSummary(calcSummaryModel)
}

func (svc OrderService) compareProductPricesAndWeights(variantSessions model.ProductVariantMap, variants model.ProductVariantMap) error {
	for id, variantSession := range variantSessions {
		variant, ok := variants[id]
		if !ok {
			continue
		}

		if !variantSession.Price.Decimal.Equal(variant.Price.Decimal) {
			return apierr.ProductPriceChanged(variant.ProductName.String, variantSession.Price.Decimal, variant.Price.Decimal)
		}

		if variantSession.Weight.Int32 != variant.Weight.Int32 {
			return apierr.ProductWeightChanged(variant.ProductName.String, int(variantSession.Weight.Int32), int(variant.Weight.Int32))
		}
	}

	return nil
}
