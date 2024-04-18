package graph

const Subsystem = "graph"

type Graph interface {
	Lock()
	Unlock()
	AddNode(node Node) (Node, error)
	GetNode(id interface{}) (Node, bool)
	DeleteNode(node Node)
	NodeExists(id interface{}) bool
	AddEdge(edge Edge) error
	GetEdge(id interface{}) (Edge, bool)
	EdgeExists(id interface{}) bool
	DeleteEdge(edge Edge)
	GetShortestPath(from Node, to Node, weightType string) (PathResult, error)
}
