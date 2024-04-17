package graph

import (
	"container/heap"
	"fmt"
	"math"
	"sync"

	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type NetworkGraph struct {
	log   *logrus.Entry
	nodes map[interface{}]Node
	edges map[interface{}]Edge
	mu    *sync.Mutex
}

func NewNetworkGraph() *NetworkGraph {
	return &NetworkGraph{
		log:   logging.DefaultLogger.WithField("subsystem", Subsystem),
		nodes: make(map[interface{}]Node),
		edges: make(map[interface{}]Edge),
		mu:    &sync.Mutex{},
	}
}

func (graph *NetworkGraph) Lock() {
	graph.mu.Lock()
}

func (graph *NetworkGraph) Unlock() {
	graph.mu.Unlock()
}

func (graph *NetworkGraph) NodeExists(id interface{}) bool {
	_, exists := graph.nodes[id]
	return exists
}

func (graph *NetworkGraph) GetNode(id interface{}) (Node, bool) {
	node, exists := graph.nodes[id]
	return node, exists
}

func (graph *NetworkGraph) AddNode(node Node) (Node, error) {
	if graph.NodeExists(node.GetId()) {
		return nil, fmt.Errorf("Node with id %d already exists", node.GetId())
	}
	graph.nodes[node.GetId()] = node
	return node, nil
}

func (graph *NetworkGraph) DeleteNode(node Node) {
	for _, edge := range node.GetEdges() {
		graph.DeleteEdge(edge)
	}
	delete(graph.nodes, node.GetId())
}

func (graph *NetworkGraph) GetEdge(id interface{}) (Edge, bool) {
	edge, exists := graph.edges[id]
	return edge, exists
}

func (graph *NetworkGraph) EdgeExists(id interface{}) bool {
	_, exists := graph.edges[id]
	return exists
}

func (graph *NetworkGraph) AddEdge(edge Edge) error {
	fromId := edge.From().GetId()
	toId := edge.To().GetId()
	graph.log.Debugf("Add edge from %s to %s with weights %v", fromId, toId, edge.GetAllWeights())
	if !graph.NodeExists(fromId) {
		return fmt.Errorf("Node with id %d does not exist", fromId)
	}
	if !graph.NodeExists(toId) {
		return fmt.Errorf("Node with id %d does not exist", toId)
	}
	edgeId := edge.GetId()
	if graph.EdgeExists(edgeId) {
		return fmt.Errorf("Edge with id %d already exists", edgeId)
	}
	graph.edges[edgeId] = edge
	edge.From().AddEdge(edge)
	return nil
}

func (graph *NetworkGraph) DeleteEdge(edge Edge) {
	from := edge.From()
	to := edge.To()
	from.DeleteEdge(edge.GetId())
	to.DeleteEdge(edge.GetId())
	delete(graph.edges, edge.GetId())
}

func (graph *NetworkGraph) GetShortestPath(from Node, to Node, weightType string) (PathResult, error) {
	distances, priorityQueue := graph.initializeDijkstra(from)
	previous := make(map[interface{}]Edge)

	graph.performDijkstra(to, weightType, distances, &priorityQueue, previous)

	path, cost, err := graph.reconstructPath(from, to, previous, weightType)
	if err != nil {
		return nil, err
	}
	return NewShortestPathResult(path, cost), nil
}

func (graph *NetworkGraph) initializeDijkstra(from Node) (map[interface{}]float64, PriorityQueue) {
	distances := make(map[interface{}]float64)
	priorityQueue := make(PriorityQueue, 0)
	for id := range graph.nodes {
		distances[id] = math.Inf(1)
	}
	heap.Init(&priorityQueue)
	distances[from.GetId()] = 0
	heap.Push(&priorityQueue, &Item{nodeId: from.GetId(), distance: 0, index: 0})

	return distances, priorityQueue
}

func (graph *NetworkGraph) performDijkstra(to Node, weightKind string, distances map[interface{}]float64, priorityQueue *PriorityQueue, previous map[interface{}]Edge) {
	for !priorityQueue.IsEmpty() {
		item := heap.Pop(priorityQueue).(*Item)
		currentId := item.GetNodeId()
		if currentId == to.GetId() {
			break
		}
		for _, edge := range graph.nodes[currentId].GetEdges() {
			neighbor := edge.To()
			weight, err := edge.GetWeight(weightKind)
			if err != nil {
				return
			}
			alternativeDistance := distances[currentId] + weight
			if alternativeDistance < distances[neighbor.GetId()] {
				neighborId := neighbor.GetId()
				distances[neighborId] = alternativeDistance
				previous[neighborId] = edge
				heap.Push(priorityQueue, &Item{nodeId: neighborId, distance: alternativeDistance})
			}
		}
	}
}

func (graph *NetworkGraph) reconstructPath(from Node, to Node, previous map[interface{}]Edge, weightType string) ([]Edge, float64, error) {
	path := make([]Edge, 0)
	current := to
	totalCost := 0.0
	for current.GetId() != from.GetId() {
		edge := previous[current.GetId()]
		path = append([]Edge{edge}, path...)
		if cost, err := edge.GetWeight(weightType); err != nil {
			return nil, 0, err
		} else {
			totalCost += cost
		}
		current = edge.From()
	}
	if len(path) == 0 {
		return nil, 0, fmt.Errorf("No path found from node %d to node %d", from.GetId(), to.GetId())
	}
	return path, totalCost, nil
}
