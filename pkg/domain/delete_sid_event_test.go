package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeleteSidEvent(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{
			name: "Test NewDeleteSidEvent",
			key:  "0_0000.0000.000b_fc00:0:b:0:1::",
		},
	}

	for _, tt := range tests {
		event := NewDeleteSidEvent(tt.key)
		assert.NotNil(t, event)
	}
}

func TestDeleteSidEvent_GetKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{
			name: "Test DeleteSidEvent GetKey",
			key:  "0_0000.0000.000b_fc00:0:b:0:1::",
		},
	}
	for _, tt := range tests {
		event := NewDeleteSidEvent(tt.key)
		assert.Equal(t, tt.key, event.GetKey())
	}
}
