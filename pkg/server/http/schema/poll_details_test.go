package schema_test

import (
	"errors"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/polls-app/pkg/domain/model/poll"
	schema "github.com/vjerci/polls-app/pkg/server/http/schema"
)

func TestPollDetailsErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		InputError    error
	}{
		{
			ExpectedError: schema.ErrPollDetailsEmptyPollID,
			InputError:    poll.ErrDetailsIDEmpty,
		},
		{
			ExpectedError: schema.ErrPollDetailsEmptyUserID,
			InputError:    poll.ErrDetailsUserIDEmpty,
		},
		{
			ExpectedError: schema.ErrPollDetailsNotFound,
			InputError:    poll.ErrDetailsNoPoll,
		},
		{
			ExpectedError: schema.ErrPollDetailsModel,
			InputError:    errors.New("test error"),
		},
	}

	for _, test := range testCases {
		schemaMap := &schema.PollDetailsSchemaMap{}

		err := schemaMap.ErrorHandler(test.InputError)

		assert.EqualValues(t, test.ExpectedError.Code, err.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, err.Message, "expected error message to match")
	}
}

func TestPollDetailsResponseConversion(t *testing.T) {
	t.Parallel()

	output := poll.DetailsResponse{
		ID:         "testID",
		Name:       "testName",
		UserAnswer: "testUserAnswer",
		Answers: []poll.DetailsAnswer{
			{
				ID:         "answerID1",
				Name:       "answerName1",
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
