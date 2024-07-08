package calculation

import (
	"math"

	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
)

type ServiceFunctionChainCalculation struct {
	BaseCalculation
	serviceFunctionChain [][]string
	routerServiceMap     map[string]string
}

func NewServiceFunctionChainCalculation(network graph.Graph, source graph.Node, destination graph.Node, weightTypes []helper.WeightKey, calculationMode CalculationMode, maxConstraints, minConstraints map[helper.WeightKey]float64, serviceFunctionChain [][]string, routerServiceMap map[string]string) *ServiceFunctionChainCalculation {
	return &ServiceFunctionChainCalculation{
		BaseCalculation:      *NewBaseCalculation(network, source, destination, weightTypes, calculationMode, maxConstraints, minConstraints),
		serviceFunctionChain: serviceFunctionChain,
		routerServiceMap:     routerServiceMap,
	}
}

func (calculation *ServiceFunctionChainCalculation) calculatePathSourceToFirstService(firstServiceRouterId string) (graph.Path, error) {
	firstServiceNode := calculation.graph.GetNode(firstServiceRouterId)
	calculation.log.Debugf("Calculating path from source node %s to first service router %s", calculation.source.GetName(), firstServiceNode.GetName())
	firstCalculation := NewShortestPathCalculation(calculation.graph, calculation.source, firstServiceNode, calculation.weightTypes, calculation.calculationMode, calculation.maxConstraints, calculation.minConstraints)
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
		serviceCalculation := NewShortestPathCalculation(calculation.graph, sourceNode, destinationNode, calculation.weightTypes, calculation.calculationMode, calculation.maxConstraints, calculation.minConstraints)
		serviceCalculation.SetInitialSourceNodeMetrics(previousPath.GetTotalCost(), previousPath.GetTotalDelay(), previousPath.GetTotalJitter(), previousPath.GetTotalPacketLoss())
		path, err := serviceCalculation.Execute()
		cost += path.GetTotalCost()
		if err != nil {
			return nil, 0, err
		}
		paths = append(paths, path)
		previousPath = path
	}
	return paths, cost, nil
}

func (calculation *ServiceFunctionChainCalculation) calculatePathLastServiceToDestination(previousPath graph.Path, lastServiceRouterId string) (graph.Path, error) {
	lastServiceNode := calculation.graph.GetNode(lastServiceRouterId)
	calculation.log.Debugf("Calculating path from last service router %s to destination node %s", lastServiceNode.GetName(), calculation.destination.GetName())
	lastCalculation := NewShortestPathCalculation(calculation.graph, lastServiceNode, calculation.destination, calculation.weightTypes, calculation.calculationMode, calculation.maxConstraints, calculation.minConstraints)
	lastCalculation.SetInitialSourceNodeMetrics(previousPath.GetTotalCost(), previousPath.GetTotalDelay(), previousPath.GetTotalJitter(), previousPath.GetTotalPacketLoss())
	path, err := lastCalculation.Execute()
	return path, err
}

func (calculation *ServiceFunctionChainCalculation) processServiceFunctionChain(serviceFunctionChain []string, bestCost float64) (float64, []graph.Path, error) {
	var paths []graph.Path

	firstPath, err := calculation.calculatePathSourceToFirstService(serviceFunctionChain[0])
	if err != nil {
		return 0, nil, err
	}
	cost := firstPath.GetTotalCost()
	if bestCost != 0 && cost > bestCost {
		return cost, nil, nil
	}
	paths = append(paths, firstPath)

	intermediatePaths, intermediateCost, err := calculation.calculatePathsBetweenServices(firstPath, serviceFunctionChain)
	if err != nil {
		return 0, nil, err
	}
	cost = intermediateCost
	if bestCost != 0 && cost > bestCost {
		return cost, nil, nil
	}
	paths = append(paths, intermediatePaths...)

	lastPath, err := calculation.calculatePathLastServiceToDestination(paths[len(paths)-1], serviceFunctionChain[len(serviceFunctionChain)-1])
	if err != nil {
		return 0, nil, err
	}
	cost = lastPath.GetTotalCost()
	if bestCost != 0 && cost > bestCost {
		return cost, nil, nil
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

func (calculation *ServiceFunctionChainCalculation) Execute() (graph.Path, error) {
	bestCost := math.Inf(1)
	var bestPath []graph.Path
	routerServiceMap := make(map[string]string)

	for index, serviceFunctionChain := range calculation.serviceFunctionChain {
		calculation.log.Debugf("Processing %d. service function chain", index+1)
		cost, paths, err := calculation.processServiceFunctionChain(serviceFunctionChain, bestCost)
		if err != nil {
			return nil, err
		}
		if len(paths) > 0 && cost < bestCost {
			bestCost = cost
			bestPath = paths
			for _, serviceRouterId := range serviceFunctionChain {
				routerServiceMap[serviceRouterId] = calculation.routerServiceMap[serviceRouterId]
			}
			calculation.log.Debugf("Found new best path with cost %f", cost)
		} else {
			calculation.log.Debugf("Path with cost %f is not better than the best path with cost %f", cost, bestCost)
		}
	}
	path := calculation.createPathFromSubPaths(bestPath)
	path.SetRouterServiceMap(routerServiceMap)
	return path, nil
}
