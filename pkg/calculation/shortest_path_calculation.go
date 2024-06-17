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
	log              *logrus.Entry
	graph            graph.Graph
	source           graph.Node
	destination      graph.Node
	weightTypes      []helper.WeightKey
	calculationType  CalculationMode
	nodeWeights      map[string]float64
	EdgeToPrevious   map[string]graph.Edge
	priorityQueue    PriorityQueue
	visitedNodes     map[string]bool
	maxConstraints   map[helper.WeightKey]float64
	nodeLatencies    map[string]float64
	nodeJitters      map[string]float64
	nodePacketLosses map[string]float64
}

func NewShortestPathCalculation(network graph.Graph, source, destination graph.Node, weightTypes []helper.WeightKey, calculationType CalculationMode, maxConstraints map[helper.WeightKey]float64) *ShortestPathCalculation {
	return &ShortestPathCalculation{
		log:              logging.DefaultLogger.WithField("subsystem", "calculation"),
		graph:            network,
		source:           source,
		destination:      destination,
		weightTypes:      weightTypes,
		calculationType:  calculationType,
		nodeWeights:      make(map[string]float64),
		EdgeToPrevious:   make(map[string]graph.Edge),
		visitedNodes:     make(map[string]bool),
		maxConstraints:   maxConstraints,
		nodeLatencies:    make(map[string]float64),
		nodeJitters:      make(map[string]float64),
		nodePacketLosses: make(map[string]float64),
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
		calculation.nodeWeights[id] = initialNodeCost
		calculation.nodeLatencies[id] = 0
		calculation.nodeJitters[id] = 0
		calculation.nodePacketLosses[id] = 0
	}
	heap.Init(&calculation.priorityQueue)
	calculation.nodeWeights[calculation.source.GetId()] = sourceNodeCost
	heap.Push(&calculation.priorityQueue, &Item{nodeId: calculation.source.GetId(), cost: calculation.nodeWeights[calculation.source.GetId()], index: 0})
}

func (calculation *ShortestPathCalculation) getMetrics(edge graph.Edge) (float64, float64, float64) {
	latency := edge.GetWeight(helper.LatencyKey)
	jitter := edge.GetWeight(helper.JitterKey)
	packetLoss := edge.GetWeight(helper.PacketLossKey)
	return latency, jitter, packetLoss
}

func (calculation *ShortestPathCalculation) updateMetricsAndPrevious(neighborId string, weight float64, edge graph.Edge) {
	latency, jitter, packetLoss := calculation.getMetrics(edge)
	latency += calculation.nodeLatencies[edge.From().GetId()]
	jitter += calculation.nodeJitters[edge.From().GetId()]
	packetLoss = 1 - ((1 - calculation.nodePacketLosses[edge.From().GetId()]) * (1 - packetLoss/100))

	metrics := map[helper.WeightKey]float64{
		helper.LatencyKey:              latency,
		helper.JitterKey:               jitter,
		helper.NormalizedPacketLossKey: packetLoss * 100,
	}

	for key, value := range metrics {
		if maxValue, ok := calculation.maxConstraints[key]; ok {
			if maxValue < value {
				calculation.log.Debugf("Path violates %s constraint", key)
				return
			}
		}
	}

	calculation.nodeWeights[neighborId] = weight
	calculation.EdgeToPrevious[neighborId] = edge
	calculation.nodeLatencies[neighborId] = latency
	calculation.nodeJitters[neighborId] = jitter
	calculation.nodePacketLosses[neighborId] = packetLoss
	heap.Push(&calculation.priorityQueue, &Item{nodeId: neighborId, cost: weight})
}

func (calculation *ShortestPathCalculation) calculateAlternativeDistance(currentId string, weight float64) float64 {
	if calculation.weightTypes[0] == helper.PacketLossKey {
		weightPercentage := weight / 100
		return calculation.nodeWeights[currentId] + -math.Log(1-weightPercentage)
	}
	return calculation.nodeWeights[currentId] + weight
}

func (calculation *ShortestPathCalculation) handleMaxCalculation(currentId, neighborId string, weight float64, edge graph.Edge) {
	minimum := math.Min(calculation.nodeWeights[currentId], weight)
	if minimum > calculation.nodeWeights[neighborId] {
		calculation.updateMetricsAndPrevious(neighborId, minimum, edge)
	}
}

func (calculation *ShortestPathCalculation) handleMinCalculation(currentId, neighborId string, weight float64, edge graph.Edge) {
	minimum := math.Min(calculation.nodeWeights[currentId], weight)
	if minimum < calculation.nodeWeights[neighborId] {
		calculation.updateMetricsAndPrevious(neighborId, minimum, edge)
	}
}

func (calculation *ShortestPathCalculation) handleDefaultCalculation(currentId, neighborId string, weight float64, edge graph.Edge) {
	alternativeDistance := calculation.calculateAlternativeDistance(currentId, weight)
	if alternativeDistance < calculation.nodeWeights[neighborId] {
		calculation.updateMetricsAndPrevious(neighborId, alternativeDistance, edge)
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

func (calculation *ShortestPathCalculation) relaxEdge(currentId string, edge graph.Edge) {
	neighbor := edge.To()
	neighborId := neighbor.GetId()
	if _, ok := calculation.visitedNodes[neighborId]; ok {
		return
	}
	weight := calculation.getWeight(edge)

	if calculation.calculationType == CalculationModeMax {
		calculation.handleMaxCalculation(currentId, neighborId, weight, edge)
	} else if calculation.calculationType == CalculationModeMin {
		calculation.handleMinCalculation(currentId, neighborId, weight, edge)
	} else {
		calculation.handleDefaultCalculation(currentId, neighborId, weight, edge)
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

func (calculation *ShortestPathCalculation) getPath(current graph.Node) ([]graph.Edge, error) {
	path := make([]graph.Edge, 0)
	for current.GetId() != calculation.source.GetId() {
		edge := calculation.EdgeToPrevious[current.GetId()]
		if edge == nil {
			return nil, fmt.Errorf("No path found from node %s to node %s", calculation.source.GetId(), calculation.destination.GetId())
		}
		path = append([]graph.Edge{edge}, path...)
		current = edge.From()
	}
	if len(path) == 0 {
		return nil, fmt.Errorf("No path found from node %s to node %s", calculation.source.GetId(), calculation.destination.GetId())
	}
	return path, nil
}

func (calculation *ShortestPathCalculation) reconstructSumPath(current graph.Node) (graph.Path, error) {
	path, err := calculation.getPath(current)
	if err != nil {
		return nil, err
	}
	totalCost := calculation.nodeWeights[calculation.destination.GetId()]
	latency := calculation.nodeLatencies[calculation.destination.GetId()]
	jitter := calculation.nodeJitters[calculation.destination.GetId()]
	packetLoss := calculation.nodePacketLosses[calculation.destination.GetId()]
	return graph.NewShortestPathWithTotalCost(path, totalCost, latency, jitter, packetLoss), nil
}

func (calculation *ShortestPathCalculation) reconstructMinPath(current graph.Node) (graph.Path, error) {
	path := make([]graph.Edge, 0) // TODO find solution for that
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
		return nil, fmt.Errorf("No path found from node %s to node %s", calculation.source.GetId(), calculation.destination.GetId())
	}
	calculation.log.Debugf("First bottleneck edge: %v with bandwidth %g: ", minEdge, minEdgeBandwidth)
	return graph.NewShortestPathWithBottleneck(path, minEdge, minEdgeBandwidth), nil
}

func (calculation ShortestPathCalculation) reconstructPath() (graph.Path, error) {
	if calculation.calculationType == CalculationModeSum {
		return calculation.reconstructSumPath(calculation.destination)
	}
	return calculation.reconstructMinPath(calculation.destination)
}
