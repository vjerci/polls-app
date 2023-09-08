package schema_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	schema "github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

func TestLoginErrorHandling(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		InputError    error
	}{
		{
			ExpectedError: schema.ErrLoginUserDoesNotExist,
			InputError:    model.ErrLoginUserNotFound,
		},
		{
			ExpectedError: schema.ErrLoginUserIDNotSet,
			InputError:    model.ErrLoginUserIDNotSet,
		},
		{
			ExpectedError: schema.ErrLoginModel,
			InputError:    errors.New("test error"),
		},
	}

	for _, test := range testCases {
		schemaMap := schema.NewSchemaMap()

		err := schemaMap.LoginError(test.InputError)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}
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

	output := model.LoginResponse{
		Token: "token",
		Name:  "Jhon",
	}

	converted := schema.NewSchemaMap().LoginResponse(&output)

	assert.EqualValues(t, output.Token, converted.Token, "expected token values to be equal")
	assert.EqualValues(t, output.Name, converted.Name, "expected name values to be equal")
}
