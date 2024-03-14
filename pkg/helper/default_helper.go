package helper

import (
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type DefaultHelper struct {
	log *logrus.Entry
}

func NewDefaultHelper() *DefaultHelper {
	return &DefaultHelper{
		log: logging.DefaultLogger.WithField("subsystem", subsystem),
	}
}

func (helper *DefaultHelper) GetLsNodeProperties() []string {
	nodeProperties := []string{"Key", "IgpRouterId", "Name"}
	helper.log.Debugln("LsNode properties", nodeProperties)
	return nodeProperties
}

func (helper *DefaultHelper) GetLsLinkProperties() []string {
	linkProperties := []string{
		"Key",
		"IgpRouterId",
		"RemoteIgpRouterId",
		"UnidirLinkDelay",
		"UnidirDelayVariation",
		"UnidirAvailableBW",
		"UnidirPacketLoss",
		"UnidirBWUtilization",
	}
	helper.log.Debugln("LsLink properties", linkProperties)
	return linkProperties
}

func (helper *DefaultHelper) GetLsPrefixProperties() []string {
	return []string{"Key", "IgpRouterId", "Prefix", "PrefixLen"}
}

func (helper *DefaultHelper) GetLsSrv6SidsProperties() []string {
	return []string{"Key", "IgpRouterId", "Srv6Sid"}
}

func (helper *DefaultHelper) GetLatencyKey() string {
	return "latency"
}

func (helper *DefaultHelper) GetJitterKey() string {
	return "jitter"
}

func (helper *DefaultHelper) GetAvailableBandwidthKey() string {
	return "availableBandwidth"
}

func (helper *DefaultHelper) GetUtilizedBandwidthKey() string {
	return "utilizedBandwith"
}

func (helper *DefaultHelper) GetPacketLossKey() string {
	return "loss"
}
