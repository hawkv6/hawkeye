package jagw

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hawkv6/hawkeye/pkg/adapter"
	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/hawkv6/hawkeye/pkg/processor"
	"github.com/jalapeno-api-gateway/jagw-go/jagw"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type JagwSubscriptionService struct {
	log                    *logrus.Entry
	jagwSubscriptionSocket string
	grpcClientConnection   *grpc.ClientConn
	subscriptionClient     jagw.SubscriptionServiceClient
	adapter                adapter.Adapter
	processor              processor.Processor
	quitChan               chan struct{}
	helper                 helper.Helper
	lsNodesSubscription    jagw.SubscriptionService_SubscribeToLsNodesClient
	lsLinksSubscription    jagw.SubscriptionService_SubscribeToLsLinksClient
	lsPrefixesSubscription jagw.SubscriptionService_SubscribeToLsPrefixesClient
	lsSrv6SidsSubscription jagw.SubscriptionService_SubscribeToLsSrv6SidsClient
}

func NewJagwSubscriptionService(config config.Config, adapter adapter.Adapter, processor processor.Processor, helper helper.Helper) *JagwSubscriptionService {
	return &JagwSubscriptionService{
		log:                    logging.DefaultLogger.WithField("subsystem", Subsystem),
		jagwSubscriptionSocket: config.GetJagwServiceAddress() + ":" + strconv.FormatUint(uint64(config.GetJagwSubscriptionPort()), 10),
		adapter:                adapter,
		processor:              processor,
		quitChan:               make(chan struct{}),
		helper:                 helper,
	}
}

func (subscriptionService *JagwSubscriptionService) Init() error {
	subscriptionService.log.Debugln("Initializing JAGW Subscription Service")
	if err := subscriptionService.createSubscriptionClient(); err != nil {
		return fmt.Errorf("Error creating subscription client: %s", err)
	}
	return nil
}

func (subscriptionService *JagwSubscriptionService) createSubscriptionClient() error {
	subscriptionService.log.Debugln("Initializing gRPC client connection")
	grpcClientConnection, err := grpc.Dial(subscriptionService.jagwSubscriptionSocket,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("Error when dialing gRPC server: %s", err)
	}
	subscriptionService.grpcClientConnection = grpcClientConnection
	subscriptionService.subscriptionClient = jagw.NewSubscriptionServiceClient(grpcClientConnection)
	return nil
}

func (subscriptionService *JagwSubscriptionService) Start() error {
	subscriptionService.log.Infoln("Starting JAGW Subscription Service")
	if err := subscriptionService.createSubscriptions(); err != nil {
		return fmt.Errorf("Error creating subscriptions: %s", err)
	}
	go subscriptionService.subscribeLsNodes()
	go subscriptionService.subscribeLsLinks()
	go subscriptionService.subscribeLsPrefixes()
	go subscriptionService.subscribeLsSrv6Sids()
	return nil
}

func (subscriptionService *JagwSubscriptionService) createSubscriptions() error {
	if err := subscriptionService.createLsNodesSubscription(); err != nil {
		return fmt.Errorf("Error initializing LsNode subscription: %s", err)
	}
	if err := subscriptionService.createLsLinksSubscription(); err != nil {
		return fmt.Errorf("Error initializing LsLink subscription: %s", err)
	}
	if err := subscriptionService.createLsPrefixesSubscription(); err != nil {
		return fmt.Errorf("Error initializing LsPrefix subscription: %s", err)
	}
	if err := subscriptionService.createLsSrv6SidsSubscription(); err != nil {
		return fmt.Errorf("Error initializing LsSrv6Sids subscription: %s", err)
	}
	return nil
}

func (subscriptionService *JagwSubscriptionService) createLsNodesSubscription() error {
	ctx := context.Background()
	subscription := &jagw.TopologySubscription{
		Properties: subscriptionService.helper.GetLsNodeProperties(),
	}
	stream, err := subscriptionService.subscriptionClient.SubscribeToLsNodes(ctx, subscription)
	if err != nil {
		return fmt.Errorf("Error when calling SubscribeToLsNodes on SubscriptionService: %s", err)
	}
	subscriptionService.lsNodesSubscription = stream
	return nil
}

func (subscriptionService *JagwSubscriptionService) convertLsNodeEvent(lsNodeEvent *jagw.LsNodeEvent) {
	event, err := subscriptionService.adapter.ConvertNodeEvent(lsNodeEvent)
	if err != nil {
		subscriptionService.log.Errorf("Error converting LsNodeEvent: %s", err)
		return
	}
	subscriptionService.log.Debugln(event)
}

func (subscriptionService *JagwSubscriptionService) subscribeLsNodes() {
	subscriptionService.log.Debugln("Subscribing to LsNodes")
	for {
		select {
		case <-subscriptionService.quitChan:
			return
		default:
			event, err := subscriptionService.lsNodesSubscription.Recv()
			if err != nil {
				subscriptionService.log.Errorf("Error when receiving LsNode event: %s", err)
			}
			subscriptionService.convertLsNodeEvent(event)
		}
	}
}

func (subscriptionService *JagwSubscriptionService) createLsLinksSubscription() error {
	subscriptionService.log.Debugln("Subscribing to LsLinks")
	ctx := context.Background()
	subscription := &jagw.TopologySubscription{
		Properties: subscriptionService.helper.GetLsLinkProperties(),
	}
	stream, err := subscriptionService.subscriptionClient.SubscribeToLsLinks(ctx, subscription)
	if err != nil {
		return fmt.Errorf("Error when calling SubscribeToLsLinks on SubscriptionService: %s", err)
	}
	subscriptionService.lsLinksSubscription = stream
	return nil
}

func (subscriptionService *JagwSubscriptionService) convertLsLink(lsLinkEvent *jagw.LsLinkEvent) {
	event, err := subscriptionService.adapter.ConvertLinkEvent(lsLinkEvent)
	if err != nil {
		subscriptionService.log.Errorf("Error converting LsLinkEvent: %s", err)
		return
	}
	subscriptionService.log.Debugln(event)
}

func (subscriptionService *JagwSubscriptionService) subscribeLsLinks() {
	for {
		select {
		case <-subscriptionService.quitChan:
			return
		default:
			event, err := subscriptionService.lsLinksSubscription.Recv()
			if err != nil {
				subscriptionService.log.Errorf("Error when receiving LsLink event: %s", err)
				continue
			}
			subscriptionService.convertLsLink(event)
		}
	}
}

func (subscriptionService *JagwSubscriptionService) createLsPrefixesSubscription() error {
	subscriptionService.log.Debugln("Subscribing to LsPrefix")
	ctx := context.Background()
	subscription := &jagw.TopologySubscription{
		Properties: subscriptionService.helper.GetLsPrefixProperties(),
	}
	stream, err := subscriptionService.subscriptionClient.SubscribeToLsPrefixes(ctx, subscription)
	if err != nil {
		return fmt.Errorf("Error when calling SubscribeToLsPrefix on SubscriptionService: %s", err)
	}
	subscriptionService.lsPrefixesSubscription = stream
	return nil
}

func (subscriptionService *JagwSubscriptionService) convertLsPrefix(lsPrefixEvent *jagw.LsPrefixEvent) {
	event, err := subscriptionService.adapter.ConvertPrefixEvent(lsPrefixEvent)
	if err != nil {
		subscriptionService.log.Errorf("Error converting LsPrefixEvent: %s", err)
		return
	}
	subscriptionService.log.Debugln(event)
}

func (subscriptionService *JagwSubscriptionService) subscribeLsPrefixes() {
	for {
		select {
		case <-subscriptionService.quitChan:
			return
		default:
			event, err := subscriptionService.lsPrefixesSubscription.Recv()
			if err != nil {
				subscriptionService.log.Errorf("Error when receiving LsPrefix event: %s", err)
				continue
			}
			subscriptionService.convertLsPrefix(event)
		}
	}
}

func (subscriptionService *JagwSubscriptionService) createLsSrv6SidsSubscription() error {
	subscriptionService.log.Debugln("Subscribing to LsSrv6Sids")
	ctx := context.Background()
	subscription := &jagw.TopologySubscription{
		Properties: subscriptionService.helper.GetLsSrv6SidsProperties(),
	}
	stream, err := subscriptionService.subscriptionClient.SubscribeToLsSrv6Sids(ctx, subscription)
	if err != nil {
		return fmt.Errorf("Error when calling SubscribeToLsSrv6Sids on SubscriptionService: %s", err)
	}
	subscriptionService.lsSrv6SidsSubscription = stream
	return nil
}

func (subscriptionService *JagwSubscriptionService) convertLsSrv6Sids(lsSrv6SidsEvent *jagw.LsSrv6SidEvent) {
	event, err := subscriptionService.adapter.ConvertSidEvent(lsSrv6SidsEvent)
	if err != nil {
		subscriptionService.log.Errorf("Error converting LsSrv6SidEvent: %s", err)
		return
	}
	subscriptionService.log.Debugln(event)
}

func (subscriptionService *JagwSubscriptionService) subscribeLsSrv6Sids() {
	for {
		select {
		case <-subscriptionService.quitChan:
			return
		default:
			event, err := subscriptionService.lsSrv6SidsSubscription.Recv()
			if err != nil {
				subscriptionService.log.Errorf("Error when receiving LsSrv6Sids event: %s", err)
				continue
			}
			subscriptionService.convertLsSrv6Sids(event)
		}
	}
}

func (subscriptionService *JagwSubscriptionService) Stop() {
	subscriptionService.log.Infoln("Closing JAGW Subscription Service")
	subscriptionService.grpcClientConnection.Close()
	close(subscriptionService.quitChan)
}
