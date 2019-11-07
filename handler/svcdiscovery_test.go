package handler

import (
	"errors"
	"github.com/gojekfarm/envoy-lb-operator/envoy"
	"github.com/gojektech/kubehandler"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

type testSvcTrigger struct {
	i         int
	eventType envoy.LBEventType
}

func (t *testSvcTrigger) trigger(eventType envoy.LBEventType, svc *v1.Service) {
	t.i += 1
	t.eventType = eventType
}

func TestSvcDiscovery_AddFuncSuccess(t *testing.T) {
	namespace := "namespace"
	name := "name"
	test := &testSvcTrigger{i: 0}
	coreClient := MockCoreClient{}
	svcInterface := MockSvcInterface{}
	svcInterface.On("Get", name, metav1.GetOptions{}).Return(&v1.Service{}, nil)
	coreClient.On("Services", namespace).Return(svcInterface)
	discovery := SvcDiscovery{
		DefaultHandler: kubehandler.DefaultHandler{},
		CoreClient:     coreClient,
		SVCTrigger:     test.trigger,
	}

	err := discovery.AddFunc(namespace, name)
	assert.NoError(t, err)
	assert.Equal(t, 1, test.i)
	assert.Equal(t, envoy.ADDED, test.eventType)
}

func TestSvcDiscovery_AddFuncFailure(t *testing.T) {
	namespace := "namespace"
	name := "name"
	test := &testSvcTrigger{i: 0}
	coreClient := MockCoreClient{}
	svcInterface := MockSvcInterface{}
	svcInterface.On("Get", name, metav1.GetOptions{}).Return(&v1.Service{}, errors.New("Error"))
	coreClient.On("Services", namespace).Return(svcInterface)
	discovery := SvcDiscovery{
		DefaultHandler: kubehandler.DefaultHandler{},
		CoreClient:     coreClient,
		SVCTrigger:     test.trigger,
	}

	err := discovery.AddFunc(namespace, name)
	assert.Error(t, err)
	assert.Equal(t, 0, test.i)
}

func TestSvcDiscovery_UpdateFuncSuccess(t *testing.T) {
	namespace := "namespace"
	name := "name"
	test := &testSvcTrigger{i: 0}
	coreClient := MockCoreClient{}
	svcInterface := MockSvcInterface{}
	svcInterface.On("Get", name, metav1.GetOptions{}).Return(&v1.Service{}, nil)
	coreClient.On("Services", namespace).Return(svcInterface)
	discovery := SvcDiscovery{
		DefaultHandler: kubehandler.DefaultHandler{},
		CoreClient:     coreClient,
		SVCTrigger:     test.trigger,
	}

	err := discovery.UpdateFunc(namespace, name)
	assert.NoError(t, err)
	assert.Equal(t, 1, test.i)
	assert.Equal(t, envoy.UPDATED, test.eventType)
}

func TestSvcDiscovery_UpdateFuncFailure(t *testing.T) {
	namespace := "namespace"
	name := "name"
	test := &testSvcTrigger{i: 0}
	coreClient := MockCoreClient{}
	svcInterface := MockSvcInterface{}
	svcInterface.On("Get", name, metav1.GetOptions{}).Return(&v1.Service{}, errors.New("Error"))
	coreClient.On("Services", namespace).Return(svcInterface)
	discovery := SvcDiscovery{
		DefaultHandler: kubehandler.DefaultHandler{},
		CoreClient:     coreClient,
		SVCTrigger:     test.trigger,
	}

	err := discovery.UpdateFunc(namespace, name)
	assert.Error(t, err)
	assert.Equal(t, 0, test.i)
}

func TestSvcDiscovery_DeleteFuncSuccess(t *testing.T) {
	namespace := "namespace"
	name := "name"
	test := &testSvcTrigger{i: 0}
	coreClient := MockCoreClient{}
	discovery := SvcDiscovery{
		DefaultHandler: kubehandler.DefaultHandler{},
		CoreClient:     coreClient,
		SVCTrigger:     test.trigger,
	}

	err := discovery.DeleteFunc(namespace, name)
	assert.NoError(t, err)
	assert.Equal(t, 1, test.i)
	assert.Equal(t, envoy.DELETED, test.eventType)
}

