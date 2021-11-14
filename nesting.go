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
	"strings"
	"sync"
)

//--------------------
// NESTED MULTIPLEXER
//--------------------

// NestedMux allows to nest handler following the pattern
// {prefix}/{resource}/{id}/{subresource}/{subresource-id}/...
type NestedMux struct {
	mu        sync.RWMutex
	prefix    string
	resources map[string]*NestedMux
	handler   http.Handler
}

// NewNestedMux creates an empty nested multiplexer.
func NewNestedMux(prefix string) *NestedMux {
	return &NestedMux{
		prefix:    prefix,
		resources: make(map[string]*NestedMux),
	}
}

// Handle registers the handler for the given resource name. Nested names are separated by a slash.
func (mux *NestedMux) Handle(name string, h http.Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	parts := strings.Split(name, "/")
	switch len(parts) {
	case 0:
		panic("empty resource name")
	case 1:
		submux, exists := mux.resources[parts[0]]
		if exists {
			panic("resource already exists")
		}
		submux.handler = h
	default:
		submux, exists := mux.resources[parts[0]]
		if !exists {
			panic("resource does not exist")
		}
		submux.Handle(strings.Join(parts[1:], "/"), h)
	}
}

// ServeHTTP implements http.Handler.
func (mux *NestedMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ids := ParseResourceIDs(r, mux.prefix)

	h := mux.retrieveHandler(ids)
	h.ServeHTTP(w, r)
}

// retrieveHandler retrieves the handler based on the resource IDs.
func (mux *NestedMux) retrieveHandler(ids ResourceIDs) http.Handler {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	switch len(ids) {
	case 0:
		return mux.handler
	case 1:
		submux, exists := mux.resources[ids[0].ID]
		if !exists {
			return http.NotFoundHandler()
		}
		return submux.handler
	default:
		submux, exists := mux.resources[ids[0].ID]
		if !exists {
			return http.NotFoundHandler()
		}
		return submux.retrieveHandler(ids[1:])
	}
}

// EOF
