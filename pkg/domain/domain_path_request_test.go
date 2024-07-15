package domain

import (
	"context"
	"reflect"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/api"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

func TestDomainPathRequest_validateMinMaxValues(t *testing.T) {
	tests := []struct {
		name       string
		intentType IntentType
		values     []Value
		wantErr    bool
	}{
		{
			name:       "Test validateMinMaxValues low latency with max value",
			intentType: IntentTypeLowLatency,
			values:     []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(10))},
			wantErr:    false,
		},
		{
			name:       "Test validateMinMaxValues low jitter with max value",
			intentType: IntentTypeLowJitter,
			values:     []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(10))},
			wantErr:    false,
		},
		{
			name:       "Test validateMinMaxValues low packet loss with max value",
			intentType: IntentTypeLowPacketLoss,
			values:     []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(10))},
			wantErr:    false,
		},
		{
			name:       "Test validateMinMaxValues high BW with min value",
			intentType: IntentTypeHighBandwidth,
			values:     []Value{getNumberValue(ValueTypeMinValue, proto.Int32(100000))},
			wantErr:    false,
		},
		{
			name:       "Test validateMinMaxValues error high BW with max value",
			intentType: IntentTypeHighBandwidth,
			values:     []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(100000))},
			wantErr:    true,
		},
		{
			name:       "Test validateMinMaxValues error low packet loss with min value",
			intentType: IntentTypeLowPacketLoss,
			values:     []Value{getNumberValue(ValueTypeMinValue, proto.Int32(100000))},
			wantErr:    true,
		},
		{
			name:       "Test validateMinMaxValues error intent type unspecified with min value",
			intentType: IntentTypeUnspecified,
			values:     []Value{getNumberValue(ValueTypeMinValue, proto.Int32(100000))},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMinMaxValues(tt.intentType, tt.values)
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}
		})
	}
}

func TestDomainPathRequest_validateDoubleAppearanceIntentType(t *testing.T) {
	tests := []struct {
		name        string
		intentType  IntentType
		intentTypes map[IntentType]bool
		wantErr     bool
	}{
		{
			name:        "Test validateDoubleAppearanceIntentType no appearance",
			intentType:  IntentTypeLowLatency,
			intentTypes: map[IntentType]bool{},
			wantErr:     false,
		},
		{
			name:        "Test validateDoubleAppearanceIntentType appearance",
			intentType:  IntentTypeLowLatency,
			intentTypes: map[IntentType]bool{IntentTypeLowLatency: true},
			wantErr:     true,
		},
		{
			name:        "Test validateDoubleAppearanceIntentType no appearance",
			intentType:  IntentTypeLowJitter,
			intentTypes: map[IntentType]bool{IntentTypeLowLatency: true},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDoubleAppearanceIntentType(tt.intentTypes, tt.intentType)
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}
		})
	}
}

func TestDomainPathRequest_validateFlexAlgoIntentType(t *testing.T) {
	tests := []struct {
		name       string
		intent     Intent
		intentType IntentType
		wantErr    bool
	}{
		{
			name:       "Test validateFlexAlgoIntentType no values",
			intent:     NewDomainIntent(IntentTypeFlexAlgo, []Value{}),
			intentType: IntentTypeFlexAlgo,
			wantErr:    true,
		},
		{
			name:       "Test validateFlexAlgoIntentType wrong value type",
			intent:     NewDomainIntent(IntentTypeFlexAlgo, []Value{getNumberValue(ValueTypeMinValue, proto.Int32(1))}),
			intentType: IntentTypeFlexAlgo,
			wantErr:    true,
		},
		{
			name:       "Test validateFlexAlgoIntentType wrong value (number less than 128)",
			intent:     NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeFlexAlgoNr, proto.Int32(1))}),
			intentType: IntentTypeFlexAlgo,
			wantErr:    true,
		},
		{
			name:       "Test validateFlexAlgoIntentType wrong value (number greater than 255)",
			intent:     NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeFlexAlgoNr, proto.Int32(256))}),
			intentType: IntentTypeFlexAlgo,
			wantErr:    true,
		},
		{
			name:       "Test validateFlexAlgoIntentType success",
			intent:     NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeFlexAlgoNr, proto.Int32(128))}),
			intentType: IntentTypeFlexAlgo,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFlexAlgoIntentType(tt.intent, tt.intentType)
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}
		})
	}
}

func TestDomainPathRequest_validateServiceFunctionChainIntentType(t *testing.T) {
	tests := []struct {
		name       string
		intent     Intent
		intentType IntentType
		wantErr    bool
	}{
		{
			name:       "Test validateServiceFunctionChainIntentType no values",
			intent:     NewDomainIntent(IntentTypeSFC, []Value{}),
			intentType: IntentTypeSFC,
			wantErr:    true,
		},
		{
			name:       "Test validateServiceFunctionChainIntentType wrong value type",
			intent:     NewDomainIntent(IntentTypeSFC, []Value{getNumberValue(ValueTypeMinValue, proto.Int32(1))}),
			intentType: IntentTypeSFC,
			wantErr:    true,
		},
		{
			name:       "Test validateServiceFunctionChainIntentType one correct value",
			intent:     NewDomainIntent(IntentTypeSFC, []Value{GetStringValue(ValueTypeSFC, proto.String("fw"))}),
			intentType: IntentTypeSFC,
			wantErr:    false,
		},
		{
			name:       "Test validateServiceFunctionChainIntentType twice the same value",
			intent:     NewDomainIntent(IntentTypeSFC, []Value{GetStringValue(ValueTypeSFC, proto.String("fw")), GetStringValue(ValueTypeSFC, proto.String("fw"))}),
			intentType: IntentTypeSFC,
			wantErr:    true,
		},
		{
			name:       "Test validateServiceFunctionChainIntentType two correct values",
			intent:     NewDomainIntent(IntentTypeSFC, []Value{GetStringValue(ValueTypeSFC, proto.String("fw")), GetStringValue(ValueTypeSFC, proto.String("ids"))}),
			intentType: IntentTypeSFC,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateServiceFunctionChainIntentType(tt.intent, tt.intentType)
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}
		})
	}
}

func TestDomainPathRequest_validateIntents(t *testing.T) {
	tests := []struct {
		name    string
		intents []Intent
		wantErr bool
	}{
		{
			name:    "Test validateIntents no intents",
			intents: []Intent{},
			wantErr: true,
		},
		{
			name: "Test validateIntents two same intents",
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{}),
				NewDomainIntent(IntentTypeLowLatency, []Value{}),
			},
			wantErr: true,
		},
		{
			name: "Test validateIntents two different intents",
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{}),
				NewDomainIntent(IntentTypeLowJitter, []Value{}),
			},
			wantErr: false,
		},
		{
			name: "Test validateIntents flex algo with wrong value",
			intents: []Intent{
				NewDomainIntent(IntentTypeFlexAlgo, []Value{getNumberValue(ValueTypeMinValue, proto.Int32(1))}),
			},
			wantErr: true,
		},
		{
			name: "Test validateIntents flex algo with correct value",
			intents: []Intent{
				NewDomainIntent(IntentTypeFlexAlgo, []Value{getNumberValue(ValueTypeFlexAlgoNr, proto.Int32(128))}),
			},
			wantErr: false,
		},
		{
			name: "Test validateIntents min max with wrong value",
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeMinValue, proto.Int32(10))}),
			},
			wantErr: true,
		},
		{
			name: "Test validateIntents flex algo with correct value",
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(10))}),
			},
			wantErr: false,
		},
		{
			name: "Test validateIntents sfc with wrong value",
			intents: []Intent{
				NewDomainIntent(IntentTypeSFC, []Value{getNumberValue(ValueTypeMinValue, proto.Int32(1))}),
			},
			wantErr: true,
		},
		{
			name: "Test validateIntents sfc with correct value",
			intents: []Intent{
				NewDomainIntent(IntentTypeSFC, []Value{GetStringValue(ValueTypeSFC, proto.String("fw"))}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateIntents(tt.intents)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateIntents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDomainPathRequest(t *testing.T) {
	tests := []struct {
		name                   string
		ipv6SourceAddress      string
		ipv6DestinationAddress string
		stream                 api.IntentController_GetIntentPathServer
		ctx                    context.Context
		intents                []Intent
		want                   *DomainPathRequest
		wantErr                bool
	}{
		{
			name:                   "Test NewDomainPathRequest no intents",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents:                []Intent{},
			want:                   nil,
			wantErr:                true,
		},
		{
			name:                   "Test NewDomainPathRequest two same intents",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{}),
				NewDomainIntent(IntentTypeLowLatency, []Value{}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:                   "Test NewDomainPathRequest two different intents no values",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{}),
				NewDomainIntent(IntentTypeLowPacketLoss, []Value{}),
			},
			want: &DomainPathRequest{
				ipv6SourceAddress:      "2001:db8::1",
				ipv6DestinationAddress: "2001:db8::2",
				intents: []Intent{
					NewDomainIntent(IntentTypeLowLatency, []Value{}),
					NewDomainIntent(IntentTypeLowPacketLoss, []Value{}),
				},
				stream: api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
				ctx:    context.Background(),
			},
			wantErr: false,
		},
		{
			name:                   "Test NewDomainPathRequest two different intents with values",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(10))}),
				NewDomainIntent(IntentTypeLowPacketLoss, []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(20))}),
			},
			want: &DomainPathRequest{
				ipv6SourceAddress:      "2001:db8::1",
				ipv6DestinationAddress: "2001:db8::2",
				intents: []Intent{
					NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(10))}),
					NewDomainIntent(IntentTypeLowPacketLoss, []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(20))}),
				},
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
				ctx:    context.Background(),
			},
			wantErr: false,
		},
		{
			name:                   "Test NewDomainPathRequest Flex Algo wrong value type ",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeFlexAlgo, []Value{getNumberValue(ValueTypeMinValue, proto.Int32(1))}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:                   "Test NewDomainPathRequest Flex Algo wrong value",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeFlexAlgo, []Value{getNumberValue(ValueTypeFlexAlgoNr, proto.Int32(1))}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:                   "Test NewDomainPathRequest Flex Algo wrong value ",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeFlexAlgo, []Value{getNumberValue(ValueTypeFlexAlgoNr, proto.Int32(128))}),
			},
			want: &DomainPathRequest{
				ipv6SourceAddress:      "2001:db8::1",
				ipv6DestinationAddress: "2001:db8::2",
				intents: []Intent{
					NewDomainIntent(IntentTypeFlexAlgo, []Value{getNumberValue(ValueTypeFlexAlgoNr, proto.Int32(128))}),
				},
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
				ctx:    context.Background(),
			},
			wantErr: false,
		},
		{
			name:                   "Test NewDomainPathRequest wrong min value",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeMinValue, proto.Int32(10))}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:                   "Test NewDomainPathRequest correct max value",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(10))}),
			},
			want: &DomainPathRequest{
				ipv6SourceAddress:      "2001:db8::1",
				ipv6DestinationAddress: "2001:db8::2",
				intents: []Intent{
					NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(10))}),
				},
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
				ctx:    context.Background(),
			},
			wantErr: false,
		},
		{
			name:                   "Test NewDomainPathRequest sfc wrong value",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeSFC, []Value{getNumberValue(ValueTypeMinValue, proto.Int32(10))}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:                   "Test NewDomainPathRequest sfc correct value",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeSFC, []Value{GetStringValue(ValueTypeSFC, proto.String("fw"))}),
			},
			want: &DomainPathRequest{
				ipv6SourceAddress:      "2001:db8::1",
				ipv6DestinationAddress: "2001:db8::2",
				intents: []Intent{
					NewDomainIntent(IntentTypeSFC, []Value{GetStringValue(ValueTypeSFC, proto.String("fw"))}),
				},
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
				ctx:    context.Background(),
			},
			wantErr: false,
		},
		{
			name:                   "Test NewDomainPathRequest complex example sfc, flex algo, low latency",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeSFC, []Value{GetStringValue(ValueTypeSFC, proto.String("fw")), GetStringValue(ValueTypeSFC, proto.String("ids"))}),
				NewDomainIntent(IntentTypeFlexAlgo, []Value{getNumberValue(ValueTypeFlexAlgoNr, proto.Int32(128))}),
				NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(10))}),
			},
			want: &DomainPathRequest{
				ipv6SourceAddress:      "2001:db8::1",
				ipv6DestinationAddress: "2001:db8::2",
				intents: []Intent{
					NewDomainIntent(IntentTypeSFC, []Value{GetStringValue(ValueTypeSFC, proto.String("fw")), GetStringValue(ValueTypeSFC, proto.String("ids"))}),
					NewDomainIntent(IntentTypeFlexAlgo, []Value{getNumberValue(ValueTypeFlexAlgoNr, proto.Int32(128))}),
					NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(10))}),
				},
				stream: api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
				ctx:    context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathRequest, err := NewDomainPathRequest(tt.ipv6SourceAddress, tt.ipv6DestinationAddress, tt.intents, tt.stream, tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDomainPathRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(pathRequest, tt.want) {
				t.Errorf("NewDomainPathRequest() = %v, want %v", pathRequest, tt.want)
			}
		})
	}
}

func TestDomainPathRequest_GetIpv6SourceAddress(t *testing.T) {
	tests := []struct {
		name                   string
		ipv6SourceAddress      string
		ipv6DestinationAddress string
		stream                 api.IntentController_GetIntentPathServer
		ctx                    context.Context
		intents                []Intent
	}{
		{
			name:                   "Test GetIpv6SourceAddress",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathRequest, err := NewDomainPathRequest(tt.ipv6SourceAddress, tt.ipv6DestinationAddress, tt.intents, tt.stream, tt.ctx)
			if err != nil {
				t.Error(err)
				return
			}
			ipv6SourceAddress := pathRequest.GetIpv6SourceAddress()
			if ipv6SourceAddress != tt.ipv6SourceAddress {
				t.Errorf("GetIpv6SourceAddress() = %v, want %v", ipv6SourceAddress, tt.ipv6SourceAddress)
			}
		})
	}
}

func TestDomainPathRequest_GetIpv6DestinationAddress(t *testing.T) {
	tests := []struct {
		name                   string
		ipv6SourceAddress      string
		ipv6DestinationAddress string
		stream                 api.IntentController_GetIntentPathServer
		ctx                    context.Context
		intents                []Intent
	}{
		{
			name:                   "Test GetIpv6DestinationAddress",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathRequest, err := NewDomainPathRequest(tt.ipv6SourceAddress, tt.ipv6DestinationAddress, tt.intents, tt.stream, tt.ctx)
			if err != nil {
				t.Error(err)
				return
			}
			ipv6DestinationAddress := pathRequest.GetIpv6DestinationAddress()
			if ipv6DestinationAddress != tt.ipv6DestinationAddress {
				t.Errorf("GetIpv6DestinationAddress() = %v, want %v", ipv6DestinationAddress, tt.ipv6DestinationAddress)
			}
		})
	}
}

func TestDomainPathRequest_GetIntents(t *testing.T) {
	tests := []struct {
		name                   string
		ipv6SourceAddress      string
		ipv6DestinationAddress string
		stream                 api.IntentController_GetIntentPathServer
		ctx                    context.Context
		intents                []Intent
	}{
		{
			name:                   "Test GetIntents single value",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{}),
			},
		},
		{
			name:                   "Test GetIntents multiple values",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(10))}),
				NewDomainIntent(IntentTypeLowPacketLoss, []Value{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathRequest, err := NewDomainPathRequest(tt.ipv6SourceAddress, tt.ipv6DestinationAddress, tt.intents, tt.stream, tt.ctx)
			if err != nil {
				t.Error(err)
				return
			}
			intents := pathRequest.GetIntents()
			if !reflect.DeepEqual(intents, tt.intents) {
				t.Errorf("GetIntents() = %v, want %v", intents, tt.intents)
			}
		})
	}
}

func TestDomainPathRequest_GetContext(t *testing.T) {
	tests := []struct {
		name                   string
		ipv6SourceAddress      string
		ipv6DestinationAddress string
		stream                 api.IntentController_GetIntentPathServer
		ctx                    context.Context
		intents                []Intent
	}{
		{
			name:                   "Test GetContext",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathRequest, err := NewDomainPathRequest(tt.ipv6SourceAddress, tt.ipv6DestinationAddress, tt.intents, tt.stream, tt.ctx)
			if err != nil {
				t.Error(err)
				return
			}
			ctx := pathRequest.GetContext()
			if ctx != tt.ctx {
				t.Errorf("GetContext() = %v, want %v", ctx, tt.ctx)
			}
		})
	}
}

func TestDomainPathRequest_GetStream(t *testing.T) {
	tests := []struct {
		name                   string
		ipv6SourceAddress      string
		ipv6DestinationAddress string
		stream                 api.IntentController_GetIntentPathServer
		ctx                    context.Context
		intents                []Intent
	}{
		{
			name:                   "Test GetStream",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathRequest, err := NewDomainPathRequest(tt.ipv6SourceAddress, tt.ipv6DestinationAddress, tt.intents, tt.stream, tt.ctx)
			if err != nil {
				t.Error(err)
				return
			}
			stream := pathRequest.GetStream()
			if stream != tt.stream {
				t.Errorf("GetStream() = %v, want %v", stream, tt.stream)
			}
		})
	}
}

func TestDomainPathRequest_Serialize(t *testing.T) {
	tests := []struct {
		name                   string
		ipv6SourceAddress      string
		ipv6DestinationAddress string
		stream                 api.IntentController_GetIntentPathServer
		ctx                    context.Context
		intents                []Intent
		want                   string
	}{
		{
			name:                   "Test Serialize single value",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{}),
			},
			want: "2001:db8::1,2001:db8::2,LowLatency",
		},
		{
			name:                   "Test Serialize multiple values",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(10))}),
				NewDomainIntent(IntentTypeLowPacketLoss, []Value{}),
			},
			want: "2001:db8::1,2001:db8::2,LowLatency,MaxValue:10,LowPacketLoss",
		},
		{
			name:                   "Test Serialize complex example sfc, flex algo, low latency",
			ipv6SourceAddress:      "2001:db8::1",
			ipv6DestinationAddress: "2001:db8::2",
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			intents: []Intent{
				NewDomainIntent(IntentTypeSFC, []Value{GetStringValue(ValueTypeSFC, proto.String("fw")), GetStringValue(ValueTypeSFC, proto.String("ids"))}),
				NewDomainIntent(IntentTypeFlexAlgo, []Value{getNumberValue(ValueTypeFlexAlgoNr, proto.Int32(128))}),
				NewDomainIntent(IntentTypeLowLatency, []Value{getNumberValue(ValueTypeMaxValue, proto.Int32(10))}),
			},
			want: "2001:db8::1,2001:db8::2,SFC,SFC:fw,SFC:ids,FlexAlgo,FlexAlgoNr:128,LowLatency,MaxValue:10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathRequest, err := NewDomainPathRequest(tt.ipv6SourceAddress, tt.ipv6DestinationAddress, tt.intents, tt.stream, tt.ctx)
			if err != nil {
				t.Error(err)
				return
			}
			serialization := pathRequest.Serialize()
			if serialization != tt.want {
				t.Errorf("Serialize() = %v, want %v", serialization, tt.want)
			}
		})
	}
}
