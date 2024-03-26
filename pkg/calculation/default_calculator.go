package calculation

import (
	"fmt"
	"net"
	"reflect"

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
		log:    logging.DefaultLogger.WithField("subsystem", subsystem),
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

func (calculator *DefaultCalculator) getShortestPath(sourceIP, destinationIp, metric string) (graph.PathResult, error) {
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
	calculator.log.Debugf("Shortest path from %s to %s with cost %f", sourceRouterId, destinationRouterId, result.GetCost())
	return result, nil
}

func (calculator *DefaultCalculator) getNodeList(result graph.PathResult) []string {
	nodeList := make([]string, 0)
	for _, edge := range result.GetEdges() {
		to := edge.To().GetId().(string)
		nodeList = append(nodeList, to)
	}
	return nodeList
}

func (calculator *DefaultCalculator) translateGraphResultToSidList(graphResult graph.PathResult) []string {
	nodeList := calculator.getNodeList(graphResult)
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

func (calculator *DefaultCalculator) cratePathResult(graphResult graph.PathResult, pathRequest domain.PathRequest) domain.PathResult {
	sidList := calculator.translateGraphResultToSidList(graphResult)
	pathResult, err := domain.NewDefaultPathResult(pathRequest, graphResult, sidList)
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

func (calculator *DefaultCalculator) checkifPathValid(pathResult domain.PathResult, weightKind string) (float64, error) {
	totalCost := 0.0
	edges := pathResult.GetEdges()
	for _, edge := range edges {
		graphEdge, exist := calculator.graph.GetEdge(edge.GetId())
		if !exist {
			return 0, fmt.Errorf("Path not valid, edge not found in graph: %s", edge.GetId())
		}
		cost, err := graphEdge.GetWeight(weightKind)
		if err != nil {
			return 0, err
		}
		totalCost += cost
	}
	return totalCost, nil
}

func (calculator *DefaultCalculator) validateCurrentResult(intent domain.Intent, currentpathResult domain.PathResult) error {
	switch intentType := intent.GetIntentType(); intentType {
	case domain.IntentTypeLowLatency:
		return calculator.validateCurrentResultForIntent(helper.NewDefaultHelper().GetLatencyKey(), currentpathResult)
	case domain.IntentTypeLowPacketLoss:
		return calculator.validateCurrentResultForIntent(helper.NewDefaultHelper().GetPacketLossKey(), currentpathResult)
	case domain.IntentTypeLowJitter:
		return calculator.validateCurrentResultForIntent(helper.NewDefaultHelper().GetJitterKey(), currentpathResult)
	default:
		return fmt.Errorf("Unknown intent type: %s", intentType)
	}
}

func (calculator *DefaultCalculator) validateCurrentResultForIntent(weightKind string, currentpathResult domain.PathResult) error {
	cost, err := calculator.checkifPathValid(currentpathResult, weightKind)
	if err != nil {
		return fmt.Errorf("Current path is not valid anymore because of: %s", err)
	}
	if cost != currentpathResult.GetCost() {
		calculator.log.Debugf("Cost of current path has changed from %f to %f", currentpathResult.GetCost(), cost)
		currentpathResult.SetCost(cost)
	}
	return nil
}

func (calcultor *DefaultCalculator) setNewPathResult(streamSession domain.StreamSession, result domain.PathResult) {
	calcultor.log.Debugf("Setting new path with SID list %v and cost %f", result.GetIpv6SidAddresses(), result.GetCost())
	streamSession.SetPathResult(result)
}

func (calculator *DefaultCalculator) UpdatePathSession(streamSession domain.StreamSession) *domain.PathResult {
	currentpathResult := streamSession.GetPathResult()
	calculator.log.Debugln("SID list of current path is", currentpathResult.GetIpv6SidAddresses())
	pathRequest := streamSession.GetPathRequest()
	result := calculator.HandlePathRequest(pathRequest)
	if !reflect.DeepEqual(result.GetIpv6SidAddresses(), currentpathResult.GetIpv6SidAddresses()) {
		calculator.log.Debugln("Better Path found, check for applicability ")
		intents := pathRequest.GetIntents()
		if len(intents) != 1 {
			// TODO change logic if there are multiple intents
			calculator.log.Errorln("Handling of several intents not yet implemented")
			return nil
		}
		calculator.log.Debugln("Validate current path and its cost")
		err := calculator.validateCurrentResult(pathRequest.GetIntents()[0], currentpathResult)
		if err != nil {
			calculator.log.Errorln("Current Path is not valid anymore: ", err)
			calculator.setNewPathResult(streamSession, result)
			return &result
		} else {
			calculator.log.Debugln("Current path is valid, check if new path is better")
		}
		currentCost := currentpathResult.GetCost()
		resultCost := result.GetCost()
		if resultCost <= currentCost*0.9 {
			calculator.log.Debugf("Cost of new path is by more than 10 percent smaller, current: %f to new: %f", currentpathResult.GetCost(), result.GetCost())
			calculator.setNewPathResult(streamSession, result)
			return &result
		} else {
			calculator.log.Debugf("Cost has not decreased by more than 10 percent current: %f new: %f", currentpathResult.GetCost(), result.GetCost())
			calculator.log.Debugln("Current path is still applicable")
		}
	} else if result.GetCost() != currentpathResult.GetCost() {
		calculator.log.Debugf("Best path unchanged but its cost has changed from %f to %f", currentpathResult.GetCost(), result.GetCost())
		currentpathResult.SetCost(result.GetCost())
	} else {
		calculator.log.Debugln("No change in path result found")
	}
	return nil
}
