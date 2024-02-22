package graph

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type DefaultGraph struct {
	log   *logrus.Entry
	nodes map[interface{}]Node
	edges map[interface{}]Edge
}

func NewDefaultGraph() *DefaultGraph {
	return &DefaultGraph{
		log:   logging.DefaultLogger.WithField("subsystem", Subsystem),
		nodes: make(map[interface{}]Node),
		edges: make(map[interface{}]Edge),
	}
}

func (graph *DefaultGraph) NodeExists(id interface{}) bool {
	_, exists := graph.nodes[id]
	return exists
}

func (graph *DefaultGraph) GetNode(id interface{}) (Node, error) {
	if !graph.NodeExists(id) {
		return nil, fmt.Errorf("node with id %d does not exist", id)
	}
	return graph.nodes[id], nil
}

func (graph *DefaultGraph) AddNode(node Node) error {
	if graph.NodeExists(node.GetId()) {
		return fmt.Errorf("node with id %d already exists", node.GetId())
	}
	graph.nodes[node.GetId()] = node
	return nil
}

func (graph *DefaultGraph) EdgeExists(id interface{}) bool {
	_, exists := graph.edges[id]
	return exists
}

func (graph *DefaultGraph) AddEdge(edge Edge) error {
	from := edge.From()
	to := edge.To()
	fromId := from.GetId()
	toId := to.GetId()
	weight, err := edge.GetWeight("delay")
	if err != nil {
		return err
	}
	graph.log.Debugf(fmt.Sprintf("Adding edge from %s to %s with weight %f", fromId, toId, weight))

	if !graph.NodeExists(edge.From().GetId()) {
		return fmt.Errorf("node with id %d does not exist", edge.From().GetId())
	}
	if !graph.NodeExists(edge.To().GetId()) {
		return fmt.Errorf("node with id %d does not exist", edge.To().GetId())
	}
	if graph.EdgeExists(edge.GetId()) {
		return fmt.Errorf("edge with id %d already exists", edge.GetId())
	}
	graph.edges[edge.GetId()] = edge
	edge.From().AddEdge(edge)
	return nil
}

func (graph *DefaultGraph) RemoveNode(node Node) {
	delete(graph.nodes, node.GetId())
}

func (graph *DefaultGraph) RemoveEdge(edge Edge) {
	delete(graph.edges, edge.GetId())
}

func (graph *DefaultGraph) GetShortestPath(from Node, to Node, weightKind string) ([]Edge, error) {
	distances, priorityQueue := graph.initializeDijkstra(from)
	previous := make(map[interface{}]Edge)

	graph.performDijkstra(from, to, weightKind, distances, &priorityQueue, previous)

	path, err := graph.reconstructPath(from, to, previous)
	if err != nil {
		return nil, err
	}

	return path, nil
}

func (graph *DefaultGraph) initializeDijkstra(from Node) (map[interface{}]float64, PriorityQueue) {
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

func (graph *DefaultGraph) performDijkstra(from Node, to Node, weightKind string, distances map[interface{}]float64, priorityQueue *PriorityQueue, previous map[interface{}]Edge) {
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

func (graph *DefaultGraph) reconstructPath(from Node, to Node, previous map[interface{}]Edge) ([]Edge, error) {
	path := make([]Edge, 0)
	current := to
	for current.GetId() != from.GetId() {
		edge := previous[current.GetId()]
		path = append([]Edge{edge}, path...)
		current = edge.From()
	}

	if len(path) == 0 {
		return nil, fmt.Errorf("no path found from node %d to node %d", from.GetId(), to.GetId())
	}

	return path, nil
}
