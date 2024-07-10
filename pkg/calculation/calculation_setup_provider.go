package calculation

import (
	"fmt"
	"net"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type CalculationSetupProvider struct {
	log   *logrus.Entry
	cache cache.Cache
	graph graph.Graph
}

func NewCalculationSetupProvider(cache cache.Cache, graph graph.Graph) *CalculationSetupProvider {
	return &CalculationSetupProvider{
		log:   logging.DefaultLogger.WithField("subsystem", subsystem),
		cache: cache,
		graph: graph,
	}
}

func (provider *CalculationSetupProvider) getNetworkAddress(ip string) (net.IP, error) {
	ipv6Address := net.ParseIP(ip)
	if ipv6Address == nil {
		return nil, fmt.Errorf("IPv6 Network address could not be parsed, invalid IPv6 address: %s", ip)
	}
	mask := net.CIDRMask(64, 128)
	return ipv6Address.Mask(mask), nil
}

func (provider *CalculationSetupProvider) getSourceNode(pathRequest domain.PathRequest) (graph.Node, error) {
	sourceIpv6, err := provider.getNetworkAddress(pathRequest.GetIpv6SourceAddress())
	if err != nil {
		return nil, err
	}
	sourceRouterId := provider.cache.GetRouterIdFromNetworkAddress(sourceIpv6.String())
	if sourceRouterId == "" {
		return nil, fmt.Errorf("Router ID not found for source IP: %s", sourceIpv6)
	}
	source := provider.graph.GetNode(sourceRouterId)
	if source == nil {
		return nil, fmt.Errorf("Source router not found")
	}
	return source, nil
}

func (provider *CalculationSetupProvider) GetDestinationNode(pathRequest domain.PathRequest) (graph.Node, error) {
	destinationIpv6, err := provider.getNetworkAddress(pathRequest.GetIpv6DestinationAddress())
	if err != nil {
		return nil, err
	}
	destinationRouterId := provider.cache.GetRouterIdFromNetworkAddress(destinationIpv6.String())
	if destinationRouterId == "" {
		return nil, fmt.Errorf("Router ID not found for destination IP: %s", destinationIpv6)
	}
	destination := provider.graph.GetNode(destinationRouterId)
	if destination == nil {
		return nil, fmt.Errorf("Destination router not found")
	}
	return destination, nil
}

func (provider *CalculationSetupProvider) getWeightKeyAndCalcMode(intentType domain.IntentType) (helper.WeightKey, CalculationMode) {
	switch intentType {
	case domain.IntentTypeSFC:
		return helper.IgpMetricKey, CalculationModeSum
	case domain.IntentTypeFlexAlgo:
		return helper.IgpMetricKey, CalculationModeSum
	case domain.IntentTypeHighBandwidth:
		return helper.AvailableBandwidthKey, CalculationModeMax
	case domain.IntentTypeLowBandwidth:
		return helper.MaximumLinkBandwidth, CalculationModeMin
	case domain.IntentTypeLowLatency:
		return helper.LatencyKey, CalculationModeSum
	case domain.IntentTypeLowPacketLoss:
		return helper.PacketLossKey, CalculationModeSum
	case domain.IntentTypeLowJitter:
		return helper.JitterKey, CalculationModeSum
	case domain.IntentLowUtilization:
		return helper.UtilizedBandwidthKey, CalculationModeSum
	default:
		return helper.UndefinedKey, CalculationModeUndefined
	}
}

func (provider *CalculationSetupProvider) getWeightKey(intentType domain.IntentType) helper.WeightKey {
	switch intentType {
	case domain.IntentTypeLowLatency:
		return helper.NormalizedLatencyKey
	case domain.IntentTypeLowJitter:
		return helper.NormalizedJitterKey
	case domain.IntentTypeLowPacketLoss:
		return helper.NormalizedPacketLossKey
	case domain.IntentTypeHighBandwidth:
		return helper.AvailableBandwidthKey
	default:
		return helper.UndefinedKey
	}
}

func (provider *CalculationSetupProvider) getWeightKeysandCalculationMode(intents []domain.Intent) ([]helper.WeightKey, CalculationMode) {
	if len(intents) == 1 {
		weightKey, calcType := provider.getWeightKeyAndCalcMode(intents[0].GetIntentType())
		return []helper.WeightKey{weightKey}, calcType
	} else {
		calculationType := CalculationModeSum
		start := 0
		if intents[0].GetIntentType() == domain.IntentTypeFlexAlgo {
			start = 1
		}
		weightKeys := make([]helper.WeightKey, len(intents)-start)
		weightKeysIndex := 0

		for i := start; i < len(intents); i++ {
			weightKey := provider.getWeightKey(intents[i].GetIntentType())
			if len(weightKeys) > weightKeysIndex {
				weightKeys[weightKeysIndex] = weightKey
			} else {
				weightKeys = append(weightKeys, weightKey)
			}
			weightKeysIndex++
		}

		return weightKeys, calculationType
	}
}

func (provider *CalculationSetupProvider) getMaxConstraints(intents []domain.Intent, weightKeys []helper.WeightKey) map[helper.WeightKey]float64 {
	maxValues := make(map[helper.WeightKey]float64)
	for index, intent := range intents {
		values := intent.GetValues()
		for _, value := range values {
			key := weightKeys[index]
			if value.GetValueType() == domain.ValueTypeMaxValue &&
				(key == helper.NormalizedLatencyKey || key == helper.NormalizedJitterKey || key == helper.NormalizedPacketLossKey) {
				maxValues[key] = float64(value.GetNumberValue())
				if key == helper.NormalizedPacketLossKey {
					maxValues[key] = maxValues[key] / 100
				}
			}
		}
	}
	return maxValues
}

func (provider *CalculationSetupProvider) getMinConstraints(intents []domain.Intent, weightKeys []helper.WeightKey) map[helper.WeightKey]float64 {
	minValues := make(map[helper.WeightKey]float64)
	for index, intent := range intents {
		values := intent.GetValues()
		for _, value := range values {
			key := weightKeys[index]
			if value.GetValueType() == domain.ValueTypeMinValue && key == helper.AvailableBandwidthKey {
				minValues[key] = float64(value.GetNumberValue())
			}
		}
	}
	return minValues
}

func (provider *CalculationSetupProvider) getServiceSids(serviceFunctionChainIntent domain.Intent) ([][]string, error) {
	serviceSids := make([][]string, 0)
	for index, value := range serviceFunctionChainIntent.GetValues() {
		value := value.GetStringValue()
		serviceSids = append(serviceSids, make([]string, 0))
		serviceSids[index] = provider.cache.GetServiceSids(value)
		if len(serviceSids[index]) == 0 {
			return nil, fmt.Errorf("No SIDs found for service: %s", value)
		}
	}
	provider.log.Debugln("Service SIDs: ", serviceSids)
	return serviceSids, nil
}

func (provider *CalculationSetupProvider) getServiceRouter(serviceSids [][]string) ([][]string, map[string]string, error) {
	serviceRouters := make([][]string, len(serviceSids))
	routerServiceMap := make(map[string]string)
	for index, sids := range serviceSids {
		for _, sid := range sids {
			routerId := provider.cache.GetRouterIdFromNetworkAddress(sid)
			if routerId == "" {
				return nil, nil, fmt.Errorf("Router ID not found for SID: %s", sid)
			}
			routerServiceMap[routerId] = sid
			serviceRouters[index] = append(serviceRouters[index], routerId)
		}
	}
	provider.log.Debugln("Service Routers: ", serviceRouters)
	return serviceRouters, routerServiceMap, nil
}

func (provider *CalculationSetupProvider) getServiceChainCombinations(serviceRouter [][]string) [][]string {
	serviceChainCombination := [][]string{{}}

	for _, routerIds := range serviceRouter {
		newQueue := [][]string{}
		for _, combination := range serviceChainCombination {
			for _, routerId := range routerIds {
				newCombination := append([]string{}, combination...)
				newCombination = append(newCombination, routerId)
				newQueue = append(newQueue, newCombination)
			}
		}
		serviceChainCombination = newQueue
	}
	provider.log.Debugln("Service Chain Combinations: ", serviceChainCombination)
	return serviceChainCombination
}

func (provider *CalculationSetupProvider) PerformServiceFunctionChainSetup(serviceFunctionChainIntent domain.Intent) (*SfcCalculationOptions, error) {
	sfcCalculationOptions := &SfcCalculationOptions{}
	serviceSids, err := provider.getServiceSids(serviceFunctionChainIntent)
	if err != nil {
		return nil, fmt.Errorf("Error getting service SIDs: %s", err)
	}

	serviceRouter, routerServiceMap, err := provider.getServiceRouter(serviceSids)
	if err != nil {
		return nil, fmt.Errorf("Error getting service routers: %s", err)
	}
	sfcCalculationOptions.routerServiceMap = routerServiceMap

	sfcCalculationOptions.serviceFunctionChain = provider.getServiceChainCombinations(serviceRouter)

	return sfcCalculationOptions, nil
}

func (provider *CalculationSetupProvider) PerformSetup(pathRequest domain.PathRequest) (*CalculationOptions, error) {
	calculationSetupOption := &CalculationOptions{}
	var err error
	calculationSetupOption.sourceNode, err = provider.getSourceNode(pathRequest)
	if err != nil {
		return nil, err
	}

	calculationSetupOption.destinationNode, err = provider.GetDestinationNode(pathRequest)
	if err != nil {
		return nil, err
	}

	intents := pathRequest.GetIntents()
	if len(intents) == 0 {
		return nil, fmt.Errorf("No intents defined in path request")
	}

	calculationSetupOption.weightKeys, calculationSetupOption.calculationMode = provider.getWeightKeysandCalculationMode(intents)
	if calculationSetupOption.calculationMode == CalculationModeUndefined {
		return nil, fmt.Errorf("Calculation mode not defined for intents")
	}
	calculationSetupOption.maxConstraints = provider.getMaxConstraints(intents, calculationSetupOption.weightKeys)
	calculationSetupOption.minConstraints = provider.getMinConstraints(intents, calculationSetupOption.weightKeys)

	return calculationSetupOption, nil
}
