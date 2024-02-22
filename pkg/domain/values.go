package domain

type Value interface {
	GetValueType() ValueType
	GetNumberValue() int32
	GetStringValue() string
}
