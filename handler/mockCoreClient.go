package handler

import (
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type MockCoreClient struct {
	mock.Mock
}

func (m MockCoreClient) RESTClient() rest.Interface {
	panic("implement me")
}

func (m MockCoreClient) ComponentStatuses() v1.ComponentStatusInterface {
	panic("implement me")
}

func (m MockCoreClient) ConfigMaps(namespace string) v1.ConfigMapInterface {
	panic("implement me")
}

func (m MockCoreClient) Endpoints(namespace string) v1.EndpointsInterface {
	args := m.Called(namespace)
	return args.Get(0).(v1.EndpointsInterface)
}

func (m MockCoreClient) Events(namespace string) v1.EventInterface {
	panic("implement me")
}

func (m MockCoreClient) LimitRanges(namespace string) v1.LimitRangeInterface {
	panic("implement me")
}

func (m MockCoreClient) Namespaces() v1.NamespaceInterface {
	panic("implement me")
}

func (m MockCoreClient) Nodes() v1.NodeInterface {
	panic("implement me")
}

func (m MockCoreClient) PersistentVolumes() v1.PersistentVolumeInterface {
	panic("implement me")
}

func (m MockCoreClient) PersistentVolumeClaims(namespace string) v1.PersistentVolumeClaimInterface {
	panic("implement me")
}

func (m MockCoreClient) Pods(namespace string) v1.PodInterface {
	panic("implement me")
}

func (m MockCoreClient) PodTemplates(namespace string) v1.PodTemplateInterface {
	panic("implement me")
}

func (m MockCoreClient) ReplicationControllers(namespace string) v1.ReplicationControllerInterface {
	panic("implement me")
}

func (m MockCoreClient) ResourceQuotas(namespace string) v1.ResourceQuotaInterface {
	panic("implement me")
}

func (m MockCoreClient) Secrets(namespace string) v1.SecretInterface {
	panic("implement me")
}

func (m MockCoreClient) Services(namespace string) v1.ServiceInterface {
	args := m.Called(namespace)
	return args.Get(0).(v1.ServiceInterface)
}

func (m MockCoreClient) ServiceAccounts(namespace string) v1.ServiceAccountInterface {
	panic("implement me")
}
