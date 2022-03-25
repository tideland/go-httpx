// Tideland Go HTTP Extensions - Middleware - Unit Tests
//
// Copyright (C) 2020-2022 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package middleware_test // import "tideland.dev/go/httpx/middleware"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/web"

	"tideland.dev/go/httpx/middleware"
)

//--------------------
// TESTING
//--------------------

// TestCORS verifies the adding of CORS headers.
func TestCORS(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	corsHeaders := middleware.CORSHeaders{
		AllowOrigin:   "http://testhost/",
		AllowMethods:  []string{http.MethodGet, http.MethodPost},
		AllowHeaders:  []string{"X-Test-Allow-Header"},
		ExposeHeaders: []string{"X-Test-Expose-Header"},
		MaxAge:        30 * time.Minute,
	}

	testhandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "content")
	})
	corswrapper := middleware.WrapCORS(corsHeaders)
	handler := middleware.Wrap(testhandler, corswrapper)
	s := web.NewSimulator(handler)

	// Step A: OPTIONS request.
	req := s.CreateRequest(http.MethodOptions, "/", nil)
	resp, err := s.Do(req)
	assert.NoError(err)
	assert.Equal(resp.Header.Get(middleware.HeaderAllowOrigin), "http://testhost/")
	assert.Equal(resp.Header.Get(middleware.HeaderAllowMethods), "GET, POST")
	assert.Equal(resp.Header.Get(middleware.HeaderAllowHeaders), "X-Test-Allow-Header")
	assert.Equal(resp.Header.Get(middleware.HeaderExposeHeaders), "X-Test-Expose-Header")
	assert.Equal(resp.Header.Get(middleware.HeaderMaxAge), "1800")
	body, err := web.BodyToString(resp)
	assert.NoError(err)
	assert.Empty(body)

	// Step B: GET request.
	resp, err = s.Get("/")
	assert.NoError(err)
	assert.Equal(resp.Header.Get(middleware.HeaderAllowOrigin), "http://testhost/")
	assert.Equal(resp.Header.Get(middleware.HeaderAllowMethods), "GET, POST")
	assert.Equal(resp.Header.Get(middleware.HeaderAllowHeaders), "X-Test-Allow-Header")
	assert.Equal(resp.Header.Get(middleware.HeaderExposeHeaders), "X-Test-Expose-Header")
	assert.Equal(resp.Header.Get(middleware.HeaderMaxAge), "1800")
	body, err = web.BodyToString(resp)
	assert.NoError(err)
	assert.Equal(body, "content")
}

// EOF
