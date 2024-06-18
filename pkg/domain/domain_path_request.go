package domain

import (
	"context"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/hawkv6/hawkeye/pkg/api"
)

type PathRequest interface {
	GetIpv6SourceAddress() string
	GetIpv6DestinationAddress() string
	GetIntents() []Intent
	GetContext() context.Context
	GetStream() api.IntentController_GetIntentPathServer
	Serialize() string
}

type DomainPathRequest struct {
	ipv6SourceAddress      string   `validate:"required,ipv6"`
	ipv6DestinationAddress string   `validate:"required,ipv6"`
	intents                []Intent `validate:"required,dive,uniqueIntentTypes"`
	stream                 api.IntentController_GetIntentPathServer
	ctx                    context.Context
}

func validateIntents(intents []Intent) error {
	intentTypes := make(map[IntentType]bool)
	for _, intent := range intents {
		if _, exists := intentTypes[intent.GetIntentType()]; exists {
			return fmt.Errorf("Intent type %v appears more than once", intent.GetIntentType())
		}
		intentTypes[intent.GetIntentType()] = true
	}
	return nil
}

func NewDomainPathRequest(ipv6SourceAddress string, ipv6DestinationAddress string, intents []Intent, stream api.IntentController_GetIntentPathServer, ctx context.Context) (*DomainPathRequest, error) {
	if err := validateIntents(intents); err != nil {
		return nil, err
	}

	pathRequest := &DomainPathRequest{
		ipv6SourceAddress:      ipv6SourceAddress,
		ipv6DestinationAddress: ipv6DestinationAddress,
		intents:                intents,
		stream:                 stream,
		ctx:                    ctx,
	}

	validator := validator.New()
	err := validator.Struct(pathRequest)
	if err != nil {
		return nil, err
	}
	return pathRequest, nil
}

func (pathRequest *DomainPathRequest) GetIpv6SourceAddress() string {
	return pathRequest.ipv6SourceAddress
}

func (pathRequest *DomainPathRequest) GetIpv6DestinationAddress() string {
	return pathRequest.ipv6DestinationAddress
}

func (pathRequest *DomainPathRequest) GetIntents() []Intent {
	return pathRequest.intents
}

func (pathRequest *DomainPathRequest) GetContext() context.Context {
	return pathRequest.ctx
}

func (pathRequest *DomainPathRequest) GetStream() api.IntentController_GetIntentPathServer {
	return pathRequest.stream
}

func (pathRequest *DomainPathRequest) Serialize() string {
	serialization := pathRequest.ipv6SourceAddress + "," + pathRequest.ipv6DestinationAddress + ","
	for i := 0; i < len(pathRequest.intents); i++ {
		if i == len(pathRequest.intents)-1 {
			serialization += pathRequest.intents[i].Serialize()
		} else {
			serialization += pathRequest.intents[i].Serialize() + ","
		}
	}
	return serialization
}
