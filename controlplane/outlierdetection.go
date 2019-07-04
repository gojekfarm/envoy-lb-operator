package controlplane

import (
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/cluster"
	"github.com/gogo/protobuf/types"
)

func OutlierDetection(baseEjectionTimeSeconds, ejectionSweepInterval int64, consecutive5xx, consecutiveGatewayFailure, enforcingConsecutive5xx, enforcingConsecutiveGatewayFailure, maxEjectionPercent uint32) *cluster.OutlierDetection {
	return &cluster.OutlierDetection{
		BaseEjectionTime:                   &types.Duration{Seconds: baseEjectionTimeSeconds},
		Interval:                           &types.Duration{Seconds: ejectionSweepInterval},
		Consecutive_5Xx:                    &types.UInt32Value{Value: consecutive5xx},
		ConsecutiveGatewayFailure:          &types.UInt32Value{Value: consecutiveGatewayFailure},
		EnforcingConsecutive_5Xx:           &types.UInt32Value{Value: enforcingConsecutive5xx},
		EnforcingConsecutiveGatewayFailure: &types.UInt32Value{Value: enforcingConsecutiveGatewayFailure},
		MaxEjectionPercent:                 &types.UInt32Value{Value: maxEjectionPercent},
	}
}
