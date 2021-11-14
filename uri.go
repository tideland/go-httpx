// Tideland Go HTTP Extension
//
// Copyright (C) 2020-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package httpx

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"strings"
)

//--------------------
// RESOURCE IDENTIFIERS
//--------------------

// ResourceID is a type that can be used to identify a resource.
type ResourceID struct {
	Name string
	ID   string
}

// ResourceIDs is a type that can be used to identify multiple resources.
type ResourceIDs []ResourceID

// ParseResourceIDs parses a new ResourceID from a URI path.
func ParseResourceIDs(r *http.Request, prefix string) ResourceIDs {
	trimmed := strings.TrimPrefix(r.URL.Path, prefix)
	parts := strings.Split(trimmed, "/")
	if len(parts) == 0 {
		return nil
	}
	var ids ResourceIDs
	var name string
	for i, part := range parts {
		switch {
		case part == "":
			continue
		case i%2 == 0:
			name = part
		case i%2 == 1:
			ids = append(ids, ResourceID{name, part})
			name = ""
		}
	}
	if name != "" {
		ids = append(ids, ResourceID{name, ""})
	}
	return ids
}

// EOF
