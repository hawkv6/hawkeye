package jagw

import (
	"fmt"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/adapter"
	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/hawkv6/hawkeye/pkg/processor"
	"github.com/jalapeno-api-gateway/jagw-go/jagw"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func setUpJagwSid(key string, igpRouterId string, sid string, sidType string, algorithm uint32) *jagw.LsSrv6Sid {
	srv6Sid := &jagw.LsSrv6Sid{}
	if key != "" {
		srv6Sid.Key = proto.String(key)
	}
	if igpRouterId != "" {
		srv6Sid.IgpRouterId = proto.String(igpRouterId)
	}
	if sid != "" {
		srv6Sid.Srv6Sid = proto.String(sid)
	}
	if sidType != "" {
		srv6Sid.Srv6EndpointBehavior = &jagw.Srv6EndpointBehavior{Algorithm: proto.Uint32(algorithm)}
	}
	return srv6Sid
}

func setUpJagwPrefix(key string, igpRouterId string, prefix string, prefixLength int32) *jagw.LsPrefix {
	lsPrefix := &jagw.LsPrefix{}
	if key != "" {
		lsPrefix.Key = proto.String(key)
	}
	if igpRouterId != "" {
		lsPrefix.IgpRouterId = proto.String(igpRouterId)
	}
	if prefix != "" {
		lsPrefix.Prefix = proto.String(prefix)
	}
	if prefixLength != 0 {
		lsPrefix.PrefixLen = proto.Int32(prefixLength)
	}
	return lsPrefix
}

func setUpJagwNode(key string, igpRouterId string, name string, srAlgorithm []uint32) *jagw.LsNode {
	lsNode := &jagw.LsNode{}
	if key != "" {
		lsNode.Key = proto.String(key)
	}
	if igpRouterId != "" {
		lsNode.IgpRouterId = proto.String(igpRouterId)
	}
	if name != "" {
		lsNode.Name = proto.String(name)
	}
	if srAlgorithm != nil {
		lsNode.SrAlgorithm = srAlgorithm
	}
	return lsNode
}

func setUpJagwLink(key string, igpRouterId string, remoteIgpRouterId string, igpMetric uint32, unidirLinkDelay uint32, unidirDelayVariation uint32, maxLinkBwKbps uint64, unidirAvailableBw uint32, unidirBwUtilization uint32, unidirPacketLossPercentage, normalizedUnidirLinkDelay, normalizedUnidirDelayVariation, normalizedUnidirPacketLoss float64) *jagw.LsLink {
	lsLink := &jagw.LsLink{}
	if key != "" {
		lsLink.Key = proto.String(key)
	}
	if igpRouterId != "" {
		lsLink.IgpRouterId = proto.String(igpRouterId)
	}
	if remoteIgpRouterId != "" {
		lsLink.RemoteIgpRouterId = proto.String(remoteIgpRouterId)
	}
	lsLink.IgpMetric = proto.Uint32(igpMetric)
	lsLink.UnidirLinkDelay = proto.Uint32(unidirLinkDelay)
	lsLink.UnidirDelayVariation = proto.Uint32(unidirDelayVariation)
	lsLink.MaxLinkBwKbps = proto.Uint64(maxLinkBwKbps)
	lsLink.UnidirAvailableBw = proto.Uint32(unidirAvailableBw)
	lsLink.UnidirBwUtilization = proto.Uint32(unidirBwUtilization)
	lsLink.UnidirPacketLossPercentage = proto.Float64(unidirPacketLossPercentage)
	lsLink.NormalizedUnidirLinkDelay = proto.Float64(normalizedUnidirLinkDelay)
	lsLink.NormalizedUnidirDelayVariation = proto.Float64(normalizedUnidirDelayVariation)
	lsLink.NormalizedUnidirPacketLoss = proto.Float64(normalizedUnidirPacketLoss)
	return lsLink
}

func getLsSrv6SidResponse() *jagw.LsSrv6SidResponse {
	return &jagw.LsSrv6SidResponse{
		LsSrv6Sids: []*jagw.LsSrv6Sid{
			setUpJagwSid("0_0000.0000.000b_fc00:0:b:0:1::", "0000.0000.000b", "fc00:0:b:0:1::", "End", 0),
			setUpJagwSid("0_0000.0000.0006_fc00:0:6:0:1::", "0000.0000.0006", "fc00:0:6:0:1::", "End", 0),
		},
	}
}

func getLsPrefixesResponse() *jagw.LsPrefixResponse {
	return &jagw.LsPrefixResponse{
		LsPrefixes: []*jagw.LsPrefix{
			setUpJagwPrefix("2_0_2_0_0_2001:db8:b::_64_0000.0000.000b", "0000.0000.000b", "2001:db8:b::", 64),
			setUpJagwPrefix("2_0_2_0_0_2001:db8:6::_64_0000.0000.0006", "0000.0000.0006", "2001:db8:6::", 64),
		},
	}
}

func getLsNodesResponse() *jagw.LsNodeResponse {
	return &jagw.LsNodeResponse{
		LsNodes: []*jagw.LsNode{
			setUpJagwNode("0_0000.0000.000b", "0000.0000.000b", "SITE-B", []uint32{0}),
			setUpJagwNode("0_0000.0000.0006", "0000.0000.0006", "XR-6", []uint32{0}),
		},
	}
}

func getLsLinksResponse() *jagw.LsLinkResponse {
	return &jagw.LsLinkResponse{
		LsLinks: []*jagw.LsLink{
			setUpJagwLink("0_0000.0000.000b_0000.0000.0006", "0000.0000.000b", "0000.0000.0006", 10, 10, 10, 1000000, 1000000, 100, 0.1, 1, 1, 1),
			setUpJagwLink("0_0000.0000.0006_0000.0000.000b", "0000.0000.0006", "0000.0000.000b", 10, 10, 10, 1000000, 1000000, 100, 0.1, 1, 1, 1),
		},
	}
}
func TestNewJagwRequestService(t *testing.T) {
	tests := []struct {
		name                  string
		jagwServiceAddress    string
		requestPort           uint16
		wantJagwRequestSocket string
	}{
		{
			name:                  "TestNewJagwRequestService",
			jagwServiceAddress:    "localhost",
			requestPort:           9902,
			wantJagwRequestSocket: "localhost:9902",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetJagwServiceAddress().Return(tt.jagwServiceAddress).Times(1)
			config.EXPECT().GetJagwRequestPort().Return(tt.requestPort).Times(1)
			adapter := adapter.NewDomainAdapter()
			processor := processor.NewMockProcessor(gomock.NewController(t))
			jagwRequestService := NewJagwRequestService(config, adapter, processor)
			assert.NotNil(t, jagwRequestService)
			assert.Equal(t, tt.wantJagwRequestSocket, jagwRequestService.jagwRequestSocket)
		})
	}
}

func TestJagwRequestService_Init(t *testing.T) {
	tests := []struct {
		name               string
		jagwServiceAddress string
		requestPort        uint16
		wantErr            bool
	}{
		{
			name:               "TestJagwRequestService_Init success",
			jagwServiceAddress: "localhost",
			requestPort:        9902,
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := adapter.NewDomainAdapter()
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetJagwServiceAddress().Return(tt.jagwServiceAddress).Times(1)
			config.EXPECT().GetJagwRequestPort().Return(tt.requestPort).Times(1)
			processor := processor.NewMockProcessor(gomock.NewController(t))
			jagwRequestService := NewJagwRequestService(config, adapter, processor)
			err := jagwRequestService.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("JagwRequestService.Init() '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
		})
	}
}

func TestJagwRequestService_Start(t *testing.T) {
	srv6Response := getLsSrv6SidResponse()
	lsPrefixesResponse := getLsPrefixesResponse()
	lsNodesResponse := getLsNodesResponse()
	lsLinksResponse := getLsLinksResponse()

	tests := []struct {
		name               string
		jagwServiceAddress string
		requestPort        uint16
		subscriptionPort   uint16
		grpcPort           uint16
		wantErr            bool
		setup              func(*JagwRequestService)
	}{
		{
			name:               "TestJagwRequestService_Start success",
			jagwServiceAddress: "localhost",
			requestPort:        9902,
			subscriptionPort:   9903,
			grpcPort:           10000,
			wantErr:            false,
			setup: func(jagwRequestService *JagwRequestService) {
				requestClient := jagw.NewMockRequestServiceClient(gomock.NewController(t))
				requestClient.EXPECT().GetLsSrv6Sids(gomock.Any(), gomock.Any()).Return(srv6Response, nil).Times(1)
				requestClient.EXPECT().GetLsPrefixes(gomock.Any(), gomock.Any()).Return(lsPrefixesResponse, nil).Times(1)
				requestClient.EXPECT().GetLsNodes(gomock.Any(), gomock.Any()).Return(lsNodesResponse, nil).Times(1)
				requestClient.EXPECT().GetLsLinks(gomock.Any(), gomock.Any()).Return(lsLinksResponse, nil).Times(1)
				jagwRequestService.requestClient = requestClient
			},
		},
		{
			name:               "TestJagwRequestService_Start lsSrv6Sids error",
			jagwServiceAddress: "localhost",
			requestPort:        9902,
			subscriptionPort:   9903,
			grpcPort:           10000,
			wantErr:            true,
			setup: func(jagwRequestService *JagwRequestService) {
				requestClient := jagw.NewMockRequestServiceClient(gomock.NewController(t))
				requestClient.EXPECT().GetLsSrv6Sids(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error to get LsSrv6Sids")).Times(1)
				jagwRequestService.requestClient = requestClient
			},
		},
		{
			name:               "TestJagwRequestService_Start lsPrefix error",
			jagwServiceAddress: "localhost",
			requestPort:        9902,
			subscriptionPort:   9903,
			grpcPort:           10000,
			wantErr:            true,
			setup: func(jagwRequestService *JagwRequestService) {
				requestClient := jagw.NewMockRequestServiceClient(gomock.NewController(t))
				requestClient.EXPECT().GetLsSrv6Sids(gomock.Any(), gomock.Any()).Return(srv6Response, nil).Times(1)
				requestClient.EXPECT().GetLsPrefixes(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error to get LsPrefixes")).Times(1)
				jagwRequestService.requestClient = requestClient
			},
		},
		{
			name:               "TestJagwRequestService_Start lsNodes error",
			jagwServiceAddress: "localhost",
			requestPort:        9902,
			subscriptionPort:   9903,
			grpcPort:           10000,
			wantErr:            true,
			setup: func(jagwRequestService *JagwRequestService) {
				requestClient := jagw.NewMockRequestServiceClient(gomock.NewController(t))
				requestClient.EXPECT().GetLsSrv6Sids(gomock.Any(), gomock.Any()).Return(srv6Response, nil).Times(1)
				requestClient.EXPECT().GetLsPrefixes(gomock.Any(), gomock.Any()).Return(lsPrefixesResponse, nil).Times(1)
				requestClient.EXPECT().GetLsNodes(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error to get LsNodes")).Times(1)
				jagwRequestService.requestClient = requestClient
			},
		},
		{
			name:               "TestJagwRequestService_Start lsLinks error",
			jagwServiceAddress: "localhost",
			requestPort:        9902,
			subscriptionPort:   9903,
			grpcPort:           10000,
			wantErr:            true,
			setup: func(jagwRequestService *JagwRequestService) {
				requestClient := jagw.NewMockRequestServiceClient(gomock.NewController(t))
				requestClient.EXPECT().GetLsSrv6Sids(gomock.Any(), gomock.Any()).Return(srv6Response, nil).Times(1)
				requestClient.EXPECT().GetLsPrefixes(gomock.Any(), gomock.Any()).Return(lsPrefixesResponse, nil).Times(1)
				requestClient.EXPECT().GetLsNodes(gomock.Any(), gomock.Any()).Return(lsNodesResponse, nil).Times(1)
				requestClient.EXPECT().GetLsLinks(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error to get LsLinks")).Times(1)
				jagwRequestService.requestClient = requestClient
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetJagwServiceAddress().Return("").Times(1)
			config.EXPECT().GetJagwRequestPort().Return(tt.requestPort).Times(1)
			adapter := adapter.NewDomainAdapter()
			processor := processor.NewMockProcessor(gomock.NewController(t))
			processor.EXPECT().ProcessSids(gomock.Any()).Return().AnyTimes()
			processor.EXPECT().ProcessPrefixes(gomock.Any()).Return().AnyTimes()
			processor.EXPECT().ProcessNodes(gomock.Any()).Return().AnyTimes()
			processor.EXPECT().ProcessLinks(gomock.Any()).Return(nil).AnyTimes()

			jagwRequestService := NewJagwRequestService(config, adapter, processor)
			tt.setup(jagwRequestService)
			err := jagwRequestService.Start()
			if (err != nil) != tt.wantErr {
				t.Errorf("JagwRequestService.Start() '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
		})
	}
}

func TestRequestService_convertLsNodes(t *testing.T) {
	lsNodesResponse := getLsNodesResponse()
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwRequestPort().Return(uint16(9902)).AnyTimes()
	adapter := adapter.NewDomainAdapter()
	processor := processor.NewMockProcessor(gomock.NewController(t))
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestRequestService_convertLsNodes success",
			wantErr: false,
		},
		{
			name:    "TestRequestService_convertLsNodes error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jagwRequestService := NewJagwRequestService(config, adapter, processor)
			if tt.wantErr {
				lsNodesResponse.LsNodes[0].Key = nil
			}
			nodes, err := jagwRequestService.convertLsNodes(lsNodesResponse.LsNodes)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestService.convertLsNodes() '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, len(lsNodesResponse.LsNodes), len(nodes))
			}
		})
	}
}

func TestRequestService_getLsNodes(t *testing.T) {
	lsNodesResponse := getLsNodesResponse()
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwRequestPort().Return(uint16(9902)).AnyTimes()
	adapter := adapter.NewDomainAdapter()
	processor := processor.NewMockProcessor(gomock.NewController(t))
	processor.EXPECT().ProcessNodes(gomock.Any()).Return().AnyTimes()
	requestClient := jagw.NewMockRequestServiceClient(gomock.NewController(t))
	jagwRequestService := NewJagwRequestService(config, adapter, processor)
	jagwRequestService.requestClient = requestClient
	tests := []struct {
		name           string
		wantRequestErr bool
		wantConvertErr bool
	}{
		{
			name:           "TestRequestService_getLsNodes success",
			wantRequestErr: false,
			wantConvertErr: false,
		},
		{
			name:           "TestRequestService_getLsNodes request error",
			wantRequestErr: true,
			wantConvertErr: false,
		},
		{
			name:           "TestRequestService_getLsNodes convert error",
			wantRequestErr: false,
			wantConvertErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantRequestErr {
				requestClient.EXPECT().GetLsNodes(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error to get LsNodes")).Times(1)
			} else {
				requestClient.EXPECT().GetLsNodes(gomock.Any(), gomock.Any()).Return(lsNodesResponse, nil).Times(1)
			}
			if tt.wantConvertErr {
				lsNodesResponse.LsNodes[0].Key = nil
			}
			err := jagwRequestService.getLsNodes()
			if (err != nil) != (tt.wantRequestErr || tt.wantConvertErr) {
				t.Errorf("RequestService.getLsNodes() '%s' error = %v, wantRequestErr %v, wantConvertError %v", tt.name, err, tt.wantRequestErr, tt.wantConvertErr)
				return
			}
		})
	}
}

func TestRequestService_convertLsLinks(t *testing.T) {
	lsLinksResponse := getLsLinksResponse()
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwRequestPort().Return(uint16(9902)).AnyTimes()
	adapter := adapter.NewDomainAdapter()
	processor := processor.NewMockProcessor(gomock.NewController(t))
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestRequestService_convertLsLinks success",
			wantErr: false,
		},
		{
			name:    "TestRequestService_convertLsLinks error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jagwRequestService := NewJagwRequestService(config, adapter, processor)
			if tt.wantErr {
				lsLinksResponse.LsLinks[0].Key = nil
			}
			links, err := jagwRequestService.convertLsLinks(lsLinksResponse.LsLinks)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestService.convertLsLinks() '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, len(lsLinksResponse.LsLinks), len(links))
			}
		})
	}
}

func TestRequestService_getLsLinks(t *testing.T) {
	tests := []struct {
		name           string
		wantRequestErr bool
		wantConvertErr bool
		wantProcessErr bool
	}{
		{
			name:           "TestRequestService_getLsLinks success",
			wantRequestErr: false,
			wantConvertErr: false,
			wantProcessErr: false,
		},
		{
			name:           "TestRequestService_getLsLinks request error",
			wantRequestErr: true,
			wantConvertErr: false,
			wantProcessErr: false,
		},
		{
			name:           "TestRequestService_getLsLinks convert error",
			wantRequestErr: false,
			wantConvertErr: true,
			wantProcessErr: false,
		},
		{
			name:           "TestRequestService_getLsLinks process error",
			wantRequestErr: false,
			wantConvertErr: false,
			wantProcessErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lsLinksResponse := getLsLinksResponse()
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
			config.EXPECT().GetJagwRequestPort().Return(uint16(9902)).AnyTimes()
			adapter := adapter.NewDomainAdapter()
			processor := processor.NewMockProcessor(gomock.NewController(t))
			requestClient := jagw.NewMockRequestServiceClient(gomock.NewController(t))
			jagwRequestService := NewJagwRequestService(config, adapter, processor)
			jagwRequestService.requestClient = requestClient
			if tt.wantRequestErr {
				requestClient.EXPECT().GetLsLinks(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error to get LsLinks")).Times(1)
			} else {
				requestClient.EXPECT().GetLsLinks(gomock.Any(), gomock.Any()).Return(lsLinksResponse, nil).Times(1)
			}
			if tt.wantConvertErr {
				lsLinksResponse.LsLinks[0].Key = nil
			}
			if tt.wantProcessErr {
				processor.EXPECT().ProcessLinks(gomock.Any()).Return(fmt.Errorf("Error to process links")).AnyTimes()
			} else {
				processor.EXPECT().ProcessLinks(gomock.Any()).Return(nil).AnyTimes()
			}
			err := jagwRequestService.getLsLinks()
			if (err != nil) != (tt.wantRequestErr || tt.wantConvertErr || tt.wantProcessErr) {
				t.Errorf("RequestService.getLsLinks() '%s' error = %v, wantRequestErr %v, wantConvertError %v, wantProcessError %v", tt.name, err, tt.wantRequestErr, tt.wantConvertErr, tt.wantProcessErr)
				return
			}
		})
	}
}

func TestRequestService_shouldSkipPrefix(t *testing.T) {
	tests := []struct {
		name       string
		hasLocator bool
		isLoopback bool
	}{
		{
			name:       "TestRequestService_shouldSkipPrefix no skip",
			hasLocator: false,
			isLoopback: false,
		},
		{
			name:       "TestRequestService_shouldSkipPrefix skip, has locator",
			hasLocator: true,
			isLoopback: false,
		},
		{
			name:       "TestRequestService_shouldSkipPrefix skip, is loopback",
			hasLocator: false,
			isLoopback: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
			config.EXPECT().GetJagwRequestPort().Return(uint16(9902)).AnyTimes()
			adapter := adapter.NewDomainAdapter()
			processor := processor.NewMockProcessor(gomock.NewController(t))
			LsPrefixResponse := getLsPrefixesResponse()
			jagwRequestService := NewJagwRequestService(config, adapter, processor)
			for _, lsPrefix := range LsPrefixResponse.LsPrefixes {
				if tt.hasLocator {
					lsPrefix.Srv6Locator = &jagw.Srv6LocatorTlv{}
				}
				if tt.isLoopback {
					lsPrefix.PrefixLen = proto.Int32(128)
				}
				skip := jagwRequestService.shouldSkipPrefix(lsPrefix)
				assert.Equal(t, tt.hasLocator || tt.isLoopback, skip)
			}
		})
	}
}

func TestRequestService_convertLsPrefix(t *testing.T) {
	lsPrefixesResponse := getLsPrefixesResponse()
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwRequestPort().Return(uint16(9902)).AnyTimes()
	adapter := adapter.NewDomainAdapter()
	processor := processor.NewMockProcessor(gomock.NewController(t))
	tests := []struct {
		name             string
		wantErr          bool
		shouldSkipPrefix bool
	}{
		{
			name:             "TestRequestService_convertLsPrefix success",
			wantErr:          false,
			shouldSkipPrefix: false,
		},
		{
			name:             "TestRequestService_convertLsPrefix error",
			wantErr:          true,
			shouldSkipPrefix: false,
		},
		{
			name:             "TestRequestService_convertLsPrefix skip prefix",
			wantErr:          false,
			shouldSkipPrefix: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jagwRequestService := NewJagwRequestService(config, adapter, processor)
			if tt.wantErr {
				lsPrefixesResponse.LsPrefixes[0].Key = nil
			}
			if tt.shouldSkipPrefix {
				lsPrefixesResponse.LsPrefixes[0].Srv6Locator = &jagw.Srv6LocatorTlv{}
			}
			prefixes, err := jagwRequestService.convertLsPrefix(lsPrefixesResponse.LsPrefixes)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestService.convertLsPrefix() '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if tt.shouldSkipPrefix {
				assert.Equal(t, len(lsPrefixesResponse.LsPrefixes)-1, len(prefixes))
				return
			}
			if !tt.wantErr {
				assert.Equal(t, len(lsPrefixesResponse.LsPrefixes), len(prefixes))
			}
		})
	}
}

func TestRequestService_getLsPrefixes(t *testing.T) {
	lsPrefixesResponse := getLsPrefixesResponse()
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwRequestPort().Return(uint16(9902)).AnyTimes()
	adapter := adapter.NewDomainAdapter()
	processor := processor.NewMockProcessor(gomock.NewController(t))
	processor.EXPECT().ProcessPrefixes(gomock.Any()).Return().AnyTimes()
	requestClient := jagw.NewMockRequestServiceClient(gomock.NewController(t))
	jagwRequestService := NewJagwRequestService(config, adapter, processor)
	jagwRequestService.requestClient = requestClient
	tests := []struct {
		name           string
		wantRequestErr bool
		wantConvertErr bool
	}{
		{
			name:           "TestRequestService_getLsPrefixes success",
			wantRequestErr: false,
			wantConvertErr: false,
		},
		{
			name:           "TestRequestService_getLsPrefixes request error",
			wantRequestErr: true,
			wantConvertErr: false,
		},
		{
			name:           "TestRequestService_getLsPrefixes convert error",
			wantRequestErr: false,
			wantConvertErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantRequestErr {
				requestClient.EXPECT().GetLsPrefixes(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error to get LsPrefixes")).Times(1)
			} else {
				requestClient.EXPECT().GetLsPrefixes(gomock.Any(), gomock.Any()).Return(lsPrefixesResponse, nil).Times(1)
			}
			if tt.wantConvertErr {
				lsPrefixesResponse.LsPrefixes[0].Key = nil
			}
			err := jagwRequestService.getLsPrefixes()
			if (err != nil) != (tt.wantRequestErr || tt.wantConvertErr) {
				t.Errorf("RequestService.getLsPrefixes() '%s' error = %v, wantRequestErr %v, wantConvertError %v", tt.name, err, tt.wantRequestErr, tt.wantConvertErr)
				return
			}
		})
	}
}

func TestRequestService_convertLsSrv6Sid(t *testing.T) {
	srv6Response := getLsSrv6SidResponse()
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwRequestPort().Return(uint16(9902)).AnyTimes()
	adapter := adapter.NewDomainAdapter()
	processor := processor.NewMockProcessor(gomock.NewController(t))
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestRequestService_convertLsSrv6Sid success",
			wantErr: false,
		},
		{
			name:    "TestRequestService_convertLsSrv6Sid error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jagwRequestService := NewJagwRequestService(config, adapter, processor)
			if tt.wantErr {
				srv6Response.LsSrv6Sids[0].Key = nil
			}
			sids, err := jagwRequestService.convertLsSrv6Sids(srv6Response.LsSrv6Sids)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestService.convertLsSrv6Sid() '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, len(srv6Response.LsSrv6Sids), len(sids))
			}
		})
	}
}

func TestRequestService_getLsSrv6Sids(t *testing.T) {
	srv6Response := getLsSrv6SidResponse()
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwRequestPort().Return(uint16(9902)).AnyTimes()
	adapter := adapter.NewDomainAdapter()
	processor := processor.NewMockProcessor(gomock.NewController(t))
	processor.EXPECT().ProcessSids(gomock.Any()).Return().AnyTimes()
	requestClient := jagw.NewMockRequestServiceClient(gomock.NewController(t))
	jagwRequestService := NewJagwRequestService(config, adapter, processor)
	jagwRequestService.requestClient = requestClient
	tests := []struct {
		name           string
		wantRequestErr bool
		wantConvertErr bool
	}{
		{
			name:           "TestRequestService_getLsSrv6Sids success",
			wantRequestErr: false,
			wantConvertErr: false,
		},
		{
			name:           "TestRequestService_getLsSrv6Sids request error",
			wantRequestErr: true,
			wantConvertErr: false,
		},
		{
			name:           "TestRequestService_getLsSrv6Sids convert error",
			wantRequestErr: false,
			wantConvertErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantRequestErr {
				requestClient.EXPECT().GetLsSrv6Sids(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error to get LsSrv6Sids")).Times(1)
			} else {
				requestClient.EXPECT().GetLsSrv6Sids(gomock.Any(), gomock.Any()).Return(srv6Response, nil).Times(1)
			}
			if tt.wantConvertErr {
				srv6Response.LsSrv6Sids[0].Key = nil
			}
			err := jagwRequestService.getSrv6Sids()
			if (err != nil) != (tt.wantRequestErr || tt.wantConvertErr) {
				t.Errorf("RequestService.getLsSrv6Sids() '%s' error = %v, want %v", tt.name, err, tt.wantRequestErr)
				return
			}
		})
	}
}

func TestRequestService_Stop(t *testing.T) {
	config := config.NewMockConfig(gomock.NewController(t))
	config.EXPECT().GetJagwServiceAddress().Return("localhost").AnyTimes()
	config.EXPECT().GetJagwRequestPort().Return(uint16(9902)).AnyTimes()
	adapter := adapter.NewDomainAdapter()
	processor := processor.NewMockProcessor(gomock.NewController(t))
	jagwRequestService := NewJagwRequestService(config, adapter, processor)
	err := jagwRequestService.Init()
	assert.NoError(t, err)
	jagwRequestService.Stop()
}
