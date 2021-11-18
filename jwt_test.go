// Tideland Go HTTP Extensions - Unit Tests
//
// Copyright (C) 2020-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package httpx_test // import "tideland.dev/go/httpx"

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"net/http"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/environments"
	"tideland.dev/go/httpx"
	"tideland.dev/go/jwt"
)

//--------------------
// TESTS
//--------------------

// TestJWTHandler tests access control by usage of the JWTHandler.
func TestJWTHandler(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
	defer wa.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add(environments.HeaderContentType, environments.ContentTypePlain)
		_, err := w.Write([]byte("request passed"))
		assert.NoError(err)
	})
	jwtWrapper := httpx.NewJWTHandler(handler, &httpx.JWTHandlerConfig{
		Key: []byte("secret"),
		Gatekeeper: func(w http.ResponseWriter, r *http.Request, claims jwt.Claims) error {
			access, ok := claims.GetString("access")
			if !ok || access != "allowed" {
				return errors.New("access is not allowed")
			}
			return nil
		},
	})

	wa.Handle("/", jwtWrapper)

	tests := []struct {
		key         string
		accessClaim string
		statusCode  int
		body        string
	}{
		{
			key:         "",
			accessClaim: "",
			statusCode:  http.StatusUnauthorized,
			body:        "request contains no authorization header",
		}, {
			key:         "unknown",
			accessClaim: "allowed",
			statusCode:  http.StatusUnauthorized,
			body:        "cannot verify the signature",
		}, {
			key:         "secret",
			accessClaim: "allowed",
			statusCode:  http.StatusOK,
			body:        "request passed",
		}, {
			key:         "unknown",
			accessClaim: "forbidden",
			statusCode:  http.StatusUnauthorized,
			body:        "cannot verify the signature",
		}, {
			key:         "secret",
			accessClaim: "forbidden",
			statusCode:  http.StatusUnauthorized,
			body:        "access rejected by gatekeeper: access is not allowed",
		},
	}
	for i, test := range tests {
		assert.Logf("test case #%d: %s / %s", i, test.key, test.accessClaim)
		wreq := wa.CreateRequest(http.MethodGet, "/")
		wreq.Header().Set(httpx.HeaderAccept, httpx.ContentTypeJSON)
		if test.key != "" && test.accessClaim != "" {
			claims := jwt.NewClaims()
			claims.Set("access", test.accessClaim)
			jwt, err := jwt.Encode(claims, []byte(test.key), jwt.HS512)
			assert.NoError(err)
			wreq.Header().Set("Authorization", "Bearer "+jwt.String())
		}
		wresp := wreq.Do()
		wresp.AssertStatusCodeEquals(test.statusCode)
		wresp.AssertBodyMatches(test.body)
	}
}

// EOF
