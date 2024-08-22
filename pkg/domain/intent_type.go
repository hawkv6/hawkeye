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
	IntentTypeLowUtilization
)

func (it IntentType) String() string {
	switch it {
	case IntentTypeUnspecified:
		return "Unspecified"
	case IntentTypeHighBandwidth:
		return "HighBandwidth"
	case IntentTypeLowBandwidth:
		return "LowBandwidth"
	case IntentTypeLowLatency:
		return "LowLatency"
	case IntentTypeLowPacketLoss:
		return "LowPacketLoss"
	case IntentTypeLowJitter:
		return "LowJitter"
	case IntentTypeFlexAlgo:
		return "FlexAlgo"
	case IntentTypeSFC:
		return "SFC"
	case IntentTypeLowUtilization:
		return "LowUtilization"
	default:
		return "Unknown"
	}
}
