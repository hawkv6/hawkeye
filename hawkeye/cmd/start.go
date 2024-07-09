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
	"github.com/hawkv6/hawkeye/pkg/jagw"
	"github.com/hawkv6/hawkeye/pkg/messaging"
	"github.com/hawkv6/hawkeye/pkg/processor"
	"github.com/hawkv6/hawkeye/pkg/service"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Hawkeye controller",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: check move other elements to config -> e.g. consul address
		config, err := config.NewFullConfig(jagwServiceAddress, jagwRequestPort, jagwSubscriptionPort, grpcPort)
		if err != nil {
			log.Fatalf("Error creating config: %v", err)
		}
		log.Infoln("Config created successfully")

		graph := graph.NewNetworkGraph()
		cache := cache.NewInMemoryCache()

		eventChan := make(chan domain.NetworkEvent)
		updateChan := make(chan struct{})

		nodeEventProcessor := processor.NewNodeEventProcessor(graph, cache)
		linkEventProcessor := processor.NewLinkEventProcessor(graph, cache)
		prefixEventProcessor := processor.NewPrefixEventProcessor(graph, cache)
		sidEventProcessor := processor.NewSidEventProcessor(graph, cache)
		eventOptions := processor.EventOptions{
			NodeEventProcessor:   nodeEventProcessor,
			LinkEventProcessor:   linkEventProcessor,
			PrefixEventProcessor: prefixEventProcessor,
			SidEventProcessor:    sidEventProcessor,
			EventDispatcher:      processor.NewEventDispatcher(nodeEventProcessor, linkEventProcessor, prefixEventProcessor, sidEventProcessor),
		}
		processor := processor.NewNetworkProcessor(graph, cache, eventChan, updateChan, eventOptions)

		adapter := adapter.NewDomainAdapter()
		requestService := jagw.NewJagwRequestService(config, adapter, processor)
		if err := requestService.Init(); err != nil {
			log.Fatalf("Error initializing JAGW Request Service: %v", err)
		}
		if err := requestService.Start(); err != nil {
			log.Fatalf("Error starting JAGW Request Service: %v", err)
		}

		serviceMonitor, err := service.NewConsulServiceMonitor(cache, updateChan)
		if err != nil {
			log.Fatalf("Error creating Consult service monitor: %v", err)
		}
		go serviceMonitor.StartMonitoring()

		messagingChannels := messaging.NewPathMessagingChannels()
		manager := calculation.NewCalculationManager(cache, graph)
		controller := controller.NewSessionController(manager, messagingChannels, updateChan)
		go controller.Start()
		go processor.Start()

		subscriptionService := jagw.NewJagwSubscriptionService(config, adapter, eventChan)
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
		serviceMonitor.StopMonitoring()
		// TODO Stop the controller
		// controller.Close()
		// TODO stop the gRPC server
		// server.Stop()
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
