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
func ReadBody(r *http.Request, v interface{}) error {
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
		return json.Unmarshal(body, v)
	case ContentTypePlain:
		pv, ok := v.(*string) // Assume v is a string pointer.
		if !ok {
			return fmt.Errorf("ReadBody: v is not a string pointer")
		}
		*pv = string(body)
	case ContentTypeXML:
		return xml.Unmarshal(body, v)
	default:
		pbs, ok := v.(*[]byte) // Assume v is a byte slice pointer.
		if !ok {
			return fmt.Errorf("ReadBody: v is not a byte slice pointer")
		}
		*pbs = body
	}
	return nil
}

// WriteBody writes the given value to the response writer. It analyzes the content type and uses the
// the appropriate encoding. Here it handles plain text, JSON, and XML. All other content types are
// written directly as byte slice.
func WriteBody(w http.ResponseWriter, contentType string, v interface{}) error {
	// Marshal based on content type.
	switch contentType {
	case ContentTypeJSON:
		body, err := json.Marshal(v)
		if err != nil {
			return err
		}
		w.Header().Set(HeaderContentType, ContentTypeJSON)
		w.Write(body)
	case ContentTypePlain:
		s, ok := v.(string)
		if !ok {
			return fmt.Errorf("WriteBody: v is not a string")
		}
		w.Header().Set(HeaderContentType, ContentTypePlain)
		w.Write([]byte(s))
	case ContentTypeXML:
		body, err := xml.Marshal(v)
		if err != nil {
			return err
		}
		w.Header().Set(HeaderContentType, ContentTypeXML)
		w.Write(body)
	default:
		bs, ok := v.([]byte)
		if !ok {
			return fmt.Errorf("WriteBody: v is not a byte slice")
		}
		w.Header().Set(HeaderContentType, contentType)
		w.Write(bs)
	}
	return nil
}

// EOF
