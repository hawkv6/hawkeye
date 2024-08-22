package domain

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestNewNumberValue(t *testing.T) {
	tests := []struct {
		name        string
		valueType   ValueType
		numberValue *int32
		wantErr     bool
	}{
		{
			name:        "Test NewNumberValue no numberValue",
			valueType:   ValueTypeUnspecified,
			numberValue: nil,
			wantErr:     true,
		},
		{
			name:        "Test NewNumberValue type min success",
			valueType:   ValueTypeMinValue,
			numberValue: proto.Int32(10),
			wantErr:     false,
		},
		{
			name:        "Test NewNumberValue type max success",
			valueType:   ValueTypeMaxValue,
			numberValue: proto.Int32(10),
			wantErr:     false,
		},
		{
			name:        "Test NewNumberValue type flex algo success",
			valueType:   ValueTypeFlexAlgoNr,
			numberValue: proto.Int32(128),
			wantErr:     false,
		},
		{
			name:        "Test NewNumberValue type flex algo negative number",
			valueType:   ValueTypeFlexAlgoNr,
			numberValue: proto.Int32(-1),
			wantErr:     true,
		},
		{
			name:        "Test NewNumberValue type max negative number",
			valueType:   ValueTypeMaxValue,
			numberValue: proto.Int32(-1),
			wantErr:     true,
		},
		{
			name:        "Test NewNumberValue type min negative number",
			valueType:   ValueTypeMinValue,
			numberValue: proto.Int32(-1),
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewNumberValue(tt.valueType, tt.numberValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNumberValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNumberValue_GetNumberValue(t *testing.T) {
	tests := []struct {
		name        string
		valueType   ValueType
		numberValue *int32
		expected    int32
	}{
		{
			name:        "Test NumberValue GetNumberValue type min",
			valueType:   ValueTypeMinValue,
			numberValue: proto.Int32(10),
			expected:    10,
		},
		{
			name:        "Test NumberValue GetNumberValue type max",
			valueType:   ValueTypeMaxValue,
			numberValue: proto.Int32(10),
			expected:    10,
		},
		{
			name:        "Test NumberValue GetNumberValue type flex algo",
			valueType:   ValueTypeFlexAlgoNr,
			numberValue: proto.Int32(128),
			expected:    128,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			numberValue, _ := NewNumberValue(tt.valueType, tt.numberValue)
			if got := numberValue.GetNumberValue(); got != tt.expected {
				t.Errorf("NumberValue.GetNumberValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNumberValue_GetStringValue(t *testing.T) {
	tests := []struct {
		name        string
		valueType   ValueType
		numberValue *int32
		expected    string
	}{
		{
			name:        "Test NumberValue GetStringValue type min",
			valueType:   ValueTypeMinValue,
			numberValue: proto.Int32(10),
			expected:    "",
		},
		{
			name:        "Test NumberValue GetStringValue type max",
			valueType:   ValueTypeMaxValue,
			numberValue: proto.Int32(10),
			expected:    "",
		},
		{
			name:        "Test NumberValue GetStringValue type flex algo",
			valueType:   ValueTypeFlexAlgoNr,
			numberValue: proto.Int32(128),
			expected:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			numberValue, _ := NewNumberValue(tt.valueType, tt.numberValue)
			if got := numberValue.GetStringValue(); got != tt.expected {
				t.Errorf("NumberValue.GetStringValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}
