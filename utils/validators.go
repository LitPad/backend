package utils

import (
	"github.com/LitPad/backend/models/choices"
	"github.com/go-playground/validator/v10"
)

// Validates if a payment type value is the correct one
func PaymentTypeValidator(fl validator.FieldLevel) bool {
	paymentTypeValue := fl.Field().Interface().(choices.PaymentType)
	switch paymentTypeValue {
	case choices.PTYPE_GPAY, choices.PTYPE_PAYPAL, choices.PTYPE_STRIPE:
		return true
	}
	return false // Error. Value doesn't match the required
}