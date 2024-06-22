package domain

import "github.com/go-playground/validator"

type Node interface {
	GetKey() string
	GetIgpRouterId() string
	GetName() string
	GetSrAlgorithm() []uint32
}

type NodeInput struct {
	Key         *string  `validate:"required"`
	IgpRouterId *string  `validate:"required"`
	Name        *string  `validate:"required"`
	SrAlgorithm []uint32 `validate:"required"`
}

type DomainNode struct {
	key         string
	igpRouterId string
	name        string
	srAlgorithm []uint32
}

func NewDomainNode(key, igpRouterId, name *string, srAlgorihm []uint32) (*DomainNode, error) {
	input := &NodeInput{
		Key:         key,
		IgpRouterId: igpRouterId,
		Name:        name,
		SrAlgorithm: srAlgorihm,
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return nil, err
	}

	defaultNode := &DomainNode{
		key:         *key,
		igpRouterId: *igpRouterId,
		name:        *name,
		srAlgorithm: srAlgorihm,
	}
	return defaultNode, nil
}

func (n *DomainNode) GetKey() string {
	return n.key
}

func (n *DomainNode) GetIgpRouterId() string {
	return n.igpRouterId
}

func (n *DomainNode) GetName() string {
	return n.name
}

func (n *DomainNode) GetSrAlgorithm() []uint32 {
	return n.srAlgorithm
}
