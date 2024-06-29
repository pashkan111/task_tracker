package utils

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateData(data interface{}) error {
	return validate.Struct(data)
}
