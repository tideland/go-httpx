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
	"log"
	"net/http"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/web"
	"tideland.dev/go/wait"

	"tideland.dev/go/httpx/middleware"
)

//--------------------
// TESTING
//--------------------

// TestThrottle verifies the throttling of requests per second.
func TestThrottle(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	testhandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	limit := wait.Limit(5)
	// Define tests.
	tests := []struct {
		name            string
		timeout         time.Duration
		rps             int
		reqs            int
		respOK          int
		respUnavailable int
		seconds         time.Duration
	}{
		{
			name:    "less-requests-than-limit",
			rps:     3,
			reqs:    9,
			respOK:  9,
			seconds: 3 * time.Second,
		},
		{
			name:    "requests-matching-limit",
			rps:     5,
			reqs:    15,
			respOK:  15,
			seconds: 3 * time.Second,
		},
		{
			name:    "more-requests-than-limit",
			rps:     20,
			reqs:    25,
			respOK:  25,
			seconds: 5 * time.Second,
		},
		{
			name:            "short-timeout",
			timeout:         20 * time.Millisecond,
			rps:             20,
			reqs:            60,
			respOK:          15,
			respUnavailable: 45,
			seconds:         3 * time.Second,
		},
		{
			name:    "long-timeout",
			timeout: time.Second,
			rps:     20,
			reqs:    25,
			respOK:  25,
			seconds: 5 * time.Second,
		},
	}
	// Run tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.SetFailable(t)
			throttlewrapper := middleware.WrapTimeoutThrottle(limit, test.timeout, log.Default())
			handler := middleware.Wrap(testhandler, throttlewrapper)
			sim := web.NewSimulator(handler)
			sleep := time.Second / time.Duration(test.rps)
			respOK := 0
			respUnavailable := 0
			begin := time.Now()
			for i := 0; i < test.reqs; i++ {
				resp, err := sim.Get("http://localhost:1234/")
				assert.NoError(err)
				switch resp.StatusCode {
				case http.StatusOK:
					respOK++
				case http.StatusServiceUnavailable:
					respUnavailable++
				default:
					assert.Fail(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
				}
				time.Sleep(sleep)
			}
			duration := time.Since(begin)
			assert.Equal(respOK, test.respOK, "status ok")
			assert.Equal(respUnavailable, test.respUnavailable, "status unavailable")
			assert.About(duration.Seconds(), test.seconds.Seconds(), 0.25, "duration")
		})
	}
}

// EOF
