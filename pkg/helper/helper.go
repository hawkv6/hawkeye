package helper

const subsystem = "helper"

type Helper interface {
	GetLsNodeProperties() []string
	GetLsLinkProperties() []string
	GetLsPrefixProperties() []string
	GetLsSrv6SidsProperties() []string
	GetLatencyKey() string
	GetJitterKey() string
	GetAvailableBandwidthKey() string
	GetUtilizedBandwidthKey() string
	GetPacketLossKey() string
}
