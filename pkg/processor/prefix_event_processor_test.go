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

func TestNewPrefixEventProcessor(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewPrefixEventProcessor",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			assert.NotNil(t, NewPrefixEventProcessor(graphMock, cacheMock))
		})
	}
}

func TestPrefixEventProcessor_clearDuplicateAnnouncedPrefix(t *testing.T) {
	tests := []struct {
		name         string
		existInCache bool
	}{
		{
			name:         "TestPrefixEventProcessor_clearDuplicateAnnouncedPrefix does not exist in cache",
			existInCache: false,
		},
		{
			name:         "TestPrefixEventProcessor_clearDuplicateAnnouncedPrefix exist in cache",
			existInCache: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewPrefixEventProcessor(graphMock, cacheMock)
			prefix, err := domain.NewDomainPrefix(proto.String("key"), proto.String("0000.0000.1"), proto.String("fc:00::"), proto.Int32(64))
			assert.NoError(t, err)
			if tt.existInCache {
				cacheMock.EXPECT().GetClientNetworkByKey(gomock.Any()).Return(prefix).AnyTimes()
			} else {
				cacheMock.EXPECT().GetClientNetworkByKey(gomock.Any()).Return(nil).AnyTimes()
				cacheMock.EXPECT().RemoveClientNetwork(gomock.Any()).Return().AnyTimes()
			}
			processor.clearDuplicateAnnouncedPrefix(prefix, "fc0:00::", 64)
		})
	}
}

func TestPrefixEventProcessor_addNetworkToCache(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestPrefixEventProcessor_addNetworkToCache",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewPrefixEventProcessor(graphMock, cacheMock)
			prefix, err := domain.NewDomainPrefix(proto.String("key"), proto.String("0000.0000.1"), proto.String("fc:00::"), proto.Int32(64))
			assert.NoError(t, err)
			cacheMock.EXPECT().StoreClientNetwork(prefix).Return().AnyTimes()
			processor.addNetworkToCache(prefix, "fc0:00::", 64)
			assert.Equal(t, 1, processor.prefixCounts["fc0:00::"])
		})
	}
}

func TestPrefixEventProcessor_updatePrefixCounts(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{
			name:  "TestPrefixEventProcessor_updatePrefixCounts count > 1",
			count: 2,
		},
		{
			name:  "TestPrefixEventProcessor_updatePrefixCounts count = 1",
			count: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewPrefixEventProcessor(graphMock, cacheMock)
			prefix, err := domain.NewDomainPrefix(proto.String("key"), proto.String("0000.0000.1"), proto.String("fc:00::"), proto.Int32(64))
			assert.NoError(t, err)
			processor.prefixCounts["fc0:00::"] = tt.count
			cacheMock.EXPECT().RemoveClientNetwork(gomock.Any()).Return().AnyTimes()
			processor.updatePrefixCounts(tt.count, "fc0:00::", 64, prefix)
			if tt.count > 1 {
				assert.Equal(t, tt.count-1, processor.prefixCounts["fc0:00::"])
				_, exists := processor.prefixCounts["fc0:00::"]
				assert.True(t, exists)
			} else {
				assert.Equal(t, 0, processor.prefixCounts["fc0:00::"])
				_, exists := processor.prefixCounts["fc0:00::"]
				assert.False(t, exists)
			}
		})
	}
}

func TestPrefixEventProcessor_deleteClientNetwork(t *testing.T) {
	key := proto.String("key")
	igpRouterId := proto.String("0000.0000.0001")
	networkAddress := "fc0:00::"
	subnetLength := uint8(64)

	tests := []struct {
		name               string
		existInCache       bool
		existInPrefixCount bool
	}{
		{
			name:               "TestPrefixEventProcessor_deleteClientNetwork does not exist in cache",
			existInCache:       false,
			existInPrefixCount: false,
		},
		{
			name:               "TestPrefixEventProcessor_deleteClientNetwork does not exist in prefix count",
			existInCache:       true,
			existInPrefixCount: false,
		},
		{
			name:               "TestPrefixEventProcessor_deleteClientNetwork exist in cache and prefix count",
			existInCache:       true,
			existInPrefixCount: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewPrefixEventProcessor(graphMock, cacheMock)
			prefix, err := domain.NewDomainPrefix(key, igpRouterId, proto.String(networkAddress), proto.Int32(int32(subnetLength)))
			assert.NoError(t, err)
			if tt.existInCache {
				cacheMock.EXPECT().GetClientNetworkByKey(gomock.Any()).Return(prefix).AnyTimes()
			} else {
				cacheMock.EXPECT().GetClientNetworkByKey(gomock.Any()).Return(nil).AnyTimes()
			}
			if tt.existInPrefixCount {
				processor.prefixCounts[networkAddress] = 1
			}
			if tt.existInCache && tt.existInPrefixCount {
				cacheMock.EXPECT().RemoveClientNetwork(gomock.Any()).Return().AnyTimes()
				assert.NoError(t, processor.deleteClientNetwork(*key))
			} else {
				assert.Error(t, processor.deleteClientNetwork(*key))
			}
		})
	}
}

func TestPrefixEventProcessor_processPrefix(t *testing.T) {
	key := proto.String("key")
	igpRouterId := proto.String("0000.0000.0001")
	networkAddress := "fc0:00::"
	subnetLength := uint8(64)
	tests := []struct {
		name               string
		existInPrefixCount bool
	}{
		{
			name:               "TestPrefixEventProcessor_processPrefix exist in prefix count",
			existInPrefixCount: true,
		},
		{
			name:               "TestPrefixEventProcessor_processPrefix does not exist in prefix count",
			existInPrefixCount: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewPrefixEventProcessor(graphMock, cacheMock)
			prefix, err := domain.NewDomainPrefix(key, igpRouterId, proto.String(networkAddress), proto.Int32(int32(subnetLength)))
			assert.NoError(t, err)
			if tt.existInPrefixCount {
				processor.prefixCounts[networkAddress] = 1
				cacheMock.EXPECT().GetClientNetworkByKey(gomock.Any()).Return(prefix).AnyTimes()
			} else {
				cacheMock.EXPECT().StoreClientNetwork(prefix).Return().AnyTimes()
			}
			processor.processPrefix(prefix)
		})
	}
}

func TestPrefixEventProcessor_ProcessPrefixes(t *testing.T) {
	key := proto.String("key")
	igpRouterId := proto.String("0000.0000.0001")
	networkAddress := "fc0:00::"
	subnetLength := uint8(64)
	tests := []struct {
		name string
	}{
		{
			name: "TestPrefixEventProcessor_ProcessPrefixes",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewPrefixEventProcessor(graphMock, cacheMock)
			prefix, err := domain.NewDomainPrefix(key, igpRouterId, proto.String(networkAddress), proto.Int32(int32(subnetLength)))
			assert.NoError(t, err)
			cacheMock.EXPECT().StoreClientNetwork(prefix).Return().AnyTimes()
			processor.ProcessPrefixes([]domain.Prefix{prefix})
		})
	}
}

func TestPrefixEventProcessor_HandleEvent(t *testing.T) {
	key := proto.String("key")
	igpRouterId := proto.String("0000.0000.0001")
	networkAddress := proto.String("fc0:00::")
	subnetLength := proto.Int32(64)
	tests := []struct {
		name      string
		eventType string
	}{
		{
			name:      "TestPrefixEventProcessor_HandleEvent add prefix ",
			eventType: "add",
		},
		{
			name:      "TestPrefixEventProcessor_HandleEvent delete prefix",
			eventType: "del",
		},
		{
			name: "TestPrefixEventProcessor_HandleEvent wrong event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphMock := graph.NewMockGraph(gomock.NewController(t))
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			processor := NewPrefixEventProcessor(graphMock, cacheMock)
			var event domain.NetworkEvent
			if tt.eventType == "add" {
				prefix, err := domain.NewDomainPrefix(key, igpRouterId, networkAddress, subnetLength)
				assert.NoError(t, err)
				event = domain.NewAddPrefixEvent(prefix)
				cacheMock.EXPECT().StoreClientNetwork(prefix).Return().AnyTimes()
			} else if tt.eventType == "del" {
				event = domain.NewDeletePrefixEvent(*key)
				cacheMock.EXPECT().GetClientNetworkByKey(gomock.Any()).Return(nil).AnyTimes()
			} else {
				event = domain.NewAddNodeEvent(nil)
			}
			assert.False(t, processor.HandleEvent(event))
		})
	}
}
