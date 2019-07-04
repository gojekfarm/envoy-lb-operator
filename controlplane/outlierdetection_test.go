package controlplane_test

import (
	cp "github.com/gojekfarm/envoy-lb-operator/controlplane"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOutlierDetection(t *testing.T) {
	var baseEjectionTimeSeconds, ejectionSweepInterval int64
	var consecutive5xx, consecutiveGatewayFailure, enforcingConsecutive5xx, enforcingConsecutiveGatewayFailure, maxEjectionPercent uint32

	baseEjectionTimeSeconds = 10
	ejectionSweepInterval = 20
	consecutive5xx = 30
	consecutiveGatewayFailure = 40
	enforcingConsecutive5xx = 50
	enforcingConsecutiveGatewayFailure = 60
	maxEjectionPercent = 70

	outlierDetection := cp.OutlierDetection(baseEjectionTimeSeconds, ejectionSweepInterval, consecutive5xx, consecutiveGatewayFailure, enforcingConsecutive5xx, enforcingConsecutiveGatewayFailure, maxEjectionPercent)

	assert.Equal(t, baseEjectionTimeSeconds, outlierDetection.BaseEjectionTime.Seconds)
	assert.Equal(t, ejectionSweepInterval, outlierDetection.Interval.Seconds)
	assert.Equal(t, consecutiveGatewayFailure, outlierDetection.ConsecutiveGatewayFailure.Value)
	assert.Equal(t, consecutive5xx, outlierDetection.Consecutive_5Xx.Value)
	assert.Equal(t, enforcingConsecutiveGatewayFailure, outlierDetection.EnforcingConsecutiveGatewayFailure.Value)
	assert.Equal(t, enforcingConsecutive5xx, outlierDetection.EnforcingConsecutive_5Xx.Value)
	assert.Equal(t, maxEjectionPercent, outlierDetection.MaxEjectionPercent.Value)
}
