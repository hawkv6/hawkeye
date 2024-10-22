package graph

import (
	"github.com/hawkv6/hawkeye/pkg/helper"
)

type NetworkEdge struct {
	id                 string
	from               Node
	to                 Node
	weights            map[helper.WeightKey]float64
	flexibleAlgorithms map[uint32]struct{}
}

func NewNetworkEdge(id string, from Node, to Node, weights map[helper.WeightKey]float64) *NetworkEdge {
	networkEdge := &NetworkEdge{
		id:                 id,
		from:               from,
		to:                 to,
		weights:            weights,
		flexibleAlgorithms: make(map[uint32]struct{}),
	}
	networkEdge.UpdateFlexibleAlgorithms()
	return networkEdge
}

func (edge *NetworkEdge) GetId() string {
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

func (edge *NetworkEdge) GetFlexibleAlgorithms() map[uint32]struct{} {
	return edge.flexibleAlgorithms
}

func (edge *NetworkEdge) UpdateFlexibleAlgorithms() {
	edge.flexibleAlgorithms = make(map[uint32]struct{})
	fromNodeAlgorithms := edge.From().GetFlexibleAlgorithms()
	toNodeAlgorithms := edge.To().GetFlexibleAlgorithms()
	for algorithm := range fromNodeAlgorithms {
		if _, exists := toNodeAlgorithms[algorithm]; exists {
			edge.flexibleAlgorithms[algorithm] = struct{}{}
		}
	}
}
