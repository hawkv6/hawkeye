package graph

type Path interface {
	GetEdges() []Edge
	GetTotalCost() float64
	SetTotalCost(float64)
	GetBottleneckEdge() Edge
	GetBottleneckValue() float64
	SetBottleneckEdge(Edge)
	SetBottleneckValue(float64)
}
