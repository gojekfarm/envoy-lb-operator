package config

import (
	"github.com/spf13/viper"
)

type CircuitBreakerConfig struct {
	MaxConnections     uint32
	MaxRequests        uint32
	MaxPendingRequests uint32
	MaxRetries         uint32
}

type OutlierDetectionConfig struct {
	BaseEjectionTimeInSeconds          int64
	EjectionSweepIntervalInSeconds     int64
	Consecutive5xx                     uint32
	ConsecutiveGatewayFailure          uint32
	EnforcingConsecutive5xx            uint32
	EnforcingConsecutiveGatewayFailure uint32
	MaxEjectionPercent                 uint32
}

type RetryConfig struct {
	RetryOn                       string
	RetryPredicate                string
	NumRetries                    uint32
	HostSelectionMaxRetryAttempts int64
}

type EnvoyClusterConfig struct {
	ConnectTimeoutMs int
	CircuitBreaker   CircuitBreakerConfig
	OutlierDetection OutlierDetectionConfig
}

type EnvoyVHostConfig struct {
	RetryConfig RetryConfig
}

type EnvoyConfig struct {
	EnvoyVHostConfig
	EnvoyClusterConfig
}

func loadEnvoyConfig() EnvoyConfig {
	return EnvoyConfig{
		EnvoyVHostConfig:   getEnvoyVHostConfig(),
		EnvoyClusterConfig: getEnvoyClusterConfig(),
	}
}

func getEnvoyClusterConfig() EnvoyClusterConfig {
	viper.SetDefault("envoy.cluster.connect_timeout_ms", 1000)

	return EnvoyClusterConfig{
		ConnectTimeoutMs: viper.GetInt("envoy.cluster.connect_timeout_ms"),
		CircuitBreaker:   getCircuitBreakerConfig(),
		OutlierDetection: getOutlierDetectionConfig(),
	}
}

func getEnvoyVHostConfig() EnvoyVHostConfig {
	return EnvoyVHostConfig{
		RetryConfig: getRetryConfig(),
	}
}

func getCircuitBreakerConfig() CircuitBreakerConfig {
	viper.SetDefault("envoy.cluster.circuit_breaker.max_connections", 1024)
	viper.SetDefault("envoy.cluster.circuit_breaker.max_requests", 50000)
	viper.SetDefault("envoy.cluster.circuit_breaker.max_pending_requests", 50000)
	viper.SetDefault("envoy.cluster.circuit_breaker.max_retries", 50000)

	return CircuitBreakerConfig{
		MaxConnections:     viper.GetUint32("envoy.cluster.circuit_breaker.max_connections"),
		MaxRequests:        viper.GetUint32("envoy.cluster.circuit_breaker.max_requests"),
		MaxPendingRequests: viper.GetUint32("envoy.cluster.circuit_breaker.max_pending_requests"),
		MaxRetries:         viper.GetUint32("envoy.cluster.circuit_breaker.max_retries"),
	}
}

func getOutlierDetectionConfig() OutlierDetectionConfig {
	viper.SetDefault("envoy.cluster.outlier_detection.base_ejection_time_in_seconds", 30)
	viper.SetDefault("envoy.cluster.outlier_detection.ejection_sweep_interval_in_seconds", 10)
	viper.SetDefault("envoy.cluster.outlier_detection.consecutive_5xx", 10000)
	viper.SetDefault("envoy.cluster.outlier_detection.consecutive_gateway_failure", 5)
	viper.SetDefault("envoy.cluster.outlier_detection.enforcing_consecutive_5xx", 0)
	viper.SetDefault("envoy.cluster.outlier_detection.enforcing_consecutive_gateway_failure", 100)
	viper.SetDefault("envoy.cluster.outlier_detection.max_ejection_percent", 50)

	return OutlierDetectionConfig{
		BaseEjectionTimeInSeconds:          viper.GetInt64("envoy.cluster.outlier_detection.base_ejection_time_in_seconds"),
		EjectionSweepIntervalInSeconds:     viper.GetInt64("envoy.cluster.outlier_detection.ejection_sweep_interval_in_seconds"),
		Consecutive5xx:                     viper.GetUint32("envoy.cluster.outlier_detection.consecutive_5xx"),
		ConsecutiveGatewayFailure:          viper.GetUint32("envoy.cluster.outlier_detection.consecutive_gateway_failure"),
		EnforcingConsecutive5xx:            viper.GetUint32("envoy.cluster.outlier_detection.enforcing_consecutive_5xx"),
		EnforcingConsecutiveGatewayFailure: viper.GetUint32("envoy.cluster.outlier_detection.enforcing_consecutive_gateway_failure"),
		MaxEjectionPercent:                 viper.GetUint32("envoy.cluster.outlier_detection.max_ejection_percent"),
	}
}

func getRetryConfig() RetryConfig {
	viper.SetDefault("envoy.vhost.retry.retry_on", "connect-failure")
	viper.SetDefault("envoy.vhost.retry.retry_predicate", "envoy.retry_host_predicates.previous_hosts")
	viper.SetDefault("envoy.vhost.retry.num_retries", 3)
	viper.SetDefault("envoy.vhost.retry.host_selection_max_retry_attempts", 3)

	return RetryConfig{
		RetryOn:                       viper.GetString("envoy.vhost.retry.retry_on"),
		RetryPredicate:                viper.GetString("envoy.vhost.retry.retry_predicate"),
		NumRetries:                    viper.GetUint32("envoy.vhost.retry.num_retries"),
		HostSelectionMaxRetryAttempts: viper.GetInt64("envoy.vhost.retry.host_selection_max_retry_attempts"),
	}
}
