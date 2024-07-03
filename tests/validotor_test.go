package tests

import (
	"errors"
	"strings"
	"task_tracker/src/errors/api_errors"
	"task_tracker/src/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Example struct to use for testing
type TestRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
}

func TestValidateRequestData__Success(t *testing.T) {
	validJson := `{"name": "John Doe", "email": "john.doe@example.com"}`
	body := strings.NewReader(validJson)

	var requestSchema TestRequest
	result, err := utils.ValidateRequestData(requestSchema, body)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "John Doe", result.Name)
	assert.Equal(t, "john.doe@example.com", result.Email)
}

func TestValidateRequestData__MissingRequiredFields(t *testing.T) {
	var requestSchema TestRequest
	invalidJson := `{"name": "John Doe"}`
	body := strings.NewReader(invalidJson)

	_, err := utils.ValidateRequestData(requestSchema, body)
	assert.Error(t, err)
	var apiErr api_errors.BadRequestError
	assert.True(t, errors.As(err, &apiErr))
	assert.Contains(t, apiErr.Detail, "Validation failed on field 'Email'")
}

func TestValidateRequestData__WrongFieldType(t *testing.T) {
	var requestSchema TestRequest

	invalidFormatJson := `{"name": "John Doe", "email": 123}`
	body := strings.NewReader(invalidFormatJson)

	result, err := utils.ValidateRequestData(requestSchema, body)
	assert.Error(t, err)
	assert.Nil(t, result)
	var apiErr api_errors.BadRequestError

	assert.True(t, errors.As(err, &apiErr))
	assert.Contains(t, apiErr.Detail, "Field 'email' must be of type string")
}
