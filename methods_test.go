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

// TestMethodHandler tests the wrapping of a handler by a
// MethodHandler.
func TestMethodHandler(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	s := web.NewSimulator(httpx.NewMethodHandler(metaHandler{}))

	tests := []struct {
		method     string
		statusCode int
		body       string
	}{
		{
			method:     http.MethodGet,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodHead,
			statusCode: http.StatusBadRequest,
			body:       "",
		}, {
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			body:       "METHOD: PUT!",
		}, {
			method:     http.MethodPatch,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodDelete,
			statusCode: http.StatusNoContent,
			body:       "",
		}, {
			method:     http.MethodConnect,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodOptions,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodTrace,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		},
	}
	for i, test := range tests {
		assert.Logf("test case #%d: %s", i, test.method)
		req, err := http.NewRequest(test.method, "/", nil)
		assert.NoError(err)
		resp, err := s.Do(req)
		assert.NoError(err)
		assert.Equal(resp.StatusCode, test.statusCode)
		body, err := web.BodyToString(resp)
		assert.NoError(err)
		assert.Contains(test.body, body)
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
	w.Header().Add(httpx.HeaderContentType, httpx.ContentTypePlain)
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
