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
	"log"
	"net/http"
	"strconv"
	"time"

	"tideland.dev/go/jwt"
)

//--------------------
// JWT HANDLER
//--------------------

// JWTHandlerConfig allows to control how the JWT handler works.
// All values are optional. In this case tokens are only decoded
// without using a cache, validated for the current time plus/minus
// a minute leeway, and there's no user defined gatekeeper function
// running afterwards.
type JWTHandlerConfig struct {
	Cache      *jwt.Cache
	Key        jwt.Key
	Leeway     time.Duration
	Gatekeeper func(w http.ResponseWriter, r *http.Request, claims jwt.Claims) error
	logger     Logger
}

// JWTHandler checks for a valid token and then runs
// a gatekeeper function.
type JWTHandler struct {
	handler    http.Handler
	cache      *jwt.Cache
	key        jwt.Key
	leeway     time.Duration
	gatekeeper func(w http.ResponseWriter, r *http.Request, claims jwt.Claims) error
	logger     Logger
}

// NewJWTHandler creates a handler checking for a valid JSON
// Web Token in each request.
func NewJWTHandler(handler http.Handler, config *JWTHandlerConfig) *JWTHandler {
	h := &JWTHandler{
		handler: handler,
		leeway:  time.Minute,
		logger:  log.Default(),
	}
	if config != nil {
		if config.Cache != nil {
			h.cache = config.Cache
		}
		if config.Key != nil {
			h.key = config.Key
		}
		if config.Leeway != 0 {
			h.leeway = config.Leeway
		}
		if config.Gatekeeper != nil {
			h.gatekeeper = config.Gatekeeper
		}
		if config.logger != nil {
			h.logger = config.logger
		}
	}
	return h
}

// WrapJWT returns a wrapper for the JWT handler with the given configuration.
func WrapJWT(config *JWTHandlerConfig) Wrapper {
	return func(h http.Handler) http.Handler {
		return NewJWTHandler(h, config)
	}
}

// ServeHTTP implements the http.Handler interface. It checks for an existing
// and valid token before calling the wrapped handler.
func (h *JWTHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.isAuthorized(w, r) {
		h.handler.ServeHTTP(w, r)
	}
}

// isAuthorized checks the request for a valid token and if configured
// asks the gatekeepr if the request may pass.
func (h *JWTHandler) isAuthorized(w http.ResponseWriter, r *http.Request) bool {
	var token *jwt.JWT
	var err error
	switch {
	case h.cache != nil && h.key != nil:
		token, err = h.cache.RequestVerify(r, h.key)
	case h.cache != nil && h.key == nil:
		token, err = h.cache.RequestDecode(r)
	case h.cache == nil && h.key != nil:
		token, err = jwt.RequestVerify(r, h.key)
	default:
		token, err = jwt.RequestDecode(r)
	}
	// Now do the checks.
	if err != nil {
		h.deny(w, r, err.Error(), http.StatusUnauthorized)
		return false
	}
	if token == nil {
		h.deny(w, r, "no JSON Web Token", http.StatusUnauthorized)
		return false
	}
	if !token.IsValid(h.leeway) {
		h.deny(w, r, "the JSON Web Token claims 'nbf' and/or 'exp' are not valid", http.StatusForbidden)
		return false
	}
	if h.gatekeeper != nil {
		err := h.gatekeeper(w, r, token.Claims())
		if err != nil {
			h.deny(w, r, "access rejected by gatekeeper: "+err.Error(), http.StatusUnauthorized)
			return false
		}
	}
	// All fine.
	return true
}

// deny sends a negative feedback to the caller.
func (h *JWTHandler) deny(w http.ResponseWriter, r *http.Request, msg string, statusCode int) {
	feedback := map[string]string{
		"statusCode": strconv.Itoa(statusCode),
		"message":    msg,
	}
	accept := r.Header.Get(HeaderAccept)
	w.WriteHeader(statusCode)
	_, err := WriteBody(w, accept, feedback)
	if err != nil {
		h.logger.Printf("JWT handler: %v", err)
	}
}

// EOF
