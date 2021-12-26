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
	"net/http"
	"testing"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/httpx"
)

//--------------------
// TESTS
//--------------------

// TestPathToResources tests the parsing of request paths into resource IDs.
func TestPathToResources(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	tests := []struct {
		name   string
		url    string
		prefix string
		ress   httpx.Resources
	}{
		{
			name: "no prefix",
			url:  "http://example.com/",
		}, {
			name:   "prefix, no resources",
			url:    "http://example.com/api",
			prefix: "/api",
		}, {
			name:   "prefix with trailing slash, no resources",
			url:    "http://example.com/api",
			prefix: "/api/",
		}, {
			name:   "prefix, one resource w/o id",
			url:    "http://example.com/api/users",
			prefix: "/api",
			ress:   httpx.Resources{{Name: "users"}},
		}, {
			name:   "prefix, one resource w/o id, trailing slash",
			url:    "http://example.com/api/users/",
			prefix: "/api",
			ress:   httpx.Resources{{Name: "users"}},
		}, {
			name:   "prefix, one resource w/ id",
			url:    "http://example.com/api/users/123",
			prefix: "/api",
			ress:   httpx.Resources{{Name: "users", ID: "123"}},
		}, {
			name:   "prefix, one resource w/ id, trailing slash",
			url:    "http://example.com/api/users/123/",
			prefix: "/api",
			ress:   httpx.Resources{{Name: "users", ID: "123"}},
		}, {
			name:   "prefix, two resources, second w/o id",
			url:    "http://example.com/api/users/123/contracts",
			prefix: "/api",
			ress:   httpx.Resources{{Name: "users", ID: "123"}, {Name: "contracts"}},
		}, {
			name:   "prefix, two resources w/ ids",
			url:    "http://example.com/api/users/123/contracts/1",
			prefix: "/api",
			ress:   httpx.Resources{{Name: "users", ID: "123"}, {Name: "contracts", ID: "1"}},
		},
	}

	for _, test := range tests {
		assert.Logf("test %q", test.name)
		req, err := http.NewRequest("GET", test.url, nil)
		assert.NoError(err)
		ress := httpx.PathToResources(req, test.prefix)
		assert.Length(ress, len(test.ress))
		for i, res := range ress {
			assert.Equal(res.Name, test.ress[i].Name)
			assert.Equal(res.ID, test.ress[i].ID)
		}
	}
}

// EOF
