package domain

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestNewStringValue(t *testing.T) {
	tests := []struct {
		name        string
		valueType   ValueType
		stringValue *string
		wantErr     bool
	}{
		{
			name:        "Test NewStringValue no stringValue",
			valueType:   ValueTypeUnspecified,
			stringValue: nil,
			wantErr:     true,
		},
		{
			name:        "Test NewStringValue success",
			valueType:   ValueTypeUnspecified,
			stringValue: proto.String("test"),
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStringValue(tt.valueType, tt.stringValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStringValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStringValue_GetStringValue(t *testing.T) {
	tests := []struct {
		name        string
		valueType   ValueType
		stringValue *string
		expected    string
	}{
		{
			name:        "Test StringValue GetStringValue",
			valueType:   ValueTypeUnspecified,
			stringValue: proto.String("test"),
			expected:    "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := NewStringValue(tt.valueType, tt.stringValue)
			if got := value.GetStringValue(); got != tt.expected {
				t.Errorf("StringValue.GetStringValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStringValue_GetNumberValue(t *testing.T) {
	tests := []struct {
		name        string
		valueType   ValueType
		stringValue *string
		expected    int32
	}{
		{
			name:        "Test StringValue GetNumberValue",
			valueType:   ValueTypeUnspecified,
			stringValue: proto.String("test"),
			expected:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := NewStringValue(tt.valueType, tt.stringValue)
			if got := value.GetNumberValue(); got != tt.expected {
				t.Errorf("StringValue.GetNumberValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}
