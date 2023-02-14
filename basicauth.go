package basicauth

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
)

var (
	// These errors are internal and will not be seen in production mode

	// ErrNoCreds is returned when no basic auth credentials are defined
	ErrNoCreds = errors.New("no basic auth credentials defined")

	// ErrAuthFail is returned when the client fails basic authentication
	ErrAuthFail = errors.New("invalid basic auth username or password")

	// ErrUnauthorized is returned in any case the basic authentication fails
	ErrUnauthorized = errors.New("Unauthorized")
)

// Authorizer is used to authenticate the basic auth username/password.
// Should return true/false and/or an error.
type Authorizer func(buffalo.Context, string, string) (bool, error)

// Middleware enables basic authentication
func Middleware(auth Authorizer) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			token := strings.SplitN(c.Request().Header.Get("Authorization"), " ", 2)
			if len(token) != 2 {
				return responseUnauthorized(c, ErrNoCreds)
			}
			b, err := base64.StdEncoding.DecodeString(token[1])
			if err != nil {
				return responseUnauthorized(c, ErrNoCreds)
			}

			pair := strings.SplitN(string(b), ":", 2)
			if len(pair) != 2 {
				return responseUnauthorized(c, ErrUnauthorized)
			}

			success, err := auth(c, pair[0], pair[1])
			if err != nil {
				// log only this situation since it is an internal error
				c.Logger().Errorf("authorizer error: %v", err)
				return responseUnauthorized(c, ErrUnauthorized)
			}
			if !success {
				return responseUnauthorized(c, ErrAuthFail)
			}

			return next(c)
		}
	}
}

func responseUnauthorized(c buffalo.Context, err error) error {
	// Always uses status 401 but internally propagate the original error.
	// The error could be used in error handlers to handle extra steps.
	c.Response().Header().Set("WWW-Authenticate", `Basic realm="Basic Authentication"`)
	return c.Error(http.StatusUnauthorized, err)
}
