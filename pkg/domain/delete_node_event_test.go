package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeleteNodeEvent(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{
			name: "Test NewDeleteNodeEvent",
			key:  "2_0_0_0000.0000.0004",
		},
	}

	for _, tt := range tests {
		event := NewDeleteNodeEvent(tt.key)
		assert.NotNil(t, event)
	}
}

func TestDeleteNodeEvent_GetKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{
			name: "Test DeleteNodeEvent GetKey",
			key:  "2_0_0_0000.0000.0004",
		},
	}
	for _, tt := range tests {
		event := NewDeleteNodeEvent(tt.key)
		assert.Equal(t, tt.key, event.GetKey())
	}
}
