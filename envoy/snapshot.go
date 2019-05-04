package envoy

import (
	"fmt"
	"sync/atomic"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	cp "github.com/gojekfarm/envoy-lb-operator/controlplane"
)

//Snapshot represents a snapshot of envoy config
//This is backed by a cache
type Snapshot struct {
	nodeID  string
	Config  cache.SnapshotCache
	Version int32
}

func (snap *Snapshot) Store() {
	atomic.AddInt32(&snap.Version, 1)
	vh := cp.VHost("local_service", []string{"*"}, []cp.Target{{Host: "svc", Regex: "/", ClusterName: "woot"}})
	cm := cp.ConnectionManager("local_route", []route.VirtualHost{vh})
	var l, err = cp.Listener("listener_grpc", "0.0.0.0", 8080, cm)

	if err != nil {
		panic(err)
	}

	snapshot := cache.NewSnapshot(fmt.Sprint(snap.Version), nil, []cache.Resource{cp.StrictDNSLRCluster("woot", "svc", uint32(443), 1000)}, nil, []cache.Resource{l})

	snap.Config.SetSnapshot(snap.nodeID, snapshot)

}

func NewSnapshot(nodeID string) *Snapshot {
	return &Snapshot{nodeID: nodeID, Config: cache.NewSnapshotCache(true, Hasher{}, logger{})}
}
