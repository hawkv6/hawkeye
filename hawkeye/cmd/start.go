package cmd

import (
	"os"
	"os/signal"

	"github.com/hawkv6/hawkeye/pkg/adapter"
	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/calculation"
	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/hawkv6/hawkeye/pkg/controller"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/jagw"
	"github.com/hawkv6/hawkeye/pkg/messaging"
	"github.com/hawkv6/hawkeye/pkg/normalizer"
	"github.com/hawkv6/hawkeye/pkg/processor"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Hawkeye controller",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.NewFullConfig(jagwServiceAddress, jagwRequestPort, jagwSubscriptionPort, grpcPort)
		if err != nil {
			log.Fatalf("Error creating config: %v", err)
		}
		log.Infoln("Config created successfully")

		eventChan := make(chan domain.NetworkEvent)
		adapter := adapter.NewDomainAdapter()
		defaultHelper := helper.NewDefaultHelper()
		graph := graph.NewNetworkGraph(defaultHelper)
		cache := cache.NewInMemoryCache()
		updateChan := make(chan struct{})
		latencyQueue := normalizer.NewNormalizationQueue(helper.RollingWindowSize)
		jitterQueue := normalizer.NewNormalizationQueue(helper.RollingWindowSize)
		packetLossQueue := normalizer.NewNormalizationQueue(helper.RollingWindowSize)
		iqrMinMaxNormalizer := normalizer.NewIQRMinMaxNormalizer(latencyQueue, jitterQueue, packetLossQueue)
		processor := processor.NewNetworkProcessor(graph, cache, iqrMinMaxNormalizer, eventChan, defaultHelper, updateChan)

		requestService := jagw.NewJagwRequestService(config, adapter, processor, defaultHelper)
		if err := requestService.Init(); err != nil {
			log.Fatalf("Error initializing JAGW Request Service: %v", err)
		}
		if err := requestService.Start(); err != nil {
			log.Fatalf("Error starting JAGW Request Service: %v", err)
		}

		messagingChannels := messaging.NewPathMessagingChannels()
		manager := calculation.NewCalculationManager(cache, graph, defaultHelper)
		controller := controller.NewSessionController(manager, messagingChannels, updateChan)
		go controller.Start()

		go processor.Start()

		subscriptionService := jagw.NewJagwSubscriptionService(config, adapter, defaultHelper, eventChan)
		if err := subscriptionService.Init(); err != nil {
			log.Fatalf("Error initializing JAGW Subscription Service: %v", err)
		}
		if err := subscriptionService.Start(); err != nil {
			log.Fatalf("Error starting JAGW Subscription Service: %v", err)
		}

		server := messaging.NewGrpcMessagingServer(adapter, config, messagingChannels)

		if err := server.Start(); err != nil {
			log.Fatalf("Error starting gRPC server: %v", err)
		}

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)

		<-signalChan
		log.Info("Received interrupt signal, shutting down")
		requestService.Stop()
		subscriptionService.Stop()
		processor.Stop()
		// TODO Stop the controller
		// controller.Close()
		// TODO stop the gRPC server
		// server.Stop()

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
