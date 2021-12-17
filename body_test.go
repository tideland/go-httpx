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
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/web"
	"tideland.dev/go/httpx"
)

//--------------------
// TESTS
//--------------------

// TestReadBody verifies the correct reading of different body and
// content types.
func TestReadBody(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	h := func(w http.ResponseWriter, r *http.Request) {
		var b body
		err := httpx.ReadBody(r, &b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Add(httpx.HeaderContentType, httpx.ContentTypePlain)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "1.: %q 2.: %d", b.FirstString, b.SecondInt)
	}
	s := web.NewFuncSimulator(h)

	tests := []struct {
		name        string
		contentType string
		body        io.Reader
		expected    string
		err         string
	}{
		{
			name:        "valid JSON",
			contentType: "application/json",
			body:        strings.NewReader(`{"first_string":"valid", "second_int":1}`),
			expected:    `1.: "valid" 2.: 1`,
		}, {
			name:        "invalid JSON",
			contentType: "application/json",
			body:        strings.NewReader(`{"first_string":"valid", "second_int":1)`),
			err:         "ReadBody: cannot unmarshal body",
		}, {
			name:        "valid XML",
			contentType: "application/xml",
			body:        strings.NewReader(`<root><fs>valid</fs><si>1</si></root>`),
			expected:    `1.: "valid" 2.: 1`,
		}, {
			name:        "invalid XML",
			contentType: "application/xml",
			body:        strings.NewReader(`<root><fs>valid</fs><si>1</si>`),
			err:         "ReadBody: cannot unmarshal body",
		},
	}
	for i, test := range tests {
		assert.Logf("test %d: %s", i, test.name)
		req, err := http.NewRequest(http.MethodPost, "/", test.body)
		assert.NoError(err)
		req.Header.Set("Content-Type", test.contentType)

		resp, err := s.Do(req)
		assert.NoError(err)

		if test.err != "" {
			assert.Equal(resp.StatusCode, http.StatusInternalServerError)
			body, err := web.BodyToString(resp)
			assert.NoError(err)
			assert.Contains(test.err, body)
		} else {
			assert.Equal(resp.StatusCode, http.StatusOK)
			body, err := web.BodyToString(resp)
			assert.NoError(err)
			assert.Equal(body, test.expected)
		}
	}
}

//--------------------
// HELPERS
//--------------------

// body will be marshalled and unmarshalled in the tests.
type body struct {
	FirstString string `json:"first_string" xml:"fs"`
	SecondInt   int    `json:"second_int" xml:"si"`
}

// EOF
