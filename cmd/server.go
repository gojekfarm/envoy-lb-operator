package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/gojekfarm/envoy-lb-operator/config"
	"github.com/gojekfarm/envoy-lb-operator/envoy"
	"github.com/gojekfarm/envoy-lb-operator/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var cliCmd = &cobra.Command{
	Use:   "envoy-lb-operator",
	Short: "envoy-lb-operator is an xds control plane for envoy",
	Long:  `This adds relevant k8s resources to XDS, LDS, ADS and CDS. Envoy config can read about on, https://www.envoyproxy.io/docs/envoy/latest/configuration/cluster_manager/cds`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
	Version: "0.0.1",
}

var serveCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the Envoy XDS management server",
	Run:   serve,
}

var masterurl string
var kubeConfig string

func init() {
	serveCmd.Flags().StringVarP(&masterurl, "master", "m", "", "Master URL for Kube API server")
	serveCmd.Flags().StringVarP(&kubeConfig, "kubeconfig", "c", "", "Help message for toggle")

	cliCmd.AddCommand(serveCmd)
}


func cancelOnInterrupt(cancelFn func()) {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	cancelFn()
}

func serve(cmd *cobra.Command, args []string) {
	log.Printf("Running application: %s\n", cmd.Version)
	cfg, err := clientcmd.BuildConfigFromFlags(masterurl, kubeConfig)
	log.Printf("got config: host: %s pat: %+v\n", cfg.Host, cfg.APIPath)
	if err != nil {
		log.Fatal(err)
	}

	envoyConfig := config.GetEnvoyConfig()
	snapshotCache := cache.NewSnapshotCache(true, envoy.Hasher{}, envoy.Logger{})
	startXdsServer(snapshotCache)
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Error creating kube client: %v", err)
	}

	for _, mapping := range config.GetDiscoveryMapping() {
		lb := envoy.NewLB(mapping.EnvoyId, envoyConfig, snapshotCache)
		go lb.HandleEvents()
		cancelFn := server.StartKubehandler(kubeClient, lb.SvcTrigger, mapping.UpstreamEndpointLabel, mapping.Namespace)
		go cancelOnInterrupt(cancelFn)

		go func() {
			for {
				lb.SnapshotRunner()
				time.Sleep(time.Duration(config.RefreshIntervalInS()) * time.Second)
			}
		}()
	}

	log.Info("Waiting in main")
	for {
		time.Sleep(1000)
	}
}

func startXdsServer(snapshotCache cache.SnapshotCache) {
	ctx := context.Background()
	// start the xDS server
	xdsServer := server.New(snapshotCache, 18000)
	go xdsServer.Run(ctx)
	xdsServer.WaitForRequests()
	go xdsServer.Report()
}

//Execute called from main
func Execute() {
	if err := cliCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
