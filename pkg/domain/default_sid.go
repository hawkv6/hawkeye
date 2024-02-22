package domain

import "github.com/go-playground/validator"

type DefaultSid struct {
	key         string `validate:"required"`
	igpRouterId string `validate:"required"`
	sid         string `validate:"required"`
}

func NewDefaultSid(key, igpRouterId, sid string) (*DefaultSid, error) {
	defaultSid := &DefaultSid{
		key:         key,
		igpRouterId: igpRouterId,
		sid:         sid,
	}
	validate := validator.New()
	if err := validate.Struct(defaultSid); err != nil {
		return nil, err
	}
	return defaultSid, nil
}

func (s *DefaultSid) GetKey() string {
	return s.key
}

func (s *DefaultSid) GetIgpRouterId() string {
	return s.igpRouterId
}

func (s *DefaultSid) GetSid() string {
	return s.sid
}
