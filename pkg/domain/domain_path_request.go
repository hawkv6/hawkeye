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

func validateMinMaxValues(intentType IntentType, values []Value) error {
	for _, value := range values {
		switch value.GetValueType() {
		case ValueTypeMaxValue:
			if intentType != IntentTypeLowLatency && intentType != IntentTypeLowJitter && intentType != IntentTypeLowPacketLoss {
				return fmt.Errorf("Max value is only allowed for Low Latency, Jitter, or Packet Loss intents")
			}
		case ValueTypeMinValue:
			if intentType != IntentTypeHighBandwidth {
				return fmt.Errorf("Min value is only allowed for High Bandwidth intents")
			}
		}
	}
	return nil
}

func validateDoubleAppearanceIntentType(intentTypes map[IntentType]bool, intentType IntentType) error {
	if _, exists := intentTypes[intentType]; exists {
		return fmt.Errorf("Intent type %v appears more than once", intentType)
	}
	return nil
}

func validateFlexAlgoIntentType(intent Intent, intentType IntentType) error {
	if intentType == IntentTypeFlexAlgo {
		values := intent.GetValues()
		if len(values) == 0 || len(values) > 1 {
			return fmt.Errorf("Flex Algo intent should have exact one VALUE_TYPE_FLEX_ALGO_NR")
		}
		valueType := values[0].GetValueType()
		if valueType != ValueTypeFlexAlgoNr {
			return fmt.Errorf("Flex Algo value number should be of type VALUE_TYPE_FLEX_ALGO_NR ")
		}
		value := values[0].GetNumberValue()
		if value < 128 || value > 255 {
			return fmt.Errorf("Flex Algo value number should be a positive number between 128 and 255")
		}
	}
	return nil
}

func validateServiceFunctionChainIntentType(intent Intent, intentType IntentType) error {
	if intentType == IntentTypeSFC {
		values := intent.GetValues()
		if len(values) == 0 {
			return fmt.Errorf("Service Function Chain intent should have at least one VALUE_TYPE_SERVICE_FUNCTION_CHAIN")
		}
		for _, value := range values {
			if value.GetValueType() != ValueTypeSFC {
				return fmt.Errorf("Service Function Chain value should be of type VALUE_TYPE_SERVICE_FUNCTION_CHAIN")
			}
		}
	}
	return nil
}

func validateIntents(intents []Intent) error {
	intentTypes := make(map[IntentType]bool)
	for _, intent := range intents {
		intentType := intent.GetIntentType()
		if err := validateDoubleAppearanceIntentType(intentTypes, intentType); err != nil {
			return err
		}
		if err := validateFlexAlgoIntentType(intent, intentType); err != nil {
			return err
		}
		if err := validateMinMaxValues(intentType, intent.GetValues()); err != nil {
			return err
		}
		if err := validateServiceFunctionChainIntentType(intent, intentType); err != nil {
			return err
		}
		intentTypes[intentType] = true
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
