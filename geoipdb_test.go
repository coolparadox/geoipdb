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

package geoipdb_test

import (
	"testing"

	"github.com/turbobytes/geoipdb"
)

// Well known IP address for testing lookups
const ip = "8.8.8.8"

// Results of some ASN lookups
var (
	asnLibGeo string
	asnIpInfo string
)

func TestInitIp(t *testing.T) {
	t.Logf("using ip '%s' for tests", ip)
}

var gh geoipdb.Handler

func TestCreateHandler(t *testing.T) {
	var err error
	gh, err = geoipdb.NewHandler()
	if err != nil {
		t.Fatalf("geoipdb.New failed: %s", err)
	}
}

func TestLibGeoipLookup(t *testing.T) {
	var asnDescr string
	asnLibGeo, asnDescr = gh.LibGeoipLookup(ip)
	if asnLibGeo == "" {
		t.Fatalf("ASN of ip '%s' is unknown by libgeoip", ip)
	}
	t.Logf("libgeoip results for %s: %s %s", ip, asnLibGeo, asnDescr)
}

func TestIpInfoLookup(t *testing.T) {
	var err error
	var asnDescr string
	asnIpInfo, asnDescr, err = gh.IpInfoLookup(ip)
	if err != nil {
		t.Fatalf("IpInfoLookup failed: %s", err)
	}
	t.Logf("ipinfo.io results for %s: %s %s", ip, asnIpInfo, asnDescr)
}

func TestCymruDnsLookup(t *testing.T) {
	if asnLibGeo != "" {
		asnDescr, err := gh.CymruDnsLookup(asnLibGeo)
		if err != nil {
			t.Fatalf("CymruDnsLookup failed for '%s': %s", asnLibGeo, err)
		}
		t.Logf("CymruDnsLookup results for %s: %s", asnLibGeo, asnDescr)
	}
	if asnIpInfo != "" && asnIpInfo != asnLibGeo {
		asnDescr, err := gh.CymruDnsLookup(asnIpInfo)
		if err != nil {
			t.Fatalf("CymruDnsLookup failed for '%s': %s", asnIpInfo, err)
		}
		t.Logf("CymruDnsLookup results for %s: %s", asnIpInfo, asnDescr)
	}
}

func TestLookupAsn(t *testing.T) {
	asn, asnDescr, err := gh.LookupAsn(ip)
	if err != nil {
		t.Fatalf("LookupAsn failed for %s: %s", ip, err)
	}
	t.Logf("LookupAsn results: %s %s", asn, asnDescr)
}
