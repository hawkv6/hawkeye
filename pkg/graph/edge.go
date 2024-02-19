package graph

type Edge interface {
	GetId() interface{}
	From() Node
	To() Node
	GetWeight(kind string) (float64, error)
	SetWeight(kind string, weight float64)
}
