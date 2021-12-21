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
	w.Header().Set("ETag", h.etag)
	if r.Method == http.MethodGet {
		etag := r.Header.Get("If-None-Match")
		if etag == h.etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}
	h.handler.ServeHTTP(w, r)
}

// EOF
