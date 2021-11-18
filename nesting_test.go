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
	"tideland.dev/go/httpx"
)

//--------------------
// TESTS
//--------------------

// TestNestedMuxNoHandler tests the mapping of requests to a
// nested handler w/o sub-handlers.
func TestNestedMuxNoHandler(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
	defer wa.Close()

	nmux := httpx.NewNestedMux("/nmux")

	wa.Handle("/nmux", nmux)

	wreq := wa.CreateRequest(http.MethodGet, "/nmux")
	wresp := wreq.Do()

	wresp.AssertStatusCodeEquals(http.StatusNotFound)
	wresp.AssertBodyMatches("")
}

// TestNestedMux tests the mapping of requests to a number of
// nested individual handlers.
func TestNestedMux(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
	defer wa.Close()

	nmux := httpx.NewNestedMux("/api/")

	wa.Handle("/api/", nmux)

	nmux.Handle("foo", makeEchoHandler(assert, "foo"))
	nmux.Handle("foo/bar", makeEchoHandler(assert, "bar"))
	nmux.Handle("foo/baz", makeEchoHandler(assert, "baz"))

	tests := []struct {
		path       string
		statusCode int
		body       string
	}{
		{
			path:       "/",
			statusCode: http.StatusNotFound,
			body:       "",
		}, {
			path:       "/api/foo/",
			statusCode: http.StatusOK,
			body:       "foo: GET /api/foo",
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
			body:       "",
		},
	}
	for i, test := range tests {
		assert.Logf("test case #%d: %s", i, test.path)
		wreq := wa.CreateRequest(http.MethodGet, test.path)
		wresp := wreq.Do()
		wresp.AssertStatusCodeEquals(test.statusCode)
		wresp.AssertBodyMatches(test.body)
	}
}

// EOF
