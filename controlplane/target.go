package controlplane

import (
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"strings"
)

func routeMatch(prefix string) route.RouteMatch {
	return route.RouteMatch{
		PathSpecifier: &route.RouteMatch_Prefix{
			Prefix: prefix,
		},
	}
}

func routeAction(target, cluster string) *route.Route_Route {
	return &route.Route_Route{
		Route: &route.RouteAction{
			HostRewriteSpecifier: &route.RouteAction_HostRewrite{
				HostRewrite: target,
			},
			ClusterSpecifier: &route.RouteAction_Cluster{
				Cluster: cluster,
			},
		},
	}
}

//Target represents a routing target criteria
type Target struct {
	Host        string
	Prefix      string
	ClusterName string
}

//Route is the route for the current target
func (t *Target) Route() []route.Route {
	var routes []route.Route
	arr := strings.Split(t.Prefix, ",")
	for _, p := range arr {
		if len(p) > 0 {
			routes = append(routes, route.Route{
				Match:  routeMatch(p),
				Action: routeAction(t.Host, t.ClusterName),
			})
		}
	}
	return routes
}
