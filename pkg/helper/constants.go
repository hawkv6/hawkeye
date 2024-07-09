package helper

const subsystem = "helper"

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
