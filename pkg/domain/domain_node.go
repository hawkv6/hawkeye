package domain

import "github.com/go-playground/validator"

type Node interface {
	GetKey() string
	GetIgpRouterId() string
	GetName() string
}

type NodeInput struct {
	Key         *string `validate:"required"`
	IgpRouterId *string `validate:"required"`
	Name        *string `validate:"required"`
}

type DomainNode struct {
	key         string
	igpRouterId string
	name        string
}

func NewDomainNode(key, igpRouterId, name *string) (*DomainNode, error) {
	input := &NodeInput{
		Key:         key,
		IgpRouterId: igpRouterId,
		Name:        name,
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return nil, err
	}

	defaultNode := &DomainNode{
		key:         *key,
		igpRouterId: *igpRouterId,
		name:        *name,
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
