package handlers

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

var ErrInvalidErrType = errors.New("invalid error type")

type Constraint struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message     string                  `json:"message"`
	Constraints map[string][]Constraint `json:"constraints,omitempty"`
}

type ValidationError struct {
	Message     string                  `json:"message"`
	Constraints map[string][]Constraint `json:"constraints"`
}

func EchoErrorFromValidation(err error) error {
	constraints := map[string][]Constraint{}
	var valErrs validation.Errors
	if !errors.As(err, &valErrs) {
		return ErrInvalidErrType
	}

	for field, valErr := range valErrs {
		var e validation.Error
		if !errors.As(valErr, &e) {
			continue
		}

		constraints[field] = append(constraints[field], Constraint{
			Code:    e.Code(),
			Message: e.Error(),
		})
	}

	return echo.NewHTTPError(http.StatusUnprocessableEntity, ValidationError{
		Message:     "Validation error",
		Constraints: constraints,
	})
}
