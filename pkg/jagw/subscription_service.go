package jagw

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/hawkv6/hawkeye/pkg/adapter"
	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/logging"
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
	eventChan              chan domain.NetworkEvent
	lsNodesSubscription    jagw.SubscriptionService_SubscribeToLsNodesClient
	lsLinksSubscription    jagw.SubscriptionService_SubscribeToLsLinksClient
	lsPrefixesSubscription jagw.SubscriptionService_SubscribeToLsPrefixesClient
	lsSrv6SidsSubscription jagw.SubscriptionService_SubscribeToLsSrv6SidsClient
	cancelFunctions        []context.CancelFunc
	wg                     sync.WaitGroup
}

func NewJagwSubscriptionService(config config.Config, adapter adapter.Adapter, eventChan chan domain.NetworkEvent) *JagwSubscriptionService {
	return &JagwSubscriptionService{
		log:                    logging.DefaultLogger.WithField("subsystem", Subsystem),
		jagwSubscriptionSocket: config.GetJagwServiceAddress() + ":" + strconv.FormatUint(uint64(config.GetJagwSubscriptionPort()), 10),
		adapter:                adapter,
		eventChan:              eventChan,
		cancelFunctions:        make([]context.CancelFunc, 0),
		wg:                     sync.WaitGroup{},
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
	grpcClientConnection, err := grpc.NewClient(subscriptionService.jagwSubscriptionSocket,
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
	ctx, cancel := context.WithCancel(context.Background())
	subscription := &jagw.TopologySubscription{
		Properties: helper.GetLsNodeProperties(),
	}
	stream, err := subscriptionService.subscriptionClient.SubscribeToLsNodes(ctx, subscription)
	if err != nil {
		cancel()
		return fmt.Errorf("Error when calling SubscribeToLsNodes on SubscriptionService: %s", err)
	}
	subscriptionService.lsNodesSubscription = stream
	subscriptionService.cancelFunctions = append(subscriptionService.cancelFunctions, cancel)
	return nil
}

func (subscriptionService *JagwSubscriptionService) enqueueNodeEvent(lsNodeEvent *jagw.LsNodeEvent) {
	event, err := subscriptionService.adapter.ConvertNodeEvent(lsNodeEvent)
	if err != nil {
		subscriptionService.log.Errorf("Error converting LsNodeEvent: %s", err)
		return
	}
	subscriptionService.eventChan <- event
}

func (subscriptionService *JagwSubscriptionService) subscribeLsNodes() {
	subscriptionService.log.Debugln("Subscribing to LsNodes")
	ctx := subscriptionService.lsNodesSubscription.Context()
	subscriptionService.wg.Add(1)
	for {
		event, err := subscriptionService.lsNodesSubscription.Recv()
		select {
		case <-ctx.Done():
			subscriptionService.log.Debugln("LsNode stream ended")
			subscriptionService.wg.Done()
			return
		default:
			if err != nil {
				subscriptionService.log.Errorf("Error when receiving LsNode event: %s", err)
			}
			subscriptionService.enqueueNodeEvent(event)
		}
	}
}

func (subscriptionService *JagwSubscriptionService) createLsLinksSubscription() error {
	ctx, cancel := context.WithCancel(context.Background())
	subscription := &jagw.TopologySubscription{
		Properties: helper.GetLsLinkProperties(),
	}
	stream, err := subscriptionService.subscriptionClient.SubscribeToLsLinks(ctx, subscription)
	if err != nil {
		cancel()
		return fmt.Errorf("Error when calling SubscribeToLsLinks on SubscriptionService: %s", err)
	}
	subscriptionService.lsLinksSubscription = stream
	subscriptionService.cancelFunctions = append(subscriptionService.cancelFunctions, cancel)
	return nil
}

func (subscriptionService *JagwSubscriptionService) enqueueLinkEvent(lsLinkEvent *jagw.LsLinkEvent) {
	event, err := subscriptionService.adapter.ConvertLinkEvent(lsLinkEvent)
	if err != nil {
		subscriptionService.log.Errorf("Error converting LsLinkEvent: %s", err)
		return
	}
	subscriptionService.eventChan <- event
}

func (subscriptionService *JagwSubscriptionService) subscribeLsLinks() {
	subscriptionService.log.Debugln("Subscribing to LsLinks")
	ctx := subscriptionService.lsLinksSubscription.Context()
	subscriptionService.wg.Add(1)
	for {
		event, err := subscriptionService.lsLinksSubscription.Recv()
		select {
		case <-ctx.Done():
			subscriptionService.log.Debugln("LsLink stream ended")
			subscriptionService.wg.Done()
			return
		default:
			if err != nil {
				subscriptionService.log.Errorf("Error when receiving LsLink event: %s", err)
				continue
			}
			subscriptionService.enqueueLinkEvent(event)
		}
	}
}

func (subscriptionService *JagwSubscriptionService) createLsPrefixesSubscription() error {
	ctx, cancel := context.WithCancel(context.Background())
	subscription := &jagw.TopologySubscription{
		Properties: helper.GetLsPrefixProperties(),
	}
	stream, err := subscriptionService.subscriptionClient.SubscribeToLsPrefixes(ctx, subscription)
	if err != nil {
		cancel()
		return fmt.Errorf("Error when calling SubscribeToLsPrefix on SubscriptionService: %s", err)
	}
	subscriptionService.lsPrefixesSubscription = stream
	subscriptionService.cancelFunctions = append(subscriptionService.cancelFunctions, cancel)
	return nil
}

func (subscriptionService *JagwSubscriptionService) enqueuePrefixEvent(lsPrefixEvent *jagw.LsPrefixEvent) {
	event, err := subscriptionService.adapter.ConvertPrefixEvent(lsPrefixEvent)
	if err != nil {
		subscriptionService.log.Errorf("Error converting LsPrefixEvent: %s", err)
		return
	}
	subscriptionService.eventChan <- event
}

func (subscriptionService *JagwSubscriptionService) subscribeLsPrefixes() {
	subscriptionService.log.Debugln("Subscribing to LsPrefix")
	ctx := subscriptionService.lsPrefixesSubscription.Context()
	subscriptionService.wg.Add(1)
	for {
		event, err := subscriptionService.lsPrefixesSubscription.Recv()
		select {
		case <-ctx.Done():
			subscriptionService.log.Debugln("LsPrefix stream ended")
			subscriptionService.wg.Done()
			return
		default:
			if err != nil {
				subscriptionService.log.Errorf("Error when receiving LsPrefix event: %s", err)
				continue
			}
			subscriptionService.enqueuePrefixEvent(event)
		}
	}
}

func (subscriptionService *JagwSubscriptionService) createLsSrv6SidsSubscription() error {
	ctx, cancel := context.WithCancel(context.Background())
	subscription := &jagw.TopologySubscription{
		Properties: helper.GetLsSrv6SidsProperties(),
	}
	stream, err := subscriptionService.subscriptionClient.SubscribeToLsSrv6Sids(ctx, subscription)
	if err != nil {
		cancel()
		return fmt.Errorf("Error when calling SubscribeToLsSrv6Sids on SubscriptionService: %s", err)
	}
	subscriptionService.lsSrv6SidsSubscription = stream
	subscriptionService.cancelFunctions = append(subscriptionService.cancelFunctions, cancel)
	return nil
}

func (subscriptionService *JagwSubscriptionService) enqueueSrv6Sids(lsSrv6SidsEvent *jagw.LsSrv6SidEvent) {
	event, err := subscriptionService.adapter.ConvertSidEvent(lsSrv6SidsEvent)
	if err != nil {
		subscriptionService.log.Errorf("Error converting LsSrv6SidsEvent: %s", err)
		return
	}
	subscriptionService.eventChan <- event
}

func (subscriptionService *JagwSubscriptionService) subscribeLsSrv6Sids() {
	subscriptionService.log.Debugln("Subscribing to LsSrv6Sids")
	ctx := subscriptionService.lsSrv6SidsSubscription.Context()
	subscriptionService.wg.Add(1)
	for {
		event, err := subscriptionService.lsSrv6SidsSubscription.Recv()
		select {
		case <-ctx.Done():
			subscriptionService.log.Debugln("LsSrv6Sids stream ended")
			subscriptionService.wg.Done()
			return
		default:
			if err == io.EOF {
				subscriptionService.log.Debugln("LsSrv6Sids stream ended")
				return
			}
			if err != nil {
				subscriptionService.log.Errorf("Error when receiving LsSrv6Sids event: %s", err)
				continue
			}
			subscriptionService.enqueueSrv6Sids(event)
		}
	}
}

func (subscriptionService *JagwSubscriptionService) Stop() {
	subscriptionService.log.Infoln("Stopping JAGW Subscription Service")
	for _, cancel := range subscriptionService.cancelFunctions {
		cancel()
	}
	subscriptionService.wg.Wait()
}
