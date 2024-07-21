package calculation

import (
	"context"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/api"
	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func TestNewCalculationSetupProvider(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test NewCalculationSetupProvider",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			assert.NotNil(t, NewCalculationSetupProvider(cacheMock, graphMock))
		})
	}
}

func TestCalculationSetupProvider_getNetworkAddress(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		network string
		wantErr bool
	}{
		{
			name:    "Test Get Network Address",
			ip:      "2001:db8::1",
			network: "2001:db8::",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			provider := NewCalculationSetupProvider(cacheMock, graphMock)
			network := provider.getNetworkAddress(tt.ip)
			assert.Equal(t, tt.network, network.String())
		})
	}
}

func TestCalculationSetupProvider_getNode(t *testing.T) {
	sourceIpv6Address := "2001:db8::1"
	destinationIpv6Address := "2001:db8::2"
	tests := []struct {
		name     string
		nodeType NodeType
		ip       string
		cacheErr bool
		graphErr bool
	}{
		{
			name:     "Test Get Source Node successfully",
			nodeType: Source,
			ip:       sourceIpv6Address,
		},
		{
			name:     "Test Get Destination Node successfully",
			nodeType: Destination,
			ip:       destinationIpv6Address,
		},
		{
			name:     "Test Get Source Node cache error",
			nodeType: Destination,
			ip:       sourceIpv6Address,
			cacheErr: true,
		},
		{
			name:     "Test Get Source Node graph error",
			nodeType: Destination,
			ip:       sourceIpv6Address,
			graphErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			provider := NewCalculationSetupProvider(cacheMock, graphMock)
			intents := []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}
			stream := api.NewMockIntentController_GetIntentPathServer(controller)
			pathRequest, err := domain.NewDomainPathRequest(sourceIpv6Address, destinationIpv6Address, intents, stream, context.Background())
			assert.NoError(t, err)
			if tt.cacheErr {
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress(gomock.Any()).Return("")
				_, err := provider.getNode(pathRequest, tt.nodeType)
				assert.Error(t, err)
				return
			}
			cacheMock.EXPECT().GetRouterIdFromNetworkAddress(gomock.Any()).Return("routerId")
			if tt.graphErr {
				graphMock.EXPECT().GetNode("routerId").Return(nil)
				_, err := provider.getNode(pathRequest, tt.nodeType)
				assert.Error(t, err)
			} else {
				graphMock.EXPECT().GetNode("routerId").Return(graph.NewMockNode(controller))
				node, err := provider.getNode(pathRequest, tt.nodeType)
				assert.NoError(t, err)
				assert.NotNil(t, node)
			}
		})
	}
}

func TestCalculationSetupProvider_getSourceNode(t *testing.T) {
	tests := []struct {
		name     string
		cacheErr bool
		graphErr bool
	}{
		{
			name: "Test Get Source Node successfully",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			provider := NewCalculationSetupProvider(cacheMock, graphMock)
			intents := []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}
			stream := api.NewMockIntentController_GetIntentPathServer(controller)
			pathRequest, err := domain.NewDomainPathRequest("2001:db8::1", "2001:db8::2", intents, stream, context.Background())
			assert.NoError(t, err)
			cacheMock.EXPECT().GetRouterIdFromNetworkAddress(gomock.Any()).Return("routerId")
			graphMock.EXPECT().GetNode("routerId").Return(graph.NewMockNode(controller))
			node, err := provider.getSourceNode(pathRequest)
			assert.NoError(t, err)
			assert.NotNil(t, node)
		})
	}
}

func TestCalculationProvider_getDestinationNode(t *testing.T) {
	tests := []struct {
		name     string
		cacheErr bool
		graphErr bool
	}{
		{
			name: "Test Get Destination Node successfully",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			provider := NewCalculationSetupProvider(cacheMock, graphMock)
			intents := []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}
			stream := api.NewMockIntentController_GetIntentPathServer(controller)
			pathRequest, err := domain.NewDomainPathRequest("2001:db8::1", "2001:db8::2", intents, stream, context.Background())
			assert.NoError(t, err)
			cacheMock.EXPECT().GetRouterIdFromNetworkAddress(gomock.Any()).Return("routerId")
			graphMock.EXPECT().GetNode("routerId").Return(graph.NewMockNode(controller))
			node, err := provider.getDestinationNode(pathRequest)
			assert.NoError(t, err)
			assert.NotNil(t, node)
		})
	}
}

func TestCalculationSetupProvider_getWeightKeyAndCalcMode(t *testing.T) {
	tests := []struct {
		name            string
		intentType      domain.IntentType
		wantWeightKey   helper.WeightKey
		calculationMode CalculationMode
	}{
		{
			name:            "Test Get Weight Key and Calculation Mode",
			intentType:      domain.IntentTypeLowLatency,
			wantWeightKey:   helper.LatencyKey,
			calculationMode: CalculationModeSum,
		},
		{
			name:            "Test Get Weight Key and Calculation Mode",
			intentType:      domain.IntentTypeLowJitter,
			wantWeightKey:   helper.JitterKey,
			calculationMode: CalculationModeSum,
		},
		{
			name:            "Test Get Weight Key and Calculation Mode",
			intentType:      domain.IntentTypeLowPacketLoss,
			wantWeightKey:   helper.PacketLossKey,
			calculationMode: CalculationModeSum,
		},
		{
			name:            "Test Get Weight Key and Calculation Mode",
			intentType:      domain.IntentTypeHighBandwidth,
			wantWeightKey:   helper.AvailableBandwidthKey,
			calculationMode: CalculationModeMax,
		},
		{
			name:            "Test Get Weight Key and Calculation Mode",
			intentType:      domain.IntentTypeLowBandwidth,
			wantWeightKey:   helper.MaximumLinkBandwidth,
			calculationMode: CalculationModeMin,
		},
		{
			name:            "Test Get Weight Key and Calculation Mode",
			intentType:      domain.IntentTypeLowUtilization,
			wantWeightKey:   helper.UtilizedBandwidthKey,
			calculationMode: CalculationModeSum,
		},
		{
			name:            "Test Get Weight Key and Calculation Mode",
			intentType:      domain.IntentTypeSFC,
			wantWeightKey:   helper.IgpMetricKey,
			calculationMode: CalculationModeSum,
		},
		{
			name:            "Test Get Weight Key and Calculation Mode",
			intentType:      domain.IntentTypeFlexAlgo,
			wantWeightKey:   helper.IgpMetricKey,
			calculationMode: CalculationModeSum,
		},
		{
			name:            "Test Get Weight Key and Calculation Mode",
			intentType:      domain.IntentTypeUnspecified,
			wantWeightKey:   helper.UndefinedKey,
			calculationMode: CalculationModeUndefined,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewCalculationSetupProvider(nil, nil)
			weightKey, calcMode := provider.getWeightKeyAndCalcMode(tt.intentType)
			assert.Equal(t, tt.wantWeightKey, weightKey)
			assert.Equal(t, tt.calculationMode, calcMode)
		})
	}
}

func TestCalculationSetupProvider_getWeightKey(t *testing.T) {
	tests := []struct {
		name          string
		intentType    domain.IntentType
		wantWeightKey helper.WeightKey
	}{
		{
			name:          "Test Get Weight Key",
			intentType:    domain.IntentTypeLowLatency,
			wantWeightKey: helper.NormalizedLatencyKey,
		},
		{
			name:          "Test Get Weight Key",
			intentType:    domain.IntentTypeLowJitter,
			wantWeightKey: helper.NormalizedJitterKey,
		},
		{
			name:          "Test Get Weight Key",
			intentType:    domain.IntentTypeLowPacketLoss,
			wantWeightKey: helper.NormalizedPacketLossKey,
		},
		{
			name:          "Test Get Weight Key",
			intentType:    domain.IntentTypeHighBandwidth,
			wantWeightKey: helper.AvailableBandwidthKey,
		},
		{
			name:          "Test Get Weight Key",
			intentType:    domain.IntentTypeUnspecified,
			wantWeightKey: helper.UndefinedKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewCalculationSetupProvider(nil, nil)
			weightKey := provider.getWeightKey(tt.intentType)
			assert.Equal(t, tt.wantWeightKey, weightKey)
		})
	}
}

func TestCalculationSetupProvider_getWeightKeys(t *testing.T) {
	sfcValue, _ := domain.NewStringValue(domain.ValueTypeSFC, proto.String("fw"))
	tests := []struct {
		name       string
		intents    []domain.Intent
		offset     int
		wantKeys   []helper.WeightKey
		wantLength int
	}{
		{
			name: "Test Get Weight Keys low latency, low jitter, low packet loss",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowJitter, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{}),
			},
			offset:     0,
			wantKeys:   []helper.WeightKey{helper.NormalizedLatencyKey, helper.NormalizedJitterKey, helper.NormalizedPacketLossKey},
			wantLength: 3,
		},
		{
			name: "Test Get Weight Keys flex algo, low packet loss",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeFlexAlgo, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{}),
			},
			offset:     1,
			wantKeys:   []helper.WeightKey{helper.NormalizedPacketLossKey},
			wantLength: 1,
		},
		{
			name: "Test Get Weight Keys flex algo, low jitter, low packet loss",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeFlexAlgo, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowJitter, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{}),
			},
			offset:     1,
			wantKeys:   []helper.WeightKey{helper.NormalizedJitterKey, helper.NormalizedPacketLossKey},
			wantLength: 2,
		},
		{
			name: "Test Get Weight Keys sfc, flex algo, low jitter, low packet loss",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeSFC, []domain.Value{sfcValue}),
				domain.NewDomainIntent(domain.IntentTypeFlexAlgo, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowJitter, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{}),
			},
			offset:     2,
			wantKeys:   []helper.WeightKey{helper.NormalizedJitterKey, helper.NormalizedPacketLossKey},
			wantLength: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewCalculationSetupProvider(nil, nil)
			keys := provider.getWeightKeys(tt.intents, tt.offset)
			assert.Equal(t, tt.wantKeys, keys)
			assert.Equal(t, tt.wantLength, len(keys))
		})
	}
}

func TestCalculationSetupProvider_GetWeightKeysandCalculationMode(t *testing.T) {
	sfcValue, _ := domain.NewStringValue(domain.ValueTypeSFC, proto.String("fw"))
	tests := []struct {
		name         string
		intents      []domain.Intent
		wantKeys     []helper.WeightKey
		wantCalcMode CalculationMode
	}{
		{
			name: "Test Get Weight Keys and Calculation Mode low latency",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
			},
			wantKeys:     []helper.WeightKey{helper.LatencyKey},
			wantCalcMode: CalculationModeSum,
		},
		{
			name: "Test Get Weight Keys and Calculation Mode low latency, low jitter",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowJitter, []domain.Value{}),
			},
			wantKeys:     []helper.WeightKey{helper.NormalizedLatencyKey, helper.NormalizedJitterKey},
			wantCalcMode: CalculationModeSum,
		},
		{
			name: "Test Get Weight Keys and Calculation Mode low latency, low jitter, low packet loss",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowJitter, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{}),
			},
			wantKeys:     []helper.WeightKey{helper.NormalizedLatencyKey, helper.NormalizedJitterKey, helper.NormalizedPacketLossKey},
			wantCalcMode: CalculationModeSum,
		},
		{
			name: "Test Get Weight Keys and Calculation Mode flex algo, low latency",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeFlexAlgo, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
			},
			wantKeys:     []helper.WeightKey{helper.LatencyKey},
			wantCalcMode: CalculationModeSum,
		},
		{
			name: "Test Get Weight Keys and Calculation Mode flex algo, low latency, low jitter, low packet loss",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeFlexAlgo, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowJitter, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{}),
			},
			wantKeys:     []helper.WeightKey{helper.NormalizedLatencyKey, helper.NormalizedJitterKey, helper.NormalizedPacketLossKey},
			wantCalcMode: CalculationModeSum,
		},
		{
			name: "Test Get Weight Keys and Calculation Mode sfc, low latency",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeSFC, []domain.Value{sfcValue}),
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
			},
			wantKeys:     []helper.WeightKey{helper.LatencyKey},
			wantCalcMode: CalculationModeSum,
		},
		{
			name: "Test Get Weight Keys and Calculation Mode sfc, low latency, low packet loss",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeSFC, []domain.Value{sfcValue}),
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{}),
			},
			wantKeys:     []helper.WeightKey{helper.NormalizedLatencyKey, helper.NormalizedPacketLossKey},
			wantCalcMode: CalculationModeSum,
		},
		{
			name: "Test Get Weight Keys and Calculation Mode sfc, flex algo, low latency",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeSFC, []domain.Value{sfcValue}),
				domain.NewDomainIntent(domain.IntentTypeFlexAlgo, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
			},
			wantKeys:     []helper.WeightKey{helper.LatencyKey},
			wantCalcMode: CalculationModeSum,
		},
		{
			name: "Test Get Weight Keys and Calculation Mode sfc, flex algo, low latency, low jitter, low packet loss",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeSFC, []domain.Value{sfcValue}),
				domain.NewDomainIntent(domain.IntentTypeFlexAlgo, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowJitter, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{}),
			},
			wantKeys:     []helper.WeightKey{helper.NormalizedLatencyKey, helper.NormalizedJitterKey, helper.NormalizedPacketLossKey},
			wantCalcMode: CalculationModeSum,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			provider := NewCalculationSetupProvider(cacheMock, graphMock)
			keys, calcMode := provider.GetWeightKeysandCalculationMode(tt.intents)
			assert.Equal(t, tt.wantKeys, keys)
			assert.Equal(t, tt.wantCalcMode, calcMode)
		})
	}
}

func TestCalculationSetupProvider_GetMaxConstraints(t *testing.T) {
	numberValue1, _ := domain.NewNumberValue(domain.ValueTypeMaxValue, proto.Int32(10))
	numberValue2, _ := domain.NewNumberValue(domain.ValueTypeMaxValue, proto.Int32(2))
	tests := []struct {
		name          string
		intents       []domain.Intent
		weightKeys    []helper.WeightKey
		wantMaxValues map[helper.WeightKey]float64
	}{
		{
			name: "Test Get Max Constraints single value",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{numberValue1}),
			},
			weightKeys: []helper.WeightKey{helper.NormalizedLatencyKey},
			wantMaxValues: map[helper.WeightKey]float64{
				helper.NormalizedLatencyKey: 10,
			},
		},
		{
			name: "Test Get Max Constraints single value",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{numberValue1}),
				domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{numberValue2}),
			},
			weightKeys: []helper.WeightKey{helper.NormalizedLatencyKey, helper.NormalizedPacketLossKey},
			wantMaxValues: map[helper.WeightKey]float64{
				helper.NormalizedLatencyKey:    10,
				helper.NormalizedPacketLossKey: 0.02,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			provider := NewCalculationSetupProvider(cacheMock, graphMock)
			maxValues := provider.getMaxConstraints(tt.intents, tt.weightKeys)
			assert.Equal(t, tt.wantMaxValues, maxValues)
		})
	}
}

func TestCalculationSetupProvider_getMinConstraints(t *testing.T) {
	numberValue1, _ := domain.NewNumberValue(domain.ValueTypeMinValue, proto.Int32(10))
	tests := []struct {
		name          string
		intents       []domain.Intent
		weightKeys    []helper.WeightKey
		wantMinValues map[helper.WeightKey]float64
	}{
		{
			name: "Test Get Min Constraint",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeHighBandwidth, []domain.Value{numberValue1}),
			},
			wantMinValues: map[helper.WeightKey]float64{
				helper.AvailableBandwidthKey: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			provider := NewCalculationSetupProvider(cacheMock, graphMock)
			minValues := provider.getMinConstraints(tt.intents, []helper.WeightKey{helper.AvailableBandwidthKey})
			assert.Equal(t, tt.wantMinValues, minValues)
		})
	}
}

func TestCalculationSetupProvider_getServiceSids(t *testing.T) {
	fw, _ := domain.NewStringValue(domain.ValueTypeSFC, proto.String("fw"))
	ids, _ := domain.NewStringValue(domain.ValueTypeSFC, proto.String("ids"))
	fwSids := []string{"fc00:0:2f::", "fc00:0:3f::"}
	idsSids := []string{"fc00:0:6f::", "fc00:0:7f::"}
	tests := []struct {
		name     string
		intents  domain.Intent
		wantSids [][]string
		wantErr  bool
	}{
		{
			name: "Test Get Service Sids success",
			intents: domain.NewDomainIntent(domain.IntentTypeSFC, []domain.Value{
				fw,
				ids,
			}),
			wantSids: [][]string{
				fwSids,
				idsSids,
			},
			wantErr: false,
		},
		{
			name: "Test Get Service Sids empty ids sids",
			intents: domain.NewDomainIntent(domain.IntentTypeSFC, []domain.Value{
				fw,
				ids,
			}),
			wantSids: [][]string{
				fwSids,
				{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			provider := NewCalculationSetupProvider(cacheMock, graphMock)
			cacheMock.EXPECT().GetServiceSids("fw").Return(fwSids)
			if tt.wantErr {
				cacheMock.EXPECT().GetServiceSids("ids").Return([]string{})
			} else {
				cacheMock.EXPECT().GetServiceSids("ids").Return(idsSids)
			}
			sids, err := provider.getServiceSids(tt.intents)
			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantSids, sids)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestCalculationSetupProvider_getServiceRouter(t *testing.T) {
	serviceSids := [][]string{
		{"fc00:0:2f::", "fc00:0:3f::"},
		{"fc00:0:6f::", "fc00:0:7f::"},
	}
	tests := []struct {
		name                 string
		routerIds            []string
		wantServiceRouters   [][]string
		wantRouterMap        map[string]string
		wantEmptyRouterIdErr bool
		algorithm            uint32
	}{
		{
			name:      "Test Get Service Router success algo 0",
			routerIds: []string{"router1", "router2", "router3", "router4"},
			wantServiceRouters: [][]string{
				{"router1", "router2"},
				{"router3", "router4"},
			},
			wantRouterMap: map[string]string{
				"router1": "fc00:0:2f::",
				"router2": "fc00:0:3f::",
				"router3": "fc00:0:6f::",
				"router4": "fc00:0:7f::",
			},
			wantEmptyRouterIdErr: false,
			algorithm:            0,
		},
		{
			name:      "Test Get Service Router success algo 128",
			routerIds: []string{"router1", "router2", "router3", "router4"},
			wantServiceRouters: [][]string{
				{"router2"},
				{"router4"},
			},
			wantRouterMap: map[string]string{
				"router2": "fc00:0:3f::",
				"router4": "fc00:0:7f::",
			},
			wantEmptyRouterIdErr: false,
			algorithm:            128,
		},
		{
			name:      "Test Get Service Router empty router id",
			routerIds: []string{"router1", "router2"},
			wantServiceRouters: [][]string{
				{"router1"},
				{},
			},
			wantRouterMap: map[string]string{
				"router1": "fc00:0:2f::",
				"router2": "fc00:0:3f::",
			},
			wantEmptyRouterIdErr: true,
			algorithm:            128,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			provider := NewCalculationSetupProvider(cacheMock, graphMock)
			if tt.wantEmptyRouterIdErr {
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress(serviceSids[0][0]).Return("").AnyTimes()
			} else {
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress(serviceSids[0][0]).Return("router1").AnyTimes()
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress(serviceSids[0][1]).Return("router2").AnyTimes()
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress(serviceSids[1][0]).Return("router3").AnyTimes()
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress(serviceSids[1][1]).Return("router4").AnyTimes()
				cacheMock.EXPECT().GetSrAlgorithmSid(tt.routerIds[0], tt.algorithm).Return(tt.wantRouterMap[tt.routerIds[0]]).AnyTimes()
				cacheMock.EXPECT().GetSrAlgorithmSid(tt.routerIds[1], tt.algorithm).Return(tt.wantRouterMap[tt.routerIds[1]]).AnyTimes()
				cacheMock.EXPECT().GetSrAlgorithmSid(tt.routerIds[2], tt.algorithm).Return(tt.wantRouterMap[tt.routerIds[2]]).AnyTimes()
				cacheMock.EXPECT().GetSrAlgorithmSid(tt.routerIds[3], tt.algorithm).Return(tt.wantRouterMap[tt.routerIds[3]]).AnyTimes()
			}
			serviceRouters, routerServiceMap, err := provider.getServiceRouter(serviceSids, tt.algorithm)
			if !tt.wantEmptyRouterIdErr {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantServiceRouters, serviceRouters)
				assert.Equal(t, tt.wantRouterMap, routerServiceMap)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestCalculationSetupProvider_getServiceChainCombinations(t *testing.T) {
	tests := []struct {
		name             string
		serviceRouters   [][]string
		wantCombinations [][]string
	}{
		{
			name: "Test Get Service Chain Combinations two services with 2 routers each",
			wantCombinations: [][]string{
				{"router1", "router3"},
				{"router1", "router4"},
				{"router2", "router3"},
				{"router2", "router4"},
			},
			serviceRouters: [][]string{
				{"router1", "router2"},
				{"router3", "router4"},
			},
		},
		{
			name: "Test Get Service Chain Combinations two services one service only one router",
			wantCombinations: [][]string{
				{"router1", "router3"},
				{"router1", "router4"},
			},
			serviceRouters: [][]string{
				{"router1"},
				{"router3", "router4"},
			},
		},
		{
			name: "Test Get Service Chain Combinations two routers only",
			wantCombinations: [][]string{
				{"router1"},
				{"router2"},
			},
			serviceRouters: [][]string{
				{"router1", "router2"},
			},
		},
		{
			name: "Test Get Service Chain Combinations 3 services, thirds service has 1 router",
			wantCombinations: [][]string{
				{"router1", "router3", "router5"},
				{"router1", "router4", "router5"},
				{"router2", "router3", "router5"},
				{"router2", "router4", "router5"},
			},
			serviceRouters: [][]string{
				{"router1", "router2"},
				{"router3", "router4"},
				{"router5"},
			},
		},
		{
			name: "Test Get Service Chain Combinations 3 services all have 2 router",
			wantCombinations: [][]string{
				{"router1", "router3", "router5"},
				{"router1", "router3", "router6"},
				{"router1", "router4", "router5"},
				{"router1", "router4", "router6"},
				{"router2", "router3", "router5"},
				{"router2", "router3", "router6"},
				{"router2", "router4", "router5"},
				{"router2", "router4", "router6"},
			},
			serviceRouters: [][]string{
				{"router1", "router2"},
				{"router3", "router4"},
				{"router5", "router6"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			provider := NewCalculationSetupProvider(cacheMock, graphMock)
			combinations := provider.getServiceChainCombinations(tt.serviceRouters)
			assert.Equal(t, tt.wantCombinations, combinations)
		})
	}
}

func TestCalculationSetupProvider_PerformServiceFunctionChainSetup(t *testing.T) {
	fwValue, _ := domain.NewStringValue(domain.ValueTypeSFC, proto.String("fw"))
	idsValue, _ := domain.NewStringValue(domain.ValueTypeSFC, proto.String("ids"))
	sfcIntent := domain.NewDomainIntent(domain.IntentTypeSFC, []domain.Value{fwValue, idsValue})
	serviceSids := [][]string{
		{"fc00:0:2f::", "fc00:0:3f::"},
		{"fc00:0:6f::", "fc00:0:7f::"},
	}

	tests := []struct {
		name              string
		wantServiceSidErr bool
		serviceRouterErr  bool
	}{
		{
			name:              "Test Perform Service Function Chain Setup success",
			wantServiceSidErr: false,
			serviceRouterErr:  false,
		},
		{
			name:              "Test Perform Service Function Chain Setup service sid error",
			wantServiceSidErr: true,
			serviceRouterErr:  false,
		},
		{
			name:              "Test Perform Service Function Chain Setup service router error",
			wantServiceSidErr: false,
			serviceRouterErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			provider := NewCalculationSetupProvider(cacheMock, graphMock)
			cacheMock.EXPECT().GetServiceSids("fw").Return(serviceSids[0]).AnyTimes()
			if tt.wantServiceSidErr {
				cacheMock.EXPECT().GetServiceSids("ids").Return([]string{}).AnyTimes()
			} else {
				cacheMock.EXPECT().GetServiceSids("ids").Return(serviceSids[1]).AnyTimes()
			}
			if tt.serviceRouterErr {
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress(serviceSids[0][0]).Return("").AnyTimes()
			} else {
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress("fc00:0:2f::").Return("router1").AnyTimes()
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress("fc00:0:3f::").Return("router2").AnyTimes()
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress("fc00:0:6f::").Return("router3").AnyTimes()
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress("fc00:0:7f::").Return("router4").AnyTimes()
				cacheMock.EXPECT().GetSrAlgorithmSid("router1", gomock.Any()).Return("fc00:0:2f::").AnyTimes()
				cacheMock.EXPECT().GetSrAlgorithmSid("router2", gomock.Any()).Return("fc00:0:3f::").AnyTimes()
				cacheMock.EXPECT().GetSrAlgorithmSid("router3", gomock.Any()).Return("fc00:0:6f::").AnyTimes()
				cacheMock.EXPECT().GetSrAlgorithmSid("router4", gomock.Any()).Return("fc00:0:7f::").AnyTimes()
			}

			if tt.wantServiceSidErr || tt.serviceRouterErr {
				sfcCalculationOptions, err := provider.PerformServiceFunctionChainSetup(sfcIntent, 0)
				assert.Nil(t, sfcCalculationOptions)
				assert.Error(t, err)
			} else {
				sfcCalculationOptions, err := provider.PerformServiceFunctionChainSetup(sfcIntent, 0)
				assert.NotNil(t, sfcCalculationOptions)
				assert.NoError(t, err)
			}
		})
	}
}

func TestCalculationSetupProvider_PerformSetup(t *testing.T) {
	sourceIpv6Address := "2001:db8:1::1"
	destinationIpv6Address := "2001:db8:2::2"
	intents := []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}
	tests := []struct {
		name               string
		wantErr            bool
		sourceNodeErr      bool
		destinationNodeErr bool
		intents            []domain.Intent
		intentErr          bool
		weightKeyErr       bool
	}{
		{
			name:    "Test Perform Setup success",
			intents: intents,
		},
		{
			name:          "Test Perform Setup source node error",
			sourceNodeErr: true,
			intents:       intents,
			wantErr:       true,
		},
		{
			name:               "Test Perform Setup destination node error",
			destinationNodeErr: true,
			intents:            intents,
			wantErr:            true,
		},
		{
			name:         "Test Perform Setup weight key error",
			weightKeyErr: true,
			intents:      []domain.Intent{domain.NewDomainIntent(domain.IntentTypeUnspecified, []domain.Value{})},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			provider := NewCalculationSetupProvider(cacheMock, graphMock)
			stream := api.NewMockIntentController_GetIntentPathServer(controller)
			pathRequest, err := domain.NewDomainPathRequest(sourceIpv6Address, destinationIpv6Address, tt.intents, stream, context.Background())
			assert.NoError(t, err)
			if tt.sourceNodeErr {
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress("2001:db8:1::").Return("").AnyTimes()
			} else {
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress("2001:db8:1::").Return("routerId").AnyTimes()
				graphMock.EXPECT().GetNode("routerId").Return(graph.NewMockNode(controller)).AnyTimes()
			}
			if tt.destinationNodeErr {
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress("2001:db8:2::").Return("").AnyTimes()
			} else {
				cacheMock.EXPECT().GetRouterIdFromNetworkAddress("2001:db8:2::").Return("routerId").AnyTimes()
				graphMock.EXPECT().GetNode("routerId").Return(graph.NewMockNode(controller)).AnyTimes()
			}
			_, err = provider.PerformSetup(pathRequest)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
