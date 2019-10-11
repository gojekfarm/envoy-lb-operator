package config

import (
	"log"

	"github.com/spf13/viper"
)

var app AppConfig

type AppConfig struct {
	envoyConfig      EnvoyConfig
	discoveryMapping map[string]string
}

func MustLoad(name, path string) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	app.envoyConfig = loadEnvoyConfig()
	app.discoveryMapping = viper.GetStringMapString("operator.envoy_discovery_mapping")
}

func GetEnvoyConfig() EnvoyConfig {
	return app.envoyConfig
}

func GetDiscoveryMapping() map[string]string {
	return app.discoveryMapping
}
