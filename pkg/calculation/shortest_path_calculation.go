package calculation

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
)

type ShortestPathCalculation struct {
	BaseCalculation
	nodeWeights      map[string]float64
	EdgeToPrevious   map[string]graph.Edge
	priorityQueue    PriorityQueue
	visitedNodes     map[string]bool
	nodeLatencies    map[string]float64
	nodeJitters      map[string]float64
	nodePacketLosses map[string]float64
	sourceNodeCost   float64
}

func NewShortestPathCalculation(options *CalculationOptions) *ShortestPathCalculation {
	return &ShortestPathCalculation{
		BaseCalculation:  *NewBaseCalculation(options),
		nodeWeights:      make(map[string]float64),
		EdgeToPrevious:   make(map[string]graph.Edge),
		visitedNodes:     make(map[string]bool),
		nodeLatencies:    make(map[string]float64),
		nodeJitters:      make(map[string]float64),
		nodePacketLosses: make(map[string]float64),
		sourceNodeCost:   0,
	}
}

func (calculation *ShortestPathCalculation) Execute() (graph.Path, error) {
	calculation.initializeDijkstra()
	calculation.performDijkstra()
	return calculation.reconstructPath()
}

func (calculation *ShortestPathCalculation) SetInitialSourceNodeMetrics(cost float64, latency, jitter, packetLoss float64) {
	source := calculation.source.GetId()
	calculation.sourceNodeCost = cost
	calculation.nodeLatencies[source] = latency
	calculation.nodeJitters[source] = jitter
	calculation.nodePacketLosses[source] = packetLoss
}

func (calculation *ShortestPathCalculation) initializeNodeMetrics(initialNodeCost float64, sourceNodeId string) {
	for id := range calculation.graph.GetNodes() {
		if id != sourceNodeId {
			calculation.nodeWeights[id] = initialNodeCost
			calculation.nodeLatencies[id] = 0
			calculation.nodeJitters[id] = 0
			calculation.nodePacketLosses[id] = 0
		}
	}
}

func (calculation *ShortestPathCalculation) initializeHeap(sourceNodeCost float64, sourceNodeId string) {
	heap.Init(&calculation.priorityQueue)
	calculation.nodeWeights[sourceNodeId] = sourceNodeCost
	newItem := &Item{
		nodeId: sourceNodeId,
		cost:   calculation.nodeWeights[sourceNodeId],
		index:  0,
	}
	heap.Push(&calculation.priorityQueue, newItem)
}

func (calculation *ShortestPathCalculation) initializeDijkstra() {
	var initialNodeCost, sourceNodeCost float64
	calculation.priorityQueue = *NewMinimumPriorityQueue()
	switch calculation.calculationMode {
	case CalculationModeMax:
		initialNodeCost = 0
		sourceNodeCost = math.Inf(1)
		calculation.priorityQueue = *NewMaximumPriorityQueue()
	case CalculationModeMin:
		initialNodeCost = math.Inf(1)
		sourceNodeCost = math.Inf(1)
	default:
		initialNodeCost = math.Inf(1)
		sourceNodeCost = calculation.sourceNodeCost
	}
	sourceNodeId := calculation.source.GetId()
	calculation.initializeNodeMetrics(initialNodeCost, sourceNodeId)
	calculation.initializeHeap(sourceNodeCost, sourceNodeId)
}

func (calculation *ShortestPathCalculation) getMetrics(edge graph.Edge, currentNodeId string) (float64, float64, float64) {
	latency := calculation.nodeLatencies[currentNodeId] + edge.GetWeight(helper.LatencyKey)
	jitter := calculation.nodeJitters[currentNodeId] + edge.GetWeight(helper.JitterKey)
	packetLoss := edge.GetWeight(helper.PacketLossKey)
	packetLoss = 1 - ((1 - calculation.nodePacketLosses[currentNodeId]) * (1 - packetLoss/100))
	return latency, jitter, packetLoss
}

func (calculation *ShortestPathCalculation) violatesMaxConstraints(edge graph.Edge, latency, jitter, packetLoss float64) bool {
	metrics := map[helper.WeightKey]float64{
		helper.NormalizedLatencyKey:    latency,
		helper.NormalizedJitterKey:     jitter,
		helper.NormalizedPacketLossKey: packetLoss,
	}
	for key, value := range metrics {
		if maxValue, ok := calculation.maxConstraints[key]; ok {
			if maxValue < value {
				calculation.log.Debugf("Edge from %s to %s violates %s constraint, returning", edge.From().GetName(), edge.To().GetName(), key)
				return true
			}
		}
	}
	return false
}

func (calculation *ShortestPathCalculation) violatesBandwidthMinConstraint(edge graph.Edge) bool {
	if minValue, ok := calculation.minConstraints[helper.AvailableBandwidthKey]; ok {
		bandwidth := edge.GetWeight(helper.AvailableBandwidthKey)
		if minValue > bandwidth {
			calculation.log.Debugf("Edge %v violates %s constraint, returning", edge, helper.AvailableBandwidthKey)
			return true
		}
	}
	return false
}

func (calculation *ShortestPathCalculation) updateMetricsAndPrevious(currentNodeId, neighborNodeId string, weight float64, edge graph.Edge) {
	latency, jitter, packetLoss := calculation.getMetrics(edge, currentNodeId)
	if !calculation.violatesMaxConstraints(edge, latency, jitter, packetLoss) && !calculation.violatesBandwidthMinConstraint(edge) {
		calculation.nodeWeights[neighborNodeId] = weight
		calculation.EdgeToPrevious[neighborNodeId] = edge
		calculation.nodeLatencies[neighborNodeId] = latency
		calculation.nodeJitters[neighborNodeId] = jitter
		calculation.nodePacketLosses[neighborNodeId] = packetLoss
		heap.Push(&calculation.priorityQueue, &Item{nodeId: neighborNodeId, cost: weight})
	}
}

func (calculation *ShortestPathCalculation) calculateAlternativeDistance(currentNodeId string, edgeWeight float64) float64 {
	if calculation.weightKeys[0] == helper.PacketLossKey {
		packetLossPercentage := edgeWeight / 100
		return calculation.nodeWeights[currentNodeId] + -math.Log(1-packetLossPercentage)
	}
	return calculation.nodeWeights[currentNodeId] + edgeWeight
}

func (calculation *ShortestPathCalculation) handleDefaultCalculation(currentNodeId, neighborNodeId string, edgeWeight float64, edge graph.Edge) {
	alternativeDistance := calculation.calculateAlternativeDistance(currentNodeId, edgeWeight)
	if alternativeDistance < calculation.nodeWeights[neighborNodeId] {
		calculation.updateMetricsAndPrevious(currentNodeId, neighborNodeId, alternativeDistance, edge)
	}
}

func (calculation *ShortestPathCalculation) handleMaxCalculation(currentNodeId, neighborNodeId string, weight float64, edge graph.Edge) {
	minimum := math.Min(calculation.nodeWeights[currentNodeId], weight)
	if minimum > calculation.nodeWeights[neighborNodeId] {
		calculation.updateMetricsAndPrevious(currentNodeId, neighborNodeId, minimum, edge)
	}
}

func (calculation *ShortestPathCalculation) handleMinCalculation(currentNodeId, neighborNodeId string, weight float64, edge graph.Edge) {
	minimum := math.Min(calculation.nodeWeights[currentNodeId], weight)
	if minimum < calculation.nodeWeights[neighborNodeId] {
		calculation.updateMetricsAndPrevious(currentNodeId, neighborNodeId, minimum, edge)
	}
}

func (calculation *ShortestPathCalculation) handleCalculation(currentNodeId, neighborNodeId string, weight float64, edge graph.Edge) {
	if calculation.calculationMode == CalculationModeSum {
		calculation.handleDefaultCalculation(currentNodeId, neighborNodeId, weight, edge)
	} else if calculation.calculationMode == CalculationModeMax {
		calculation.handleMaxCalculation(currentNodeId, neighborNodeId, weight, edge)
	} else {
		calculation.handleMinCalculation(currentNodeId, neighborNodeId, weight, edge)
	}
}

func (calculation *ShortestPathCalculation) getEdgeWeight(edge graph.Edge) float64 {
	weight := 0.0
	if len(calculation.weightKeys) == 2 {
		weight = edge.GetWeight(calculation.weightKeys[0])*float64(helper.TwoFactorWeights[0]) + edge.GetWeight(calculation.weightKeys[1])*float64(helper.TwoFactorWeights[1])

	} else if len(calculation.weightKeys) == 3 {
		weight = edge.GetWeight(calculation.weightKeys[0])*float64(helper.ThreeFactorWeights[0]) + edge.GetWeight(calculation.weightKeys[1])*float64(helper.ThreeFactorWeights[1]) + edge.GetWeight(calculation.weightKeys[2])*float64(helper.ThreeFactorWeights[2])
	} else {
		weight = edge.GetWeight(calculation.weightKeys[0])
	}
	return weight
}

func (calculation *ShortestPathCalculation) relaxEdge(currentNodeId string, edge graph.Edge) {
	neighborNode := edge.To()
	neighborNodeId := neighborNode.GetId()
	if _, ok := calculation.visitedNodes[neighborNodeId]; ok {
		return
	}
	edgeWeight := calculation.getEdgeWeight(edge)
	calculation.handleCalculation(currentNodeId, neighborNodeId, edgeWeight, edge)
}

func (calculation *ShortestPathCalculation) performDijkstra() {
	for !calculation.priorityQueue.IsEmpty() {
		nodeItem := heap.Pop(&calculation.priorityQueue).(*Item)
		calculation.visitedNodes[nodeItem.GetNodeId()] = true
		currentNodeId := nodeItem.GetNodeId()
		if currentNodeId == calculation.destination.GetId() {
			break
		}
		currentNode := calculation.graph.GetNode(currentNodeId)
		for _, edge := range currentNode.GetEdges() {
			calculation.relaxEdge(currentNodeId, edge)
		}
	}
}

func (calculation *ShortestPathCalculation) getPath(current graph.Node) ([]graph.Edge, graph.Edge, error) {
	bottleneckBandwidth := math.Inf(1)
	var bottleneckEdge graph.Edge
	path := make([]graph.Edge, 0)
	for current.GetId() != calculation.source.GetId() {
		edge := calculation.EdgeToPrevious[current.GetId()]
		if edge == nil {
			return nil, nil, fmt.Errorf("No path found from node %s to node %s", calculation.source.GetId(), calculation.destination.GetId())
		}
		if edge.GetWeight(helper.AvailableBandwidthKey) < bottleneckBandwidth {
			calculation.log.Debugln("Bottleneck edge changed")
			bottleneckBandwidth = edge.GetWeight(helper.AvailableBandwidthKey)
			bottleneckEdge = edge
		}
		path = append([]graph.Edge{edge}, path...)
		current = edge.From()
	}
	if len(path) == 0 {
		return nil, nil, fmt.Errorf("No path found from node %s to node %s", calculation.source.GetId(), calculation.destination.GetId())
	}
	return path, bottleneckEdge, nil
}

func (calculation ShortestPathCalculation) reconstructPath() (graph.Path, error) {
	path, bottleneckEdge, err := calculation.getPath(calculation.destination)
	if err != nil {
		return nil, err
	}
	totalCost := calculation.nodeWeights[calculation.destination.GetId()]
	latency := calculation.nodeLatencies[calculation.destination.GetId()]
	jitter := calculation.nodeJitters[calculation.destination.GetId()]
	packetLoss := calculation.nodePacketLosses[calculation.destination.GetId()]
	bottleneckBandwidth := bottleneckEdge.GetWeight(helper.AvailableBandwidthKey)
	calculation.log.Debugln("Calculation finished - shortest path found")
	calculation.log.Debugf("Total cost %g, total latency %gus, total jitter %gus, total packet loss %f -> %f%%", totalCost, latency, jitter, packetLoss, packetLoss*100)
	calculation.log.Debugf("Available bandwidth %g, bottleneck edge form %s to %s", bottleneckBandwidth, bottleneckEdge.From().GetName(), bottleneckEdge.To().GetName())
	return graph.NewShortestPath(path, totalCost, latency, jitter, packetLoss, bottleneckBandwidth, bottleneckEdge), nil
}
