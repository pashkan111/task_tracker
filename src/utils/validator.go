package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"task_tracker/src/errors/api_errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateRequestData[T any](request_schema T, body io.Reader) (*T, error) {
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&request_schema); err != nil {
		var unmarshal_err *json.UnmarshalTypeError
		if errors.As(err, &unmarshal_err) {
			return nil, api_errors.BadRequestError{Detail: getJsonUnmarshalError(unmarshal_err)}
		}
		return nil, api_errors.BadRequestError{Detail: err.Error()}
	}

	err := validateData(&request_schema)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldError := fmt.Sprintf("Validation failed on field '%s', condition: '%s'", err.Field(), err.Tag())
			return nil, api_errors.BadRequestError{Detail: fieldError}
		}
		// TODO: add logging
		return nil, api_errors.BadRequestError{Detail: "Validation failed"}
	}
	return &request_schema, nil
}

func validateData(data interface{}) error {
	return validate.Struct(data)
}

func getJsonUnmarshalError(err *json.UnmarshalTypeError) string {
	return "Field '" + err.Field + "' must be of type " + err.Type.String()
}
