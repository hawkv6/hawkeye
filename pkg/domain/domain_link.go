package domain

import (
	"github.com/go-playground/validator"
)

type Link interface {
	GetKey() string
	GetIgpRouterId() string
	GetRemoteIgpRouterId() string
	GetIgpMetric() uint32
	GetUnidirLinkDelay() uint32
	GetUnidirDelayVariation() uint32
	GetMaxLinkBWKbps() uint64
	GetUnidirAvailableBandwidth() uint32
	GetUnidirPacketLoss() float64
	GetUnidirBandwidthUtilization() uint32
	SetNormalizedUnidirLinkDelay(float64)
	SetNormalizedUnidirDelayVariation(float64)
	SetNormalizedPacketLoss(float64)
	GetNormalizedUnidirLinkDelay() float64
	GetNormalizedUnidirDelayVariation() float64
	GetNormalizedUnidirPacketLoss() float64
}

type LinkInput struct {
	Key                            *string  `validate:"required"`
	IgpRouterId                    *string  `validate:"required"`
	RemoteIgpRouterId              *string  `validate:"required"`
	IgpMetric                      *uint32  `validate:"required"`
	UnidirLinkDelay                *uint32  `validate:"required"`
	UnidirDelayVariation           *uint32  `validate:"required"`
	MaxLinkBWKbps                  *uint64  `validate:"required,min=1"`
	UnidirAvailableBw              *uint32  `validate:"required"`
	UnidirPacketLoss               *float64 `validate:"required,min=0"`
	UnidirBandwidthUtilization     *uint32  `validate:"required"`
	NormalizedUnidirLinkDelay      *float64 `validate:"required,min=0,max=1"`
	NormalizedUnidirDelayVariation *float64 `validate:"required,min=0,max=1"`
	NormalizedUnidirPacketLoss     *float64 `validate:"required,min=0,max=1"`
}

type DomainLink struct {
	key                            string
	igpRouterId                    string
	remoteIgpRouterId              string
	igpMetric                      uint32
	unidirLinkDelay                uint32
	unidirDelayVariation           uint32
	maxLinkBWKbps                  uint64
	unidirAvailableBandwidth       uint32
	unidirPacketLoss               float64
	unidirBandwidthUtilization     uint32
	normalizedUnidirLinkDelay      float64
	normalizedUnidirDelayVariation float64
	normalizedUnidirPacketLoss     float64
}

func NewDomainLink(key, igpRouterId, remoteIgpRouterId *string, igpMetric, unidirLinkDelay, unidirDelayVariation *uint32, maxLinkBWKbps *uint64, unidirAvailableBandwidth, unidirBandwidthUtilization *uint32, unidirPacketLoss, normalizedUnidirLinkDelay, normalizedUnidirDelayVariation, normalizedUnidirPacketLoss *float64) (*DomainLink, error) {
	input := &LinkInput{
		Key:                            key,
		IgpRouterId:                    igpRouterId,
		RemoteIgpRouterId:              remoteIgpRouterId,
		IgpMetric:                      igpMetric,
		UnidirLinkDelay:                unidirLinkDelay,
		UnidirDelayVariation:           unidirDelayVariation,
		MaxLinkBWKbps:                  maxLinkBWKbps,
		UnidirAvailableBw:              unidirAvailableBandwidth,
		UnidirPacketLoss:               unidirPacketLoss,
		UnidirBandwidthUtilization:     unidirBandwidthUtilization,
		NormalizedUnidirLinkDelay:      normalizedUnidirLinkDelay,
		NormalizedUnidirDelayVariation: normalizedUnidirDelayVariation,
		NormalizedUnidirPacketLoss:     normalizedUnidirPacketLoss,
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return nil, err
	}

	defaultLink := &DomainLink{
		key:                            *key,
		igpRouterId:                    *igpRouterId,
		igpMetric:                      *igpMetric,
		remoteIgpRouterId:              *remoteIgpRouterId,
		unidirLinkDelay:                *unidirLinkDelay,
		unidirDelayVariation:           *unidirDelayVariation,
		maxLinkBWKbps:                  *maxLinkBWKbps,
		unidirAvailableBandwidth:       *unidirAvailableBandwidth,
		unidirPacketLoss:               *unidirPacketLoss,
		unidirBandwidthUtilization:     *unidirBandwidthUtilization,
		normalizedUnidirLinkDelay:      *normalizedUnidirLinkDelay,
		normalizedUnidirDelayVariation: *normalizedUnidirDelayVariation,
		normalizedUnidirPacketLoss:     *normalizedUnidirPacketLoss,
	}

	return defaultLink, nil
}

func (link *DomainLink) GetKey() string {
	return link.key
}

func (link *DomainLink) GetIgpRouterId() string {
	return link.igpRouterId
}

func (link *DomainLink) GetRemoteIgpRouterId() string {
	return link.remoteIgpRouterId
}

func (link *DomainLink) GetIgpMetric() uint32 {
	return link.igpMetric
}

func (link *DomainLink) GetUnidirLinkDelay() uint32 {
	return link.unidirLinkDelay
}

func (link *DomainLink) GetUnidirDelayVariation() uint32 {
	return link.unidirDelayVariation
}

func (link *DomainLink) GetMaxLinkBWKbps() uint64 {
	return link.maxLinkBWKbps
}

func (link *DomainLink) GetUnidirAvailableBandwidth() uint32 {
	return link.unidirAvailableBandwidth
}

func (link *DomainLink) GetUnidirPacketLoss() float64 {
	return link.unidirPacketLoss
}

func (link *DomainLink) GetUnidirBandwidthUtilization() uint32 {
	return link.unidirBandwidthUtilization
}

func (link *DomainLink) SetNormalizedUnidirLinkDelay(normalizedUnidirLinkDelay float64) {
	link.normalizedUnidirLinkDelay = normalizedUnidirLinkDelay
}

func (link *DomainLink) SetNormalizedUnidirDelayVariation(normalizedUnidirDelayVariation float64) {
	link.normalizedUnidirDelayVariation = normalizedUnidirDelayVariation
}

func (link *DomainLink) SetNormalizedPacketLoss(normalizedPacketLoss float64) {
	link.normalizedUnidirPacketLoss = normalizedPacketLoss
}

func (link *DomainLink) GetNormalizedUnidirLinkDelay() float64 {
	return link.normalizedUnidirLinkDelay
}

func (link *DomainLink) GetNormalizedUnidirDelayVariation() float64 {
	return link.normalizedUnidirDelayVariation
}

func (link *DomainLink) GetNormalizedUnidirPacketLoss() float64 {
	return link.normalizedUnidirPacketLoss
}
