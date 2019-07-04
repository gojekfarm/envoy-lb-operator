package controlplane_test

import (
	"testing"
	"time"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	cp "github.com/gojekfarm/envoy-lb-operator/controlplane"
	"github.com/stretchr/testify/assert"
)

func TestConnectionManager(t *testing.T) {

	vhosts := []route.VirtualHost{cp.VHost("foo", []string{"*"}, []cp.Target{
		{
			Host:        "foo",
			Prefix:      "/foo",
			ClusterName: "foo_cluster",
		},
	},
		cp.RetryPolicy("xxx", "retry_predicate", 10, 20),
	)}

	duration := 10 * time.Millisecond
	cm := cp.ConnectionManager("route1234", vhosts, &duration)
	assert.Equal(t, &duration, cm.DrainTimeout)
	assert.Equal(t, hcm.AUTO, cm.CodecType)
	assert.Equal(t, "ingress_route1234", cm.StatPrefix)
	assert.Equal(t, 1, len(cm.HttpFilters))
	assert.Equal(t, util.Router, cm.HttpFilters[0].Name)
	assert.Equal(t, "route1234", cm.RouteSpecifier.(*hcm.HttpConnectionManager_RouteConfig).RouteConfig.Name)
	assert.Equal(t, vhosts, cm.RouteSpecifier.(*hcm.HttpConnectionManager_RouteConfig).RouteConfig.VirtualHosts)
}
