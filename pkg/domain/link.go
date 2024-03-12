package domain

type Link interface {
	GetKey() string
	GetIgpRouterId() string
	GetRemoteIgpRouterId() string
	GetUnidirLinkDelay() uint32
	GetUnidirDelayVariation() uint32
	GetUnidirAvailableBandwidth() uint32
	GetUnidirPacketLoss() float32
	GetUnidirBandwidthUtilization() uint32
}
