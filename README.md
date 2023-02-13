# mw-basicauth

[![Standard Test](https://github.com/gobuffalo/mw-basicauth/actions/workflows/standard-go-test.yml/badge.svg)](https://github.com/gobuffalo/mw-basicauth/actions/workflows/standard-go-test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/gobuffalo/mw-basicauth.svg)](https://pkg.go.dev/github.com/gobuffalo/mw-basicauth)
[![Go Report Card](https://goreportcard.com/badge/github.com/gobuffalo/mw-basicauth)](https://goreportcard.com/report/github.com/gobuffalo/mw-basicauth)

[Basic HTTP Authentication](https://tools.ietf.org/html/rfc7617) Middleware
for [Buffalo](https://github.com/gobuffalo/buffalo)

## Installation

```console
$ go get github.com/gobuffalo/mw-basicauth
```

## Usage

```go
auth := func(c buffalo.Context, u, p string) (bool, error) {
    return (u == "username" && p == "password"), nil
}

app.Use(basicauth.Middleware(auth))
```

## Hitting protected endpoints

1. Base64 Encode `username:password`, which becomes `dXNlcm5hbWU6cGFzc3dvcmQK` in the aforementioned example

2. Then pass the following HTTP header along with all requests to protected endpoints: `Authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQK`
