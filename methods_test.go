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
	"tideland.dev/go/audit/environments"
	"tideland.dev/go/httpx"
)

//--------------------
// TESTS
//--------------------

// TestMethodHandler tests the wrapping of a handler by a
// MethodHandler.
func TestMethodHandler(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
	defer wa.Close()

	wa.Handle("/mh/", httpx.NewMethodHandler(metaHandler{}))

	tests := []struct {
		method     string
		statusCode int
		body       string
	}{
		{
			method:     http.MethodGet,
			statusCode: http.StatusBadRequest,
			body:       "Bad Request",
		}, {
			method:     http.MethodHead,
			statusCode: http.StatusBadRequest,
			body:       "",
		}, {
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			body:       "Bad Request",
		}, {
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			body:       "METHOD: PUT!",
		}, {
			method:     http.MethodPatch,
			statusCode: http.StatusBadRequest,
			body:       "Bad Request",
		}, {
			method:     http.MethodDelete,
			statusCode: http.StatusNoContent,
			body:       "",
		}, {
			method:     http.MethodConnect,
			statusCode: http.StatusBadRequest,
			body:       "Bad Request",
		}, {
			method:     http.MethodOptions,
			statusCode: http.StatusBadRequest,
			body:       "Bad Request",
		}, {
			method:     http.MethodTrace,
			statusCode: http.StatusBadRequest,
			body:       "Bad Request",
		},
	}
	for i, test := range tests {
		assert.Logf("test case #%d: %s", i, test.method)
		wreq := wa.CreateRequest(test.method, "/mh/")
		wresp := wreq.Do()
		wresp.AssertStatusCodeEquals(test.statusCode)
		wresp.AssertBodyMatches(test.body)
	}
}

//--------------------
// HELPING META HANDLER
//--------------------

// metaHandler provides some of the methods for the MethodHandler
// just for testing.
type metaHandler struct{}

func (h metaHandler) ServeHTTPPut(w http.ResponseWriter, r *http.Request) {
	reply := "METHOD: " + r.Method + "!"
	w.Header().Add(environments.HeaderContentType, environments.ContentTypePlain)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(reply)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h metaHandler) ServeHTTPDelete(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNoContent)
}

func (h metaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "bad request", http.StatusBadRequest)
}

// EOF
