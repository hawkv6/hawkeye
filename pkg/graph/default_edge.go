package graph

import "fmt"

type DefaultEdge struct {
	id      interface{}
	from    Node
	to      Node
	weights map[string]float64
}

func NewDefaultEdge(id interface{}, from Node, to Node, weights map[string]float64) *DefaultEdge {
	return &DefaultEdge{
		id:      id,
		from:    from,
		to:      to,
		weights: weights,
	}
}

func (edge *DefaultEdge) GetId() interface{} {
	return edge.id
}

func (edge *DefaultEdge) From() Node {
	return edge.from
}

func (edge *DefaultEdge) To() Node {
	return edge.to
}

func (edge *DefaultEdge) GetWeight(kind string) (float64, error) {
	weight, ok := edge.weights[kind]
	if !ok {
		return 0, fmt.Errorf("weight kind %s not found", kind)
	}
	return weight, nil
}

func (edge *DefaultEdge) SetWeight(kind string, weight float64) {
	edge.weights[kind] = weight
}
