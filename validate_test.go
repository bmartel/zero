package zero_test

import (
	"testing"

	"github.com/bmartel/zero"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string `valid:"required,min=3"`
	Age  int    `valid:"gte=18"`
}

func (User) Validates() map[string]string {
	return map[string]string{
		"name.required": "name is required",
		"name.min":      "name must be at least 3 characters",
		"age.gte":       "age must be at least 18",
	}
}

func TestValidatorCreate(t *testing.T) {
	v := zero.New("valid")

	assert.IsType(t, &zero.Zero{}, v, "it should create a zero validation struct")
}

func TestValidatorValidationFail(t *testing.T) {
	user := User{
		Age: 17,
	}

	v := zero.New("valid")
	msgs, isValid := v.Validate(user)

	expectedMsgs := map[string][]string{
		"name": []string{"name is required"},
		"age":  []string{"age must be at least 18"},
	}

	assert.False(t, isValid, "it should fail validation")
	assert.Equal(t, msgs, expectedMsgs, "it should contain the custom error messages")
}

func TestValidatorValidationSuccess(t *testing.T) {
	user := User{
		Name: "test",
		Age:  18,
	}

	v := zero.New("valid")
	msgs, isValid := v.Validate(user)

	assert.Empty(t, msgs)
	assert.True(t, isValid, "it should pass validation")
}
