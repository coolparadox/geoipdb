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
	"sync"
)

// cacheTTL is the expiration time of a cache entry.
const cacheTTL = time.Hour * 24

// cacheEntry is the data we want to keep cached.
type cacheEntry struct {
	// ASN number
	asn string
	// ASN description
	descr string
	// Due date of this entry
	due time.Time
}

// cache allows manipulating cached data.
type cache struct {
	// Concurrent access control to maps
	sync.RWMutex
	// IP to ASN data
	ip map[string]cacheEntry
	// ASN to IP list
	asn map[string]map[string]interface{}
}

// newCache returns an empty initialized cache.
func newCache() cache {
	return cache{
		ip:  make(map[string]cacheEntry),
		asn: make(map[string]map[string]interface{}),
	}
}

// store updates the cache.
func (c cache) store(ip string, asn string, descr string) {
	c.Lock()
	defer c.Unlock()
	// Purge ASN map of given ip
	for _, ips := range c.asn {
		delete(ips, ip)
	}
	// Purge ASN map of empty entries
	for asn, ips := range c.asn {
		if len(ips) < 1 {
			delete(c.asn, asn)
		}
	}
	// Update IP map
	c.ip[ip] = cacheEntry{
		asn:   asn,
		descr: descr,
		due:   time.Now().Add(cacheTTL),
	}
	// Update ASN map
	if c.asn[asn] == nil {
		c.asn[asn] = make(map[string]interface{})
	}
	c.asn[asn][ip] = nil
}

// lookupByIP retrieves cached data by IP address.
//
// Returns
// the ASN identification and description,
// if cached data is expired,
// and if ip was found in cache.
func (c cache) lookupByIP(ip string) (asn string, descr string, expired bool, found bool) {
	c.RLock()
	defer c.RUnlock()
	entry, ok := c.ip[ip]
	if !ok {
		return "", "", false, false
	}
	return entry.asn, entry.descr, time.Now().After(entry.due), true
}

// lookupByASN retrieves the list of cached IPs associated with a given ASN.
//
// Returns a non nil list of IP addresses.
func (c cache) lookupByASN(asn string) map[string]interface{} {
	c.RLock()
	defer c.RUnlock()
	answer, ok := c.asn[asn]
	if !ok || answer == nil {
		return make(map[string]interface{})
	}
	return answer
}

// purgeASN removes from the cache all information related to a given ASN.
func (c cache) purgeASN(asn string) {
	c.Lock()
	defer c.Unlock()
	// Purge ip map of given asn
	for ip, entry := range c.ip {
		if entry.asn == asn {
			delete(c.ip, ip)
		}
	}
	// Purge asn map of given asn
	delete(c.asn, asn)
}

// purgeAll removes all entries from the cache
func (c cache) purgeAll() {
	c.Lock()
	defer c.Unlock()
	for ip, _ := range c.ip {
		delete(c.ip, ip)
	}
	for asn, _ := range c.asn {
		delete(c.asn, asn)
	}
}

// asnList retrieves all ASNs known to the cache.
//
// Returns a non nil list of ASNs.
func (c cache) asnList() []string {
	c.RLock()
	defer c.RUnlock()
	answer := make([]string, len(c.asn))
	var i int
	for asn := range c.asn {
		answer[i] = asn
		i++
	}
	return answer
}
