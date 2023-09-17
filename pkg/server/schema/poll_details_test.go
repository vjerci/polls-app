package schema_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	schema "github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

func TestPollDetailsErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		InputError    error
	}{
		{
			ExpectedError: schema.ErrPollDetailsEmptyPollID,
			InputError:    model.ErrPollDetailsIDEmpty,
		},
		{
			ExpectedError: schema.ErrPollDetailsEmptyUserID,
			InputError:    model.ErrPollDetailsUserIDEmpty,
		},
		{
			ExpectedError: schema.ErrPollDetailsNotFound,
			InputError:    model.ErrPollDetailsNoPoll,
		},
		{
			ExpectedError: schema.ErrPollDetailsModel,
			InputError:    errors.New("test error"),
		},
	}

	for _, test := range testCases {
		schemaMap := &schema.PollDetailsSchemaMap{}

		err := schemaMap.ErrorHandler(test.InputError)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}
	}
}

func TestPollDetailsResponseConversion(t *testing.T) {
	t.Parallel()

	output := model.PollDetailsResponse{
		ID:         "testID",
		Name:       "testName",
		UserAnswer: "testUserAnswer",
		Answers: []model.PollDetailsAnswer{
			{
				Name:       "answerName1",
				ID:         "answerID1",
				VotesCount: 1,
			},
		},
	}

	schemaMap := &schema.PollDetailsSchemaMap{}

	converted := schemaMap.Response(&output)

	assert.EqualValues(t, output.Name, converted.Name, "expected name values to be equal")
	assert.EqualValues(t, output.ID, converted.ID, "expected id values to be equal")
	assert.EqualValues(t, output.UserAnswer, converted.UserVote, "expected user answer values to be equal")

	for i, answer := range converted.Answers {
		assert.EqualValues(t, answer.Name, converted.Answers[i].Name, "expected answers name to match")
		assert.EqualValues(t, answer.ID, converted.Answers[i].ID, "expected answers ids to match")
		assert.EqualValues(t, answer.VotesCount, converted.Answers[i].VotesCount, "expected answers votes counts to match")
	}
}
