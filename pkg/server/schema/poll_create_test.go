package schema_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	schema "github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

func TestPollCreateErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		InputError    error
	}{
		{
			ExpectedError: schema.ErrPollCreateAnswerEmpty,
			InputError:    model.ErrPollCreateAnswerEmpty,
		},
		{
			ExpectedError: schema.ErrPollCreateAnswersLen,
			InputError:    model.ErrPollCreateAnswersLen,
		},
		{
			ExpectedError: schema.ErrPollCreateNameEmpty,
			InputError:    model.ErrPollCreateNameEmpty,
		},
		{
			ExpectedError: schema.ErrPollCreateModel,
			InputError:    errors.New("test error"),
		},
	}

	for _, test := range testCases {
		schemaMap := schema.NewSchemaMap()

		err := schemaMap.PollCreateError(test.InputError)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}
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

	output := model.PollCreateResponse{
		PollID:     "pollID",
		AnswersIDS: []string{"answer1", "answer2"},
	}

	converted := schema.NewSchemaMap().PollCreateResponse(&output)

	assert.EqualValues(t, output.PollID, converted.ID, "expected poll id values to be equal")

	for i, answer := range output.AnswersIDS {
		assert.EqualValues(t, answer, converted.AnswersIDS[i], "expected answer id values to be equal")
	}
}
