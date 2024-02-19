package commands

import (
	"os"

	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/jagw"
	"github.com/spf13/cobra"
)

var (
	jagwServiceAddress   string
	jagwRequestPort      string
	jagwSubscriptionPort string
	grpcPort             string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Hawkeye controller",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.NewDefaultConfig(jagwServiceAddress, jagwRequestPort, jagwSubscriptionPort, grpcPort)
		if err != nil {
			log.Fatalf("Error creating config: %v", err)
		}
		log.Printf("Config created successfully %v", config)

		requestService := jagw.NewDefaultJagwRequestService(config)
		if err := requestService.Init(); err != nil {
			log.Fatalf("Error initializing JAGW Request Service: %v", err)
		}
		graph := graph.NewDefaultGraph()
		if err := requestService.GetLsLinks(graph); err != nil {
			log.Fatalf("Error getting LsLinks from JAGW: %v", err)
		}
		source, err := graph.GetNode("0000.0000.000a")
		if err != nil {
			log.Fatalf("Error getting source node: %v", err)
		}
		destination, err := graph.GetNode("0000.0000.000c")
		if err != nil {
			log.Fatalf("Error getting destination node: %v", err)
		}
		path, err := graph.GetShortestPath(source, destination, "delay")
		if err != nil {
			log.Fatalf("Error getting shortest path: %v", err)
		}
		log.Infoln("Shortest path from ", source.GetId(), " to ", destination.GetId(), " is: ")
		for _, edge := range path {
			log.Infoln("Edge: ", edge.From().GetId(), " -> ", edge.To().GetId())
		}

		// TODO subscribe for lslinkedge events https://github.com/hawkv6/hawkeye/issues/2
		// TODO get all linksedge from jagw https://github.com/hawkv6/hawkeye/issues/1
		// TODO start grpc server (handle streams): https://github.com/hawkv6/hawkeye/issues/3
		// TODO calculate based on intents: https://github.com/hawkv6/hawkeye/issues/4
		// TODO recaluclate based on events: https://github.com/hawkv6/hawkeye/issues/5
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&jagwServiceAddress, "jagw-service-address", "j", os.Getenv("HAWKEYE_JAGW_SERVICE_ADDRESS"), "JAGW Service Address e.g. localhost or 127.0.0.1")
	startCmd.Flags().StringVarP(&jagwRequestPort, "jagw-request-port", "r", os.Getenv("HAWKEYE_JAGW_REQUEST_PORT"), "JAGW Request Port e.g. 9903")
	startCmd.Flags().StringVarP(&jagwSubscriptionPort, "jagw-subscription-port", "s", os.Getenv("HAWKEYE_JAGW_SUBSCRIPTION_PORT"), "JAGW Subscription Port e.g. 9902")
	startCmd.Flags().StringVarP(&grpcPort, "grpc-port", "p", os.Getenv("HAWKEYE_GRPC_PORT"), "gRPC Port e.g. 10000")
	markRequiredFlags(startCmd, []string{"jagw-service-address", "jagw-request-port", "jagw-subscription-port", "grpc-port"})
}
