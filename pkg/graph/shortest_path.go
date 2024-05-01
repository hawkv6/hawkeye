package graph

type ShortestPath struct {
	edges []Edge
	cost  float64
}

func NewShortestPath(edges []Edge, cost float64) *ShortestPath {
	return &ShortestPath{
		edges: edges,
		cost:  cost,
	}
}

func (path *ShortestPath) GetEdges() []Edge {
	return path.edges
}

func (path *ShortestPath) GetCost() float64 {
	return path.cost
}

func (path *ShortestPath) SetCost(cost float64) {
	path.cost = cost
}
