package zero

import (
	"fmt"
	"strings"

	"gopkg.in/go-playground/validator.v8"
)

// New ... Create a validator instance and bind custom validation types
func New(tagName string) *Zero {
	// Validator instance
	return &Zero{validator.New(&validator.Config{TagName: tagName}), messages}
}

// Zero is a convenience wrapper for the go-playground/validator
type Zero struct {
	validator *validator.Validate
	messages  map[string]string
}

// SetMessages overrides all default field error messages
func (z *Zero) SetMessages(messages map[string]string) {
	z.messages = messages
}

// SetMessage sets a single message for a field error
func (z *Zero) SetMessage(field string, message string) {
	z.messages[field] = message
}

// AddValidator adds a custom validation func
func (z *Zero) AddValidator(name string, validatorFunc validator.Func, message string) {
	z.validator.RegisterValidation(name, validatorFunc)
	z.SetMessage(name, message)
}

// AddValidators adds a map of custom validation funcs
func (z *Zero) AddValidators(validators map[string]ValidatorFunc) {
	// Custom Validations
	for name, validation := range validators {
		z.AddValidator(name, validation.Func, validation.Message)
	}
}

// Validate will validate a Validator type and return any error messages as well as validation status
func (z *Zero) Validate(v Validator) (map[string][]string, bool) {
	err := z.validator.Struct(v)
	if err == nil {
		return make(map[string][]string), true
	}

	return z.errors(err, v.Validates()), false
}

// Errors formats the errors returned from validation failure
func (z *Zero) errors(err error, validationMessages map[string]string) map[string][]string {
	messages := make(map[string][]string, 0)

	switch validationErrors := err.(type) {
	case validator.ValidationErrors:
		for _, validationError := range validationErrors {
			field := toSnake(validationError.Field)

			// Use a struct overriden message
			msg, found := validationMessages[field+"."+validationError.Tag]
			if !found {
				// Use a default message if there is one
				msg, found = z.messages[validationError.Tag]
				if !found {
					continue
				}
			}

			switch strings.Count(msg, "%s") {
			case 1:
				msg = fmt.Sprintf(msg, strings.ToLower(validationError.Field))
			case 2:
				msg = fmt.Sprintf(msg, strings.ToLower(validationError.Field), validationError.Param)
			case 3:
				msg = fmt.Sprintf(msg, strings.ToLower(validationError.Field), validationError.Param, validationError.Value)
			}

			messages[field] = append(messages[field], msg)
		}
	}

	return messages
}

// ValidatorFunc describes a validation action and associated error message
type ValidatorFunc struct {
	Message string
	validator.Func
}

// Validator is the interface that must be adhered for any types being validated
type Validator interface {
	Validates() map[string]string
}

// Validation is a default struct for completing the validator interface
type Validation struct{}

// Validates ...
func (Validation) Validates() map[string]string {
	return map[string]string{}
}
