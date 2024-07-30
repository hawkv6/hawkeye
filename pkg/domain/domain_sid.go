package domain

import "github.com/go-playground/validator"

type Sid interface {
	GetKey() string
	GetIgpRouterId() string
	GetSid() string
	GetAlgorithm() uint32
}

type SidInput struct {
	Key         *string `validate:"required"`
	IgpRouterId *string `validate:"required"`
	Sid         *string `validate:"required"`
	Algorithm   *uint32 `validate:"required,min=0,max=255"`
}

type DomainSid struct {
	key         string
	igpRouterId string
	sid         string
	algorithm   uint32
}

func NewDomainSid(key, igpRouterId, sid *string, algorithm *uint32) (*DomainSid, error) {
	input := &SidInput{
		Key:         key,
		IgpRouterId: igpRouterId,
		Sid:         sid,
		Algorithm:   algorithm,
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return nil, err
	}

	defaultSid := &DomainSid{
		key:         *key,
		igpRouterId: *igpRouterId,
		sid:         *sid,
		algorithm:   *algorithm,
	}
	return defaultSid, nil
}

func (sid *DomainSid) GetKey() string {
	return sid.key
}

func (sid *DomainSid) GetIgpRouterId() string {
	return sid.igpRouterId
}

func (sid *DomainSid) GetSid() string {
	return sid.sid
}

func (sid *DomainSid) GetAlgorithm() uint32 {
	return sid.algorithm
}
