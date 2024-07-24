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

// Validates if a rating choice value is the correct one
func RatingChoiceValidator(fl validator.FieldLevel) bool {
	ratingChoiceValue := fl.Field().Interface().(choices.RatingChoice)
	switch ratingChoiceValue {
	case choices.RC_1, choices.RC_2, choices.RC_3, choices.RC_4, choices.RC_5:
		return true
	}
	return false // Error. Value doesn't match the required
}

// Validates if a chapter status value is the correct one
func ChapterStatusValidator(fl validator.FieldLevel) bool {
	chapterStatusValue := fl.Field().Interface().(choices.ChapterStatus)
	switch chapterStatusValue {
	case choices.CS_DRAFT, choices.CS_PUBLISHED, choices.CS_TRASH:
		return true
	}
	return false // Error. Value doesn't match the required
}

// Validates if a age discretion value is the correct one
func AgeDiscretionValidator(fl validator.FieldLevel) bool {
	paymentTypeValue := fl.Field().Interface().(choices.AgeType)
	switch paymentTypeValue {
	case choices.ATYPE_FOUR, choices.ATYPE_TWELVE, choices.ATYPE_SIXTEEN, choices.ATYPE_EIGHTEEN:
		return true
	}
	return false // Error. Value doesn't match the required
}