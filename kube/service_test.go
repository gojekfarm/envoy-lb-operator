package kube_test

import (
	"testing"

	"github.com/gojekfarm/envoy-lb-operator/kube"
	"github.com/stretchr/testify/assert"
)

func TestServiceClusterName(t *testing.T) {
	grpccl := kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC}.Cluster()
	assert.Equal(t, "foo_cluster", grpccl.Name)
	httpcl := kube.Service{Address: "bar", Port: uint32(8000), Type: kube.HTTP}.Cluster()
	assert.Equal(t, "bar_cluster", httpcl.Name)
}

func TestGRPCIsHttp2Cluster(t *testing.T) {
	assert.Nil(t, kube.Service{Address: "foo", Port: uint32(8000), Type: kube.HTTP}.Cluster().Http2ProtocolOptions)
	assert.NotNil(t, kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC}.Cluster().Http2ProtocolOptions)
}

func TestDefaultTarget(t *testing.T) {
	target := kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC}.DefaultTarget()
	assert.Equal(t, "foo", target.Host)
	assert.Equal(t, "/", target.Prefix)
	assert.Equal(t, "foo_cluster", target.ClusterName)
}

func TestTarget(t *testing.T) {
	target := kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC}.Target("/foo")
	assert.Equal(t, "foo", target.Host)
	assert.Equal(t, "/foo", target.Prefix)
	assert.Equal(t, "foo_cluster", target.ClusterName)
}
