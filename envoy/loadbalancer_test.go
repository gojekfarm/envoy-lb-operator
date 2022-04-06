package envoy_test

import (
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	"testing"

	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/gojekfarm/envoy-lb-operator/config"

	"github.com/gojekfarm/envoy-lb-operator/envoy"
	"github.com/gojekfarm/envoy-lb-operator/kube"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSnapshotVersion(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{}, cache.NewSnapshotCache(true, envoy.Hasher{}, envoy.Logger{}), false)
	assert.Equal(t, int32(0), lb.GetCacheVersion())
}

func TestSnapshotVersionIncrementsOnHandleEvents(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{}, cache.NewSnapshotCache(true, envoy.Hasher{}, envoy.Logger{}), false)
	assert.Equal(t, int32(0), lb.GetCacheVersion())
	lb.Trigger(envoy.LBEvent{
		Svc:       kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC, Path: "/foo", Domain: "*"},
		EventType: envoy.ADDED,
	})
	lb.Close()
	lb.HandleEvents()
	assert.Equal(t, int32(1), lb.GetCacheVersion())
}

func TestSnapshotVersionDoesNotIncrementOnSnapshotRunnerIfAutoRefreshIsDisabled(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{}, cache.NewSnapshotCache(true, envoy.Hasher{}, envoy.Logger{}), false)
	assert.Equal(t, int32(0), lb.GetCacheVersion())
	lb.SnapshotRunner()
	assert.Equal(t, int32(0), lb.GetCacheVersion())
}

func TestSnapshotVersionIncrementsOnSnapshotRunnerIfAutoRefreshIsEnabled(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{}, cache.NewSnapshotCache(true, envoy.Hasher{}, envoy.Logger{}), true)
	assert.Equal(t, int32(0), lb.GetCacheVersion())
	lb.SnapshotRunner()
	assert.Equal(t, int32(1), lb.GetCacheVersion())
}

func TestSnapshotVersionIncrementsOnEndpointTrigger(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{}, cache.NewSnapshotCache(true, envoy.Hasher{}, envoy.Logger{}), false)
	assert.Equal(t, int32(0), lb.GetCacheVersion())
	lb.EndpointTrigger()
	assert.Equal(t, int32(1), lb.GetCacheVersion())
}

func TestInitialState(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{}, cache.NewSnapshotCache(true, envoy.Hasher{}, envoy.Logger{}), false)
	lb.SnapshotRunner()
	sn, _ := lb.GetCache().GetSnapshot("node1")
	assert.Equal(t, 1, len(sn.Listeners.Items))
	assert.Equal(t, 0, len(sn.Clusters.Items))
}

func TestAddedUpstream(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{}, cache.NewSnapshotCache(true, envoy.Hasher{}, envoy.Logger{}), false)
	lb.SnapshotRunner()
	lb.Trigger(envoy.LBEvent{
		Svc:       kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC, Path: "/foo", Domain: "*"},
		EventType: envoy.ADDED,
	})
	lb.Close()
	lb.HandleEvents()
	lb.SnapshotRunner()
	sn, _ := lb.GetCache().GetSnapshot("node1")
	assert.Equal(t, 1, len(sn.Listeners.Items))
	assert.Equal(t, 1, len(sn.Clusters.Items))
}

func TestAddUpdatedUpstream(t *testing.T) {
	config.MustLoad("application", "../")
	envoyConfig := config.GetEnvoyConfig()
	lb := envoy.NewLB("node1", envoyConfig, cache.NewSnapshotCache(true, envoy.Hasher{}, envoy.Logger{}), false)
	lb.SnapshotRunner()
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
	lb.SnapshotRunner()
	sn, _ := lb.GetCache().GetSnapshot("node1")
	assert.Equal(t, 1, len(sn.Listeners.Items))
	assert.Equal(t, 1, len(sn.Clusters.Items))
	cfg, _ := json.Marshal(sn.Clusters.Items["foo_cluster"])
	assert.Equal(t, `{"name":"foo_cluster","ClusterDiscoveryType":{"Type":1},"connect_timeout":1000000000,"lb_policy":1,"load_assignment":{"cluster_name":"foo_cluster","endpoints":[{"lb_endpoints":[{"HostIdentifier":{"Endpoint":{"address":{"Address":{"SocketAddress":{"address":"foo","PortSpecifier":{"PortValue":8001}}}}}}}]}]},"circuit_breakers":{"thresholds":[{"max_connections":{"value":1024},"max_pending_requests":{"value":50000},"max_requests":{"value":50000},"max_retries":{"value":50000}}]},"http2_protocol_options":{},"dns_lookup_family":1,"outlier_detection":{"consecutive_5xx":{"value":10000},"interval":{"seconds":10},"base_ejection_time":{"seconds":30},"max_ejection_percent":{"value":50},"enforcing_consecutive_5xx":{},"consecutive_gateway_failure":{"value":5},"enforcing_consecutive_gateway_failure":{"value":100}},"LbConfig":null}`, string(cfg))
}

func TestDeletedUpstream(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{}, cache.NewSnapshotCache(true, envoy.Hasher{}, envoy.Logger{}), false)
	lb.SnapshotRunner()
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
	lb.SnapshotRunner()
	sn, _ := lb.GetCache().GetSnapshot("node1")
	assert.Equal(t, 1, len(sn.Listeners.Items))
	assert.Equal(t, 0, len(sn.Clusters.Items))
}

func TestSingleVhostDifferentPaths(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{}, cache.NewSnapshotCache(true, envoy.Hasher{}, envoy.Logger{}), false)
	lb.SnapshotRunner()
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
	lb.SnapshotRunner()
	sn, _ := lb.GetCache().GetSnapshot("node1")
	assert.Equal(t, 1, len(sn.Listeners.Items))
	//No Easy way to assert
	//cfg, _ := json.Marshal(sn.Listeners.Items["assert.Equal(t, "", string(cfg))
}

func TestMultipleVhostsDifferentPaths(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{}, cache.NewSnapshotCache(true, envoy.Hasher{}, envoy.Logger{}), false)
	lb.SnapshotRunner()
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
	lb.SnapshotRunner()
	sn, _ := lb.GetCache().GetSnapshot("node1")
	assert.Equal(t, 1, len(sn.Listeners.Items))
	//No Easy way to assert
	//cfg, _ := json.Marshal(sn.Listeners.Items["assert.Equal(t, "", string(cfg))
}

func TestInitializeMultipleUpstreamsOnStart(t *testing.T) {
	lb := envoy.NewLB("node1", config.EnvoyConfig{}, cache.NewSnapshotCache(true, envoy.Hasher{}, envoy.Logger{}), false)
	svc1 := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Port: int32(1234),
			}},
		},
	}

	svc2 := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "bar",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Port: int32(1234),
			}},
		},
	}
	svcList := &corev1.ServiceList{Items: []corev1.Service{svc1, svc2}}
	expectedSvc1 := kube.Service(kube.Service{Address: "foo", Port: 0x0, Type: 0, Path: "/", Domain: "*"})
	expectedSvc2 := kube.Service(kube.Service{Address: "bar", Port: 0x0, Type: 0, Path: "/", Domain: "*"})

	lb.InitializeUpstream(svcList)

	assert.Equal(t, int32(1), lb.GetCacheVersion())
	assert.Equal(t, 2, len(lb.GetUpstreams()))
	assert.Equal(t, expectedSvc1, lb.GetUpstreams()["foo"])
	assert.Equal(t, expectedSvc2, lb.GetUpstreams()["bar"])
}
