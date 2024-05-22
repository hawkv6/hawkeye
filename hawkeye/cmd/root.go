package cmd

import (
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/spf13/cobra"
)

var (
	log                  = logging.DefaultLogger.WithField("subsystem", "cmd")
	jagwServiceAddress   string
	jagwRequestPort      string
	jagwSubscriptionPort string
	grpcPort             string
)

func markRequiredFlags(cmd *cobra.Command, flags []string) {
	for _, flag := range flags {
		if err := cmd.MarkFlagRequired(flag); err != nil {
			log.Fatal(err)
		}
	}

}

var rootCmd = &cobra.Command{
	Use:   "hawkeye",
	Short: "Controller for Enabling Intent-Driven End-to-End SRv6 Networking",
	Long: `
Hawkeye is a controller that enables intent-driven end-to-end SRv6 networking.

Start Hawkeye by running the following command:
$ hawkeye start --jagw-request-service localhost:9903 --jagw-subscription-service localhost:9902 --grpc-port 10000
or
$ hawkeye start -r localhost:9903 -s localhost:9902 -p 10000`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
