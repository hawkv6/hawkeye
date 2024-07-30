package calculation

import (
	"testing"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func TestNewCalculationUpdaterService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewCalculationUpdaterService",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			assert.NotNil(t, NewCalculationUpdaterService(cacheMock, graphMock))

		})
	}
}

func TestCalculationUpdateService_getInitialTotalCost(t *testing.T) {
	tests := []struct {
		name        string
		weightTypes []helper.WeightKey
		want        float64
	}{
		{
			name:        "Test getInitialTotalCost with one weight type packet loss",
			weightTypes: []helper.WeightKey{helper.PacketLossKey},
			want:        1.0,
		},
		{
			name:        "Test getInitialTotalCost with one weight type other than packet loss",
			weightTypes: []helper.WeightKey{helper.LatencyKey},
			want:        0.0,
		},
		{
			name:        "Test getInitialTotalCost with multiple weight types",
			weightTypes: []helper.WeightKey{helper.LatencyKey, helper.PacketLossKey},
			want:        0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := CalculationUpdaterService{}
			assert.Equal(t, tt.want, service.getInitialTotalCost(tt.weightTypes))
		})
	}
}

func TestCalculationUpdateService_getUpdatedTotalCost(t *testing.T) {
	tests := []struct {
		name             string
		weightTypes      []helper.WeightKey
		initialTotalCost float64
		want             float64
		weights          map[helper.WeightKey][]float64
		wantErr          bool
	}{
		{
			name:             "Test getNewTotalCost with one weight type packet error - edge not found",
			weightTypes:      []helper.WeightKey{helper.PacketLossKey},
			initialTotalCost: 1.0,
			want:             0.0,
			wantErr:          true,
		},
		{
			name:             "Test getNewTotalCost with one weight type packet loss success",
			weightTypes:      []helper.WeightKey{helper.PacketLossKey},
			initialTotalCost: 1.0,
			want:             1 * (1 - 0.01) * (1 - 0.02),
			wantErr:          false,
			weights:          map[helper.WeightKey][]float64{helper.PacketLossKey: {0.01, 0.02}},
		},
		{
			name:             "Test getNewTotalCost with one weight type latency success",
			weightTypes:      []helper.WeightKey{helper.LatencyKey},
			initialTotalCost: 0.0,
			want:             3.0,
			wantErr:          false,
			weights:          map[helper.WeightKey][]float64{helper.LatencyKey: {1.0, 2.0}},
		},
		{
			name:             "Test getNewTotalCost with two weight types latency + packet loss success",
			weightTypes:      []helper.WeightKey{helper.NormalizedLatencyKey, helper.NormalizedPacketLossKey},
			initialTotalCost: 0.0,
			want:             (float64(helper.TwoFactorWeights[0])*0.1 + float64(helper.TwoFactorWeights[1])*0.1) + (float64(helper.TwoFactorWeights[0])*0.2 + float64(helper.TwoFactorWeights[1])*0.2),
			wantErr:          false,
			weights:          map[helper.WeightKey][]float64{helper.NormalizedLatencyKey: {0.1, 0.2}, helper.NormalizedPacketLossKey: {0.1, 0.2}},
		},
		{
			name:             "Test getNewTotalCost with two weight types latency + packet loss, jitter success",
			weightTypes:      []helper.WeightKey{helper.NormalizedLatencyKey, helper.NormalizedPacketLossKey, helper.NormalizedJitterKey},
			initialTotalCost: 0.0,
			want:             (float64(helper.ThreeFactorWeights[0])*0.1 + float64(helper.ThreeFactorWeights[1])*0.1 + float64(helper.ThreeFactorWeights[2])*0.1) + (float64(helper.ThreeFactorWeights[0])*0.2 + float64(helper.ThreeFactorWeights[1])*0.2 + float64(helper.ThreeFactorWeights[2])*0.2),
			wantErr:          false,
			weights:          map[helper.WeightKey][]float64{helper.NormalizedLatencyKey: {0.1, 0.2}, helper.NormalizedPacketLossKey: {0.1, 0.2}, helper.NormalizedJitterKey: {0.1, 0.2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			if tt.wantErr {
				testGraph := graph.NewMockGraph(controller)
				edgeMock := graph.NewMockEdge(controller)
				testGraph.EXPECT().GetEdge(gomock.Any()).Return(nil)
				service := NewCalculationUpdaterService(cacheMock, testGraph)
				path := graph.NewMockPath(controller)
				path.EXPECT().GetEdges().Return([]graph.Edge{edgeMock})
				edgeMock.EXPECT().GetId().Return("1").AnyTimes()
				pathRequest := domain.NewMockPathRequest(controller)
				pathResult, err := domain.NewDomainPathResult(pathRequest, path, []string{})
				assert.NoError(t, err)
				totalCost, err := service.getUpdatedTotalCost(pathResult.GetEdges(), tt.weightTypes, tt.initialTotalCost)
				assert.Error(t, err)
				assert.Equal(t, 0.0, totalCost)
				return
			}
			testGraph := graph.NewMockGraph(controller)
			edgeMock1 := graph.NewMockEdge(controller)
			edgeMock2 := graph.NewMockEdge(controller)
			service := NewCalculationUpdaterService(cacheMock, testGraph)
			path := graph.NewMockPath(controller)
			path.EXPECT().GetEdges().Return([]graph.Edge{edgeMock1, edgeMock2})
			edgeMock1.EXPECT().GetId().Return("1").AnyTimes()
			edgeMock2.EXPECT().GetId().Return("2").AnyTimes()
			testGraph.EXPECT().GetEdge("1").Return(edgeMock1)
			testGraph.EXPECT().GetEdge("2").Return(edgeMock2)
			if len(tt.weightTypes) == 1 {
				edgeMock1.EXPECT().GetWeight(tt.weightTypes[0]).Return(tt.weights[tt.weightTypes[0]][0]).AnyTimes()
				edgeMock2.EXPECT().GetWeight(tt.weightTypes[0]).Return(tt.weights[tt.weightTypes[0]][1]).AnyTimes()
			} else if len(tt.weightTypes) == 2 {
				edgeMock1.EXPECT().GetWeight(tt.weightTypes[0]).Return(tt.weights[tt.weightTypes[0]][0]).AnyTimes()
				edgeMock2.EXPECT().GetWeight(tt.weightTypes[0]).Return(tt.weights[tt.weightTypes[0]][1]).AnyTimes()
				edgeMock1.EXPECT().GetWeight(tt.weightTypes[1]).Return(tt.weights[tt.weightTypes[1]][0]).AnyTimes()
				edgeMock2.EXPECT().GetWeight(tt.weightTypes[1]).Return(tt.weights[tt.weightTypes[1]][1]).AnyTimes()
			} else {
				edgeMock1.EXPECT().GetWeight(tt.weightTypes[0]).Return(tt.weights[tt.weightTypes[0]][0]).AnyTimes()
				edgeMock2.EXPECT().GetWeight(tt.weightTypes[0]).Return(tt.weights[tt.weightTypes[0]][1]).AnyTimes()
				edgeMock1.EXPECT().GetWeight(tt.weightTypes[1]).Return(tt.weights[tt.weightTypes[1]][0]).AnyTimes()
				edgeMock2.EXPECT().GetWeight(tt.weightTypes[1]).Return(tt.weights[tt.weightTypes[1]][1]).AnyTimes()
				edgeMock1.EXPECT().GetWeight(tt.weightTypes[2]).Return(tt.weights[tt.weightTypes[2]][0]).AnyTimes()
				edgeMock2.EXPECT().GetWeight(tt.weightTypes[2]).Return(tt.weights[tt.weightTypes[2]][1]).AnyTimes()
			}

			pathRequest := domain.NewMockPathRequest(controller)
			pathResult, err := domain.NewDomainPathResult(pathRequest, path, []string{})
			assert.NoError(t, err)
			totalCost, err := service.getUpdatedTotalCost(pathResult.GetEdges(), tt.weightTypes, tt.initialTotalCost)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, totalCost)
		})
	}
}

func TestCalculationUpdateService_updateTotalCost(t *testing.T) {
	tests := []struct {
		name             string
		weightTypes      []helper.WeightKey
		weights          map[helper.WeightKey][]float64
		wantErr          bool
		totalCostChanged bool
	}{
		{
			name:        "Test calculateTotalCost with one weight type packet loss error - edge not found",
			weightTypes: []helper.WeightKey{helper.PacketLossKey},
			wantErr:     true,
		},
		{
			name:             "Test calculateTotalCost with one weight type packet loss success totalCost changed",
			weightTypes:      []helper.WeightKey{helper.PacketLossKey},
			wantErr:          false,
			weights:          map[helper.WeightKey][]float64{helper.PacketLossKey: {0.01, 0.02}},
			totalCostChanged: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			if tt.wantErr {
				testGraph := graph.NewMockGraph(controller)
				edgeMock := graph.NewMockEdge(controller)
				testGraph.EXPECT().GetEdge(gomock.Any()).Return(nil)
				service := NewCalculationUpdaterService(cacheMock, testGraph)
				path := graph.NewMockPath(controller)
				path.EXPECT().GetEdges().Return([]graph.Edge{edgeMock})
				edgeMock.EXPECT().GetId().Return("1").AnyTimes()
				pathRequest := domain.NewMockPathRequest(controller)
				pathResult, err := domain.NewDomainPathResult(pathRequest, path, []string{})
				assert.NoError(t, err)
				err = service.updateTotalCost(pathResult, tt.weightTypes)
				assert.Error(t, err)
				return
			}
			testGraph := graph.NewMockGraph(controller)
			edgeMock1 := graph.NewMockEdge(controller)
			edgeMock2 := graph.NewMockEdge(controller)
			service := NewCalculationUpdaterService(cacheMock, testGraph)
			path := graph.NewMockPath(controller)
			path.EXPECT().GetEdges().Return([]graph.Edge{edgeMock1, edgeMock2})
			edgeMock1.EXPECT().GetId().Return("1").AnyTimes()
			edgeMock2.EXPECT().GetId().Return("2").AnyTimes()
			testGraph.EXPECT().GetEdge("1").Return(edgeMock1)
			testGraph.EXPECT().GetEdge("2").Return(edgeMock2)
			edgeMock1.EXPECT().GetWeight(tt.weightTypes[0]).Return(tt.weights[tt.weightTypes[0]][0]).AnyTimes()
			edgeMock2.EXPECT().GetWeight(tt.weightTypes[0]).Return(tt.weights[tt.weightTypes[0]][1]).AnyTimes()

			pathRequest := domain.NewMockPathRequest(controller)
			pathResult, err := domain.NewDomainPathResult(pathRequest, path, []string{})
			path.EXPECT().GetTotalCost().Return(0.0)
			path.EXPECT().SetTotalCost(gomock.Any()).AnyTimes()
			assert.NoError(t, err)
			err = service.updateTotalCost(pathResult, tt.weightTypes)
			assert.NoError(t, err)
		})
	}
}

func TestCalculationUpdateService_getUpdatedBottleneckValues(t *testing.T) {
	tests := []struct {
		name       string
		weightType helper.WeightKey
		minValue   float64
		wantErr    bool
		newValue   float64
	}{
		{
			name:       "Test getUpdatedBottleneckValues with edge not found",
			weightType: helper.PacketLossKey,
			wantErr:    true,
		},
		{
			name:       "Test getUpdatedBottleneckValues with edge found",
			weightType: helper.LatencyKey,
			wantErr:    false,
			minValue:   20,
			newValue:   10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			testGraph := graph.NewMockGraph(controller)
			edgeMock := graph.NewMockEdge(controller)
			service := NewCalculationUpdaterService(cacheMock, testGraph)
			if tt.wantErr {
				testGraph.EXPECT().GetEdge(gomock.Any()).Return(nil)
				edgeMock.EXPECT().GetId().Return("1").AnyTimes()
				_, _, err := service.getUpdatedBottleneckValues([]graph.Edge{edgeMock}, tt.weightType, 0.0, edgeMock)
				assert.Error(t, err)
				return
			}
			testGraph.EXPECT().GetEdge(gomock.Any()).Return(edgeMock)
			edgeMock.EXPECT().GetId().Return("1").AnyTimes()
			edgeMock.EXPECT().GetWeight(tt.weightType).Return(tt.newValue).AnyTimes()
			edgeMock.EXPECT().GetWeight(tt.weightType).Return(float64(100)).AnyTimes()
			minValue, bottleneckEdge, err := service.getUpdatedBottleneckValues([]graph.Edge{edgeMock}, tt.weightType, tt.minValue, edgeMock)
			assert.NoError(t, err)
			assert.NotNil(t, bottleneckEdge)
			assert.Equal(t, tt.newValue, minValue)
		})
	}
}

func TestCalculationUpdateService_updateBottlneckValues(t *testing.T) {
	tests := []struct {
		name              string
		weightType        helper.WeightKey
		minValue          float64
		bottleNeckChanged bool
		minimumChanged    bool
		oldValue          float64
		newValue          float64
	}{
		{
			name:              "Test updateBottleneckValues bottleneck edge changed",
			weightType:        helper.PacketLossKey,
			bottleNeckChanged: true,
			newValue:          10,
		},
		{
			name:       "Test updateBottleneckValues bottleneck edge not changed but its minimum value",
			weightType: helper.LatencyKey,
			oldValue:   20,
			newValue:   10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			testGraph := graph.NewMockGraph(controller)
			service := NewCalculationUpdaterService(cacheMock, testGraph)
			pathResult := domain.NewMockPathResult(controller)
			mockEdge := graph.NewMockEdge(controller)
			pathResult.EXPECT().GetBottleneckValue().Return(tt.oldValue).AnyTimes()
			pathResult.EXPECT().SetBottleneckValue(tt.newValue).AnyTimes()
			if tt.bottleNeckChanged {
				mockEdge2 := graph.NewMockEdge(controller)
				pathResult.EXPECT().GetBottleneckEdge().Return(mockEdge2).AnyTimes()
				pathResult.EXPECT().SetBottleneckEdge(mockEdge)
			} else {
				pathResult.EXPECT().GetBottleneckEdge().Return(mockEdge)
			}
			service.updateBottleneckValues(pathResult, mockEdge, tt.newValue)
		})
	}
}

func TestCalculationUpdateService_updateMinimumValue(t *testing.T) {
	tests := []struct {
		name       string
		weightType helper.WeightKey
		wantErr    bool
	}{
		{
			name:       "Test updateMinimumValue with error",
			weightType: helper.LatencyKey,
			wantErr:    true,
		},
		{
			name:       "Test updateMinimumValue without error",
			weightType: helper.LatencyKey,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			testGraph := graph.NewMockGraph(controller)
			edgeMock := graph.NewMockEdge(controller)
			service := NewCalculationUpdaterService(cacheMock, testGraph)
			if tt.wantErr {
				testGraph.EXPECT().GetEdge(gomock.Any()).Return(nil)
				edgeMock.EXPECT().GetId().Return("1").AnyTimes()
				pathResult := domain.NewMockPathResult(controller)
				pathResult.EXPECT().GetEdges().Return([]graph.Edge{edgeMock})
				err := service.updateMinimumValue(pathResult, tt.weightType)
				assert.Error(t, err)
				return
			}
			testGraph.EXPECT().GetEdge(gomock.Any()).Return(edgeMock)
			edgeMock.EXPECT().GetId().Return("1").AnyTimes()
			edgeMock.EXPECT().GetWeight(tt.weightType).Return(float64(100)).AnyTimes()
			pathResult := domain.NewMockPathResult(controller)
			pathResult.EXPECT().GetEdges().Return([]graph.Edge{edgeMock})
			pathResult.EXPECT().GetBottleneckEdge().Return(edgeMock).AnyTimes()
			pathResult.EXPECT().GetBottleneckValue().Return(float64(100)).AnyTimes()
			pathResult.EXPECT().SetBottleneckEdge(gomock.Any()).AnyTimes()
			pathResult.EXPECT().SetBottleneckValue(gomock.Any()).AnyTimes()
			err := service.updateMinimumValue(pathResult, tt.weightType)
			assert.NoError(t, err)
		})
	}
}

func TestCalculationUpdateService_updateCurrentResult(t *testing.T) {
	tests := []struct {
		name            string
		calculationMode CalculationMode
		wantErr         bool
	}{
		{
			name:            "Test updateCurrentResult with CalculationModeSum",
			calculationMode: CalculationModeSum,
			wantErr:         false,
		},
		{
			name:            "Test updateCurrentResult with CalculationModeMin",
			calculationMode: CalculationModeMin,
			wantErr:         false,
		},
		{
			name:            "Test updateCurrentResult with CalculationModeSum error",
			calculationMode: CalculationModeSum,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			testGraph := graph.NewMockGraph(controller)
			service := NewCalculationUpdaterService(cacheMock, testGraph)
			pathResult := domain.NewMockPathResult(controller)
			edgeMock := graph.NewMockEdge(controller)

			pathResult.EXPECT().GetEdges().Return([]graph.Edge{edgeMock}).AnyTimes()
			if tt.wantErr {
				testGraph.EXPECT().GetEdge(gomock.Any()).Return(nil)
				edgeMock.EXPECT().GetId().Return("1").AnyTimes()
				err := service.updateCurrentResult([]helper.WeightKey{helper.LatencyKey}, tt.calculationMode, pathResult)
				assert.Error(t, err)
			}

			edgeMock.EXPECT().GetId().Return("1").AnyTimes()
			testGraph.EXPECT().GetEdge("1").Return(edgeMock)
			edgeMock.EXPECT().GetWeight(gomock.Any()).Return(float64(100)).AnyTimes()
			if tt.calculationMode == CalculationModeSum {
				pathResult.EXPECT().GetTotalCost().Return(float64(100)).AnyTimes()
				pathResult.EXPECT().SetTotalCost(gomock.Any()).AnyTimes()
			} else {
				pathResult.EXPECT().GetBottleneckEdge().Return(edgeMock).AnyTimes()
				pathResult.EXPECT().GetBottleneckValue().Return(float64(100)).AnyTimes()
				pathResult.EXPECT().SetBottleneckEdge(gomock.Any()).AnyTimes()
				pathResult.EXPECT().SetBottleneckValue(gomock.Any()).AnyTimes()
			}
			err := service.updateCurrentResult([]helper.WeightKey{helper.LatencyKey}, tt.calculationMode, pathResult)
			assert.NoError(t, err)
		})
	}
}

func TestCalculationUpdateService_updatePathIfCostImproved(t *testing.T) {
	tests := []struct {
		name         string
		oldTotalCost float64
		newTotalCost float64
	}{
		{
			name:         "Test updatePathIfCostImproved apply new path",
			oldTotalCost: 20,
			newTotalCost: 10,
		},
		{
			name:         "Test updatePathIfCostImproved keep old path",
			oldTotalCost: 10,
			newTotalCost: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			testGraph := graph.NewMockGraph(controller)
			service := NewCalculationUpdaterService(cacheMock, testGraph)
			newPathResult := domain.NewMockPathResult(controller)
			oldPathResult := domain.NewMockPathResult(controller)
			newPathResult.EXPECT().GetTotalCost().Return(tt.newTotalCost).AnyTimes()
			oldPathResult.EXPECT().GetTotalCost().Return(tt.oldTotalCost).AnyTimes()
			newPathResult.EXPECT().SetTotalCost(gomock.Any()).AnyTimes()
			pathRequest := domain.NewMockPathRequest(controller)
			streamSession := domain.NewDomainStreamSession(pathRequest, oldPathResult)
			if tt.oldTotalCost > tt.newTotalCost*(1-helper.FlappingThreshold) {
				pathResult := service.updatePathIfCostImproved(oldPathResult, newPathResult, streamSession)
				assert.Equal(t, newPathResult, pathResult)
				return
			}
			assert.Nil(t, service.updatePathIfCostImproved(oldPathResult, newPathResult, streamSession))
		})
	}
}

func TestCalculationUpdateService_updatePathIfMinimumImproved(t *testing.T) {
	tests := []struct {
		name            string
		oldMinimumValue float64
		newMinimumValue float64
	}{
		{
			name:            "Test updatePathIfMinimumImproved apply new path",
			oldMinimumValue: 20,
			newMinimumValue: 10,
		},
		{
			name:            "Test updatePathIfMinimumImproved keep old path",
			oldMinimumValue: 10,
			newMinimumValue: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			testGraph := graph.NewMockGraph(controller)
			service := NewCalculationUpdaterService(cacheMock, testGraph)
			newPathResult := domain.NewMockPathResult(controller)
			oldPathResult := domain.NewMockPathResult(controller)
			newPathResult.EXPECT().GetBottleneckValue().Return(tt.newMinimumValue).AnyTimes()
			oldPathResult.EXPECT().GetBottleneckValue().Return(tt.oldMinimumValue).AnyTimes()
			if tt.oldMinimumValue > tt.newMinimumValue*(1-helper.FlappingThreshold) {
				newBottleneckEdge := graph.NewMockEdge(controller)
				oldBottleneckEdge := graph.NewMockEdge(controller)
				oldPathResult.EXPECT().GetBottleneckEdge().Return(oldBottleneckEdge).AnyTimes()
				newPathResult.EXPECT().GetBottleneckEdge().Return(newBottleneckEdge).AnyTimes()
				newPathResult.EXPECT().SetBottleneckValue(gomock.Any()).AnyTimes()
				pathRequest := domain.NewMockPathRequest(controller)
				streamSession := domain.NewDomainStreamSession(pathRequest, oldPathResult)
				pathResult := service.updatePathIfMinimumImproved(oldPathResult, newPathResult, streamSession)
				assert.Equal(t, newPathResult, pathResult)
				return
			}
			pathRequest := domain.NewMockPathRequest(controller)
			streamSession := domain.NewDomainStreamSession(pathRequest, oldPathResult)
			assert.Nil(t, service.updatePathIfMinimumImproved(oldPathResult, newPathResult, streamSession))
		})
	}
}

func TestCalculationUpdateService_currentServicesStillValid(t *testing.T) {
	tests := []struct {
		name           string
		serviceSidList []string
		want           bool
	}{
		{
			name:           "Test areCurrentServicesStillValid with valid services",
			serviceSidList: []string{"1", "2"},
			want:           true,
		},
		{
			name:           "Test areCurrentServicesStillValid with invalid services",
			serviceSidList: []string{"1", "2"},
			want:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			service := NewCalculationUpdaterService(cacheMock, nil)
			if tt.want {
				cacheMock.EXPECT().DoesServiceSidExist(gomock.Any()).Return(true).AnyTimes()
			} else {
				cacheMock.EXPECT().DoesServiceSidExist(gomock.Any()).Return(false)
			}
			assert.Equal(t, tt.want, service.currentServicesStillValid(tt.serviceSidList))
		})
	}
}

func TestCalculationUpdateService_currentServicesNotValidAnymore(t *testing.T) {
	tests := []struct {
		name           string
		serviceSidList []string
		want           bool
	}{
		{
			name:           "Test areCurrentServicesStillValid with valid services",
			serviceSidList: []string{"1", "2"},
			want:           false,
		},
		{
			name:           "Test areCurrentServicesStillValid with invalid services",
			serviceSidList: []string{"1", "2"},
			want:           true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			service := NewCalculationUpdaterService(cacheMock, nil)
			if tt.want {
				cacheMock.EXPECT().DoesServiceSidExist(gomock.Any()).Return(false).AnyTimes()
			} else {
				cacheMock.EXPECT().DoesServiceSidExist(gomock.Any()).Return(true).AnyTimes()
			}
			fwValue, err := domain.NewStringValue(domain.ValueTypeSFC, proto.String("fw"))
			assert.NoError(t, err)
			firstIntent := domain.NewDomainIntent(domain.IntentTypeSFC, []domain.Value{fwValue})
			assert.NoError(t, err)
			currentPathResult := domain.NewMockPathResult(controller)
			currentPathResult.EXPECT().GetServiceSidList().Return(tt.serviceSidList).AnyTimes()
			assert.NoError(t, err)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, service.currentServicesNotValidAnymore(firstIntent, currentPathResult))
			assert.Equal(t, tt.want, service.currentServicesNotValidAnymore(firstIntent, currentPathResult))

		})
	}
}

func TestCalculationUpdateService_currentPathNotValidAnymore(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Test currentPathNotValidAnymore with valid path",
			want: false,
		},
		{
			name: "Test currentPathNotValidAnymore with invalid path",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			service := NewCalculationUpdaterService(cacheMock, graphMock)
			currentPathResult := domain.NewMockPathResult(controller)
			edgeMock := graph.NewMockEdge(controller)
			currentPathResult.EXPECT().GetEdges().Return([]graph.Edge{edgeMock}).AnyTimes()
			weightKeys := []helper.WeightKey{helper.LatencyKey}
			calculationMode := CalculationModeSum
			edgeMock.EXPECT().GetId().Return("1").AnyTimes()
			if tt.want {
				graphMock.EXPECT().GetEdge("1").Return(nil)
			} else {
				graphMock.EXPECT().GetEdge("1").Return(edgeMock)
				edgeMock.EXPECT().GetWeight(gomock.Any()).Return(float64(100)).AnyTimes()
				currentPathResult.EXPECT().GetTotalCost().Return(float64(100)).AnyTimes()
				newPathResult := domain.NewMockPathResult(controller)
				newPathResult.EXPECT().GetTotalCost().Return(float64(20)).AnyTimes()
			}
			assert.Equal(t, tt.want, service.currentPathNotValidAnymore(weightKeys, calculationMode, currentPathResult))
		})
	}
}

func TestCalculationUpdateService_handlePathChange(t *testing.T) {
	tests := []struct {
		name            string
		calculationMode CalculationMode
		pathValid       bool
	}{
		{
			name:            "Test handlePathChange with invalid path",
			calculationMode: CalculationModeSum,
			pathValid:       false,
		},
		{
			name:            "Test handlePathChange with valid path CalculationModeSum",
			calculationMode: CalculationModeSum,
			pathValid:       true,
		},
		{
			name:            "Test handlePathChange with valid path CalculationModeMin",
			calculationMode: CalculationModeMin,
			pathValid:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			service := NewCalculationUpdaterService(cacheMock, graphMock)
			pathRequest := domain.NewMockPathRequest(controller)
			calculationMode := tt.calculationMode
			weightKey := []helper.WeightKey{helper.LatencyKey}
			currentPathResult := domain.NewMockPathResult(controller)
			newPathResult := domain.NewMockPathResult(controller)
			streamSession := domain.NewDomainStreamSession(pathRequest, currentPathResult)
			if !tt.pathValid {
				fwValue, err := domain.NewStringValue(domain.ValueTypeSFC, proto.String("fw"))
				assert.NoError(t, err)
				sfcIntent := domain.NewDomainIntent(domain.IntentTypeSFC, []domain.Value{fwValue})
				pathRequest.EXPECT().GetIntents().Return([]domain.Intent{sfcIntent}).AnyTimes()
				currentPathResult.EXPECT().GetServiceSidList().Return([]string{"1", "2"}).AnyTimes()
				cacheMock.EXPECT().DoesServiceSidExist(gomock.Any()).Return(false).AnyTimes()
				pathResult := service.handlePathChange(weightKey, calculationMode, currentPathResult, newPathResult, streamSession)
				assert.Equal(t, newPathResult, pathResult)
				return
			}
			firstIntent := domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})
			intents := []domain.Intent{firstIntent}
			pathRequest.EXPECT().GetIntents().Return(intents).AnyTimes()
			currentPathResult.EXPECT().GetEdges().Return([]graph.Edge{}).AnyTimes()
			currentPathResult.EXPECT().GetTotalCost().Return(float64(100)).AnyTimes()
			currentPathResult.EXPECT().SetTotalCost(gomock.Any()).AnyTimes()
			edgeMock := graph.NewMockEdge(controller)
			currentPathResult.EXPECT().GetBottleneckEdge().Return(edgeMock).AnyTimes()
			currentPathResult.EXPECT().GetBottleneckValue().Return(float64(100)).AnyTimes()
			newPathResult.EXPECT().GetTotalCost().Return(float64(100)).AnyTimes()
			newPathResult.EXPECT().GetBottleneckEdge().Return(edgeMock).AnyTimes()
			newPathResult.EXPECT().GetBottleneckValue().Return(float64(100)).AnyTimes()
			currentPathResult.EXPECT().SetBottleneckEdge(gomock.Any()).AnyTimes()
			currentPathResult.EXPECT().SetBottleneckValue(gomock.Any()).AnyTimes()
			assert.Nil(t, service.handlePathChange(weightKey, calculationMode, currentPathResult, newPathResult, streamSession))
		})
	}
}

func TestCalculationUpdateService_UpdateCalculation(t *testing.T) {
	tests := []struct {
		name             string
		wantErr          bool
		identicalSidList bool
	}{
		{
			name:             "Test UpdateCalculation with different sid list",
			wantErr:          false,
			identicalSidList: false,
		},
		{
			name:             "Test UpdateCalculation with error",
			wantErr:          true,
			identicalSidList: true,
		},
		{
			name:             "Test UpdateCalculation with identical sid list",
			wantErr:          false,
			identicalSidList: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cacheMock := cache.NewMockCache(controller)
			graphMock := graph.NewMockGraph(controller)
			service := NewCalculationUpdaterService(cacheMock, graphMock)
			pathRequest := domain.NewMockPathRequest(controller)
			currentPathResult := domain.NewMockPathResult(controller)
			newPathResult := domain.NewMockPathResult(controller)
			streamSession := domain.NewDomainStreamSession(pathRequest, currentPathResult)

			calculationUpdateOptions := CalculationUpdateOptions{
				currentPathResult:     currentPathResult,
				currentAppliedSidList: []string{"2001:db8::1", "2001:db8::2"},
				weightKeys:            []helper.WeightKey{helper.LatencyKey},
				calculationMode:       CalculationModeSum,
				newPathResult:         newPathResult,
				pathRequest:           pathRequest,
				streamSession:         streamSession,
			}
			if !tt.identicalSidList {
				newPathResult.EXPECT().GetIpv6SidAddresses().Return([]string{"2001:db8::3", "2001:db8::4"}).AnyTimes()
				intent := domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})
				pathRequest.EXPECT().GetIntents().Return([]domain.Intent{intent}).AnyTimes()
				edgeMock := graph.NewMockEdge(controller)
				edgeMock.EXPECT().GetId().Return("1").AnyTimes()
				graphMock.EXPECT().GetEdge("1").Return(nil).AnyTimes()
				currentPathResult.EXPECT().GetEdges().Return([]graph.Edge{edgeMock}).AnyTimes()
				pathResult, err := service.UpdateCalculation(&calculationUpdateOptions)
				assert.NoError(t, err)
				assert.NotNil(t, pathResult)
				return
			}
			newPathResult.EXPECT().GetIpv6SidAddresses().Return([]string{"2001:db8::1", "2001:db8::2"}).AnyTimes()
			if tt.wantErr {
				intent := domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})
				pathRequest.EXPECT().GetIntents().Return([]domain.Intent{intent}).AnyTimes()
				edgeMock := graph.NewMockEdge(controller)
				edgeMock.EXPECT().GetId().Return("1").AnyTimes()
				graphMock.EXPECT().GetEdge(gomock.Any()).Return(nil).AnyTimes()
				currentPathResult.EXPECT().GetEdges().Return([]graph.Edge{edgeMock}).AnyTimes()
				pathResult, err := service.UpdateCalculation(&calculationUpdateOptions)
				assert.Error(t, err)
				assert.Nil(t, pathResult)
				return
			} else {
				edgeMock := graph.NewMockEdge(controller)
				edgeMock.EXPECT().GetId().Return("1").AnyTimes()
				edgeMock.EXPECT().GetWeight(gomock.Any()).Return(float64(100)).AnyTimes()
				currentPathResult.EXPECT().GetEdges().Return([]graph.Edge{edgeMock}).AnyTimes()
				currentPathResult.EXPECT().GetTotalCost().Return(float64(100)).AnyTimes()
				newPathResult.EXPECT().GetTotalCost().Return(float64(100)).AnyTimes()
				graphMock.EXPECT().GetEdge("1").Return(edgeMock).AnyTimes()
				pathResult, err := service.UpdateCalculation(&calculationUpdateOptions)
				assert.NoError(t, err)
				assert.Nil(t, pathResult)
				return
			}
		})
	}
}
