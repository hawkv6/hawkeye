package domain

import (
	"github.com/go-playground/validator"
)

type StringValue struct {
	BaseValue
	stringValue string
}

type StringValueInput struct {
	StringValue *string `validate:"required"`
}

func NewStringValue(valueType ValueType, stringValue *string) (*StringValue, error) {
	stringValueInput := &StringValueInput{
		StringValue: stringValue,
	}
	validate := validator.New()
	if err := validate.Struct(stringValueInput); err != nil {
		return nil, err
	}
	value := &StringValue{
		BaseValue:   *NewBaseValue(valueType),
		stringValue: *stringValue,
	}
	return value, nil
}

func (stringValue *StringValue) GetNumberValue() int32 {
	return 0
}

func (stringValue *StringValue) GetStringValue() string {
	return stringValue.stringValue
}
