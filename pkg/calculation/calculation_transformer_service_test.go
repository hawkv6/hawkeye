package calculation

import (
	"testing"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewCalculationTransformerService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewCalculationTransformerService",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cache := cache.NewMockCache(controller)
			assert.NotNil(t, NewCalculationTransformerService(cache))
		})
	}
}

func TestCalculationTransformerService_getNodesFromPath(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestCalculationTransformerService_getNodesFromPath",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cache := cache.NewMockCache(controller)
			service := NewCalculationTransformerService(cache)
			path := graph.NewMockPath(controller)
			edge := graph.NewMockEdge(controller)
			path.EXPECT().GetEdges().Return([]graph.Edge{edge})
			to := graph.NewMockNode(controller)
			to.EXPECT().GetId().Return("to")
			edge.EXPECT().To().Return(to)
			assert.Equal(t, []string{"to"}, service.getNodesFromPath(path))
		})
	}
}

func TestCalculationTransformerService_translatePathToSidList(t *testing.T) {
	tests := []struct {
		name       string
		nodeId     string
		nodeSid    string
		serviceSid string
	}{
		{
			name:       "TestCalculationTransformerService_translatePathToSidList service and node sid",
			nodeId:     "to",
			nodeSid:    "2001:db8:1::",
			serviceSid: "2001:db8:f::",
		},
		{
			name:    "TestCalculationTransformerService_translatePathToSidList only node sid",
			nodeId:  "to",
			nodeSid: "2001:db8:1::",
		},
		{
			name:   "TestCalculationTransformerService_translatePathToSidList no node sid",
			nodeId: "to",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cache := cache.NewMockCache(controller)
			service := NewCalculationTransformerService(cache)
			path := graph.NewMockPath(controller)
			edge := graph.NewMockEdge(controller)
			path.EXPECT().GetEdges().Return([]graph.Edge{edge})
			to := graph.NewMockNode(controller)
			to.EXPECT().GetId().Return(tt.nodeId)
			edge.EXPECT().To().Return(to)
			path.EXPECT().GetRouterServiceMap().Return(map[string]string{tt.nodeId: tt.serviceSid})
			cache.EXPECT().GetSrAlgorithmSid(tt.nodeId, uint32(0)).Return(tt.nodeSid)
			gotSidList, gotServiceSidList := service.translatePathToSidList(path, uint32(0))
			if tt.nodeSid == "" {
				assert.Len(t, gotSidList, 0)
				return
			}
			assert.Equal(t, []string{tt.nodeSid, tt.serviceSid}, gotSidList)
			assert.Equal(t, []string{tt.serviceSid}, gotServiceSidList)
		})
	}
}

func TestCalculationTransformerService_TransformResult(t *testing.T) {
	// sourceIpv6Address := "2001:db8:1::"
	destinationIpv6Address := "2001:db8:2::"
	// intents := domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})
	nodeId := "to"
	nodeSid := "2001:db8:1::"
	serviceSid := "2001:db8:f::"
	tests := []struct {
		name      string
		pathIsNil bool
		wantErr   bool
	}{
		{
			name:      "TestCalculationTransformerService_TransformResult path is nil",
			pathIsNil: true,
			wantErr:   false,
		},
		{
			name:      "TestCalculationTransformerService_TransformResult path is not nil",
			pathIsNil: false,
			wantErr:   false,
		},
		{
			name:    "TestCalculationTransformerService_TransformResult path error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			cache := cache.NewMockCache(controller)
			service := NewCalculationTransformerService(cache)
			if tt.pathIsNil {
				pathRequest := domain.NewMockPathRequest(controller)
				pathRequest.EXPECT().GetIpv6DestinationAddress().Return(destinationIpv6Address)
				pathResult := service.TransformResult(nil, pathRequest, uint32(0))
				assert.NotNil(t, pathResult)
				assert.Equal(t, destinationIpv6Address, pathResult.GetIpv6SidAddresses()[0])
				return
			} else {
				path := graph.NewMockPath(controller)
				edge := graph.NewMockEdge(controller)
				path.EXPECT().GetEdges().Return([]graph.Edge{edge})
				to := graph.NewMockNode(controller)
				to.EXPECT().GetId().Return(nodeId)
				edge.EXPECT().To().Return(to)
				path.EXPECT().GetRouterServiceMap().Return(map[string]string{nodeId: serviceSid})
				cache.EXPECT().GetSrAlgorithmSid(nodeId, uint32(0)).Return(nodeSid)
				var pathResult domain.PathResult
				if !tt.wantErr {
					pathResult = service.TransformResult(path, domain.NewMockPathRequest(controller), uint32(0))
					assert.NotNil(t, pathResult)
					assert.Equal(t, []string{nodeSid, serviceSid}, pathResult.GetIpv6SidAddresses())
				} else {
					assert.Nil(t, service.TransformResult(path, nil, uint32(0)))
				}
			}
		})
	}
}
