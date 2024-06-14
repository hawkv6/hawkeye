package calculation

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type CalculationMode int

const (
	CalculationModeUndefined CalculationMode = iota
	CalculationModeSum
	CalculationModeMin
	CalculationModeMax
)

type ShortestPathCalculation struct {
	log             *logrus.Entry
	graph           graph.Graph
	source          graph.Node
	destination     graph.Node
	weightTypes     []helper.WeightKey
	calculationType CalculationMode
	weights         map[interface{}]float64
	EdgeToPrevious  map[interface{}]graph.Edge
	priorityQueue   PriorityQueue
	visitedNodes    map[interface{}]bool
}

func NewShortestPathCalculation(network graph.Graph, source, destination graph.Node, weightTypes []helper.WeightKey, calculationType CalculationMode) *ShortestPathCalculation {
	return &ShortestPathCalculation{
		log:             logging.DefaultLogger.WithField("subsystem", "calculation"),
		graph:           network,
		source:          source,
		destination:     destination,
		weightTypes:     weightTypes,
		calculationType: calculationType,
		weights:         make(map[interface{}]float64),
		EdgeToPrevious:  make(map[interface{}]graph.Edge),
		visitedNodes:    make(map[interface{}]bool),
	}
}

func (calculation *ShortestPathCalculation) Execute() (graph.Path, error) {
	calculation.initializeDijkstra()
	calculation.performDijkstra()
	return calculation.reconstructPath()
}

func (calculation *ShortestPathCalculation) initializeDijkstra() {
	var initialNodeCost float64
	var sourceNodeCost float64
	if calculation.calculationType == CalculationModeMax {
		initialNodeCost = 0
		sourceNodeCost = math.Inf(1)
		calculation.priorityQueue = *NewMaximumPriorityQueue()
	} else if calculation.calculationType == CalculationModeMin {
		initialNodeCost = math.Inf(1)
		sourceNodeCost = math.Inf(1)
		calculation.priorityQueue = *NewMinimumPriorityQueue()
	} else {
		initialNodeCost = math.Inf(1)
		sourceNodeCost = 0
		calculation.priorityQueue = *NewMinimumPriorityQueue()
	}
	for id := range calculation.graph.GetNodes() {
		calculation.weights[id] = initialNodeCost
	}
	heap.Init(&calculation.priorityQueue)
	calculation.weights[calculation.source.GetId()] = sourceNodeCost
	heap.Push(&calculation.priorityQueue, &Item{nodeId: calculation.source.GetId(), cost: calculation.weights[calculation.source.GetId()], index: 0})
}

func (calculation *ShortestPathCalculation) updateMetricsAndPrevious(neighborId interface{}, value float64, edge graph.Edge) {
	calculation.weights[neighborId] = value
	calculation.EdgeToPrevious[neighborId] = edge
	heap.Push(&calculation.priorityQueue, &Item{nodeId: neighborId, cost: value})
}

func (calculation *ShortestPathCalculation) calculateAlternativeDistance(currentId interface{}, weight float64) float64 {
	if calculation.weightTypes[0] == helper.PacketLossKey {
		weightPercentage := weight / 100
		return calculation.weights[currentId] + -math.Log(1-weightPercentage)
	}
	return calculation.weights[currentId] + weight
}

func (calculation *ShortestPathCalculation) handleMaxCalculation(currentId interface{}, weight float64, neighborId interface{}, edge graph.Edge) {
	if _, ok := calculation.visitedNodes[neighborId]; !ok {
		minimum := math.Min(calculation.weights[currentId], weight)
		if minimum > calculation.weights[neighborId] {
			calculation.updateMetricsAndPrevious(neighborId, minimum, edge)
		}
	}
}

func (calculation *ShortestPathCalculation) handleMinCalculation(currentId interface{}, weight float64, neighborId interface{}, edge graph.Edge) {
	if _, ok := calculation.visitedNodes[neighborId]; !ok {
		minimum := math.Min(calculation.weights[currentId], weight)
		if minimum < calculation.weights[neighborId] {
			calculation.updateMetricsAndPrevious(neighborId, minimum, edge)
		}
	}
}

func (calculation *ShortestPathCalculation) handleDefaultCalculation(currentId interface{}, weight float64, neighborId interface{}, edge graph.Edge) {
	if _, ok := calculation.visitedNodes[neighborId]; !ok {
		alternativeDistance := calculation.calculateAlternativeDistance(currentId, weight)
		if alternativeDistance < calculation.weights[neighborId] {
			calculation.updateMetricsAndPrevious(neighborId, alternativeDistance, edge)
		}
	}
}

func (calculation *ShortestPathCalculation) getWeight(edge graph.Edge) float64 {
	weight := 0.0
	if len(calculation.weightTypes) == 2 {
		weight = edge.GetWeight(calculation.weightTypes[0])*float64(helper.TwoFactorWeights[0]) + edge.GetWeight(calculation.weightTypes[1])*float64(helper.TwoFactorWeights[1])

	} else if len(calculation.weightTypes) == 3 {
		weight = edge.GetWeight(calculation.weightTypes[0])*float64(helper.ThreeFactorWeights[0]) + edge.GetWeight(calculation.weightTypes[1])*float64(helper.ThreeFactorWeights[1]) + edge.GetWeight(calculation.weightTypes[2])*float64(helper.ThreeFactorWeights[2])
	} else {
		weight = edge.GetWeight(calculation.weightTypes[0])
	}
	return weight
}

func (calculation *ShortestPathCalculation) relaxEdge(currentId interface{}, edge graph.Edge) {
	neighbor := edge.To()
	neighborId := neighbor.GetId()
	weight := calculation.getWeight(edge)

	if calculation.calculationType == CalculationModeMax {
		calculation.handleMaxCalculation(currentId, weight, neighborId, edge)
	} else if calculation.calculationType == CalculationModeMin {
		calculation.handleMinCalculation(currentId, weight, neighborId, edge)
	} else {
		calculation.handleDefaultCalculation(currentId, weight, neighborId, edge)
	}
}

func (calculation *ShortestPathCalculation) performDijkstra() {
	for !calculation.priorityQueue.IsEmpty() {
		item := heap.Pop(&calculation.priorityQueue).(*Item)
		calculation.visitedNodes[item.GetNodeId()] = true
		currentId := item.GetNodeId()
		if currentId == calculation.destination.GetId() {
			break
		}
		currentNode, _ := calculation.graph.GetNode(currentId)
		for _, edge := range currentNode.GetEdges() {
			calculation.relaxEdge(currentId, edge)
		}
	}
}

func (calculation *ShortestPathCalculation) reconstructSumPath(current graph.Node, path []graph.Edge) (graph.Path, error) {
	costMap := make(map[helper.WeightKey]float64, len(path))

	for _, weightType := range calculation.weightTypes {
		if _, ok := costMap[weightType]; !ok {
			costMap[weightType] = 0
		}
	}
	costMap[helper.LatencyKey] = 0
	costMap[helper.JitterKey] = 0
	costMap[helper.PacketLossKey] = 1

	totalCost := 0.0
	for current.GetId() != calculation.source.GetId() {
		edge := calculation.EdgeToPrevious[current.GetId()]
		path = append([]graph.Edge{edge}, path...)
		costMap[helper.LatencyKey] += edge.GetWeight(helper.LatencyKey)
		costMap[helper.JitterKey] += edge.GetWeight(helper.JitterKey)
		costMap[helper.PacketLossKey] *= edge.GetWeight(helper.PacketLossKey)
		if len(calculation.weightTypes) == 2 {
			totalCost += edge.GetWeight(calculation.weightTypes[0])*float64(helper.TwoFactorWeights[0]) + edge.GetWeight(calculation.weightTypes[1])*float64(helper.TwoFactorWeights[1])
		} else if len(calculation.weightTypes) == 3 {
			totalCost += edge.GetWeight(calculation.weightTypes[0])*float64(helper.ThreeFactorWeights[0]) + edge.GetWeight(calculation.weightTypes[1])*float64(helper.ThreeFactorWeights[1]) + edge.GetWeight(calculation.weightTypes[2])*float64(helper.ThreeFactorWeights[2])
		}
		current = edge.From()
	}
	if len(calculation.weightTypes) == 1 {
		if calculation.weightTypes[0] == helper.PacketLossKey {
			totalCost = 1 - costMap[helper.PacketLossKey]
		} else {
			totalCost = costMap[calculation.weightTypes[0]]
		}
	}

	if len(path) == 0 {
		return nil, fmt.Errorf("No path found from node %d to node %d", calculation.source.GetId(), calculation.destination.GetId())
	}
	return graph.NewShortestPathWithTotalCost(path, totalCost, costMap[helper.LatencyKey], costMap[helper.JitterKey], costMap[helper.PacketLossKey]), nil
}

func (calculation *ShortestPathCalculation) reconstructMinPath(current graph.Node, path []graph.Edge) (graph.Path, error) {
	var minEdge graph.Edge
	minEdgeBandwidth := math.Inf(1)
	for current.GetId() != calculation.source.GetId() {
		edge := calculation.EdgeToPrevious[current.GetId()]
		if edge.GetWeight(helper.AvailableBandwidthKey) < minEdgeBandwidth {
			minEdgeBandwidth = edge.GetWeight(helper.AvailableBandwidthKey)
			minEdge = edge
		}
		path = append([]graph.Edge{edge}, path...)
		current = edge.From()
	}
	if len(path) == 0 {
		return nil, fmt.Errorf("No path found from node %d to node %d", calculation.source.GetId(), calculation.destination.GetId())
	}
	calculation.log.Debugf("First bottleneck edge: %v with bandwidth %g: ", minEdge, minEdgeBandwidth)
	return graph.NewShortestPathWithBottleneck(path, minEdge, minEdgeBandwidth), nil
}

func (calculation ShortestPathCalculation) reconstructPath() (graph.Path, error) {
	path := make([]graph.Edge, 0)
	if calculation.calculationType == CalculationModeSum {
		return calculation.reconstructSumPath(calculation.destination, path)
	}
	return calculation.reconstructMinPath(calculation.destination, path)
}
