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

func filterEvents(endpointLabel string) func(*metav1.ListOptions) {
	return func(opt *metav1.ListOptions) {
		opt.LabelSelector = endpointLabel
	}
}

func StartSvcKubeHandler(client *kubernetes.Clientset, triggerfunc func(eventType envoy.LBEventType, svc *v1.Service), upstreamLabel, namespace string) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	kubeInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(client, time.Second*1, kubeinformers.WithNamespace(namespace), kubeinformers.WithTweakListOptions(filterEvents(upstreamLabel)))
	informer := kubeInformerFactory.Core().V1().Services().Informer()
	discoveryHandler := &handler.SvcDiscovery{
		CoreClient: client.CoreV1(),
		SVCTrigger: triggerfunc,
		DefaultHandler: kubehandler.DefaultHandler{
			Informer: informer,
			Synced:   informer.HasSynced,
		},
	}
	loop := kubehandler.NewEventLoop("service_queue")
	loop.Register(discoveryHandler)
	go loop.Run(20, ctx.Done())
	go kubeInformerFactory.Start(ctx.Done())

	return cancel
}

func StartEndpointKubeHandler(client *kubernetes.Clientset, triggerfunc func(), endpointLabel, namespace string) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	kubeInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(client, time.Second*1, kubeinformers.WithNamespace(namespace), kubeinformers.WithTweakListOptions(filterEvents(endpointLabel)))
	informer := kubeInformerFactory.Core().V1().Endpoints().Informer()
	endpointHandler := &handler.EndpointDiscovery{
		CoreClient: client.CoreV1(),
		Trigger:    triggerfunc,
		DefaultHandler: kubehandler.DefaultHandler{
			Informer: informer,
			Synced:   informer.HasSynced,
		},
	}
	loop := kubehandler.NewEventLoop("endpoint_queue")
	loop.Register(endpointHandler)
	go loop.Run(20, ctx.Done())
	go kubeInformerFactory.Start(ctx.Done())

	return cancel
}
