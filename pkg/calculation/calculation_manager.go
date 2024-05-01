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

func (manager *CalculationManager) getNodesFromPath(path graph.Path) []string {
	nodeList := make([]string, 0)
	for _, edge := range path.GetEdges() {
		to := edge.To().GetId().(string)
		nodeList = append(nodeList, to)
	}
	return nodeList
}

func (manager *CalculationManager) translatePathToSidList(path graph.Path) []string {
	nodeList := manager.getNodesFromPath(path)
	manager.log.Debugln("Node in Path: ", nodeList)
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

func (manger *CalculationManager) getSourceNode(pathRequest domain.PathRequest) (graph.Node, error) {
	sourceIpv6 := manger.getNetworkAddress(pathRequest.GetIpv6SourceAddress())
	sourceRouterId, ok := manger.cache.GetRouterIdFromNetworkAddress(sourceIpv6.String())
	if !ok {
		return nil, fmt.Errorf("Router ID not found for source IP: %s", sourceRouterId)
	}
	source, exist := manger.graph.GetNode(sourceRouterId)
	if !exist {
		return nil, fmt.Errorf("Source router not found")
	}
	return source, nil
}

func (manager *CalculationManager) GetDestinationNode(pathRequest domain.PathRequest) (graph.Node, error) {
	destinationIpv6 := manager.getNetworkAddress(pathRequest.GetIpv6DestinationAddress())
	destinationRouterId, ok := manager.cache.GetRouterIdFromNetworkAddress(destinationIpv6.String())
	if !ok {
		return nil, fmt.Errorf("Router ID not found for destination IP: %s", destinationRouterId)
	}
	destination, exist := manager.graph.GetNode(destinationRouterId)
	if !exist {
		return nil, fmt.Errorf("Destination router not found")
	}
	return destination, nil
}

func (manager *CalculationManager) getWeightKeyAndCalcType(intentType domain.IntentType) (helper.WeightKey, CalculationType) {
	switch intentType {
	case domain.IntentTypeHighBandwidth:
		return helper.AvailableBandwidthKey, CalculationTypeMax
	case domain.IntentTypeLowBandwidth:
		return helper.AvailableBandwidthKey, CalculationTypeMin
	case domain.IntentTypeLowLatency:
		return helper.LatencyKey, CalculationTypeSum
	case domain.IntentTypeLowPacketLoss:
		return helper.PacketLossKey, CalculationTypeSum
	case domain.IntentTypeLowJitter:
		return helper.JitterKey, CalculationTypeSum
	default:
		return "", ""
	}
}

func (manager *CalculationManager) createPathResult(path graph.Path, pathRequest domain.PathRequest) domain.PathResult {
	sidList := manager.translatePathToSidList(path)
	pathResult, err := domain.NewDomainPathResult(pathRequest, path, sidList)
	if err != nil {
		manager.log.Errorln("Error creating path result: ", err)
		return nil
	}
	return pathResult
}

func (manager *CalculationManager) findPathForSingleIntent(intent domain.Intent, pathRequest domain.PathRequest) domain.PathResult {
	sourceNode, err := manager.getSourceNode(pathRequest)
	if err != nil {
		manager.log.Errorln("Error getting source node: ", err)
		return nil
	}
	destinationNode, err := manager.GetDestinationNode(pathRequest)
	if err != nil {
		manager.log.Errorln("Error getting destination node: ", err)
		return nil
	}
	weightKey, calcType := manager.getWeightKeyAndCalcType(intent.GetIntentType())
	if weightKey == "" || calcType == "" {
		return nil
	}
	calculation := NewShortestPathCalculation(manager.graph, sourceNode, destinationNode, weightKey, calcType)
	path, err := calculation.Execute()
	if err != nil {
		manager.log.Errorln("Error getting shortest path: ", err)
		return nil
	}
	manager.log.Debugf("Shortest Path found with cost %g: ", path.GetCost())
	return manager.createPathResult(path, pathRequest)

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

func (manager *CalculationManager) checkifPathValid(pathResult domain.PathResult, weightType helper.WeightKey) (float64, error) {
	totalCost := 0.0
	edges := pathResult.GetEdges()
	for _, edge := range edges {
		graphEdge, exist := manager.graph.GetEdge(edge.GetId())
		if !exist {
			return 0, fmt.Errorf("Path not valid, edge not found in graph: %s", edge.GetId())
		}
		cost, err := graphEdge.GetWeight(weightType)
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
		return manager.validateCurrentResultForIntent(helper.LatencyKey, currentpathResult)
	case domain.IntentTypeLowPacketLoss:
		return manager.validateCurrentResultForIntent(helper.PacketLossKey, currentpathResult)
	case domain.IntentTypeLowJitter:
		return manager.validateCurrentResultForIntent(helper.JitterKey, currentpathResult)
	default:
		return fmt.Errorf("Unknown intent type: %s", intentType)
	}
}

func (manager *CalculationManager) validateCurrentResultForIntent(weightType helper.WeightKey, currentpathResult domain.PathResult) error {
	cost, err := manager.checkifPathValid(currentpathResult, weightType)
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
