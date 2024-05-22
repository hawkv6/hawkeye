package main

import "github.com/hawkv6/hawkeye/hawkeye/cmd"

var (
	jagwServiceAddress   string
	jagwRequestPort      string
	jagwSubscriptionPort string
	grpcPort             string
)

func main() {
	cmd.Execute()
}
