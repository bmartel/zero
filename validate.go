package zero

import "gopkg.in/go-playground/validator.v8"

// New ... Create a validator instance and bind custom validation types
func New(tagName string) *Zero {
	valid := validator.New(&validator.Config{TagName: tagName})

	// Validator instance
	return &Zero{valid}
}

// Zero is a convenience wrapper for the go-playground/validator
type Zero struct {
	validator *validator.Validate
}

// Register custom validation funcs
func (z *Zero) Register(validators map[string]validator.Func) {
	// Custom Validations
	for validatorName, validatorFunc := range validators {
		z.validator.RegisterValidation(validatorName, validatorFunc)
	}
}

// Validate will validate a Validator type and return any error messages as well as validation status
func (z *Zero) Validate(v Validator) (map[string][]string, bool) {
	err := z.validator.Struct(v)
	if err == nil {
		return make(map[string][]string), true
	}

	return errors(err, v.Validates()), false
}

// Validator is the interface that must be adhered for any types being validated
type Validator interface {
	Validates() map[string]string
}

// Errors formats the errors returned from validation failure
func errors(err error, validationMessages map[string]string) map[string][]string {
	messages := make(map[string][]string, 0)

	switch validationErrors := err.(type) {
	case validator.ValidationErrors:
		for _, validationError := range validationErrors {
			field := toSnake(validationError.Field)
			if msg := validationMessages[field+"."+validationError.Tag]; msg != "" {
				messages[field] = append(messages[field], msg)
			}
		}
	}

	return messages
}
