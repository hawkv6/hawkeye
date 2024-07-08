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

type DomainAdapter struct {
	log *logrus.Entry
}

func NewDomainAdapter() *DomainAdapter {
	return &DomainAdapter{
		log: logging.DefaultLogger.WithField("subsystem", Subsystem),
	}
}

func (adapter *DomainAdapter) ConvertNode(lsNode *jagw.LsNode) (domain.Node, error) {
	return domain.NewDomainNode(lsNode.Key, lsNode.IgpRouterId, lsNode.Name, lsNode.SrAlgorithm)
}

func (adapter *DomainAdapter) ConvertNodeEvent(lsNodeEvent *jagw.LsNodeEvent) (domain.NetworkEvent, error) {
	if lsNodeEvent == nil || lsNodeEvent.Action == nil || lsNodeEvent.Key == nil {
		return nil, fmt.Errorf("LsNodeEvent, action or key is nil")
	}
	if *lsNodeEvent.Action == "del" {
		return domain.NewDeleteNodeEvent(*lsNodeEvent.Key), nil
	} else if *lsNodeEvent.Action == "add" {
		node, err := adapter.ConvertNode(lsNodeEvent.LsNode)
		if err != nil {
			return nil, fmt.Errorf("Error converting LsNode to Node: %s", err)
		}
		return domain.NewAddNodeEvent(node), nil
	}
	if *lsNodeEvent.Action == "update" {
		node, err := adapter.ConvertNode(lsNodeEvent.LsNode)
		if err != nil {
			return nil, fmt.Errorf("Error converting LsNode to Node: %s", err)
		}
		return domain.NewUpdateNodeEvent(node), nil
	} else {

		return nil, fmt.Errorf("Unknown action: %s", *lsNodeEvent.Action)
	}
}

func (adapter *DomainAdapter) ConvertLink(lsLink *jagw.LsLink) (domain.Link, error) {
	return domain.NewDomainLink(lsLink.Key, lsLink.IgpRouterId, lsLink.RemoteIgpRouterId, lsLink.IgpMetric, lsLink.UnidirLinkDelay, lsLink.UnidirDelayVariation, lsLink.MaxLinkBwKbps, lsLink.UnidirAvailableBw, lsLink.UnidirBwUtilization, lsLink.UnidirPacketLossPercentage, lsLink.NormalizedUnidirLinkDelay, lsLink.NormalizedUnidirDelayVariation, lsLink.NormalizedUnidirPacketLoss)
}

func (adapter *DomainAdapter) ConvertLinkEvent(lsLinkEvent *jagw.LsLinkEvent) (domain.NetworkEvent, error) {
	if lsLinkEvent == nil || lsLinkEvent.Action == nil || lsLinkEvent.Key == nil {
		return nil, fmt.Errorf("LsLinkEvent, action or key is nil")
	}
	if *lsLinkEvent.Action == "del" {
		return domain.NewDeleteLinkEvent(*lsLinkEvent.Key), nil
	} else if *lsLinkEvent.Action == "add" {
		link, err := adapter.ConvertLink(lsLinkEvent.LsLink)
		if err != nil {
			return nil, fmt.Errorf("Error converting LsLink to Link: %s", err)
		}
		return domain.NewAddLinkEvent(link), nil
	} else if *lsLinkEvent.Action == "update" {
		link, err := adapter.ConvertLink(lsLinkEvent.LsLink)
		if err != nil {
			return nil, fmt.Errorf("Error converting LsLink to Link: %s", err)
		}
		return domain.NewUpdateLinkEvent(link), nil
	} else {
		return nil, fmt.Errorf("Unknown action: %s", *lsLinkEvent.Action)
	}
}

func (adapter *DomainAdapter) ConvertPrefix(lsPrefix *jagw.LsPrefix) (domain.Prefix, error) {
	return domain.NewDomainPrefix(lsPrefix.Key, lsPrefix.IgpRouterId, lsPrefix.Prefix, lsPrefix.PrefixLen)

}

func (adapter *DomainAdapter) ConvertPrefixEvent(lsPrefixEvent *jagw.LsPrefixEvent) (domain.NetworkEvent, error) {
	if lsPrefixEvent == nil || lsPrefixEvent.Action == nil || lsPrefixEvent.Key == nil {
		return nil, fmt.Errorf("LsPrefixEvent, action or key is nil")
	}
	if *lsPrefixEvent.Action == "del" {
		return domain.NewDeletePrefixEvent(*lsPrefixEvent.Key), nil
	} else if *lsPrefixEvent.Action == "add" {
		prefix, err := adapter.ConvertPrefix(lsPrefixEvent.LsPrefix)
		if err != nil {
			return nil, fmt.Errorf("Error converting LsPrefix to Prefix: %s", err)
		}
		return domain.NewAddPrefixEvent(prefix), nil
	} else {
		return nil, fmt.Errorf("Unknown action: %s", *lsPrefixEvent.Action)
	}
}

func (adapter *DomainAdapter) ConvertSid(lsSrv6Sid *jagw.LsSrv6Sid) (domain.Sid, error) {
	return domain.NewDomainSid(lsSrv6Sid.Key, lsSrv6Sid.IgpRouterId, lsSrv6Sid.Srv6Sid, lsSrv6Sid.Srv6EndpointBehavior.Algorithm)
}

func (adapter *DomainAdapter) ConvertSidEvent(lsSrv6SidEvent *jagw.LsSrv6SidEvent) (domain.NetworkEvent, error) {
	if lsSrv6SidEvent == nil || lsSrv6SidEvent.Action == nil || lsSrv6SidEvent.Key == nil {
		return nil, fmt.Errorf("LsSrv6SidEvent, action or key is nil")
	}
	if *lsSrv6SidEvent.Action == "del" {
		return domain.NewDeleteSidEvent(*lsSrv6SidEvent.Key), nil
	} else if *lsSrv6SidEvent.Action == "add" {
		sid, err := adapter.ConvertSid(lsSrv6SidEvent.LsSrv6Sid)
		if err != nil {
			return nil, fmt.Errorf("Error converting LsSrv6Sid to Sid: %s", err)
		}
		return domain.NewAddSidEvent(sid), nil
	} else {
		return nil, fmt.Errorf("Unknown action: %s", *lsSrv6SidEvent.Action)
	}
}

func (adapter *DomainAdapter) convertValuesToDomain(apiValues []*api.Value) ([]domain.Value, error) {
	valueList := make([]domain.Value, 0)
	for _, apiValue := range apiValues {
		var value domain.Value
		var err error

		switch apiValue.Type {
		case api.ValueType_VALUE_TYPE_MIN_VALUE:
			value, err = domain.NewNumberValue(domain.ValueTypeMinValue, apiValue.NumberValue)
		case api.ValueType_VALUE_TYPE_MAX_VALUE:
			value, err = domain.NewNumberValue(domain.ValueTypeMaxValue, apiValue.NumberValue)
		case api.ValueType_VALUE_TYPE_FLEX_ALGO_NR:
			value, err = domain.NewNumberValue(domain.ValueTypeFlexAlgoNr, apiValue.NumberValue)
		case api.ValueType_VALUE_TYPE_SFC:
			value, err = domain.NewStringValue(domain.ValueTypeSFC, apiValue.StringValue)
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

func (adapter *DomainAdapter) convertIntentTypeToDomain(apiIntentType api.IntentType) (domain.IntentType, error) {
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
	case api.IntentType_INTENT_TYPE_LOW_UTILIZATION:
		return domain.IntentLowUtilization, nil
	default:
		return domain.IntentTypeUnspecified, fmt.Errorf("Intent type unspecified")
	}
}

func (adapter *DomainAdapter) convertIntentsToDomain(apiIntents []*api.Intent) ([]domain.Intent, error) {
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
		intent, err := domain.NewDomainIntent(intentType, values)
		if err != nil {
			adapter.log.Errorln("Error creating intent: ", err)
			return nil, err
		}
		intentList = append(intentList, intent)
	}
	return intentList, nil
}

func (adapter *DomainAdapter) ConvertPathRequest(pathRequest *api.PathRequest, stream api.IntentController_GetIntentPathServer, ctx context.Context) (domain.PathRequest, error) {
	intents, err := adapter.convertIntentsToDomain(pathRequest.Intents)
	if err != nil {
		adapter.log.Errorln("Error converting intents: ", err)
		return nil, err
	}
	return domain.NewDomainPathRequest(pathRequest.Ipv6SourceAddress, pathRequest.Ipv6DestinationAddress, intents, stream, ctx)
}

func (adapter *DomainAdapter) convertValuesToApi(values []domain.Value) []*api.Value {
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

func (adapter *DomainAdapter) convertIntentsToApi(intents []domain.Intent) []*api.Intent {
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

func (adapter *DomainAdapter) ConvertPathResult(pathResult domain.PathResult) (*api.PathResult, error) {
	if pathResult == nil {
		return nil, fmt.Errorf("PathResult could not be calculated due to error")
	}
	ipv6SidAddresses := pathResult.GetIpv6SidAddresses()
	apiPathResult := &api.PathResult{
		Ipv6SourceAddress:      pathResult.GetIpv6SourceAddress(),
		Ipv6DestinationAddress: pathResult.GetIpv6DestinationAddress(),
		Ipv6SidAddresses:       ipv6SidAddresses,
		Intents:                adapter.convertIntentsToApi(pathResult.GetIntents()),
	}
	return apiPathResult, nil
}
