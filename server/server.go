package server

import (
	"context"
	"fmt"
	"net"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//Xds has the grpc wrapper for xds api on a given port
type Xds struct {
	Server   xds.Server
	Port     uint
	cb       *callbacks
	cbSignal chan struct{}
}

const grpcMaxConcurrentStreams = 1000000

//Run starts the xDS server.
func (xd *Xds) Run(ctx context.Context) {
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams))
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", xd.Port))
	if err != nil {
		log.WithError(err).Fatal("failed to listen")
	}

	// register services
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, xd.Server)
	v2.RegisterEndpointDiscoveryServiceServer(grpcServer, xd.Server)
	v2.RegisterClusterDiscoveryServiceServer(grpcServer, xd.Server)
	v2.RegisterRouteDiscoveryServiceServer(grpcServer, xd.Server)
	v2.RegisterListenerDiscoveryServiceServer(grpcServer, xd.Server)

	log.WithFields(log.Fields{"port": xd.Port}).Info("management server listening")
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Error(err)
		}
	}()
	<-ctx.Done()

	grpcServer.GracefulStop()
}

func (xd *Xds) WaitForRequests() {
	<-xd.cbSignal
}

func (xd *Xds) Report() {
	<-xd.cbSignal
	xd.cb.Report()
}

func New(config cache.SnapshotCache, port uint) *Xds {
	signal := make(chan struct{})
	cb := &callbacks{
		signal:   signal,
		fetches:  0,
		requests: 0,
	}
	return &Xds{Server: xds.NewServer(config, cb), Port: port, cbSignal: signal, cb: cb}
}
