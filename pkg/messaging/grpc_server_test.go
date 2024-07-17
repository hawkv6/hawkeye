package messaging

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/hawkv6/hawkeye/pkg/adapter"
	"github.com/hawkv6/hawkeye/pkg/api"
	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewGrpcMessagingServer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewGrpcMessagingServer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetGrpcPort().Return(uint16(10000))
			adapter := adapter.NewMockAdapter(gomock.NewController(t))
			channels := NewPathMessagingChannels()
			assert.NotNil(t, NewGrpcMessagingServer(adapter, config, channels))
		})
	}
}

func TestGrpcMessagingServer_Start(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestGrpcMessagingServer_Start success",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetGrpcPort().Return(uint16(10000)).AnyTimes()
			adapter := adapter.NewMockAdapter(gomock.NewController(t))
			channels := NewPathMessagingChannels()
			server := NewGrpcMessagingServer(adapter, config, channels)
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				_ = server.Start()
				wg.Done()
			}()
			time.Sleep(100 * time.Millisecond)
			server.stopChan <- struct{}{}
			wg.Wait()
		})
	}
}

func TestGrpcMessagingServer_GetIntentPath(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestGrpcMessagingServer_GetIntentPath no error",
			wantErr: false,
		},
		{
			name:    "TestGrpcMessagingServer_GetIntentPath error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetGrpcPort().Return(uint16(10000)).AnyTimes()
			adapter := adapter.NewMockAdapter(gomock.NewController(t))
			channels := NewPathMessagingChannels()
			server := NewGrpcMessagingServer(adapter, config, channels)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			stream := api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t))
			stream.EXPECT().Context().Return(ctx).AnyTimes()
			if tt.wantErr {
				stream.EXPECT().Recv().Return(nil, assert.AnError).AnyTimes()
			} else {
				cancel()
			}
			go func() {
				err := server.GetIntentPath(stream)
				if (err != nil) != tt.wantErr {
					t.Errorf("GrpcMessagingServer.GetIntentPath() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			time.Sleep(100 * time.Millisecond)
		})
	}
}

func TestGrpcMessagingServer_handleIncomingPathRequests(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestGrpcMessagingServer_handleIncomingPathRequests no error",
			wantErr: false,
		},
		{
			name:    "TestGrpcMessagingServer_handleIncomingPathRequests error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetGrpcPort().Return(uint16(10000)).AnyTimes()
			adapter := adapter.NewMockAdapter(gomock.NewController(t))
			channels := NewPathMessagingChannels()
			server := NewGrpcMessagingServer(adapter, config, channels)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			stream := api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t))
			stream.EXPECT().Context().Return(ctx).AnyTimes()
			if tt.wantErr {
				stream.EXPECT().Recv().Return(nil, assert.AnError).AnyTimes()
			} else {
				cancel()
			}
			go func() {
				server.handleIncomingPathRequests(stream, nil, ctx)
			}()
			time.Sleep(100 * time.Millisecond)
		})
	}
}

func TestGrpcMessagingServer_processStream(t *testing.T) {
	tests := []struct {
		name           string
		wantReceiveErr bool
		receiveErr     error
		wantConvertErr bool
		apiRequest     *api.PathRequest
	}{
		{
			name:           "TestGrpcMessagingServer_processStream no error",
			wantReceiveErr: false,
			wantConvertErr: false,
			apiRequest: &api.PathRequest{
				Ipv6SourceAddress:      "2001:db8::1",
				Ipv6DestinationAddress: "2001:db8::2",
				Intents: []*api.Intent{
					{
						Type: api.IntentType_INTENT_TYPE_LOW_LATENCY,
					},
				},
			},
		},
		{
			name:           "TestGrpcMessagingServer_processStream receive error EOF",
			wantReceiveErr: true,
			wantConvertErr: false,
			receiveErr:     io.EOF,
		},
		{
			name:           "TestGrpcMessagingServer_processStream error not EOF",
			wantReceiveErr: true,
			wantConvertErr: false,
			receiveErr:     assert.AnError,
		},
		{
			name:           "TestGrpcMessagingServer_processStream convert error",
			wantReceiveErr: false,
			wantConvertErr: true,
			apiRequest: &api.PathRequest{
				Ipv6SourceAddress:      "2001:db8::1",
				Ipv6DestinationAddress: "2001:db8::2",
				Intents:                []*api.Intent{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetGrpcPort().Return(uint16(10000)).AnyTimes()
			adapter := adapter.NewDomainAdapter()
			channels := NewPathMessagingChannels()
			server := NewGrpcMessagingServer(adapter, config, channels)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			stream := api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t))
			stream.EXPECT().Context().Return(ctx).AnyTimes()
			if tt.wantReceiveErr {
				stream.EXPECT().Recv().Return(nil, tt.receiveErr).AnyTimes()
			} else {
				stream.EXPECT().Recv().Return(tt.apiRequest, nil).AnyTimes()
				cancel()
			}
			go func() {
				err := server.processStream(stream, nil, ctx)
				if (err != nil) != (tt.wantReceiveErr || tt.wantConvertErr) {
					t.Errorf("GrpcMessagingServer.processStream() error = %v, wantReceiveErr %v, wantConvertErr %v", err, tt.wantReceiveErr, tt.wantConvertErr)
				}
			}()
			time.Sleep(100 * time.Millisecond)
		})
	}
}

func TestGrpcMessagingServer_processPathResult(t *testing.T) {
	tests := []struct {
		name           string
		wantConvertErr bool
		wantStreamErr  bool
	}{
		{
			name:           "TestGrpcMessagingServer_processPathResult no error",
			wantConvertErr: false,
			wantStreamErr:  false,
		},
		{
			name:           "TestGrpcMessagingServer_processPathResult convert err",
			wantConvertErr: true,
			wantStreamErr:  false,
		},
		{
			name:           "TestGrpcMessagingServer_processPathResult stream err",
			wantConvertErr: false,
			wantStreamErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetGrpcPort().Return(uint16(10000)).AnyTimes()
			adapter := adapter.NewMockAdapter(gomock.NewController(t))
			channels := NewPathMessagingChannels()
			server := NewGrpcMessagingServer(adapter, config, channels)
			stream := api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t))
			if !tt.wantConvertErr {
				adapter.EXPECT().ConvertPathResult(gomock.Any()).Return(&api.PathResult{}, nil).AnyTimes()
			} else {
				adapter.EXPECT().ConvertPathResult(gomock.Any()).Return(nil, assert.AnError).AnyTimes()
			}
			if !tt.wantStreamErr {
				stream.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
			} else {
				stream.EXPECT().Send(gomock.Any()).Return(assert.AnError).AnyTimes()
			}
			err := server.processPathResult(stream, nil)
			if (err != nil) != (tt.wantConvertErr || tt.wantStreamErr) {
				t.Errorf("GrpcMessagingServer.processPathResult() error = %v, wantConvertErr %v, wantStreamErr %v", err, tt.wantConvertErr, tt.wantStreamErr)
			}
		})
	}
}

func TestGrpcMessagingServer_handleIntentPathResponse(t *testing.T) {
	tests := []struct {
		name           string
		wantProcessErr bool
		wantOtherErr   bool
	}{
		{
			name:           "TestGrpcMessagingServer_handleIntentPathResponse success",
			wantProcessErr: false,
			wantOtherErr:   false,
		},
		{
			name:           "TestGrpcMessagingServer_handleIntentPathResponse process error",
			wantProcessErr: true,
			wantOtherErr:   false,
		},
		{
			name:           "TestGrpcMessagingServer_handleIntentPathResponse other error",
			wantProcessErr: false,
			wantOtherErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetGrpcPort().Return(uint16(10000)).AnyTimes()
			adapter := adapter.NewMockAdapter(gomock.NewController(t))
			channels := NewPathMessagingChannels()
			server := NewGrpcMessagingServer(adapter, config, channels)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			stream := api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t))
			stream.EXPECT().Context().Return(ctx).AnyTimes()
			go func() {
				server.handleIntentPathResponse(stream, ctx)
			}()
			if !tt.wantProcessErr && !tt.wantOtherErr {
				adapter.EXPECT().ConvertPathResult(gomock.Any()).Return(&api.PathResult{}, nil).AnyTimes()
				stream.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
				channels.GetPathResponseChan() <- nil
			} else if tt.wantProcessErr {
				adapter.EXPECT().ConvertPathResult(gomock.Any()).Return(nil, assert.AnError).AnyTimes()
				channels.GetErrorChan() <- assert.AnError
			} else if tt.wantOtherErr {
				channels.GetErrorChan() <- assert.AnError
			}
			time.Sleep(100 * time.Millisecond)
		})
	}
}

func TestGrpcMessagingServer_Stop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestGrpcMessagingServer_Stop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.NewMockConfig(gomock.NewController(t))
			config.EXPECT().GetGrpcPort().Return(uint16(10000)).AnyTimes()
			adapter := adapter.NewMockAdapter(gomock.NewController(t))
			channels := NewPathMessagingChannels()
			server := NewGrpcMessagingServer(adapter, config, channels)
			server.Stop()
		})
	}
}
