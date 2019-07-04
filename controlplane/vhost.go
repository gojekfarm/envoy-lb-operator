package controlplane

import (
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
)

func VHost(name string, domains []string, targets []Target, retryPolicy *route.RetryPolicy) route.VirtualHost {
	var routes []route.Route
	for _, t := range targets {
		routes = append(routes, t.Route())
	}

	var retryPredicateArray []*route.RetryPolicy_RetryHostPredicate
	retryPredicate := route.RetryPolicy_RetryHostPredicate{
		Name: "envoy.retry_host_predicates.previous_hosts",
	}
	retryPredicateArray = append(retryPredicateArray, &retryPredicate)
	return route.VirtualHost{
		Name:        name,
		Domains:     domains,
		Routes:      routes,
		RetryPolicy: retryPolicy,
	}
}
