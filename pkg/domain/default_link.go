package domain

import "github.com/go-playground/validator"

type DefaultLink struct {
	key               string  `validate:"required"`
	igpRouterId       string  `validate:"required"`
	remoteIgpRouterId string  `validate:"required"`
	unidirLinkDelay   float64 `validate:"required,min=0"`
}

func NewDefaultLink(key, igpRouterId, remoteIgpRouterId string, unidirLinkDelay uint32) (*DefaultLink, error) {
	defaultLink := &DefaultLink{
		key:               key,
		igpRouterId:       igpRouterId,
		remoteIgpRouterId: remoteIgpRouterId,
		unidirLinkDelay:   float64(unidirLinkDelay),
	}
	validate := validator.New()
	if err := validate.Struct(defaultLink); err != nil {
		return nil, err
	}
	return defaultLink, nil
}

func (l *DefaultLink) GetKey() string {
	return l.key
}

func (l *DefaultLink) GetIgpRouterId() string {
	return l.igpRouterId
}

func (l *DefaultLink) GetRemoteIgpRouterId() string {
	return l.remoteIgpRouterId
}

func (l *DefaultLink) GetUnidirLinkDelay() float64 {
	return l.unidirLinkDelay
}
