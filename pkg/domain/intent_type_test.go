package domain

import (
	"testing"
)

func TestIntentType_String(t *testing.T) {
	tests := []struct {
		name     string
		intent   IntentType
		expected string
	}{
		{"Unspecified", IntentTypeUnspecified, "Unspecified"},
		{"HighBandwidth", IntentTypeHighBandwidth, "HighBandwidth"},
		{"LowBandwidth", IntentTypeLowBandwidth, "LowBandwidth"},
		{"LowLatency", IntentTypeLowLatency, "LowLatency"},
		{"LowPacketLoss", IntentTypeLowPacketLoss, "LowPacketLoss"},
		{"LowJitter", IntentTypeLowJitter, "LowJitter"},
		{"FlexAlgo", IntentTypeFlexAlgo, "FlexAlgo"},
		{"SFC", IntentTypeSFC, "SFC"},
		{"LowUtilization", IntentTypeLowUtilization, "LowUtilization"},
		{"Unknown", IntentType(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.intent.String(); got != tt.expected {
				t.Errorf("IntentType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
