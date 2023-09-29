package schema_test

import (
	"errors"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	schema "github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

func TestPollListErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		InputError    error
	}{
		{
			ExpectedError: schema.ErrPollListInvalidPage,
			InputError:    model.ErrPollListInvalidPage,
		},
		{
			ExpectedError: schema.ErrPollListNoData,
			InputError:    model.ErrPollListNoPolls,
		},
		{
			ExpectedError: schema.ErrPollListModel,
			InputError:    errors.New("test error"),
		},
	}

	for _, test := range testCases {
		schemaMap := &schema.PollListSchemaMap{}

		err := schemaMap.ErrorHandler(test.InputError)

		assert.EqualValues(t, test.ExpectedError.Code, err.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, err.Message, "expected error message to match")
	}
}

func TestPollListInputConversion(t *testing.T) {
	t.Parallel()

	input := schema.PollListRequest{
		Page: 1,
	}

	converted := input.ToModel()

	assert.EqualValues(t, input.Page, converted.Page, "expected page values to be equal")
}

func TestPollListResponseConversion(t *testing.T) {
	t.Parallel()

	output := model.PollListResponse{
		Polls: []model.GeneralPollInfo{
			{
				Name: "test name 1",
				ID:   "test id 1",
			},
			{
				Name: "test name 2",
				ID:   "test id 2",
			},
		},
		HasNext: true,
	}

	schemaMap := &schema.PollListSchemaMap{}

	converted := schemaMap.Response(&output)

	assert.EqualValues(t, output.HasNext, converted.HasNext, "expected has next page value to be equal")

	for i, v := range output.Polls {
		assert.EqualValues(t, v.Name, converted.Polls[i].Name, "expected name values to be equal")
		assert.EqualValues(t, v.ID, converted.Polls[i].ID, "expected id values to be equal")
	}
}
