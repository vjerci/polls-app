package schema_test

import (
	"errors"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model/poll"
	schema "github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/schema"
)

func TestPollCreateErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		InputError    error
	}{
		{
			ExpectedError: schema.ErrPollCreateAnswerEmpty,
			InputError:    poll.ErrCreateAnswerEmpty,
		},
		{
			ExpectedError: schema.ErrPollCreateAnswersLen,
			InputError:    poll.ErrCreateAnswersLen,
		},
		{
			ExpectedError: schema.ErrPollCreateNameEmpty,
			InputError:    poll.ErrCreateNameEmpty,
		},
		{
			ExpectedError: schema.ErrPollCreateModel,
			InputError:    errors.New("test error"),
		},
	}

	for _, test := range testCases {
		schemaMap := &schema.PollCreateSchemaMap{}

		err := schemaMap.ErrorHandler(test.InputError)

		assert.EqualValues(t, test.ExpectedError.Code, err.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, err.Message, "expected error message to match")
	}
}

func TestPollCreateInputConversion(t *testing.T) {
	t.Parallel()

	input := schema.PollCreateRequest{
		Name:    "pollName",
		Answers: []string{"answer1", "answer2"},
	}

	converted := input.ToModel()

	assert.EqualValues(t, input.Name, converted.Name, "expected poll name values to be equal")

	for i, answer := range input.Answers {
		assert.EqualValues(t, answer, converted.Answers[i], "expected answers name values to be equal")
	}
}

func TestPollCreateResponseConversion(t *testing.T) {
	t.Parallel()

	output := poll.CreateResponse{
		PollID:     "pollID",
		AnswersIDS: []string{"answer1", "answer2"},
	}

	schemaMap := &schema.PollCreateSchemaMap{}

	converted := schemaMap.Response(&output)

	assert.EqualValues(t, output.PollID, converted.ID, "expected poll id values to be equal")

	for i, answer := range output.AnswersIDS {
		assert.EqualValues(t, answer, converted.AnswersIDS[i], "expected answer id values to be equal")
	}
}
