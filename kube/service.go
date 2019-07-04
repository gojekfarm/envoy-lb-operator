package kube

import (
	"fmt"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	cp "github.com/gojekfarm/envoy-lb-operator/controlplane"

	corev1 "k8s.io/api/core/v1"
)

//Type is the type of service
type Type int32

const (
	//HTTP are plain old http services
	HTTP Type = iota
	//GRPC (see https://grpc.io/)
	GRPC
)

//Service represents a headless k8s service that needs a loadbalancer
type Service struct {
	Address string
	Port    uint32
	Type    Type
	Path    string
	Domain  string
}

func (s Service) clusterName() string {
	return fmt.Sprintf("%s_cluster", s.Address)
}

//Cluster returns envoy control plane config for a headless strict dns lookup
func (s Service) Cluster() *v2.Cluster {
	circuitBreaker := cp.CircuitBreaker(1024, 50000, 50000, 50000)
	outlierDetection := cp.OutlierDetection(30, 10, 10000, 5, 0, 100, 50)
	if s.Type == GRPC {
		return cp.StrictDNSLRHttp2Cluster(s.clusterName(), s.Address, s.Port, 1000, circuitBreaker, outlierDetection)
	}
	return cp.StrictDNSLRCluster(s.clusterName(), s.Address, s.Port, 1000, circuitBreaker, outlierDetection)
}

//DefaultTarget represents the vhost target
func (s Service) DefaultTarget() cp.Target {
	return cp.Target{Host: s.Address, Prefix: s.Path, ClusterName: s.clusterName()}
}

func ServiceType(svc *corev1.Service) Type {
	serviceTypeAnnotation := svc.GetAnnotations()["envoy-lb-operator.gojektech.k8s.io/service-type"]
	if serviceTypeAnnotation == "grpc" {
		return GRPC
	}
	return HTTP

}

func ServicePath(svc *corev1.Service) string {
	servicePath := svc.GetAnnotations()["envoy-lb-operator.gojektech.k8s.io/service-path"]
	if servicePath == "" {
		return "/"
	}
	return servicePath
}

func ServiceDomain(svc *corev1.Service) string {
	serviceDomain := svc.GetAnnotations()["envoy-lb-operator.gojektech.k8s.io/service-domain"]
	if serviceDomain == "" {
		return "*"
	}
	return serviceDomain
}
