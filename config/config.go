package config

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var app AppConfig

type Log struct {
	level log.Level
}

type DiscoveryMap struct {
	EnvoyId       string `mapstructure:"envoy_id"`
	UpstreamLabel string `mapstructure:"upstream_label"`
	EndpointLabel string `mapstructure:"endpoint_label"`
	Namespace     string `mapstructure:"namespace"`
}

type AppConfig struct {
	envoyConfig      EnvoyConfig
	discoveryMapping []DiscoveryMap
	refreshInterval  int
	autoRefreshConn  bool
	Log
}

func MustLoad(name, path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetDefault("operator.log.level", "info")
	viper.SetDefault("operator.refresh_interval_in_s", 10)
	viper.SetDefault("operator.auto_refresh_conn", false)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	app.envoyConfig = loadEnvoyConfig()

	err = viper.UnmarshalKey("operator.envoy_discovery_mapping", &app.discoveryMapping)
	if err != nil || app.discoveryMapping == nil {
		return errors.New(fmt.Sprintf("Error loading envoy discovery mapping config - %v",  err))
	}
	for _, mapping := range app.discoveryMapping {
		if mapping.EnvoyId == "" || mapping.Namespace == "" || mapping.UpstreamLabel == "" {
			return errors.New("invalid Configuration of envoy discovery mapping. Please check if envoy_id, namespace and upstream_label are configured")
		}
	}

	logLevel := viper.GetString("operator.log.level")
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return errors.New(fmt.Sprintf("Error loading logger config: %v", err))
	}
	app.Log = Log{
		level: level,
	}

	app.refreshInterval = viper.GetInt("operator.refresh_interval_in_s")
	app.autoRefreshConn = viper.GetBool("operator.auto_refresh_conn")
	return nil
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

func AutoRefreshConn() bool {
	return app.autoRefreshConn
}

func Clear() {
	app = AppConfig{}
}
