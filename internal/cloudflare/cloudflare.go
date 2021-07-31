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

// Package cloudflare ...
package cloudflare

import (
	"context"
	"time"

	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
)

// Cloudflare .
type Cloudflare struct {
	// Auth .
	Auth Auth

	// graphql .
	graphql *graphql.Client
}

// New .
func New(auth Auth) (*Cloudflare, error) {
	return &Cloudflare{
		Auth: auth,

		graphql: graphql.NewClient("https://api.cloudflare.com/client/v4/graphql"),
	}, nil
}

// Zone .
func (cf *Cloudflare) Zone(ctx context.Context, zones []string) (Response, error) {
	r := graphql.NewRequest(`
		query ($zoneIDs: [String!], $mintime: Time!, $maxtime: Time!, $limit: Int!) {
			viewer {
				zones (filter: { zoneTag_in: $zoneIDs }) {
					zoneTag

					httpRequests1mGroups (limit: $limit, filter: { datetime: $maxtime }) {
						uniq {
							uniques
						}

						sum {
							browserMap {
								pageViews
								uaBrowserFamily
							}

							bytes
							cachedBytes
							cachedRequests

							clientHTTPVersionMap {
								clientHTTPProtocol
								requests
							}

							clientSSLMap {
								clientSSLProtocol
								requests
							}

							contentTypeMap {
								bytes
								requests
								edgeResponseContentTypeName
							}

							countryMap {
								bytes
								clientCountryName
								requests
								threats
							}

							encryptedBytes
							encryptedRequests

							ipClassMap {
								ipType
								requests
							}

							pageViews
							requests

							responseStatusMap {
								edgeResponseStatus
								requests
							}

							threatPathingMap {
								requests
								threatPathingName
							}

							threats
						}

						dimensions {
							datetime
						}
					}

					httpRequestsAdaptiveGroups (limit: $limit, filter: { datetime_geq: $mintime, datetime_lt: $maxtime }) {
						count

						avg {
							sampleInterval
						}

						dimensions {
							coloCode
							datetime
						}

						sum {
							edgeResponseBytes
							visits
						}
					}
				}
			}
		}
	`)
	/*`
		query ($zoneIDs: [String!], $mintime: Time!, $maxtime: Time!, $limit: Int!) {
			viewer {
				zones (filter: { zoneTag_in: $zoneIDs }) {
					zoneTag

					healthCheckEventsAdaptive (limit: $limit, filter: { datetime_geq: $mintime, datetime_lt: $maxtime }) {
						datetime
						eventId
						expectedResponseCodes
						failureReason
						fqdn
						healthChanged
						healthCheckId
						healthCheckName
						healthStatus
						originIP
						originResponseStatus
						region
						rttMs
						sampleInterval
						scope
						tcpConnMs
						timeToFirstByteMs
						tlsHandshakeMs
					}

					httpRequests1mGroups (limit: $limit, filter: { datetime: $maxtime }) {
						uniq {
							uniques
						}

						sum {
							browserMap {
								pageViews
								uaBrowserFamily
							}

							bytes
							cachedBytes
							cachedRequests

							clientHTTPVersionMap {
								clientHTTPProtocol
								requests
							}

							clientSSLMap {
								clientSSLProtocol
								requests
							}

							contentTypeMap {
								bytes
								requests
								edgeResponseContentTypeName
							}

							countryMap {
								bytes
								clientCountryName
								requests
								threats
							}

							encryptedBytes
							encryptedRequests

							ipClassMap {
								ipType
								requests
							}

							pageViews
							requests

							responseStatusMap {
								edgeResponseStatus
								requests
							}

							threatPathingMap {
								requests
								threatPathingName
							}

							threats
						}

						dimensions {
							datetime
						}
					}

					httpRequestsAdaptiveGroups (limit: $limit, filter: { datetime_geq: $mintime, datetime_lt: $maxtime }) {
						count

						avg {
							sampleInterval
						}

						dimensions {
							coloCode
							datetime
						}

						sum {
							edgeResponseBytes
							visits
						}
					}

					firewallEventsAdaptiveGroups (limit: $limit, filter: { datetime_geq: $mintime, datetime_lt: $maxtime }) {
						count

						dimensions {
						  	action
						  	source
						  	clientRequestHTTPHost
						  	clientCountryName
						}
					}

					loadBalancingRequestsAdaptive (limit: $limit, filter: { datetime_geq: $mintime, datetime_lt: $maxtime }) {
						coloCode
						datetime
						errorType
						lbName
						numberOriginsSelected

						origins {
							fqdn
							health
							ipv4
							ipv6
							originName
							selected
							weight
						}

						pools {
							avgRttMs
							healthCheckEnabled
							healthy
							id
							poolName
						}

						region
						sampleInterval
						sessionAffinity
						steeringPolicy
					}
				}
			}
		}
	`*/
	if err := cf.Auth.Authorize(ctx, r.Header); err != nil {
		return Response{}, errors.Wrap(err, "cloudflare: failed to authorize request")
	}

	r.Header.Set("Cache-Control", "no-cache")

	now := time.Now().Add(-180 * time.Second).UTC()
	s := 60 * time.Second
	now = now.Truncate(s)
	now1mAgo := now.Add(-60 * time.Second)
	r.Var("limit", 10)
	r.Var("maxtime", now)
	r.Var("mintime", now1mAgo)
	r.Var("zoneIDs", zones)

	var resp Response
	if err := cf.graphql.Run(ctx, r, &resp); err != nil {
		return Response{}, errors.Wrap(err, "cloudflare: failed to get data")
	}
	return resp, nil
}
