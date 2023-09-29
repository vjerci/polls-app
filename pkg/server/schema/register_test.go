package schema_test

import (
	"errors"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

func TestRegisterErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		InputError    error
	}{

		{
			ExpectedError: schema.ErrRegisterUserIDNotSet,
			InputError:    model.ErrRegisterUserIDNotSet,
		},
		{
			ExpectedError: schema.ErrRegisterNameNotSet,
			InputError:    model.ErrRegisterNameNotSet,
		},
		{
			ExpectedError: schema.ErrRegisterUserDuplicate,
			InputError:    model.ErrRegisterDuplicate,
		},
		{
			ExpectedError: schema.ErrRegisterModel,
			InputError:    errors.New("generic error from db"),
		},
	}

	for _, test := range testCases {
		schemaMap := &schema.RegisterSchemaMap{}

		err := schemaMap.ErrorHandler(test.InputError)

		assert.EqualValues(t, test.ExpectedError.Code, err.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, err.Message, "expected error message to match")
	}
}

func TestRegisterInputConversion(t *testing.T) {
	t.Parallel()

	input := schema.RegisterRequest{
		UserID: "test_user_id",
		Name:   "Jhon",
	}

	converted := input.ToModel()

	assert.EqualValues(t, input.UserID, converted.UserID, "expected userIDs to match")
	assert.EqualValues(t, input.Name, converted.Name, "expected names to match")
}

func TestRegisterResponseConversion(t *testing.T) {
	t.Parallel()

	output := auth.AccessToken("test")

	schemaMap := &schema.RegisterSchemaMap{}

	converted := schemaMap.Response(output)

	assert.EqualValues(t, string(output), converted, "expected converted value to be equal to output")
}
