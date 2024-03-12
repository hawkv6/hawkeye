package domain

import "github.com/go-playground/validator"

type PrefixInput struct {
	Key          *string `validate:"required"`
	IgpRouterId  *string `validate:"required"`
	Prefix       *string `validate:"required"`
	PrefixLength *int32  `validate:"required,min=0,max=128"`
}
type DefaultPrefix struct {
	key          string
	igpRouterId  string
	prefix       string
	prefixLength uint8
}

func NewDefaultPrefix(key, igpRouterId, prefix *string, prefixLength *int32) (*DefaultPrefix, error) {
	input := &PrefixInput{
		Key:          key,
		IgpRouterId:  igpRouterId,
		Prefix:       prefix,
		PrefixLength: prefixLength,
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return nil, err
	}

	lsPrefixAdapter := &DefaultPrefix{
		key:          *key,
		igpRouterId:  *igpRouterId,
		prefix:       *prefix,
		prefixLength: uint8(*prefixLength),
	}
	return lsPrefixAdapter, nil
}

func (p *DefaultPrefix) GetKey() string {
	return p.key
}

func (p *DefaultPrefix) GetIgpRouterId() string {
	return p.igpRouterId
}

func (p *DefaultPrefix) GetPrefix() string {
	return p.prefix
}

func (p *DefaultPrefix) GetPrefixLength() uint8 {
	return p.prefixLength
}
