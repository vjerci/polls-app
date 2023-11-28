package mapper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	pollmodel "github.com/vjerci/polls-app/pkg/domain/model/poll"
	"github.com/vjerci/polls-app/pkg/server/grpc/mapper"
)

func TestPollListInputConversion(t *testing.T) {
	t.Parallel()

	schemaMap := &mapper.PollListSchemaMap{}

	input := &pollmodel.ListResponse{
		Polls: []pollmodel.GeneralInfo{
			{
				Name: "testPollName",
				ID:   "testPollID",
			},
		},
		HasNext: true,
	}

	response := schemaMap.Response(input)

	assert.EqualValues(t, input.HasNext, response.HasNextPage, "expected poll has next page values to be equal")

	for i, poll := range input.Polls {
		assert.EqualValues(t, poll.ID, response.Polls[i].Id, "expected polls ids to be equal")
		assert.EqualValues(t, poll.Name, response.Polls[i].Name, "expected poll names to be equal")
	}
}
