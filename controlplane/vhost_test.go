package controlplane_test

import (
	"testing"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	cp "github.com/gojekfarm/envoy-lb-operator/controlplane"
	"github.com/stretchr/testify/assert"
)

func TestVHost(t *testing.T) {
	vhost := cp.VHost("foo", []string{"*"}, []cp.Target{
		{
			Host:        "foo",
			Regex:       "/foo",
			ClusterName: "foo_cluster",
		},
		{
			Host:        "bar",
			Regex:       "/bar",
			ClusterName: "bar_cluster",
		},
	})

	assert.Equal(t, "foo", vhost.Name)
	assert.Equal(t, 1, len(vhost.Domains))
	assert.Equal(t, "*", vhost.Domains[0])
	assert.Equal(t, 2, len(vhost.Routes))
	assert.Equal(t, "/foo", vhost.Routes[0].Match.PathSpecifier.(*route.RouteMatch_Regex).Regex)
	assert.Equal(t, "foo", vhost.Routes[0].Action.(*route.Route_Route).Route.HostRewriteSpecifier.(*route.RouteAction_HostRewrite).HostRewrite)
	assert.Equal(t, "foo_cluster", vhost.Routes[0].Action.(*route.Route_Route).Route.ClusterSpecifier.(*route.RouteAction_Cluster).Cluster)

	assert.Equal(t, "/bar", vhost.Routes[1].Match.PathSpecifier.(*route.RouteMatch_Regex).Regex)
	assert.Equal(t, "bar", vhost.Routes[1].Action.(*route.Route_Route).Route.HostRewriteSpecifier.(*route.RouteAction_HostRewrite).HostRewrite)
	assert.Equal(t, "bar_cluster", vhost.Routes[1].Action.(*route.Route_Route).Route.ClusterSpecifier.(*route.RouteAction_Cluster).Cluster)

}
