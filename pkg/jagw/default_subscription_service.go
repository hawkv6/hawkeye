package jagw

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hawkv6/hawkeye/pkg/adapter"
	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/hawkv6/hawkeye/pkg/processor"
	"github.com/jalapeno-api-gateway/jagw-go/jagw"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DefaultJagwSubscriptionService struct {
	log                    *logrus.Entry
	jagwSubscriptionSocket string
	grpcClientConnection   *grpc.ClientConn
	subscriptonClient      jagw.SubscriptionServiceClient
	adapter                adapter.Adapter
	processor              processor.Processor
	quitChan               chan struct{}
}

func NewDefaultJagwSubscriptionService(config config.Config, adapter adapter.Adapter, processor processor.Processor) *DefaultJagwSubscriptionService {
	return &DefaultJagwSubscriptionService{
		log:                    logging.DefaultLogger.WithField("subsystem", Subsystem),
		jagwSubscriptionSocket: config.GetJagwServiceAddress() + ":" + strconv.FormatUint(uint64(config.GetJagwSubscriptionPort()), 10),
		adapter:                adapter,
		processor:              processor,
		quitChan:               make(chan struct{}),
	}
}

func (subscriptionService *DefaultJagwSubscriptionService) Start() error {
	subscriptionService.log.Debugln("Initializing JAGW Subscription Service")
	grpcClientConnection, err := grpc.Dial(subscriptionService.jagwSubscriptionSocket,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	subscriptionService.grpcClientConnection = grpcClientConnection
	subscriptionService.subscriptonClient = jagw.NewSubscriptionServiceClient(grpcClientConnection)
	return nil
}

func prettyPrint(any interface{}) {
	s, _ := json.MarshalIndent(any, "", "  ")
	fmt.Printf("%s\n\n", string(s))
}

func (subscriptionService *DefaultJagwSubscriptionService) SubscribeLsLinks() error {
	subscriptionService.log.Infoln("Subscribing to LsLinks")
	ctx := context.Background()
	subscription := &jagw.TopologySubscription{}
	stream, err := subscriptionService.subscriptonClient.SubscribeToLsLinks(ctx, subscription)
	if err != nil {
		return fmt.Errorf("Error when calling SubscribeToLsLinks on SubscriptionService: %s", err)
	}

	for {
		event, err := stream.Recv()
		if err != nil {
			subscriptionService.log.Errorf("Error when receiving LsLink event: %s", err)
		}
		prettyPrint(event)
	}

	// go func() {
	// 	for {
	// 		select {
	// 		case <-subscriptionService.quitChan:
	// 			return
	// 		default:
	// 			event, err := stream.Recv()
	// 			if err != nil {
	// 				subscriptionService.log.Errorf("Error when receiving LsLink event: %s", err)
	// 				return
	// 			}
	// 			prettyPrint(event)
	// 		}
	// 	}
	// }()
	// return nil
}

func (subscriptionService *DefaultJagwSubscriptionService) Stop() {
	subscriptionService.log.Infoln("Closing JAGW Subscription Service")
	close(subscriptionService.quitChan)
	subscriptionService.grpcClientConnection.Close()
}
