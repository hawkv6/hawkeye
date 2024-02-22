package domain

type ValueType int

const (
	ValueTypeUnspecified ValueType = iota
	ValueTypeMinValue
	ValueTypeMaxValue
	ValueTypeSFC
	ValueTypeFlexAlgoNr
)
