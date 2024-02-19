package graph

const Subsystem = "graph"

type Graph interface {
	GetShortestPath(from Node, to Node, weightKind string) ([]Edge, error)
	AddNode(node Node) error
	GetNode(id interface{}) (Node, error)
	NodeExists(id interface{}) bool
	AddEdge(edge Edge) error
	EdgeExists(id interface{}) bool
	RemoveNode(node Node)
	RemoveEdge(edge Edge)
}
