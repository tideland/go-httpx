// Tideland Go HTTP Extension - Middleware
//
// Copyright (C) 2020-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package middleware // import "tideland.dev/go/httpx/middleware"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"net/http"
)

//--------------------
// RECOVERY
//--------------------

// RecoveryHandler is able to recover from panics of wrapped handlers.
type RecoveryHandler struct {
	handler http.Handler
	logger  Logger
}

// NewRecoveryHandler creates a new handler able to recover from panics of the
// wrapped handler. An internal server error is returned to the client, the
// panic message is logged.
func NewRevoeryHandler(handler http.Handler, logger Logger) *RecoveryHandler {
	return &RecoveryHandler{
		handler: handler,
		logger:  logger,
	}
}

// WrapRecovering returns a wrapper using the recovery handler.
func WrapRecovering(logger Logger) Wrapper {
	return func(handler http.Handler) http.Handler {
		return NewRevoeryHandler(handler, logger)
	}
}

// ServeHTTP implements the http.Handler interface.
func (h *RecoveryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if hr := recover(); hr != nil {
			msg := fmt.Sprintf("RecoveryHandler: panic during serving %s %s: %v", r.Method, r.URL.Path, hr)
			h.logger.Printf(msg)
			http.Error(w, msg, http.StatusInternalServerError)
		}
	}()
	h.handler.ServeHTTP(w, r)
}

// EOF
