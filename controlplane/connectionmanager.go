package controlplane

import (
	"fmt"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"time"
)

//ConnectionManager for a given set of virtual hosts
func ConnectionManager(routeName string, vhosts []route.VirtualHost, drainTimeout *time.Duration) *hcm.HttpConnectionManager {
	return &hcm.HttpConnectionManager{
		CodecType:    hcm.AUTO,
		DrainTimeout: drainTimeout,
		StatPrefix:   fmt.Sprintf("ingress_%s", routeName),
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
