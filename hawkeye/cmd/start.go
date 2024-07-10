package cmd

import (
	"os"
	"os/signal"
	"sync"

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

func initializeNetworkProcessor(graph graph.Graph, cache cache.Cache, eventChan chan domain.NetworkEvent, updateChan chan struct{}) *processor.NetworkProcessor {
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
	return processor.NewNetworkProcessor(graph, cache, eventChan, updateChan, eventOptions)
}

func requestNetworkElements(config *config.FullConfig, adapter adapter.Adapter, networkProcessor *processor.NetworkProcessor) {
	requestService := jagw.NewJagwRequestService(config, adapter, networkProcessor)
	if err := requestService.Init(); err != nil {
		log.Fatalf("Error initializing JAGW Request Service: %v", err)
	}
	if err := requestService.Start(); err != nil {
		log.Fatalf("Error starting JAGW Request Service: %v", err)
	}
	requestService.Stop()
}

func startServiceMonitoring(cache cache.Cache, updateChan chan struct{}, wg *sync.WaitGroup) *service.ConsulServiceMonitor {
	serviceMonitor, err := service.NewConsulServiceMonitor(cache, updateChan, consulServerAddress)
	if err != nil {
		log.Fatalf("Error creating Consult service monitor: %v", err)
	}
	go func() {
		serviceMonitor.Start()
		wg.Done()
	}()
	return serviceMonitor
}

func startController(cache cache.Cache, graph graph.Graph, updateChan chan struct{}, wg *sync.WaitGroup) (*messaging.PathMessagingChannels, *controller.SessionController) {
	messagingChannels := messaging.NewPathMessagingChannels()
	calculationSetupProvider := calculation.NewCalculationSetupProvider(cache, graph)
	manager := calculation.NewCalculationManager(cache, graph, calculationSetupProvider)
	controller := controller.NewSessionController(manager, messagingChannels, updateChan)
	wg.Add(1)
	go func() {
		controller.Start()
		wg.Done()
	}()
	return messagingChannels, controller
}

func startNetworkProcessor(networkProcesor *processor.NetworkProcessor, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		networkProcesor.Start()
		wg.Done()
	}()
}

func startSubscriptionService(config *config.FullConfig, adapter adapter.Adapter, eventChan chan domain.NetworkEvent) *jagw.JagwSubscriptionService {
	subscriptionService := jagw.NewJagwSubscriptionService(config, adapter, eventChan)
	if err := subscriptionService.Init(); err != nil {
		log.Fatalf("Error initializing JAGW Subscription Service: %v", err)
	}
	if err := subscriptionService.Start(); err != nil {
		log.Fatalf("Error starting JAGW Subscription Service: %v", err)
	}
	return subscriptionService
}

func startGrpcServer(adapter adapter.Adapter, config *config.FullConfig, messagingChannels messaging.MessagingChannels, wg *sync.WaitGroup) *messaging.GrpcMessagingServer {
	server := messaging.NewGrpcMessagingServer(adapter, config, messagingChannels)
	wg.Add(1)
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Error starting gRPC server: %v", err)
		}
	}()
	return server
}

func listenForInterruptSignal(server *messaging.GrpcMessagingServer, subscriptionService *jagw.JagwSubscriptionService, serviceMonitor *service.ConsulServiceMonitor, networkProcessor *processor.NetworkProcessor, controller *controller.SessionController, wg *sync.WaitGroup) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	log.Info("Received interrupt signal, shutting down")
	server.Stop()
	subscriptionService.Stop()
	serviceMonitor.Stop()
	networkProcessor.Stop()
	controller.Stop()
	wg.Wait()
	log.Infoln("All services stopped successfully")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Hawkeye controller",
	Run: func(cmd *cobra.Command, args []string) {
		graph := graph.NewNetworkGraph()
		cache := cache.NewInMemoryCache()
		eventChan := make(chan domain.NetworkEvent)
		updateChan := make(chan struct{})
		networkProcessor := initializeNetworkProcessor(graph, cache, eventChan, updateChan)

		config, err := config.NewFullConfig(jagwServiceAddress, jagwRequestPort, jagwSubscriptionPort, grpcPort)
		if err != nil {
			log.Fatalf("Error creating config: %v", err)
		}
		log.Infoln("Config created successfully")
		requestNetworkElements(config, adapter.NewDomainAdapter(), networkProcessor)

		adapter := adapter.NewDomainAdapter()
		wg := sync.WaitGroup{}
		serviceMonitor := startServiceMonitoring(cache, updateChan, &wg)
		messagingChannels, controller := startController(cache, graph, updateChan, &wg)
		startNetworkProcessor(networkProcessor, &wg)

		subscriptionService := startSubscriptionService(config, adapter, eventChan)

		server := startGrpcServer(adapter, config, messagingChannels, &wg)

		listenForInterruptSignal(server, subscriptionService, serviceMonitor, networkProcessor, controller, &wg)

	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&jagwServiceAddress, "jagw-service-address", "j", os.Getenv("HAWKEYE_JAGW_SERVICE_ADDRESS"), "JAGW Service Address e.g. localhost or 127.0.0.1")
	startCmd.Flags().StringVarP(&jagwRequestPort, "jagw-request-port", "r", os.Getenv("HAWKEYE_JAGW_REQUEST_PORT"), "JAGW Request Port e.g. 9903")
	startCmd.Flags().StringVarP(&jagwSubscriptionPort, "jagw-subscription-port", "s", os.Getenv("HAWKEYE_JAGW_SUBSCRIPTION_PORT"), "JAGW Subscription Port e.g. 9902")
	startCmd.Flags().StringVarP(&grpcPort, "grpc-port", "p", os.Getenv("HAWKEYE_GRPC_PORT"), "gRPC Port e.g. 10000")
	startCmd.Flags().StringVarP(&consulServerAddress, "consul-server-address", "c", os.Getenv("HAWKEYE_CONSUL_SERVER_ADDRESS"), "Consul Server Address e.g. consul-hawkv6.stud.network.garden")

	markRequiredFlags(startCmd, []string{"jagw-service-address", "jagw-request-port", "jagw-subscription-port", "grpc-port"})
}
