package graph

const Subsystem = "graph"

type Graph interface {
	Lock()
	Unlock()
	AddNode(node Node) (Node, error)
	GetNode(id string) (Node, bool)
	GetNodes() map[string]Node
	DeleteNode(node Node)
	NodeExists(id string) bool
	AddEdge(edge Edge) error
	GetEdge(id string) (Edge, bool)
	GetEdges() map[string]Edge
	EdgeExists(id string) bool
	DeleteEdge(edge Edge)
}
