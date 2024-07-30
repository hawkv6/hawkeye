package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestNewDomainIntent(t *testing.T) {
	tests := []struct {
		name       string
		intentType IntentType
		values     []Value
	}{
		{
			name:       "Test NewDomainIntent IntentTypeFlexAlgo",
			intentType: IntentTypeFlexAlgo,
			values:     []Value{},
		},
		{
			name:       "Test NewDomainIntent IntentTypeHighBandwidth",
			intentType: IntentTypeHighBandwidth,
			values:     []Value{},
		},
		{
			name:       "Test NewDomainIntent IntentTypeLowBandwidth",
			intentType: IntentTypeLowBandwidth,
			values:     []Value{},
		},
		{
			name:       "Test NewDomainIntent IntentTypeLowJitter",
			intentType: IntentTypeLowJitter,
			values:     []Value{},
		},
		{
			name:       "Test NewDomainIntent IntentTypeLowLatency",
			intentType: IntentTypeLowLatency,
			values:     []Value{},
		},
		{
			name:       "Test NewDomainIntent IntentTypeLowUtilization",
			intentType: IntentTypeLowUtilization,
			values:     []Value{},
		},
		{
			name:       "Test NewDomainIntent IntentTypeSFC",
			intentType: IntentTypeSFC,
			values:     []Value{},
		},
		{
			name:       "Test NewDomainIntent IntentTypeUnspecified",
			intentType: IntentTypeUnspecified,
			values:     []Value{},
		},
	}

	for _, tt := range tests {
		intent := NewDomainIntent(tt.intentType, tt.values)
		assert.NotNil(t, intent)
	}
}

func TestDomainIntent_GetIntentType(t *testing.T) {
	tests := []struct {
		name       string
		intentType IntentType
		values     []Value
	}{
		{
			name:       "Test DomainIntent GetIntentType IntentTypeFlexAlgo",
			intentType: IntentTypeFlexAlgo,
			values:     []Value{},
		},
		{
			name:       "Test DomainIntent GetIntentType IntentTypeHighBandwidth",
			intentType: IntentTypeHighBandwidth,
			values:     []Value{},
		},
		{
			name:       "Test DomainIntent GetIntentType IntentTypeLowBandwidth",
			intentType: IntentTypeLowBandwidth,
			values:     []Value{},
		},
		{
			name:       "Test DomainIntent GetIntentType IntentTypeLowJitter",
			intentType: IntentTypeLowJitter,
			values:     []Value{},
		},
		{
			name:       "Test DomainIntent GetIntentType IntentTypeLowLatency",
			intentType: IntentTypeLowLatency,
			values:     []Value{},
		},
		{
			name:       "Test DomainIntent GetIntentType IntentTypeLowUtilization",
			intentType: IntentTypeLowUtilization,
			values:     []Value{},
		},
		{
			name:       "Test DomainIntent GetIntentType IntentTypeSFC",
			intentType: IntentTypeSFC,
			values:     []Value{},
		},
		{
			name:       "Test DomainIntent GetIntentType IntentTypeUnspecified",
			intentType: IntentTypeUnspecified,
			values:     []Value{},
		},
	}

	for _, tt := range tests {
		intent := NewDomainIntent(tt.intentType, tt.values)
		assert.Equal(t, tt.intentType, intent.GetIntentType())
	}
}

func getNumberValue(valueType ValueType, value *int32) *NumberValue {
	number, _ := NewNumberValue(valueType, value)
	return number
}

func GetStringValue(valueType ValueType, value *string) *StringValue {
	stringValue, _ := NewStringValue(valueType, value)
	return stringValue
}

func TestDomainIntent_GetValues(t *testing.T) {
	tests := []struct {
		name       string
		intentType IntentType
		values     []Value
	}{
		{
			name:       "Test DomainIntent GetValues no values",
			intentType: IntentTypeFlexAlgo,
			values:     []Value{},
		},
		{
			name:       "Test DomainIntent GetValues NumberValue",
			intentType: IntentTypeHighBandwidth,
			values: []Value{
				getNumberValue(ValueTypeMinValue, proto.Int32(1)),
			},
		},
		{
			name:       "Test DomainIntent GetValues string",
			intentType: IntentTypeHighBandwidth,
			values: []Value{
				GetStringValue(ValueTypeSFC, proto.String("fw")),
			},
		},
	}

	for _, tt := range tests {
		intent := NewDomainIntent(tt.intentType, tt.values)
		assert.Equal(t, tt.values, intent.GetValues())
	}
}

func TestDomainIntent_convertValue(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  string
	}{
		{
			name:  "Test convertValue ValueTypeMinValue",
			value: getNumberValue(ValueTypeMinValue, proto.Int32(1)),
			want:  "MinValue:1",
		},
		{
			name:  "Test convertValue ValueTypeMaxValue",
			value: getNumberValue(ValueTypeMaxValue, proto.Int32(1)),
			want:  "MaxValue:1",
		},
		{
			name:  "Test convertValue ValueTypeFlexAlgoNr",
			value: getNumberValue(ValueTypeFlexAlgoNr, proto.Int32(128)),
			want:  "FlexAlgoNr:128",
		},
		{
			name:  "Test convertValue ValueTypeSFC",
			value: GetStringValue(ValueTypeSFC, proto.String("fw")),
			want:  "SFC:fw",
		},
		{
			name:  "Test convertValue default",
			value: getNumberValue(ValueTypeUnspecified, proto.Int32(1)),
			want:  "Unspecified",
		},
	}

	for _, tt := range tests {
		intent := NewDomainIntent(IntentTypeFlexAlgo, []Value{})
		assert.Equal(t, tt.want, intent.convertValue(tt.value))
	}
}

func TestDomainIntent_Serialize(t *testing.T) {
	tests := []struct {
		name       string
		intentType IntentType
		values     []Value
		want       string
	}{
		{
			name:       "Test DomainIntent Serialize no values",
			intentType: IntentTypeLowJitter,
			values:     []Value{},
			want:       "LowJitter",
		},
		{
			name:       "Test DomainIntent Serialize NumberValue",
			intentType: IntentTypeHighBandwidth,
			values: []Value{
				getNumberValue(ValueTypeMinValue, proto.Int32(1)),
			},
			want: "HighBandwidth,MinValue:1",
		},
		{
			name:       "Test DomainIntent Serialize string",
			intentType: IntentTypeHighBandwidth,
			values: []Value{
				GetStringValue(ValueTypeSFC, proto.String("fw")),
			},
			want: "HighBandwidth,SFC:fw",
		},
	}

	for _, tt := range tests {
		intent := NewDomainIntent(tt.intentType, tt.values)
		assert.Equal(t, tt.want, intent.Serialize())
	}
}
