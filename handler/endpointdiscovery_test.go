package handler

import (
	"errors"
	"github.com/gojektech/kubehandler"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

type testTrigger struct {
	i int
}

func (t *testTrigger) increment() {
	t.i += 1
}

func (t *testTrigger) get() int {
	return t.i
}

func TestEndpointDiscovery_UpdateFuncSuccess(t *testing.T) {
	namespace := "namespace"
	name := "name"
	test := testTrigger{i: 0}
	coreClient := MockCoreClient{}
	endpointsInterface := MockEndpointsInterface{}
	endpointsInterface.On("Get", name, metav1.GetOptions{}).Return(&v1.Endpoints{}, nil)
	coreClient.On("Endpoints", namespace).Return(endpointsInterface)
	discovery := EndpointDiscovery{
		DefaultHandler: kubehandler.DefaultHandler{},
		CoreClient:     coreClient,
		Trigger:        test.increment,
	}

	err := discovery.UpdateFunc(namespace, name)
	assert.NoError(t, err)
	assert.Equal(t, 1, test.get())
}

func TestEndpointDiscovery_UpdateFuncFailure(t *testing.T) {
	namespace := "namespace"
	name := "name"
	test := testTrigger{i: 0}
	coreClient := MockCoreClient{}
	endpointsInterface := MockEndpointsInterface{}
	endpointsInterface.On("Get", name, metav1.GetOptions{}).Return(&v1.Endpoints{}, errors.New("Error"))
	coreClient.On("Endpoints", namespace).Return(endpointsInterface)
	discovery := EndpointDiscovery{
		DefaultHandler: kubehandler.DefaultHandler{},
		CoreClient:     coreClient,
		Trigger:        test.increment,
	}

	err := discovery.UpdateFunc(namespace, name)
	assert.Error(t, err)
	assert.Equal(t, 0, test.get())
}