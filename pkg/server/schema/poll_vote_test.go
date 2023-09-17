package schema_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	schema "github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

func TestPollVoteErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExpectedError error
		InputError    error
	}{
		{
			ExpectedError: schema.ErrPollVoteInvalidVote,
			InputError:    model.ErrPollVoteAnswerNotFound,
		},
		{
			ExpectedError: schema.ErrPollVoteInvalidPollID,
			InputError:    model.ErrPollVotePollIDEmpty,
		},
		{
			ExpectedError: schema.ErrPollVoteInvalidAnswerID,
			InputError:    model.ErrPollVoteAnswerIDEmpty,
		},
		{
			ExpectedError: schema.ErrPollVoteInvalidUserID,
			InputError:    model.ErrPollVoteUserIDEmpty,
		},
		{
			ExpectedError: schema.ErrPollVoteModel,
			InputError:    errors.New("test error"),
		},
	}

	for _, test := range testCases {
		schemaMap := &schema.PollVoteSchemaMap{}

		err := schemaMap.ErrorHandler(test.InputError)

		if !errors.Is(err, test.ExpectedError) {
			t.Fatalf(`expected to get error "%s" got "%s" instead`, test.ExpectedError, err)
		}
	}
}

func TestPollVoteResponseConversion(t *testing.T) {
	t.Parallel()

	output := model.PollVoteResponse{
		ModifiedAnswer: true,
	}

	schemaMap := &schema.PollVoteSchemaMap{}

	converted := schemaMap.Response(&output)

	assert.EqualValues(t,
		output.ModifiedAnswer,
		converted.ModifiedAnswer,
		"expected properties modified answer to match after conversion")
}
