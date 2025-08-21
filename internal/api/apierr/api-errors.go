package apierr

import (
	"errors"
	"fmt"
	"strings"
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

func InvalidRequest(causedBy error) *Error {
	return &Error{
		StatusCode: 500,
		Code:       InvalidRequestErrorCode,
		Message:    InvalidRequestErrorMessage,
		CausedBy:   causedBy,
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

func replacePlaceholders(message string, params Params) (result string) {
	for key, val := range params {
		result = strings.ReplaceAll(message, fmt.Sprintf("{%s}", key), val)
	}
	return
}

type ValidationError map[string][]string

func NewValidationError() ValidationError {
	return make(ValidationError)
}

func (v ValidationError) Add(key string, message string) {
	v[key] = append(v[key], message)
}

func (v ValidationError) Error() *Error {
	if len(v) > 0 {
		return Validation(v)
	}
	return nil
}
