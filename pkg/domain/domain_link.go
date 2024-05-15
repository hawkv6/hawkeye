package domain

import (
	"github.com/go-playground/validator"
)

type Link interface {
	GetKey() string
	GetIgpRouterId() string
	GetRemoteIgpRouterId() string
	GetUnidirLinkDelay() uint32
	GetUnidirDelayVariation() uint32
	GetMaxLinkBWKbps() uint64
	GetUnidirAvailableBandwidth() uint32
	GetUnidirPacketLoss() float32
	GetUnidirBandwidthUtilization() uint32
}

type LinkInput struct {
	Key                        *string  `validate:"required"`
	IgpRouterId                *string  `validate:"required"`
	RemoteIgpRouterId          *string  `validate:"required"`
	UnidirLinkDelay            *uint32  `validate:"required"`
	UnidirDelayVariation       *uint32  `validate:"required"`
	MaxLinkBWKbps              *uint64  `validate:"required,min=1"`
	UnidirAvailableBw          *uint32  `validate:"required"`
	UnidirPacketLoss           *float32 `validate:"required,min=0"`
	UnidirBandwidthUtilization *uint32  `validate:"required"`
}

type DomainLink struct {
	key                        string
	igpRouterId                string
	remoteIgpRouterId          string
	unidirLinkDelay            uint32
	unidirDelayVariation       uint32
	maxLinkBWKbps              uint64
	unidirAvailableBandwidth   uint32
	unidirPacketLoss           float32
	unidirBandwidthUtilization uint32
}

func NewDomainLink(key, igpRouterId, remoteIgpRouterId *string, unidirLinkDelay, unidirDelayVariation *uint32, maxLinkBWKbps *uint64, unidirAvailableBandwidth, unidirBandwidthUtilization *uint32, unidirPacketLoss *float32) (*DomainLink, error) {
	input := &LinkInput{
		Key:                        key,
		IgpRouterId:                igpRouterId,
		RemoteIgpRouterId:          remoteIgpRouterId,
		UnidirLinkDelay:            unidirLinkDelay,
		UnidirDelayVariation:       unidirDelayVariation,
		MaxLinkBWKbps:              maxLinkBWKbps,
		UnidirAvailableBw:          unidirAvailableBandwidth,
		UnidirPacketLoss:           unidirPacketLoss,
		UnidirBandwidthUtilization: unidirBandwidthUtilization,
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return nil, err
	}

	defaultLink := &DomainLink{
		key:                        *key,
		igpRouterId:                *igpRouterId,
		remoteIgpRouterId:          *remoteIgpRouterId,
		unidirLinkDelay:            *unidirLinkDelay,
		unidirDelayVariation:       *unidirDelayVariation,
		maxLinkBWKbps:              *maxLinkBWKbps,
		unidirAvailableBandwidth:   *unidirAvailableBandwidth,
		unidirPacketLoss:           *unidirPacketLoss,
		unidirBandwidthUtilization: *unidirBandwidthUtilization,
	}

	return defaultLink, nil
}

func (l *DomainLink) GetKey() string {
	return l.key
}

func (l *DomainLink) GetIgpRouterId() string {
	return l.igpRouterId
}

func (l *DomainLink) GetRemoteIgpRouterId() string {
	return l.remoteIgpRouterId
}

func (l *DomainLink) GetUnidirLinkDelay() uint32 {
	return l.unidirLinkDelay
}

func (l *DomainLink) GetUnidirDelayVariation() uint32 {
	return l.unidirDelayVariation
}

func (l *DomainLink) GetMaxLinkBWKbps() uint64 {
	return l.maxLinkBWKbps
}

func (l *DomainLink) GetUnidirAvailableBandwidth() uint32 {
	return l.unidirAvailableBandwidth
}

func (l *DomainLink) GetUnidirPacketLoss() float32 {
	return l.unidirPacketLoss
}

func (l *DomainLink) GetUnidirBandwidthUtilization() uint32 {
	return l.unidirBandwidthUtilization
}
