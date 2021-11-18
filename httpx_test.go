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
	"fmt"
	"net/http"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/environments"
)

//--------------------
// WEB ASSERTER AND HELPERS
//--------------------

// StartTestServer initialises and starts the asserter for the tests.
func startWebAsserter(assert *asserts.Asserts) *environments.WebAsserter {
	wa := environments.NewWebAsserter(assert)
	return wa
}

// makeEchoHandler creates a handler echoing the HTTP method and the path.
func makeEchoHandler(assert *asserts.Asserts, id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reply := fmt.Sprintf("%s: %s %s", id, r.Method, r.URL.Path)
		w.Header().Add(environments.HeaderContentType, environments.ContentTypePlain)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(reply)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// EOF
