package domain

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/hawkv6/hawkeye/pkg/api"
)

type DefaultPathRequest struct {
	ipv6SourceAddress      string   `validate:"required,ipv6"`
	ipv6DestinationAddress string   `validate:"required,ipv6"`
	intents                []Intent `validate:"required,dive"`
	stream                 api.IntentController_GetIntentPathServer
	ctx                    context.Context
}

func NewDefaultPathRequest(ipv6SourceAddress string, ipv6DestinationAddress string, intents []Intent, stream api.IntentController_GetIntentPathServer, ctx context.Context) (*DefaultPathRequest, error) {
	defaultPathRequest := &DefaultPathRequest{
		ipv6SourceAddress:      ipv6SourceAddress,
		ipv6DestinationAddress: ipv6DestinationAddress,
		intents:                intents,
		stream:                 stream,
		ctx:                    ctx,
	}
	validator := validator.New()
	err := validator.Struct(defaultPathRequest)
	if err != nil {
		return nil, err
	}
	return defaultPathRequest, nil
}

func (defaultPathRequest *DefaultPathRequest) GetIpv6SourceAddress() string {
	return defaultPathRequest.ipv6SourceAddress
}

func (defaultPathRequest *DefaultPathRequest) GetIpv6DestinationAddress() string {
	return defaultPathRequest.ipv6DestinationAddress
}

func (defaultPathRequest *DefaultPathRequest) GetIntents() []Intent {
	return defaultPathRequest.intents
}

func (defaultPathRequest *DefaultPathRequest) GetContext() context.Context {
	return defaultPathRequest.ctx
}

func (defaultPathRequest *DefaultPathRequest) GetStream() api.IntentController_GetIntentPathServer {
	return defaultPathRequest.stream
}

func (defaultPathRequest *DefaultPathRequest) Serialize() string {
	serialization := defaultPathRequest.ipv6SourceAddress + "," + defaultPathRequest.ipv6DestinationAddress + ","
	for i := 0; i < len(defaultPathRequest.intents); i++ {
		if i == len(defaultPathRequest.intents)-1 {
			serialization += defaultPathRequest.intents[i].Serialize()
		} else {
			serialization += defaultPathRequest.intents[i].Serialize() + ","
		}
	}
	return serialization
}
