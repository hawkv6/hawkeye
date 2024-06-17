package graph

import (
	"fmt"
	"sync"

	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type NetworkGraph struct {
	log      *logrus.Entry
	nodes    map[string]Node
	edges    map[string]Edge
	mu       *sync.Mutex
	helper   helper.Helper
	isLocked bool
}

func NewNetworkGraph(helper helper.Helper) *NetworkGraph {
	return &NetworkGraph{
		log:      logging.DefaultLogger.WithField("subsystem", Subsystem),
		nodes:    make(map[string]Node),
		edges:    make(map[string]Edge),
		mu:       &sync.Mutex{},
		helper:   helper,
		isLocked: false,
	}
}

func (graph *NetworkGraph) Lock() {
	graph.mu.Lock()
}

func (graph *NetworkGraph) Unlock() {
	graph.mu.Unlock()
}

func (graph *NetworkGraph) NodeExists(id string) bool {
	_, exists := graph.nodes[id]
	return exists
}

func (graph *NetworkGraph) GetNode(id string) (Node, bool) {
	node, exists := graph.nodes[id]
	return node, exists
}

func (graph *NetworkGraph) GetNodes() map[string]Node {
	return graph.nodes
}

func (graph *NetworkGraph) AddNode(node Node) (Node, error) {
	if graph.NodeExists(node.GetId()) {
		return nil, fmt.Errorf("Node with id %s already exists", node.GetId())
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

func (graph *NetworkGraph) GetEdge(id string) (Edge, bool) {
	edge, exists := graph.edges[id]
	return edge, exists
}

func (graph *NetworkGraph) GetEdges() map[string]Edge {
	return graph.edges
}

func (graph *NetworkGraph) EdgeExists(id string) bool {
	_, exists := graph.edges[id]
	return exists
}

func (graph *NetworkGraph) AddEdge(edge Edge) error {
	fromId := edge.From().GetId()
	toId := edge.To().GetId()
	graph.log.Debugf("Add edge from %s to %s with weights %v", fromId, toId, edge.GetAllWeights())
	if !graph.NodeExists(fromId) {
		return fmt.Errorf("Node with id %s does not exist", fromId)
	}
	if !graph.NodeExists(toId) {
		return fmt.Errorf("Node with id %s does not exist", toId)
	}
	edgeId := edge.GetId()
	if graph.EdgeExists(edgeId) {
		return fmt.Errorf("Edge with id %s already exists", edgeId)
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
