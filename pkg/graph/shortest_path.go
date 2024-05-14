package graph

type ShortestPath struct {
	edges           []Edge
	totalCost       float64
	bottleneckEdge  Edge
	bottleneckValue float64
}

func NewShortestPathWithTotalCost(edges []Edge, cost float64) *ShortestPath {
	return &ShortestPath{
		edges:     edges,
		totalCost: cost,
	}
}
func NewShortestPathWithBottleneck(edges []Edge, bottleneckEdge Edge, bottleneckValue float64) *ShortestPath {
	return &ShortestPath{
		edges:           edges,
		bottleneckEdge:  bottleneckEdge,
		bottleneckValue: bottleneckValue,
	}
}

func (path *ShortestPath) GetEdges() []Edge {
	return path.edges
}

func (path *ShortestPath) GetTotalCost() float64 {
	return path.totalCost
}

func (path *ShortestPath) SetTotalCost(cost float64) {
	path.totalCost = cost
}

func (path *ShortestPath) GetBottleneckEdge() Edge {
	return path.bottleneckEdge
}

func (path *ShortestPath) GetBottleneckValue() float64 {
	return path.bottleneckValue
}

func (path *ShortestPath) SetBottleneckEdge(edge Edge) {
	path.bottleneckEdge = edge
}

func (path *ShortestPath) SetBottleneckValue(value float64) {
	path.bottleneckValue = value
}
