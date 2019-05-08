package kube

import (
	"fmt"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
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
}

func (s Service) clusterName() string {
	return fmt.Sprintf("%s_cluster", s.Address)
}

//Cluster returns envoy control plane config for a headless strict dns lookup
func (s Service) Cluster() *v2.Cluster {
	if s.Type == GRPC {
		return cp.StrictDNSLRHttp2Cluster(s.clusterName(), s.Address, s.Port, 1000)
	}
	return cp.StrictDNSLRCluster(s.clusterName(), s.Address, s.Port, 1000)
}

//Target represents the vhost target
func (s Service) Target(prefix string) cp.Target {
	return cp.Target{Host: s.Address, Prefix: prefix, ClusterName: s.clusterName()}
}

//DefaultTarget represents the vhost target
func (s Service) DefaultTarget() cp.Target {
	return cp.Target{Host: s.Address, Prefix: "/", ClusterName: s.clusterName()}
}

func ServiceType(svc *corev1.Service) Type {
	serviceTypeAnnotation := svc.GetAnnotations()["envoy-lb-operator.gojektech.k8s.io/service-type"]
	if serviceTypeAnnotation == "grpc" {
		return GRPC
	}
	return HTTP

}
