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
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/turbobytes/geoipdb"
	"github.com/turbobytes/geoipdb/iputils"
	"gopkg.in/mgo.v2"
)

// Well known IP address for testing lookups
const ip = "8.8.8.8"

// Results of some ASN lookups
var (
	asnLibGeo    string
	asnIpInfo    string
	asnLookupAsn string
)

func TestInitIp(t *testing.T) {
	t.Logf("using ip '%s' for tests", ip)
}

var gh geoipdb.Handler

func TestNewHandler(t *testing.T) {
	var err error
	gh, err = geoipdb.NewHandler(nil, time.Second*5)
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
	var err error
	var asnDescr string
	asnLookupAsn, asnDescr, err = gh.LookupAsn(ip)
	if err != nil {
		t.Fatalf("LookupAsn failed for %s: %s", ip, err)
	}
	t.Logf("LookupAsn results: %s %s", asnLookupAsn, asnDescr)
	verifyAsn(t, asnLookupAsn, asnDescr)
}

func TestAsnCachePurge(t *testing.T) {
	gh.AsnCachePurge()
	TestLookupAsn(t)
}

func Example_lookupAsn() {
	ip := "8.8.8.8"
	gh, err := geoipdb.NewHandler(nil, time.Second*5)
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

func TestOverridesLookupNilOverrides(t *testing.T) {
	_, err := gh.OverridesLookup(asnLookupAsn)
	if err != geoipdb.OverridesNilCollectionError {
		t.Fatalf("OverridesLookup returned unexpected error: %s", err)
	}
}

var (
	mgS *mgo.Session
	mgD *mgo.Database
	mgC *mgo.Collection
)

const (
	mgUrl        = "127.0.0.1"
	mgDatabase   = "dnsdist"
	mgCollection = "geoipdb_test"
)

func TestNewHandlerWithOverrides(t *testing.T) {
	var err error
	mgS, err = mgo.Dial(mgUrl)
	if err != nil {
		t.Fatalf("cannot dial to mongodb in '%s': %s", mgUrl, err)
	}
	mgD = mgS.DB(mgDatabase)
	mgC = mgD.C(mgCollection)
	mgC.DropCollection()
	gh, err = geoipdb.NewHandler(mgC, time.Second*5)
	if err != nil {
		t.Fatalf("cannot create geoipdb handler: %s", err)
	}
}

func TestOverridesListEmpty(t *testing.T) {
	overrides, err := gh.OverridesList()
	if err != nil {
		t.Fatalf("OverridesList failed: %s", err)
	}
	if overrides == nil {
		t.Fatalf("expected an empty list, got nil")
	}
	if len(overrides) != 0 {
		t.Fatalf("expected an empty list, got: %v", overrides)
	}
}

func TestOverridesLookupUnknownOverride(t *testing.T) {
	_, err := gh.OverridesLookup(asnLookupAsn)
	if err != geoipdb.OverridesAsnNotFoundError {
		t.Fatalf("OverridesLookup returned unexpected error: %s", err)
	}
}

func TestOverridesSetMalformedAsn(t *testing.T) {
	err := gh.OverridesSet("qwerty", "l33t")
	if err != geoipdb.OverridesMalformedAsnError {
		t.Fatalf("OverridesSet returned unexpected error: %s", err)
	}
}

const overridenDescr = "TurboBytes geoipdb rules!!"

func TestOverridesSet(t *testing.T) {
	err := gh.OverridesSet(asnLookupAsn, overridenDescr)
	if err != nil {
		t.Fatalf("OverridesSet failed: %s", err)
	}
}

func TestOverridesListNotEmpty(t *testing.T) {
	overrides, err := gh.OverridesList()
	if err != nil {
		t.Fatalf("OverridesList failed: %s", err)
	}
	t.Logf("overrides list: %v", overrides)
	expected := []geoipdb.AsnOverride{{Asn: asnLookupAsn, Name: overridenDescr}}
	if !reflect.DeepEqual(overrides, expected) {
		t.Fatalf("unexpected return value, expected: %v", expected)
	}
}

func TestOverridesLookupKnownOverride(t *testing.T) {
	descr, err := gh.OverridesLookup(asnLookupAsn)
	if err != nil {
		t.Fatalf("OverridesLookup failed: %s", err)
	}
	if descr != overridenDescr {
		t.Fatalf("overriden description mismatch: expected '%s', received '%s'", overridenDescr, descr)
	}
}

func TestLookupAsnWithOverride(t *testing.T) {
	_, descr, err := gh.LookupAsn(ip)
	if err != nil {
		t.Fatalf("LookupAsn failed for %s: %s", ip, err)
	}
	t.Logf("LookupAsn results: %s %s", asnLookupAsn, descr)
	if descr != overridenDescr {
		t.Fatalf("overriden description mismatch: expected '%s'", descr)
	}
}

func TestOverridesRemove(t *testing.T) {
	err := gh.OverridesRemove(asnLookupAsn)
	if err != nil {
		t.Fatalf("OverridesRemove failed: %s", err)
	}
	TestOverridesLookupUnknownOverride(t)
	TestLookupAsn(t)
	TestOverridesListEmpty(t)
}

func TestLookupIp(t *testing.T) {
	expected := []string{ip}
	ips := gh.LookupIp(asnLookupAsn)
	if !reflect.DeepEqual(ips, expected) {
		t.Fatalf("LookupIp result mismatch for %s: expected %v, got %v", asnLookupAsn, expected, ips)
	}
}

func TestAsnCacheList(t *testing.T) {
	expected := []string{asnLookupAsn}
	asns := gh.AsnCacheList()
	if !reflect.DeepEqual(asns, expected) {
		t.Fatalf("AsnCacheList result mismatch: expected %v, got %v", expected, asns)
	}
}

func TestIsLocalIP(t *testing.T) {
	publicIps := []string{
		"8.8.8.8",
		"74.125.130.100",
		"1.1.1.1",
		"45.45.45.45",
		"120.222.111.222",
		"2404:6800:4003:c01::64",
	}
	privateIps := []string{
		"127.0.0.1",
		"10.5.6.4",
		"192.168.5.99",
		"100.66.55.66",
		"fd07:a47c:3742:823e:3b02:76:982b:463",
		"::1",
	}
	for _, ip := range publicIps {
		if iputils.IsLocalIP(net.ParseIP(ip)) {
			t.Fatalf("IsLocalIp(%s) returned %s", ip, "true")
		}
	}
	for _, ip := range privateIps {
		if !iputils.IsLocalIP(net.ParseIP(ip)) {
			t.Fatalf("IsLocalIp(%s) returned %s", ip, "false")
		}
	}
}

func TestLookupAsnMalformedIP(t *testing.T) {
	ip := "192.168.0"
	_, _, err := gh.LookupAsn(ip)
	if err != geoipdb.MalformedIPError {
		t.Fatalf("unexpected LookupAsn error: %v", err)
	}
}

func TestLookupAsnIPv6(t *testing.T) {
	ip := "fd07:a47c:3742:823e:3b02:76:982b:463"
	_, _, err := gh.LookupAsn(ip)
	if true {
	//if err != geoipdb.IPv6NotSupportedError {
		t.Fatalf("unexpected LookupAsn error: %v", err)
	}
}

func TestLookupAsnPrivateIP(t *testing.T) {
	ip := "192.168.0.101"
	_, _, err := gh.LookupAsn(ip)
	if err != geoipdb.PrivateIPError {
		t.Fatalf("unexpected LookupAsn error: %v", err)
	}
}

type ipTestData struct {
	ip    string
	asn   string
	descr string
	err   string
}

func TestLookupAsnOtherIPs(t *testing.T) {
	tests := []ipTestData{
		ipTestData{"1.1.1.1", "", "", "unknown ASN for ip '1.1.1.1'"},
		ipTestData{"8.8.8.8", "AS15169", "Google Inc.", ""},
		ipTestData{"10.0.45.98", "", "", "private IP address"},
		ipTestData{"74.125.130.100", "AS15169", "Google Inc.", ""},
		ipTestData{"80.10.246.2", "", "", "unknown ASN for ip '80.10.246.2'"},
		ipTestData{"80.10.246.129", "", "", "unknown ASN for ip '80.10.246.129'"},
		ipTestData{"127.0.0.1", "", "", "private IP address"},
		ipTestData{"192.168.0.102", "", "", "private IP address"},
		ipTestData{"2001:4860:1004::876:102", "", "", "IPv6 not yet supported"},
		ipTestData{"2404:6800:4003:c01::64", "", "", "IPv6 not yet supported"},
		ipTestData{"ns1.google.com", "", "", "malformed IP address"},
		ipTestData{"ns2.google.com", "", "", "malformed IP address"},
	}
	for _, test := range tests {
		asn, descr, err := gh.LookupAsn(test.ip)
		if asn != test.asn {
			t.Fatalf("unexpected ASN returned by LookupAsn(\"%s\"): %s", test.ip, asn)
		}
		if descr != test.descr {
			t.Fatalf("unexpected AS name returned by LookupAsn(\"%s\"): %s", test.ip, descr)
		}
		if err == nil {
			if test.err != "" {
				t.Fatalf("unexpected error returned by LookupAsn(\"%s\"): %s", test.ip, "nil")
			}
		} else {
			if err.Error() != test.err {
				t.Fatalf("unexpected error returned by LookupAsn(\"%s\"): %s", test.ip, err)
			}
		}
	}
}
