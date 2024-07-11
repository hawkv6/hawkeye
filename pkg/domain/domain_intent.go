package domain

import (
	"strconv"
)

type Intent interface {
	GetIntentType() IntentType
	GetValues() []Value
	Serialize() string
}

type DomainIntent struct {
	intentType IntentType
	values     []Value
}

func NewDomainIntent(intentType IntentType, values []Value) *DomainIntent {
	return &DomainIntent{
		intentType: intentType,
		values:     values,
	}
}

func (intent *DomainIntent) GetIntentType() IntentType {
	return intent.intentType
}

func (intent *DomainIntent) GetValues() []Value {
	return intent.values
}

func (intent *DomainIntent) convertValue(value Value) string {
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

func (intent *DomainIntent) Serialize() string {
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
