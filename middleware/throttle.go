// Tideland Go HTTP Extensions - Middleware
//
// Copyright (C) 2020-2022 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package middleware // import "tideland.dev/go/httpx/middleware"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"net/http"

	"tideland.dev/go/wait"
)

//--------------------
// THROTTLED HANDLER
//--------------------

// ThrottleHandler allows to limit the number of handled requests per second.
type ThrottledHandler struct {
	handler  http.Handler
	throttle *wait.Throttle
	logger   Logger
}

// NewThrottledHandler create a new handler wrapping the given handler and requesting the
// number of requests per seconds to the given limit.
func NewThrottledHandler(handler http.Handler, limit wait.Limit, logger Logger) *ThrottledHandler {
	return &ThrottledHandler{
		handler:  handler,
		throttle: wait.NewThrottle(limit, 1),
		logger:   logger,
	}
}

// WrapThrottle returns a wrapper for the throttled handler with the given limit.
func WrapThrottle(limit wait.Limit, logger Logger) Wrapper {
	return func(h http.Handler) http.Handler {
		return NewThrottledHandler(h, limit, logger)
	}
}

// ServeHTTP implements http.Handler.
func (h *ThrottledHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	evt := func() error {
		h.handler.ServeHTTP(w, r)
		return nil
	}
	if err := h.throttle.Process(context.Background(), evt); err != nil {
		msg := fmt.Sprintf("ThrottledHandler: error during serving %s %s: %v", r.Method, r.URL.Path, err)
		h.logger.Printf(msg)
		http.Error(w, msg, http.StatusServiceUnavailable)
	}
}

// EOF
