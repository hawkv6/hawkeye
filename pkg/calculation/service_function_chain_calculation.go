package calculation

import (
	"fmt"
	"math"

	"github.com/hawkv6/hawkeye/pkg/graph"
)

type ServiceFunctionChainCalculation struct {
	BaseCalculation
	serviceFunctionChain [][]string
	routerServiceMap     map[string]string
}

func NewServiceFunctionChainCalculation(options *CalculationOptions, sfcCalculationOptions *SfcCalculationOptions) *ServiceFunctionChainCalculation {
	return &ServiceFunctionChainCalculation{
		BaseCalculation:      *NewBaseCalculation(options),
		serviceFunctionChain: sfcCalculationOptions.serviceFunctionChain,
		routerServiceMap:     sfcCalculationOptions.routerServiceMap,
	}
}

func (calculation *ServiceFunctionChainCalculation) calculatePathSourceToFirstService(firstServiceRouterId string) (graph.Path, error) {
	firstServiceNode := calculation.graph.GetNode(firstServiceRouterId)
	calculation.log.Debugf("Calculating path from source node %s to first service router %s", calculation.source.GetName(), firstServiceNode.GetName())
	calculationOptions := &CalculationOptions{calculation.graph, calculation.source, firstServiceNode, calculation.weightKeys, calculation.calculationMode, calculation.maxConstraints, calculation.minConstraints}
	firstCalculation := NewShortestPathCalculation(calculationOptions)
	path, err := firstCalculation.Execute()
	return path, err
}

func (calculation *ServiceFunctionChainCalculation) calculatePathsBetweenServices(previousPath graph.Path, serviceFunctionChain []string) ([]graph.Path, float64, error) {
	calculation.log.Debugln("Calculating paths between services")
	cost := float64(0)
	paths := make([]graph.Path, 0)
	for i := 0; i < len(serviceFunctionChain)-1; i++ {
		sourceNode := calculation.graph.GetNode(serviceFunctionChain[i])
		destinationNode := calculation.graph.GetNode(serviceFunctionChain[i+1])
		calculation.log.Debugf("Calculating path from service router %s to service router %s", sourceNode.GetName(), destinationNode.GetName())
		calculationOptions := &CalculationOptions{calculation.graph, sourceNode, destinationNode, calculation.weightKeys, calculation.calculationMode, calculation.maxConstraints, calculation.minConstraints}
		serviceCalculation := NewShortestPathCalculation(calculationOptions)
		serviceCalculation.SetInitialSourceNodeMetrics(previousPath.GetTotalCost(), previousPath.GetTotalDelay(), previousPath.GetTotalJitter(), previousPath.GetTotalPacketLoss())
		path, err := serviceCalculation.Execute()
		if err != nil {
			return nil, 0, err
		}
		cost += path.GetTotalCost()
		paths = append(paths, path)
		previousPath = path
	}
	return paths, cost, nil
}

func (calculation *ServiceFunctionChainCalculation) calculatePathLastServiceToDestination(previousPath graph.Path, lastServiceRouterId string) (graph.Path, error) {
	lastServiceNode := calculation.graph.GetNode(lastServiceRouterId)
	calculation.log.Debugf("Calculating path from last service router %s to destination node %s", lastServiceNode.GetName(), calculation.destination.GetName())
	calculationOptions := &CalculationOptions{calculation.graph, lastServiceNode, calculation.destination, calculation.weightKeys, calculation.calculationMode, calculation.maxConstraints, calculation.minConstraints}
	lastCalculation := NewShortestPathCalculation(calculationOptions)
	lastCalculation.SetInitialSourceNodeMetrics(previousPath.GetTotalCost(), previousPath.GetTotalDelay(), previousPath.GetTotalJitter(), previousPath.GetTotalPacketLoss())
	path, err := lastCalculation.Execute()
	return path, err
}

func (calculation *ServiceFunctionChainCalculation) evaluateFistSubPath(serviceFunctionChain []string, bestCost float64, cost *float64) (graph.Path, bool, error) {
	firstPath, err := calculation.calculatePathSourceToFirstService(serviceFunctionChain[0])
	if err != nil {
		return nil, true, err
	}
	*cost = firstPath.GetTotalCost()
	if *cost > bestCost {
		return nil, true, nil
	}
	return firstPath, false, nil
}

func (calculation *ServiceFunctionChainCalculation) evaluateSecondSubPath(firstPath graph.Path, serviceFunctionChain []string, cost *float64, bestCost float64) ([]graph.Path, bool, error) {
	intermediatePaths, intermediateCost, err := calculation.calculatePathsBetweenServices(firstPath, serviceFunctionChain)
	if err != nil {
		return nil, true, err
	}
	*cost = intermediateCost
	if *cost > bestCost {
		return nil, true, nil
	}
	return intermediatePaths, false, nil
}

func (calculation *ServiceFunctionChainCalculation) evaluateLastSubPath(lastIntermediatePath graph.Path, serviceFunctionChain []string, bestCost float64, cost *float64) (graph.Path, bool, error) {
	lastPath, err := calculation.calculatePathLastServiceToDestination(lastIntermediatePath, serviceFunctionChain[len(serviceFunctionChain)-1])
	if err != nil {
		return nil, true, err
	}
	*cost = lastPath.GetTotalCost()
	if *cost > bestCost {
		return nil, true, nil
	}
	return lastPath, false, nil
}

func (calculation *ServiceFunctionChainCalculation) processServiceFunctionChain(serviceFunctionChain []string, bestCost float64) (float64, []graph.Path, error) {
	paths := make([]graph.Path, 0)
	cost := float64(0)
	firstPath, shouldReturn, err := calculation.evaluateFistSubPath(serviceFunctionChain, bestCost, &cost)
	if shouldReturn {
		return cost, nil, err
	}
	paths = append(paths, firstPath)

	var intermediatePaths []graph.Path
	intermediatePaths, shouldReturn, err = calculation.evaluateSecondSubPath(firstPath, serviceFunctionChain, &cost, bestCost)
	if shouldReturn {
		return cost, nil, err
	}
	paths = append(paths, intermediatePaths...)

	lastPath, shouldReturn, err := calculation.evaluateLastSubPath(paths[len(paths)-1], serviceFunctionChain, bestCost, &cost)
	if shouldReturn {
		return cost, nil, err
	}
	paths = append(paths, lastPath)
	return cost, paths, nil
}

func (calculation *ServiceFunctionChainCalculation) createPathFromSubPaths(subPaths []graph.Path) graph.Path {
	edges := make([]graph.Edge, 0)
	var bottleneckEdge graph.Edge
	bottleneckValue := math.Inf(1)
	for _, subPath := range subPaths {
		edges = append(edges, subPath.GetEdges()...)
		if subPath.GetBottleneckValue() < bottleneckValue {
			bottleneckEdge = subPath.GetBottleneckEdge()
			bottleneckValue = subPath.GetBottleneckValue()
		}
	}
	lastPath := subPaths[len(subPaths)-1]
	return graph.NewShortestPath(edges, lastPath.GetTotalCost(), lastPath.GetTotalDelay(), lastPath.GetTotalJitter(), lastPath.GetTotalPacketLoss(), bottleneckValue, bottleneckEdge)
}

func (calculation *ServiceFunctionChainCalculation) updateBestPath(currentCost float64, bestCost *float64, currentPaths []graph.Path, bestSubPaths *[]graph.Path) bool {
	if len(currentPaths) > 0 && currentCost < *bestCost {
		*bestCost = currentCost
		*bestSubPaths = currentPaths
		calculation.log.Debugf("Found new best path with cost %f", currentCost)
		return true
	}
	calculation.log.Debugf("Path with cost %f is not better than the best path with cost %f", currentCost, *bestCost)
	return false
}

func (calculation *ServiceFunctionChainCalculation) updateRouterServiceMap(serviceFunctionChain []string, routerServiceMap map[string]string) {
	for _, serviceRouterId := range serviceFunctionChain {
		routerServiceMap[serviceRouterId] = calculation.routerServiceMap[serviceRouterId]
	}
}

func (calculation *ServiceFunctionChainCalculation) calculateBestServiceChain() ([]graph.Path, map[string]string, error) {
	bestCost := math.Inf(1)
	var bestSubPaths []graph.Path
	bestServiceFunctionChain := make([]string, 0)

	for index, serviceFunctionChain := range calculation.serviceFunctionChain {
		calculation.log.Debugf("Processing %d. service function chain", index+1)
		cost, paths, err := calculation.processServiceFunctionChain(serviceFunctionChain, bestCost)
		if err != nil {
			calculation.log.Errorf("Error processing service function chain: %s", err)
			continue
		}
		if calculation.updateBestPath(cost, &bestCost, paths, &bestSubPaths) {
			bestServiceFunctionChain = serviceFunctionChain
		}
	}
	if len(bestSubPaths) == 0 {
		return nil, nil, fmt.Errorf("No valid path for service function chain found")
	}

	routerServiceMap := make(map[string]string)
	calculation.updateRouterServiceMap(bestServiceFunctionChain, routerServiceMap)
	return bestSubPaths, routerServiceMap, nil
}

func (calculation *ServiceFunctionChainCalculation) Execute() (graph.Path, error) {
	bestSubPaths, routerServiceMap, err := calculation.calculateBestServiceChain()
	if err != nil {
		return nil, err
	}
	path := calculation.createPathFromSubPaths(bestSubPaths)
	path.SetRouterServiceMap(routerServiceMap)
	return path, nil
}
