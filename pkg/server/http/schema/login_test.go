package schema_test

import (
	"errors"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/polls-app/pkg/domain/model/auth"
	schema "github.com/vjerci/polls-app/pkg/server/http/schema"
)

func TestLoginErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		InputError    error
	}{
		{
			ExpectedError: schema.ErrLoginUserDoesNotExist,
			InputError:    auth.ErrLoginUserNotFound,
		},
		{
			ExpectedError: schema.ErrLoginUserIDNotSet,
			InputError:    auth.ErrLoginUserIDNotSet,
		},
		{
			ExpectedError: schema.ErrLoginModel,
			InputError:    errors.New("test error"),
		},
	}

	for _, test := range testCases {
		schemaMap := &schema.LoginSchemaMap{}

		err := schemaMap.ErrorHandler(test.InputError)

		assert.EqualValues(t, test.ExpectedError.Code, err.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, err.Message, "expected error message to match")
	}
}

func TestLoginInputConversion(t *testing.T) {
	t.Parallel()

	input := schema.LoginRequest{
		UserID: "test",
	}

	converted := input.ToModel()

	assert.EqualValues(t, input.UserID, converted.UserID, "expected userID values to be equal")
}

func TestLoginResponseConversion(t *testing.T) {
	t.Parallel()

	output := auth.LoginResponse{
		Token: "token",
		Name:  "Jhon",
	}

	schemaMap := &schema.LoginSchemaMap{}

	converted := schemaMap.Response(&output)

	assert.EqualValues(t, output.Token, converted.Token, "expected token values to be equal")
	assert.EqualValues(t, output.Name, converted.Name, "expected name values to be equal")
}
