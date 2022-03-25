// Tideland Go HTTP Extension
//
// Copyright (C) 2020-2022 Frank Mueller / Tideland / Oldenburg / Germany
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
// RESOURCES
//--------------------

// Resource identifies a resource in a URI path by name and ID.
type Resource struct {
	Name string
	ID   string
}

// Resources is a number or resources in a URI path.
type Resources []Resource

// Path returns the number of resource names concatenated with slashes
// like they are stored in the nested multiplexer.
func (ress Resources) Path() string {
	names := make([]string, len(ress))
	for i, res := range ress {
		names[i] = res.Name
	}
	return strings.Join(names, "/")
}

// IsPath check if the resources path matches a given path.
func (ress Resources) IsPath(path string) bool {
	return ress.Path() == path
}

// PathID first checks if the resources path matches a given path and
// then returns the ID of that resource.
func (ress Resources) PathID(path string) string {
	names := strings.Split(path, "/")
	for i, name := range names {
		if ress[i].Name == name {
			return ress[i].ID
		}
	}
	return ""
}

// PathToResources parses a new Resource from a URI path.
func PathToResources(r *http.Request, prefix string) Resources {
	// Remove prefix with and without trailing slash.
	prefix = strings.TrimSuffix(prefix, "/")
	trimmed := strings.TrimPrefix(r.URL.Path, prefix)
	trimmed = strings.TrimPrefix(trimmed, "/")
	// Now split the path.
	parts := strings.Split(trimmed, "/")
	if len(parts) == 0 {
		return nil
	}
	var ress Resources
	var name string
	for i, part := range parts {
		switch {
		case i%2 == 0:
			name = part
		case i%2 == 1:
			ress = append(ress, Resource{name, part})
			name = ""
		}
	}
	if name != "" {
		ress = append(ress, Resource{name, ""})
	}
	return ress
}

// EOF
