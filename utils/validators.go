package utils

import (
	"github.com/LitPad/backend/models/choices"
	"github.com/go-playground/validator/v10"
)

// Validates if a account type value is the correct one
func AccountTypeValidator(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(choices.AccType).IsValid()
}

// Validates if a payment type value is the correct one
func PaymentTypeValidator(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(choices.PaymentType).IsValid()
}

// Validates if a subscription type value is the correct one
func SubscriptionTypeValidator(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(choices.SubscriptionTypeChoice).IsValid()
}

// Validates if a rating choice value is the correct one
func RatingChoiceValidator(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(choices.RatingChoice).IsValid()
}

// Validates if a contract type choice value is the correct one
func ContractTypeChoiceValidator(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(choices.ContractTypeChoice).IsValid()
}

// Validates if a contract idtype choice value is the correct one
func ContractIDTypeChoiceValidator(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(choices.ContractIDTypeChoice).IsValid()
}

// Validates if a contract status choice value is the correct one
func ContractStatusChoiceValidator(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(choices.ContractStatusChoice).IsValid()
}

// Validates if a age discretion value is the correct one
func AgeDiscretionValidator(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(choices.AgeType).IsValid()
}

// Validates if a reply type value is the correct one
func ReplyTypeValidator(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(choices.ReplyType).IsValid()
}

// Validates if a device type value is the correct one
func DeviceTypeValidator(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(choices.DeviceType).IsValid()
}