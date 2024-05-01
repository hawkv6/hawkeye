package graph

type Path interface {
	GetEdges() []Edge
	GetCost() float64
	SetCost(float64)
}
