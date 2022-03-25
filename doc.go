// Tideland Go HTTP Extensions
//
// Copyright (C) 2020-2022 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package httpx contains helper functions extending the standard
// net/http package for the daily work with HTTP.
//
// Most important tasks are the wrapping of standard net/http handlers with
// non-functional tasks like logging, recovery, authorization, throttling, and
// adding of CORS and ETag headers. Additionally the package contains a convenient
// way to map the HTTP methods to individual handler methods, like a POST to the
// method ServerHTTPPost(). In case the handler does not implement the matching
// method, the error code `http.ErrMethodNotAllowed` will be returned to the client.
//
// For RESTful APIs also the nesting of handlers and the parsing of paths for URIs
// like /api/v1/users/{user-id}/orders/{order-id} are supported. Additionally the
// work with different content types is simplified.
package httpx // import "tideland.dev/go/httpx"

// EOF
