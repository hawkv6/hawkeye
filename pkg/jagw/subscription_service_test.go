package jagw

import (
	"context"
	"fmt"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/adapter"
	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/jalapeno-api-gateway/jagw-go/jagw"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func getExampleDeleteLsNodeEvent() *jagw.LsNodeEvent {
	return &jagw.LsNodeEvent{
		Action: proto.String("del"),
		Key:    proto.String("2_0_0_0000.0000.0004"),
	}
}

func getExampleDeleteLsLinkEvent() *jagw.LsLinkEvent {
	return &jagw.LsLinkEvent{
		Action: proto.String("del"),
		Key:    proto.String("2_0_2_0_0000.0000.0007_2001:db8:57::7_0000.0000.0005_2001:db8:57::5"),
	}
}

func getExampleDeleteLsPrefixEvent() *jagw.LsPrefixEvent {
	return &jagw.LsPrefixEvent{
		Action: proto.String("del"),
		Key:    proto.String("2_0_2_0_0_2001:db8:e5::_64_0000.0000.0005"),
	}
}

func getExampleDeleteLsSrv6SidEvent() *jagw.LsSrv6SidEvent {
	return &jagw.LsSrv6SidEvent{
		Action: proto.String("del"),
		Key:    proto.String("0_0000.0000.000a_fc00:0:a:0:1::"),
	}
}

func TestNewJagwSubscriptionService(t *testing.T) {
	tests := []struct {
		name                       string
		jagwServiceAddress         string
		subscriptionPort           uint16
		wantJagwSubscriptionSocket string
	}{
		{
			name:                       "TestNewJagwRequestService",
			jagwServiceAddress:         "localhost",
			subscriptionPort:           9903,
			wantJagwSubscriptionSocket: "localhost:9903",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetJagwServiceAddress().Return(tt.jagwServiceAddress).Times(1)
			config.EXPECT().GetJagwSubscriptionPort().Return(tt.subscriptionPort).Times(1)
			adapter := adapter.NewDomainAdapter()
			jagwSubscriptionService := NewJagwSubscriptionService(config, adapter, make(chan domain.NetworkEvent))
			assert.NotNil(t, jagwSubscriptionService)
			assert.Equal(t, tt.wantJagwSubscriptionSocket, jagwSubscriptionService.jagwSubscriptionSocket)
		})
	}
}

func TestJagwSubscriptionService_Init(t *testing.T) {
	tests := []struct {
		name               string
		jagwServiceAddress string
		subscriptionPort   uint16
		wantErr            bool
	}{
		{
			name:               "TestJagwSubscriptionService_Init success",
			jagwServiceAddress: "localhost",
			subscriptionPort:   9903,
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := adapter.NewDomainAdapter()
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetJagwServiceAddress().Return(tt.jagwServiceAddress).Times(1)
			config.EXPECT().GetJagwSubscriptionPort().Return(tt.subscriptionPort).Times(1)
			jagwSubscriptionService := NewJagwSubscriptionService(config, adapter, make(chan domain.NetworkEvent))
			err := jagwSubscriptionService.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("JagwSubscriptionService.Init() '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
		})
	}
}

func TestJagwSubscriptionService_Start_createSubscriptions(t *testing.T) {
	lsNodesSubscription := jagw.NewMockSubscriptionService_SubscribeToLsNodesClient(gomock.NewController(t))
	lsLinksSubscription := jagw.NewMockSubscriptionService_SubscribeToLsLinksClient(gomock.NewController(t))
	lsPrefixesSubscription := jagw.NewMockSubscriptionService_SubscribeToLsPrefixesClient(gomock.NewController(t))
	lsSrv6SidsSubscription := jagw.NewMockSubscriptionService_SubscribeToLsSrv6SidsClient(gomock.NewController(t))
	lsNodesSubscriptionContext, lsNodesSubscriptionCancel := context.WithCancel(context.Background())
	lsLinksSubscriptionContext, lsLinksSubscriptionCancel := context.WithCancel(context.Background())
	lsPrefixesSubscriptionContext, lsPrefixesSubscriptionCancel := context.WithCancel(context.Background())
	lsSrv6SidsSubscriptionContext, lsSrv6SidsSubscriptionCancel := context.WithCancel(context.Background())
	lsNodesSubscription.EXPECT().Context().Return(lsNodesSubscriptionContext).AnyTimes()
	lsLinksSubscription.EXPECT().Context().Return(lsLinksSubscriptionContext).AnyTimes()
	lsPrefixesSubscription.EXPECT().Context().Return(lsPrefixesSubscriptionContext).AnyTimes()
	lsSrv6SidsSubscription.EXPECT().Context().Return(lsSrv6SidsSubscriptionContext).AnyTimes()
	lsNodesSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("error receiving lsnode event")).AnyTimes()
	lsLinksSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("error receiving lslink event")).AnyTimes()
	lsPrefixesSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("error receiving lsprefix event")).AnyTimes()
	lsSrv6SidsSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("error receiving lssrv6sid event")).AnyTimes()
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwSubscriptionPort().Return(uint16(9903)).AnyTimes()
	adapter := adapter.NewDomainAdapter()

	tests := []struct {
		name    string
		wantErr bool
		setup   func(*JagwSubscriptionService)
	}{
		{
			name: "TestJagwSubscriptionService_Start and createSubscriptions success",
			setup: func(jagwSubscriptionService *JagwSubscriptionService) {
				subscriptionClient := jagw.NewMockSubscriptionServiceClient(gomock.NewController(t))
				subscriptionClient.EXPECT().SubscribeToLsNodes(gomock.Any(), gomock.Any()).Return(lsNodesSubscription, nil).AnyTimes()
				subscriptionClient.EXPECT().SubscribeToLsLinks(gomock.Any(), gomock.Any()).Return(lsLinksSubscription, nil).AnyTimes()
				subscriptionClient.EXPECT().SubscribeToLsPrefixes(gomock.Any(), gomock.Any()).Return(lsPrefixesSubscription, nil).AnyTimes()
				subscriptionClient.EXPECT().SubscribeToLsSrv6Sids(gomock.Any(), gomock.Any()).Return(lsSrv6SidsSubscription, nil).AnyTimes()

				jagwSubscriptionService.subscriptionClient = subscriptionClient
			},
		},
		{
			name:    "TestJagwSubscriptionService_Start and createSubscriptions lsnodes error",
			wantErr: true,
			setup: func(jagwRequestService *JagwSubscriptionService) {
				subscriptionClient := jagw.NewMockSubscriptionServiceClient(gomock.NewController(t))
				subscriptionClient.EXPECT().SubscribeToLsNodes(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error get lsnodes")).Times(1)
				jagwRequestService.subscriptionClient = subscriptionClient
			},
		},
		{
			name:    "TestJagwSubscriptionService_Start and createSubscriptions lslinks error",
			wantErr: true,
			setup: func(jagwRequestService *JagwSubscriptionService) {
				subscriptionClient := jagw.NewMockSubscriptionServiceClient(gomock.NewController(t))
				subscriptionClient.EXPECT().SubscribeToLsNodes(gomock.Any(), gomock.Any()).Return(lsNodesSubscription, nil).Times(1)
				subscriptionClient.EXPECT().SubscribeToLsLinks(gomock.Any(), gomock.Any()).Return(lsLinksSubscription, fmt.Errorf("Error get lslinks")).Times(1)
				jagwRequestService.subscriptionClient = subscriptionClient
			},
		},
		{
			name:    "TestJagwSubscriptionService_Start and createSubscriptions lsprefixes error",
			wantErr: true,
			setup: func(jagwRequestService *JagwSubscriptionService) {
				subscriptionClient := jagw.NewMockSubscriptionServiceClient(gomock.NewController(t))
				subscriptionClient.EXPECT().SubscribeToLsNodes(gomock.Any(), gomock.Any()).Return(lsNodesSubscription, nil).Times(1)
				subscriptionClient.EXPECT().SubscribeToLsLinks(gomock.Any(), gomock.Any()).Return(lsLinksSubscription, nil).Times(1)
				subscriptionClient.EXPECT().SubscribeToLsPrefixes(gomock.Any(), gomock.Any()).Return(lsPrefixesSubscription, fmt.Errorf("Error get lsprefixes")).Times(1)
				jagwRequestService.subscriptionClient = subscriptionClient
			},
		},
		{
			name:    "TestJagwSubscriptionService_Start and createSubscriptions lssrv6sids error",
			wantErr: true,
			setup: func(jagwRequestService *JagwSubscriptionService) {
				subscriptionClient := jagw.NewMockSubscriptionServiceClient(gomock.NewController(t))
				subscriptionClient.EXPECT().SubscribeToLsNodes(gomock.Any(), gomock.Any()).Return(lsNodesSubscription, nil).Times(1)
				subscriptionClient.EXPECT().SubscribeToLsLinks(gomock.Any(), gomock.Any()).Return(lsLinksSubscription, nil).Times(1)
				subscriptionClient.EXPECT().SubscribeToLsPrefixes(gomock.Any(), gomock.Any()).Return(lsPrefixesSubscription, nil).Times(1)
				subscriptionClient.EXPECT().SubscribeToLsSrv6Sids(gomock.Any(), gomock.Any()).Return(lsSrv6SidsSubscription, fmt.Errorf("Error get lssrv6sids")).Times(1)
				jagwRequestService.subscriptionClient = subscriptionClient
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jagwSubscriptionService := NewJagwSubscriptionService(config, adapter, make(chan domain.NetworkEvent))
			tt.setup(jagwSubscriptionService)
			err := jagwSubscriptionService.Start()
			if (err != nil) != tt.wantErr {
				t.Errorf("JagwRequestService.Start() '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
		})
	}
	lsNodesSubscriptionCancel()
	lsLinksSubscriptionCancel()
	lsPrefixesSubscriptionCancel()
	lsSrv6SidsSubscriptionCancel()
}

func TestJagwSubscriptionService_subcribeLsNodes_enqueueNodeEvent(t *testing.T) {
	nodeEvent := getExampleDeleteLsNodeEvent()
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwSubscriptionPort().Return(uint16(9903)).AnyTimes()
	adapter := adapter.NewDomainAdapter()
	lsNodesSubscription := jagw.NewMockSubscriptionService_SubscribeToLsNodesClient(gomock.NewController(t))
	lsNodesSubscriptionContext, lsNodesSubscriptionCancel := context.WithCancel(context.Background())
	lsNodesSubscription.EXPECT().Context().Return(lsNodesSubscriptionContext).AnyTimes()
	tests := []struct {
		name           string
		wantConvertErr bool
		setup          func(*JagwSubscriptionService)
	}{
		{
			name:           "TestJagwSubscriptionService_subcribeLsNodes_enqueueNodeEvent success",
			wantConvertErr: false,
			setup: func(jagwSubscriptionService *JagwSubscriptionService) {
				lsNodesSubscription.EXPECT().Recv().Return(nodeEvent, nil).AnyTimes()
				jagwSubscriptionService.lsNodesSubscription = lsNodesSubscription
			},
		},
		{
			name:           "TestJagwSubscriptionService_subcribeLsNodes_enqueueNodeEvent convert error",
			wantConvertErr: true,
			setup: func(jagwSubscriptionService *JagwSubscriptionService) {
				nodeEventCopy := nodeEvent
				nodeEventCopy.Key = nil
				lsNodesSubscription.EXPECT().Recv().Return(nodeEventCopy, nil).AnyTimes()
				jagwSubscriptionService.lsNodesSubscription = lsNodesSubscription
			},
		},
		{
			name:           "TestJagwSubscriptionService_subcribeLsNodes_enqueueNodeEvent receive error",
			wantConvertErr: true,
			setup: func(jagwSubscriptionService *JagwSubscriptionService) {
				lsNodesSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("error receiving lsnode event")).AnyTimes()
				jagwSubscriptionService.lsNodesSubscription = lsNodesSubscription
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jagwSubscriptionService := NewJagwSubscriptionService(config, adapter, make(chan domain.NetworkEvent))
			tt.setup(jagwSubscriptionService)
			go jagwSubscriptionService.subscribeLsNodes()
			if !tt.wantConvertErr {
				event := <-jagwSubscriptionService.eventChan
				assert.NotNil(t, event)
			}
		})
	}
	lsNodesSubscriptionCancel()
}

func TestJagwSubscriptionService_subcribeLsLinks_enqueueLinkEvent(t *testing.T) {
	linkEvent := getExampleDeleteLsLinkEvent()
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwSubscriptionPort().Return(uint16(9903)).AnyTimes()
	adapter := adapter.NewDomainAdapter()
	lsLinkSubscription := jagw.NewMockSubscriptionService_SubscribeToLsLinksClient(gomock.NewController(t))
	lsLinkSubscriptionContext, lsLinkSubscriptionCancel := context.WithCancel(context.Background())
	lsLinkSubscription.EXPECT().Context().Return(lsLinkSubscriptionContext).AnyTimes()
	tests := []struct {
		name           string
		wantConvertErr bool
		setup          func(*JagwSubscriptionService)
	}{
		{
			name:           "TestJagwSubscriptionService_subscribeLsLinks_enqueueLinkEvent success",
			wantConvertErr: false,
			setup: func(jagwSubscriptionService *JagwSubscriptionService) {
				lsLinkSubscription.EXPECT().Recv().Return(linkEvent, nil).AnyTimes()
				jagwSubscriptionService.lsLinksSubscription = lsLinkSubscription
			},
		},
		{
			name:           "TestJagwSubscriptionServic_subscribeLsLinks_enqueueLinkEevent convert error",
			wantConvertErr: true,
			setup: func(jagwSubscriptionService *JagwSubscriptionService) {
				linkEventCopy := linkEvent
				linkEventCopy.Key = nil
				lsLinkSubscription.EXPECT().Recv().Return(linkEventCopy, nil).AnyTimes()
				jagwSubscriptionService.lsLinksSubscription = lsLinkSubscription
			},
		},
		{
			name:           "TestJagwSubscriptionService_subscribeLsLinks_enqueueLinkEvent receive error",
			wantConvertErr: true,
			setup: func(jagwSubscriptionService *JagwSubscriptionService) {
				lsLinkSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("error receiving lslink event")).AnyTimes()
				jagwSubscriptionService.lsLinksSubscription = lsLinkSubscription
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jagwSubscriptionService := NewJagwSubscriptionService(config, adapter, make(chan domain.NetworkEvent))
			tt.setup(jagwSubscriptionService)
			go jagwSubscriptionService.subscribeLsLinks()
			if !tt.wantConvertErr {
				event := <-jagwSubscriptionService.eventChan
				assert.NotNil(t, event)
			}
		})
	}
	lsLinkSubscriptionCancel()
}

func TestJagwSubscriptionService_subscribeLsPrefixes_enqueuePrefixEvent(t *testing.T) {
	prefixEvent := getExampleDeleteLsPrefixEvent()
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwSubscriptionPort().Return(uint16(9903)).AnyTimes()
	adapter := adapter.NewDomainAdapter()
	lsPrefixesSubscription := jagw.NewMockSubscriptionService_SubscribeToLsPrefixesClient(gomock.NewController(t))
	lsPrefixesSubscriptionContext, lsPrefixesSubscriptionCancel := context.WithCancel(context.Background())
	lsPrefixesSubscription.EXPECT().Context().Return(lsPrefixesSubscriptionContext).AnyTimes()
	tests := []struct {
		name           string
		wantConvertErr bool
		setup          func(*JagwSubscriptionService)
	}{
		{
			name:           "TestJagwSubscriptionService_subscribeLsPrefixes_enqueuePrefixEvent success",
			wantConvertErr: false,
			setup: func(jagwSubscriptionService *JagwSubscriptionService) {
				lsPrefixesSubscription.EXPECT().Recv().Return(prefixEvent, nil).AnyTimes()
				jagwSubscriptionService.lsPrefixesSubscription = lsPrefixesSubscription
			},
		},
		{
			name:           "TestJagwSubscriptionService_subscribe LsPrefixes_enqueuePrefixEvent convert error",
			wantConvertErr: true,
			setup: func(jagwSubscriptionService *JagwSubscriptionService) {
				prefixEventCopy := prefixEvent
				prefixEventCopy.Key = nil
				lsPrefixesSubscription.EXPECT().Recv().Return(prefixEventCopy, nil).AnyTimes()
				jagwSubscriptionService.lsPrefixesSubscription = lsPrefixesSubscription
			},
		},
		{
			name:           "TestJagwSubscriptionService_subscribeLsPrefixes_enqueuePrefixEvent receive error",
			wantConvertErr: true,
			setup: func(jagwSubscriptionService *JagwSubscriptionService) {
				lsPrefixesSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("error receiving lsprefix event")).AnyTimes()
				jagwSubscriptionService.lsPrefixesSubscription = lsPrefixesSubscription
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jagwSubscriptionService := NewJagwSubscriptionService(config, adapter, make(chan domain.NetworkEvent))
			tt.setup(jagwSubscriptionService)
			go jagwSubscriptionService.subscribeLsPrefixes()
			if !tt.wantConvertErr {
				event := <-jagwSubscriptionService.eventChan
				assert.NotNil(t, event)
			}
		})
	}
	lsPrefixesSubscriptionCancel()
}

func TestJagwSubscriptionService_subscribeLsSrv6Sids_enqueueSrv6SidEvent(t *testing.T) {
	srv6SidEvent := getExampleDeleteLsSrv6SidEvent()
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwSubscriptionPort().Return(uint16(9903)).AnyTimes()
	adapter := adapter.NewDomainAdapter()
	lsSrv6SidsSubscription := jagw.NewMockSubscriptionService_SubscribeToLsSrv6SidsClient(gomock.NewController(t))
	lsSrv6SidsSubscriptionContext, lsSrv6SidsSubscriptionCancel := context.WithCancel(context.Background())
	lsSrv6SidsSubscription.EXPECT().Context().Return(lsSrv6SidsSubscriptionContext).AnyTimes()
	tests := []struct {
		name           string
		wantConvertErr bool
		setup          func(*JagwSubscriptionService)
	}{
		{
			name:           "TestJagwSubscriptionService_subscribeLsSrv6Sids_enqueueSrv6SidEvent success",
			wantConvertErr: false,
			setup: func(jagwSubscriptionService *JagwSubscriptionService) {
				lsSrv6SidsSubscription.EXPECT().Recv().Return(srv6SidEvent, nil).AnyTimes()
				jagwSubscriptionService.lsSrv6SidsSubscription = lsSrv6SidsSubscription
			},
		},
		{
			name:           "TestJagwSubscriptionService_subscribeLsSrv6Sids_enqueueSrv6SidEvent convert error",
			wantConvertErr: true,
			setup: func(jagwSubscriptionService *JagwSubscriptionService) {
				srv6SidEventCopy := srv6SidEvent
				srv6SidEventCopy.Key = nil
				lsSrv6SidsSubscription.EXPECT().Recv().Return(srv6SidEventCopy, nil).AnyTimes()
				jagwSubscriptionService.lsSrv6SidsSubscription = lsSrv6SidsSubscription
			},
		},
		{
			name:           "TestJagwSubscriptionService_subscribeLsSrv6Sids_enqueueSrv6SidEvent receive error",
			wantConvertErr: true,
			setup: func(jagwSubscriptionService *JagwSubscriptionService) {
				lsSrv6SidsSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("error receiving lssrv6sid event")).AnyTimes()
				jagwSubscriptionService.lsSrv6SidsSubscription = lsSrv6SidsSubscription
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jagwSubscriptionService := NewJagwSubscriptionService(config, adapter, make(chan domain.NetworkEvent))
			tt.setup(jagwSubscriptionService)
			go jagwSubscriptionService.subscribeLsSrv6Sids()
			if !tt.wantConvertErr {
				event := <-jagwSubscriptionService.eventChan
				assert.NotNil(t, event)
			}
		})
	}
	lsSrv6SidsSubscriptionCancel()
}
