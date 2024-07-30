package helper

import "github.com/hawkv6/hawkeye/pkg/logging"

var log = logging.DefaultLogger.WithField("subsystem", subsystem)

func GetLsNodeProperties() []string {
	lsNodeProperties := []string{PropertyKey, PropertyIgpRouterId, PropertyName, PropertySrAlgorithm}
	log.Debugln("LsNode properties", lsNodeProperties)
	return lsNodeProperties
}

func GetLsLinkProperties() []string {
	lsLinkProperties := []string{PropertyKey, PropertyIgpRouterId, PropertyRemoteIgpRouterId, PropertyIgpMetric, PropertyUnidirLinkDelay, PropertyUnidirDelayVariation, PropertyMaxLinkBwKbps, PropertyUnidirAvailableBw, PropertyUnidirPacketLoss, PropertyUnidirBwUtilization, PropertyNormalizedUnidirLinkDelay, PropertyNormalizedUnidirDelayVariation, PropertyNormalizedUnidirPacketLoss}
	log.Debugln("LsLink properties", lsLinkProperties)
	return lsLinkProperties
}

func GetLsPrefixProperties() []string {
	lsPrefixProperties := []string{PropertyKey, PropertyIgpRouterId, PropertyPrefix, PropertyPrefixLen, PropertySrv6Locator}
	log.Debugln("LsPrefix properties", lsPrefixProperties)
	return lsPrefixProperties
}

func GetLsSrv6SidsProperties() []string {
	lsSrv6SidsProperties := []string{PropertyKey, PropertyIgpRouterId, PropertySrv6Sid, PropertySrv6EndpointBehavior}
	log.Debugln("LsSrv6Sids properties", lsSrv6SidsProperties)
	return lsSrv6SidsProperties
}
