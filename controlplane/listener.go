package controlplane

import (
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
)

//Listener for given address and connectionManager
func Listener(name, address string, port uint32, connectionManager *hcm.HttpConnectionManager) (*v2.Listener, error) {
	connectionManagerAny, err := types.MarshalAny(connectionManager)
	if err != nil {
		return nil, err
	}
	return &v2.Listener{
		Name:    name,
		Address: TCPAddress(address, port),
		FilterChains: []listener.FilterChain{{
			Filters: []listener.Filter{{
				Name:       util.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{TypedConfig: connectionManagerAny},
			}},
		}},
	}, nil
}
