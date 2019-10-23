package config_test

import (
	"testing"

	"github.com/gojekfarm/envoy-lb-operator/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadEnvoyConfigSetsDefaultsForOptionalConfigsIfNotConfigured(t *testing.T) {
	expectedEnvoyConfig := config.EnvoyConfig{
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

	config.MustLoad("sample", "./testdata")

	assert.Equal(t, expectedEnvoyConfig, config.GetEnvoyConfig())
	assert.Equal(t, 10, config.RefreshIntervalInS())
	assert.Equal(t, "info", config.LogLevel().String())
}

func TestShouldLoadEnvoyDiscoveryMapping(t *testing.T) {
	expectedEnvoyDiscoveryMap := []config.DiscoveryMap{
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


	assert.Equal(t, expectedEnvoyDiscoveryMap, config.GetDiscoveryMapping())
}
