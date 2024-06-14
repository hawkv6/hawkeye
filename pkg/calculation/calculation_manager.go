package calculation

import (
	"fmt"
	"math"
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

func (calcultor *CalculationManager) getNetworkAddress(ip string) (net.IP, error) {
	ipv6Address := net.ParseIP(ip)
	if ipv6Address == nil {
		return nil, fmt.Errorf("IPv6 Network address could not be parsed, invalid IPv6 address: %s", ip)
	}
	mask := net.CIDRMask(64, 128)
	return ipv6Address.Mask(mask), nil
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
	var sidList []string
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
	sourceIpv6, err := manger.getNetworkAddress(pathRequest.GetIpv6SourceAddress())
	if err != nil {
		return nil, err
	}
	sourceRouterId, ok := manger.cache.GetRouterIdFromNetworkAddress(sourceIpv6.String())
	if !ok {
		return nil, fmt.Errorf("Router ID not found for source IP: %s", sourceIpv6)
	}
	source, exist := manger.graph.GetNode(sourceRouterId)
	if !exist {
		return nil, fmt.Errorf("Source router not found")
	}
	return source, nil
}

func (manager *CalculationManager) GetDestinationNode(pathRequest domain.PathRequest) (graph.Node, error) {
	destinationIpv6, err := manager.getNetworkAddress(pathRequest.GetIpv6DestinationAddress())
	if err != nil {
		return nil, err
	}
	destinationRouterId, ok := manager.cache.GetRouterIdFromNetworkAddress(destinationIpv6.String())
	if !ok {
		return nil, fmt.Errorf("Router ID not found for destination IP: %s", destinationIpv6)
	}
	destination, exist := manager.graph.GetNode(destinationRouterId)
	if !exist {
		return nil, fmt.Errorf("Destination router not found")
	}
	return destination, nil
}

func (manager *CalculationManager) getWeightKeyAndCalcType(intentType domain.IntentType) (helper.WeightKey, CalculationMode) {
	switch intentType {
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

func (manager *CalculationManager) getNormalizedWeightKey(intentType domain.IntentType) helper.WeightKey {

	switch intentType {
	case domain.IntentTypeLowLatency:
		return helper.NormalizedLatencyKey
	case domain.IntentTypeLowJitter:
		return helper.NormalizedJitterKey
	case domain.IntentTypeLowPacketLoss:
		return helper.NormalizedPacketLossKey
	default:
		return helper.UndefinedKey
	}
}

func (manager *CalculationManager) createPathResult(path graph.Path, pathRequest domain.PathRequest) domain.PathResult {
	var sidList []string
	if path == nil {
		manager.log.Errorln("No path found, return destination IPv6 address as SID list")
		sidList = []string{pathRequest.GetIpv6DestinationAddress()}
	} else {
		sidList = manager.translatePathToSidList(path)
	}
	pathResult, err := domain.NewDomainPathResult(pathRequest, path, sidList)
	if err != nil {
		manager.log.Errorln("Error creating path result: ", err)
		return nil
	}
	return pathResult
}

func (manager *CalculationManager) lockElements() {
	manager.log.Debugln("Locking cache and graph mutexes")
	manager.graph.Lock()
	manager.cache.Lock()
}

func (manager *CalculationManager) unlockElements() {
	manager.log.Debugln("Unlocking cache and graph mutexes")
	manager.graph.Unlock()
	manager.cache.Unlock()
}

func (manager *CalculationManager) CalculateBestPath(pathRequest domain.PathRequest) (domain.PathResult, error) {
	manager.lockElements()
	defer manager.unlockElements()

	sourceNode, err := manager.getSourceNode(pathRequest)
	if err != nil {
		return nil, err
	}
	destinationNode, err := manager.GetDestinationNode(pathRequest)
	if err != nil {
		return nil, err
	}

	weightKeys := make([]helper.WeightKey, 0)
	var calculationType CalculationMode
	intents := pathRequest.GetIntents()
	if len(intents) == 0 {
		return nil, fmt.Errorf("No intents found in path request")
	} else if len(intents) == 1 {
		weightKey, calcType := manager.getWeightKeyAndCalcType(intents[0].GetIntentType())
		calculationType = calcType
		weightKeys = append(weightKeys, weightKey)
	} else {
		calculationType = CalculationModeSum
		for _, intent := range pathRequest.GetIntents() {
			weightKeys = append(weightKeys, manager.getNormalizedWeightKey(intent.GetIntentType()))
		}
	}
	if calculationType == CalculationModeUndefined {
		return nil, fmt.Errorf("Calculation mode not defined for intent type: %s", intents[0].GetIntentType())
	}
	calculation := NewShortestPathCalculation(manager.graph, sourceNode, destinationNode, weightKeys, calculationType)
	path, err := calculation.Execute()
	if err != nil {
		return nil, err

	} else {
		if calculationType == CalculationModeSum {
			manager.log.Debugf("Shortest Path found with total cost %g, delay: %gus, jitter %gus, packet loss %g%%: ", path.GetTotalCost(), path.GetTotalDelay(), path.GetTotalJitter(), path.GetTotalPacketLoss())
		} else {
			manager.log.Debugf("Shortest Path found with bottleneck %g: ", path.GetBottleneckValue())
			manager.log.Debugf("Bottleneck edge for this new path is: %v", path.GetBottleneckEdge())
		}
	}
	return manager.createPathResult(path, pathRequest), nil
}

func (manager *CalculationManager) calculateTotalCost(pathResult domain.PathResult, weightType helper.WeightKey) error {
	totalCost := 0.0
	edges := pathResult.GetEdges()
	for _, edge := range edges {
		graphEdge, exist := manager.graph.GetEdge(edge.GetId())
		if !exist {
			return fmt.Errorf("Path not valid, edge not found in graph: %s", edge.GetId())
		}
		cost := graphEdge.GetWeight(weightType)
		totalCost += cost
	}
	if pathResult.GetTotalCost() != totalCost {
		manager.log.Debugf("Total cost of current applied path changed from %f to %f", pathResult.GetTotalCost(), totalCost)
		pathResult.SetTotalCost(totalCost)
	}
	return nil
}

func (manager *CalculationManager) calculateMinimumValue(pathResult domain.PathResult, weightType helper.WeightKey) error {
	minValue := math.Inf(1)
	var bottleneckEdge graph.Edge
	edges := pathResult.GetEdges()
	for _, edge := range edges {
		graphEdge, exist := manager.graph.GetEdge(edge.GetId())
		if !exist {
			return fmt.Errorf("Path not valid, edge not found in graph: %s", edge.GetId())
		}
		cost := graphEdge.GetWeight(weightType)
		if cost < minValue {
			bottleneckEdge = graphEdge
			minValue = cost
		}
	}
	if pathResult.GetBottleneckEdge() != bottleneckEdge {
		manager.log.Debugf("Bottleneck edge of current applied path changed from %v to %v", pathResult.GetBottleneckEdge(), bottleneckEdge)
		manager.log.Debugf("Bottleneck value of current applied path changed from %f to %f", pathResult.GetBottleneckValue(), minValue)
		pathResult.SetBottleneckEdge(bottleneckEdge)
		pathResult.SetBottleneckValue(minValue)
	} else if pathResult.GetBottleneckValue() != minValue {
		manager.log.Debugf("Bottleneck value of current applied path changed from %f to %f", pathResult.GetBottleneckValue(), minValue)
		pathResult.SetBottleneckValue(minValue)
	}
	return nil
}

func (manager *CalculationManager) validateCurrentResult(currentpathResult domain.PathResult, weightKey helper.WeightKey, calcType CalculationMode) error {
	if calcType == CalculationModeSum {
		return manager.calculateTotalCost(currentpathResult, weightKey)
	}
	return manager.calculateMinimumValue(currentpathResult, weightKey)
}

func (manager *CalculationManager) CalculatePathUpdate(streamSession domain.StreamSession) (*domain.PathResult, error) {
	manager.lockElements()
	defer manager.unlockElements()

	currentPathResult := streamSession.GetPathResult()
	currentAppliedSidList := currentPathResult.GetIpv6SidAddresses()
	manager.log.Debugln("SID list of current path is", currentAppliedSidList)

	pathRequest := streamSession.GetPathRequest()
	intents := pathRequest.GetIntents()
	manager.log.Debugln("Recalculate path with new network state")

	newPathResult, err := manager.CalculateBestPath(pathRequest)
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(newPathResult.GetIpv6SidAddresses(), currentAppliedSidList) {
		manager.log.Debugln("Better Path found, check for applicability.")
		manager.log.Debugln("Validate current path and its cost")
		weightKey, calcType := manager.getWeightKeyAndCalcType(intents[0].GetIntentType())
		if err := manager.validateCurrentResult(currentPathResult, weightKey, calcType); err != nil { // TODO adapt if there several intents
			manager.log.Errorln("Current Path is not valid anymore, new path will be applied: ", err)
			streamSession.SetPathResult(newPathResult)
			return &newPathResult, nil
		}

		manager.log.Debugln("Current path is still valid, check if new path is better")
		if calcType == CalculationModeSum {
			newPathTotalCost := newPathResult.GetTotalCost()
			currentPathTotalCost := currentPathResult.GetTotalCost()
			if newPathTotalCost < currentPathTotalCost*(1-helper.FlappingThreshold) {
				manager.log.Debugf("New path will be applied, cost of new path is by more than 10 percent smaller, current: %f to new: %f", currentPathTotalCost, newPathTotalCost)
				streamSession.SetPathResult(newPathResult)
				return &newPathResult, nil
			} else {
				manager.log.Debugf("No path changes, cost of new path is not smaller by more than 10 percent, current: %f to new: %f", currentPathTotalCost, newPathTotalCost)
			}
		} else {
			newPathMinimumValue := newPathResult.GetBottleneckValue()
			newPathMinimumEdge := newPathResult.GetBottleneckEdge()
			currentPathMinimumValue := currentPathResult.GetBottleneckValue()
			currentPathMinimumEdge := currentPathResult.GetBottleneckEdge()
			if newPathMinimumValue < currentPathMinimumValue*(1-helper.FlappingThreshold) {
				manager.log.Debugf("Bottleneck in old path was %v with value %g: ", currentPathMinimumEdge, currentPathMinimumValue)
				manager.log.Debugf("Bottleneck in new path is %v with value %g: ", newPathMinimumEdge, newPathMinimumValue)
				manager.log.Debugf("New Path will be applied, bottleneck of new path is by more than 10 percent smaller, current: %f to new: %f", currentPathMinimumValue, newPathMinimumValue)
				streamSession.SetPathResult(newPathResult)
				return &newPathResult, nil
			} else {
				manager.log.Debugf("No path changes, cost of new path is not smaller by more than 10 percent, current: %f to new: %f", currentPathMinimumValue, newPathMinimumValue)
			}
		}
	} else {
		manager.log.Debugln("SID List is equal - no change in path found")
		weightKey, calcType := manager.getWeightKeyAndCalcType(intents[0].GetIntentType())
		if err := manager.validateCurrentResult(currentPathResult, weightKey, calcType); err != nil { // TODO adapt if there several intents
			return nil, fmt.Errorf("Current Path is not valid anymore, new path will be applied: %s", err)
		}
	}
	return nil, fmt.Errorf("No path changes, cost of new path is not smaller by more than 10 percent")
}
