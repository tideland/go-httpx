// Tideland Go HTTP Extensions - Middleware - Unit Tests
//
// Copyright (C) 2020-2021 Frank Mueller / Tideland / Oldenburg / Germany
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

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/web"

	"tideland.dev/go/httpx/middleware"
)

//--------------------
// TESTS
//--------------------

// TestLoggingHandler tests wrapping with the logging handler.
func TestLoggingHandler(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	logger := &bufferedLogger{}
	testhandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	logwrapper := middleware.WrapLogging(logger)
	handler := middleware.Wrap(testhandler, logwrapper)
	s := web.NewSimulator(handler)

	for i := 0; i < 5; i++ {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:12345/", nil)
		assert.NoError(err)
		resp, err := s.Do(req)
		assert.NoError(err)
		assert.Equal(resp.StatusCode, http.StatusOK)
	}

	assert.Length(logger.lines, 5)
	for _, line := range logger.lines {
		assert.Equal(line, "GET /")
	}
}

//--------------------
// HELPER
//--------------------

// bufferedLogger simply collects the logged lines.
type bufferedLogger struct {
	lines []string
}

// Printf implements the logger interface.
func (l *bufferedLogger) Printf(format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	l.lines = append(l.lines, line)

	fmt.Println(line)
}

// EOF
