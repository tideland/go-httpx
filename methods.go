// Tideland Go HTTP Extensions
//
// Copyright (C) 2020-2022 Frank Mueller / Tideland / Oldenburg / Germany
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
// METHOD HANDLER INTERFACES
//--------------------

// GetHandler has to be implemented by a handler for GET requests
// dispatched through the MethodHandler.
type GetHandler interface {
	ServeHTTPGet(w http.ResponseWriter, r *http.Request)
}

// HeadHandler has to be implemented by a handler for HEAD requests
// dispatched through the MethodHandler.
type HeadHandler interface {
	ServeHTTPHead(w http.ResponseWriter, r *http.Request)
}

// PostHandler has to be implemented by a handler for POST requests
// dispatched through the MethodHandler.
type PostHandler interface {
	ServeHTTPPost(w http.ResponseWriter, r *http.Request)
}

// PutHandler has to be implemented by a handler for PUT requests
// dispatched through the MethodHandler.
type PutHandler interface {
	ServeHTTPPut(w http.ResponseWriter, r *http.Request)
}

// PatchHandler has to be implemented by a handler for PATCH requests
// dispatched through the MethodHandler.
type PatchHandler interface {
	ServeHTTPPatch(w http.ResponseWriter, r *http.Request)
}

// DeleteHandler has to be implemented by a handler for DELETE requests
// dispatched through the MethodHandler.
type DeleteHandler interface {
	ServeHTTPDelete(w http.ResponseWriter, r *http.Request)
}

// ConnectHandler has to be implemented by a handler for CONNECT requests
// dispatched through the MethodHandler.
type ConnectHandler interface {
	ServeHTTPConnect(w http.ResponseWriter, r *http.Request)
}

// OptionsHandler has to be implemented by a handler for OPTIONS requests
// dispatched through the MethodHandler.
type OptionsHandler interface {
	ServeHTTPOptions(w http.ResponseWriter, r *http.Request)
}

// TraceHandler has to be implemented by a handler for TRACE requests
// dispatched through the MethodHandler.
type TraceHandler interface {
	ServeHTTPTrace(w http.ResponseWriter, r *http.Request)
}

//--------------------
// METHOD HANDLER
//--------------------

// MethodHandler wraps a http.Handler implementing also individual httpx handler
// interfaces. It distributes the requests to the handler methods if those are
// implemented.
type MethodHandler struct {
	handler http.Handler
}

// NewMethodHandler returns a new method handler.
func NewMethodHandler(h http.Handler) *MethodHandler {
	return &MethodHandler{
		handler: h,
	}
}

// ServeHTTP implements the http.Handler interface. If the wrapped handler implements
// the matching interface for the HTTP request method the according ServeHTTP<method>()
// method will be called. Other it simply calls the default ServeHTTP() method.
func (h *MethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if hh, ok := h.handler.(GetHandler); ok {
			hh.ServeHTTPGet(w, r)
			return
		}
	case http.MethodPost:
		if hh, ok := h.handler.(PostHandler); ok {
			hh.ServeHTTPPost(w, r)
			return
		}
	case http.MethodPut:
		if hh, ok := h.handler.(PutHandler); ok {
			hh.ServeHTTPPut(w, r)
			return
		}
	case http.MethodDelete:
		if hh, ok := h.handler.(DeleteHandler); ok {
			hh.ServeHTTPDelete(w, r)
			return
		}
	case http.MethodHead:
		if hh, ok := h.handler.(HeadHandler); ok {
			hh.ServeHTTPHead(w, r)
			return
		}
	case http.MethodPatch:
		if hh, ok := h.handler.(PatchHandler); ok {
			hh.ServeHTTPPatch(w, r)
			return
		}
	case http.MethodConnect:
		if hh, ok := h.handler.(ConnectHandler); ok {
			hh.ServeHTTPConnect(w, r)
			return
		}
	case http.MethodOptions:
		if hh, ok := h.handler.(OptionsHandler); ok {
			hh.ServeHTTPOptions(w, r)
			return
		}
	case http.MethodTrace:
		if hh, ok := h.handler.(TraceHandler); ok {
			hh.ServeHTTPTrace(w, r)
			return
		}
	}
	// Fall back to default for no matching handler method or any
	// other HTTP method.
	h.handler.ServeHTTP(w, r)
}

// EOF
