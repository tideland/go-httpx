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
	"log"
	"net/http"
	"os"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/web"

	"tideland.dev/go/httpx"
)

//--------------------
// TESTING
//--------------------

// TestRecoveryHandler tests recovering from panics.
func TestRecoveryHandler(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	testhandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/panic/" {
			panic("ouch, a panic")
		}
		w.WriteHeader(http.StatusOK)
	})
	recoveringwrapper := httpx.WrapRecovering(log.New(os.Stdout, "[test] ", log.LstdFlags))
	handler := httpx.Wrap(testhandler, recoveringwrapper)
	s := web.NewSimulator(handler)

	// First a non-panic request.
	req, err := http.NewRequest(http.MethodGet, "http://localhost:1234/all-fine/", nil)
	assert.NoError(err)
	resp, err := s.Do(req)
	assert.NoError(err)
	assert.Equal(resp.StatusCode, http.StatusOK)

	// And now one with a panic.
	req, err = http.NewRequest(http.MethodGet, "http://localhost:1234/panic/", nil)
	assert.NoError(err)
	resp, err = s.Do(req)
	assert.NoError(err)
	assert.Equal(resp.StatusCode, http.StatusInternalServerError)
	body, err := web.BodyToString(resp)
	assert.NoError(err)
	assert.Equal(body, "RecoveryHandler: panic during serving GET /panic/: ouch, a panic\n")
}

// EOF
