package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	jagwRequestService      string
	jagwSubscriptionService string
	grpcPort                int
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Hawkeye controller",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting Hawkeye with JAGW Request Service: %s, JAGW Subscription Service: %s, and gRPC Port: %d\n", jagwRequestService, jagwSubscriptionService, grpcPort)
		// subscribe for lslinkedge events
		// get all linksedge from jagw
		// start grpc server (handle streams)
		// calculate based on intents
		// recaluclate based on events
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&jagwRequestService, "jagw-request-service", "r", "", "JAGW Request Service e.g. localhost:9093")
	startCmd.Flags().StringVarP(&jagwSubscriptionService, "jagw-subscription-service", "s", "", "JAGW Subscription Service e.g. localhost:9092")
	startCmd.Flags().IntVarP(&grpcPort, "grpc-port", "p", 0, "gRPC Port e.g. 10000")
}
