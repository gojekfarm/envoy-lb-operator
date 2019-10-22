package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var app AppConfig

type Log struct {
	level log.Level
}

type DiscoveryMap struct {
	EnvoyId               string `mapstructure:"envoy_id"`
	UpstreamEndpointLabel string `mapstructure:"upstream_endpoint_label"`
	Namespace             string `mapstructure:"namespace"`
}

type AppConfig struct {
	envoyConfig      EnvoyConfig
	discoveryMapping []DiscoveryMap
	refreshInterval  int
	Log
}

func MustLoad(name, path string) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	app.envoyConfig = loadEnvoyConfig()

	err = viper.UnmarshalKey("operator.envoy_discovery_mapping", &app.discoveryMapping)
	if err != nil {
		log.Fatalf("Error loading envoy discovery mapping config - %v\n", err)
	}

	logLevel := viper.GetString("operator.log.level")
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("Error loading logger config: %v", err)
	}
	app.Log = Log{
		level: level,
	}

	app.refreshInterval = viper.GetInt("operator.refresh_interval_in_s")
}

func GetEnvoyConfig() EnvoyConfig {
	return app.envoyConfig
}

func GetDiscoveryMapping() []DiscoveryMap {
	return app.discoveryMapping
}

func LogLevel() log.Level {
	return app.Log.level
}

func RefreshIntervalInS() int {
	return app.refreshInterval
}
