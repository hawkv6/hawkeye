package calculation

import (
	"context"
	"fmt"
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

func TestNewCalculationManager(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewCalculationManager",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			calculationSetup := NewMockCalculationSetup(controller)
			calculationTransformer := NewMockCalculationTransformer(controller)
			calculationUpdater := NewMockCalculationUpdater(controller)
			got := NewCalculationManager(cacheMock, graphMock, calculationSetup, calculationTransformer, calculationUpdater)
			if got == nil {
				t.Errorf("NewCalculationManager() = %v, want %v", got, nil)
			}
		})
	}
}

func TestCalculationManager_lockElements(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestCalculationManager_lockElements",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			calculationSetup := NewMockCalculationSetup(controller)
			calculationTransformer := NewMockCalculationTransformer(controller)
			calculationUpdater := NewMockCalculationUpdater(controller)
			manager := NewCalculationManager(cacheMock, graphMock, calculationSetup, calculationTransformer, calculationUpdater)
			cacheMock.EXPECT().Lock().Return()
			graphMock.EXPECT().Lock().Return()
			manager.lockElements()
		})
	}
}

func TestCalculationManager_unlockElements(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestCalculationManager_unlockElements",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			calculationSetup := NewMockCalculationSetup(controller)
			calculationTransformer := NewMockCalculationTransformer(controller)
			calculationUpdater := NewMockCalculationUpdater(controller)
			manager := NewCalculationManager(cacheMock, graphMock, calculationSetup, calculationTransformer, calculationUpdater)
			cacheMock.EXPECT().Unlock().Return()
			graphMock.EXPECT().Unlock().Return()
			manager.unlockElements()
		})
	}
}

func TestCalculationManager_getGraphAndAlgorithm(t *testing.T) {
	tests := []struct {
		name      string
		algorithm int32
	}{
		{
			name:      "TestCalculationManager_getGraphAndAlgorithm without flex algo",
			algorithm: 0,
		},
		{
			name:      "TestCalculationManager_getGraphAndAlgorithm with flex algo",
			algorithm: 128,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			calculationSetup := NewMockCalculationSetup(controller)
			calculationTransformer := NewMockCalculationTransformer(controller)
			calculationUpdater := NewMockCalculationUpdater(controller)
			manager := NewCalculationManager(cacheMock, graphMock, calculationSetup, calculationTransformer, calculationUpdater)
			var intent domain.Intent
			if tt.algorithm > 1 {
				flexAlgoValue, err := domain.NewNumberValue(domain.ValueTypeFlexAlgoNr, proto.Int32(tt.algorithm))
				assert.NoError(t, err)
				intent = domain.NewDomainIntent(domain.IntentTypeFlexAlgo, []domain.Value{flexAlgoValue})
				graphMock.EXPECT().GetSubGraph(gomock.Any()).Return(graphMock)
			} else {
				intent = domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})
			}
			got, gotAlgorithm := manager.getGraphAndAlgorithm(graphMock, intent)
			assert.Equal(t, graphMock, got)
			assert.Equal(t, uint32(tt.algorithm), gotAlgorithm)
		})
	}
}

func TestCalculationManager_setUpCalculation(t *testing.T) {
	ipv6SourceAddress := "2001:db8::1"
	ipv6DestinationAddress := "2001:db8::2"
	sfcValue, _ := domain.NewStringValue(domain.ValueTypeSFC, proto.String("sfc"))
	flexAlgoValue, _ := domain.NewNumberValue(domain.ValueTypeFlexAlgoNr, proto.Int32(128))
	tests := []struct {
		name             string
		wantErr          bool
		algorithm        int32
		firstIntentType  domain.IntentType
		firstIntentValue domain.Value
	}{
		{
			name:            "TestCalculationManager_setUpCalculation perform setup error",
			wantErr:         true,
			algorithm:       0,
			firstIntentType: domain.IntentTypeLowLatency,
		},
		{
			name:             "TestCalculationManager_setUpCalculation SFC",
			wantErr:          false,
			algorithm:        0,
			firstIntentType:  domain.IntentTypeSFC,
			firstIntentValue: sfcValue,
		},
		{
			name:            "TestCalculationManager_setUpCalculation ShortestPath",
			wantErr:         false,
			algorithm:       0,
			firstIntentType: domain.IntentTypeLowLatency,
		},
		{
			name:             "TestCalculationManager_setUpCalculation FlexAlgo",
			wantErr:          false,
			algorithm:        128,
			firstIntentType:  domain.IntentTypeFlexAlgo,
			firstIntentValue: flexAlgoValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			calculationSetup := NewMockCalculationSetup(controller)
			calculationTransformer := NewMockCalculationTransformer(controller)
			calculationUpdater := NewMockCalculationUpdater(controller)
			manager := NewCalculationManager(cacheMock, graphMock, calculationSetup, calculationTransformer, calculationUpdater)
			stream := api.NewMockIntentController_GetIntentPathServer(controller)
			ctx := context.Background()
			var intents []domain.Intent
			if tt.firstIntentValue != nil {
				intents = []domain.Intent{domain.NewDomainIntent(tt.firstIntentType, []domain.Value{tt.firstIntentValue})}
			} else {
				intents = []domain.Intent{domain.NewDomainIntent(tt.firstIntentType, []domain.Value{})}
			}
			pathRequest, err := domain.NewDomainPathRequest(ipv6SourceAddress, ipv6DestinationAddress, intents, stream, ctx)
			assert.NoError(t, err)
			if tt.wantErr {
				calculationSetup.EXPECT().PerformSetup(pathRequest).Return(nil, fmt.Errorf("set up failed due to something"))
				err := manager.setUpCalculation(pathRequest)
				assert.Error(t, err)
			} else {
				calculationOptions := &CalculationOptions{
					graph: graphMock,
				}
				calculationSetup.EXPECT().PerformSetup(pathRequest).Return(calculationOptions, nil)
				if tt.firstIntentType == domain.IntentTypeSFC {
					SfcCalculationOptions := &SfcCalculationOptions{}
					calculationSetup.EXPECT().PerformServiceFunctionChainSetup(gomock.Any()).Return(SfcCalculationOptions, nil)
				}
				if tt.firstIntentType == domain.IntentTypeFlexAlgo {
					graphMock.EXPECT().GetSubGraph(gomock.Any()).Return(graphMock)
				}
				err := manager.setUpCalculation(pathRequest)
				assert.NoError(t, err)
				assert.Equal(t, uint32(tt.algorithm), manager.algorithm)
			}
		})
	}
}

func TestCalculationManager_getFirstNonSfcIntent(t *testing.T) {
	tests := []struct {
		name    string
		intents []domain.Intent
		want    domain.Intent
	}{
		{
			name: "TestCalculationManager_getFirstNonSfcIntent",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeSFC, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
			},
			want: domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
		},
		{
			name: "TestCalculationManager_getFirstNonSfcIntent with only LowLatency",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
			},
			want: domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
		},
		{
			name: "TestCalculationManager_getFirstNonSfcIntent with several non-SFC intents",
			intents: []domain.Intent{
				domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{}),
				domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{}),
			},
			want: domain.NewDomainIntent(domain.IntentTypeLowPacketLoss, []domain.Value{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			calculationSetup := NewMockCalculationSetup(controller)
			calculationTransformer := NewMockCalculationTransformer(controller)
			calculationUpdater := NewMockCalculationUpdater(controller)
			manager := NewCalculationManager(cacheMock, graphMock, calculationSetup, calculationTransformer, calculationUpdater)
			got := manager.getFirstNonSfcIntent(tt.intents)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCalculationManager_setupServiceFunctionChainCalculation(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestCalculationManager_setupServiceFunctionChainCalculation",
			wantErr: false,
		},
		{
			name:    "TestCalculationManager_setupServiceFunctionChainCalculation error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			calculationSetup := NewMockCalculationSetup(controller)
			calculationTransformer := NewMockCalculationTransformer(controller)
			calculationUpdater := NewMockCalculationUpdater(controller)
			manager := NewCalculationManager(cacheMock, graphMock, calculationSetup, calculationTransformer, calculationUpdater)
			intent := domain.NewDomainIntent(domain.IntentTypeSFC, []domain.Value{})
			manager.algorithm = 0
			calculationOptions := &CalculationOptions{
				graph: graphMock,
			}
			if tt.wantErr {
				calculationSetup.EXPECT().PerformServiceFunctionChainSetup(gomock.Any()).Return(nil, fmt.Errorf("SFC setup failed"))
				err := manager.setupServiceFunctionChainCalculation(intent, calculationOptions)
				assert.Error(t, err)
			} else {
				SfcCalculationOptions := &SfcCalculationOptions{}
				calculationSetup.EXPECT().PerformServiceFunctionChainSetup(gomock.Any()).Return(SfcCalculationOptions, nil)
				err := manager.setupServiceFunctionChainCalculation(intent, calculationOptions)
				assert.NoError(t, err)
			}
		})
	}
}

func TestCalculationManager_CalculateBestPath(t *testing.T) {
	srAlgorithm := []uint32{0}
	nodes := map[int]graph.Node{
		1: graph.NewNetworkNode("1", "1", srAlgorithm),
		2: graph.NewNetworkNode("2", "2", srAlgorithm),
		3: graph.NewNetworkNode("3", "3", srAlgorithm),
		4: graph.NewNetworkNode("4", "4", srAlgorithm),
		5: graph.NewNetworkNode("5", "5", srAlgorithm),
		6: graph.NewNetworkNode("6", "6", srAlgorithm),
		7: graph.NewNetworkNode("7", "7", srAlgorithm),
		8: graph.NewNetworkNode("8", "8", srAlgorithm),
	}
	// 	     [1]
	//      / | \
	//    1/ 2|  \1
	//    /   |   \
	//  [2]  [3]  [4]
	//   |1   |4   \1
	//   |    |     \
	//  [5]  [6]-1-[7]
	//   \    |    /
	//    \6  |1  /5
	//     \  |  /
	//       [8]

	edges := map[int]graph.Edge{
		1:  graph.NewNetworkEdge("1", nodes[1], nodes[2], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency 1ms jitter 1us loss 1%
		2:  graph.NewNetworkEdge("2", nodes[1], nodes[3], map[helper.WeightKey]float64{helper.LatencyKey: 2000, helper.JitterKey: 20, helper.PacketLossKey: 2}),  // latency 2ms jitter 2us loss 2%
		3:  graph.NewNetworkEdge("3", nodes[1], nodes[4], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency 1ms jitter 1us loss 1%
		4:  graph.NewNetworkEdge("4", nodes[2], nodes[5], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency 1ms jitter 1us loss 1%
		5:  graph.NewNetworkEdge("5", nodes[3], nodes[5], map[helper.WeightKey]float64{helper.LatencyKey: 3000, helper.JitterKey: 30, helper.PacketLossKey: 3}),  // latency 3ms jitter 3us loss 3%
		6:  graph.NewNetworkEdge("6", nodes[3], nodes[6], map[helper.WeightKey]float64{helper.LatencyKey: 4000, helper.JitterKey: 40, helper.PacketLossKey: 4}),  // latency 4ms jitter 4us loss 4%
		7:  graph.NewNetworkEdge("7", nodes[4], nodes[7], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency ms jitter 1us loss 1%
		8:  graph.NewNetworkEdge("8", nodes[5], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 6000, helper.JitterKey: 60, helper.PacketLossKey: 6}),  // latency 6ms jitter 6us loss 6%
		9:  graph.NewNetworkEdge("9", nodes[6], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency 1ms jitter 1us loss 1%
		10: graph.NewNetworkEdge("10", nodes[7], nodes[6], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}), // latency 1ms jitter 1us loss 1%
		11: graph.NewNetworkEdge("11", nodes[7], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 5000, helper.JitterKey: 50, helper.PacketLossKey: 5}), // latency 5ms jitter 5us loss 5%
	}
	tests := []struct {
		name               string
		wantSetupErr       bool
		wantCalculationErr bool
	}{
		{
			name:               "TestCalculationManager_CalculateBestPath setup error",
			wantSetupErr:       true,
			wantCalculationErr: false,
		},
		{
			name:               "TestCalculationManager_CalculateBestPath calculation error",
			wantCalculationErr: true,
			wantSetupErr:       false,
		},
		{
			name:               "TestCalculationManager_CalculateBestPath success",
			wantSetupErr:       false,
			wantCalculationErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			calculationSetup := NewMockCalculationSetup(controller)
			calculationTransformer := NewMockCalculationTransformer(controller)
			calculationUpdater := NewMockCalculationUpdater(controller)
			manager := NewCalculationManager(cacheMock, graphMock, calculationSetup, calculationTransformer, calculationUpdater)
			stream := api.NewMockIntentController_GetIntentPathServer(controller)
			ctx := context.Background()
			sourceAddress := "2001:db8::1"
			destinationAddress := "2001:db8::2"
			intents := []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}
			pathRequest, err := domain.NewDomainPathRequest(sourceAddress, destinationAddress, intents, stream, ctx)
			assert.NoError(t, err)
			graphMock.EXPECT().Lock().Return().AnyTimes()
			graphMock.EXPECT().Unlock().Return().AnyTimes()
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()

			if tt.wantSetupErr {
				calculationSetup.EXPECT().PerformSetup(pathRequest).Return(nil, fmt.Errorf("setup failed"))
				_, err := manager.CalculateBestPath(pathRequest)
				assert.Error(t, err)
				return
			}
			if tt.wantCalculationErr {
				graph, err := setupGraph(nodes, map[int]graph.Edge{})
				assert.Nil(t, err)
				manager.graph = graph
				calculationOptions := &CalculationOptions{
					graph:           graph,
					sourceNode:      nodes[1],
					destinationNode: nodes[8],
					weightKeys:      []helper.WeightKey{helper.LatencyKey},
					calculationMode: CalculationModeSum,
				}
				calculationSetup.EXPECT().PerformSetup(pathRequest).Return(calculationOptions, nil)
				_, err = manager.CalculateBestPath(pathRequest)
				assert.Error(t, err)
			} else {
				graph, err := setupGraph(nodes, edges)
				assert.Nil(t, err)
				manager.graph = graph
				calculationOptions := &CalculationOptions{
					graph:           graph,
					sourceNode:      nodes[1],
					destinationNode: nodes[8],
					weightKeys:      []helper.WeightKey{helper.LatencyKey},
					calculationMode: CalculationModeSum,
				}
				calculationSetup.EXPECT().PerformSetup(pathRequest).Return(calculationOptions, nil)
				calculationTransformer.EXPECT().TransformResult(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				_, err = manager.CalculateBestPath(pathRequest)
				assert.NoError(t, err)
			}
		})
	}
}

func TestCalculationManager_getCalculationUpdateOptinos(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestCalculationManager_getCalculationUpdateOptinos",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			calculationSetup := NewMockCalculationSetup(controller)
			calculationTransformer := NewMockCalculationTransformer(controller)
			calculationUpdater := NewMockCalculationUpdater(controller)
			manager := NewCalculationManager(cacheMock, graphMock, calculationSetup, calculationTransformer, calculationUpdater)
			currentPathResult := domain.NewMockPathResult(controller)
			pathRequest := domain.NewMockPathRequest(controller)
			streamSession := domain.NewDomainStreamSession(pathRequest, currentPathResult)
			currentPathResult.EXPECT().GetIpv6SidAddresses().Return([]string{"2001:db8::1", "2001:db8::2"})
			pathRequest.EXPECT().GetIntents().Return([]domain.Intent{})
			calculationSetup.EXPECT().GetWeightKeysandCalculationMode(gomock.Any()).Return([]helper.WeightKey{}, CalculationModeSum)
			calculationUpdateOptions := manager.getCalculationUpdateOptions(streamSession)
			assert.NotNil(t, calculationUpdateOptions)
			assert.Equal(t, []string{"2001:db8::1", "2001:db8::2"}, calculationUpdateOptions.currentAppliedSidList)
			assert.Equal(t, []helper.WeightKey{}, calculationUpdateOptions.weightKeys)
			assert.Equal(t, CalculationModeSum, calculationUpdateOptions.calculationMode)
			assert.Equal(t, pathRequest, calculationUpdateOptions.pathRequest)
			assert.Equal(t, currentPathResult, calculationUpdateOptions.currentPathResult)
		})
	}
}

func TestCalculationManager_CalculatePathUpdate(t *testing.T) {
	srAlgorithm := []uint32{0}
	nodes := map[int]graph.Node{
		1: graph.NewNetworkNode("1", "1", srAlgorithm),
		2: graph.NewNetworkNode("2", "2", srAlgorithm),
		3: graph.NewNetworkNode("3", "3", srAlgorithm),
		4: graph.NewNetworkNode("4", "4", srAlgorithm),
		5: graph.NewNetworkNode("5", "5", srAlgorithm),
		6: graph.NewNetworkNode("6", "6", srAlgorithm),
		7: graph.NewNetworkNode("7", "7", srAlgorithm),
		8: graph.NewNetworkNode("8", "8", srAlgorithm),
	}
	// 	     [1]
	//      / | \
	//    1/ 2|  \1
	//    /   |   \
	//  [2]  [3]  [4]
	//   |1   |4   \1
	//   |    |     \
	//  [5]  [6]-1-[7]
	//   \    |    /
	//    \6  |1  /5
	//     \  |  /
	//       [8]

	edges := map[int]graph.Edge{
		1:  graph.NewNetworkEdge("1", nodes[1], nodes[2], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency 1ms jitter 1us loss 1%
		2:  graph.NewNetworkEdge("2", nodes[1], nodes[3], map[helper.WeightKey]float64{helper.LatencyKey: 2000, helper.JitterKey: 20, helper.PacketLossKey: 2}),  // latency 2ms jitter 2us loss 2%
		3:  graph.NewNetworkEdge("3", nodes[1], nodes[4], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency 1ms jitter 1us loss 1%
		4:  graph.NewNetworkEdge("4", nodes[2], nodes[5], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency 1ms jitter 1us loss 1%
		5:  graph.NewNetworkEdge("5", nodes[3], nodes[5], map[helper.WeightKey]float64{helper.LatencyKey: 3000, helper.JitterKey: 30, helper.PacketLossKey: 3}),  // latency 3ms jitter 3us loss 3%
		6:  graph.NewNetworkEdge("6", nodes[3], nodes[6], map[helper.WeightKey]float64{helper.LatencyKey: 4000, helper.JitterKey: 40, helper.PacketLossKey: 4}),  // latency 4ms jitter 4us loss 4%
		7:  graph.NewNetworkEdge("7", nodes[4], nodes[7], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency ms jitter 1us loss 1%
		8:  graph.NewNetworkEdge("8", nodes[5], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 6000, helper.JitterKey: 60, helper.PacketLossKey: 6}),  // latency 6ms jitter 6us loss 6%
		9:  graph.NewNetworkEdge("9", nodes[6], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency 1ms jitter 1us loss 1%
		10: graph.NewNetworkEdge("10", nodes[7], nodes[6], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}), // latency 1ms jitter 1us loss 1%
		11: graph.NewNetworkEdge("11", nodes[7], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 5000, helper.JitterKey: 50, helper.PacketLossKey: 5}), // latency 5ms jitter 5us loss 5%
	}
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestCalculationManager_CalculatePathUpdate success",
			wantErr: false,
		},
		{
			name:    "TestCalculationManager_CalculatePathUpdate error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			calculationSetup := NewMockCalculationSetup(controller)
			calculationTransformer := NewMockCalculationTransformer(controller)
			calculationUpdater := NewMockCalculationUpdater(controller)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			network, err := setupGraph(nodes, edges)
			assert.Nil(t, err)
			manager := NewCalculationManager(cacheMock, network, calculationSetup, calculationTransformer, calculationUpdater)
			stream := api.NewMockIntentController_GetIntentPathServer(controller)
			intents := []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}
			pathRequest, err := domain.NewDomainPathRequest("2001:db8::1", "2001:db8::2", intents, stream, context.Background())
			assert.NoError(t, err)
			calculationMode := CalculationModeSum
			weightKeys := []helper.WeightKey{helper.LatencyKey}
			calculationOptions := &CalculationOptions{
				sourceNode:      nodes[1],
				destinationNode: nodes[8],
				weightKeys:      weightKeys,
				calculationMode: calculationMode,
			}
			calculationSetup.EXPECT().PerformSetup(pathRequest).Return(calculationOptions, nil).AnyTimes()
			path := graph.NewMockPath(controller)
			pathResult, err := domain.NewDomainPathResult(pathRequest, path, []string{"2001:db8::1", "2001:db8::2"})
			assert.NoError(t, err)
			calculationSetup.EXPECT().GetWeightKeysandCalculationMode(gomock.Any()).Return(weightKeys, calculationMode)
			streamSession := domain.NewDomainStreamSession(pathRequest, pathResult)
			if tt.wantErr {
				for _, edge := range nodes[1].GetEdges() {
					network.DeleteEdge(edge)
				}
				_, err := manager.CalculatePathUpdate(streamSession)
				assert.Error(t, err)
			} else {
				calculationTransformer.EXPECT().TransformResult(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				calculationUpdater.EXPECT().UpdateCalculation(gomock.Any()).Return(pathResult, nil)
				_, err := manager.CalculatePathUpdate(streamSession)
				assert.NoError(t, err)
			}
		})
	}
}
