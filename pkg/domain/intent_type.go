package domain

type IntentType int

const (
	IntentTypeUnspecified IntentType = iota
	IntentTypeHighBandwidth
	IntentTypeLowBandwidth
	IntentTypeLowLatency
	IntentTypeLowPacketLoss
	IntentTypeLowJitter
	IntentTypeFlexAlgo
	IntentTypeSFC
)
