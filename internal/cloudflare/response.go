//
// Copyright (c) 2021 Matthew Penner
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
//

package cloudflare

import (
	"time"
)

// Response .
type Response struct {
	// Viewer .
	Viewer ResponseViewer `json:"viewer"`
}

// ResponseViewer .
type ResponseViewer struct {
	// Zones .
	Zones []Zone `json:"zones"`
}

// Zone .
type Zone struct {
	// ZoneID .
	ZoneID string `json:"zoneTag"`

	// FirewallEventsAdaptiveGroups .
	FirewallEventsAdaptiveGroups []FirewallEvent `json:"firewallEventsAdaptiveGroups"`

	// HealthCheckEventsAdaptive .
	HealthCheckEventsAdaptive []HealthCheckEvent `json:"healthCheckEventsAdaptive"`

	// HTTPRequests1mGroups .
	HTTPRequests1mGroups []HTTPRequest1m `json:"httpRequests1mGroups"`

	// HTTPRequestsAdaptiveGroups .
	HTTPRequestsAdaptiveGroups []HTTPRequestAdaptive `json:"httpRequestsAdaptiveGroups"`

	// LoadBalancingRequestsAdaptive .
	LoadBalancingRequestsAdaptive []LoadBalancingRequest `json:"loadBalancingRequestsAdaptive"`
}

// FirewallEvent .
type FirewallEvent struct {
	Count      uint64                  `json:"count"`
	Dimensions FirewallEventDimensions `json:"dimensions"`
}

// FirewallEventDimensions .
type FirewallEventDimensions struct {
	Action                string `json:"action"`
	ClientCountryName     string `json:"clientCountryName"`
	ClientRequestHTTPHost string `json:"clientRequestHTTPHost"`
}

// HealthCheckEvent .
type HealthCheckEvent struct {
	DateTime              time.Time `json:"datetime"`
	EventID               string    `json:"eventId"`
	ExpectedResponseCodes string    `json:"expectedResponseCodes"`
	FailureReason         string    `json:"failureReason"`
	FQDN                  string    `json:"fqdn"`
	HealthChanged         uint8     `json:"healthChanged"`
	HealthCheckID         string    `json:"healthCheckId"`
	HealthCheckName       string    `json:"healthCheckName"`
	HealthStatus          string    `json:"healthStatus"`
	OriginIP              string    `json:"originIP"`
	OriginResponseStatus  uint16    `json:"originResponseStatus"`
	Region                string    `json:"region"`
	RttMs                 uint64    `json:"rttMs"`
	SampleInterval        uint32    `json:"sampleInterval"`
	Scope                 string    `json:"scope"`
	TCPConnMs             uint32    `json:"tcpConnMs"`
	TimeToFirstByteMs     uint32    `json:"timeToFirstByteMs"`
	TLSHandshakeMs        uint32    `json:"tlsHandshakeMs"`
}

// HTTPRequest1m .
type HTTPRequest1m struct {
	Dimensions HTTPRequest1mDimensions `json:"dimensions"`
	Sum        HTTPRequest1mSum        `json:"sum"`
	Unique     HTTPRequest1mUnique     `json:"uniq"`
}

// HTTPRequest1mDimensions .
type HTTPRequest1mDimensions struct {
	DateTime time.Time `json:"datetime"`
}

// HTTPRequest1mSum .
type HTTPRequest1mSum struct {
	BrowserMap []struct {
		PageViews       uint64 `json:"pageViews"`
		UABrowserFamily string `json:"uaBrowserFamily"`
	} `json:"browserMap"`
	Bytes                uint64 `json:"bytes"`
	CachedBytes          uint64 `json:"cachedBytes"`
	CachedRequests       uint64 `json:"cachedRequests"`
	ClientHTTPVersionMap []struct {
		Protocol string `json:"clientHTTPProtocol"`
		Requests uint64 `json:"requests"`
	} `json:"clientHTTPVersionMap"`
	ClientSSLMap []struct {
		Protocol string `json:"clientSSLProtocol"`
		Requests uint64 `json:"requests"`
	} `json:"clientSSLMap"`
	ContentTypeMap []struct {
		Bytes                   uint64 `json:"bytes"`
		EdgeResponseContentType string `json:"edgeResponseContentTypeName"`
		Requests                uint64 `json:"requests"`
	} `json:"contentTypeMap"`
	CountryMap []struct {
		Bytes             uint64 `json:"bytes"`
		ClientCountryName string `json:"clientCountryName"`
		Requests          uint64 `json:"requests"`
		Threats           uint64 `json:"threats"`
	} `json:"countryMap"`
	EncryptedBytes    uint64 `json:"encryptedBytes"`
	EncryptedRequests uint64 `json:"encryptedRequests"`
	IPClassMap        []struct {
		Type     string `json:"ipType"`
		Requests uint64 `json:"requests"`
	} `json:"ipClassMap"`
	PageViews         uint64 `json:"pageViews"`
	Requests          uint64 `json:"requests"`
	ResponseStatusMap []struct {
		EdgeResponseStatus int    `json:"edgeResponseStatus"`
		Requests           uint64 `json:"requests"`
	} `json:"responseStatusMap"`
	ThreatPathingMap []struct {
		Name     string `json:"threatPathingName"`
		Requests uint64 `json:"requests"`
	} `json:"threatPathingMap"`
	Threats uint64 `json:"threats"`
}

// HTTPRequest1mUnique .
type HTTPRequest1mUnique struct {
	Uniques uint64 `json:"uniques"`
}

// HTTPRequestAdaptive .
type HTTPRequestAdaptive struct {
	Count      uint64                        `json:"count"`
	Average    HTTPRequestAdaptiveAverage    `json:"avg"`
	Dimensions HTTPRequestAdaptiveDimensions `json:"dimensions"`
	Sum        HTTPRequestAdaptiveSum        `json:"sum"`
}

// HTTPRequestAdaptiveAverage .
type HTTPRequestAdaptiveAverage struct {
	SampleInterval float64 `json:"sampleInterval"`
}

// HTTPRequestAdaptiveDimensions .
type HTTPRequestAdaptiveDimensions struct {
	ColoCode string    `json:"coloCode"`
	DateTime time.Time `json:"datetime"`
}

// HTTPRequestAdaptiveSum .
type HTTPRequestAdaptiveSum struct {
	EdgeResponseBytes uint64 `json:"edgeResponseBytes"`
	Visits            uint64 `json:"visits"`
}

// LoadBalancingRequest .
type LoadBalancingRequest struct {
	ColoCode              string    `json:"coloCode"`
	DateTime              time.Time `json:"datetime"`
	ErrorType             string    `json:"errorType"`
	LBName                string    `json:"lbName"`
	NumberOriginsSelected uint16    `json:"numberOriginsSelected"`
	Origins               []Origin  `json:"origins"`
	Pools                 []Pool    `json:"pools"`
	Region                string    `json:"region"`
	SampleInterval        uint32    `json:"sampleInterval"`
	SessionAffinity       string    `json:"sessionAffinity"`
	SteeringPolicy        string    `json:"steeringPolicy"`
}

// Origin .
type Origin struct {
	FQDN       string  `json:"fqdn"`
	Health     uint8   `json:"health"`
	IPv4       string  `json:"ipv4"`
	IPv6       string  `json:"ipv6"`
	OriginName string  `json:"originName"`
	Selected   uint8   `json:"selected"`
	Weight     float64 `json:"weight"`
}

// Pool .
type Pool struct {
	AvgRttMs           uint64 `json:"avgRttMs"`
	HealthCheckEnabled uint8  `json:"healthCheckEnabled"`
	Healthy            uint8  `json:"healthy"`
	ID                 string `json:"id"`
	PoolName           string `json:"poolName"`
}
