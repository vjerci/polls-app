package mapper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	pollmodel "github.com/vjerci/polls-app/pkg/domain/model/poll"
	"github.com/vjerci/polls-app/pkg/server/grpc/mapper"
)

func TestPollVoteInputConversion(t *testing.T) {
	t.Parallel()

	schemaMap := &mapper.PollVoteSchemaMap{}

	input := &pollmodel.VoteResponse{
		ModifiedAnswer: true,
	}

	response := schemaMap.Response(input)

	assert.EqualValues(t,
		input.ModifiedAnswer,
		response.ModifiedAnswer,
		"expected modified answer flag values to be equal",
	)
}
