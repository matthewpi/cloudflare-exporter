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

// Package metrics ...
package metrics

import (
	"io"

	"github.com/VictoriaMetrics/metrics"
)

// WritePrometheus writes all the registered metrics in Prometheus format to w.
//
// If exposeProcessMetrics is true, then various `go_*` and `process_*` metrics
// are exposed for the current process.
//
// The WritePrometheus func is usually called inside "/metrics" handler:
//
//     http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
//         metrics.WritePrometheus(w, true)
//     })
//
func WritePrometheus(w io.Writer, exposeProcessMetrics bool) {
	metrics.WritePrometheus(w, exposeProcessMetrics)
}

// ZoneRequestsTotal .
func ZoneRequestsTotal(zone string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_requests_total{" +
			"zone=\"" + zone + "\"" +
			"}",
	)
}

// ZoneRequestsCached .
func ZoneRequestsCached(zone string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_requests_cached{" +
			"zone=\"" + zone + "\"" +
			"}",
	)
}

// ZoneRequestsEncrypted .
func ZoneRequestsEncrypted(zone string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_requests_encrypted{" +
			"zone=\"" + zone + "\"" +
			"}",
	)
}

// ZoneRequestsContentType .
func ZoneRequestsContentType(zone, contentType string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_requests_content_type{" +
			"zone=\"" + zone + "\"," +
			"content_type=\"" + contentType + "\"" +
			"}",
	)
}

// ZoneRequestsCountry .
func ZoneRequestsCountry(zone, country string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_requests_country{" +
			"zone=\"" + zone + "\"," +
			"country=\"" + country + "\"" +
			"}",
	)
}

// ZoneRequestsStatus .
func ZoneRequestsStatus(zone, status string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_requests_status{" +
			"zone=\"" + zone + "\"," +
			"status=\"" + status + "\"" +
			"}",
	)
}

// ZoneBandwidthTotal .
func ZoneBandwidthTotal(zone string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_bandwidth_total{" +
			"zone=\"" + zone + "\"" +
			"}",
	)
}

// ZoneBandwidthCached .
func ZoneBandwidthCached(zone string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_bandwidth_cached{" +
			"zone=\"" + zone + "\"" +
			"}",
	)
}

// ZoneBandwidthEncrypted .
func ZoneBandwidthEncrypted(zone string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_bandwidth_encrypted{" +
			"zone=\"" + zone + "\"" +
			"}",
	)
}

// ZoneBandwidthContentType .
func ZoneBandwidthContentType(zone, contentType string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_bandwidth_content_type{" +
			"zone=\"" + zone + "\"," +
			"content_type=\"" + contentType + "\"" +
			"}",
	)
}

// ZoneBandwidthCountry .
func ZoneBandwidthCountry(zone, country string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_bandwidth_country{" +
			"zone=\"" + zone + "\"," +
			"country=\"" + country + "\"" +
			"}",
	)
}

// ZoneColocationVisits .
func ZoneColocationVisits(zone, colocation string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_colocation_visits{" +
			"zone=\"" + zone + "\"," +
			"colocation=\"" + colocation + "\"" +
			"}",
	)
}

// ZoneColocationResponseBytes .
func ZoneColocationResponseBytes(zone, colocation string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_colocation_response_bytes{" +
			"zone=\"" + zone + "\"," +
			"colocation=\"" + colocation + "\"" +
			"}",
	)
}

// ZoneThreatsTotal .
func ZoneThreatsTotal(zone string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_threats_total{" +
			"zone=\"" + zone + "\"" +
			"}",
	)
}

// ZoneThreatsCountry .
func ZoneThreatsCountry(zone, country string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_threats_country{" +
			"zone=\"" + zone + "\"," +
			"country=\"" + country + "\"" +
			"}",
	)
}

// ZoneThreatsType .
func ZoneThreatsType(zone, threatType string) *metrics.Counter {
	return metrics.GetOrCreateCounter(
		"cloudflare_zone_threats_type{" +
			"zone=\"" + zone + "\"," +
			"type=\"" + threatType + "\"" +
			"}",
	)
}
