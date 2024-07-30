package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeletePrefixEvent(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{
			name: "Test NewDeletePrefixEvent",
			key:  "2_0_2_0_0_fc00:0:c:129::_64_0000.0000.000c",
		},
	}

	for _, tt := range tests {
		event := NewDeletePrefixEvent(tt.key)
		assert.NotNil(t, event)
	}
}

func TestDeletePrefixEvent_GetKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{
			name: "Test DeletePrefixEvent GetKey",
			key:  "2_0_2_0_0_fc00:0:c:129::_64_0000.0000.000c",
		},
	}
	for _, tt := range tests {
		event := NewDeletePrefixEvent(tt.key)
		assert.Equal(t, tt.key, event.GetKey())
	}
}
