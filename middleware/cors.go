// Tideland Go HTTP Extension - Middleware
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
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
)

//--------------------
// CONSTANTS
//--------------------

const (
	HeaderAllowOrigin   = "Access-Control-Allow-Origin"
	HeaderAllowMethods  = "Access-Control-Allow-Methods"
	HeaderAllowHeaders  = "Access-Control-Allow-Headers"
	HeaderExposeHeaders = "Access-Control-Expose-Headers"
	HeaderMaxAge        = "Access-Control-Max-Age"
)

//--------------------
// CORS
//--------------------

// CORSHeaders contains and adds the configured headers to a response writer.
type CORSHeaders struct {
	AllowOrigin   string
	AllowMethods  []string
	AllowHeaders  []string
	ExposeHeaders []string
	MaxAge        time.Duration
}

// add adds the configured headers to a response writer.
func (h CORSHeaders) add(w http.ResponseWriter) {
	if h.AllowOrigin != "" {
		w.Header().Add(HeaderAllowOrigin, h.AllowOrigin)
	}
	if len(h.AllowMethods) > 0 {
		meths := strings.Join(h.AllowMethods, ", ")
		w.Header().Add(HeaderAllowMethods, meths)
	}
	if len(h.AllowHeaders) > 0 {
		hdrs := strings.Join(h.AllowHeaders, ", ")
		w.Header().Add(HeaderAllowHeaders, hdrs)
	}
	if len(h.ExposeHeaders) > 0 {
		hdrs := strings.Join(h.ExposeHeaders, ", ")
		w.Header().Add(HeaderExposeHeaders, hdrs)
	}
	if h.MaxAge.Seconds() > 0 {
		sec := int(math.Round(h.MaxAge.Seconds()))
		w.Header().Add(HeaderMaxAge, fmt.Sprintf("%d", sec))
	}

}

// CORSHandler adds cross-origin resource sharing headers for the wrapped handler.
type CORSHandler struct {
	handler http.Handler
	headers CORSHeaders
}

// NewCORSHandler creates a new CORSHandler wrapping the given handler and using
// the given header values.
func NewCORSHandler(handler http.Handler, headers CORSHeaders) *CORSHandler {
	return &CORSHandler{
		handler: handler,
		headers: headers,
	}
}

// WrapCORS returns a wrapper using the CORS handler.
func WrapCORS(headers CORSHeaders) Wrapper {
	return func(handler http.Handler) http.Handler {
		return NewCORSHandler(handler, headers)
	}
}

// ServeHTTP implements http.Server interface.
func (h *CORSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.headers.add(w)

	if r.Method == http.MethodOptions {
		// No further content in case of OPTIONS request.
		w.WriteHeader(http.StatusNoContent)
		return
	}

	h.handler.ServeHTTP(w, r)
}

// EOF
