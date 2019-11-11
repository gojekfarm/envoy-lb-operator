package envoy

import (
	"fmt"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/gojekfarm/envoy-lb-operator/config"
	cp "github.com/gojekfarm/envoy-lb-operator/controlplane"
	"github.com/gojekfarm/envoy-lb-operator/kube"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
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
	events          chan LBEvent
	upstreams       map[string]kube.Service
	nodeID          string
	cache           cache.SnapshotCache
	cacheVersion    int32
	envoyConfig     config.EnvoyConfig
	autoRefreshConn bool
}

func (lb *LoadBalancer) Trigger(evt LBEvent) {
	lb.events <- evt
}

func (lb *LoadBalancer) InitializeUpstream(serviceList *corev1.ServiceList) {
	lb.incrementVersion()
	for _, service := range serviceList.Items {
		svc := lb.getService(&service)
		lb.upstreams[svc.Address] = svc
	}
	log.Debug("Populated all existing upstreams.")
}

func (lb *LoadBalancer) SvcTrigger(eventType LBEventType, svc *corev1.Service) {
	log.Debugf("Received event: %s eventtype: %+v for node: %s", svc, eventType, lb.nodeID)
	if svc.Spec.ClusterIP == v1.ClusterIPNone {
		lb.Trigger(LBEvent{EventType: eventType, Svc: lb.getService(svc)})
	}
}

func (lb *LoadBalancer) Close() {
	log.Debug("Closing lb operator")
	close(lb.events)
}

func (lb *LoadBalancer) HandleEvents() {
	for evt := range lb.events {
		lb.incrementVersion()
		switch evt.EventType {
		case DELETED:
			delete(lb.upstreams, evt.Svc.Address)
			log.Debugf("Deleting upstream: %v id: %s", evt.Svc, lb.nodeID)
		default:
			lb.upstreams[evt.Svc.Address] = evt.Svc
			log.Debugf("Adding upstream: %v id: %s\n", evt.Svc, lb.nodeID)
		}
	}
}

func (lb *LoadBalancer) EndpointTrigger() {
	lb.incrementVersion()
}

func (lb *LoadBalancer) SnapshotRunner() {
	log.Debug("Executing Snapshot Runner...")
	if lb.autoRefreshConn {
		lb.incrementVersion()
	}
	var clusters []cache.Resource

	targetsByDomain := make(map[string][]cp.Target)
	if len(lb.upstreams) > 0 {
		for _, svc := range lb.upstreams {

			clusters = append(clusters, svc.Cluster(lb.envoyConfig.ConnectTimeoutMs, lb.envoyConfig.CircuitBreaker, lb.envoyConfig.OutlierDetection))
			if targetsByDomain[svc.Domain] == nil {
				targetsByDomain[svc.Domain] = []cp.Target{svc.DefaultTarget()}
			} else {
				targetsByDomain[svc.Domain] = append(targetsByDomain[svc.Domain], svc.DefaultTarget())
			}
		}
	}
	vhosts := []route.VirtualHost{}
	for domain, targets := range targetsByDomain {
		retryConfig := lb.envoyConfig.RetryConfig
		vhosts = append(vhosts, cp.VHost(fmt.Sprintf("local_service_%s", domain), []string{domain}, targets, cp.RetryPolicy(retryConfig.RetryOn, retryConfig.RetryPredicate, retryConfig.NumRetries, retryConfig.HostSelectionMaxRetryAttempts)))
	}

	drainTimeoutInMs := 10 * time.Millisecond
	cm := cp.ConnectionManager("local_route", vhosts, &drainTimeoutInMs)
	var listener, err = cp.Listener("listener_grpc", "0.0.0.0", 80, cm)

	if err != nil {
		log.Errorf("Error %v", err)
		panic(err)
	}
	snapshot := cache.NewSnapshot(fmt.Sprint(lb.cacheVersion), nil, clusters, nil, []cache.Resource{listener})
	err = lb.cache.SetSnapshot(lb.nodeID, snapshot)
	if err != nil {
		log.Errorf("snapshot error: %s", err.Error())
	}
}

func NewLB(nodeID string, envoyConfig config.EnvoyConfig, snapshotCache cache.SnapshotCache, autoRefreshConn bool) *LoadBalancer {
	return &LoadBalancer{events: make(chan LBEvent, 10), upstreams: make(map[string]kube.Service), nodeID: nodeID, cache: snapshotCache, envoyConfig: envoyConfig, autoRefreshConn: autoRefreshConn}
}

func (lb *LoadBalancer) incrementVersion() {
	atomic.AddInt32(&lb.cacheVersion, 1)
	log.Infof("Incrementing snapshot version to %v\n", lb.cacheVersion)
}

func (lb *LoadBalancer) getService(svc *corev1.Service) kube.Service {
	return kube.Service{Address: svc.Name, Port: uint32(svc.Spec.Ports[0].TargetPort.IntVal), Type: kube.ServiceType(svc), Path: kube.ServicePath(svc), Domain: kube.ServiceDomain(svc)}
}

func (lb *LoadBalancer) GetCacheVersion() int32 {
	return lb.cacheVersion
}

func (lb *LoadBalancer) GetCache() cache.SnapshotCache {
	return lb.cache
}

func (lb *LoadBalancer) GetUpstreams() map[string]kube.Service {
	return lb.upstreams
}