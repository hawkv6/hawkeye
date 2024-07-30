package domain

type ValueType int

const (
	ValueTypeUnspecified ValueType = iota
	ValueTypeMinValue
	ValueTypeMaxValue
	ValueTypeSFC
	ValueTypeFlexAlgoNr
)

func (vt ValueType) String() string {
	switch vt {
	case ValueTypeUnspecified:
		return "Unspecified"
	case ValueTypeMinValue:
		return "MinValue"
	case ValueTypeMaxValue:
		return "MaxValue"
	case ValueTypeSFC:
		return "SFC"
	case ValueTypeFlexAlgoNr:
		return "FlexAlgoNr"
	default:
		return "Unknown"
	}
}
