// Tideland Go HTTP Extensions - Unit Tests
//
// Copyright (C) 2020-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package httpx_test // import "tideland.dev/go/httpx"

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/web"

	"tideland.dev/go/httpx"
)

//--------------------
// TESTS
//--------------------

// TestNestedMuxNoHandler tests the mapping of requests to a
// nested handler w/o sub-handlers.
func TestNestedMuxNoHandler(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	nmux := httpx.NewNestedMux("/")

	s := web.NewSimulator(nmux)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	resp, err := s.Do(req)
	assert.NoError(err)
	assert.Equal(resp.StatusCode, http.StatusNotFound)
}

// TestNestedMux tests the mapping of requests to a number of
// nested individual handlers.
func TestNestedMux(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	nmux := httpx.NewNestedMux("/api/")

	nmux.Handle("foo", makeEchoHandler(assert, "foo"))
	nmux.Handle("foo/bar", makeEchoHandler(assert, "bar"))
	nmux.Handle("foo/baz", makeEchoHandler(assert, "baz"))

	s := web.NewSimulator(nmux)

	tests := []struct {
		path       string
		statusCode int
		body       string
	}{
		{
			path:       "/",
			statusCode: http.StatusNotFound,
			body:       "404 page not found\n",
		}, {
			path:       "/api/foo/",
			statusCode: http.StatusOK,
			body:       "foo: GET /api/foo/",
		}, {
			path:       "/api/foo/4711",
			statusCode: http.StatusOK,
			body:       "foo: GET /api/foo/4711",
		}, {
			path:       "/api/foo/4711/bar",
			statusCode: http.StatusOK,
			body:       "bar: GET /api/foo/4711/bar",
		}, {
			path:       "/api/foo/4711/bar/1",
			statusCode: http.StatusOK,
			body:       "bar: GET /api/foo/4711/bar/1",
		}, {
			path:       "/api/foo/4711/baz",
			statusCode: http.StatusOK,
			body:       "baz: GET /api/foo/4711/baz",
		}, {
			path:       "/api/foo/4711/baz/2",
			statusCode: http.StatusOK,
			body:       "baz: GET /api/foo/4711/baz/2",
		}, {
			path:       "/api/foo/4711/bar/1/nothingelse",
			statusCode: http.StatusNotFound,
			body:       "404 page not found\n",
		},
	}
	for i, test := range tests {
		assert.Logf("test case #%d: %s", i, test.path)
		resp, err := s.Get(test.path)
		assert.NoError(err)
		assert.Equal(resp.StatusCode, test.statusCode)
		body, err := web.BodyToString(resp)
		assert.NoError(err)
		assert.Contains(test.body, body)
	}
}

// EOF
