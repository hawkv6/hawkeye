package helper

import (
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

const (
	PropertyKey                            = "Key"
	PropertyIgpRouterId                    = "IgpRouterId"
	PropertyIgpMetric                      = "IgpMetric"
	PropertyName                           = "Name"
	PropertyRemoteIgpRouterId              = "RemoteIgpRouterId"
	PropertyUnidirLinkDelay                = "UnidirLinkDelay"
	PropertyUnidirDelayVariation           = "UnidirDelayVariation"
	PropertyMaxLinkBwKbps                  = "MaxLinkBwKbps"
	PropertyUnidirAvailableBw              = "UnidirAvailableBw"
	PropertyUnidirPacketLoss               = "UnidirPacketLossPercentage"
	PropertyUnidirBwUtilization            = "UnidirBwUtilization"
	PropertyNormalizedUnidirLinkDelay      = "NormalizedUnidirLinkDelay"
	PropertyNormalizedUnidirDelayVariation = "NormalizedUnidirDelayVariation"
	PropertyNormalizedUnidirPacketLoss     = "NormalizedUnidirPacketLoss"
	PropertyPrefix                         = "Prefix"
	PropertyPrefixLen                      = "PrefixLen"
	PropertySrv6Sid                        = "Srv6Sid"
	PropertySrAlgorithm                    = "SrAlgorithm"
	PropertySrv6Locator                    = "Srv6Locator"
	PropertySrv6EndpointBehavior           = "Srv6EndpointBehavior"
)

type WeightKey string

const (
	UndefinedKey            WeightKey = ""
	IgpMetricKey            WeightKey = PropertyIgpMetric
	LatencyKey              WeightKey = PropertyUnidirLinkDelay
	JitterKey               WeightKey = PropertyUnidirDelayVariation
	MaximumLinkBandwidth    WeightKey = PropertyMaxLinkBwKbps
	AvailableBandwidthKey   WeightKey = PropertyUnidirAvailableBw
	UtilizedBandwidthKey    WeightKey = PropertyUnidirBwUtilization
	PacketLossKey           WeightKey = PropertyUnidirPacketLoss
	NormalizedLatencyKey    WeightKey = PropertyNormalizedUnidirLinkDelay
	NormalizedJitterKey     WeightKey = PropertyNormalizedUnidirDelayVariation
	NormalizedPacketLossKey WeightKey = PropertyNormalizedUnidirPacketLoss
)

const RollingWindowSize = 5 // make it configurable via env variable

type DefaultHelper struct {
	log                  *logrus.Entry
	lsNodeProperties     []string
	lsLinkProperties     []string
	lsPrefixProperties   []string
	lsSrv6SidsProperties []string
}

func NewDefaultHelper() *DefaultHelper {
	return &DefaultHelper{
		log:                  logging.DefaultLogger.WithField("subsystem", subsystem),
		lsNodeProperties:     []string{PropertyKey, PropertyIgpRouterId, PropertyName, PropertySrAlgorithm},
		lsLinkProperties:     []string{PropertyKey, PropertyIgpRouterId, PropertyRemoteIgpRouterId, PropertyIgpMetric, PropertyUnidirLinkDelay, PropertyUnidirDelayVariation, PropertyMaxLinkBwKbps, PropertyUnidirAvailableBw, PropertyUnidirPacketLoss, PropertyUnidirBwUtilization, PropertyNormalizedUnidirLinkDelay, PropertyNormalizedUnidirDelayVariation, PropertyNormalizedUnidirPacketLoss},
		lsPrefixProperties:   []string{PropertyKey, PropertyIgpRouterId, PropertyPrefix, PropertyPrefixLen, PropertySrv6Locator},
		lsSrv6SidsProperties: []string{PropertyKey, PropertyIgpRouterId, PropertySrv6Sid, PropertySrv6EndpointBehavior},
	}
}

func (helper *DefaultHelper) GetLsNodeProperties() []string {
	helper.log.Debugln("LsNode properties", helper.lsNodeProperties)
	return helper.lsNodeProperties
}

func (helper *DefaultHelper) GetLsLinkProperties() []string {
	helper.log.Debugln("LsLink properties", helper.lsLinkProperties)
	return helper.lsLinkProperties
}

func (helper *DefaultHelper) GetLsPrefixProperties() []string {
	helper.log.Debugln("LsPrefix properties", helper.lsPrefixProperties)
	return helper.lsPrefixProperties
}

func (helper *DefaultHelper) GetLsSrv6SidsProperties() []string {
	helper.log.Debugln("LsSrv6Sids properties", helper.lsSrv6SidsProperties)
	return helper.lsSrv6SidsProperties
}
