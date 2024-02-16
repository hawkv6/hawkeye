package commands

import (
	"os"

	"github.com/hawkv6/hawkeye/pkg/config"
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

		// subscribe for lslinkedge events
		// get all linksedge from jagw
		// start grpc server (handle streams)
		// calculate based on intents
		// recaluclate based on events
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
