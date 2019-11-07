package handler

import (
	"github.com/gojektech/kubehandler"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1types "k8s.io/client-go/kubernetes/typed/core/v1"
)

//SvcDiscovery handles addition, updation and deletion of headless services
type EndpointDiscovery struct {
	kubehandler.DefaultHandler
	CoreClient corev1types.CoreV1Interface
	Trigger   func()
}

func (d *EndpointDiscovery) UpdateFunc(namespace, name string) error {
	ep, err := d.CoreClient.Endpoints(namespace).Get(name, v1.GetOptions{})
	if err != nil {
		return err
	}
	log.Infof("Received update event for endpoint - %v. Refreshing...\n", ep)
	d.Trigger()
	return nil
}
