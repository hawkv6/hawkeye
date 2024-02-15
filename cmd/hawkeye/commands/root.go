package commands

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hawkeye",
	Short: "Controller for Enabling Intent-Driven End-to-End SRv6 Networking",
	Long: `
Hawkeye is a controller that enables intent-driven end-to-end SRv6 networking.

Start Hawkeye by running the following command:
$ hawkeye start --jagw-request-service localhost:9093 --jagw-subscription-service localhost:9092 --grpc-port 10000
or
$ hawkeye start -r localhost:9093 -s localhost:9092 -p 10000`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
