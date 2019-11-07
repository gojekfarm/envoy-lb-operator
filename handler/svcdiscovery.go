package handler

import (
	"github.com/gojekfarm/envoy-lb-operator/envoy"
	"github.com/gojektech/kubehandler"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	corev1types "k8s.io/client-go/kubernetes/typed/core/v1"
)

//SvcDiscovery handles addition, updation and deletion of headless services
type SvcDiscovery struct {
	kubehandler.DefaultHandler
	CoreClient corev1types.CoreV1Interface
	SVCTrigger func(eventType envoy.LBEventType, svc *corev1.Service)
}

func (d *SvcDiscovery) AddFunc(namespace, name string) error {
	svc, err := d.CoreClient.Services(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return err
	}
	d.SVCTrigger(envoy.ADDED, svc)
	return nil
}

func (d *SvcDiscovery) UpdateFunc(namespace, name string) error {
	svc, err := d.CoreClient.Services(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return err
	}
	d.SVCTrigger(envoy.UPDATED, svc)
	return nil
}

func (d *SvcDiscovery) DeleteFunc(namespace, name string) error {
	//After deletion from kubernetes we won't be getting the service object.
	//We pass the dummy service with the name so that the deletion action can be triggered.
	svc := &corev1.Service{
		ObjectMeta: v1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: corev1.ServiceSpec{ClusterIP: corev1.ClusterIPNone,
			Ports: []corev1.ServicePort{{TargetPort: intstr.IntOrString{IntVal: 1234},}}},
	}
	d.SVCTrigger(envoy.DELETED, svc)
	return nil
}
