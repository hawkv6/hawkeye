package domain

import "github.com/go-playground/validator"

type StringValue struct {
	BaseValue
	stringValue string `validate:"required"`
}

func NewStringValue(valueType ValueType, stringValue string) (*StringValue, error) {
	value := &StringValue{
		BaseValue:   *NewBaseValue(valueType),
		stringValue: stringValue,
	}
	validate := validator.New()
	if err := validate.Struct(value); err != nil {
		return nil, err
	}
	return value, nil
}

func (stringValue *StringValue) GetNumberValue() int32 {
	return 0
}

func (stringValue *StringValue) GetStringValue() string {
	return stringValue.stringValue
}
