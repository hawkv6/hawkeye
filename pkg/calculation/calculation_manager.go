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

type CalculationManager struct {
	log    *logrus.Entry
	cache  cache.Cache
	graph  graph.Graph
	helper helper.Helper
}

func NewCalculationManager(cache cache.Cache, graph graph.Graph, helper helper.Helper) *CalculationManager {
	return &CalculationManager{
		log:    logging.DefaultLogger.WithField("subsystem", subsystem),
		cache:  cache,
		graph:  graph,
		helper: helper,
	}
}

func (calcultor *CalculationManager) getNetworkAddress(ip string) net.IP {
	ipv6Address := net.ParseIP(ip)
	mask := net.CIDRMask(64, 128)
	return ipv6Address.Mask(mask)
}

func (manager *CalculationManager) getShortestPath(sourceIP, destinationIp, metric string) (graph.PathResult, error) {
	sourceIpv6 := manager.getNetworkAddress(sourceIP)
	destinationIpv6 := manager.getNetworkAddress(destinationIp)
	sourceRouterId, ok := manager.cache.GetRouterIdFromNetworkAddress(sourceIpv6.String())
	if !ok {
		return nil, fmt.Errorf("Router ID not found for source IP: %s", sourceRouterId)
	}
	destinationRouterId, ok := manager.cache.GetRouterIdFromNetworkAddress(destinationIpv6.String())
	if !ok {
		return nil, fmt.Errorf("Router ID not found for destination IP: %s", destinationRouterId)
	}
	manager.log.Debugln("Source Router ID: ", sourceRouterId)
	manager.log.Debugln("Destination Router ID: ", destinationRouterId)

	source, exist := manager.graph.GetNode(sourceRouterId)
	if !exist {
		return nil, fmt.Errorf("Source router not found")
	}
	destination, exist := manager.graph.GetNode(destinationRouterId)
	if !exist {
		return nil, fmt.Errorf("Destination router not found")
	}
	result, err := manager.graph.GetShortestPath(source, destination, metric)
	if err != nil {
		return nil, err
	}
	manager.log.Debugf("Shortest path from %s to %s with cost %f", sourceRouterId, destinationRouterId, result.GetCost())
	return result, nil
}

func (manager *CalculationManager) getNodeList(result graph.PathResult) []string {
	nodeList := make([]string, 0)
	for _, edge := range result.GetEdges() {
		to := edge.To().GetId().(string)
		nodeList = append(nodeList, to)
	}
	return nodeList
}

func (manager *CalculationManager) translateGraphResultToSidList(graphResult graph.PathResult) []string {
	nodeList := manager.getNodeList(graphResult)
	manager.log.Debugln("Node List: ", nodeList)
	sidList := make([]string, 0)
	for _, node := range nodeList {
		sid, ok := manager.cache.GetSidFromRouterId(node)
		if !ok {
			manager.log.Errorln("SID not found for router: ", node)
			continue
		}
		sidList = append(sidList, sid)
	}

	manager.log.Debugln("Translated SID List: ", sidList)
	return sidList
}

func (manager *CalculationManager) cratePathResult(graphResult graph.PathResult, pathRequest domain.PathRequest) domain.PathResult {
	sidList := manager.translateGraphResultToSidList(graphResult)
	pathResult, err := domain.NewDomainPathResult(pathRequest, graphResult, sidList)
	if err != nil {
		manager.log.Errorln("Error creating path result: ", err)
		return nil
	}
	return pathResult
}

func (manager *CalculationManager) findPathForSingleIntent(intent domain.Intent, pathRequest domain.PathRequest) domain.PathResult {
	ipv6SourceAddress := pathRequest.GetIpv6SourceAddress()
	ipv6DestinationAddress := pathRequest.GetIpv6DestinationAddress()
	switch intentType := intent.GetIntentType(); intentType {
	case domain.IntentTypeLowLatency:
		nodeIds, err := manager.getShortestPath(ipv6SourceAddress, ipv6DestinationAddress, helper.NewDefaultHelper().GetLatencyKey())
		if err != nil {
			manager.log.Errorln("Error getting shortest path: ", err)
			return nil
		}
		return manager.cratePathResult(nodeIds, pathRequest)
	case domain.IntentTypeLowPacketLoss:
		nodeIds, err := manager.getShortestPath(ipv6SourceAddress, ipv6DestinationAddress, helper.NewDefaultHelper().GetPacketLossKey())
		if err != nil {
			manager.log.Errorln("Error getting shortest path: ", err)
			return nil
		}
		return manager.cratePathResult(nodeIds, pathRequest)
	case domain.IntentTypeLowJitter:
		nodeIds, err := manager.getShortestPath(ipv6SourceAddress, ipv6DestinationAddress, helper.NewDefaultHelper().GetJitterKey())
		if err != nil {
			manager.log.Errorln("Error getting shortest path: ", err)
			return nil
		}
		return manager.cratePathResult(nodeIds, pathRequest)
	default:
		return nil
	}
}

func (manager *CalculationManager) CalculateBestPath(pathRequest domain.PathRequest) domain.PathResult {
	manager.graph.Lock()
	manager.cache.Lock()
	defer manager.graph.Unlock()
	defer manager.cache.Unlock()
	intents := pathRequest.GetIntents()
	if len(intents) == 1 {
		return manager.findPathForSingleIntent(intents[0], pathRequest)
	} else {
		// TODO: change logic if there are multiple intents
		manager.log.Errorln("Handling of several intents not yet implemented")
		for _, intent := range intents {
			manager.log.Debugln("Received Intent: ", intent)
		}
	}
	return nil
}

func (manager *CalculationManager) checkifPathValid(pathResult domain.PathResult, weightKind string) (float64, error) {
	totalCost := 0.0
	edges := pathResult.GetEdges()
	for _, edge := range edges {
		graphEdge, exist := manager.graph.GetEdge(edge.GetId())
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

func (manager *CalculationManager) validateCurrentResult(intent domain.Intent, currentpathResult domain.PathResult) error {
	switch intentType := intent.GetIntentType(); intentType {
	case domain.IntentTypeLowLatency:
		return manager.validateCurrentResultForIntent(helper.NewDefaultHelper().GetLatencyKey(), currentpathResult)
	case domain.IntentTypeLowPacketLoss:
		return manager.validateCurrentResultForIntent(helper.NewDefaultHelper().GetPacketLossKey(), currentpathResult)
	case domain.IntentTypeLowJitter:
		return manager.validateCurrentResultForIntent(helper.NewDefaultHelper().GetJitterKey(), currentpathResult)
	default:
		return fmt.Errorf("Unknown intent type: %s", intentType)
	}
}

func (manager *CalculationManager) validateCurrentResultForIntent(weightKind string, currentpathResult domain.PathResult) error {
	cost, err := manager.checkifPathValid(currentpathResult, weightKind)
	if err != nil {
		return fmt.Errorf("Current path is not valid anymore because of: %s", err)
	}
	if cost != currentpathResult.GetCost() {
		manager.log.Debugf("Cost of current path has changed from %f to %f", currentpathResult.GetCost(), cost)
		currentpathResult.SetCost(cost)
	}
	return nil
}

func (calcultor *CalculationManager) setNewPathResult(streamSession domain.StreamSession, result domain.PathResult) {
	calcultor.log.Debugf("Setting new path with SID list %v and cost %f", result.GetIpv6SidAddresses(), result.GetCost())
	streamSession.SetPathResult(result)
}

func (manager *CalculationManager) CalculatePathUpdate(streamSession domain.StreamSession) *domain.PathResult {
	currentpathResult := streamSession.GetPathResult()
	manager.log.Debugln("SID list of current path is", currentpathResult.GetIpv6SidAddresses())
	pathRequest := streamSession.GetPathRequest()
	result := manager.CalculateBestPath(pathRequest)
	if !reflect.DeepEqual(result.GetIpv6SidAddresses(), currentpathResult.GetIpv6SidAddresses()) {
		manager.log.Debugln("Better Path found, check for applicability ")
		intents := pathRequest.GetIntents()
		if len(intents) != 1 {
			// TODO change logic if there are multiple intents
			manager.log.Errorln("Handling of several intents not yet implemented")
			return nil
		}
		manager.log.Debugln("Validate current path and its cost")
		err := manager.validateCurrentResult(pathRequest.GetIntents()[0], currentpathResult)
		if err != nil {
			manager.log.Errorln("Current Path is not valid anymore: ", err)
			manager.setNewPathResult(streamSession, result)
			return &result
		} else {
			manager.log.Debugln("Current path is valid, check if new path is better")
		}
		currentCost := currentpathResult.GetCost()
		resultCost := result.GetCost()
		if resultCost <= currentCost*0.9 {
			manager.log.Debugf("Cost of new path is by more than 10 percent smaller, current: %f to new: %f", currentpathResult.GetCost(), result.GetCost())
			manager.setNewPathResult(streamSession, result)
			return &result
		} else {
			manager.log.Debugf("Cost has not decreased by more than 10 percent current: %f new: %f", currentpathResult.GetCost(), result.GetCost())
			manager.log.Debugln("Current path is still applicable")
		}
	} else if result.GetCost() != currentpathResult.GetCost() {
		manager.log.Debugf("Best path unchanged but its cost has changed from %f to %f", currentpathResult.GetCost(), result.GetCost())
		currentpathResult.SetCost(result.GetCost())
	} else {
		manager.log.Debugln("No change in path result found")
	}
	return nil
}
