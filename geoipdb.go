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

Basics

Get a geoipdb Handler with NewHandler, and use its lookup methods at will.

Lookup of Autonomous System Numbers

For looking up the autonomous system number of an IP address, use LookupAsn
as it wraps more than one search method.

If you want a specific service to be queried for ASN,
see other Handler lookup methods.
*/
package geoipdb

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/abh/geoip"
	"github.com/miekg/dns"
	"github.com/turbobytes/geoipdb/iputils"
	"gopkg.in/mgo.v2"
)

func init() {
	// reASN is a regexp for matching against an ASN.
	reASN = regexp.MustCompilePOSIX("^AS[[:digit:]]+$")
	// reDNSFilter is a regexp for matching content in DNS answers
	// that is not part of ASN description.
	reDNSFilter = regexp.MustCompilePOSIX(".*\\|")
}

// Pre-compiled regular expressions, see init() body source.
var (
	reASN       *regexp.Regexp
	reDNSFilter *regexp.Regexp
)

var (
	// MalformedIPError is returned on parse failure of IP parameter.
	MalformedIPError = errors.New("malformed IP address")
	// PrivateIPError is returned on AS lookup of a private IP address.
	PrivateIPError = errors.New("private IP address")
)

// Handler is a handler to TurboBytes GeoIP helper functions.
type Handler struct {
	geoip4    *geoip.GeoIP
	geoip6    *geoip.GeoIP
	cymru     cymruClient
	timeout   time.Duration
	overrides *mgo.Collection
	cache     cache
}

// NewHandler creates a handler
// for accessing geoipdb features.
//
// Parameter overrides, if not nil,
// is used to access a collection of overrides of ASN descriptions.
// (See Overrides<...> methods.)
//
// Parameter timeout is honored by methods that access external services.
// Pass zero to disable timeout.
//
// Returns a geoipdb handler.
func NewHandler(overrides *mgo.Collection, timeout time.Duration) (Handler, error) {
	ge4, err := geoip.OpenType(geoip.GEOIP_ASNUM_EDITION)
	if err != nil {
		return Handler{}, fmt.Errorf("cannot open GeoIP database: %s", err)
	}
	ge6, err := geoip.OpenType(geoip.GEOIP_ASNUM_EDITION_V6)
	if err != nil {
		return Handler{}, fmt.Errorf("cannot open GeoIP database: %s", err)
	}
	cy := newCymruClient(timeout)
	return Handler{
		geoip4:    ge4,
		geoip6:    ge6,
		cymru:     cy,
		timeout:   timeout,
		overrides: overrides,
		cache:     newCache(),
	}, nil
}

// LibGeoipLookup queries the libgeoip database for the ASN of a given ip address.
//
// Returns
// an ASN identification
// and the corresponding description.
func (h Handler) LibGeoipLookup(ip string) (string, string) {
	var name string
	ipAddr, isIPv4 := iputils.ParseIP(ip)
	if ipAddr == nil {
		return "", ""
	}
	if isIPv4 {
		name, _ = h.geoip4.GetName(ip)
	} else {
		name, _ = h.geoip6.GetNameV6(ip)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ""
	}
	answer := strings.SplitN(name, " ", 2)
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
// Particularly, the overrides collection (see NewHandler)
// takes precedence for querying ASN descriptions.
//
// Data returned by LookupAsn is cached with a 1 day TTL.
// Also see: AsnCachePurge.
//
// Returns
// an ASN identification
// and the corresponding description.
func (h Handler) LookupAsn(ip string) (string, string, error) {
	// Sanity check input
	ipAddr, _ := iputils.ParseIP(ip)
	if ipAddr == nil {
		return "", "", MalformedIPError
	}
	if iputils.IsLocalIP(ipAddr) {
		return "", "", PrivateIPError
	}
	// Try cache
	asn, descr, expired, found := h.cache.lookupByIP(ip)
	if found && !expired {
		return asn, descr, nil
	}
	log.Printf("(geoipdb) cache miss for %s\n", ip)
	// Try uncached lookup
	var err error
	asn, descr, err = h.lookupAsnUncached(ip)
	if err == nil {
		// Update cache
		h.cache.store(ip, asn, descr)
	}
	return asn, descr, err
}

// lookupAsnUncached is the uncached version of LookupAsn.
func (h Handler) lookupAsnUncached(ip string) (string, string, error) {
	// Try libgeoip
	asnGi, asnDescr := h.LibGeoipLookup(ip)
	if asnGi != "" && asnDescr != "" {
		// libgeoip returned an ASN and description.
		return asnGi, h.getOverridenDescr(asnGi, asnDescr), nil
	}
	if asnGi == "" {
		log.Printf("warning: libgeoip lookup failed for ip '%s'\n", ip)
	}
	// Try ipinfo.io
	asnDescr = ""
	asnIp, asnDescr, errIp := h.IpInfoLookup(ip)
	if errIp == nil {
		if asnIp != "" && asnDescr != "" {
			// ipinfo.io returned an ASN and description.
			return asnIp, h.getOverridenDescr(asnIp, asnDescr), nil
		}
	} else {
		log.Printf("warning: ipinfo lookup failed for ip '%s': %s\n", ip, errIp)
	}
	var asn string
	if asnGi != "" {
		asn = asnGi
	} else if errIp == nil && asnIp != "" {
		asn = asnIp
	} else {
		// Cannot find an ASN. Give up.
		return "", "", fmt.Errorf("unknown ASN for ip '%v'", ip)
	}
	// We found an ASN, but no description for it.
	// Try getting one from cymru's dns service.
	asnDescr, err := h.CymruDnsLookup(asn)
	if err != nil {
		log.Printf("warning: cymru lookup failed for asn '%s': %s\n", asn, err)
		return asn, h.getOverridenDescr(asn, ""), nil
	}
	return asn, h.getOverridenDescr(asn, asnDescr), nil
}

// IpInfoLookup queries ipinfo.io for the ASN of a given ip address.
//
// Returns
// an ASN identification
// and the corresponding description.
func (h Handler) IpInfoLookup(ip string) (string, string, error) {
	client := &http.Client{
		Timeout: h.timeout,
	}
	url := fmt.Sprintf("http://ipinfo.io/%s/org", ip)
	resp, err := client.Get(url)
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
	// ipinfo.io returns errors as regular text (no out-of-band error codes).
	// Let's try to be smart and identify them.
	if !reASN.MatchString(answer[0]) {
		return "", "", fmt.Errorf("ipinfo.io lookup failed for '%s': %s", ip, asnData)
	}
	if len(answer) < 2 {
		return answer[0], "", nil
	}
	return answer[0], answer[1], nil
}

// CymruDnsLookup performs a query to Team Cymru's DNS service
// for the description of a given ASN.
//
// Returns the ASN description.
func (h Handler) CymruDnsLookup(asn string) (string, error) {
	return h.cymru.lookup(asn)
}

// cymruClient can do DNS queries to Team Cymru's database
// for retrieving ASN descriptions.
type cymruClient struct {
	dnsClient *dns.Client
	reFilter  *regexp.Regexp
}

// newCymruClient creates an initialized cymruClient.
func newCymruClient(timeout time.Duration) cymruClient {
	c := new(dns.Client)
	c.Timeout = timeout
	return cymruClient{
		dnsClient: c,
		reFilter:  reDNSFilter.Copy(),
	}
}

// lookup retrieves the description of a given ASN
// by reaching Team Cymru's DNS database.
//
// Returns the ASN description.
func (cc cymruClient) lookup(asn string) (string, error) {
	if asn == "" {
		return "", fmt.Errorf("empty asn parameter")
	}
	if cc.dnsClient == nil {
		return "", fmt.Errorf("cymruClient not initialized")
	}
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = true
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = dns.Question{
		Name:   asn + ".asn.cymru.com.",
		Qtype:  dns.TypeTXT,
		Qclass: dns.ClassINET,
	}
	// Send query to Google public dns server
	msg, _, err := cc.dnsClient.Exchange(msg, "8.8.8.8:53")
	if err != nil {
		return "", fmt.Errorf("failed to query dns: %s", err)
	}
	for _, ans := range msg.Answer {
		if t, ok := ans.(*dns.TXT); ok {
			return strings.TrimSpace(cc.reFilter.ReplaceAllString(t.Txt[0], "")), nil
		}
	}
	return "", fmt.Errorf("not yet implemented")
}

// getOverridenDescr answers the ASN description
// taken from the override collection, if found.
// Otherwise, answers the fallback parameter.
func (h Handler) getOverridenDescr(asn string, fallback string) string {
	descr, err := h.OverridesLookup(asn)
	if err != nil {
		if err != OverridesNilCollectionError && err != OverridesAsnNotFoundError {
			log.Printf("warning: %s\n", err)
		}
		return fallback
	}
	return descr
}

// AsnCachePurge erases all LookupAsn cached data.
func (h Handler) AsnCachePurge() {
	log.Println("(geoipdb) cache purge")
	h.cache.purgeAll()
}

// LookupIp searches the cache
// for all IP addresses associated with a given ASN.
//
// Returns a non nil list of IP addresses.
func (h Handler) LookupIp(asn string) []string {
	ips := h.cache.lookupByASN(asn)
	answer := make([]string, len(ips))
	var i int
	for ip, _ := range ips {
		answer[i] = ip
		i++
	}
	return answer
}

// AsnCacheList retrieves all ASNs known to the cache.
//
// Returns a non nil list of ASNs.
func (h Handler) AsnCacheList() []string {
	return h.cache.asnList()
}
