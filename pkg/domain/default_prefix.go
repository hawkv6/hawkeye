package domain

import "github.com/go-playground/validator"

type DefaultPrefix struct {
	key          string `validate:"required"`
	igpRouterId  string `validate:"required"`
	prefix       string `validate:"required"`
	prefixLength uint8  `validate:"required,min=0,max=128"`
}

func NewDefaultPrefix(key, igpRouterId, prefix string, prefixLength int32) (*DefaultPrefix, error) {
	lsPrefixAdapter := &DefaultPrefix{
		key:          key,
		igpRouterId:  igpRouterId,
		prefix:       prefix,
		prefixLength: uint8(prefixLength),
	}
	validate := validator.New()
	if err := validate.Struct(lsPrefixAdapter); err != nil {
		return nil, err
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
