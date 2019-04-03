package basicauth_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/httptest"
	basicauth "github.com/gobuffalo/mw-basicauth"
	"github.com/stretchr/testify/require"
)

func app() *buffalo.App {
	h := func(c buffalo.Context) error {
		return c.Render(200, nil)
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
	r := require.New(t)

	w := httptest.New(app())

	// missing authorization
	res := w.HTML("/").Get()
	r.Equal(401, res.Code)
	r.Contains(res.Header().Get("WWW-Authenticate"), `Basic realm="Basic Authentication"`)
	r.Contains(res.Body.String(), "Unauthorized")

	// bad header value, not Basic
	req := w.HTML("/")
	req.Headers["Authorization"] = "badcreds"
	res = req.Get()
	r.Equal(401, res.Code)
	r.Contains(res.Body.String(), "Unauthorized")

	// bad cred values
	req = w.HTML("/")
	req.Headers["Authorization"] = "bad creds"
	res = req.Get()
	r.Equal(401, res.Code)
	r.Contains(res.Body.String(), "Unauthorized")

	// invalid cred values in authorization
	creds := base64.StdEncoding.EncodeToString([]byte("badcredvalue"))
	req = w.HTML("/")
	req.Headers["Authorization"] = fmt.Sprintf("Basic %s", creds)
	res = req.Get()
	r.Equal(401, res.Code)
	r.Contains(res.Body.String(), "Unauthorized")

	// wrong cred values in authorization
	creds = base64.StdEncoding.EncodeToString([]byte("foo:bar"))
	req = w.HTML("/")
	req.Headers["Authorization"] = fmt.Sprintf("Basic %s", creds)
	res = req.Get()
	r.Equal(401, res.Code)
	r.Contains(res.Body.String(), "Unauthorized")

	// valid cred values
	creds = base64.StdEncoding.EncodeToString([]byte("tester:pass123"))
	req = w.HTML("/")
	req.Headers["Authorization"] = fmt.Sprintf("Basic %s", creds)
	res = req.Get()
	r.Equal(200, res.Code)
}
