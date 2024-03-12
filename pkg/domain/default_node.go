package domain

import "github.com/go-playground/validator"

type NodeInput struct {
	Key         *string `validate:"required"`
	IgpRouterId *string `validate:"required"`
	Name        *string `validate:"required"`
}

type DefaultNode struct {
	key         string
	igpRouterId string
	name        string
}

func NewDefaultNode(key, igpRouterId, name *string) (*DefaultNode, error) {
	input := &NodeInput{
		Key:         key,
		IgpRouterId: igpRouterId,
		Name:        name,
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return nil, err
	}

	defaultNode := &DefaultNode{
		key:         *key,
		igpRouterId: *igpRouterId,
		name:        *name,
	}
	return defaultNode, nil
}

func (n *DefaultNode) GetKey() string {
	return n.key
}

func (n *DefaultNode) GetIgpRouterId() string {
	return n.igpRouterId
}

func (n *DefaultNode) GetName() string {
	return n.name
}
