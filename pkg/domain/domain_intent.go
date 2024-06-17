package domain

import (
	"strconv"

	"github.com/go-playground/validator"
)

type Intent interface {
	GetIntentType() IntentType
	GetValues() []Value
	Serialize() string
}

type DefaultIntent struct {
	intentType IntentType `validate:"required, min=0"`
	values     []Value
}

func NewDefaultIntent(intentType IntentType, values []Value) (*DefaultIntent, error) {
	intent := &DefaultIntent{
		intentType: intentType,
		values:     values,
	}
	validate := validator.New()
	err := validate.Struct(intent)
	if err != nil {
		return nil, err
	}
	return intent, nil
}

func (intent *DefaultIntent) GetIntentType() IntentType {
	return intent.intentType
}

func (intent *DefaultIntent) GetValues() []Value {
	return intent.values
}

func (intent *DefaultIntent) convertValue(value Value) string {
	valueType := value.GetValueType()
	switch valueType {
	case ValueTypeMinValue:
		return valueType.String() + ":" + strconv.Itoa(int(value.GetNumberValue()))
	case ValueTypeMaxValue:
		return valueType.String() + ":" + strconv.Itoa(int(value.GetNumberValue()))
	case ValueTypeFlexAlgoNr:
		return valueType.String() + ":" + strconv.Itoa((int(value.GetNumberValue())))
	case ValueTypeSFC:
		return valueType.String() + ":" + value.GetStringValue()
	default:
		return IntentTypeUnspecified.String()
	}
}

func (intent *DefaultIntent) Serialize() string {
	if len(intent.values) == 0 {
		return intent.intentType.String()
	}
	serialization := intent.intentType.String() + ","
	for i := 0; i < len(intent.values); i++ {
		if i == len(intent.values)-1 {
			serialization += intent.convertValue(intent.values[i])
		} else {
			serialization += intent.convertValue(intent.values[i]) + ","
		}
	}
	return serialization
}
