package graph

type Path interface {
	GetEdges() []Edge
	GetTotalCost() float64
	SetTotalCost(float64)
	GetTotalDelay() float64
	GetTotalJitter() float64
	GetTotalPacketLoss() float64
	SetTotalDelay(float64)
	SetTotalJitter(float64)
	SetTotalPacketLoss(float64)
	GetBottleneckEdge() Edge
	GetBottleneckValue() float64
	SetBottleneckEdge(Edge)
	SetBottleneckValue(float64)
	SetRouterServiceMap(map[string]string)
	GetRouterServiceMap() map[string]string
}
