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

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/matthewpi/cloudflare-exporter/internal/cloudflare"
	"github.com/matthewpi/cloudflare-exporter/internal/metrics"
)

var cf *cloudflare.Cloudflare

var zones []string
var zoneIDMap = map[string]string{}

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.String("bind", ":8089", "")
	fs.String("token", "", "")
	fs.String("email", "", "")
	fs.String("key", "", "")
	fs.String("zones", "", "comma separated list of zone_id:domain")
	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	fs.VisitAll(func(f *flag.Flag) {
		if f.Value != nil && f.Value.String() != "" {
			return
		}
		if f.Name == "zones" {
			fmt.Println("Zones: " + f.Value.String())
			return
		}

		if f.Name != "token" && f.Name != "email" && f.Name != "key" {
			return
		}

		if err := f.Value.Set(os.Getenv("CF_" + strings.ToUpper(f.Name))); err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
	})

	zonesFlag := fs.Lookup("zones").Value.String()
	if zonesFlag == "" {
		fmt.Println("no zones specified")
		os.Exit(1)
		return
	}
	if strings.Contains(zonesFlag, ",") {
		for _, z := range strings.Split(fs.Lookup("zones").Value.String(), ",") {
			s := strings.SplitN(z, ":", 2)
			if len(s) != 2 {
				fmt.Printf("invalid zone \"%s\": missing `:` (zone_id:domain)\n", z)
				os.Exit(1)
				return
			}
			zones = append(zones, s[0])
			zoneIDMap[s[0]] = s[1]
		}
	} else {
		s := strings.SplitN(zonesFlag, ":", 2)
		if len(s) != 2 {
			fmt.Printf("invalid zone \"%s\": missing `:` (zone_id:domain)\n", zonesFlag)
			os.Exit(1)
			return
		}
		zones = append(zones, s[0])
		zoneIDMap[s[0]] = s[1]
	}

	var (
		auth cloudflare.Auth
		err  error
	)
	email, key := fs.Lookup("email"), fs.Lookup("key")
	if email != nil && key != nil && email.Value.String() != "" && key.Value.String() != "" {
		auth, err = cloudflare.NewKeyAuthorization(email.Value.String(), key.Value.String())
	} else if token := fs.Lookup("token"); token != nil && token.Value.String() != "" {
		auth, err = cloudflare.NewTokenAuthorization(token.Value.String())
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	if auth == nil {
		fmt.Println("no authentication method specified")
		os.Exit(1)
		return
	}
	cf, err = cloudflare.New(auth)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// Create a context that is cancelled by an interrupt signal.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Start scraping metrics from Cloudflare.
	go updateTask(ctx)

	// Define a /metrics route.
	http.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}

		metrics.WritePrometheus(w, false)
	}))

	// Start the http server.
	go func(ctx context.Context, bind string) {
		fmt.Println("listening on :8089")
		var lc net.ListenConfig
		l, err := lc.Listen(ctx, "tcp", bind)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		s := &http.Server{
			Addr:    bind,
			Handler: nil,
		}
		if err := s.Serve(l); err != nil &&
			err != http.ErrServerClosed &&
			!strings.HasSuffix(err.Error(), " use of closed network connection") {
			fmt.Println(err)
		}
	}(ctx, fs.Lookup("bind").Value.String())

	// Block until we receive a signal.
	<-ctx.Done()
	fmt.Println("received signal")
	cancel()
}

func updateTask(ctx context.Context) {
	// Initially fetch the metrics.
	if err := fetchMetrics(ctx); err != nil {
		fmt.Printf("failed to fetch metrics: %v\n", err)
	}

	// Make the ticker start at 0 seconds so it runs exactly when the minute
	// changes.
	time.Sleep(time.Duration(60-time.Now().Second()) * time.Second)

	t := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-t.C:
				go func() {
					if err := fetchMetrics(ctx); err != nil {
						fmt.Printf("failed to fetch metrics: %v\n", err)
					}
				}()
			}
		}
	}()
}

func fetchMetrics(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	r, err := cf.Zone(
		ctx,
		zones,
	)
	if err != nil {
		return err
	}
	cancel()

	for _, z := range r.Viewer.Zones {
		zone := zoneIDMap[z.ZoneID]
		// for _, e := range z.FirewallEventsAdaptiveGroups {}
		// for _, e := range z.HealthCheckEventsAdaptive {}

		// HTTPRequests1mGroups
		for _, e := range z.HTTPRequests1mGroups {
			metrics.ZoneRequestsTotal(zone).Add(int(e.Sum.Requests))
			metrics.ZoneRequestsCached(zone).Add(int(e.Sum.CachedRequests))
			metrics.ZoneRequestsEncrypted(zone).Add(int(e.Sum.EncryptedRequests))

			metrics.ZoneBandwidthTotal(zone).Add(int(e.Sum.Bytes))
			metrics.ZoneBandwidthCached(zone).Add(int(e.Sum.CachedBytes))
			metrics.ZoneBandwidthEncrypted(zone).Add(int(e.Sum.EncryptedBytes))

			metrics.ZoneThreatsTotal(zone).Add(int(e.Sum.Threats))

			for _, ct := range e.Sum.ContentTypeMap {
				metrics.ZoneRequestsContentType(zone, ct.EdgeResponseContentType).
					Add(int(ct.Requests))
				metrics.ZoneBandwidthContentType(zone, ct.EdgeResponseContentType).
					Add(int(ct.Bytes))
			}

			for _, c := range e.Sum.CountryMap {
				metrics.ZoneRequestsCountry(zone, c.ClientCountryName).
					Add(int(c.Requests))
				metrics.ZoneBandwidthCountry(zone, c.ClientCountryName).
					Add(int(c.Bytes))
				metrics.ZoneThreatsCountry(zone, c.ClientCountryName).
					Add(int(c.Threats))
			}

			for _, s := range e.Sum.ResponseStatusMap {
				metrics.ZoneRequestsStatus(zone, strconv.Itoa(s.EdgeResponseStatus)).
					Add(int(s.Requests))
			}

			for _, t := range e.Sum.ThreatPathingMap {
				metrics.ZoneThreatsType(zone, t.Name).Add(int(t.Requests))
			}
		}
		// END HTTPRequests1mGroups

		// HTTPRequestsAdaptiveGroups
		for _, e := range z.HTTPRequestsAdaptiveGroups {
			metrics.ZoneColocationVisits(zone, e.Dimensions.ColoCode).Add(int(e.Sum.Visits))
			metrics.ZoneColocationResponseBytes(zone, e.Dimensions.ColoCode).
				Add(int(e.Sum.EdgeResponseBytes))
		}
		// END HTTPRequestsAdaptiveGroups

		// for _, e := range z.LoadBalancingRequestsAdaptive {}
	}
	return nil
}
