package domain

type Intent interface {
	GetIntentType() IntentType
	GetValues() []Value
	Serialize() string
}
