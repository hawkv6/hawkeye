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
		"MaxLinkBWKbps",
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

type WeightKey string

const (
	DefaultKey            WeightKey = "default"
	LatencyKey            WeightKey = "latency"
	JitterKey             WeightKey = "jitter"
	MaximumLinkBandwidth  WeightKey = "maximumLinkBandwidth"
	AvailableBandwidthKey WeightKey = "availableBandwidth"
	UtilizedBandwidthKey  WeightKey = "utilizedBandwith"
	PacketLossKey         WeightKey = "loss"
	normalizedLatencyKey  WeightKey = "normalizedLatency"
	normalizedJitterKey   WeightKey = "normalizedJitter"
	normalizedLossKey     WeightKey = "normalizedLoss"
)

const FlappingThreshold = 0.1 // If better path is found with 10% less cost, then it is considered as better path
