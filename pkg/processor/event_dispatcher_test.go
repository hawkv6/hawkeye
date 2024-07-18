package processor

import (
	"testing"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/stretchr/testify/assert"
)

func TestNewEventDispatcher(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewEventDispatcher",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			nodeEventHandler := &NodeEventProcessor{}
			linkEventHandler := &LinkEventProcessor{}
			prefixEventHandler := &PrefixEventProcessor{}
			sidEventHandler := &SidEventProcessor{}
			dispatcher := NewEventDispatcher(nodeEventHandler, linkEventHandler, prefixEventHandler, sidEventHandler)
			assert.NotNil(t, dispatcher)
			assert.Len(t, dispatcher.eventHandlers, 10)
		})
	}
}

func TestEventDispatcher_Dispatch(t *testing.T) {
	tests := []struct {
		name  string
		event domain.NetworkEvent
		want  bool
	}{
		{
			name:  "TestEventDispatcher_Dispatch success",
			event: domain.NewDeleteNodeEvent("key"),
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := graph.NewNetworkGraph()
			cache := cache.NewInMemoryCache()
			nodeEventHandler := NewNodeEventProcessor(graph, cache)
			dispatcher := NewEventDispatcher(nodeEventHandler, nil, nil, nil)
			success := dispatcher.Dispatch(tt.event)
			assert.Equal(t, tt.want, success)
		})
	}
}
