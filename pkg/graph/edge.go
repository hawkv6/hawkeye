package graph

import "github.com/hawkv6/hawkeye/pkg/helper"

type Edge interface {
	GetId() string
	From() Node
	To() Node
	GetAllWeights() map[helper.WeightKey]float64
	GetWeight(kind helper.WeightKey) float64
	SetWeight(kind helper.WeightKey, weight float64)
	GetFlexibleAlgorithms() map[uint32]struct{}
	UpdateFlexibleAlgorithms()
}
