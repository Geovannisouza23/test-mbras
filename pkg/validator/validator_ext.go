package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrResponse struct {
	Errors []ErrorObject `json:"errors"`
}

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

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
