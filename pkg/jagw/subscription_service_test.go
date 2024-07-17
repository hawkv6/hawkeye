package jagw

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hawkv6/hawkeye/pkg/adapter"
	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/jalapeno-api-gateway/jagw-go/jagw"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

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
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwSubscriptionPort().Return(uint16(9903)).AnyTimes()
	tests := []struct {
		name           string
		wantConvertErr bool
		wantReceiveErr bool
	}{
		{
			name:           "TestJagwSubscriptionService_subcribeLsNodes_enqueueNodeEvent success",
			wantConvertErr: false,
			wantReceiveErr: false,
		},
		{
			name:           "TestJagwSubscriptionService_subcribeLsNodes_enqueueNodeEvent convert error",
			wantConvertErr: true,
			wantReceiveErr: false,
		},
		{
			name:           "TestJagwSubscriptionService_subcribeLsNodes_enqueueNodeEvent receive error",
			wantConvertErr: false,
			wantReceiveErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := adapter.NewMockAdapter(gomock.NewController(t))
			jagwSubscriptionService := NewJagwSubscriptionService(config, adapter, make(chan domain.NetworkEvent))
			lsNodesSubscription := jagw.NewMockSubscriptionService_SubscribeToLsNodesClient(gomock.NewController(t))
			if !tt.wantReceiveErr && !tt.wantConvertErr {
				adapter.EXPECT().ConvertNodeEvent(gomock.Any()).Return(domain.NewDeleteNodeEvent("key"), nil).AnyTimes()
				lsNodesSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("Closed connection")).AnyTimes().After(lsNodesSubscription.EXPECT().Recv().Return(nil, nil).Times(1))
			}
			if tt.wantConvertErr && !tt.wantReceiveErr {
				adapter.EXPECT().ConvertNodeEvent(gomock.Any()).Return(nil, fmt.Errorf("error converting lsnode to node")).AnyTimes()
				lsNodesSubscription.EXPECT().Recv().Return(nil, nil).AnyTimes()
			}
			if !tt.wantConvertErr && tt.wantReceiveErr {
				lsNodesSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("error getting lsnode event")).AnyTimes()
			}

			jagwSubscriptionService.lsNodesSubscription = lsNodesSubscription
			lsNodesSubscriptionContext, lsNodesSubscriptionCancel := context.WithCancel(context.Background())
			lsNodesSubscription.EXPECT().Context().Return(lsNodesSubscriptionContext).AnyTimes()

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				jagwSubscriptionService.subscribeLsNodes()
				wg.Done()
			}()
			if !tt.wantConvertErr && !tt.wantReceiveErr {
				event := <-jagwSubscriptionService.eventChan
				assert.NotNil(t, event)
			}
			time.Sleep(100 * time.Millisecond)
			lsNodesSubscriptionCancel()
			wg.Wait()
		})
	}
}

func TestJagwSubscriptionService_subcribeLsLinks_enqueueLinkEvent(t *testing.T) {
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwSubscriptionPort().Return(uint16(9903)).AnyTimes()
	tests := []struct {
		name            string
		wantConvertErr  bool
		wantReceivceErr bool
	}{
		{
			name:            "TestJagwSubscriptionService_subscribeLsLinks_enqueueLinkEvent success",
			wantConvertErr:  false,
			wantReceivceErr: false,
		},
		{
			name:            "TestJagwSubscriptionServic_subscribeLsLinks_enqueueLinkEevent convert error",
			wantConvertErr:  true,
			wantReceivceErr: false,
		},
		{
			name:            "TestJagwSubscriptionService_subscribeLsLinks_enqueueLinkEvent receive error",
			wantConvertErr:  false,
			wantReceivceErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := adapter.NewMockAdapter(gomock.NewController(t))
			jagwSubscriptionService := NewJagwSubscriptionService(config, adapter, make(chan domain.NetworkEvent))
			lsLinkSubscription := jagw.NewMockSubscriptionService_SubscribeToLsLinksClient(gomock.NewController(t))
			if !tt.wantConvertErr && !tt.wantReceivceErr {
				adapter.EXPECT().ConvertLinkEvent(gomock.Any()).Return(domain.NewDeleteLinkEvent("key"), nil).AnyTimes()
				lsLinkSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("Closed connection")).AnyTimes().After(lsLinkSubscription.EXPECT().Recv().Return(nil, nil).Times(1))
			}
			if tt.wantConvertErr && !tt.wantReceivceErr {
				adapter.EXPECT().ConvertLinkEvent(gomock.Any()).Return(nil, fmt.Errorf("error converting lslink to link")).AnyTimes()
				lsLinkSubscription.EXPECT().Recv().Return(nil, nil).AnyTimes()
			}
			if !tt.wantConvertErr && tt.wantReceivceErr {
				lsLinkSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("error receiving lslink event")).AnyTimes()
			}
			jagwSubscriptionService.lsLinksSubscription = lsLinkSubscription

			lsLinkSubscriptionContext, lsLinkSubscriptionCancel := context.WithCancel(context.Background())
			lsLinkSubscription.EXPECT().Context().Return(lsLinkSubscriptionContext).AnyTimes()
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				jagwSubscriptionService.subscribeLsLinks()
				wg.Done()
			}()
			if !tt.wantConvertErr && !tt.wantReceivceErr {
				event := <-jagwSubscriptionService.eventChan
				assert.NotNil(t, event)
			}
			time.Sleep(100 * time.Millisecond)
			lsLinkSubscriptionCancel()
			wg.Wait()
		})
	}
}

func TestJagwSubscriptionService_subscribeLsPrefixes_enqueuePrefixEvent(t *testing.T) {
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwSubscriptionPort().Return(uint16(9903)).AnyTimes()
	tests := []struct {
		name           string
		wantConvertErr bool
		wantReceiveErr bool
	}{
		{
			name:           "TestJagwSubscriptionService_subscribeLsPrefixes_enqueuePrefixEvent success",
			wantConvertErr: false,
			wantReceiveErr: false,
		},
		{
			name:           "TestJagwSubscriptionService_subscribe LsPrefixes_enqueuePrefixEvent convert error",
			wantConvertErr: true,
			wantReceiveErr: false,
		},
		{
			name:           "TestJagwSubscriptionService_subscribeLsPrefixes_enqueuePrefixEvent receive error",
			wantConvertErr: false,
			wantReceiveErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := adapter.NewMockAdapter(gomock.NewController(t))
			jagwSubscriptionService := NewJagwSubscriptionService(config, adapter, make(chan domain.NetworkEvent))
			lsPrefixesSubscription := jagw.NewMockSubscriptionService_SubscribeToLsPrefixesClient(gomock.NewController(t))
			if !tt.wantConvertErr && !tt.wantReceiveErr {
				adapter.EXPECT().ConvertPrefixEvent(gomock.Any()).Return(domain.NewDeletePrefixEvent("key"), nil).AnyTimes()
				lsPrefixesSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("Closed connection")).AnyTimes().After(lsPrefixesSubscription.EXPECT().Recv().Return(nil, nil).Times(1))
			}
			if tt.wantConvertErr && !tt.wantReceiveErr {
				adapter.EXPECT().ConvertPrefixEvent(gomock.Any()).Return(nil, fmt.Errorf("error converting lsprefix to prefix")).AnyTimes()
				lsPrefixesSubscription.EXPECT().Recv().Return(nil, nil).AnyTimes()
			}
			if !tt.wantConvertErr && tt.wantReceiveErr {
				lsPrefixesSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("error receiving lsprefix event")).AnyTimes()
			}
			jagwSubscriptionService.lsPrefixesSubscription = lsPrefixesSubscription
			lsPrefixesSubscriptionContext, lsPrefixesSubscriptionCancel := context.WithCancel(context.Background())
			lsPrefixesSubscription.EXPECT().Context().Return(lsPrefixesSubscriptionContext).AnyTimes()

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				jagwSubscriptionService.subscribeLsPrefixes()
				wg.Done()
			}()
			if !tt.wantConvertErr && !tt.wantReceiveErr {
				event := <-jagwSubscriptionService.eventChan
				assert.NotNil(t, event)
			}
			time.Sleep(100 * time.Millisecond)
			lsPrefixesSubscriptionCancel()
		})
	}
}

func TestJagwSubscriptionService_subscribeLsSrv6Sids_enqueueSrv6SidEvent(t *testing.T) {
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwSubscriptionPort().Return(uint16(9903)).AnyTimes()
	tests := []struct {
		name           string
		wantConvertErr bool
		wantReceiveErr bool
	}{
		{
			name:           "TestJagwSubscriptionService_subscribeLsSrv6Sids_enqueueSrv6SidEvent success",
			wantConvertErr: false,
			wantReceiveErr: false,
		},
		{
			name:           "TestJagwSubscriptionService_subscribeLsSrv6Sids_enqueueSrv6SidEvent convert error",
			wantConvertErr: true,
			wantReceiveErr: false,
		},
		{
			name:           "TestJagwSubscriptionService_subscribeLsSrv6Sids_enqueueSrv6SidEvent receive error",
			wantConvertErr: false,
			wantReceiveErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := adapter.NewMockAdapter(gomock.NewController(t))
			lsSrv6SidsSubscription := jagw.NewMockSubscriptionService_SubscribeToLsSrv6SidsClient(gomock.NewController(t))
			jagwSubscriptionService := NewJagwSubscriptionService(config, adapter, make(chan domain.NetworkEvent))
			if !tt.wantConvertErr && !tt.wantReceiveErr {
				adapter.EXPECT().ConvertSidEvent(gomock.Any()).Return(domain.NewDeleteSidEvent("key"), nil).AnyTimes()
				lsSrv6SidsSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("Closed connection")).AnyTimes().After(lsSrv6SidsSubscription.EXPECT().Recv().Return(nil, nil).Times(1))
			}
			if tt.wantConvertErr && !tt.wantReceiveErr {

				adapter.EXPECT().ConvertSidEvent(gomock.Any()).Return(nil, fmt.Errorf("error converting lssrv6sid to sid")).AnyTimes()
				lsSrv6SidsSubscription.EXPECT().Recv().Return(nil, nil).AnyTimes()
			}
			if !tt.wantConvertErr && tt.wantReceiveErr {
				lsSrv6SidsSubscription.EXPECT().Recv().Return(nil, fmt.Errorf("error receiving lssrv6sid event")).AnyTimes()
			}

			jagwSubscriptionService.lsSrv6SidsSubscription = lsSrv6SidsSubscription
			lsSrv6SidsSubscriptionContext, lsSrv6SidsSubscriptionCancel := context.WithCancel(context.Background())
			lsSrv6SidsSubscription.EXPECT().Context().Return(lsSrv6SidsSubscriptionContext).AnyTimes()

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				jagwSubscriptionService.subscribeLsSrv6Sids()
				wg.Done()
			}()
			if !tt.wantConvertErr && !tt.wantReceiveErr {
				event := <-jagwSubscriptionService.eventChan
				assert.NotNil(t, event)
			}
			time.Sleep(100 * time.Millisecond)
			lsSrv6SidsSubscriptionCancel()
			wg.Wait()
		})
	}
}

func TestJagwSubscriptionService_Stop(t *testing.T) {
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwSubscriptionPort().Return(uint16(9903)).AnyTimes()
	adapter := adapter.NewMockAdapter(gomock.NewController(t))
	jagwSubscriptionService := NewJagwSubscriptionService(config, adapter, make(chan domain.NetworkEvent))

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
	jagwSubscriptionService.cancelFunctions = append(jagwSubscriptionService.cancelFunctions, lsNodesSubscriptionCancel, lsLinksSubscriptionCancel, lsPrefixesSubscriptionCancel, lsSrv6SidsSubscriptionCancel)

	subscriptionClient := jagw.NewMockSubscriptionServiceClient(gomock.NewController(t))
	subscriptionClient.EXPECT().SubscribeToLsNodes(gomock.Any(), gomock.Any()).Return(lsNodesSubscription, nil).AnyTimes()
	subscriptionClient.EXPECT().SubscribeToLsLinks(gomock.Any(), gomock.Any()).Return(lsLinksSubscription, nil).AnyTimes()
	subscriptionClient.EXPECT().SubscribeToLsPrefixes(gomock.Any(), gomock.Any()).Return(lsPrefixesSubscription, nil).AnyTimes()
	subscriptionClient.EXPECT().SubscribeToLsSrv6Sids(gomock.Any(), gomock.Any()).Return(lsSrv6SidsSubscription, nil).AnyTimes()
	err := jagwSubscriptionService.Init()
	assert.NoError(t, err)
	jagwSubscriptionService.subscriptionClient = subscriptionClient
	err = jagwSubscriptionService.Start()
	assert.NoError(t, err)
	jagwSubscriptionService.Stop()
}
