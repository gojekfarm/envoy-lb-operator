package controlplane

import (
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
)

func VHost(name string, domains []string, targets []Target) route.VirtualHost {
	var routes []route.Route
	for _, t := range targets {
		routes = append(routes, t.Route())
	}
	return route.VirtualHost{
		Name:    name,
		Domains: domains,
		Routes:  routes,
	}

}
