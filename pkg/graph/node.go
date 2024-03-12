package graph

type Node interface {
	GetId() interface{}
	GetEdges() map[interface{}]Edge
	AddEdge(Edge)
	DeleteEdge(interface{})
	GetName() string
	SetName(string)
}
