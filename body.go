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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

//--------------------
// CONSTANTS
//--------------------

const (
	HeaderAccept      = "Accept"
	HeaderContentType = "Content-Type"

	ContentTypeJSON  = "application/json"
	ContentTypePlain = "text/plain"
	ContentTypeXML   = "application/xml"
)

//--------------------
// BODY HANDLING
//--------------------

// ReadBody reads and unmarshals the body of the request into the given interface. It analyzes the
// content type and uses the appropriate unmarshaler. Here it handles plain text, JSON, and XML. All
// other content types are returned directly as byte slice.
func ReadBody(r *http.Request, value interface{}) error {
	// Read content type and body.
	contentType := r.Header.Get(HeaderContentType)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = r.Body.Close(); err != nil {
		return err
	}

	// Unmarshal based on content type.
	switch contentType {
	case ContentTypeJSON:
		err := json.Unmarshal(body, value)
		if err != nil {
			return fmt.Errorf("ReadBody: cannot unmarshal body: %v", err)
		}
	case ContentTypePlain:
		pv, ok := value.(*string) // Assume v is a string pointer.
		if !ok {
			return fmt.Errorf("ReadBody: value is not a string pointer")
		}
		*pv = string(body)
	case ContentTypeXML:
		err := xml.Unmarshal(body, value)
		if err != nil {
			return fmt.Errorf("ReadBody: cannot unmarshal body: %v", err)
		}
	default:
		pbs, ok := value.(*[]byte) // Assume v is a byte slice pointer.
		if !ok {
			return fmt.Errorf("ReadBody: value is not a byte slice pointer")
		}
		*pbs = body
	}
	return nil
}

// WriteBody writes the given value to the response writer. It analyzes the content type and uses the
// the appropriate encoding. Here it handles plain text, JSON, and XML. All other content types are
// written directly as byte slice.
func WriteBody(w http.ResponseWriter, contentType string, value interface{}) (int, error) {
	// Marshal based on content type.
	switch contentType {
	case ContentTypeJSON:
		body, err := json.Marshal(value)
		if err != nil {
			return 0, fmt.Errorf("WriteBody: cannot marshal value: %v", err)
		}
		w.Header().Set(HeaderContentType, ContentTypeJSON)
		return w.Write(body)
	case ContentTypePlain:
		s, ok := value.(string)
		if !ok {
			return 0, fmt.Errorf("WriteBody: value is not a string")
		}
		w.Header().Set(HeaderContentType, ContentTypePlain)
		return w.Write([]byte(s))
	case ContentTypeXML:
		body, err := xml.Marshal(value)
		if err != nil {
			return 0, fmt.Errorf("WriteBody: cannot marshal value: %v", err)
		}
		w.Header().Set(HeaderContentType, ContentTypeXML)
		return w.Write(body)
	default:
		bs, ok := value.([]byte)
		if !ok {
			return 0, fmt.Errorf("WriteBody: value is not a byte slice")
		}
		w.Header().Set(HeaderContentType, contentType)
		return w.Write(bs)
	}
}

// EOF
