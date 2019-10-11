package main

import (
	"github.com/gojekfarm/envoy-lb-operator/cmd"
	"github.com/gojekfarm/envoy-lb-operator/config"
)

func main() {
	config.MustLoad("application", "./")
	cmd.Execute()
}
