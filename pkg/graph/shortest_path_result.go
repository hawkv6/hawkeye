package graph

type ShortestPathResult struct {
	edges []Edge
	cost  float64
}

func NewShortestPathResult(edges []Edge, cost float64) *ShortestPathResult {
	return &ShortestPathResult{
		edges: edges,
		cost:  cost,
	}
}

func (result *ShortestPathResult) GetEdges() []Edge {
	return result.edges
}

func (result *ShortestPathResult) GetCost() float64 {
	return result.cost
}

func (result *ShortestPathResult) SetCost(cost float64) {
	result.cost = cost
}
