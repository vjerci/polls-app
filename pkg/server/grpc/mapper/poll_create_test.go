package mapper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	pollmodel "github.com/vjerci/polls-app/pkg/domain/model/poll"
	"github.com/vjerci/polls-app/pkg/server/grpc/mapper"
)

func TestPollCreateInputConversion(t *testing.T) {
	t.Parallel()

	schemaMap := &mapper.PollCreateSchemaMap{}

	input := &pollmodel.CreateResponse{
		PollID:     "testPollID",
		AnswersIDS: []string{"testAnswer1", "testAnswer2"},
	}

	response := schemaMap.Response(input)

	assert.EqualValues(t, input.PollID, response.Id, "expected poll id values to be equal")
	assert.EqualValues(t, input.AnswersIDS, response.AnswerIds, "expected answer id values to be equal")
}
