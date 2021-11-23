// Tideland Go HTTP Extensions
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
	"sync"
)

//--------------------
// NESTED MULTIPLEXER
//--------------------

// NestedMux allows to nest handler following the RESTful API pattern
// {prefix}/{resource}/{id}/{subresource}/{subresource-id}/...
type NestedMux struct {
	mu       sync.RWMutex
	prefix   string
	handlers map[string]http.Handler
}

// NewNestedMux creates an empty nested multiplexer.
func NewNestedMux(prefix string) *NestedMux {
	return &NestedMux{
		prefix:   prefix,
		handlers: make(map[string]http.Handler),
	}
}

// Handle registers the handler for the given resource name. Nested names are separated by a slash.
func (mux *NestedMux) Handle(path string, h http.Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	mux.handlers[path] = h
}

// ServeHTTP implements http.Handler.
func (mux *NestedMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	ress := PathToResources(r, mux.prefix)
	path := ress.Path()
	h, exists := mux.handlers[path]

	if !exists {
		h = http.NotFoundHandler()
	}

	h.ServeHTTP(w, r)
}

// EOF
