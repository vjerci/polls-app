package mapper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	authmodel "github.com/vjerci/polls-app/pkg/domain/model/auth"
	"github.com/vjerci/polls-app/pkg/server/grpc/mapper"
)

func TestGoogleLoginInputConversion(t *testing.T) {
	t.Parallel()

	schemaMap := &mapper.GoogleLoginSchemaMap{}

	input := &authmodel.GoogleLoginResponse{
		Token: "testToken",
		Name:  "testName",
	}

	response := schemaMap.Response(input)

	assert.EqualValues(t, input.Name, response.Name, "expected name values to be equal")
	assert.EqualValues(t, input.Token, response.Token, "expected name values to be equal")
}
