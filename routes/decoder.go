package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/LitPad/backend/utils"
	"github.com/gofiber/fiber/v2"
)

var validator = utils.Validator()

func DecodeJSONBody(c *fiber.Ctx, dst interface{}) (int, *utils.ErrorResponse) {
	var errData *utils.ErrorResponse
	code := 200
	if !strings.Contains(c.Get("Content-Type"), "application/json") {
		errD := utils.RequestErr(utils.ERR_INVALID_REQUEST, "Content-Type header is not application/json")
		errData = &errD
		return code, errData
	}

	dec := json.NewDecoder(bytes.NewReader(c.Body()))
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	msg := "Invalid Entry"
	fieldErrors := make(map[string]string)
	status_code := 422
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		errStr := err.Error()
		switch {
		case errors.As(err, &syntaxError):
			msg = fmt.Sprintf(
				"Request body contains badly-formed JSON (at position %d)",
				syntaxError.Offset,
			)

		case errors.Is(err, io.ErrUnexpectedEOF):
			status_code = http.StatusBadRequest
			msg = "Request body contains badly-formed JSON"

		case errors.As(err, &unmarshalTypeError):
			fieldName := unmarshalTypeError.Field
			fieldErrors[fieldName] = "Invalid format"
		case strings.HasPrefix(errStr, "json: unknown field "):
			fieldName := strings.TrimPrefix(errStr, "json: unknown field ")
			fieldErrors[fieldName] = "Unknown field"
		case errors.Is(err, io.EOF):
			status_code = http.StatusBadRequest
			msg = "Request body must not be empty"

		case errStr == "http: request body too large":
			status_code = http.StatusRequestEntityTooLarge
			msg = "Request body must not be larger than 1MB"

		default:
			status_code = 400
			msg = "Invalid request"
		}
		errData := utils.RequestErr(utils.ERR_INVALID_REQUEST, msg)
		if len(fieldErrors) > 0 {
			errData.Data = &fieldErrors
		}
		code = status_code
		return code, &errData
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		errData := utils.RequestErr(utils.ERR_INVALID_REQUEST, "Request body must only contain a single JSON object")
		return 400, &errData
	}
	return code, nil
}

func ValidateRequest(c *fiber.Ctx, data interface{}) (*int, *utils.ErrorResponse) {
	if errCode, errData := DecodeJSONBody(c, &data); errData != nil {
		return &errCode, errData
	}
	if errData := validator.Validate(data); errData != nil {
		errCode := 422
		return &errCode, errData
	}
	return nil, nil
}

func ValidateFormRequest(c *fiber.Ctx, data interface{}) (*int, *utils.ErrorResponse) {
	errC := 400
	if !strings.Contains(c.Get("Content-Type"), "multipart/form-data") {
		errD := utils.RequestErr(utils.ERR_INVALID_REQUEST, "Content-Type header is not multipart/form-data")
		return &errC, &errD
	}

	if err := c.BodyParser(data); err != nil {
		errD := utils.RequestErr(utils.ERR_INVALID_REQUEST, "Unable to parse form body")
		return &errC, &errD
	}
	if errData := validator.Validate(data); errData != nil {
		errC = 422
		return &errC, errData
	}
	return nil, nil
}

func ValidatePathParams(c *fiber.Ctx, pathParams map[string]string) (*int, *utils.ErrorResponse) {
	for paramName, paramValue := range pathParams {
		if paramValue == "" {
			errData := utils.RequestErr(utils.ERR_INVALID_REQUEST, fmt.Sprintf("Missing or invalid value for path parameter: %s", paramName))
			errCode := 400
			return &errCode, &errData
		}
	}

	return nil, nil
}
