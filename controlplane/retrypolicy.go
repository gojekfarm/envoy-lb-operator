package controlplane

import (
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/gogo/protobuf/types"
)

func RetryPolicy(retryOn, retryPredicate string, numRetries uint32, hostSelectionMaxRetryAttempts int64) *route.RetryPolicy {
	return &route.RetryPolicy{
		RetryOn:    retryOn,
		NumRetries: &types.UInt32Value{Value: numRetries},
		RetryHostPredicate: []*route.RetryPolicy_RetryHostPredicate{{
			Name: retryPredicate,
		}},
		HostSelectionRetryMaxAttempts: hostSelectionMaxRetryAttempts,
	}
}
