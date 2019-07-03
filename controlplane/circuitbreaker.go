package controlplane

import (
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/cluster"
	"github.com/gogo/protobuf/types"
)

func CircuitBreaker(maxConnections, maxRequests, maxPendingRequests, maxRetries uint32) *cluster.CircuitBreakers {
	return &cluster.CircuitBreakers{
		Thresholds: []*cluster.CircuitBreakers_Thresholds{{
			MaxConnections:     &types.UInt32Value{Value: maxConnections},
			MaxRequests:        &types.UInt32Value{Value: maxRequests},
			MaxPendingRequests: &types.UInt32Value{Value: maxPendingRequests},
			MaxRetries:         &types.UInt32Value{Value: maxRetries},
		}},
	}
}
