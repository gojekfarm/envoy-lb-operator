package handler

import (
	"github.com/stretchr/testify/mock"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

type MockEndpointsInterface struct {
	mock.Mock
}

func (m MockEndpointsInterface) Create(*v1.Endpoints) (*v1.Endpoints, error) {
	panic("implement me")
}

func (m MockEndpointsInterface) Update(*v1.Endpoints) (*v1.Endpoints, error) {
	panic("implement me")
}

func (m MockEndpointsInterface) Delete(name string, options *metav1.DeleteOptions) error {
	panic("implement me")
}

func (m MockEndpointsInterface) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	panic("implement me")
}

func (m MockEndpointsInterface) Get(name string, options metav1.GetOptions) (*v1.Endpoints, error) {
	args := m.Called(name, options)
	return args.Get(0).(*v1.Endpoints), args.Error(1)
}

func (m MockEndpointsInterface) List(opts metav1.ListOptions) (*v1.EndpointsList, error) {
	panic("implement me")
}

func (m MockEndpointsInterface) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	panic("implement me")
}

func (m MockEndpointsInterface) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Endpoints, err error) {
	panic("implement me")
}

