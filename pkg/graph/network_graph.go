package graph

import (
	"fmt"
	"sync"

	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type NetworkGraph struct {
	log       *logrus.Entry
	nodes     map[string]Node
	edges     map[string]Edge
	mu        *sync.Mutex
	isLocked  bool
	subGraphs map[uint32]*NetworkGraph
}

func NewNetworkGraph() *NetworkGraph {
	return &NetworkGraph{
		log:      logging.DefaultLogger.WithField("subsystem", Subsystem),
		nodes:    make(map[string]Node),
		edges:    make(map[string]Edge),
		mu:       &sync.Mutex{},
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

func (graph *NetworkGraph) GetNode(id string) Node {
	return graph.nodes[id]
}

func (graph *NetworkGraph) GetNodes() map[string]Node {
	return graph.nodes
}

func (graph *NetworkGraph) AddNode(node Node) Node {
	graph.nodes[node.GetId()] = node
	return node
}

func (graph *NetworkGraph) DeleteNode(node Node) {
	for _, edge := range node.GetEdges() {
		graph.DeleteEdge(edge)
	}
	delete(graph.nodes, node.GetId())
}

func (graph *NetworkGraph) GetEdge(id string) Edge {
	return graph.edges[id]
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

func (graph *NetworkGraph) addNodesToSubgraph(newSubGraphs map[uint32]*NetworkGraph) {
	for _, node := range graph.nodes {
		for flexAlgo := range node.GetFlexibleAlgorithms() {
			if _, exists := newSubGraphs[flexAlgo]; !exists {
				newSubGraphs[flexAlgo] = NewNetworkGraph()
			}
			newSubGraphs[flexAlgo].AddNode(node)
		}
	}
}

func (graph *NetworkGraph) addEdgesToSubgraph(newSubGraphs map[uint32]*NetworkGraph) {
	for _, edge := range graph.edges {
		for flexAlgo := range edge.GetFlexibleAlgorithms() {
			if _, exists := newSubGraphs[flexAlgo]; !exists {
				graph.log.Errorf("Subgraph for flex algo %d does not exist, but was found in edge %s", flexAlgo, edge.GetId())
			}
			if err := newSubGraphs[flexAlgo].AddEdge(edge); err != nil {
				graph.log.Errorf("Failed to add edge %s to subgraph %d: %v", edge.GetId(), flexAlgo, err)
			}
		}
	}
}

func (graph *NetworkGraph) UpdateSubGraphs() {
	graph.mu.Lock()
	defer graph.mu.Unlock()
	graph.log.Debugln("Updating subgraphs")
	newSubGraphs := make(map[uint32]*NetworkGraph)
	graph.addNodesToSubgraph(newSubGraphs)
	graph.addEdgesToSubgraph(newSubGraphs)
	graph.subGraphs = newSubGraphs
}

func (graph *NetworkGraph) GetSubGraph(algorithm uint32) *NetworkGraph {
	return graph.subGraphs[algorithm]
}
