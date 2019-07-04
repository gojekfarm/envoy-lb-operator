package envoy

import (
	"fmt"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	cp "github.com/gojekfarm/envoy-lb-operator/controlplane"
	"github.com/gojekfarm/envoy-lb-operator/kube"
	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"sync/atomic"
	"time"
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
	if svc.Spec.ClusterIP == v1.ClusterIPNone {
		lb.Trigger(LBEvent{EventType: eventType, Svc: kube.Service{Address: svc.Name, Port: uint32(svc.Spec.Ports[0].TargetPort.IntVal), Type: kube.ServiceType(svc), Path: kube.ServicePath(svc), Domain: kube.ServiceDomain(svc)}})
	}
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
	var clusters []cache.Resource

	targetsByDomain := make(map[string][]cp.Target)
	if len(lb.upstreams) > 0 {
		for _, svc := range lb.upstreams {

			clusters = append(clusters, svc.Cluster())
			if targetsByDomain[svc.Domain] == nil {
				targetsByDomain[svc.Domain] = []cp.Target{svc.DefaultTarget()}
			} else {
				targetsByDomain[svc.Domain] = append(targetsByDomain[svc.Domain], svc.DefaultTarget())
			}
		}
	}
	vhosts := []route.VirtualHost{}
	for domain, targets := range targetsByDomain {
		vhosts = append(vhosts, cp.VHost(fmt.Sprintf("local_service_%s", domain), []string{domain}, targets, cp.RetryPolicy("connect-failure", "envoy.retry_host_predicates.previous_hosts", 3, 3)))
	}

	drainTimeoutInMs := 10 * time.Millisecond
	cm := cp.ConnectionManager("local_route", vhosts, &drainTimeoutInMs)
	var listener, err = cp.Listener("listener_grpc", "0.0.0.0", 80, cm)

	if err != nil {
		panic(err)
	}
	snapshot := cache.NewSnapshot(fmt.Sprint(lb.ConfigVersion), nil, clusters, nil, []cache.Resource{listener})
	lb.Config.SetSnapshot(lb.nodeID, snapshot)
}

func NewLB(nodeID string) *LoadBalancer {
	return &LoadBalancer{events: make(chan LBEvent, 10), upstreams: make(map[string]kube.Service), nodeID: nodeID, Config: cache.NewSnapshotCache(true, Hasher{}, logger{})}
}
