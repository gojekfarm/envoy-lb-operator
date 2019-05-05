package handler

import (
	"github.com/gojekfarm/envoy-lb-operator/envoy"
	"github.com/gojektech/kubehandler"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1types "k8s.io/client-go/kubernetes/typed/core/v1"
)

//Discovery handles addition, updation and deletion of headless services
type Discovery struct {
	kubehandler.DefaultHandler
	CoreClient corev1types.CoreV1Interface
	SVCTrigger func(eventType envoy.LBEventType, svc *corev1.Service)
}

func (d *Discovery) AddFunc(namespace, name string) error {
	svc, err := d.CoreClient.Services(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return err
	}
	d.SVCTrigger(envoy.ADDED, svc)
	return nil
}

func (d *Discovery) UpdateFunc(namespace, name string) error {
	svc, err := d.CoreClient.Services(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return err
	}
	d.SVCTrigger(envoy.UPDATED, svc)
	return nil
}

func (d *Discovery) DeleteFunc(namespace, name string) error {
	svc, err := d.CoreClient.Services(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return err
	}
	d.SVCTrigger(envoy.DELETED, svc)
	return nil
}
