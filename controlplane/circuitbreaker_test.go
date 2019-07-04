package controlplane_test

import (
	cp "github.com/gojekfarm/envoy-lb-operator/controlplane"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCircuitBreaker(t *testing.T) {
	var maxConnections, maxRequests, maxPendingRequests, maxRetries uint32
	maxConnections = 10
	maxRequests = 20
	maxPendingRequests = 30
	maxRetries = 40

	circuitBreaker := cp.CircuitBreaker(maxConnections, maxRequests, maxPendingRequests, maxRetries)

	assert.Equal(t, maxConnections, circuitBreaker.Thresholds[0].MaxConnections.Value)
	assert.Equal(t, maxRequests, circuitBreaker.Thresholds[0].MaxRequests.Value)
	assert.Equal(t, maxPendingRequests, circuitBreaker.Thresholds[0].MaxPendingRequests.Value)
	assert.Equal(t, maxRetries, circuitBreaker.Thresholds[0].MaxRetries.Value)
}
