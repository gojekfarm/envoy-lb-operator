package envoy

import (
	"fmt"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/gojekfarm/envoy-lb-operator/config"
	cp "github.com/gojekfarm/envoy-lb-operator/controlplane"
	"github.com/gojekfarm/envoy-lb-operator/kube"
	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"log"
	"sync"
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
	sync.RWMutex
	events        chan LBEvent
	upstreams     map[string]kube.Service
	nodeID        string
	Config        cache.SnapshotCache
	ConfigVersion int32
	EnvoyConfig   config.EnvoyConfig
}

func (lb *LoadBalancer) Trigger(evt LBEvent) {
	fmt.Printf("Before Adding event to channel - %v for node - %s\n", evt, lb.nodeID)
	lb.events <- evt
	fmt.Printf("After Adding event to channel - %v for node - %s\n", evt, lb.nodeID)
}

func (lb *LoadBalancer) SvcTrigger(eventType LBEventType, svc *corev1.Service) {
	log.Printf("!!! Received event: %s eventtype: %+v for node: %s", svc, eventType, lb.nodeID)
	if svc.Spec.ClusterIP == v1.ClusterIPNone {
		lb.Trigger(LBEvent{EventType: eventType, Svc: kube.Service{Address: svc.Name, Port: uint32(svc.Spec.Ports[0].TargetPort.IntVal), Type: kube.ServiceType(svc), Path: kube.ServicePath(svc), Domain: kube.ServiceDomain(svc)}})
	}
}

func (lb *LoadBalancer) Close() {
	fmt.Printf("!!!!!!! Closing lb")
	close(lb.events)
}

func (lb *LoadBalancer) HandleEvents() {
	fmt.Printf("loadbalancer %s handling events", lb.nodeID)
	for evt := range lb.events {
		switch evt.EventType {
		case DELETED:
			delete(lb.upstreams, evt.Svc.Address)
			fmt.Printf("deleting upstream: %v id: %s", evt.Svc, lb.nodeID)
		default:
			lb.upstreams[evt.Svc.Address] = evt.Svc
			fmt.Printf("adding upstream: %v id: %s\n", evt.Svc, lb.nodeID)
		}
	}
	fmt.Printf("!!!! Completed handling events: %s", lb.nodeID)
}

func (lb *LoadBalancer) SnapshotRunner() {
	fmt.Printf("Running snapshot runner.... %s\n", lb.nodeID)
	lb.RLock()
	defer lb.RUnlock()
	atomic.AddInt32(&lb.ConfigVersion, 1)
	var clusters []cache.Resource

	targetsByDomain := make(map[string][]cp.Target)
	if len(lb.upstreams) > 0 {
		for _, svc := range lb.upstreams {

			clusters = append(clusters, svc.Cluster(lb.EnvoyConfig.ConnectTimeoutMs, lb.EnvoyConfig.CircuitBreaker, lb.EnvoyConfig.OutlierDetection))
			if targetsByDomain[svc.Domain] == nil {
				targetsByDomain[svc.Domain] = []cp.Target{svc.DefaultTarget()}
			} else {
				targetsByDomain[svc.Domain] = append(targetsByDomain[svc.Domain], svc.DefaultTarget())
			}
		}
	}
	vhosts := []route.VirtualHost{}
	for domain, targets := range targetsByDomain {
		retryConfig := lb.EnvoyConfig.RetryConfig
		vhosts = append(vhosts, cp.VHost(fmt.Sprintf("local_service_%s", domain), []string{domain}, targets, cp.RetryPolicy(retryConfig.RetryOn, retryConfig.RetryPredicate, retryConfig.NumRetries, retryConfig.HostSelectionMaxRetryAttempts)))
	}

	drainTimeoutInMs := 10 * time.Millisecond
	cm := cp.ConnectionManager("local_route", vhosts, &drainTimeoutInMs)
	var listener, err = cp.Listener("listener_grpc", "0.0.0.0", 80, cm)

	if err != nil {
		fmt.Printf("Error %v", err)
		panic(err)
	}
	snapshot := cache.NewSnapshot(fmt.Sprint(lb.ConfigVersion), nil, clusters, nil, []cache.Resource{listener})
	err = lb.Config.SetSnapshot(lb.nodeID, snapshot)
	if err != nil {
		fmt.Printf("snapshot error: %s", err.Error())
	}
}

func NewLB(nodeID string, envoyConfig config.EnvoyConfig, snapshotCache cache.SnapshotCache) *LoadBalancer {
	return &LoadBalancer{events: make(chan LBEvent, 10), upstreams: make(map[string]kube.Service), nodeID: nodeID, Config: snapshotCache, EnvoyConfig: envoyConfig}
}
