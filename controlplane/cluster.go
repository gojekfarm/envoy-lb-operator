package controlplane

import (
	"time"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
)

//StrictDNSLRCluster creates a strict dns cluster with lb policy as least request.
func StrictDNSLRCluster(name, svcAddress string, port uint32, timeoutms int) *v2.Cluster {
	endpointAddress := TCPAddress(svcAddress, port)
	return &v2.Cluster{
		Name:                 name,
		ConnectTimeout:       time.Duration(timeoutms) * time.Millisecond,
		ClusterDiscoveryType: &v2.Cluster_Type{Type: v2.Cluster_STRICT_DNS},
		DnsLookupFamily:      v2.Cluster_V4_ONLY,
		LbPolicy:             v2.Cluster_LEAST_REQUEST,
		LoadAssignment: &v2.ClusterLoadAssignment{
			ClusterName: name,
			Endpoints: []endpoint.LocalityLbEndpoints{{
				LbEndpoints: []endpoint.LbEndpoint{{
					HostIdentifier: &endpoint.LbEndpoint_Endpoint{
						Endpoint: &endpoint.Endpoint{
							Address: &endpointAddress,
						},
					},
				}},
			}},
		},
	}
}

//StrictDNSLRHttp2Cluster creates an http2 strict dns cluster with lb policy as least request.
func StrictDNSLRHttp2Cluster(name, svcAddress string, port uint32, timeoutms int) *v2.Cluster {
	c := StrictDNSLRCluster(name, svcAddress, port, timeoutms)
	c.Http2ProtocolOptions = &core.Http2ProtocolOptions{}
	return c
}
