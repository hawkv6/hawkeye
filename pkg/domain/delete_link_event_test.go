package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeleteLinkEvent(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{
			name: "Test NewDeleteLinkEvent",
			key:  "2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6",
		},
	}

	for _, tt := range tests {
		event := NewDeleteLinkEvent(tt.key)
		assert.NotNil(t, event)
	}
}

func TestDeleteLinkEvent_GetKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{
			name: "Test DeleteLinkEvent GetKey",
			key:  "2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6",
		},
	}
	for _, tt := range tests {
		event := NewDeleteLinkEvent(tt.key)
		assert.Equal(t, tt.key, event.GetKey())
	}
}
