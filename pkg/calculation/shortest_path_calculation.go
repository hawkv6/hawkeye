package calculation

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
)

type CalculationType string

const (
	CalculationTypeMin CalculationType = "min"
	CalculationTypeMax CalculationType = "max"
	CalculationTypeSum CalculationType = "sum"
)

type ShortestPathCalculation struct {
	graph           graph.Graph
	source          graph.Node
	destination     graph.Node
	weightType      helper.WeightKey
	calculationType CalculationType
	distances       map[interface{}]float64
	previous        map[interface{}]graph.Edge
	priorityQueue   PriorityQueue
}

func NewShortestPathCalculation(network graph.Graph, source, destination graph.Node, weightType helper.WeightKey, calculationType CalculationType) *ShortestPathCalculation {
	return &ShortestPathCalculation{
		graph:           network,
		source:          source,
		destination:     destination,
		weightType:      weightType,
		calculationType: calculationType,
		distances:       make(map[interface{}]float64),
		previous:        make(map[interface{}]graph.Edge),
		priorityQueue:   make(PriorityQueue, 0),
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
	for id := range calculation.graph.GetNodes() {
		calculation.distances[id] = math.Inf(1)
	}
	heap.Init(&calculation.priorityQueue)
	calculation.distances[calculation.source.GetId()] = 0
	heap.Push(&calculation.priorityQueue, &Item{nodeId: calculation.source.GetId(), distance: calculation.distances[calculation.source.GetId()], index: 0})
}

func (calculation *ShortestPathCalculation) performDijkstra() error {
	for !calculation.priorityQueue.IsEmpty() {
		item := heap.Pop(&calculation.priorityQueue).(*Item)
		currentId := item.GetNodeId()
		if currentId == calculation.destination.GetId() {
			break
		}
		currentNode, _ := calculation.graph.GetNode(currentId)
		for _, edge := range currentNode.GetEdges() {
			neighbor := edge.To()
			weight, err := edge.GetWeight(calculation.weightType)
			if err != nil {
				return err
			}
			var alternativeDistance float64
			if calculation.weightType == helper.PacketLossKey {
				packetLossTransform := -math.Log(1 - weight)
				alternativeDistance = calculation.distances[currentId] + packetLossTransform
			} else {
				alternativeDistance = calculation.distances[currentId] + weight
			}

			if alternativeDistance < calculation.distances[neighbor.GetId()] {
				neighborId := neighbor.GetId()
				calculation.distances[neighborId] = alternativeDistance
				calculation.previous[neighborId] = edge
				heap.Push(&calculation.priorityQueue, &Item{nodeId: neighborId, distance: alternativeDistance})
			}
		}
	}
	return nil
}

func (calculation ShortestPathCalculation) reconstructPath() (graph.Path, error) {
	path := make([]graph.Edge, 0)
	current := calculation.destination
	var totalCost float64
	if calculation.weightType != helper.PacketLossKey {
		totalCost = 0
	} else {
		totalCost = 1
	}
	for current.GetId() != calculation.source.GetId() {
		edge := calculation.previous[current.GetId()]
		path = append([]graph.Edge{edge}, path...)
		cost, err := edge.GetWeight(calculation.weightType)
		if err != nil {
			return nil, err
		} else if calculation.weightType != helper.PacketLossKey {
			totalCost += cost
		} else {
			totalCost *= cost
		}
		current = edge.From()
	}
	if len(path) == 0 {
		return nil, fmt.Errorf("No path found from node %d to node %d", calculation.source.GetId(), calculation.destination.GetId())
	}
	return graph.NewShortestPath(path, totalCost), nil
}
