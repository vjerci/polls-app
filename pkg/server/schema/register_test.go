package schema_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/login"
	schema "github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

func TestRegisterErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
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
		schemaMap := schema.NewSchemaMap()

		err := schemaMap.RegisterError(test.InputError)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}
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

	output := login.AccessToken("test")

	schemaMap := schema.NewSchemaMap()

	converted := schemaMap.RegisterResponse(output)

	assert.EqualValues(t, string(output), converted, "expected converted value to be equal to output")
}
