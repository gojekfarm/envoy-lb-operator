package envoy

import (
	"fmt"
	"sync/atomic"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	cp "github.com/gojekfarm/envoy-lb-operator/controlplane"
	kube "github.com/gojekfarm/envoy-lb-operator/kube"

	corev1 "k8s.io/api/core/v1"
)

//LBEventType is the type of event impacting the LB
type LBEventType int32

const (
	//ADDED represents a service create event
	ADDED LBEventType = iota
	//UPDATED represents a service updated event
	UPDATED
	//DELETED represents a service updated event
	DELETED
)

//LBEvent is the event triggered by kubernetes service changes
type LBEvent struct {
	Svc       kube.Service
	EventType LBEventType
}

//LoadBalancer represents the current state of upstreams for a load balancer
type LoadBalancer struct {
	events        chan LBEvent
	upstreams     map[string]kube.Service
	nodeID        string
	Config        cache.SnapshotCache
	ConfigVersion int32
}

func (lb *LoadBalancer) Trigger(evt LBEvent) {
	lb.events <- evt
}

func (lb *LoadBalancer) SvcTrigger(eventType LBEventType, svc *corev1.Service) {
	lb.Trigger(LBEvent{EventType: eventType, Svc: kube.Service{Address: svc.Name, Port: uint32(svc.Spec.Ports[0].TargetPort.IntVal), Type: kube.ServiceType(svc)}})
}

func (lb *LoadBalancer) Close() {
	close(lb.events)
}

func (lb *LoadBalancer) HandleEvents() {
	for evt := range lb.events {
		switch evt.EventType {
		case DELETED:
			delete(lb.upstreams, evt.Svc.Address)
		default:
			lb.upstreams[evt.Svc.Address] = evt.Svc
		}
	}
}

func (lb *LoadBalancer) Snapshot() {
	atomic.AddInt32(&lb.ConfigVersion, 1)
	// svc := kube.Service{Address: "svc", Port: uint32(443), Type: kube.GRPC}
	var targets []cp.Target
	var clusters []cache.Resource
	for _, svc := range lb.upstreams {
		clusters = append(clusters, svc.Cluster())
		targets = append(targets, svc.DefaultTarget())
	}

	vh := cp.VHost("local_service", []string{"*"}, targets)
	cm := cp.ConnectionManager("local_route", []route.VirtualHost{vh})
	var l, err = cp.Listener("listener_grpc", "0.0.0.0", 8080, cm)

	if err != nil {
		panic(err)
	}

	snapshot := cache.NewSnapshot(fmt.Sprint(lb.ConfigVersion), nil, clusters, nil, []cache.Resource{l})

	lb.Config.SetSnapshot(lb.nodeID, snapshot)
}

func NewLB(nodeID string) *LoadBalancer {
	return &LoadBalancer{events: make(chan LBEvent, 10), upstreams: make(map[string]kube.Service), nodeID: nodeID, Config: cache.NewSnapshotCache(true, Hasher{}, logger{})}
}
