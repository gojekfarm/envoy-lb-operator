package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gojekfarm/envoy-lb-operator/envoy"
	"github.com/gojekfarm/envoy-lb-operator/kube"
	"github.com/gojekfarm/envoy-lb-operator/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

func init() {
	cliCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "use debug level logs")
	cobra.OnInitialize(initConfig)
	cliCmd.AddCommand(serveCmd)
}

func initConfig() {
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

func serve(cmd *cobra.Command, args []string) {
	log.Printf("Starting control plane")

	lb := envoy.NewLB("nodeID")

	ctx := context.Background()

	// start the xDS server
	xdsServer := server.New(lb.Config, 18000)
	go xdsServer.Run(ctx)
	xdsServer.WaitForRequests()
	go xdsServer.Report()
	go lb.HandleEvents()

	//test data for now.
	lb.Trigger(envoy.LBEvent{
		Svc:       kube.Service{Address: "foo", Port: uint32(8000), Type: kube.GRPC},
		EventType: envoy.ADDED,
	})
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
