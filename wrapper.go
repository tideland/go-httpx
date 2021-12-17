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

import "net/http"

//--------------------
// WRAPPER
//--------------------

// Wrapper defines a function wrapping a handler with another handler.
type Wrapper func(h http.Handler) http.Handler

// Wrap wraps the given handler with all listed wrappers. So it returns
// a stack of handlers able to pre- and post-process the request and response.
func Wrap(h http.Handler, wrappers ...Wrapper) http.Handler {
	for _, wrap := range wrappers {
		h = wrap(h)
	}
	return h
}

// EOF