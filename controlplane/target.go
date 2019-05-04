package controlplane

import "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"

func routeMatch(regex string) route.RouteMatch {
	return route.RouteMatch{
		PathSpecifier: &route.RouteMatch_Regex{
			Regex: regex,
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
	Regex       string
	ClusterName string
}

//Route is the route for the current target
func (t *Target) Route() route.Route {
	return route.Route{
		Match:  routeMatch(t.Regex),
		Action: routeAction(t.Host, t.ClusterName),
	}

}
