package service

import (
	"ecommerce-app/internal/api/apierr"
	"ecommerce-app/internal/api/entity/request"
	"ecommerce-app/internal/api/entity/response"
	"ecommerce-app/internal/api/model"
	"ecommerce-app/internal/api/repository"
	"ecommerce-app/internal/constants"
	"errors"
	"strings"

	"github.com/shopspring/decimal"
)

type CalculateService struct {
	productVariantRepository  *repository.ProductVariantRepository
	rajaOngkirRepository      *repository.RajaOngkirRepository
	checkoutSessionRepository *repository.CheckoutSessionRepository
	productVariantService     *ProductVariantService
}

func NewCalculateService(
	productVariantRepository *repository.ProductVariantRepository,
	rajaOngkirRepository *repository.RajaOngkirRepository,
	checkoutSessionRepository *repository.CheckoutSessionRepository,
	productVariantService *ProductVariantService,
) *CalculateService {
	return &CalculateService{
		productVariantRepository:  productVariantRepository,
		rajaOngkirRepository:      rajaOngkirRepository,
		checkoutSessionRepository: checkoutSessionRepository,
		productVariantService:     productVariantService,
	}
}

func (svc CalculateService) CalculateSummaries(body *request.CalculateSummaryRequest) (*response.CalculateSummary, error) {
	checkoutSession, err := svc.checkoutSessionRepository.GetByUserId(1)
	if err != nil {
		return nil, err
	}

	if checkoutSession == nil {
		return nil, apierr.DataNotFound("Checkout session")
	}

	calcSummaryModel := &model.CalculateSummary{
		IsInitial:             body.IsInitial,
		Variants:              checkoutSession.ProductVariants,
		ShipperDestinationId:  body.ShipperDestinationId,
		ReceiverDestinationId: body.ReceiverDestinationId,
		Courier:               body.Courier,
		CourierType:           body.CourierType,
		CourierService:        body.CourierService,
	}

	summaryModel, err := svc.CalculateSummary(calcSummaryModel)
	if err != nil {
		return nil, err
	}

	return svc.toCalculateSummaryResponse(summaryModel), nil
}

func (svc CalculateService) CalculateSummary(calcModel *model.CalculateSummary) (*model.SummaryModel, error) {
	// Calculate product summaries
	calculatedProducts, subTotal, weight := svc.calculateProductSummaries(calcModel)

	// Fetch Raja Ongkir API to calculate shipping cost
	shipping, err := svc.rajaOngkirRepository.CalculateShippingCost(&model.RajaOngkirCalculateShippingCost{
		ShipperDestinationId:  calcModel.ShipperDestinationId,
		ReceiverDestinationId: calcModel.ReceiverDestinationId,
		Weight:                float64(weight) / 1000, // convert to kilogram
		ItemValue:             subTotal,
	})
	if err != nil {
		return nil, apierr.InternalServer(err)
	}

	shipment, shipmentType, err := svc.findEligibleShipment(shipping.Data, calcModel, calcModel.IsInitial)
	if err != nil {
		return nil, err
	}

	result := &model.SummaryModel{
		AvailableShipments: shipping.Data,
		Products:           calculatedProducts,
		Weight:             weight,
		Courier:            shipment.ShippingName,
		CourierType:        shipmentType,
		CourierService:     shipment.ServiceName,
		SubTotal:           subTotal,
		ShippingCost:       shipment.ShippingCost,
		Total:              subTotal.Add(shipment.ShippingCost),
	}

	return result, nil
}

func (svc CalculateService) calculateProductSummaries(calcModel *model.CalculateSummary) (
	calculatedProducts []model.ProductSummary,
	subTotal decimal.Decimal,
	weight int,
) {
	subTotal = decimal.Zero

	for _, variant := range calcModel.Variants {
		var (
			total = decimal.Zero
		)

		// Calculate total and weight
		qty := decimal.NewFromInt32(variant.Qty.Int32)
		total = total.Add(variant.Price.Decimal.Mul(qty))
		weightTotal := int(variant.Weight.Int32) * int(variant.Qty.Int32)

		calculatedProducts = append(calculatedProducts, model.ProductSummary{
			ProductVariant: variant,
			Qty:            int(variant.Qty.Int32),
			Total:          total,
			Weight:         weightTotal,
		})

		// Sum subtotal and weight
		subTotal = subTotal.Add(total)
		weight += weightTotal
	}

	return
}

func (svc CalculateService) findEligibleShipment(rajaOngkirCalculate *model.RajaOngkirCalculateResponse, calcModel *model.CalculateSummary, isInitial bool) (*model.RajaOngkirShipment, string, error) {
	// Default shipment
	if isInitial {
		if len(rajaOngkirCalculate.CalculateRegular) == 0 {
			return nil, "", apierr.InternalServer(errors.New("cannot find default shipment type"))
		}

		return svc.findCheapestShipment(rajaOngkirCalculate.CalculateRegular), constants.CourierTypeRegular, nil
	}

	var (
		shipments    []model.RajaOngkirShipment
		shipmentType string
		courierType  = strings.ToLower(calcModel.CourierType)
	)

	// Regular shipments
	if courierType == constants.CourierTypeRegular {
		shipments = rajaOngkirCalculate.CalculateRegular
		shipmentType = constants.CourierTypeRegular
	}

	// Cargo shipments
	if courierType == constants.CourierTypeCargo {
		shipments = rajaOngkirCalculate.CalculateCargo
		shipmentType = constants.CourierTypeCargo
	}

	// Instant shipments
	if courierType == constants.CourierTypeInstant {
		shipments = rajaOngkirCalculate.CalculateInstant
		shipmentType = constants.CourierTypeInstant
	}

	shipment := svc.findShipmentByCourier(shipments, calcModel.Courier, calcModel.CourierService)
	if shipment == nil {
		return nil, "", apierr.DataNotFound("Courier")
	}

	return shipment, shipmentType, nil
}

func (svc CalculateService) findCheapestShipment(shipments []model.RajaOngkirShipment) *model.RajaOngkirShipment {
	if len(shipments) == 0 {
		return nil
	}

	result := shipments[0]
	for _, shipment := range shipments {
		if shipment.ShippingCost.LessThan(result.ShippingCost) {
			result = shipment
		}
	}

	return &result
}

func (svc CalculateService) findShipmentByCourier(shipments []model.RajaOngkirShipment, courierName, serviceName string) *model.RajaOngkirShipment {
	for _, shipment := range shipments {
		if strings.ToLower(shipment.ShippingName) == strings.ToLower(courierName) &&
			strings.ToLower(shipment.ServiceName) == strings.ToLower(serviceName) {
			return &shipment
		}
	}

	return nil
}

func (svc CalculateService) validateProductVariantList(products []model.ProductVariant, productVariantRequests []request.ProductVariant) *apierr.Error {
	for _, productVariant := range productVariantRequests {
		result := svc.findProductVariantById(products, productVariant.ProductVariantId)
		if result == nil {
			return apierr.IdNotFound("Product Variant ID", productVariant.ProductVariantId)
		}
	}

	return nil
}

func (svc CalculateService) findProductVariantById(products []model.ProductVariant, id int64) *model.ProductVariant {
	for _, product := range products {
		if product.ID.Int64 == id {
			return &product
		}
	}

	return nil
}

func (svc CalculateService) toCalculateSummaryResponse(summary *model.SummaryModel) *response.CalculateSummary {
	var (
		productSummaries   []response.ProductSummary
		availableShipments response.AvailableShipments
	)

	// Convert products
	for _, product := range summary.Products {
		productSummaries = append(productSummaries, response.ProductSummary{
			ProductVariantId: product.ProductVariant.ID.Int64,
			Name:             product.ProductVariant.ProductName.String,
			Price:            product.ProductVariant.Price.Decimal,
			Qty:              product.Qty,
			Total:            product.Total,
			Weight:           product.Weight,
		})
	}

	// Convert shipments
	// -- Regular shipment
	for _, shipment := range summary.AvailableShipments.CalculateRegular {
		availableShipments.Regular = append(availableShipments.Regular, response.Shipment{
			CourierName:  shipment.ShippingName,
			ServiceName:  shipment.ServiceName,
			ShippingCost: shipment.ShippingCost,
			Etd:          shipment.Etd,
		})
	}

	// -- Instant shipment
	for _, shipment := range summary.AvailableShipments.CalculateInstant {
		availableShipments.Instant = append(availableShipments.Instant, response.Shipment{
			CourierName:  shipment.ShippingName,
			ServiceName:  shipment.ServiceName,
			ShippingCost: shipment.ShippingCost,
			Etd:          shipment.Etd,
		})
	}

	// -- Cargo shipment
	for _, shipment := range summary.AvailableShipments.CalculateCargo {
		availableShipments.Cargo = append(availableShipments.Regular, response.Shipment{
			CourierName:  shipment.ShippingName,
			ServiceName:  shipment.ServiceName,
			ShippingCost: shipment.ShippingCost,
			Etd:          shipment.Etd,
		})
	}

	result := &response.CalculateSummary{
		AvailableShipments: availableShipments,
		Products:           productSummaries,
		WeightTotal:        summary.Weight,
		Courier:            summary.Courier,
		CourierType:        summary.CourierType,
		CourierService:     summary.CourierService,
		SubTotal:           summary.SubTotal,
		ShippingCost:       summary.ShippingCost,
		Total:              summary.Total,
	}

	return result
}
