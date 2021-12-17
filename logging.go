// Tideland Go HTTP Extension
//
// Copyright (C) 2020-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package httpx // import "tideland.dev/go/httpx"

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
)

//--------------------
// IMPORTS
//--------------------

// Logger defines an interface for many different loggers. One of it
// is the standard logger log.Logger.
type Logger interface {
	Printf(format string, v ...interface{})
}

// LoggingHandler wraps a handler and logs the requests to it.
type LoggingHandler struct {
	handler http.Handler
	logger  Logger
}

// NewLoggingHandler creates a new logging handler with the given logger and handler.
func NewLoggingHandler(handler http.Handler, logger Logger) *LoggingHandler {
	return &LoggingHandler{
		handler: handler,
		logger:  logger,
	}
}

// WrapLogging returns a wrapper for the logging handler with the given logger.
func WrapLogging(logger Logger) Wrapper {
	return func(h http.Handler) http.Handler {
		return NewLoggingHandler(h, logger)
	}
}

// ServeHTTP logs the request and calls the wrapped handler.
func (h *LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("%s %s", r.Method, r.URL.Path)
	h.handler.ServeHTTP(w, r)
}

// EOF
