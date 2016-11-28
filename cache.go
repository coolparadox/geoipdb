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
	"time"
)

// cachedData is the data we want to keep cached.
type cachedData struct {
	// ASN number
	asn string
	// ASN description
	descr string
	// Due date of this cached information
	due time.Time
}

// cache allows manipulating cached data.
type cache struct {
	// IP to ASN data
	ip map[string]cachedData
	// ASN to IP list
	asn map[string][]string
}

// newCache returns an empty initialized cache.
func newCache() cache {
	return cache{
		ip:  make(map[string]cachedData),
		asn: make(map[string][]string),
	}
}

// store stores data to the cache.
func (c cache) store(ip string, asn string, descr string) {
}

// lookupByIP retrieves cached data by IP address.
//
// Returns
// the ASN identification and description,
// if cached data is expired,
// and if ip was found in cache.
func (c cache) lookupByIP(ip string) (asn string, descr string, expired bool, found bool) {
	return "", "", false, false
}

// lookupByASN retrieves the list of cached IPs associated with a given ASN.
//
// Returns a non nil list of IP addresses.
func (c cache) lookupByASN(asn string) []string {
	return make([]string, 0)
}

// purgeIP removes from cache all information related to a given IP.
func (c cache) purgeIP(ip string) {
}

// purgeASN removes from the cache all information related to a given ASN.
func (c cache) purgeASN(asn string) {
}
