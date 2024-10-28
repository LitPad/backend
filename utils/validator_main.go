package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	customValidator *validator.Validate
	translator      ut.Translator
)

func RegisterTagName() {
	customValidator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})
}

// Initialize the custom validator and translator
func init() {
	customValidator = validator.New()
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	translator, _ = uni.GetTranslator("en")

	// Register Custom Validators
	customValidator.RegisterValidation("payment_type_validator", PaymentTypeValidator)
	customValidator.RegisterValidation("subscription_type_validator", SubscriptionTypeValidator)
	customValidator.RegisterValidation("rating_choice_validator", RatingChoiceValidator)
	customValidator.RegisterValidation("age_discretion_validator", AgeDiscretionValidator)
	customValidator.RegisterValidation("chapter_status_validator", ChapterStatusValidator)
	customValidator.RegisterValidation("contract_type_validator", ContractTypeChoiceValidator)
	customValidator.RegisterValidation("contract_id_type_validator", ContractIDTypeChoiceValidator)
	customValidator.RegisterValidation("contract_status_validator", ContractStatusChoiceValidator)

	RegisterTagName()
}

// Register translations
func registerTranslations(param string) {
	registerTranslation := func(tag, translation string, translator ut.Translator) {
		customValidator.RegisterTranslation(tag, translator, func(ut ut.Translator) error {
			return ut.Add(tag, translation, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, fe.Field())
			return t
		})
	}

	registerTranslation("gt", "Value is too small!", translator)
	registerTranslation("required", "This field is required.", translator)
	registerTranslation("required_if", "This field is required.", translator)
	registerTranslation("required_without", "This field is required.", translator)
	registerTranslation("payment_type_validator", "Invalid payment type", translator)
	registerTranslation("subscription_type_validator", "Invalid subscription type", translator)
	registerTranslation("rating_choice_validator", "Invalid rating choice", translator)
	registerTranslation("age_discretion_validator", "Invalid age discretion. Choices are 4, 12, 16, 18", translator)
	registerTranslation("chapter_status_validator", "Invalid chapter status. Choices are DRAFT, PUBLISHED, TRASH", translator)
	registerTranslation("contract_type_validator", "Invalid contract type. Choices are EXCLUSIVE, NON-EXCLUSIVE, ONLY-EXCLUSIVE", translator)
	registerTranslation("contract_id_type_validator", "Invalid ID type. Choices are DRIVERS-LICENSE, GOVERNMENT-ID, PASSPORT", translator)
	registerTranslation("contract_status_validator", "Invalid status type. Choices are PENDING, APPROVED, DECLINED, UPDATED", translator)

	minErrMsg := fmt.Sprintf("%s characters min", param)
	registerTranslation("min", minErrMsg, translator)
	maxErrMsg := fmt.Sprintf("%s characters max", param)
	registerTranslation("max", maxErrMsg, translator)
	registerTranslation("email", "Invalid Email", translator)
	eqErrMsg := fmt.Sprintf("Must be %s", param)
	registerTranslation("eq", eqErrMsg, translator)
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

// CustomValidator is a custom validator that uses "github.com/go-playground/validator/v10"
type CustomValidator struct{}

// Validate performs the validation of the given struct
func (cv *CustomValidator) Validate(i interface{}) *ErrorResponse {
	if err := customValidator.Struct(i); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return cv.translateValidationErrors(errs)
		}
		return &ErrorResponse{Message: "Validation error"}
	}
	return nil
}

// translateValidationErrors translates the validation errors to custom errors
func (cv *CustomValidator) translateValidationErrors(errs validator.ValidationErrors) *ErrorResponse {
	errData := make(map[string]string)
	for _, err := range errs {
		errParam := err.Param()
		registerTranslations(errParam)
		errMsg := err.Translate(translator)
		errField := err.Field()
		errData[errField] = errMsg
	}
	errResp := RequestErr(ERR_INVALID_ENTRY, "Invalid Entry", errData)
	return &errResp
}

// New creates a new instance of CustomValidator
func Validator() *CustomValidator {
	return &CustomValidator{}
}
