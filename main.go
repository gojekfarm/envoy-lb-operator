package main

import (
	"github.com/gojekfarm/envoy-lb-operator/cmd"
	"github.com/gojekfarm/envoy-lb-operator/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	config.MustLoad("application", "./")
	log.SetLevel(config.LogLevel())
	cmd.Execute()
}
