package graph

import (
	"fmt"
	"sync"

	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type NetworkGraph struct {
	log    *logrus.Entry
	nodes  map[interface{}]Node
	edges  map[interface{}]Edge
	mu     *sync.Mutex
	helper helper.Helper
}

func NewNetworkGraph(helper helper.Helper) *NetworkGraph {
	return &NetworkGraph{
		log:    logging.DefaultLogger.WithField("subsystem", Subsystem),
		nodes:  make(map[interface{}]Node),
		edges:  make(map[interface{}]Edge),
		mu:     &sync.Mutex{},
		helper: helper,
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

func (graph *NetworkGraph) GetNodes() map[interface{}]Node {
	return graph.nodes
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
