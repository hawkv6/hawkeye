package domain

import "github.com/go-playground/validator"

type NumberValue struct {
	BaseValue   `validate:"required"`
	numberValue int32 `validate:"required,min=0"`
}

func NewNumberValue(valueType ValueType, numberValue int32) (*NumberValue, error) {
	value := &NumberValue{
		BaseValue:   *NewBaseValue(valueType),
		numberValue: numberValue,
	}
	validate := validator.New()
	if err := validate.Struct(value); err != nil {
		return nil, err
	}
	return value, nil
}

func (numberValue *NumberValue) GetNumberValue() int32 {
	return numberValue.numberValue
}

func (numberValue *NumberValue) GetStringValue() string {
	return ""
}
