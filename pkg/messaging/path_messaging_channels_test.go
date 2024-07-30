package messaging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPathMessagingChannels(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewPathMessagingChannels",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, NewPathMessagingChannels())
		})
	}
}

func TestPathMessagingChannels_GetPathRequestChan(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestPathMessagingChannels_GetPathRequestChan",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channels := NewPathMessagingChannels()
			assert.NotNil(t, channels.GetPathRequestChan())
		})
	}
}

func TestPathMessagingChannels_GetPathResponseChan(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestPathMessagingChannels_GetPathResponseChan",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channels := NewPathMessagingChannels()
			assert.NotNil(t, channels.GetPathResponseChan())
		})
	}
}

func TestPathMessagingChannels_GetErrorChan(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestPathMessagingChannels_GetErrorChan",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channels := NewPathMessagingChannels()
			assert.NotNil(t, channels.GetErrorChan())
		})
	}
}
