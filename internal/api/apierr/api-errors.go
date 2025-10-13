package apierr

import (
	"ecommerce-app/internal/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type Params map[string]string

type Error struct {
	StatusCode int
	Code       string
	Message    string
	CausedBy   error
	Details    interface{}
}

func (err *Error) Error() string {
	if err.CausedBy != nil {
		return err.CausedBy.Error()
	}

	return ""
}

func InternalServer(causedBy error) *Error {
	return &Error{
		StatusCode: 500,
		Code:       InternalServerErrorCode,
		Message:    InternalServerErrorMessage,
		CausedBy:   causedBy,
	}
}

func DataNotFound(entity string) *Error {
	return &Error{
		StatusCode: 400,
		Code:       DataNotFoundErrorCode,
		Message: replacePlaceholders(DataNotFoundErrorMessage, Params{
			"entity": entity,
		}),
		CausedBy: errors.New("data not found"),
	}
}

func IdNotFound(field string, id int64) *Error {
	strId := strconv.Itoa(int(id))
	return &Error{
		StatusCode: 400,
		Code:       IdNotFoundErrorCode,
		Message: replacePlaceholders(IdNotFoundErrorMessage, Params{
			"field": field,
			"id":    strId,
		}),
		CausedBy: errors.New("id not found"),
		Details: map[string]interface{}{
			"field": field,
			"id":    id,
		},
	}
}

func EmptyField(entity string) *Error {
	return &Error{
		StatusCode: 400,
		Code:       EmptyFieldErrorCode,
		Message: replacePlaceholders(EmptyFieldErrorMessage, Params{
			"field": entity,
		}),
		CausedBy: errors.New("data not found"),
	}
}

func InvalidRequest(causedBy error) *Error {
	return &Error{
		StatusCode: 500,
		Code:       InvalidRequestErrorCode,
		Message:    InvalidRequestErrorMessage,
		CausedBy:   causedBy,
	}
}

func ProductPriceChanged(productName string, oldPrice decimal.Decimal, newPrice decimal.Decimal) *Error {
	oldPriceFloat, _ := oldPrice.Float64()
	newPriceFloat, _ := newPrice.Float64()

	message := replacePlaceholders(ProductPriceChangedMessage, Params{
		"productName": productName,
		"oldPrice":    utils.FormatRupiah(oldPriceFloat),
		"newPrice":    utils.FormatRupiah(newPriceFloat),
	})

	return &Error{
		StatusCode: 409,
		Code:       ProductPriceChangedErrorCode,
		Message:    message,
		CausedBy:   errors.New(message),
		Details: map[string]interface{}{
			"product_name": productName,
			"old_price":    oldPrice,
			"new_price":    newPrice,
		},
	}
}

func ProductWeightChanged(productName string, oldWeight, newWeight int) *Error {
	message := replacePlaceholders(ProductWeightChangedMessage, Params{
		"productName": productName,
		"oldWeight":   strconv.Itoa(oldWeight),
		"newWeight":   strconv.Itoa(newWeight),
	})

	return &Error{
		StatusCode: 409,
		Code:       ProductWeightChangedErrorCode,
		Message:    message,
		CausedBy:   errors.New(message),
		Details: map[string]interface{}{
			"product_name": productName,
			"old_weight":   oldWeight,
			"new_weight":   newWeight,
		},
	}
}

func ShippingCostChanged(courier, courierService string, oldCost decimal.Decimal, newCost decimal.Decimal) *Error {
	oldCostFloat, _ := oldCost.Float64()
	newCostFloat, _ := newCost.Float64()

	message := replacePlaceholders(ShippingCostChangedMessage, Params{
		"courier":        courier,
		"courierService": courierService,
		"oldCost":        utils.FormatRupiah(oldCostFloat),
		"newCost":        utils.FormatRupiah(newCostFloat),
	})

	return &Error{
		StatusCode: 409,
		Code:       ShippingCostChangedErrorCode,
		Message:    message,
		CausedBy:   errors.New(message),
		Details: map[string]interface{}{
			"courier":         courier,
			"courier_service": courierService,
			"old_cost":        oldCost,
			"new_cost":        newCost,
		},
	}
}

func Validation(validationError ValidationError) *Error {
	return &Error{
		StatusCode: 400,
		Code:       ValidationErrorCode,
		Message:    ValidationErrorMessage,
		CausedBy:   errors.New("validation error"),
		Details:    validationError,
	}
}

func IsAPIError(err error) bool {
	var apiError *Error
	return errors.As(err, &apiError)
}

func replacePlaceholders(message string, params Params) string {
	for key, val := range params {
		message = strings.ReplaceAll(message, fmt.Sprintf("{%s}", key), val)
	}

	return message
}

type ValidationError map[string][]string

func NewValidationError() ValidationError {
	return make(ValidationError)
}

func (v ValidationError) Add(key string, message string) {
	v[key] = append(v[key], message)
}

func (v ValidationError) GetError() *Error {
	if len(v) > 0 {
		return Validation(v)
	}
	return nil
}
