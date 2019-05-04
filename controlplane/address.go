package controlplane

import "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

//TCPAddress builds a envoy tcp address
func TCPAddress(address string, port uint32) core.Address {
	return core.Address{Address: &core.Address_SocketAddress{
		SocketAddress: &core.SocketAddress{
			Address:  address,
			Protocol: core.TCP,
			PortSpecifier: &core.SocketAddress_PortValue{
				PortValue: port,
			},
		},
	}}
}
