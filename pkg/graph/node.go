package graph

type Node interface {
	GetId() string
	GetEdges() map[string]Edge
	AddEdge(Edge)
	DeleteEdge(string)
	GetName() string
	SetName(string)
	GetFlexibleAlgorithms() map[uint32]struct{}
	SetFlexibleAlgorithms([]uint32)
}
