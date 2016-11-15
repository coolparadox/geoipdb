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

/*
Package geoipdb is a library of GeoIP related helper functions for TurboBytes
stack.
*/

package geoipdb

import (
	"fmt"

	"github.com/abh/geoip"
)

// Handler is a handler to TurboBytes GeoIP helper functions.
type Handler struct {
	gip *geoip.GeoIP
}

// NewHandler creates and returns a geoipdb handler.
func NewHandler() (Handler, error) {
	gip, err := geoip.OpenType(geoip.GEOIP_ASNUM_EDITION)
	if err != nil {
		return Handler{}, fmt.Errorf("cannot open GeoIP database: %s", err)
	}
	return Handler{gip: gip}, nil
}

// LookupAsn answers the Autonomous System Number (ASN) of a valid IP address.
//
// Returns
// an ASN identification (eg "AS15169")
// and the corresponding network name (eg. "Google Inc.").
func (h Handler) LookupAsn(ip string) (string, string, error) {
	return "", "", fmt.Errorf("not yet implemented")
}
