package server

import (
	"context"
	"time"

	"github.com/gojekfarm/envoy-lb-operator/envoy"
	"github.com/gojekfarm/envoy-lb-operator/handler"
	"github.com/gojektech/kubehandler"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

func filterServices(endpointLabel string) func(*metav1.ListOptions) {
	return func(opt *metav1.ListOptions) {
		opt.LabelSelector = endpointLabel
	}
}

func StartKubehandler(client *kubernetes.Clientset, triggerfunc func(eventType envoy.LBEventType, svc *v1.Service), endpointLabel, namespace string) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	kubeInformerFactory := kubeinformers.NewFilteredSharedInformerFactory(client, time.Second*1, namespace, filterServices(endpointLabel))
	informer := kubeInformerFactory.Core().V1().Services().Informer()
	discoveryHandler := &handler.Discovery{
		CoreClient: client.CoreV1(),
		SVCTrigger: triggerfunc,
		DefaultHandler: kubehandler.DefaultHandler{
			Informer: informer,
			Synced:   informer.HasSynced,
		},
	}
	loop := kubehandler.NewEventLoop("discovery_queue")
	loop.Register(discoveryHandler)
	go kubeInformerFactory.Start(ctx.Done())
	go loop.Run(20, ctx.Done())

	// Initialise for the beginning
	serviceList, _ := client.CoreV1().Services(v1.NamespaceAll).List(metav1.ListOptions{LabelSelector: endpointLabel})
	for _, svc := range serviceList.Items {
		triggerfunc(envoy.ADDED, &svc)
	}

	return cancel
}
