package adapter

import (
	"context"
	"fmt"

	"github.com/hawkv6/hawkeye/pkg/api"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/jalapeno-api-gateway/jagw-go/jagw"
	"github.com/sirupsen/logrus"
)

type DefaultAdapter struct {
	log *logrus.Entry
}

func NewDefaultAdapter() *DefaultAdapter {
	return &DefaultAdapter{
		log: logging.DefaultLogger.WithField("subsystem", Subsystem),
	}
}

func (adapter *DefaultAdapter) ConvertNode(lsNode *jagw.LsNode) (domain.Node, error) {
	return domain.NewDefaultNode(lsNode.Key, lsNode.IgpRouterId, lsNode.Name)
}

func (adapter *DefaultAdapter) ConvertLink(lsLink *jagw.LsLink) (domain.Link, error) {
	return domain.NewDefaultLink(lsLink.Key, lsLink.IgpRouterId, lsLink.RemoteIgpRouterId, lsLink.UnidirLinkDelay, lsLink.UnidirDelayVariation, lsLink.UnidirAvailableBw, lsLink.UnidirBwUtilization, lsLink.UnidirPacketLoss)
}

func (adapter *DefaultAdapter) ConvertPrefix(lsPrefix *jagw.LsPrefix) (domain.Prefix, error) {
	return domain.NewDefaultPrefix(lsPrefix.Key, lsPrefix.IgpRouterId, lsPrefix.Prefix, lsPrefix.PrefixLen)
}

func (adapter *DefaultAdapter) ConvertSid(lsSrv6Sid *jagw.LsSrv6Sid) (domain.Sid, error) {
	return domain.NewDefaultSid(lsSrv6Sid.Key, lsSrv6Sid.IgpRouterId, lsSrv6Sid.Srv6Sid)
}

func (adapter *DefaultAdapter) convertValuesToDomain(apiValues []*api.Value) ([]domain.Value, error) {
	valueList := make([]domain.Value, 0)
	for _, apiValue := range apiValues {
		var value domain.Value
		var err error

		switch apiValue.Type {
		case api.ValueType_VALUE_TYPE_MIN_VALUE:
			value, err = domain.NewNumberValue(domain.ValueTypeMinValue, *apiValue.NumberValue)
		case api.ValueType_VALUE_TYPE_MAX_VALUE:
			value, err = domain.NewNumberValue(domain.ValueTypeMaxValue, *apiValue.NumberValue)
		case api.ValueType_VALUE_TYPE_FLEX_ALGO_NR:
			value, err = domain.NewNumberValue(domain.ValueTypeFlexAlgoNr, *apiValue.NumberValue)
		case api.ValueType_VALUE_TYPE_SFC:
			value, err = domain.NewStringValue(domain.ValueTypeSFC, *apiValue.StringValue)
		default:
			return nil, fmt.Errorf("Value type unspecified")
		}

		if err != nil {
			adapter.log.Errorln("Error creating value: ", err)
			return nil, err
		}

		valueList = append(valueList, value)
	}

	return valueList, nil
}

func (adapter *DefaultAdapter) convertIntentTypeToDomain(apiIntentType api.IntentType) (domain.IntentType, error) {
	switch apiIntentType {
	case api.IntentType_INTENT_TYPE_HIGH_BANDWIDTH:
		return domain.IntentTypeHighBandwidth, nil
	case api.IntentType_INTENT_TYPE_LOW_BANDWIDTH:
		return domain.IntentTypeLowBandwidth, nil
	case api.IntentType_INTENT_TYPE_LOW_LATENCY:
		return domain.IntentTypeLowLatency, nil
	case api.IntentType_INTENT_TYPE_LOW_PACKET_LOSS:
		return domain.IntentTypeLowPacketLoss, nil
	case api.IntentType_INTENT_TYPE_LOW_JITTER:
		return domain.IntentTypeLowJitter, nil
	case api.IntentType_INTENT_TYPE_FLEX_ALGO:
		return domain.IntentTypeFlexAlgo, nil
	case api.IntentType_INTENT_TYPE_SFC:
		return domain.IntentTypeSFC, nil
	default:
		return domain.IntentTypeUnspecified, fmt.Errorf("Intent type unspecified")
	}
}

func (adapter *DefaultAdapter) convertIntentsToDomain(apiIntents []*api.Intent) ([]domain.Intent, error) {
	intentList := make([]domain.Intent, 0)
	for _, apiIntent := range apiIntents {
		values, err := adapter.convertValuesToDomain(apiIntent.Values)
		if err != nil {
			return nil, err
		}
		intentType, err := adapter.convertIntentTypeToDomain(apiIntent.Type)
		if err != nil {
			return nil, err
		}
		intent, err := domain.NewDefaultIntent(intentType, values)
		if err != nil {
			adapter.log.Errorln("Error creating intent: ", err)
			return nil, err
		}
		intentList = append(intentList, intent)
	}
	return intentList, nil
}

func (adapter *DefaultAdapter) ConvertPathRequest(pathRequest *api.PathRequest, stream api.IntentController_GetIntentPathServer, ctx context.Context) (domain.PathRequest, error) {
	intents, err := adapter.convertIntentsToDomain(pathRequest.Intents)
	if err != nil {
		return nil, err
	}
	return domain.NewDefaultPathRequest(pathRequest.Ipv6SourceAddress, pathRequest.Ipv6DestinationAddress, intents, stream, ctx)
}

func (adapter *DefaultAdapter) convertValuesToApi(values []domain.Value) []*api.Value {
	apiValues := make([]*api.Value, 0, len(values))
	for _, value := range values {
		var apiValue *api.Value
		switch value.GetValueType() {
		case domain.ValueTypeMinValue:
			numberValue := value.GetNumberValue()
			apiValue = &api.Value{
				Type:        api.ValueType_VALUE_TYPE_MIN_VALUE,
				NumberValue: &numberValue,
			}
		case domain.ValueTypeMaxValue:
			numberValue := value.GetNumberValue()
			apiValue = &api.Value{
				Type:        api.ValueType_VALUE_TYPE_MAX_VALUE,
				NumberValue: &numberValue,
			}
		case domain.ValueTypeFlexAlgoNr:
			numberValue := value.GetNumberValue()
			apiValue = &api.Value{
				Type:        api.ValueType_VALUE_TYPE_FLEX_ALGO_NR,
				NumberValue: &numberValue,
			}
		case domain.ValueTypeSFC:
			stringValue := value.GetStringValue()
			apiValue = &api.Value{
				Type:        api.ValueType_VALUE_TYPE_SFC,
				StringValue: &stringValue,
			}
		}
		apiValues = append(apiValues, apiValue)
	}
	return apiValues
}

func (adapter *DefaultAdapter) convertIntentsToApi(intents []domain.Intent) []*api.Intent {
	apiIntents := make([]*api.Intent, len(intents))
	for index, intent := range intents {
		values := adapter.convertValuesToApi(intent.GetValues())
		apiIntent := &api.Intent{
			Type:   api.IntentType(intent.GetIntentType()),
			Values: values,
		}
		apiIntents[index] = apiIntent
	}
	return apiIntents
}

func (adapter *DefaultAdapter) ConvertPathResult(pathResult domain.PathResult) (*api.PathResult, error) {
	ipv6SidAddresses := pathResult.GetIpv6SidAddresses()
	apiPathResult := &api.PathResult{
		Ipv6SourceAddress:      pathResult.GetIpv6SourceAddress(),
		Ipv6DestinationAddress: pathResult.GetIpv6DestinationAddress(),
		Ipv6SidAddresses:       ipv6SidAddresses,
		Intents:                adapter.convertIntentsToApi(pathResult.GetIntents()),
	}
	return apiPathResult, nil
}
