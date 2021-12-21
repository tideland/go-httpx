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
// TESTING
//--------------------

// TestETag verifies the adding of an ETag header.
func TestETag(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	testhandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	etagwrapper := httpx.WrapETag("ABC123")
	handler := httpx.Wrap(testhandler, etagwrapper)
	s := web.NewSimulator(handler)

	// Without an If-None-Match header.
	req, err := http.NewRequest(http.MethodGet, "http://localhost:1234/", nil)
	assert.NoError(err)
	resp, err := s.Do(req)
	assert.NoError(err)
	assert.Equal(resp.StatusCode, http.StatusOK)
	assert.Equal(resp.Header.Get("ETag"), "ABC123")

	// With a non-matching If-None-Match header.
	req, err = http.NewRequest(http.MethodGet, "http://localhost:1234/", nil)
	assert.NoError(err)
	req.Header.Set("If-None-Match", "321CBA")
	resp, err = s.Do(req)
	assert.NoError(err)
	assert.Equal(resp.StatusCode, http.StatusOK)
	assert.Equal(resp.Header.Get("ETag"), "ABC123")

	// With a matching If-None-Match header.
	req, err = http.NewRequest(http.MethodGet, "http://localhost:1234/", nil)
	assert.NoError(err)
	req.Header.Set("If-None-Match", "ABC123")
	resp, err = s.Do(req)
	assert.NoError(err)
	assert.Equal(resp.StatusCode, http.StatusNotModified)
	assert.Equal(resp.Header.Get("ETag"), "ABC123")
}

// EOF
