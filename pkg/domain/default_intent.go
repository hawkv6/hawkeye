package domain

import "github.com/go-playground/validator"

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

func (intent DefaultIntent) GetIntentType() IntentType {
	return intent.intentType
}

func (intent DefaultIntent) GetValues() []Value {
	return intent.values
}
