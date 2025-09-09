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

	"github.com/doug-martin/goqu/v9"
	"github.com/shopspring/decimal"
)

type OrderService struct {
	db                       *goqu.Database
	orderRepository          *repository.OrderRepository
	productVariantRepository *repository.ProductVariantRepository
	rajaOngkirRepository     *repository.RajaOngkirRepository
}

func NewOrderService(
	db *goqu.Database,
	orderRepository *repository.OrderRepository,
	productVariantRepository *repository.ProductVariantRepository,
	rajaOngkirRepository *repository.RajaOngkirRepository,
) *OrderService {
	return &OrderService{
		db:                       db,
		orderRepository:          orderRepository,
		productVariantRepository: productVariantRepository,
		rajaOngkirRepository:     rajaOngkirRepository,
	}
}

func (svc OrderService) CalculateSummaries(body *request.CalculateSummaryRequest) (*response.CalculateSummary, error) {
	productVariantMap, err := svc.getProductVariantMap(body)
	if err != nil {
		return nil, err
	}

	calcSummaryModel := &model.CalculateSummary{
		Variants:              productVariantMap,
		ShipperDestinationId:  body.ShipperDestinationId,
		ReceiverDestinationId: body.ReceiverDestinationId,
		Courier:               body.Courier,
		CourierType:           body.CourierType,
		CourierService:        body.CourierService,
	}

	summaryModel, err := svc.calculateSummary(body.IsInitial, calcSummaryModel, body.Products)
	if err != nil {
		return nil, err
	}

	return svc.toCalculateSummaryResponse(summaryModel), nil
}

func (svc OrderService) calculateSummary(isInitial bool, calcModel *model.CalculateSummary, productRequests []request.ProductVariant) (*model.SummaryModel, error) {
	// Calculate product summaries
	calculatedProducts, subTotal, weight, err := svc.calculateProductSummaries(calcModel, productRequests)
	if err != nil {
		return nil, err
	}

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

	shipment, shipmentType, err := svc.findEligibleShipment(shipping.Data, calcModel, isInitial)
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

func (svc OrderService) calculateProductSummaries(calcModel *model.CalculateSummary, productRequests []request.ProductVariant) (
	calculatedProducts []model.ProductSummary,
	subTotal decimal.Decimal,
	weight int,
	err error,
) {
	subTotal = decimal.Zero

	for _, productRequest := range productRequests {
		var (
			total = decimal.Zero
		)

		// Find product by product variant id
		variantModel, exists := calcModel.Variants[productRequest.ProductVariantId]
		if !exists {
			return nil, decimal.Zero, 0, apierr.IdNotFound("Product Variant ID", productRequest.ProductVariantId)
		}

		// Calculate total and weight
		qty := decimal.NewFromInt(int64(productRequest.Qty))
		total = total.Add(variantModel.Price.Decimal.Mul(qty))
		weightTotal := int(variantModel.Weight.Int32) * productRequest.Qty

		calculatedProducts = append(calculatedProducts, model.ProductSummary{
			ProductVariant: variantModel,
			Qty:            productRequest.Qty,
			Total:          total,
			Weight:         weightTotal,
		})

		// Sum subtotal and weight
		subTotal = subTotal.Add(total)
		weight += weightTotal
	}

	return
}

func (svc OrderService) findEligibleShipment(rajaOngkirCalculate *model.RajaOngkirCalculateResponse, calcModel *model.CalculateSummary, isInitial bool) (*model.RajaOngkirShipment, string, error) {
	// Default shipment
	if isInitial {
		if len(rajaOngkirCalculate.CalculateRegular) == 0 {
			return nil, "", apierr.InternalServer(errors.New("cannot find default shipment type"))
		}

		return &rajaOngkirCalculate.CalculateRegular[0], constants.CourierTypeRegular, nil
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

func (svc OrderService) findShipmentByCourier(shipments []model.RajaOngkirShipment, courierName, serviceName string) *model.RajaOngkirShipment {
	for _, shipment := range shipments {
		if strings.ToLower(shipment.ShippingName) == strings.ToLower(courierName) &&
			strings.ToLower(shipment.ServiceName) == strings.ToLower(serviceName) {
			return &shipment
		}
	}

	return nil
}

func (svc OrderService) getProductVariantMap(body *request.CalculateSummaryRequest) (model.ProductVariantMap, error) {
	ids := svc.getProductVariantIds(body)

	variants, err := svc.productVariantRepository.FindByIds(ids)
	if err != nil {
		return nil, apierr.InternalServer(err)
	}

	result := make(model.ProductVariantMap)
	for _, variant := range variants {
		result[variant.ID.Int64] = &variant
	}

	return result, nil
}

func (svc OrderService) getProductVariantIds(body *request.CalculateSummaryRequest) (result []int64) {
	for _, product := range body.Products {
		result = append(result, product.ProductVariantId)
	}

	return
}

func (svc OrderService) validateProductVariantList(products []model.ProductVariant, productVariantRequests []request.ProductVariant) *apierr.Error {
	for _, productVariant := range productVariantRequests {
		result := svc.findProductVariantById(products, productVariant.ProductVariantId)
		if result == nil {
			return apierr.IdNotFound("Product Variant ID", productVariant.ProductVariantId)
		}
	}

	return nil
}

func (svc OrderService) findProductVariantById(products []model.ProductVariant, id int64) *model.ProductVariant {
	for _, product := range products {
		if product.ID.Int64 == id {
			return &product
		}
	}

	return nil
}

func (svc OrderService) toCalculateSummaryResponse(summary *model.SummaryModel) *response.CalculateSummary {
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
		})
	}

	// -- Instant shipment
	for _, shipment := range summary.AvailableShipments.CalculateInstant {
		availableShipments.Instant = append(availableShipments.Instant, response.Shipment{
			CourierName:  shipment.ShippingName,
			ServiceName:  shipment.ServiceName,
			ShippingCost: shipment.ShippingCost,
		})
	}

	// -- Cargo shipment
	for _, shipment := range summary.AvailableShipments.CalculateCargo {
		availableShipments.Cargo = append(availableShipments.Regular, response.Shipment{
			CourierName:  shipment.ShippingName,
			ServiceName:  shipment.ServiceName,
			ShippingCost: shipment.ShippingCost,
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
