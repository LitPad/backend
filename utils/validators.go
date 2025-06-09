package utils

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/LitPad/backend/models/choices"
	"github.com/go-playground/validator/v10"
)

// Validates if a account type value is the correct one
func AccountTypeValidator(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(choices.AccType).IsValid()
}

// Validates if a featured content location choice value is the correct one
func FeaturedContentLocationChoiceValidator(fl validator.FieldLevel) bool {
	return fl.Field().Interface().(choices.FeaturedContentLocationChoice).IsValid()
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

func CountWords(text string) int {
    if strings.TrimSpace(text) == "" {
        return 0
    }
    
    // Split by whitespace and filter out empty strings
    fields := strings.FieldsFunc(text, func(c rune) bool {
        return unicode.IsSpace(c)
    })
    
    return len(fields)
}

// WordCountMinValidator validates minimum word count
func WordCountMinValidator(fl validator.FieldLevel) bool {
    text := fl.Field().String()
    wordCount := CountWords(text)
    minCount := fl.Param()
    
    var min int
    if _, err := fmt.Sscanf(minCount, "%d", &min); err != nil {
        return false
    }
    
    return wordCount >= min
}

// WordCountMaxValidator validates maximum word count
func WordCountMaxValidator(fl validator.FieldLevel) bool {
    text := fl.Field().String()
    wordCount := CountWords(text)
    maxCount := fl.Param()
    
    var max int
    if _, err := fmt.Sscanf(maxCount, "%d", &max); err != nil {
        return false
    }
    
    return wordCount <= max
}