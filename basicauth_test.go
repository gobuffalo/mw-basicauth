package basicauth_test

import (
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/httptest"
	basicauth "github.com/gobuffalo/mw-basicauth"
	"github.com/stretchr/testify/require"
)

func app() *buffalo.App {
	h := func(c buffalo.Context) error {
		return c.Render(200, render.String("Welcome"))
	}
	auth := func(c buffalo.Context, u, p string) (bool, error) {
		return (u == "tester" && p == "pass123"), nil
	}
	a := buffalo.New(buffalo.Options{})
	a.Use(basicauth.Middleware(auth))
	a.GET("/", h)
	return a
}

func TestBasicAuth(t *testing.T) {
	tests := []struct {
		status  int
		name    string
		auth    string
		message string
	}{
		{http.StatusUnauthorized, "missing", "MISSING", "no basic auth credentials defined"},
		{http.StatusUnauthorized, "empty", "", "no basic auth credentials defined"},
		{http.StatusUnauthorized, "badcreds", "badcreds", "no basic auth credentials defined"},
		{http.StatusUnauthorized, "bad creds", "bad creds", "no basic auth credentials defined"},
		{http.StatusUnauthorized, "invalid", "Basic " + base64.StdEncoding.EncodeToString([]byte("badcredvalue")), "Unauthorized"},
		{http.StatusUnauthorized, "wrong", "Basic " + base64.StdEncoding.EncodeToString([]byte("foo:bar")), "invalid basic auth username"},
		{http.StatusOK, "valid", "Basic " + base64.StdEncoding.EncodeToString([]byte("tester:pass123")), "Welcome"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)
			w := httptest.New(app())

			// missing authorization
			req := w.HTML("/")
			if tt.auth != "MISSING" {
				req.Headers["Authorization"] = tt.auth
			}
			res := req.Get()
			r.Equal(tt.status, res.Code)
			if tt.status == http.StatusUnauthorized {
				r.Contains(res.Header().Get("WWW-Authenticate"), `Basic realm="Basic Authentication"`)
			}
			r.Contains(res.Body.String(), tt.message)
		})
	}
}
