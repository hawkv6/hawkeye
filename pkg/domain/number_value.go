package domain

import (
	"fmt"

	"github.com/go-playground/validator"
)

type NumberValue struct {
	BaseValue
	numberValue int32
}

type NumberValueInput struct {
	BaseValue   `validate:"required"`
	NumberValue int32 `validate:"required,min=1"`
}

func NewNumberValue(valueType ValueType, numberValue *int32) (*NumberValue, error) {
	if numberValue == nil {
		return nil, fmt.Errorf("Number value was not provided")
	}
	numberValueInput := &NumberValueInput{
		BaseValue:   *NewBaseValue(valueType),
		NumberValue: *numberValue,
	}
	validate := validator.New()
	if err := validate.Struct(numberValueInput); err != nil {
		return nil, err
	}
	newNumberValue := &NumberValue{
		BaseValue:   numberValueInput.BaseValue,
		numberValue: numberValueInput.NumberValue,
	}
	return newNumberValue, nil
}

func (numberValue *NumberValue) GetNumberValue() int32 {
	return numberValue.numberValue
}

func (numberValue *NumberValue) GetStringValue() string {
	return ""
}
