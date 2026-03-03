package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// ErrorObject represents a single validation error.
type ErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrResponse represents a response containing multiple validation errors.
type ErrResponse struct {
	Errors []ErrorObject `json:"errors"`
}

var validate = validator.New()

// ValidateStruct validates a struct using the validator package.
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// ToErrResponse converts a validation error into an ErrResponse.
func ToErrResponse(err error) ErrResponse {
	var errs []ErrorObject
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			errs = append(errs, ErrorObject{
				Code:    "INVALID_" + strings.ToUpper(fe.Field()),
				Message: fe.Error(),
			})
		}
	} else {
		errs = append(errs, ErrorObject{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
	}
	return ErrResponse{Errors: errs}
}
