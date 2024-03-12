package domain

import (
	"github.com/go-playground/validator"
)

type LinkInput struct {
	Key                        *string  `validate:"required"`
	IgpRouterId                *string  `validate:"required"`
	RemoteIgpRouterId          *string  `validate:"required"`
	UnidirLinkDelay            *uint32  `validate:"required"`
	UnidirDelayVariation       *uint32  `validate:"required"`
	UnidirAvailableBw          *uint32  `validate:"required"`
	UnidirPacketLoss           *float32 `validate:"required,min=0"`
	UnidirBandwidthUtilization *uint32  `validate:"required"`
}

type DefaultLink struct {
	key                        string
	igpRouterId                string
	remoteIgpRouterId          string
	unidirLinkDelay            uint32
	unidirDelayVariation       uint32
	unidirAvailableBandwidth   uint32
	unidirPacketLoss           float32
	unidirBandwidthUtilization uint32
}

func NewDefaultLink(key, igpRouterId, remoteIgpRouterId *string, unidirLinkDelay, unidirDelayVariation, unidirAvailableBandwidth, unidirBandwidthUtilization *uint32, unidirPacketLoss *float32) (*DefaultLink, error) {
	input := &LinkInput{
		Key:                        key,
		IgpRouterId:                igpRouterId,
		RemoteIgpRouterId:          remoteIgpRouterId,
		UnidirLinkDelay:            unidirLinkDelay,
		UnidirDelayVariation:       unidirDelayVariation,
		UnidirAvailableBw:          unidirAvailableBandwidth,
		UnidirPacketLoss:           unidirPacketLoss,
		UnidirBandwidthUtilization: unidirBandwidthUtilization,
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return nil, err
	}

	defaultLink := &DefaultLink{
		key:                        *key,
		igpRouterId:                *igpRouterId,
		remoteIgpRouterId:          *remoteIgpRouterId,
		unidirLinkDelay:            *unidirLinkDelay,
		unidirDelayVariation:       *unidirDelayVariation,
		unidirAvailableBandwidth:   *unidirAvailableBandwidth,
		unidirPacketLoss:           *unidirPacketLoss,
		unidirBandwidthUtilization: *unidirBandwidthUtilization,
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

func (l *DefaultLink) GetUnidirLinkDelay() uint32 {
	return l.unidirLinkDelay
}

func (l *DefaultLink) GetUnidirDelayVariation() uint32 {
	return l.unidirDelayVariation
}

func (l *DefaultLink) GetUnidirAvailableBandwidth() uint32 {
	return l.unidirAvailableBandwidth
}

func (l *DefaultLink) GetUnidirPacketLoss() float32 {
	return l.unidirPacketLoss
}

func (l *DefaultLink) GetUnidirBandwidthUtilization() uint32 {
	return l.unidirBandwidthUtilization
}
