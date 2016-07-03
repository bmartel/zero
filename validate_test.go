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
		"name.min": "%s must be at least %s characters",
		"age.gte":  "%s must be at least %s",
	}
}

type Post struct {
	zero.Validation
	Title string `valid:"required,min=3,max=64"`
	Body  string `valid:"required,min=10,max=512"`
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

func TestValidationFailWithDefault(t *testing.T) {
	post := Post{
		Title: "t",
		Body:  "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	}

	expectedMsgs := map[string][]string{
		"title": []string{"title must have minimum size 3"},
		"body":  []string{"body must have maximum size 512"},
	}

	v := zero.New("valid")
	msgs, isValid := v.Validate(post)

	assert.False(t, isValid, "it should fail validation")
	assert.Equal(t, msgs, expectedMsgs, "it should contain the custom error messages")
}

func TestValidationSuccessWithDefault(t *testing.T) {
	post := Post{
		Title: "test",
		Body:  "this is the body content",
	}

	v := zero.New("valid")
	msgs, isValid := v.Validate(post)

	assert.Empty(t, msgs)
	assert.True(t, isValid, "it should pass validation")
}
