package domain

import "github.com/go-playground/validator"

type Prefix interface {
	GetKey() string
	GetIgpRouterId() string
	GetPrefix() string
	GetPrefixLength() uint8
}

type PrefixInput struct {
	Key          *string `validate:"required"`
	IgpRouterId  *string `validate:"required"`
	Prefix       *string `validate:"required"`
	PrefixLength *int32  `validate:"required,min=0,max=128"`
}
type DomainPrefix struct {
	key          string
	igpRouterId  string
	prefix       string
	prefixLength uint8
}

func NewDomainPrefix(key, igpRouterId, prefix *string, prefixLength *int32) (*DomainPrefix, error) {
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

	lsPrefixAdapter := &DomainPrefix{
		key:          *key,
		igpRouterId:  *igpRouterId,
		prefix:       *prefix,
		prefixLength: uint8(*prefixLength),
	}
	return lsPrefixAdapter, nil
}

func (prefix *DomainPrefix) GetKey() string {
	return prefix.key
}

func (prefix *DomainPrefix) GetIgpRouterId() string {
	return prefix.igpRouterId
}

func (prefix *DomainPrefix) GetPrefix() string {
	return prefix.prefix
}

func (prefix *DomainPrefix) GetPrefixLength() uint8 {
	return prefix.prefixLength
}
