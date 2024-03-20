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

type DefaultCalculator struct {
	log    *logrus.Entry
	cache  cache.CacheService
	graph  graph.Graph
	helper helper.Helper
}

func NewDefaultCalculator(cache cache.CacheService, graph graph.Graph, helper helper.Helper) *DefaultCalculator {
	return &DefaultCalculator{
		log:    logging.DefaultLogger.WithField("subsystem", "calculator"),
		cache:  cache,
		graph:  graph,
		helper: helper,
	}
}

func (calcultor *DefaultCalculator) getNetworkAddress(ip string) net.IP {
	ipv6Address := net.ParseIP(ip)
	mask := net.CIDRMask(64, 128)
	return ipv6Address.Mask(mask)
}

func (calculator *DefaultCalculator) getShortestPath(sourceIP, destinationIp, metric string) ([]string, error) {
	sourceIpv6 := calculator.getNetworkAddress(sourceIP)
	destinationIpv6 := calculator.getNetworkAddress(destinationIp)
	sourceRouterId, ok := calculator.cache.GetRouterIdFromNetworkAddress(sourceIpv6.String())
	if !ok {
		return nil, fmt.Errorf("Router ID not found for source IP: %s", sourceRouterId)
	}
	destinationRouterId, ok := calculator.cache.GetRouterIdFromNetworkAddress(destinationIpv6.String())
	if !ok {
		return nil, fmt.Errorf("Router ID not found for destination IP: %s", destinationRouterId)
	}
	calculator.log.Debugln("Source Router ID: ", sourceRouterId)
	calculator.log.Debugln("Destination Router ID: ", destinationRouterId)

	source, exist := calculator.graph.GetNode(sourceRouterId)
	if !exist {
		return nil, fmt.Errorf("Source router not found")
	}
	destination, exist := calculator.graph.GetNode(destinationRouterId)
	if !exist {
		return nil, fmt.Errorf("Destination router not found")
	}
	result, err := calculator.graph.GetShortestPath(source, destination, metric)
	if err != nil {
		return nil, err
	}
	nodeList := make([]string, 0)
	for _, edge := range result.GetEdges() {
		to := edge.To().GetId().(string)
		nodeList = append(nodeList, to)
	}
	calculator.log.Debugf("Shortest path from %s to %s: %v with cost %f", sourceRouterId, destinationRouterId, nodeList, result.GetCost())
	return nodeList, nil
}

func (calculator *DefaultCalculator) translateNodesToSidList(nodeList []string) []string {
	calculator.log.Debugln("Node List: ", nodeList)
	sidList := make([]string, 0)
	for _, node := range nodeList {
		sid, ok := calculator.cache.GetSidFromRouterId(node)
		if !ok {
			calculator.log.Errorln("SID not found for router: ", node)
			continue
		}
		sidList = append(sidList, sid)
	}

	calculator.log.Debugln("Translated SID List: ", sidList)
	return sidList
}

func (calculator *DefaultCalculator) cratePathResult(nodeIds []string, pathRequest domain.PathRequest) domain.PathResult {
	sidList := calculator.translateNodesToSidList(nodeIds)
	pathResult, err := domain.NewDefaultPathResult(pathRequest, sidList)
	if err != nil {
		calculator.log.Errorln("Error creating path result: ", err)
		return nil
	}
	return pathResult
}

func (calculator *DefaultCalculator) handleSingleIntent(intent domain.Intent, pathRequest domain.PathRequest) domain.PathResult {
	ipv6SourceAddress := pathRequest.GetIpv6SourceAddress()
	ipv6DestinationAddress := pathRequest.GetIpv6DestinationAddress()
	switch intentType := intent.GetIntentType(); intentType {
	case domain.IntentTypeLowLatency:
		nodeIds, err := calculator.getShortestPath(ipv6SourceAddress, ipv6DestinationAddress, helper.NewDefaultHelper().GetLatencyKey())
		if err != nil {
			calculator.log.Errorln("Error getting shortest path: ", err)
			return nil
		}
		return calculator.cratePathResult(nodeIds, pathRequest)
	case domain.IntentTypeLowPacketLoss:
		nodeIds, err := calculator.getShortestPath(ipv6SourceAddress, ipv6DestinationAddress, helper.NewDefaultHelper().GetPacketLossKey())
		if err != nil {
			calculator.log.Errorln("Error getting shortest path: ", err)
			return nil
		}
		return calculator.cratePathResult(nodeIds, pathRequest)
	case domain.IntentTypeLowJitter:
		nodeIds, err := calculator.getShortestPath(ipv6SourceAddress, ipv6DestinationAddress, helper.NewDefaultHelper().GetJitterKey())
		if err != nil {
			calculator.log.Errorln("Error getting shortest path: ", err)
			return nil
		}
		return calculator.cratePathResult(nodeIds, pathRequest)
	default:
		return nil
	}
}

func (calculator *DefaultCalculator) HandlePathRequest(pathRequest domain.PathRequest) domain.PathResult {
	calculator.graph.Lock()
	calculator.cache.Lock()
	defer calculator.graph.Unlock()
	defer calculator.cache.Unlock()
	intents := pathRequest.GetIntents()
	if len(intents) == 1 {
		return calculator.handleSingleIntent(intents[0], pathRequest)
	} else {
		// TODO: change logic if there are multiple intents
		calculator.log.Errorln("Handling of several intents not yet implemented")
		for _, intent := range intents {
			calculator.log.Debugln("Received Intent: ", intent)
		}
	}
	return nil
}
