package controlplane_test

import (
	"testing"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/stretchr/testify/assert"

	cp "github.com/gojekfarm/envoy-lb-operator/controlplane"
)

func TestTCPAddress(t *testing.T) {
	addr := cp.TCPAddress("foo", 443)
	socketAddress := addr.Address.(*core.Address_SocketAddress).SocketAddress
	assert.Equal(t, "foo", socketAddress.Address)
	assert.Equal(t, core.TCP, socketAddress.Protocol)
	assert.Equal(t, uint32(443), socketAddress.PortSpecifier.(*core.SocketAddress_PortValue).PortValue)
}
