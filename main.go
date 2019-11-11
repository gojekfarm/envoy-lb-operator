package main

import (
	"github.com/gojekfarm/envoy-lb-operator/cmd"
	"github.com/gojekfarm/envoy-lb-operator/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := config.MustLoad("application", "./")
	if err != nil {
		log.Fatalf("Error while loading config - %v\n", err)
	}
	log.SetLevel(config.LogLevel())
	cmd.Execute()
}
