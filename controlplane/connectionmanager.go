package controlplane

import (
	"fmt"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/pkg/util"
)

//ConnectionManager for a given set of virtual hosts
func ConnectionManager(routeName string, vhosts []route.VirtualHost) *hcm.HttpConnectionManager {
	return &hcm.HttpConnectionManager{
		CodecType:  hcm.AUTO,
		StatPrefix: fmt.Sprintf("ingress_%s", routeName),
		RouteSpecifier: &hcm.HttpConnectionManager_RouteConfig{
			RouteConfig: &v2.RouteConfiguration{
				Name:         routeName,
				VirtualHosts: vhosts,
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name: util.Router,
		}},
	}
}
