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
	"net"
	"testing"

	"github.com/turbobytes/geoipdb"
)

const host = "www.turbobytes.com"

var ip string
var asn string

func TestInitIp(t *testing.T) {
	ips, err := net.LookupIP(host)
	if err != nil {
		t.Fatalf("failed to lookup ip addresses for '%s': %s", host, err)
	}
	if len(ips) < 1 {
		t.Fatalf("ip address lookup for '%s' returned empty", host)
	}
	ip = ips[0].String()
	t.Logf("using ip %s (%s) for tests", ip, host)
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
	asn, asnDescr := gh.LibGeoipLookup(ip)
	if asn == "" {
		t.Fatalf("ASN of ip '%s' is unknown by libgeoip", ip)
	}
	t.Logf("libgeoip results: %s %s", asn, asnDescr)
}

func TestIpInfoLookup(t *testing.T) {
	asn, asnDescr, err := gh.IpInfoLookup(ip)
	if err != nil {
		t.Fatalf("IpInfoLookup failed: %s", err)
	}
	t.Logf("ipinfo.io results: %s %s", asn, asnDescr)
}

func TestCymruDnsLookup(t *testing.T) {
	//if asn == "" {
		//t.Fatalf("missing asn")
	//}
	asnDescr, err := gh.CymruDnsLookup(asn)
	if err != nil {
		t.Fatalf("CymruDnsLookup failed: %s", err)
	}
	t.Logf("CymruDnsLookup results: %s", asnDescr)
}

func TestLookupAsn(t *testing.T) {
	asn, asnDescr, err := gh.LookupAsn(ip)
	if err != nil {
		t.Fatalf("LookupAsn failed for %s: %s", ip, err)
	}
	t.Logf("LookupAsn results: %s %s", asn, asnDescr)
}

