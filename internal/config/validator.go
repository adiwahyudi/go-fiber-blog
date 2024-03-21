package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func NewValidator() *validator.Validate {
	v := validator.New()
	if err := v.RegisterValidation("validateEmail", validateEmail); err != nil {
		// Handle error if registration fails
		logrus.Panic(err)
	}
	return v
}

func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	if email == "" {
		return true // No need to validate if empty
	}
	// Validate email format
	return validator.New().Var(email, "email") == nil
}
