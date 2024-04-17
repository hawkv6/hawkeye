package graph

import "fmt"

type NetworkEdge struct {
	id      interface{}
	from    Node
	to      Node
	weights map[string]float64
}

func NewNetworkEdge(id interface{}, from Node, to Node, weights map[string]float64) *NetworkEdge {
	return &NetworkEdge{
		id:      id,
		from:    from,
		to:      to,
		weights: weights,
	}
}

func (edge *NetworkEdge) GetId() interface{} {
	return edge.id
}

func (edge *NetworkEdge) From() Node {
	return edge.from
}

func (edge *NetworkEdge) To() Node {
	return edge.to
}

func (edge *NetworkEdge) GetAllWeights() map[string]float64 {
	return edge.weights
}

func (edge *NetworkEdge) GetWeight(kind string) (float64, error) {
	weight, ok := edge.weights[kind]
	if !ok {
		return 0, fmt.Errorf("weight kind %s not found", kind)
	}
	return weight, nil
}

func (edge *NetworkEdge) SetWeight(kind string, weight float64) {
	edge.weights[kind] = weight
}
