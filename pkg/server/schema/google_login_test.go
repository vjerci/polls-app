package schema_test

import (
	"errors"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/auth"
	schema "github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

func TestGoogleLoginErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		InputError    error
	}{
		{
			ExpectedError: schema.ErrGoogleLoginTokenNotSet,
			InputError:    auth.ErrGoogleLoginTokenNotSet,
		},
		{
			ExpectedError: schema.ErrLoginModel,
			InputError:    errors.New("test error"),
		},
	}

	for _, test := range testCases {
		schemaMap := &schema.GoogleLoginSchemaMap{}

		err := schemaMap.ErrorHandler(test.InputError)

		assert.EqualValues(t, test.ExpectedError.Code, err.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, err.Message, "expected error message to match")
	}
}

func TestGoogleLoginInputConversion(t *testing.T) {
	t.Parallel()

	input := schema.GoogleLoginRequest{
		Token: "test",
	}

	converted := input.ToModel()

	assert.EqualValues(t, input.Token, converted.Token, "expected token values to be equal")
}

func TestGoogleLoginResponseConversion(t *testing.T) {
	t.Parallel()

	output := auth.GoogleLoginResponse{
		Token: "token",
		Name:  "Jhon",
	}

	schemaMap := &schema.GoogleLoginSchemaMap{}

	converted := schemaMap.Response(&output)

	assert.EqualValues(t, output.Token, converted.Token, "expected token values to be equal")
	assert.EqualValues(t, output.Name, converted.Name, "expected name values to be equal")
}
