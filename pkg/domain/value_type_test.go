package domain

import "testing"

func TestValeType_String(t *testing.T) {
	tests := []struct {
		name     string
		value    ValueType
		expected string
	}{
		{"Unspecified", ValueTypeUnspecified, "Unspecified"},
		{"MinValue", ValueTypeMinValue, "MinValue"},
		{"MaxValue", ValueTypeMaxValue, "MaxValue"},
		{"SFC", ValueTypeSFC, "SFC"},
		{"FlexAlgoNr", ValueTypeFlexAlgoNr, "FlexAlgoNr"},
		{"Unknown", ValueType(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.expected {
				t.Errorf("ValueType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
