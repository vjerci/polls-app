package login_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/login"
)

var signingKey = "y5XYbfKoqOqGa+XVggMZs4yvcE6SwWS9Q9UBZviMw0A="

func TestCreateToken(t *testing.T) {
	t.Parallel()

	loginClient := login.New(signingKey)
	token, err := loginClient.CreateToken("userID")

	if err != nil {
		t.Fatalf(`expected no error got "%s" instead`, err)
	}

	if token == "" {
		t.Fatalf("got empty token expected a functional one")
	}
}

func TestDecodeToken(t *testing.T) {
	t.Parallel()

	startingUserID := "userID"

	loginClient := login.New(signingKey)
	token, err := loginClient.CreateToken(startingUserID)

	if err != nil {
		t.Fatalf(`expected no error while encoding got "%s" instead`, err)
	}

	userID, err := loginClient.Decode(token)

	if err != nil {
		t.Fatalf(`expected no error while decoding token got "%s" instead`, err)
	}

	assert.EqualValues(t, startingUserID, userID, "expected user_id's to match")
}
