package controlplane_test

import (
	"testing"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	cp "github.com/gojekfarm/envoy-lb-operator/controlplane"
	"github.com/stretchr/testify/assert"
)

func TestTargetRoute(t *testing.T) {
	target := &cp.Target{Host: "test", Regex: "/", ClusterName: "cluster_a"}
	r := target.Route()
	assert.Equal(t, "/", r.Match.PathSpecifier.(*route.RouteMatch_Regex).Regex)
	assert.Equal(t, "test", r.Action.(*route.Route_Route).Route.HostRewriteSpecifier.(*route.RouteAction_HostRewrite).HostRewrite)
	assert.Equal(t, "cluster_a", r.Action.(*route.Route_Route).Route.ClusterSpecifier.(*route.RouteAction_Cluster).Cluster)
}
