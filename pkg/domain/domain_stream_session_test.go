package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestNewDomainStreamSession(t *testing.T) {
	tests := []struct {
		name         string
		pathRequest  PathRequest
		pathResponse PathResult
		wantErr      bool
	}{
		{
			name:         "Test NewDomainStreamSession success",
			pathRequest:  NewMockPathRequest(gomock.NewController(t)),
			pathResponse: NewMockPathResult(gomock.NewController(t)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			streamSession := NewDomainStreamSession(tt.pathRequest, tt.pathResponse)
			assert.NotNil(t, streamSession)
		})
	}
}

func TestDomainStreamSession_GetContext(t *testing.T) {
	tests := []struct {
		name         string
		pathResponse PathResult
	}{
		{
			name:         "Test DomainStreamSession GetContext",
			pathResponse: NewMockPathResult(gomock.NewController(t)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			pathRequest := NewMockPathRequest(ctrl)
			pathRequest.EXPECT().GetContext().Return(context.Background())
			streamSession := NewDomainStreamSession(pathRequest, tt.pathResponse)
			if got := streamSession.GetContext(); got != context.Background() {
				t.Errorf("DomainStreamSession.GetContext() = %v, want %v", got, context.Background())
			}
		})
	}
}

func TestDomainStreamSession_GetPathRequest(t *testing.T) {
	tests := []struct {
		name         string
		pathResponse PathResult
		pathRequest  PathRequest
	}{
		{
			name:         "Test DomainStreamSession GetPathRequest",
			pathRequest:  NewMockPathRequest(gomock.NewController(t)),
			pathResponse: NewMockPathResult(gomock.NewController(t)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			streamSession := NewDomainStreamSession(tt.pathRequest, tt.pathResponse)
			assert.NotNil(t, streamSession.GetPathRequest())
		})
	}
}

func TestDomainStreamSession_GetPathResult(t *testing.T) {
	tests := []struct {
		name         string
		pathResponse PathResult
		pathRequest  PathRequest
	}{
		{
			name:         "Test DomainStreamSession GetPathResult",
			pathRequest:  NewMockPathRequest(gomock.NewController(t)),
			pathResponse: NewMockPathResult(gomock.NewController(t)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			streamSession := NewDomainStreamSession(tt.pathRequest, tt.pathResponse)
			assert.NotNil(t, streamSession.GetPathResult())
		})
	}
}

func TestDomainStreamSession_SetPathResult(t *testing.T) {
	tests := []struct {
		name         string
		pathResponse PathResult
		pathRequest  PathRequest
	}{
		{
			name:         "Test DomainStreamSession SetPathResult",
			pathRequest:  NewMockPathRequest(gomock.NewController(t)),
			pathResponse: NewMockPathResult(gomock.NewController(t)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			streamSession := NewDomainStreamSession(tt.pathRequest, tt.pathResponse)
			streamSession.pathResult = nil
			streamSession.SetPathResult(tt.pathResponse)
			assert.NotNil(t, streamSession.GetPathResult())
		})
	}
}
