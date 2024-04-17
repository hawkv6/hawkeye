package domain

import "github.com/go-playground/validator"

type Sid interface {
	GetKey() string
	GetIgpRouterId() string
	GetSid() string
}

type SidInput struct {
	Key         *string `validate:"required"`
	IgpRouterId *string `validate:"required"`
	Sid         *string `validate:"required"`
}

type DomainSid struct {
	key         string `validate:"required"`
	igpRouterId string `validate:"required"`
	sid         string `validate:"required"`
}

func NewDomainSid(key, igpRouterId, sid *string) (*DomainSid, error) {
	input := &SidInput{
		Key:         key,
		IgpRouterId: igpRouterId,
		Sid:         sid,
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return nil, err
	}

	defaultSid := &DomainSid{
		key:         *key,
		igpRouterId: *igpRouterId,
		sid:         *sid,
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
