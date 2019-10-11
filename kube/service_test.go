package kube_test

import (
	"testing"

	"github.com/gojekfarm/envoy-lb-operator/config"

	"github.com/gojekfarm/envoy-lb-operator/kube"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestServiceClusterName(t *testing.T) {
	config.MustLoad("application", "../")
	envoyConfig := config.GetEnvoyConfig()
	grpccl := kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC, Path: "/"}.Cluster(envoyConfig.ConnectTimeoutMs, envoyConfig.CircuitBreaker, envoyConfig.OutlierDetection)
	assert.Equal(t, "foo_cluster", grpccl.Name)
	httpcl := kube.Service{Address: "bar", Port: uint32(8000), Type: kube.HTTP, Path: "/"}.Cluster(envoyConfig.ConnectTimeoutMs, envoyConfig.CircuitBreaker, envoyConfig.OutlierDetection)
	assert.Equal(t, "bar_cluster", httpcl.Name)
}

func TestGRPCIsHttp2Cluster(t *testing.T) {
	assert.Nil(t, kube.Service{Address: "foo", Port: uint32(8000), Type: kube.HTTP, Path: "/"}.Cluster(1000, config.CircuitBreakerConfig{}, config.OutlierDetectionConfig{}).Http2ProtocolOptions)
	assert.NotNil(t, kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC, Path: "/"}.Cluster(1000, config.CircuitBreakerConfig{}, config.OutlierDetectionConfig{}).Http2ProtocolOptions)
}

func TestDefaultTarget(t *testing.T) {
	target := kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC, Path: "/foo"}.DefaultTarget()
	assert.Equal(t, "foo", target.Host)
	assert.Equal(t, "/foo", target.Prefix)
	assert.Equal(t, "foo_cluster", target.ClusterName)
}

func TestDefaultServiceType(t *testing.T) {
	assert.Equal(t, kube.HTTP, kube.ServiceType(&corev1.Service{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"foo": "bar"}}}))
}

func TestAnnotatedServiceType(t *testing.T) {
	assert.Equal(t, kube.HTTP, kube.ServiceType(&corev1.Service{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"envoy-lb-operator.gojektech.k8s.io/service-type": "bar"}}}))
	assert.Equal(t, kube.GRPC, kube.ServiceType(&corev1.Service{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"envoy-lb-operator.gojektech.k8s.io/service-type": "grpc"}}}))
}

func TestDefaultServiceDomain(t *testing.T) {
	assert.Equal(t, "*", kube.ServiceDomain(&corev1.Service{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"foo": "bar"}}}))
}

func TestAnnotatedServiceDomain(t *testing.T) {
	assert.Equal(t, "bar.com", kube.ServiceDomain(&corev1.Service{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"envoy-lb-operator.gojektech.k8s.io/service-domain": "bar.com"}}}))
}

func TestDefaultServicePath(t *testing.T) {
	assert.Equal(t, "/", kube.ServicePath(&corev1.Service{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"foo": "bar"}}}))
}

func TestAnnotatedServicePath(t *testing.T) {
	assert.Equal(t, "/foo", kube.ServicePath(&corev1.Service{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"envoy-lb-operator.gojektech.k8s.io/service-path": "/foo"}}}))
}
