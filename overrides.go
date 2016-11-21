// Copyright (c) 2016 turbobytes
//
// This file is part of geoipdb, a library of GeoIP related helper functions
// for TurboBytes stack.
//
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package geoipdb

import (
	"errors"
	"fmt"
)

// OverridesNilCollectionError is returned by Overrides<...> methods
// when Handler was created without an overrides collection
// (see NewHandler).
var OverridesNilCollectionError = errors.New("nil overrides collection")

// OverridesAsnNotFoundError is returned by OverridesLookup
// when there is no override defined.
var OverridesAsnNotFoundError = errors.New("ASN not found")

// OverridesLookup queries the database of local overrides
// for the description of a given ASN.
//
// Returns the ASN description.
func (h Handler) OverridesLookup(asn string) (string, error) {
	if h.overrides == nil {
		return "", OverridesNilCollectionError
	}
	return "", fmt.Errorf("not yet implemented")
}

// OverridesSet stores a user defined description for a given ASN
// in the database of local overrides.
func (h Handler) OverridesSet(asn string, descr string) error {
	if h.overrides == nil {
		return OverridesNilCollectionError
	}
	return fmt.Errorf("not yet implemented")
}

// OverridesRemove removes the description for a given ASN
// from the database of local overrides.
func (h Handler) OverridesRemove(asn string) error {
	if h.overrides == nil {
		return OverridesNilCollectionError
	}
	return fmt.Errorf("not yet implemented")
}
