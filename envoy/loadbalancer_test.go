package envoy_test

import (
	"encoding/json"
	"github.com/gojekfarm/envoy-lb-operator/config"
	"testing"

	"github.com/gojekfarm/envoy-lb-operator/envoy"
	kube "github.com/gojekfarm/envoy-lb-operator/kube"
	"github.com/stretchr/testify/assert"
)

func TestSnapshotVersion(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{})
	assert.Equal(t, int32(0), lb.ConfigVersion)
}

func TestSnapshotVersionIncrementsOnStore(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{})
	assert.Equal(t, int32(0), lb.ConfigVersion)
	lb.Snapshot()
	assert.Equal(t, int32(1), lb.ConfigVersion)
	lb.Snapshot()
	assert.Equal(t, int32(2), lb.ConfigVersion)
}

func TestInitialState(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{})
	lb.Snapshot()
	sn, _ := lb.Config.GetSnapshot("node1")
	assert.Equal(t, 1, len(sn.Listeners.Items))
	assert.Equal(t, 0, len(sn.Clusters.Items))
}

func TestAddedUpstream(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{})
	lb.Snapshot()
	lb.Trigger(envoy.LBEvent{
		Svc:       kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC, Path: "/foo", Domain: "*"},
		EventType: envoy.ADDED,
	})
	lb.Close()
	lb.HandleEvents()
	lb.Snapshot()
	sn, _ := lb.Config.GetSnapshot("node1")
	assert.Equal(t, 1, len(sn.Listeners.Items))
	assert.Equal(t, 1, len(sn.Clusters.Items))
}

func TestAddUpdatedUpstream(t *testing.T) {
	envoyConfig, _ := config.LoadEnvoyConfig("application", "../")
	lb := envoy.NewLB("node1", envoyConfig)
	lb.Snapshot()
	lb.Trigger(envoy.LBEvent{
		Svc:       kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC, Path: "/foo", Domain: "*"},
		EventType: envoy.ADDED,
	})
	lb.Trigger(envoy.LBEvent{
		Svc:       kube.Service{Address: "foo", Port: uint32(8001), Type: kube.GRPC, Path: "/foo", Domain: "*"},
		EventType: envoy.UPDATED,
	})
	lb.Close()
	lb.HandleEvents()
	lb.Snapshot()
	sn, _ := lb.Config.GetSnapshot("node1")
	assert.Equal(t, 1, len(sn.Listeners.Items))
	assert.Equal(t, 1, len(sn.Clusters.Items))
	cfg, _ := json.Marshal(sn.Clusters.Items["foo_cluster"])
	assert.Equal(t, `{"name":"foo_cluster","ClusterDiscoveryType":{"Type":1},"connect_timeout":1000000000,"lb_policy":1,"load_assignment":{"cluster_name":"foo_cluster","endpoints":[{"lb_endpoints":[{"HostIdentifier":{"Endpoint":{"address":{"Address":{"SocketAddress":{"address":"foo","PortSpecifier":{"PortValue":8001}}}}}}}]}]},"circuit_breakers":{"thresholds":[{"max_connections":{"value":1024},"max_pending_requests":{"value":50000},"max_requests":{"value":50000},"max_retries":{"value":50000}}]},"http2_protocol_options":{},"dns_lookup_family":1,"outlier_detection":{"consecutive_5xx":{"value":10000},"interval":{"seconds":10},"base_ejection_time":{"seconds":30},"max_ejection_percent":{"value":50},"enforcing_consecutive_5xx":{},"consecutive_gateway_failure":{"value":5},"enforcing_consecutive_gateway_failure":{"value":100}},"LbConfig":null}`, string(cfg))
}

func TestDeletedUpstream(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{})
	lb.Snapshot()
	lb.Trigger(envoy.LBEvent{
		Svc:       kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC, Path: "/foo", Domain: "*"},
		EventType: envoy.ADDED,
	})
	lb.Trigger(envoy.LBEvent{
		Svc:       kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC, Path: "/foo", Domain: "*"},
		EventType: envoy.DELETED,
	})
	lb.Close()
	lb.HandleEvents()
	lb.Snapshot()
	sn, _ := lb.Config.GetSnapshot("node1")
	assert.Equal(t, 1, len(sn.Listeners.Items))
	assert.Equal(t, 0, len(sn.Clusters.Items))
}

func TestSingleVhostDifferentPaths(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{})
	lb.Snapshot()
	lb.Trigger(envoy.LBEvent{
		Svc:       kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC, Path: "/foo", Domain: "*"},
		EventType: envoy.ADDED,
	})
	lb.Trigger(envoy.LBEvent{
		Svc:       kube.Service{Address: "bar", Port: uint32(8000), Type: kube.GRPC, Path: "/bar", Domain: "*"},
		EventType: envoy.ADDED,
	})
	lb.Close()
	lb.HandleEvents()
	lb.Snapshot()
	sn, _ := lb.Config.GetSnapshot("node1")
	assert.Equal(t, 1, len(sn.Listeners.Items))
	//No Easy way to assert
	//cfg, _ := json.Marshal(sn.Listeners.Items["assert.Equal(t, "", string(cfg))
}

func TestMultipleVhostsDifferentPaths(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{})
	lb.Snapshot()
	lb.Trigger(envoy.LBEvent{
		Svc:       kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC, Path: "/", Domain: "foo.abc.com"},
		EventType: envoy.ADDED,
	})
	lb.Trigger(envoy.LBEvent{
		Svc:       kube.Service{Address: "bar", Port: uint32(8000), Type: kube.GRPC, Path: "/", Domain: "bar.abc.com"},
		EventType: envoy.ADDED,
	})
	lb.Close()
	lb.HandleEvents()
	lb.Snapshot()
	sn, _ := lb.Config.GetSnapshot("node1")
	assert.Equal(t, 1, len(sn.Listeners.Items))
	//No Easy way to assert
	//cfg, _ := json.Marshal(sn.Listeners.Items["assert.Equal(t, "", string(cfg))
}
