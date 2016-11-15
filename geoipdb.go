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
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/abh/geoip"
)

// reASN is a regexp for matching against an ASN.
var reASN *regexp.Regexp

func init() {
	reASN = regexp.MustCompilePOSIX("^AS[[:digit:]]+$")
}

// Handler is a handler to TurboBytes GeoIP helper functions.
type Handler struct {
	gi *geoip.GeoIP
}

// NewHandler creates and returns a geoipdb handler.
func NewHandler() (Handler, error) {
	gi, err := geoip.OpenType(geoip.GEOIP_ASNUM_EDITION)
	if err != nil {
		return Handler{}, fmt.Errorf("cannot open GeoIP database: %s", err)
	}
	return Handler{gi: gi}, nil
}

// LibGeoipLookup queries the libgeoip database for the ASN of a given ip address.
//
// If found, returns
// an ASN identification
// and the corresponding description.
func (h Handler) LibGeoipLookup(ip string) (string, string) {
	tmp, _ := h.gi.GetName(ip)
	tmp = strings.TrimSpace(tmp)
	if tmp == "" {
		return "", ""
	}
	answer := strings.SplitN(tmp, " ", 2)
	if len(answer) < 2 {
		return answer[0], ""
	}
	return answer[0], answer[1]
}

// LookupAsn searches for the Autonomous System Number (ASN)
// of a valid IP address.
//
// This is the preferred ASN lookup function to be used by clients,
// as it queries several resources for finding proper answers.
//
// Returns
// an ASN identification
// and the corresponding description.
func (h Handler) LookupAsn(ip string) (string, string, error) {
	asnGi, asnDescr := h.LibGeoipLookup(ip)
	if asnGi != "" && asnDescr != "" {
		return asnGi, asnDescr, nil
	}
	asnDescr = ""
	asnIp, asnDescr, errIp := h.IpInfoLookup(ip)
	if errIp == nil {
		if asnIp != "" && asnDescr != "" {
			return asnIp, asnDescr, nil
		}
	} else {
		log.Println(errIp)
	}
	if errIp == nil && asnIp != "" {
		return asnIp, "", nil
	}
	if asnGi != "" {
		return asnGi, "", nil
	}
	return "", "", fmt.Errorf("unknown ASN for ip '%v'", ip)
}

// IpInfoLookup queries ipinfo.io for the ASN of a given ip address.
//
// Returns
// an ASN identification
// and the corresponding description.
func (h Handler) IpInfoLookup(ip string) (string, string, error) {
	url := fmt.Sprintf("http://ipinfo.io/%s/org", ip)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("failed to GET '%s': %s", url, err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read ipinfo.io response: %s", err)
	}
	asnData := strings.TrimSpace(string(data))
	if asnData == "" {
		return "", "", fmt.Errorf("GET '%s' returned an empty answer", url)
	}
	answer := strings.SplitN(asnData, " ", 2)
	// ipinfo.io returns errors as regular text (no outband error codes).
	// Let's try to be smart and identify them.
	if !reASN.MatchString(answer[0]) {
		return "", "", fmt.Errorf("ipinfo.io lookup failed for '%s': %s", ip, asnData)
	}
	if len(answer) < 2 {
		return answer[0], "", nil
	}
	return answer[0], answer[1], nil
}
