package graph

type PathResult interface {
	GetEdges() []Edge
	GetCost() float64
}
