# Tideland Go HTTP Extensions

[![GitHub release](https://img.shields.io/github/release/tideland/go-httpx.svg)](https://github.com/tideland/go-httpx)
[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/tideland/go-httpx/master/LICENSE)
[![Go Module](https://img.shields.io/github/go-mod/go-version/tideland/go-httpx)](https://github.com/tideland/go-httpx/blob/master/go.mod)
[![GoDoc](https://godoc.org/tideland.dev/go/httpx?status.svg)](https://pkg.go.dev/mod/tideland.dev/go/httpx?tab=packages)
![Workflow](https://github.com/tideland/go-httpx/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/tideland/go-httpx)](https://goreportcard.com/report/tideland.dev/go/httpx)

## Description

**Tideland Go HTTP Extensions** contains helper functions extending the standard 
`net/http` package for the daily work with HTTP.

Most important tasks are the wrapping of standard `net/http` handlers with
non-functional tasks like logging, recovery, authorization, throttling, and
adding of CORS and ETag headers. Additionally the package contains a convenient
way to map the HTTP methods to individual handler methods, like a `POST` to the
method `ServerHTTPPost()`. In case the handler does not implement the matching 
method, the error code `http.ErrMethodNotAllowed` will be returned to the client.

For RESTful APIs also the nesting of handlers and the parsing of paths for URIs
like `/api/v1/users/{user-id}/orders/{order-id}` are supported. Additionally the
work with different content types is simplified.

## Contributors

- Frank Mueller (https://github.com/themue / https://github.com/tideland / https://tideland.dev)

