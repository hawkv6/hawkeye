package domain

import (
	"fmt"

	"github.com/go-playground/validator"
)

type NumberValue struct {
	BaseValue
	numberValue int32
}

type NumberValueBuilder struct {
	BaseValue   `validate:"required"`
	NumberValue int32 `validate:"required,min=1"`
}

func NewNumberValue(valueType ValueType, numberValue *int32) (*NumberValue, error) {
	if numberValue == nil {
		return nil, fmt.Errorf("Number value was not provided")
	}
	value := &NumberValueBuilder{
		BaseValue:   *NewBaseValue(valueType),
		NumberValue: *numberValue,
	}
	validate := validator.New()
	if err := validate.Struct(value); err != nil {
		return nil, err
	}
	newNumberValue := &NumberValue{
		BaseValue:   value.BaseValue,
		numberValue: value.NumberValue,
	}
	return newNumberValue, nil
}

func (numberValue *NumberValue) GetNumberValue() int32 {
	return numberValue.numberValue
}

func (numberValue *NumberValue) GetStringValue() string {
	return ""
}
