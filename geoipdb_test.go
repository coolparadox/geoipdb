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
	"fmt"
	"strings"
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

const (
	asnGoogle = "AS15169"
	asnLevel3 = "AS3356"
)

const (
	subStrGoogle = "Google Inc"
	subStrLevel3 = "Level 3 Communications"
)

func verifyAsn(t *testing.T, asn string, descr string) {
	switch asn {
	case asnGoogle:
		if !strings.Contains(descr, subStrGoogle) {
			t.Fatalf("%s description does not contain '%s': %s", asn, subStrGoogle, descr)
		}
	case asnLevel3:
		if !strings.Contains(descr, subStrLevel3) {
			t.Fatalf("%s description does not contain '%s': %s", asn, subStrLevel3, descr)
		}
	default:
		t.Fatalf("unexpected ASN identification '%s'", asn)
	}
}

func TestLibGeoipLookup(t *testing.T) {
	var asnDescr string
	asnLibGeo, asnDescr = gh.LibGeoipLookup(ip)
	if asnLibGeo == "" {
		t.Fatalf("ASN of ip '%s' is unknown by libgeoip", ip)
	}
	t.Logf("libgeoip result for %s: %s %s", ip, asnLibGeo, asnDescr)
	verifyAsn(t, asnLibGeo, asnDescr)
}

func TestIpInfoLookup(t *testing.T) {
	var err error
	var asnDescr string
	asnIpInfo, asnDescr, err = gh.IpInfoLookup(ip)
	if err != nil {
		t.Fatalf("IpInfoLookup failed: %s", err)
	}
	t.Logf("ipinfo result for %s: %s %s", ip, asnIpInfo, asnDescr)
	verifyAsn(t, asnIpInfo, asnDescr)
}

func TestCymruDnsLookup(t *testing.T) {
	if asnLibGeo != "" {
		asnDescr, err := gh.CymruDnsLookup(asnLibGeo)
		if err != nil {
			t.Fatalf("CymruDnsLookup failed for '%s': %s", asnLibGeo, err)
		}
		t.Logf("cymru result for %s: %s", asnLibGeo, asnDescr)
		verifyAsn(t, asnLibGeo, asnDescr)
	}
	if asnIpInfo != "" && asnIpInfo != asnLibGeo {
		asnDescr, err := gh.CymruDnsLookup(asnIpInfo)
		if err != nil {
			t.Fatalf("CymruDnsLookup failed for '%s': %s", asnIpInfo, err)
		}
		t.Logf("cymru result for %s: %s", asnIpInfo, asnDescr)
		verifyAsn(t, asnIpInfo, asnDescr)
	}
}

func TestLookupAsn(t *testing.T) {
	asn, asnDescr, err := gh.LookupAsn(ip)
	if err != nil {
		t.Fatalf("LookupAsn failed for %s: %s", ip, err)
	}
	t.Logf("LookupAsn results: %s %s", asn, asnDescr)
	verifyAsn(t, asn, asnDescr)
}

func Example_lookupAsn() {

	ip := "8.8.8.8"
	gh, err := geoipdb.NewHandler()
	if err != nil {
		panic(err)
	}
	asn, descr, err := gh.LookupAsn(ip)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ASN for %s: %s (%s)\n", ip, asn, descr)


	// Output:
	// ASN for 8.8.8.8: AS15169 (Google Inc.)

}
