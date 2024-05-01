package graph

const Subsystem = "graph"

type Graph interface {
	Lock()
	Unlock()
	AddNode(node Node) (Node, error)
	GetNode(id interface{}) (Node, bool)
	GetNodes() map[interface{}]Node
	DeleteNode(node Node)
	NodeExists(id interface{}) bool
	AddEdge(edge Edge) error
	GetEdge(id interface{}) (Edge, bool)
	EdgeExists(id interface{}) bool
	DeleteEdge(edge Edge)
}
