package calculation

import (
	"fmt"
	"math"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type CalculationManager struct {
	log              *logrus.Entry
	cache            cache.Cache
	graph            graph.Graph
	calculationSetup CalculationSetup
}

func NewCalculationManager(cache cache.Cache, graph graph.Graph, calculationSetup CalculationSetup) *CalculationManager {
	return &CalculationManager{
		log:              logging.DefaultLogger.WithField("subsystem", subsystem),
		cache:            cache,
		graph:            graph,
		calculationSetup: calculationSetup,
	}
}

func (manager *CalculationManager) getNodesFromPath(path graph.Path) []string {
	nodeList := make([]string, 0)
	for _, edge := range path.GetEdges() {
		to := edge.To().GetId()
		nodeList = append(nodeList, to)
	}
	return nodeList
}

func (manager *CalculationManager) translatePathToSidList(path graph.Path, algorithm uint32) ([]string, []string) {
	nodeList := manager.getNodesFromPath(path)
	serviceSidList := make([]string, 0)
	routerServiceMap := path.GetRouterServiceMap()
	manager.log.Debugln("Node in Path: ", nodeList)
	var sidList []string
	for _, node := range nodeList {
		sid := manager.cache.GetSrAlgorithmSid(node, algorithm)
		if sid == "" {
			manager.log.Errorln("SID not found for router: ", node)
			continue
		}
		sidList = append(sidList, sid)
		if serviceSid, ok := routerServiceMap[node]; ok {
			sidList = append(sidList, serviceSid)
			serviceSidList = append(serviceSidList, serviceSid)
		}
	}
	manager.log.Debugln("Translated SID List: ", sidList)
	return sidList, serviceSidList
}

func (manager *CalculationManager) createPathResult(path graph.Path, pathRequest domain.PathRequest, algorithm uint32) domain.PathResult {
	var sidList []string
	var serviceSidList []string
	if path == nil {
		manager.log.Errorln("No path found, return destination IPv6 address as SID list")
		sidList = []string{pathRequest.GetIpv6DestinationAddress()}
	} else {
		sidList, serviceSidList = manager.translatePathToSidList(path, algorithm)
	}
	pathResult, err := domain.NewDomainPathResult(pathRequest, path, sidList)
	if serviceSidList != nil {
		pathResult.SetServiceSidList(serviceSidList)
	}
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

func (manager *CalculationManager) getGraphAndAlgorithm(graph graph.Graph, firstIntent domain.Intent) (graph.Graph, uint32) {
	if firstIntent.GetIntentType() != domain.IntentTypeFlexAlgo {
		return graph, 0
	} else {
		algorithm := uint32(firstIntent.GetValues()[0].GetNumberValue())
		return graph.GetSubGraph(algorithm), algorithm
	}
}

func (manager *CalculationManager) CalculateBestPath(pathRequest domain.PathRequest) (domain.PathResult, error) {
	manager.lockElements()
	defer manager.unlockElements()

	calculationOptions, err := manager.calculationSetup.PerformSetup(pathRequest)
	intents := pathRequest.GetIntents()

	var (
		algorithm            uint32
		calculation          Calculation
		serviceFunctionChain [][]string
		routerServiceMap     map[string]string
	)

	intentToUse := intents[0]
	if len(intents) > 1 && intents[0].GetIntentType() == domain.IntentTypeSFC {
		intentToUse = intents[1]
	}

	calculationOptions.graph, algorithm = manager.getGraphAndAlgorithm(manager.graph, intentToUse)

	if intentToUse.GetIntentType() == domain.IntentTypeSFC {
		sfcCalculationOptions, err := manager.calculationSetup.PerformServiceFunctionChainSetup(intentToUse)
		if err != nil {
			return nil, err
		}
		manager.log.Debugln("Service Function Chain: ", serviceFunctionChain)
		manager.log.Debugln("Router Service Map: ", routerServiceMap)
		calculation = NewServiceFunctionChainCalculation(calculationOptions, sfcCalculationOptions)
	} else {
		calculation = NewShortestPathCalculation(calculationOptions)
	}

	path, err := calculation.Execute()
	if err != nil {
		return nil, err
	}
	return manager.createPathResult(path, pathRequest, algorithm), nil
}

func (manager *CalculationManager) calculateTotalCost(pathResult domain.PathResult, weightTypes []helper.WeightKey) error {
	var newTotalCost float64
	if len(weightTypes) == 1 && weightTypes[0] == helper.PacketLossKey {
		newTotalCost = 1.0
	} else {
		newTotalCost = 0.0
	}
	for _, edge := range pathResult.GetEdges() {
		graphEdge := manager.graph.GetEdge(edge.GetId())
		if graphEdge == nil {
			return fmt.Errorf("Path not valid, edge not found in graph: %s", edge.GetId())
		}
		if len(weightTypes) == 1 {
			if weightTypes[0] == helper.PacketLossKey {
				newTotalCost *= 1 - graphEdge.GetWeight(weightTypes[0])
			} else {
				newTotalCost += graphEdge.GetWeight(weightTypes[0])
			}
		} else if len(weightTypes) == 2 {
			newTotalCost += edge.GetWeight(weightTypes[0])*float64(helper.TwoFactorWeights[0]) + edge.GetWeight(weightTypes[1])*float64(helper.TwoFactorWeights[1])
		} else {
			newTotalCost += edge.GetWeight(weightTypes[0])*float64(helper.ThreeFactorWeights[0]) + edge.GetWeight(weightTypes[1])*float64(helper.ThreeFactorWeights[1]) + edge.GetWeight(weightTypes[2])*float64(helper.ThreeFactorWeights[2])
		}
	}
	if len(weightTypes) == 1 && weightTypes[0] == helper.PacketLossKey {
		newTotalCost = (1 - newTotalCost) * 100
	}
	currentTotalCost := pathResult.GetTotalCost()
	if currentTotalCost != newTotalCost {
		manager.log.Debugf("Total cost of current applied path changed from %f to %f", currentTotalCost, newTotalCost)
		pathResult.SetTotalCost(newTotalCost)
	}
	return nil
}

func (manager *CalculationManager) calculateMinimumValue(pathResult domain.PathResult, weightType helper.WeightKey) error {
	minValue := math.Inf(1)
	var bottleneckEdge graph.Edge
	edges := pathResult.GetEdges()
	for _, edge := range edges {
		graphEdge := manager.graph.GetEdge(edge.GetId())
		if graphEdge == nil {
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

func (manager *CalculationManager) updateCurrentResult(weightKeys []helper.WeightKey, calculationMode CalculationMode, currentPathResult domain.PathResult) error {
	var err error
	if calculationMode == CalculationModeSum {
		err = manager.calculateTotalCost(currentPathResult, weightKeys)
	} else {
		err = manager.calculateMinimumValue(currentPathResult, weightKeys[0])
	}
	if err != nil {
		return fmt.Errorf("Current Path is not valid anymore, new path will be applied: %s", err)
	}
	return nil
}

func (manager *CalculationManager) handleSumCalculationMode(currentPathResult, newPathResult domain.PathResult, streamSession domain.StreamSession) (*domain.PathResult, error) {
	newPathTotalCost := newPathResult.GetTotalCost()
	currentPathTotalCost := currentPathResult.GetTotalCost()
	if newPathTotalCost < currentPathTotalCost*(1-helper.FlappingThreshold) {
		manager.log.Debugf("New path will be applied, cost of new path is by more than 10 percent smaller, current: %f to new: %f", currentPathTotalCost, newPathTotalCost)
		streamSession.SetPathResult(newPathResult)
		return &newPathResult, nil
	} else {
		manager.log.Debugf("No path changes, cost of new path is not smaller by more than 10 percent, current: %f to new: %f", currentPathTotalCost, newPathTotalCost)
	}
	return nil, nil
}

func (manager *CalculationManager) handleNonSumCalculationMode(currentPathResult, newPathResult domain.PathResult, streamSession domain.StreamSession) (*domain.PathResult, error) {
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
	return nil, nil
}

func (manager *CalculationManager) AreServicesStillValid(serviceSidList []string) bool {
	manager.log.Debugln("Check if current services are still valid")
	for _, sid := range serviceSidList {
		if !manager.cache.DoesServiceSidExist(sid) {
			manager.log.Debugf("Service SID %s not found in cache, new path will be applied", sid)
			return false
		}
	}
	return true
}

func (manager *CalculationManager) handlePathChange(weightKeys []helper.WeightKey, calculationMode CalculationMode, currentPathResult, newPathResult domain.PathResult, streamSession domain.StreamSession) (*domain.PathResult, error) {
	manager.log.Debugln("Better Path found, check for applicability")
	manager.log.Debugln("Validate current path and its cost")

	firstIntent := streamSession.GetPathRequest().GetIntents()[0]
	if firstIntent.GetIntentType() == domain.IntentTypeSFC {
		if !manager.AreServicesStillValid(currentPathResult.GetServiceSidList()) {
			streamSession.SetPathResult(newPathResult)
			return &newPathResult, nil
		}
	}

	if err := manager.updateCurrentResult(weightKeys, calculationMode, currentPathResult); err != nil {
		manager.log.Errorln("Current Path is not valid anymore, new path will be applied: ", err)
		streamSession.SetPathResult(newPathResult)
		return &newPathResult, nil
	}

	manager.log.Debugln("Current path is still valid, check if new path is better")
	if calculationMode == CalculationModeSum {
		return manager.handleSumCalculationMode(currentPathResult, newPathResult, streamSession)
	} else {
		return manager.handleNonSumCalculationMode(currentPathResult, newPathResult, streamSession)
	}
}

func (manager *CalculationManager) CalculatePathUpdate(streamSession domain.StreamSession) (*domain.PathResult, error) {
	// TODO find solution for update
	// currentPathResult := streamSession.GetPathResult()
	// currentAppliedSidList := currentPathResult.GetIpv6SidAddresses()
	// manager.log.Debugln("SID list of current path is", currentAppliedSidList)

	// pathRequest := streamSession.GetPathRequest()
	// intents := pathRequest.GetIntents()
	// manager.log.Debugln("Recalculate path with new network state")

	// newPathResult, err := manager.CalculateBestPath(pathRequest)
	// if err != nil {
	// 	return nil, err
	// }

	// weightKeys, calculationMode := manager.getWeightKeysandCalculationType(intents)
	// if calculationMode == CalculationModeUndefined {
	// 	return nil, fmt.Errorf("Calculation mode not defined for intents")
	// }

	// if !reflect.DeepEqual(newPathResult.GetIpv6SidAddresses(), currentAppliedSidList) {
	// 	return manager.handlePathChange(weightKeys, calculationMode, currentPathResult, newPathResult, streamSession)
	// } else {
	// 	manager.log.Debugln("No changes in path detected, update current path with new path cost")
	// 	if err := manager.updateCurrentResult(weightKeys, calculationMode, currentPathResult); err != nil {
	// 		return nil, err
	// 	}
	// }
	manager.log.Debugln("No path changes, current path is still valid")
	return nil, nil
}
