package schema

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var ErrAuthMissing = &echo.HTTPError{
	Message:  "authentication header missing",
	Code:     http.StatusUnauthorized,
	Internal: nil,
}

var ErrAuthInvalid = &echo.HTTPError{
	Message:  "authentication invalid",
	Code:     http.StatusUnauthorized,
	Internal: nil,
}

var ErrAuthUserDoesNotExist = &echo.HTTPError{
	Message:  "authentication invalid, user does not exist in db",
	Code:     http.StatusUnauthorized,
	Internal: nil,
}
