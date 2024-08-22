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
	consulServerAddress  string
)

var rootCmd = &cobra.Command{
	Use:   "hawkeye",
	Short: "Controller for Enabling Intent-Based Networking in SRv6",
	Long: `
Hawkeye is a controller that enables Intent-Based Networking in SRv6.

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
