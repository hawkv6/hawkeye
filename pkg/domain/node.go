package domain

type Node interface {
	GetKey() string
	GetIgpRouterId() string
	GetName() string
}
