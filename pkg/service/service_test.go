package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConcreteService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewConcreteService",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, NewConcreteService("servicetype", "serviceId", "sid", true))
		})
	}
}

func TestConcreteService_GetType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestConcreteService_GetType",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewConcreteService("servicetype", "serviceId", "sid", true)
			assert.Equal(t, "servicetype", service.GetType())
		})
	}
}

func TestConcreteService_GetId(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestConcreteService_GetId",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewConcreteService("servicetype", "serviceId", "sid", true)
			assert.Equal(t, "serviceId", service.GetId())
		})
	}
}

func TestConcreteService_GetSid(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestConcreteService_GetSid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewConcreteService("servicetype", "serviceId", "sid", true)
			assert.Equal(t, "sid", service.GetSid())
		})
	}
}

func TestConcreteService_IsHealthy(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestConcreteService_IsHealthy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewConcreteService("servicetype", "serviceId", "sid", true)
			assert.True(t, service.IsHealty())
		})
	}
}
