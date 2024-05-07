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

type CalculationType string

const (
	CalculationTypeMin CalculationType = "min"
	CalculationTypeMax CalculationType = "max"
	CalculationTypeSum CalculationType = "sum"
)

type ShortestPathCalculation struct {
	log             *logrus.Entry
	graph           graph.Graph
	source          graph.Node
	destination     graph.Node
	weightType      helper.WeightKey
	calculationType CalculationType
	weights         map[interface{}]float64
	EdgeToPrevious  map[interface{}]graph.Edge
	priorityQueue   PriorityQueue
	visitedNodes    map[interface{}]bool
}

func NewShortestPathCalculation(network graph.Graph, source, destination graph.Node, weightType helper.WeightKey, calculationType CalculationType) *ShortestPathCalculation {
	return &ShortestPathCalculation{
		log:             logging.DefaultLogger.WithField("subsystem", "calculation"),
		graph:           network,
		source:          source,
		destination:     destination,
		weightType:      weightType,
		calculationType: calculationType,
		weights:         make(map[interface{}]float64),
		EdgeToPrevious:  make(map[interface{}]graph.Edge),
		visitedNodes:    make(map[interface{}]bool),
	}
}

func (calculation *ShortestPathCalculation) Execute() (graph.Path, error) {
	calculation.initializeDijkstra()
	if err := calculation.performDijkstra(); err != nil {
		return nil, err
	}
	return calculation.reconstructPath()
}

func (calculation *ShortestPathCalculation) initializeDijkstra() {
	var initialNodeCost float64
	var sourceNodeCost float64
	if calculation.calculationType == CalculationTypeMax {
		initialNodeCost = 0
		sourceNodeCost = math.Inf(1)
		calculation.priorityQueue = *NewMaximumPriorityQueue()
	} else if calculation.calculationType == CalculationTypeMin {
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
	if calculation.weightType == helper.PacketLossKey {
		weightPercentage := weight / 100
		return calculation.weights[currentId] + -math.Log(1-weightPercentage)
	}
	return calculation.weights[currentId] + weight
}

func (calculation *ShortestPathCalculation) handleMaxCalculation(currentId interface{}, weight float64, neighborId interface{}, edge graph.Edge) error {
	if calculation.weightType != helper.AvailableBandwidthKey {
		return fmt.Errorf("CalculationTypeMax is only supported for High Bandwidth Path Calculation")
	}
	if _, ok := calculation.visitedNodes[neighborId]; !ok {
		minimum := math.Min(calculation.weights[currentId], weight)
		if minimum > calculation.weights[neighborId] {
			calculation.updateMetricsAndPrevious(neighborId, minimum, edge)
		}
	}
	return nil
}

func (calculation *ShortestPathCalculation) handleMinCalculation(currentId interface{}, weight float64, neighborId interface{}, edge graph.Edge) error {
	if calculation.weightType != helper.AvailableBandwidthKey {
		return fmt.Errorf("CalculationTypeMin is only supported for Low Bandwidth Path Calculation")
	}
	if _, ok := calculation.visitedNodes[neighborId]; !ok {
		minimum := math.Min(calculation.weights[currentId], weight)
		if minimum < calculation.weights[neighborId] {
			calculation.updateMetricsAndPrevious(neighborId, minimum, edge)
		}
	}
	return nil
}

func (calculation *ShortestPathCalculation) handleDefaultCalculation(currentId interface{}, weight float64, neighborId interface{}, edge graph.Edge) {
	if _, ok := calculation.visitedNodes[neighborId]; !ok {
		alternativeDistance := calculation.calculateAlternativeDistance(currentId, weight)
		if alternativeDistance < calculation.weights[neighborId] {
			calculation.updateMetricsAndPrevious(neighborId, alternativeDistance, edge)
		}
	}
}

func (calculation *ShortestPathCalculation) relaxEdge(currentId interface{}, edge graph.Edge) error {
	neighbor := edge.To()
	weight := edge.GetWeight(calculation.weightType)
	neighborId := neighbor.GetId()

	if calculation.calculationType == CalculationTypeMax {
		return calculation.handleMaxCalculation(currentId, weight, neighborId, edge)
	} else if calculation.calculationType == CalculationTypeMin {
		return calculation.handleMinCalculation(currentId, weight, neighborId, edge)
	} else {
		calculation.handleDefaultCalculation(currentId, weight, neighborId, edge)
	}
	return nil
}

func (calculation *ShortestPathCalculation) performDijkstra() error {
	for !calculation.priorityQueue.IsEmpty() {
		item := heap.Pop(&calculation.priorityQueue).(*Item)
		calculation.visitedNodes[item.GetNodeId()] = true
		currentId := item.GetNodeId()
		if currentId == calculation.destination.GetId() {
			break
		}
		currentNode, _ := calculation.graph.GetNode(currentId)
		for _, edge := range currentNode.GetEdges() {
			if err := calculation.relaxEdge(currentId, edge); err != nil {
				return err
			}
		}
	}
	return nil
}

func (calculation ShortestPathCalculation) reconstructPath() (graph.Path, error) {
	path := make([]graph.Edge, 0)
	current := calculation.destination
	if calculation.calculationType == CalculationTypeSum {
		var totalCost float64
		if calculation.weightType != helper.PacketLossKey {
			totalCost = 0
		} else {
			totalCost = 1
		}
		for current.GetId() != calculation.source.GetId() {
			edge := calculation.EdgeToPrevious[current.GetId()]
			path = append([]graph.Edge{edge}, path...)
			cost := edge.GetWeight(calculation.weightType)
			if calculation.weightType != helper.PacketLossKey {
				totalCost += cost
			} else {
				totalCost *= 1 - cost
			}
			current = edge.From()
		}
		if calculation.weightType == helper.PacketLossKey {
			totalCost = 1 - totalCost
		}
		if len(path) == 0 {
			return nil, fmt.Errorf("No path found from node %d to node %d", calculation.source.GetId(), calculation.destination.GetId())
		}
		return graph.NewShortestPath(path, totalCost), nil
	}
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
	calculation.log.Debugf("Minimum Bandwidth Edge found with %v bandwidth %g: ", minEdge, minEdgeBandwidth)

	if len(path) == 0 {
		return nil, fmt.Errorf("No path found from node %d to node %d", calculation.source.GetId(), calculation.destination.GetId())
	}
	return graph.NewShortestPath(path, minEdgeBandwidth), nil
}
