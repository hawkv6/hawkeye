package processor

import (
	"sync"
	"testing"
	"time"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func TestNewNetworkProcessor(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cache := cache.NewMockCache(gomock.NewController(t))
			nodeEventProcessor := NewNodeEventProcessor(graphMock, cache)
			linkEventProcessor := NewLinkEventProcessor(graphMock, cache)
			prefixEventProcessor := NewPrefixEventProcessor(graphMock, cache)
			sidEventProcessor := NewSidEventProcessor(graphMock, cache)
			eventDispatcher := NewEventDispatcher(nodeEventProcessor, linkEventProcessor, prefixEventProcessor, sidEventProcessor)
			eventOptions := EventOptions{
				NodeEventProcessor:   nodeEventProcessor,
				LinkEventProcessor:   linkEventProcessor,
				PrefixEventProcessor: prefixEventProcessor,
				SidEventProcessor:    sidEventProcessor,
				EventDispatcher:      eventDispatcher,
			}
			networkProcessor := NewNetworkProcessor(graphMock, cache, nil, nil, eventOptions)
			assert.NotNil(t, networkProcessor)
		})
	}
}

func TestNetworkProcessor_ProcessNodes(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNetworkProcessor_ProcessNodes",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cache := cache.NewMockCache(gomock.NewController(t))
			nodeEventProcessor := NewNodeEventProcessor(graphMock, cache)
			linkEventProcessor := NewLinkEventProcessor(graphMock, cache)
			prefixEventProcessor := NewPrefixEventProcessor(graphMock, cache)
			sidEventProcessor := NewSidEventProcessor(graphMock, cache)
			eventDispatcher := NewEventDispatcher(nodeEventProcessor, linkEventProcessor, prefixEventProcessor, sidEventProcessor)
			eventOptions := EventOptions{
				NodeEventProcessor:   nodeEventProcessor,
				LinkEventProcessor:   linkEventProcessor,
				PrefixEventProcessor: prefixEventProcessor,
				SidEventProcessor:    sidEventProcessor,
				EventDispatcher:      eventDispatcher,
			}
			networkProcessor := NewNetworkProcessor(graphMock, cache, nil, nil, eventOptions)
			networkProcessor.ProcessNodes(nil)
		})
	}
}

func TestNetworkProcessor_ProcessLinks(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNetworkProcessor_ProcessLinks",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cache := cache.NewMockCache(gomock.NewController(t))
			nodeEventProcessor := NewNodeEventProcessor(graphMock, cache)
			linkEventProcessor := NewLinkEventProcessor(graphMock, cache)
			prefixEventProcessor := NewPrefixEventProcessor(graphMock, cache)
			sidEventProcessor := NewSidEventProcessor(graphMock, cache)
			eventDispatcher := NewEventDispatcher(nodeEventProcessor, linkEventProcessor, prefixEventProcessor, sidEventProcessor)
			eventOptions := EventOptions{
				NodeEventProcessor:   nodeEventProcessor,
				LinkEventProcessor:   linkEventProcessor,
				PrefixEventProcessor: prefixEventProcessor,
				SidEventProcessor:    sidEventProcessor,
				EventDispatcher:      eventDispatcher,
			}
			networkProcessor := NewNetworkProcessor(graphMock, cache, nil, nil, eventOptions)
			assert.NoError(t, networkProcessor.ProcessLinks(nil))
		})
	}
}

func TestNetworkProcessor_ProcessPrefixes(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNetworkProcessor_ProcessPrefixes",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cache := cache.NewMockCache(gomock.NewController(t))
			nodeEventProcessor := NewNodeEventProcessor(graphMock, cache)
			linkEventProcessor := NewLinkEventProcessor(graphMock, cache)
			prefixEventProcessor := NewPrefixEventProcessor(graphMock, cache)
			sidEventProcessor := NewSidEventProcessor(graphMock, cache)
			eventDispatcher := NewEventDispatcher(nodeEventProcessor, linkEventProcessor, prefixEventProcessor, sidEventProcessor)
			eventOptions := EventOptions{
				NodeEventProcessor:   nodeEventProcessor,
				LinkEventProcessor:   linkEventProcessor,
				PrefixEventProcessor: prefixEventProcessor,
				SidEventProcessor:    sidEventProcessor,
				EventDispatcher:      eventDispatcher,
			}
			networkProcessor := NewNetworkProcessor(graphMock, cache, nil, nil, eventOptions)
			networkProcessor.ProcessPrefixes(nil)
		})
	}
}

func TestNetworkProcessor_ProcessSids(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNetworkProcessor_ProcessSids",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cache := cache.NewMockCache(gomock.NewController(t))
			nodeEventProcessor := NewNodeEventProcessor(graphMock, cache)
			linkEventProcessor := NewLinkEventProcessor(graphMock, cache)
			prefixEventProcessor := NewPrefixEventProcessor(graphMock, cache)
			sidEventProcessor := NewSidEventProcessor(graphMock, cache)
			eventDispatcher := NewEventDispatcher(nodeEventProcessor, linkEventProcessor, prefixEventProcessor, sidEventProcessor)
			eventOptions := EventOptions{
				NodeEventProcessor:   nodeEventProcessor,
				LinkEventProcessor:   linkEventProcessor,
				PrefixEventProcessor: prefixEventProcessor,
				SidEventProcessor:    sidEventProcessor,
				EventDispatcher:      eventDispatcher,
			}
			networkProcessor := NewNetworkProcessor(graphMock, cache, nil, nil, eventOptions)
			networkProcessor.ProcessSids(nil)
		})
	}
}

func TestNetworkProcessor_dispatchEvent(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNetworkProcessor_dispatchEvent set update false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			nodeEventProcessor := NewNodeEventProcessor(graphMock, cacheMock)
			linkEventProcessor := NewLinkEventProcessor(graphMock, cacheMock)
			prefixEventProcessor := NewPrefixEventProcessor(graphMock, cacheMock)
			sidEventProcessor := NewSidEventProcessor(graphMock, cacheMock)
			eventDispatcher := NewEventDispatcher(nodeEventProcessor, linkEventProcessor, prefixEventProcessor, sidEventProcessor)
			eventOptions := EventOptions{
				NodeEventProcessor:   nodeEventProcessor,
				LinkEventProcessor:   linkEventProcessor,
				PrefixEventProcessor: prefixEventProcessor,
				SidEventProcessor:    sidEventProcessor,
				EventDispatcher:      eventDispatcher,
			}
			networkProcessor := NewNetworkProcessor(graphMock, cacheMock, nil, nil, eventOptions)

			cacheMock.EXPECT().Lock().Return().AnyTimes()
			graphMock.EXPECT().Lock().Return().AnyTimes()
			deleteNodeEvent := domain.NewDeleteNodeEvent("node key")
			node, err := domain.NewDomainNode(proto.String("node key"), proto.String("igp router id"), proto.String("node name"), []uint32{})
			assert.Nil(t, err)
			cacheMock.EXPECT().GetNodeByKey(gomock.Any()).Return(node).AnyTimes()
			cacheMock.EXPECT().RemoveNode(gomock.Any()).Return().AnyTimes()
			graphMock.EXPECT().GetNode(gomock.Any()).Return(nil).AnyTimes()
			graphMock.EXPECT().NodeExists(gomock.Any()).Return(true).AnyTimes()
			graphMock.EXPECT().DeleteNode(gomock.Any()).Return().AnyTimes()

			holdTime := helper.NetworkProcessorHoldTime
			timer := time.NewTimer(holdTime)
			networkProcessor.dispatchEvent(deleteNodeEvent, timer, holdTime)
		})
	}
}

func TestNetworkProcessor_triggerUpdates(t *testing.T) {
	tests := []struct {
		name                string
		needsSubgraphUpdate bool
		mutexLocked         bool
	}{
		{
			name:                "TestNetworkProcessor_triggerUpdates no subgraph update",
			needsSubgraphUpdate: false,
			mutexLocked:         true,
		},
		{
			name:                "TestNetworkProcessor_triggerUpdates subgraph update",
			needsSubgraphUpdate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			nodeEventProcessor := NewNodeEventProcessor(graphMock, cacheMock)
			linkEventProcessor := NewLinkEventProcessor(graphMock, cacheMock)
			prefixEventProcessor := NewPrefixEventProcessor(graphMock, cacheMock)
			sidEventProcessor := NewSidEventProcessor(graphMock, cacheMock)
			eventDispatcher := NewEventDispatcher(nodeEventProcessor, linkEventProcessor, prefixEventProcessor, sidEventProcessor)
			eventOptions := EventOptions{
				NodeEventProcessor:   nodeEventProcessor,
				LinkEventProcessor:   linkEventProcessor,
				PrefixEventProcessor: prefixEventProcessor,
				SidEventProcessor:    sidEventProcessor,
				EventDispatcher:      eventDispatcher,
			}
			networkProcessor := NewNetworkProcessor(graphMock, cacheMock, nil, make(chan struct{}), eventOptions)
			if tt.mutexLocked {
				networkProcessor.mutexesLocked = true
				cacheMock.EXPECT().Unlock().Return().AnyTimes()
				graphMock.EXPECT().Unlock().Return().AnyTimes()
			}
			if tt.needsSubgraphUpdate {
				networkProcessor.needsSubgraphUpdate = true
				graphMock.EXPECT().UpdateSubGraphs().Return().AnyTimes()
			}
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				networkProcessor.triggerUpdates()
				wg.Done()
			}()

			<-networkProcessor.updateChan
			wg.Wait()

		})
	}
}

func TestNetworkProcessor_Start(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNetworkProcessor_Start",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			nodeEventProcessor := NewNodeEventProcessor(graphMock, cacheMock)
			linkEventProcessor := NewLinkEventProcessor(graphMock, cacheMock)
			prefixEventProcessor := NewPrefixEventProcessor(graphMock, cacheMock)
			sidEventProcessor := NewSidEventProcessor(graphMock, cacheMock)
			eventDispatcher := NewEventDispatcher(nodeEventProcessor, linkEventProcessor, prefixEventProcessor, sidEventProcessor)
			eventOptions := EventOptions{
				NodeEventProcessor:   nodeEventProcessor,
				LinkEventProcessor:   linkEventProcessor,
				PrefixEventProcessor: prefixEventProcessor,
				SidEventProcessor:    sidEventProcessor,
				EventDispatcher:      eventDispatcher,
			}
			eventChan := make(chan domain.NetworkEvent)
			networkProcessor := NewNetworkProcessor(graphMock, cacheMock, eventChan, make(chan struct{}), eventOptions)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			graphMock.EXPECT().Lock().Return().AnyTimes()
			graphMock.EXPECT().Unlock().Return().AnyTimes()
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				networkProcessor.Start()
				wg.Done()
			}()
			// trigerr dispatchEvent
			eventChan <- nil
			//trigger updates
			time.Sleep(1*time.Second + helper.NetworkProcessorHoldTime)
			<-networkProcessor.updateChan
			networkProcessor.quitChan <- struct{}{}
			wg.Wait()
		})
	}
}

func TestNetworkProcessor_Stop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNetworkProcessor_Stop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			nodeEventProcessor := NewNodeEventProcessor(graphMock, cacheMock)
			linkEventProcessor := NewLinkEventProcessor(graphMock, cacheMock)
			prefixEventProcessor := NewPrefixEventProcessor(graphMock, cacheMock)
			sidEventProcessor := NewSidEventProcessor(graphMock, cacheMock)
			eventDispatcher := NewEventDispatcher(nodeEventProcessor, linkEventProcessor, prefixEventProcessor, sidEventProcessor)
			eventOptions := EventOptions{
				NodeEventProcessor:   nodeEventProcessor,
				LinkEventProcessor:   linkEventProcessor,
				PrefixEventProcessor: prefixEventProcessor,
				SidEventProcessor:    sidEventProcessor,
				EventDispatcher:      eventDispatcher,
			}
			eventChan := make(chan domain.NetworkEvent)
			networkProcessor := NewNetworkProcessor(graphMock, cacheMock, eventChan, make(chan struct{}), eventOptions)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			graphMock.EXPECT().Lock().Return().AnyTimes()
			graphMock.EXPECT().Unlock().Return().AnyTimes()
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				networkProcessor.Start()
				wg.Done()
			}()
			time.Sleep(100 * time.Millisecond)
			networkProcessor.Stop()
			wg.Wait()
		})
	}
}
