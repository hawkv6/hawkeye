package processor

import (
	"testing"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func TestNewNodeEventProcessor(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewNodeEventProcessor",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			assert.NotNil(t, NewNodeEventProcessor(graphMock, cacheMock))
		})
	}
}

func TestNodeEventProcessor_updateNode(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNodeEventProcessor_updateNode success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewNodeEventProcessor(graphMock, cacheMock)
			node := graph.NewMockNode(gomock.NewController(t))
			node.EXPECT().GetId().Return("1").AnyTimes()
			node.EXPECT().SetName("name").Return()
			node.EXPECT().SetFlexibleAlgorithms([]uint32{1, 2}).Return()
			processor.updateNode(node, "name", []uint32{1, 2})
		})
	}
}

func TestNodeEventProcessor_updateNodeInGraph(t *testing.T) {
	tests := []struct {
		name       string
		nodeExists bool
	}{
		{
			name:       "TestNodeEventProcessor_updateNodeInGraph node does not exist",
			nodeExists: false,
		},

		{
			name:       "TestNodeEventProcessor_updateNodeInGraph node exists",
			nodeExists: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewNodeEventProcessor(graphMock, cacheMock)
			node := graph.NewMockNode(gomock.NewController(t))
			if tt.nodeExists {
				graphMock.EXPECT().NodeExists("1").Return(true).AnyTimes()
				node.EXPECT().GetId().Return("1").AnyTimes()
				node.EXPECT().SetName("name").Return().AnyTimes()
				node.EXPECT().SetFlexibleAlgorithms([]uint32{1, 2}).Return().AnyTimes()
				graphMock.EXPECT().GetNode("1").Return(node).AnyTimes()
			} else {
				graphMock.EXPECT().NodeExists("1").Return(false).AnyTimes()
				graphMock.EXPECT().AddNode(gomock.Any()).Return(node).AnyTimes()
			}
			processor.updateNodeInGraph("1", "name", []uint32{1, 2})
		})
	}
}

func TestNodeEventProcessor_addNodeToCache(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNodeEventProcessor_addNodeToCache",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewNodeEventProcessor(graphMock, cacheMock)
			node, err := domain.NewDomainNode(proto.String("1"), proto.String("igp router id"), proto.String("name"), []uint32{1, 2})
			assert.NoError(t, err)
			cacheMock.EXPECT().StoreNode(node).Return()
			processor.addNodeToCache(node)
		})
	}
}

func TestNodeEventProcessor_removeNodeFromGraph(t *testing.T) {
	tests := []struct {
		name         string
		existInGraph bool
	}{
		{
			name:         "TestNodeEventProcessor_removeNodeFromGraph node exists",
			existInGraph: true,
		},
		{
			name:         "TestNodeEventProcessor_removeNodeFromGraph node does not exist",
			existInGraph: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewNodeEventProcessor(graphMock, cacheMock)
			node, err := domain.NewDomainNode(proto.String("1"), proto.String("igp router id"), proto.String("name"), []uint32{1, 2})
			assert.Nil(t, err)
			if tt.existInGraph {
				graphNode := graph.NewMockNode(gomock.NewController(t))
				graphMock.EXPECT().NodeExists(gomock.Any()).Return(true).AnyTimes()
				graphMock.EXPECT().GetNode(gomock.Any()).Return(graphNode).AnyTimes()
				graphMock.EXPECT().DeleteNode(gomock.Any()).Return()
			} else {
				graphMock.EXPECT().NodeExists(gomock.Any()).Return(false).AnyTimes()
			}
			processor.removeNodeFromGraph(node)
		})
	}
}

func TestNodeEventProcessor_deleteNode(t *testing.T) {
	tests := []struct {
		name          string
		existsInCache bool
	}{
		{
			name:          "TestNodeEventProcessor_deleteNode exists in cache",
			existsInCache: true,
		},
		{
			name:          "TestNodeEventProcessor_deleteNode does not exist in cache",
			existsInCache: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewNodeEventProcessor(graphMock, cacheMock)
			node, err := domain.NewDomainNode(proto.String("1"), proto.String("igp router id"), proto.String("name"), []uint32{1, 2})
			assert.NoError(t, err)
			if tt.existsInCache {
				cacheMock.EXPECT().GetNodeByKey(gomock.Any()).Return(node).AnyTimes()
				cacheMock.EXPECT().RemoveNode(gomock.Any()).Return()
				graphMock.EXPECT().NodeExists(gomock.Any()).Return(false).AnyTimes()
			} else {
				cacheMock.EXPECT().GetNodeByKey(gomock.Any()).Return(nil).AnyTimes()
			}
			processor.deleteNode("1")
		})
	}
}

func TestNodeEventProcessor_addOrUpdateNodeInGraphAndCache(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNodeEventProcessor_addOrUpdateNodeInGraphAndCache",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewNodeEventProcessor(graphMock, cacheMock)
			node, err := domain.NewDomainNode(proto.String("1"), proto.String("igp router id"), proto.String("name"), []uint32{1, 2})
			assert.NoError(t, err)
			cacheMock.EXPECT().StoreNode(node).Return().AnyTimes()
			graphMock.EXPECT().NodeExists(gomock.Any()).Return(false).AnyTimes()
			graphMock.EXPECT().AddNode(gomock.Any()).Return(nil).AnyTimes()
			processor.addOrUpdateNodeInGraphAndCache(node)
		})
	}
}

func TestNodeEventProcessor_ProcessNodes(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNodeEventProcessor_addOrUpdateNodeInGraphAndCache",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewNodeEventProcessor(graphMock, cacheMock)
			node, err := domain.NewDomainNode(proto.String("1"), proto.String("igp router id"), proto.String("name"), []uint32{1, 2})
			assert.NoError(t, err)
			cacheMock.EXPECT().StoreNode(node).Return().AnyTimes()
			graphMock.EXPECT().NodeExists(gomock.Any()).Return(false).AnyTimes()
			graphMock.EXPECT().AddNode(gomock.Any()).Return(nil).AnyTimes()
			processor.ProcessNodes([]domain.Node{node})
		})
	}
}

func TestNodeEventProcessor_HandleEvent(t *testing.T) {
	key := proto.String("1")
	igpRouterId := proto.String("0000.0000.0001")
	nodeName := proto.String("XR-1")
	algorithms := []uint32{1, 2}
	tests := []struct {
		name        string
		eventType   string
		want        bool
		key         *string
		igpRouterId *string
		nodeName    *string
		algorithms  []uint32
	}{
		{
			name:        "TestNodeEventProcessor_HandleEvent AddNodeEvent",
			eventType:   "add",
			want:        true,
			key:         key,
			igpRouterId: igpRouterId,
			nodeName:    nodeName,
			algorithms:  algorithms,
		},
		{
			name:        "TestNodeEventProcessor_HandleEvent UpdateNodeEvent",
			eventType:   "update",
			want:        true,
			key:         key,
			igpRouterId: igpRouterId,
			nodeName:    nodeName,
			algorithms:  algorithms,
		},
		{
			name:      "TestNodeEventProcessor_HandleEvent DeleteNodeEvent",
			eventType: "del",
			want:      true,
			key:       key,
		},
		{
			name:      "TestNodeEventProcessor_HandleEvent Unknown event type",
			eventType: "unknown",
			want:      false,
			key:       key,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewNodeEventProcessor(graphMock, cacheMock)
			var event domain.NetworkEvent
			if tt.eventType == "add" {
				node, err := domain.NewDomainNode(tt.key, tt.igpRouterId, tt.nodeName, tt.algorithms)
				assert.NoError(t, err)
				event = domain.NewAddNodeEvent(node)
				graphNode := graph.NewMockNode(gomock.NewController(t))
				graphMock.EXPECT().NodeExists(gomock.Any()).Return(false).AnyTimes()
				graphMock.EXPECT().AddNode(gomock.Any()).Return(graphNode).AnyTimes()
				cacheMock.EXPECT().StoreNode(gomock.Any()).AnyTimes()
			} else if tt.eventType == "update" {
				node, err := domain.NewDomainNode(tt.key, tt.igpRouterId, tt.nodeName, tt.algorithms)
				assert.NoError(t, err)
				event = domain.NewUpdateNodeEvent(node)
				graphNode := graph.NewMockNode(gomock.NewController(t))
				graphMock.EXPECT().NodeExists(gomock.Any()).Return(true).AnyTimes()
				graphMock.EXPECT().GetNode(gomock.Any()).Return(graphNode).AnyTimes()
				graphNode.EXPECT().GetId().Return(*tt.igpRouterId).AnyTimes()
				graphNode.EXPECT().SetName(gomock.Any()).Return().AnyTimes()
				graphNode.EXPECT().SetFlexibleAlgorithms(gomock.Any()).Return().AnyTimes()
				cacheMock.EXPECT().StoreNode(gomock.Any()).AnyTimes()
			} else if tt.eventType == "del" {
				event = domain.NewDeleteNodeEvent(*tt.key)
				cacheMock.EXPECT().GetNodeByKey(gomock.Any()).Return(nil).AnyTimes()
			} else {
				event = domain.NewDeleteLinkEvent(*tt.key)
			}
			assert.Equal(t, tt.want, processor.HandleEvent(event))
		})
	}
}
