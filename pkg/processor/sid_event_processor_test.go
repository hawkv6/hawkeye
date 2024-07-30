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

func TestNewSidEventProcessor(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewSidEventProcessor",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			assert.NotNil(t, NewSidEventProcessor(graphMock, cacheMock))
		})
	}
}

func TestSidEventProcessor_addSidtoCache(t *testing.T) {
	key := proto.String("0_0000.0000.000b_fc00:0:b:0:1::")
	igpRouterId := proto.String("0000.0000.000b")
	sid := proto.String("fc00:0:b:0:1::")
	algorithm := proto.Uint32(0)
	tests := []struct {
		name string
	}{
		{
			name: "TestSidEventProcessor_addSidtoCache",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewSidEventProcessor(graphMock, cacheMock)
			sid, err := domain.NewDomainSid(key, igpRouterId, sid, algorithm)
			assert.NoError(t, err)
			cacheMock.EXPECT().StoreSid(sid).Return().AnyTimes()
			processor.addSidtoCache(sid)
		})
	}
}

func TestSidEventProcessor_deleteSidFromCache(t *testing.T) {
	key := proto.String("0_0000.0000.000b_fc00:0:b:0:1::")
	tests := []struct {
		name         string
		existInCache bool
	}{
		{
			name:         "TestSidEventProcessor_deleteSidFromCache exist in cache",
			existInCache: true,
		},
		{
			name:         "TestSidEventProcessor_deleteSidFromCache does not exist in cache",
			existInCache: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewSidEventProcessor(graphMock, cacheMock)
			sid, err := domain.NewDomainSid(key, proto.String("0000.0000.000b"), proto.String("fc00:0:b:0:1::"), proto.Uint32(0))
			assert.NoError(t, err)
			if tt.existInCache {
				cacheMock.EXPECT().GetSidByKey(gomock.Any()).Return(sid).AnyTimes()
				cacheMock.EXPECT().RemoveSid(sid).Return()
			} else {
				cacheMock.EXPECT().GetSidByKey(gomock.Any()).Return(nil).AnyTimes()
			}
			processor.deleteSidFromCache(*key)
		})
	}
}

func TestSidEventProcessor_ProcessSids(t *testing.T) {
	key := proto.String("0_0000.0000.000b_fc00:0:b:0:1::")
	igpRouterId := proto.String("0000.0000.000b")
	sid := proto.String("fc00:0:b:0:1::")
	algorithm := proto.Uint32(0)
	tests := []struct {
		name string
	}{
		{
			name: "TestSidEventProcessor_ProcessSids",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewSidEventProcessor(graphMock, cacheMock)
			sid, err := domain.NewDomainSid(key, igpRouterId, sid, algorithm)
			assert.NoError(t, err)
			cacheMock.EXPECT().StoreSid(sid).Return().AnyTimes()
			processor.ProcessSids([]domain.Sid{sid})
		})
	}
}

func TestSidEventProcessor_HandleEvent(t *testing.T) {
	key := proto.String("0_0000.0000.000b_fc00:0:b:0:1::")
	igpRouterId := proto.String("0000.0000.000b")
	sid := proto.String("fc00:0:b:0:1::")
	algorithm := proto.Uint32(0)
	tests := []struct {
		name      string
		eventType string
	}{
		{
			name:      "TestSidEventProcessor_HandleEvent AddSidEvent",
			eventType: "add",
		},
		{
			name:      "TestSidEventProcessor_HandleEvent DeleteSidEvent",
			eventType: "delete",
		},
		{
			name:      "TestSidEventProcessor_HandleEvent No handler found",
			eventType: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewSidEventProcessor(graphMock, cacheMock)
			sid, err := domain.NewDomainSid(key, igpRouterId, sid, algorithm)
			var event domain.NetworkEvent
			assert.NoError(t, err)
			if tt.eventType == "add" {
				event = domain.NewAddSidEvent(sid)
				cacheMock.EXPECT().StoreSid(sid).Return().AnyTimes()
			} else if tt.eventType == "delete" {
				cacheMock.EXPECT().GetSidByKey(gomock.Any()).Return(sid).AnyTimes()
				cacheMock.EXPECT().RemoveSid(sid).Return().AnyTimes()
				event = domain.NewDeleteSidEvent(*key)
			} else {
				event = domain.NewDeleteNodeEvent("")
			}
			assert.False(t, processor.HandleEvent(event))
		})
	}
}
