package config_test

import (
	"os"
	"testing"

	"github.com/gojekfarm/envoy-lb-operator/config"
	"github.com/stretchr/testify/assert"
)

var expectedEnvoyConfig = config.EnvoyConfig{
	EnvoyClusterConfig: config.EnvoyClusterConfig{
		ConnectTimeoutMs: 1000,
		CircuitBreaker: config.CircuitBreakerConfig{
			MaxConnections:     1024,
			MaxRequests:        50000,
			MaxPendingRequests: 50000,
			MaxRetries:         50000,
		},
		OutlierDetection: config.OutlierDetectionConfig{
			BaseEjectionTimeInSeconds:          30,
			EjectionSweepIntervalInSeconds:     10,
			Consecutive5xx:                     10000,
			ConsecutiveGatewayFailure:          5,
			EnforcingConsecutive5xx:            0,
			EnforcingConsecutiveGatewayFailure: 100,
			MaxEjectionPercent:                 50,
		},
	},
	EnvoyVHostConfig: config.EnvoyVHostConfig{
		RetryConfig: config.RetryConfig{
			RetryOn:                       "connect-failure",
			RetryPredicate:                "envoy.retry_host_predicates.previous_hosts",
			NumRetries:                    3,
			HostSelectionMaxRetryAttempts: 3,
		},
	},
}

func TestLoadEnvoyConfigSetsDefaultsIfNotConfigured(t *testing.T) {
	os.Create("application.yml")
	defer func() {
		os.Remove("application.yml")
	}()

	config.MustLoad("application", "./")

	assert.Equal(t, expectedEnvoyConfig, config.GetEnvoyConfig())
}

func TestShouldLoadEnvoyUpstreamMapping(t *testing.T) {
	expectedEnvoyConfig := []config.DiscoveryMap{
		{
			EnvoyId:               "stream_1",
			UpstreamEndpointLabel: "endpoint_1",
			Namespace:             "namespace_1",
		},
		{
			EnvoyId:               "stream_2",
			UpstreamEndpointLabel: "endpoint_2",
			Namespace:             "namespace_2",
		},
	}

	config.MustLoad("sample", "./testdata")

	assert.Equal(t, expectedEnvoyConfig, config.GetDiscoveryMapping())
}
