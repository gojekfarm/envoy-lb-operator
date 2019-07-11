package cmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gojekfarm/envoy-lb-operator/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gojekfarm/envoy-lb-operator/envoy"
	"github.com/gojekfarm/envoy-lb-operator/kube"
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
}

var serveCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the Envoy XDS management server",
	Run:   serve,
}

var debug bool
var masterurl string
var kubeConfig string

func init() {
	cliCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "use debug level logs")
	serveCmd.Flags().StringVarP(&masterurl, "master", "m", "", "Master URL for Kube API server")
	serveCmd.Flags().StringVarP(&kubeConfig, "kubeconfig", "c", "", "Help message for toggle")

	cobra.OnInitialize(initConfig)
	cliCmd.AddCommand(serveCmd)
}

func initConfig() {
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

func cancelOnInterrupt(cancelFn func()) {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	cancelFn()
}

func serve(cmd *cobra.Command, args []string) {
	log.Printf("Starting control plane")
	cfg, err := clientcmd.BuildConfigFromFlags(masterurl, kubeConfig)
	if err != nil {
		log.Fatal(err)
	}

	envoyConfig, err := config.LoadDefaultEnvoyConfig()
	if err != nil {
		log.Fatal(err)
	}
	kubeClient, err := kubernetes.NewForConfig(cfg)
	lb := envoy.NewLB("nodeID", envoyConfig)
	if err != nil {
		//test data for now.
		lb.Trigger(envoy.LBEvent{
			Svc:       kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC},
			EventType: envoy.ADDED,
		})
	} else {
		go cancelOnInterrupt(server.StartKubehandler(kubeClient, lb.SvcTrigger))
	}

	ctx := context.Background()
	// start the xDS server
	xdsServer := server.New(lb.Config, 18000)
	go xdsServer.Run(ctx)
	xdsServer.WaitForRequests()
	go xdsServer.Report()
	go lb.HandleEvents()

	for {
		lb.Snapshot()
		time.Sleep(10 * time.Second)
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
	}
}

//Execute called from main
func Execute() {
	if err := cliCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
