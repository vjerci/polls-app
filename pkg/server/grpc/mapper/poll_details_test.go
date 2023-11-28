package mapper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	pollmodel "github.com/vjerci/polls-app/pkg/domain/model/poll"
	"github.com/vjerci/polls-app/pkg/server/grpc/mapper"
)

func TestPollDetailsInputConversion(t *testing.T) {
	t.Parallel()

	schemaMap := &mapper.PollDetailsSchemaMap{}

	input := &pollmodel.DetailsResponse{
		ID:         "testPollID",
		Name:       "testPollName",
		UserAnswer: "testUserAnswerID",
		Answers: []pollmodel.DetailsAnswer{
			{
				Name:       "testAnswer1Name",
				ID:         "testAnswer1ID",
				VotesCount: 25,
			},
		},
	}

	response := schemaMap.Response(input)

	assert.EqualValues(t, input.ID, response.Id, "expected poll id values to be equal")
	assert.EqualValues(t, input.Name, response.Name, "expected poll name values to be equal")
	assert.EqualValues(t, input.UserAnswer, response.UserVote, "expected user answer values to be equal")

	for i, answer := range input.Answers {
		assert.EqualValues(t, answer.ID, response.Answers[i].Id, "expected answer ids to be equal")
		assert.EqualValues(t, answer.Name, response.Answers[i].Name, "expected answer names to be equal")
		assert.EqualValues(t, answer.VotesCount, response.Answers[i].VotesCount, "expected answer names to be equal")
	}
}
