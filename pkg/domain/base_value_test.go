package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAddBaseValue(t *testing.T) {
	tests := []struct {
		name      string
		valueType ValueType
	}{
		{
			name:      "Test NewAddBaseValue ValueTypeUnspecified",
			valueType: ValueTypeUnspecified,
		},
		{
			name:      "Test NewAddBaseValue ValueTypeFlexAlgoNr",
			valueType: ValueTypeFlexAlgoNr,
		},
		{
			name:      "Test NewAddBaseValue ValueTypeMaxValue",
			valueType: ValueTypeMaxValue,
		},
		{
			name:      "Test NewAddBaseValue ValueTypeMinValue",
			valueType: ValueTypeMinValue,
		},
		{
			name:      "Test NewAddBaseValue ValueTypeSFC",
			valueType: ValueTypeSFC,
		},
	}

	for _, tt := range tests {
		baseValue := NewBaseValue(tt.valueType)
		assert.NotNil(t, baseValue)
	}
}

func TestBaseValue_GetValueType(t *testing.T) {
	tests := []struct {
		name      string
		valueType ValueType
	}{
		{
			name:      "Test BaseValue GetValueType ValueTypeUnspecified",
			valueType: ValueTypeUnspecified,
		},
		{
			name:      "Test BaseValue GetValueType ValueTypeFlexAlgoNr",
			valueType: ValueTypeFlexAlgoNr,
		},
		{
			name:      "Test BaseValue GetValueType ValueTypeMaxValue",
			valueType: ValueTypeMaxValue,
		},
		{
			name:      "Test BaseValue GetValueType ValueTypeMinValue",
			valueType: ValueTypeMinValue,
		},
		{
			name:      "Test BaseValue GetValueType ValueTypeSFC",
			valueType: ValueTypeSFC,
		},
	}

	for _, tt := range tests {
		baseValue := NewBaseValue(tt.valueType)
		assert.NotNil(t, baseValue)
		assert.Equal(t, tt.valueType, baseValue.GetValueType())
	}
}
