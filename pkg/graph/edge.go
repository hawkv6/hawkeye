package graph

import "github.com/hawkv6/hawkeye/pkg/helper"

type Edge interface {
	GetId() interface{}
	From() Node
	To() Node
	GetAllWeights() map[helper.WeightKey]float64
	GetWeight(kind helper.WeightKey) (float64, error)
	SetWeight(kind helper.WeightKey, weight float64)
}
