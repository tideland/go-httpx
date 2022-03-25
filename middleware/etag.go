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
	"net/http"
)

//--------------------
// CONSTANTS
//--------------------

const (
	HeaderETag        = "etag"
	HeaderIfNoneMatch = "If-None-Match"
)

//--------------------
// ETAG
//--------------------

// ETagHandler adds an ETag header for client caching. Additionally it checks
// if the client already has the latest version of the resource. In this case
// code 304 ("not modified") is returned.
type ETagHandler struct {
	handler http.Handler
	etag    string
}

// NewETagHandler creates a new handler able to add an ETag header for client
// caching. Additionally it checks if the client already has the latest version
// of the resource.
func NewETagHandler(handler http.Handler, etag string) *ETagHandler {
	return &ETagHandler{
		handler: handler,
		etag:    etag,
	}
}

// WrapETag returns a wrapper using the ETag handler.
func WrapETag(etag string) Wrapper {
	return func(handler http.Handler) http.Handler {
		return NewETagHandler(handler, etag)
	}
}

// ServeHTTP implements the http.Handler interface.
func (h *ETagHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(HeaderETag, h.etag)
	if r.Method == http.MethodGet {
		etag := r.Header.Get(HeaderIfNoneMatch)
		if etag == h.etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}
	h.handler.ServeHTTP(w, r)
}

// EOF
