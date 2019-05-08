package controlplane_test

import (
	"testing"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
	cp "github.com/gojekfarm/envoy-lb-operator/controlplane"
	"github.com/stretchr/testify/assert"
)

func TestListener(t *testing.T) {

	vhosts := []route.VirtualHost{cp.VHost("foo", []string{"*"}, []cp.Target{
		{
			Host:        "foo",
			Prefix:      "/foo",
			ClusterName: "foo_cluster",
		},
	})}

	cm := cp.ConnectionManager("route1234", vhosts)

	l, err := cp.Listener("foo", "0.0.0.0", uint32(9000), cm)
	assert.NoError(t, err)
	assert.NotNil(t, l)
	assert.Equal(t, "foo", l.Name)
	socketAddress := l.Address.Address.(*core.Address_SocketAddress).SocketAddress
	assert.Equal(t, "0.0.0.0", socketAddress.Address)
	assert.Equal(t, core.TCP, socketAddress.Protocol)
	assert.Equal(t, uint32(9000), socketAddress.PortSpecifier.(*core.SocketAddress_PortValue).PortValue)

	assert.Equal(t, 1, len(l.FilterChains))
	assert.Equal(t, util.HTTPConnectionManager, l.FilterChains[0].Filters[0].Name)
	anycm, _ := types.MarshalAny(cm)
	assert.Equal(t, anycm, l.FilterChains[0].Filters[0].ConfigType.(*listener.Filter_TypedConfig).TypedConfig)
}
