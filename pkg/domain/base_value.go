package domain

type BaseValue struct {
	valueType ValueType
}

func NewBaseValue(valueType ValueType) *BaseValue {
	return &BaseValue{valueType: valueType}
}

func (baseValue *BaseValue) GetValueType() ValueType {
	return baseValue.valueType
}
