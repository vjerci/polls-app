package schema_test

import (
	"errors"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/polls-app/pkg/domain/model/poll"
	schema "github.com/vjerci/polls-app/pkg/server/http/schema"
)

func TestPollVoteErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError *echo.HTTPError
		InputError    error
	}{
		{
			ExpectedError: schema.ErrPollVoteInvalidVote,
			InputError:    poll.ErrVoteAnswerNotFound,
		},
		{
			ExpectedError: schema.ErrPollVoteInvalidPollID,
			InputError:    poll.ErrVotePollIDEmpty,
		},
		{
			ExpectedError: schema.ErrPollVoteInvalidAnswerID,
			InputError:    poll.ErrVoteAnswerIDEmpty,
		},
		{
			ExpectedError: schema.ErrPollVoteInvalidUserID,
			InputError:    poll.ErrVoteUserIDEmpty,
		},
		{
			ExpectedError: schema.ErrPollVoteModel,
			InputError:    errors.New("test error"),
		},
	}

	for _, test := range testCases {
		schemaMap := &schema.PollVoteSchemaMap{}

		err := schemaMap.ErrorHandler(test.InputError)

		assert.EqualValues(t, test.ExpectedError.Code, err.Code, "expected http status code to match")
		assert.EqualValues(t, test.ExpectedError.Message, err.Message, "expected error message to match")
	}
}

func TestPollVoteResponseConversion(t *testing.T) {
	t.Parallel()

	output := poll.VoteResponse{
		ModifiedAnswer: true,
	}

	schemaMap := &schema.PollVoteSchemaMap{}

	converted := schemaMap.Response(&output)

	assert.EqualValues(t,
		output.ModifiedAnswer,
		converted.ModifiedAnswer,
		"expected properties modified answer to match after conversion")
}
