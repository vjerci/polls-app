package mapper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/polls-app/pkg/domain/util/auth"
	"github.com/vjerci/polls-app/pkg/server/grpc/mapper"
)

func TestRegisterInputConversion(t *testing.T) {
	t.Parallel()

	schemaMap := &mapper.RegisterSchemaMap{}

	input := auth.AccessToken("testToken")
	response := schemaMap.Response(input)

	assert.EqualValues(t, input, response.Token, "expected token values to be equal")
}
