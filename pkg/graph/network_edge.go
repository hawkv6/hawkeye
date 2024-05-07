package graph

import (
	"github.com/hawkv6/hawkeye/pkg/helper"
)

type NetworkEdge struct {
	id      interface{}
	from    Node
	to      Node
	weights map[helper.WeightKey]float64
}

func NewNetworkEdge(id interface{}, from Node, to Node, weights map[helper.WeightKey]float64) *NetworkEdge {
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

func (edge *NetworkEdge) GetAllWeights() map[helper.WeightKey]float64 {
	return edge.weights
}

func (edge *NetworkEdge) GetWeight(kind helper.WeightKey) float64 {
	return edge.weights[kind]
}

func (edge *NetworkEdge) SetWeight(kind helper.WeightKey, weight float64) {
	edge.weights[kind] = weight
}
