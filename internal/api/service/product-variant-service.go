package service

import (
	"ecommerce-app/internal/api/apierr"
	"ecommerce-app/internal/api/entity/request"
	"ecommerce-app/internal/api/model"
	"ecommerce-app/internal/api/repository"
)

type ProductVariantService struct {
	productVariantRepository *repository.ProductVariantRepository
}

func NewProductVariantService(
	productVariantRepository *repository.ProductVariantRepository,
) *ProductVariantService {
	return &ProductVariantService{
		productVariantRepository,
	}
}

func (svc ProductVariantService) GetProductVariantMap(products []request.ProductVariant) (model.ProductVariantMap, error) {
	ids := svc.getProductVariantIds(products)

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

func (svc ProductVariantService) getProductVariantIds(productVariants []request.ProductVariant) (result []int64) {
	for _, variant := range productVariants {
		result = append(result, variant.ProductVariantId)
	}

	return
}
