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
	"time"

	"tideland.dev/go/wait"
)

//--------------------
// THROTTLED HANDLER
//--------------------

// ThrottleHandler allows to limit the number of handled requests per second.
type ThrottledHandler struct {
	handler  http.Handler
	throttle *wait.Throttle
	timeout  time.Duration
	logger   Logger
}

// NewTimeoutThrottledHandler creates a new handler wrapping the given handler and limiting the
// number of requests per seconds to the given limit. The timeout defines the maximum time
// to wait for a request.
func NewTimeoutThrottledHandler(handler http.Handler, limit wait.Limit, timeout time.Duration, logger Logger) *ThrottledHandler {
	return &ThrottledHandler{
		handler:  handler,
		throttle: wait.NewThrottle(limit, 1),
		timeout:  timeout,
		logger:   logger,
	}
}

// NewThrottledHandler creates a new handler wrapping the given handler and limiting the
// number of requests per seconds to the given limit. There's no timeout.
func NewThrottledHandler(handler http.Handler, limit wait.Limit, logger Logger) *ThrottledHandler {
	return NewTimeoutThrottledHandler(handler, limit, 0, logger)
}

// WrapTimeoutThrottle returns a wrapper for the throttled handler with the given limit. The
// timeout defines the maximum time to wait for a request.
func WrapTimeoutThrottle(limit wait.Limit, timeout time.Duration, logger Logger) Wrapper {
	return func(h http.Handler) http.Handler {
		return NewTimeoutThrottledHandler(h, limit, timeout, logger)
	}
}

// WrapThrottle returns a wrapper for the throttled handler with the given limit. There's no
// timeout.
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
	ctx := context.Background()
	cancel := func() {}
	if h.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, h.timeout)
	}
	defer cancel()
	if err := h.throttle.Process(ctx, evt); err != nil {
		msg := fmt.Sprintf("ThrottledHandler: error during serving %s %s: %v", r.Method, r.URL.Path, err)
		h.logger.Printf(msg)
		http.Error(w, msg, http.StatusServiceUnavailable)
	}
}

// EOF
