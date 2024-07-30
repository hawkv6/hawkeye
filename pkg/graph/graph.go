package graph

const Subsystem = "graph"

type Graph interface {
	Lock()
	Unlock()
	AddNode(Node) Node
	GetNode(string) Node
	GetNodes() map[string]Node
	DeleteNode(Node)
	NodeExists(string) bool
	AddEdge(Edge) error
	GetEdge(string) Edge
	GetEdges() map[string]Edge
	EdgeExists(string) bool
	DeleteEdge(Edge)
	UpdateSubGraphs()
	GetSubGraph(uint32) Graph
}
