package graph

type Edge interface {
	GetId() interface{}
	From() Node
	To() Node
	GetAllWeights() map[string]float64
	GetWeight(kind string) (float64, error)
	SetWeight(kind string, weight float64)
}
