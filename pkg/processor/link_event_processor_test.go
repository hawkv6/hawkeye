package processor

import (
	"fmt"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func TestNewLinkEventProcessor(t *testing.T) {
	tests := []struct {
		name  string
		graph graph.Graph
		cache cache.Cache
	}{
		{
			name:  "TestNewLinkEventProcessor",
			graph: graph.NewNetworkGraph(),
			cache: cache.NewInMemoryCache(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			processor := NewLinkEventProcessor(tt.graph, tt.cache)
			assert.NotNil(t, processor)
		})
	}
}

func TestLinkEventProcessor_getCurrentLinkWeights(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
	}{
		{
			name:                           "TestLinkEventProcessor_getCurrentLinkWeights",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewLinkEventProcessor(graph.NewNetworkGraph(), cache.NewInMemoryCache())
			link, err := domain.NewDomainLink(tt.key, tt.igpRouterId, tt.remoteIgpRouterId, tt.igpMetric, tt.unidirLinkDelay, tt.unidirDelayVariation, tt.maxLinkBWKbps, tt.unidirAvailableBandwidth, tt.unidirBandwidthUtilization, tt.unidirPacketLoss, tt.normalizedUnidirLinkDelay, tt.normalizedUnidirDelayVariation, tt.normalizedUnidirPacketLoss)
			assert.NoError(t, err)
			weights := processor.getCurrentLinkWeights(link)
			assert.Len(t, weights, 10)
			assert.Equal(t, weights[helper.IgpMetricKey], float64(10))
			assert.Equal(t, weights[helper.LatencyKey], float64(2000))
			assert.Equal(t, weights[helper.JitterKey], float64(100))
			assert.Equal(t, weights[helper.MaximumLinkBandwidthKey], float64(1000000))
			assert.Equal(t, weights[helper.AvailableBandwidthKey], float64(99766))
			assert.Equal(t, weights[helper.UtilizedBandwidthKey], float64(234))
			assert.Equal(t, weights[helper.PacketLossKey], float64(3.0059316283477027))
			assert.Equal(t, weights[helper.NormalizedLatencyKey], float64(0.05))
			assert.Equal(t, weights[helper.NormalizedJitterKey], float64(0.016452169298129225))
			assert.Equal(t, weights[helper.NormalizedPacketLossKey], float64(1e-10))
		})
	}
}

func TestLinkEventProcessor_deleteEdge(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		edgeExists bool
	}{
		{
			name: "TestLinkEventProcessor_deleteEdge edge exists",
			key:  "2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6",
		},
		{
			name: "TestLinkEventProcessor_deleteEdge edge does not exist",
			key:  "2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			processor := NewLinkEventProcessor(graphMock, cache.NewInMemoryCache())
			if tt.edgeExists {
				graphMock.EXPECT().EdgeExists(tt.key).Return(true)
				edge := graph.NewMockEdge(gomock.NewController(t))
				graphMock.EXPECT().GetEdge(tt.key).Return(edge)
				from := graph.NewMockNode(gomock.NewController(t))
				to := graph.NewMockNode(gomock.NewController(t))
				edge.EXPECT().From().Return(from)
				edge.EXPECT().To().Return(to)
				from.EXPECT().GetName().Return("from")
				to.EXPECT().GetName().Return("to")
				graphMock.EXPECT().DeleteEdge(edge).Return()
			} else {
				graphMock.EXPECT().EdgeExists(tt.key).Return(false)
			}
			processor.deleteEdge(tt.key)
		})
	}
}

func TestLinkEventProcessor_addEdgeToGraph(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestLinkEventProcessor_addEdgeToGraph success",
			wantErr: false,
		},
		{
			name:    "TestLinkEventProcessor_addEdgeToGraph error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			processor := NewLinkEventProcessor(graphMock, cache.NewInMemoryCache())
			edge := graph.NewMockEdge(gomock.NewController(t))
			edge.EXPECT().GetId().Return("edge id").AnyTimes()
			from := graph.NewMockNode(gomock.NewController(t))
			to := graph.NewMockNode(gomock.NewController(t))
			edge.EXPECT().From().Return(from).AnyTimes()
			edge.EXPECT().To().Return(to).AnyTimes()
			from.EXPECT().GetName().Return("from").AnyTimes()
			to.EXPECT().GetName().Return("to").AnyTimes()
			edge.EXPECT().GetAllWeights().Return(map[helper.WeightKey]float64{})
			if !tt.wantErr {
				graphMock.EXPECT().AddEdge(edge).Return(nil).AnyTimes()
			} else {
				graphMock.EXPECT().AddEdge(edge).Return(fmt.Errorf("Edge with id %s already exists", edge.GetId())).AnyTimes()
			}
			err := processor.addEdgeToGraph(edge)
			if (err != nil) != tt.wantErr {
				t.Errorf("addEdgeToGraph() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLinkEventProcessor_getOrCreateNode(t *testing.T) {
	nodeId := "0000.0000.000b"
	tests := []struct {
		name              string
		nodeExistsInGraph bool
		nodeExistsInCache bool
	}{
		{
			name:              "TestLinkEventProcessor_getOrCreateNode node exists",
			nodeExistsInGraph: true,
			nodeExistsInCache: true,
		},
		{
			name:              "TestLinkEventProcessor_getOrCreateNode node does not exist in graph but in cache",
			nodeExistsInGraph: false,
			nodeExistsInCache: true,
		},
		{
			name:              "TestLinkEventProcessor_getOrCreateNode node does not exist in graph nor in cache",
			nodeExistsInGraph: false,
			nodeExistsInCache: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewLinkEventProcessor(graphMock, cacheMock)
			node := graph.NewMockNode(gomock.NewController(t))
			if tt.nodeExistsInGraph {
				graphMock.EXPECT().NodeExists(nodeId).Return(true).AnyTimes()
				graphMock.EXPECT().GetNode(nodeId).Return(node).AnyTimes()
			} else {
				graphMock.EXPECT().NodeExists(nodeId).Return(false).AnyTimes()
				if tt.nodeExistsInCache {
					domainNode, _ := domain.NewDomainNode(proto.String("key"), proto.String("igp router id"), proto.String("name"), []uint32{})
					cacheMock.EXPECT().GetNodeByIgpRouterId(nodeId).Return(domainNode).AnyTimes()
					graphMock.EXPECT().AddNode(gomock.Any()).Return(node).AnyTimes()
				} else {
					cacheMock.EXPECT().GetNodeByIgpRouterId(nodeId).Return(nil).AnyTimes()
					graphMock.EXPECT().AddNode(gomock.Any()).Return(node).AnyTimes()
				}
			}
			returnNode := processor.getOrCreateNode(nodeId)
			assert.NotNil(t, returnNode)
		})
	}
}

func TestLinkEventProcessor_addLinkToGraph(t *testing.T) {
	key := "2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"
	igpRouterId := "0000.0000.000b"
	remoteIgpRouterId := "0000.0000.0006"
	igpMetric := uint32(10)
	unidirLinkDelay := uint32(2000)
	undirDelayVariation := uint32(100)
	maxLinkBWKbps := uint64(1000000)
	unidirAvailableBandwidth := uint32(99766)
	undirBandwidthUtilization := uint32(234)
	unidirPacketLoss := float64(3.0059316283477027)
	normalizedUnidirLinkDelay := float64(0.05)
	normalizedUnidirDelayVariation := float64(0.016452169298129225)
	normalizedUnidirPacketLoss := float64(1e-10)
	tests := []struct {
		name                           string
		edgeExistsInGraph              bool
		wantErr                        bool
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
	}{
		{
			name:                           "TestLinkEventProcessor_addLinkToGraph edge already exists",
			edgeExistsInGraph:              true,
			wantErr:                        false,
			key:                            proto.String(key),
			igpRouterId:                    proto.String(igpRouterId),
			remoteIgpRouterId:              proto.String(remoteIgpRouterId),
			igpMetric:                      proto.Uint32(igpMetric),
			unidirLinkDelay:                proto.Uint32(unidirLinkDelay),
			unidirDelayVariation:           proto.Uint32(undirDelayVariation),
			maxLinkBWKbps:                  proto.Uint64(maxLinkBWKbps),
			unidirAvailableBandwidth:       proto.Uint32(unidirAvailableBandwidth),
			unidirBandwidthUtilization:     proto.Uint32(undirBandwidthUtilization),
			unidirPacketLoss:               proto.Float64(unidirPacketLoss),
			normalizedUnidirLinkDelay:      proto.Float64(normalizedUnidirLinkDelay),
			normalizedUnidirDelayVariation: proto.Float64(normalizedUnidirDelayVariation),
			normalizedUnidirPacketLoss:     proto.Float64(normalizedUnidirPacketLoss),
		},
		{
			name:                           "TestLinkEventProcessor_addLinkToGraph does not exist, create successfully",
			edgeExistsInGraph:              false,
			wantErr:                        false,
			key:                            proto.String(key),
			igpRouterId:                    proto.String(igpRouterId),
			remoteIgpRouterId:              proto.String(remoteIgpRouterId),
			igpMetric:                      proto.Uint32(igpMetric),
			unidirLinkDelay:                proto.Uint32(unidirLinkDelay),
			unidirDelayVariation:           proto.Uint32(undirDelayVariation),
			maxLinkBWKbps:                  proto.Uint64(maxLinkBWKbps),
			unidirAvailableBandwidth:       proto.Uint32(unidirAvailableBandwidth),
			unidirBandwidthUtilization:     proto.Uint32(undirBandwidthUtilization),
			unidirPacketLoss:               proto.Float64(unidirPacketLoss),
			normalizedUnidirLinkDelay:      proto.Float64(normalizedUnidirLinkDelay),
			normalizedUnidirDelayVariation: proto.Float64(normalizedUnidirDelayVariation),
			normalizedUnidirPacketLoss:     proto.Float64(normalizedUnidirPacketLoss),
		},
		{
			name:                           "TestLinkEventProcessor_addLinkToGraph does not exist, create successfully",
			edgeExistsInGraph:              false,
			wantErr:                        true,
			key:                            proto.String(key),
			igpRouterId:                    proto.String(igpRouterId),
			remoteIgpRouterId:              proto.String(remoteIgpRouterId),
			igpMetric:                      proto.Uint32(0),
			unidirLinkDelay:                proto.Uint32(unidirLinkDelay),
			unidirDelayVariation:           proto.Uint32(undirDelayVariation),
			maxLinkBWKbps:                  proto.Uint64(maxLinkBWKbps),
			unidirAvailableBandwidth:       proto.Uint32(unidirAvailableBandwidth),
			unidirBandwidthUtilization:     proto.Uint32(undirBandwidthUtilization),
			unidirPacketLoss:               proto.Float64(unidirPacketLoss),
			normalizedUnidirLinkDelay:      proto.Float64(normalizedUnidirLinkDelay),
			normalizedUnidirDelayVariation: proto.Float64(normalizedUnidirDelayVariation),
			normalizedUnidirPacketLoss:     proto.Float64(normalizedUnidirPacketLoss),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			processor := NewLinkEventProcessor(graphMock, cache.NewInMemoryCache())
			link, err := domain.NewDomainLink(tt.key, tt.igpRouterId, tt.remoteIgpRouterId, tt.igpMetric, tt.unidirLinkDelay, tt.unidirDelayVariation, tt.maxLinkBWKbps, tt.unidirAvailableBandwidth, tt.unidirBandwidthUtilization, tt.unidirPacketLoss, tt.normalizedUnidirLinkDelay, tt.normalizedUnidirDelayVariation, tt.normalizedUnidirPacketLoss)
			assert.Nil(t, err)
			if tt.edgeExistsInGraph {
				graphMock.EXPECT().EdgeExists(gomock.Any()).Return(true).AnyTimes()
			} else {
				graphMock.EXPECT().EdgeExists(gomock.Any()).Return(false).AnyTimes()
				if !tt.wantErr {
					graphMock.EXPECT().NodeExists(igpRouterId).Return(true).AnyTimes()
					graphMock.EXPECT().NodeExists(remoteIgpRouterId).Return(true).AnyTimes()
					from := graph.NewMockNode(gomock.NewController(t))
					to := graph.NewMockNode(gomock.NewController(t))
					graphMock.EXPECT().GetNode(igpRouterId).Return(from).AnyTimes()
					graphMock.EXPECT().GetNode(remoteIgpRouterId).Return(to).AnyTimes()
					edge := graph.NewMockEdge(gomock.NewController(t))
					edge.EXPECT().GetId().Return(key).AnyTimes().AnyTimes()
					edge.EXPECT().From().Return(from).AnyTimes().AnyTimes()
					edge.EXPECT().To().Return(to).AnyTimes()
					from.EXPECT().GetName().Return("from").AnyTimes()
					to.EXPECT().GetName().Return("to").AnyTimes()
					graphMock.EXPECT().AddEdge(gomock.Any()).Return(nil).AnyTimes()
					from.EXPECT().GetFlexibleAlgorithms().Return(map[uint32]struct{}{}).AnyTimes()
					to.EXPECT().GetFlexibleAlgorithms().Return(map[uint32]struct{}{}).AnyTimes()
				} else {
					graphMock.EXPECT().NodeExists(igpRouterId).Return(false).AnyTimes()
					graphMock.EXPECT().NodeExists(remoteIgpRouterId).Return(false).AnyTimes()
				}
			}
			err = processor.addLinkToGraph(link)
			if (err != nil) != tt.wantErr {
				t.Errorf("addLinkToGraph() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLinkEventProcessor_setEdgeWeight(t *testing.T) {
	tests := []struct {
		name         string
		value        float64
		currentValue float64
		wantErr      bool
	}{
		{
			name:         "TestLinkEventProcessor_setEdgeWeight value is 0",
			value:        0,
			currentValue: 1,
			wantErr:      true,
		},
		{
			name:         "TestLinkEventProcessor_setEdgeWeight value is not 0",
			value:        1,
			currentValue: 1,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			processor := NewLinkEventProcessor(graphMock, cache.NewInMemoryCache())
			edge := graph.NewMockEdge(gomock.NewController(t))
			key := helper.LatencyKey
			edge.EXPECT().GetWeight(key).Return(tt.currentValue).AnyTimes()
			if tt.wantErr {
				err := processor.setEdgeWeight(edge, key, tt.value)
				assert.Error(t, err)
			} else {
				edge.EXPECT().SetWeight(key, tt.value).AnyTimes()
				err := processor.setEdgeWeight(edge, key, tt.value)
				assert.NoError(t, err)
			}
		})
	}
}

func TestLinkEventProcessor_updateLinkInGraph(t *testing.T) {
	key := "2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"
	igpRouterId := "0000.0000.000b"
	remoteIgpRouterId := "0000.0000.0006"
	igpMetric := uint32(10)
	unidirLinkDelay := uint32(2000)
	undirDelayVariation := uint32(100)
	maxLinkBWKbps := uint64(1000000)
	unidirAvailableBandwidth := uint32(99766)
	undirBandwidthUtilization := uint32(234)
	unidirPacketLoss := float64(3.0059316283477027)
	normalizedUnidirLinkDelay := float64(0.05)
	normalizedUnidirDelayVariation := float64(0.016452169298129225)
	normalizedUnidirPacketLoss := float64(1e-10)
	tests := []struct {
		name                           string
		edgeExistsInGraph              bool
		wantErr                        bool
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
	}{
		{
			name:                           "TestLinkEventProcessor_updateLinkInGraph edge already exists no error in setting weights",
			edgeExistsInGraph:              true,
			wantErr:                        false,
			key:                            proto.String(key),
			igpRouterId:                    proto.String(igpRouterId),
			remoteIgpRouterId:              proto.String(remoteIgpRouterId),
			igpMetric:                      proto.Uint32(igpMetric),
			unidirLinkDelay:                proto.Uint32(unidirLinkDelay),
			unidirDelayVariation:           proto.Uint32(undirDelayVariation),
			maxLinkBWKbps:                  proto.Uint64(maxLinkBWKbps),
			unidirAvailableBandwidth:       proto.Uint32(unidirAvailableBandwidth),
			unidirBandwidthUtilization:     proto.Uint32(undirBandwidthUtilization),
			unidirPacketLoss:               proto.Float64(unidirPacketLoss),
			normalizedUnidirLinkDelay:      proto.Float64(normalizedUnidirLinkDelay),
			normalizedUnidirDelayVariation: proto.Float64(normalizedUnidirDelayVariation),
			normalizedUnidirPacketLoss:     proto.Float64(normalizedUnidirPacketLoss),
		},
		{
			name:                           "TestLinkEventProcessor_updateLinkInGraph edge already exists error in setting weights",
			edgeExistsInGraph:              true,
			wantErr:                        true,
			key:                            proto.String(key),
			igpRouterId:                    proto.String(igpRouterId),
			remoteIgpRouterId:              proto.String(remoteIgpRouterId),
			igpMetric:                      proto.Uint32(0),
			unidirLinkDelay:                proto.Uint32(unidirLinkDelay),
			unidirDelayVariation:           proto.Uint32(undirDelayVariation),
			maxLinkBWKbps:                  proto.Uint64(maxLinkBWKbps),
			unidirAvailableBandwidth:       proto.Uint32(unidirAvailableBandwidth),
			unidirBandwidthUtilization:     proto.Uint32(undirBandwidthUtilization),
			unidirPacketLoss:               proto.Float64(unidirPacketLoss),
			normalizedUnidirLinkDelay:      proto.Float64(normalizedUnidirLinkDelay),
			normalizedUnidirDelayVariation: proto.Float64(normalizedUnidirDelayVariation),
			normalizedUnidirPacketLoss:     proto.Float64(normalizedUnidirPacketLoss),
		},
		{
			name:                           "TestLinkEventProcessor_updateLinkInGraph edge does not exists no error in creating edge",
			edgeExistsInGraph:              false,
			wantErr:                        false,
			key:                            proto.String(key),
			igpRouterId:                    proto.String(igpRouterId),
			remoteIgpRouterId:              proto.String(remoteIgpRouterId),
			igpMetric:                      proto.Uint32(0),
			unidirLinkDelay:                proto.Uint32(unidirLinkDelay),
			unidirDelayVariation:           proto.Uint32(undirDelayVariation),
			maxLinkBWKbps:                  proto.Uint64(maxLinkBWKbps),
			unidirAvailableBandwidth:       proto.Uint32(unidirAvailableBandwidth),
			unidirBandwidthUtilization:     proto.Uint32(undirBandwidthUtilization),
			unidirPacketLoss:               proto.Float64(unidirPacketLoss),
			normalizedUnidirLinkDelay:      proto.Float64(normalizedUnidirLinkDelay),
			normalizedUnidirDelayVariation: proto.Float64(normalizedUnidirDelayVariation),
			normalizedUnidirPacketLoss:     proto.Float64(normalizedUnidirPacketLoss),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			processor := NewLinkEventProcessor(graphMock, cache.NewInMemoryCache())
			link, err := domain.NewDomainLink(tt.key, tt.igpRouterId, tt.remoteIgpRouterId, tt.igpMetric, tt.unidirLinkDelay, tt.unidirDelayVariation, tt.maxLinkBWKbps, tt.unidirAvailableBandwidth, tt.unidirBandwidthUtilization, tt.unidirPacketLoss, tt.normalizedUnidirLinkDelay, tt.normalizedUnidirDelayVariation, tt.normalizedUnidirPacketLoss)
			assert.NoError(t, err)
			edgeMock := graph.NewMockEdge(gomock.NewController(t))
			if tt.edgeExistsInGraph {
				graphMock.EXPECT().GetEdge(gomock.Any()).Return(edgeMock).AnyTimes()
				edgeMock.EXPECT().GetWeight(gomock.Any()).Return(float64(1)).AnyTimes()
				edgeMock.EXPECT().SetWeight(gomock.Any(), gomock.Any()).Return().AnyTimes()

			} else {
				graphMock.EXPECT().GetEdge(gomock.Any()).Return(nil).AnyTimes()
				graphMock.EXPECT().EdgeExists(gomock.Any()).Return(true).AnyTimes()
			}
			err = processor.updateLinkInGraph(link)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateLinkInGraph() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLinkEventProcessor_ProcessLinks(t *testing.T) {
	key := "2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"
	igpRouterId := "0000.0000.000b"
	remoteIgpRouterId := "0000.0000.0006"
	igpMetric := uint32(10)
	unidirLinkDelay := uint32(2000)
	undirDelayVariation := uint32(100)
	maxLinkBWKbps := uint64(1000000)
	unidirAvailableBandwidth := uint32(99766)
	undirBandwidthUtilization := uint32(234)
	unidirPacketLoss := float64(3.0059316283477027)
	normalizedUnidirLinkDelay := float64(0.05)
	normalizedUnidirDelayVariation := float64(0.016452169298129225)
	normalizedUnidirPacketLoss := float64(1e-10)
	tests := []struct {
		name                           string
		edgeExistsInGraph              bool
		wantErr                        bool
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
	}{
		{
			name:                           "TestLinkEventProcessor_ProcessLinks edge already exists",
			edgeExistsInGraph:              true,
			wantErr:                        false,
			key:                            proto.String(key),
			igpRouterId:                    proto.String(igpRouterId),
			remoteIgpRouterId:              proto.String(remoteIgpRouterId),
			igpMetric:                      proto.Uint32(igpMetric),
			unidirLinkDelay:                proto.Uint32(unidirLinkDelay),
			unidirDelayVariation:           proto.Uint32(undirDelayVariation),
			maxLinkBWKbps:                  proto.Uint64(maxLinkBWKbps),
			unidirAvailableBandwidth:       proto.Uint32(unidirAvailableBandwidth),
			unidirBandwidthUtilization:     proto.Uint32(undirBandwidthUtilization),
			unidirPacketLoss:               proto.Float64(unidirPacketLoss),
			normalizedUnidirLinkDelay:      proto.Float64(normalizedUnidirLinkDelay),
			normalizedUnidirDelayVariation: proto.Float64(normalizedUnidirDelayVariation),
			normalizedUnidirPacketLoss:     proto.Float64(normalizedUnidirPacketLoss),
		},
		{
			name:                           "TestLinkEventProcessor_ProcessLinks edge does not exists return error",
			edgeExistsInGraph:              false,
			wantErr:                        true,
			key:                            proto.String(key),
			igpRouterId:                    proto.String(igpRouterId),
			remoteIgpRouterId:              proto.String(remoteIgpRouterId),
			igpMetric:                      proto.Uint32(0),
			unidirLinkDelay:                proto.Uint32(unidirLinkDelay),
			unidirDelayVariation:           proto.Uint32(undirDelayVariation),
			maxLinkBWKbps:                  proto.Uint64(maxLinkBWKbps),
			unidirAvailableBandwidth:       proto.Uint32(unidirAvailableBandwidth),
			unidirBandwidthUtilization:     proto.Uint32(undirBandwidthUtilization),
			unidirPacketLoss:               proto.Float64(unidirPacketLoss),
			normalizedUnidirLinkDelay:      proto.Float64(normalizedUnidirLinkDelay),
			normalizedUnidirDelayVariation: proto.Float64(normalizedUnidirDelayVariation),
			normalizedUnidirPacketLoss:     proto.Float64(normalizedUnidirPacketLoss),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			processor := NewLinkEventProcessor(graphMock, cache.NewInMemoryCache())
			link, err := domain.NewDomainLink(tt.key, tt.igpRouterId, tt.remoteIgpRouterId, tt.igpMetric, tt.unidirLinkDelay, tt.unidirDelayVariation, tt.maxLinkBWKbps, tt.unidirAvailableBandwidth, tt.unidirBandwidthUtilization, tt.unidirPacketLoss, tt.normalizedUnidirLinkDelay, tt.normalizedUnidirDelayVariation, tt.normalizedUnidirPacketLoss)
			assert.NoError(t, err)
			if tt.edgeExistsInGraph {
				graphMock.EXPECT().GetEdge(gomock.Any()).Return(nil).AnyTimes()
				graphMock.EXPECT().EdgeExists(gomock.Any()).Return(true).AnyTimes()
			} else {
				graphMock.EXPECT().GetEdge(gomock.Any()).Return(nil).AnyTimes()
				graphMock.EXPECT().EdgeExists(gomock.Any()).Return(false).AnyTimes()
			}
			graphMock.EXPECT().UpdateSubGraphs().Return().AnyTimes()
			err = processor.ProcessLinks([]domain.Link{link})
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessLinks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLinkEventProcessor_handleAddLinkEvent(t *testing.T) {
	key := "2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"
	igpRouterId := "0000.0000.000b"
	remoteIgpRouterId := "0000.0000.0006"
	igpMetric := uint32(10)
	unidirLinkDelay := uint32(2000)
	undirDelayVariation := uint32(100)
	maxLinkBWKbps := uint64(1000000)
	unidirAvailableBandwidth := uint32(99766)
	undirBandwidthUtilization := uint32(234)
	unidirPacketLoss := float64(3.0059316283477027)
	normalizedUnidirLinkDelay := float64(0.05)
	normalizedUnidirDelayVariation := float64(0.016452169298129225)
	normalizedUnidirPacketLoss := float64(1e-10)
	tests := []struct {
		name                           string
		want                           bool
		existInGraph                   bool
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
	}{
		{
			name:                           "TestLinkEventProcessor_handleAddLinkEvent success",
			want:                           true,
			existInGraph:                   true,
			key:                            proto.String(key),
			igpRouterId:                    proto.String(igpRouterId),
			remoteIgpRouterId:              proto.String(remoteIgpRouterId),
			igpMetric:                      proto.Uint32(igpMetric),
			unidirLinkDelay:                proto.Uint32(unidirLinkDelay),
			unidirDelayVariation:           proto.Uint32(undirDelayVariation),
			maxLinkBWKbps:                  proto.Uint64(maxLinkBWKbps),
			unidirAvailableBandwidth:       proto.Uint32(unidirAvailableBandwidth),
			unidirBandwidthUtilization:     proto.Uint32(undirBandwidthUtilization),
			unidirPacketLoss:               proto.Float64(unidirPacketLoss),
			normalizedUnidirLinkDelay:      proto.Float64(normalizedUnidirLinkDelay),
			normalizedUnidirDelayVariation: proto.Float64(normalizedUnidirDelayVariation),
			normalizedUnidirPacketLoss:     proto.Float64(normalizedUnidirPacketLoss),
		},
		{
			name:                           "TestLinkEventProcessor_handleAddLinkEvent error",
			want:                           false,
			existInGraph:                   false,
			key:                            proto.String(key),
			igpRouterId:                    proto.String(igpRouterId),
			remoteIgpRouterId:              proto.String(remoteIgpRouterId),
			igpMetric:                      proto.Uint32(0),
			unidirLinkDelay:                proto.Uint32(unidirLinkDelay),
			unidirDelayVariation:           proto.Uint32(undirDelayVariation),
			maxLinkBWKbps:                  proto.Uint64(maxLinkBWKbps),
			unidirAvailableBandwidth:       proto.Uint32(unidirAvailableBandwidth),
			unidirBandwidthUtilization:     proto.Uint32(undirBandwidthUtilization),
			unidirPacketLoss:               proto.Float64(unidirPacketLoss),
			normalizedUnidirLinkDelay:      proto.Float64(normalizedUnidirLinkDelay),
			normalizedUnidirDelayVariation: proto.Float64(normalizedUnidirDelayVariation),
			normalizedUnidirPacketLoss:     proto.Float64(normalizedUnidirPacketLoss),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			processor := NewLinkEventProcessor(graphMock, cache.NewInMemoryCache())
			link, err := domain.NewDomainLink(tt.key, tt.igpRouterId, tt.remoteIgpRouterId, tt.igpMetric, tt.unidirLinkDelay, tt.unidirDelayVariation, tt.maxLinkBWKbps, tt.unidirAvailableBandwidth, tt.unidirBandwidthUtilization, tt.unidirPacketLoss, tt.normalizedUnidirLinkDelay, tt.normalizedUnidirDelayVariation, tt.normalizedUnidirPacketLoss)
			assert.NoError(t, err)
			event := domain.NewAddLinkEvent(link)
			graphMock.EXPECT().GetEdge(gomock.Any()).Return(nil).AnyTimes()
			if tt.existInGraph {
				graphMock.EXPECT().EdgeExists(gomock.Any()).Return(true).AnyTimes()
			} else {
				graphMock.EXPECT().EdgeExists(gomock.Any()).Return(false).AnyTimes()
			}
			needsUpdate := processor.handleAddLinkEvent(event)
			if needsUpdate != tt.want {
				t.Errorf("HandleEvent() = %v, want %v", needsUpdate, tt.want)
			}
		})
	}
}

func TestLinkEventProcessor_handleUpdateLinkEvent(t *testing.T) {
	key := "2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"
	igpRouterId := "0000.0000.000b"
	remoteIgpRouterId := "0000.0000.0006"
	igpMetric := uint32(10)
	unidirLinkDelay := uint32(2000)
	undirDelayVariation := uint32(100)
	maxLinkBWKbps := uint64(1000000)
	unidirAvailableBandwidth := uint32(99766)
	undirBandwidthUtilization := uint32(234)
	unidirPacketLoss := float64(3.0059316283477027)
	normalizedUnidirLinkDelay := float64(0.05)
	normalizedUnidirDelayVariation := float64(0.016452169298129225)
	normalizedUnidirPacketLoss := float64(1e-10)
	tests := []struct {
		name                           string
		want                           bool
		existInGraph                   bool
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
	}{
		{
			name:                           "TestLinkEventProcessor_handleUpdateLinkEvent success",
			want:                           false,
			existInGraph:                   true,
			key:                            proto.String(key),
			igpRouterId:                    proto.String(igpRouterId),
			remoteIgpRouterId:              proto.String(remoteIgpRouterId),
			igpMetric:                      proto.Uint32(igpMetric),
			unidirLinkDelay:                proto.Uint32(unidirLinkDelay),
			unidirDelayVariation:           proto.Uint32(undirDelayVariation),
			maxLinkBWKbps:                  proto.Uint64(maxLinkBWKbps),
			unidirAvailableBandwidth:       proto.Uint32(unidirAvailableBandwidth),
			unidirBandwidthUtilization:     proto.Uint32(undirBandwidthUtilization),
			unidirPacketLoss:               proto.Float64(unidirPacketLoss),
			normalizedUnidirLinkDelay:      proto.Float64(normalizedUnidirLinkDelay),
			normalizedUnidirDelayVariation: proto.Float64(normalizedUnidirDelayVariation),
			normalizedUnidirPacketLoss:     proto.Float64(normalizedUnidirPacketLoss),
		},
		{
			name:                           "TestLinkEventProcessor_handleUpdateLinkEvent error",
			want:                           false,
			existInGraph:                   false,
			key:                            proto.String(key),
			igpRouterId:                    proto.String(igpRouterId),
			remoteIgpRouterId:              proto.String(remoteIgpRouterId),
			igpMetric:                      proto.Uint32(0),
			unidirLinkDelay:                proto.Uint32(unidirLinkDelay),
			unidirDelayVariation:           proto.Uint32(undirDelayVariation),
			maxLinkBWKbps:                  proto.Uint64(maxLinkBWKbps),
			unidirAvailableBandwidth:       proto.Uint32(unidirAvailableBandwidth),
			unidirBandwidthUtilization:     proto.Uint32(undirBandwidthUtilization),
			unidirPacketLoss:               proto.Float64(unidirPacketLoss),
			normalizedUnidirLinkDelay:      proto.Float64(normalizedUnidirLinkDelay),
			normalizedUnidirDelayVariation: proto.Float64(normalizedUnidirDelayVariation),
			normalizedUnidirPacketLoss:     proto.Float64(normalizedUnidirPacketLoss),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			processor := NewLinkEventProcessor(graphMock, cache.NewInMemoryCache())
			link, err := domain.NewDomainLink(tt.key, tt.igpRouterId, tt.remoteIgpRouterId, tt.igpMetric, tt.unidirLinkDelay, tt.unidirDelayVariation, tt.maxLinkBWKbps, tt.unidirAvailableBandwidth, tt.unidirBandwidthUtilization, tt.unidirPacketLoss, tt.normalizedUnidirLinkDelay, tt.normalizedUnidirDelayVariation, tt.normalizedUnidirPacketLoss)
			assert.NoError(t, err)
			event := domain.NewUpdateLinkEvent(link)
			graphMock.EXPECT().GetEdge(gomock.Any()).Return(nil).AnyTimes()
			if tt.existInGraph {
				graphMock.EXPECT().EdgeExists(gomock.Any()).Return(true).AnyTimes()
			} else {
				graphMock.EXPECT().EdgeExists(gomock.Any()).Return(false).AnyTimes()
			}
			needsUpdate := processor.handleUpdateLinkEvent(event)
			if needsUpdate != tt.want {
				t.Errorf("HandleEvent() = %v, want %v", needsUpdate, tt.want)
			}
		})
	}
}

func TestLinkEventProcessor_handleDeleteLinkEvent(t *testing.T) {
	tests := []struct {
		name          string
		want          bool
		existsInGraph bool
	}{
		{
			name:          "TestLinkEventProcessor_handleDeleteLinkEvent success",
			want:          true,
			existsInGraph: true,
		},
		{
			name:          "TestLinkEventProcessor_handleDeleteLinkEvent does not exist in graph",
			want:          false,
			existsInGraph: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			processor := NewLinkEventProcessor(graphMock, cache.NewInMemoryCache())
			event := domain.NewDeleteLinkEvent("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6")
			if tt.existsInGraph {
				graphMock.EXPECT().EdgeExists(gomock.Any()).Return(true).AnyTimes()
				edge := graph.NewMockEdge(gomock.NewController(t))
				from := graph.NewMockNode(gomock.NewController(t))
				to := graph.NewMockNode(gomock.NewController(t))
				graphMock.EXPECT().GetEdge(gomock.Any()).Return(edge).AnyTimes()
				edge.EXPECT().From().Return(from).AnyTimes()
				edge.EXPECT().To().Return(to).AnyTimes()
				from.EXPECT().GetName().Return("from").AnyTimes()
				to.EXPECT().GetName().Return("to").AnyTimes()
				graphMock.EXPECT().DeleteEdge(gomock.Any()).Return().AnyTimes()
			} else {
				graphMock.EXPECT().EdgeExists(gomock.Any()).Return(false).AnyTimes()
			}
			needsUpdate := processor.handleDeleteLinkEvent(event)
			if needsUpdate != tt.want {
				t.Errorf("HandleEvent() = %v, want %v", needsUpdate, tt.want)
			}
		})
	}
}

func TestLinkEventProcessor_HandleEvent(t *testing.T) {
	key := "2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"
	igpRouterId := "0000.0000.000b"
	remoteIgpRouterId := "0000.0000.0006"
	igpMetric := uint32(10)
	unidirLinkDelay := uint32(2000)
	undirDelayVariation := uint32(100)
	maxLinkBWKbps := uint64(1000000)
	unidirAvailableBandwidth := uint32(99766)
	undirBandwidthUtilization := uint32(234)
	unidirPacketLoss := float64(3.0059316283477027)
	normalizedUnidirLinkDelay := float64(0.05)
	normalizedUnidirDelayVariation := float64(0.016452169298129225)
	normalizedUnidirPacketLoss := float64(1e-10)
	tests := []struct {
		name      string
		eventType string
		want      bool
	}{
		{
			name:      "TestLinkEventProcessor_HandleEvent AddLinkEvent",
			eventType: "add",
			want:      true,
		},
		{
			name:      "TestLinkEventProcessor_HandleEvent UpdateLinkEvent",
			eventType: "update",
			want:      true,
		},
		{
			name:      "TestLinkEventProcessor_HandleEvent DeleteLinkEvent",
			eventType: "del",
			want:      true,
		},
		{
			name: "TestLinkEventProcessor_HandleEvent DeleteNodeEvent",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			processor := NewLinkEventProcessor(graphMock, cache.NewInMemoryCache())
			var event domain.NetworkEvent
			link, err := domain.NewDomainLink(&key, &igpRouterId, &remoteIgpRouterId, &igpMetric, &unidirLinkDelay, &undirDelayVariation, &maxLinkBWKbps, &unidirAvailableBandwidth, &undirBandwidthUtilization, &unidirPacketLoss, &normalizedUnidirLinkDelay, &normalizedUnidirDelayVariation, &normalizedUnidirPacketLoss)
			assert.NoError(t, err)
			if tt.eventType == "add" {
				event = domain.NewAddLinkEvent(link)
				graphMock.EXPECT().EdgeExists(gomock.Any()).Return(true).AnyTimes()
			} else if tt.eventType == "update" {
				event = domain.NewUpdateLinkEvent(link)
				edge := graph.NewMockEdge(gomock.NewController(t))
				graphMock.EXPECT().GetEdge(gomock.Any()).Return(edge).AnyTimes()
				edge.EXPECT().GetWeight(gomock.Any()).Return(float64(1)).AnyTimes()
				edge.EXPECT().SetWeight(gomock.Any(), gomock.Any()).Return().AnyTimes()
			} else if tt.eventType == "del" {
				event = domain.NewDeleteLinkEvent("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6")
				graphMock.EXPECT().EdgeExists(gomock.Any()).Return(false).AnyTimes()
			} else {
				event = domain.NewDeleteNodeEvent("0000.0000.000b")
			}

			needsUpdate := processor.HandleEvent(event)
			if needsUpdate != tt.want {
				t.Errorf("HandleEvent() = %v, want %v", needsUpdate, tt.want)
			}
		})
	}
}
