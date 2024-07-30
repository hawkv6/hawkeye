package domain

type Value interface {
	GetValueType() ValueType
	GetNumberValue() int32
	GetStringValue() string
}

type BaseValue struct {
	valueType ValueType
}

func NewBaseValue(valueType ValueType) *BaseValue {
	return &BaseValue{valueType: valueType}
}

func (baseValue *BaseValue) GetValueType() ValueType {
	return baseValue.valueType
}
