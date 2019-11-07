package handler

import (
	"github.com/stretchr/testify/mock"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type MockSvcInterface struct {
	mock.Mock
}

func (m MockSvcInterface) Create(*v1.Service) (*v1.Service, error) {
	panic("implement me")
}

func (m MockSvcInterface) Update(*v1.Service) (*v1.Service, error) {
	panic("implement me")
}

func (m MockSvcInterface) UpdateStatus(*v1.Service) (*v1.Service, error) {
	panic("implement me")
}

func (m MockSvcInterface) Delete(name string, options *metav1.DeleteOptions) error {
	panic("implement me")
}

func (m MockSvcInterface) Get(name string, options metav1.GetOptions) (*v1.Service, error) {
	args := m.Called(name, options)
	return args.Get(0).(*v1.Service), args.Error(1)
}

func (m MockSvcInterface) List(opts metav1.ListOptions) (*v1.ServiceList, error) {
	panic("implement me")
}

func (m MockSvcInterface) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	panic("implement me")
}

func (m MockSvcInterface) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Service, err error) {
	panic("implement me")
}

func (m MockSvcInterface) ProxyGet(scheme, name, port, path string, params map[string]string) rest.ResponseWrapper {
	panic("implement me")
}


