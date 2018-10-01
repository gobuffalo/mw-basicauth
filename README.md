<p align="center"><img src="https://github.com/gobuffalo/buffalo/blob/master/logo.svg" width="360"></p>

<p align="center">
  <a href="https://godoc.org/github.com/gobuffalo/mw-basicauth"><img src="https://godoc.org/github.com/gobuffalo/mw-basicauth?status.svg" alt="GoDoc"></a>
  <a href="https://travis-ci.org/gobuffalo/mw-basicauth"><img src="https://travis-ci.org/gobuffalo/mw-basicauth.svg?branch=master" alt="Build Status"></a>
  <a href="https://goreportcard.com/report/github.com/gobuffalo/mw-basicauth"><img src="https://goreportcard.com/badge/github.com/gobuffalo/mw-basicauth" alt="Go Report Card" /></a>
</p>

# [Basic HTTP Authentication](https://tools.ietf.org/html/rfc7617) Middleware for [Buffalo](https://github.com/gobuffalo/buffalo)

## Installation

```bash
$ go get -u github.com/gobuffalo/mw-basicauth
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

2. Then pass the following HTTP header along with all requests to protected endpoints: `Authorization: Basic: dXNlcm5hbWU6cGFzc3dvcmQK`
