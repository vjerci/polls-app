package middleware_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/api/middleware"
)

func TestErrorHandling(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name string

		InputError error

		ExpectedResponse string
		ExpectedHTTPCode int

		ExpectedLog []byte
	}{
		{
			Name: "Http error - with string message ",

			InputError: &echo.HTTPError{
				Code:     http.StatusBadGateway,
				Message:  "serverError",
				Internal: nil,
			},

			ExpectedResponse: `{"success":false,"error":"serverError"}` + "\n",
			ExpectedHTTPCode: http.StatusBadGateway,

			ExpectedLog: []byte(`serving http error "serverError":""`),
		},
		{
			Name: "Http error - with int message ",

			InputError: &echo.HTTPError{
				Code:     http.StatusBadGateway,
				Message:  2,
				Internal: nil,
			},

			ExpectedResponse: `{"success":false,"error":"internal server error"}` + "\n",
			ExpectedHTTPCode: http.StatusBadGateway,
		},
		{
			Name: "Generic error ",

			InputError: errors.New("testError"),

			ExpectedResponse: `{"success":false,"error":"internal server error"}` + "\n",
			ExpectedHTTPCode: http.StatusInternalServerError,
		},
	}

	for _, test := range testCases {
		test := test

		t.Run(t.Name()+"-"+test.Name, func(t *testing.T) {
			t.Parallel()

			echoServer := echo.New()

			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()

			echoContext := echoServer.NewContext(req, rec)

			client := middleware.Client{}

			client.ErrorHandler(test.InputError, echoContext)

			resBytes, err := io.ReadAll(rec.Body)
			if err != nil {
				t.Fatalf("expected response to be readable")
			}

			assert.EqualValues(t, test.ExpectedResponse, string(resBytes), "response does not match")
			assert.EqualValues(t, test.ExpectedHTTPCode, rec.Code, "response status code does not match")
		})
	}
}
