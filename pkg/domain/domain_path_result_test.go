package domain

import (
	"reflect"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/graph"
	"go.uber.org/mock/gomock"
)

func TestNewDomainPathResult(t *testing.T) {
	tests := []struct {
		name             string
		pathRequest      PathRequest
		shortestPath     graph.Path
		ipv6SidAddresses []string
		wantErr          bool
	}{
		{
			name:             "New Domain Path Result",
			pathRequest:      NewMockPathRequest(gomock.NewController(t)),
			shortestPath:     graph.NewMockPath(gomock.NewController(t)),
			ipv6SidAddresses: []string{"2001:db8:0:1::1"},
			wantErr:          false,
		},
		{
			name:             "New Domain Path Result",
			pathRequest:      NewMockPathRequest(gomock.NewController(t)),
			shortestPath:     graph.NewMockPath(gomock.NewController(t)),
			ipv6SidAddresses: []string{"2001:db8:0:1::1", "2001:db8:0:1::2"},
			wantErr:          false,
		},
		{
			name:             "New Domain Path Result invalid sid addresses",
			pathRequest:      NewMockPathRequest(gomock.NewController(t)),
			shortestPath:     graph.NewMockPath(gomock.NewController(t)),
			ipv6SidAddresses: []string{"2001:db8:0:1::1", "not an ipv6"},
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDomainPathResult(tt.pathRequest, tt.shortestPath, tt.ipv6SidAddresses)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDomainPathResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDomainPathResult_GetIpv6SidAddresses(t *testing.T) {
	tests := []struct {
		name             string
		pathRequest      PathRequest
		shortestPath     graph.Path
		ipv6SidAddresses []string
	}{
		{
			name:             "Get IPv6 SID Addresses",
			pathRequest:      NewMockPathRequest(gomock.NewController(t)),
			shortestPath:     graph.NewMockPath(gomock.NewController(t)),
			ipv6SidAddresses: []string{"2001:db8:0:1::1", "2001:db8:0:1::2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathResult, err := NewDomainPathResult(tt.pathRequest, tt.shortestPath, tt.ipv6SidAddresses)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(tt.ipv6SidAddresses, pathResult.GetIpv6SidAddresses()) {
				t.Errorf("Expected %v, got %v", tt.ipv6SidAddresses, pathResult.GetIpv6SidAddresses())
			}
		})
	}
}

func TestDomainPathResult_GetServiceSidList(t *testing.T) {
	tests := []struct {
		name             string
		pathRequest      PathRequest
		shortestPath     graph.Path
		ipv6SidAddresses []string
	}{
		{
			name:             "Get Service SID List",
			pathRequest:      NewMockPathRequest(gomock.NewController(t)),
			shortestPath:     graph.NewMockPath(gomock.NewController(t)),
			ipv6SidAddresses: []string{"2001:db8:0:1::1", "2001:db8:0:1::2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathResult, err := NewDomainPathResult(tt.pathRequest, tt.shortestPath, tt.ipv6SidAddresses)
			if err != nil {
				t.Error(err)
			}
			pathResult.serviceSidAddresses = tt.ipv6SidAddresses
			if !reflect.DeepEqual(tt.ipv6SidAddresses, pathResult.GetServiceSidList()) {
				t.Errorf("Expected %v, got %v", tt.ipv6SidAddresses, pathResult.GetServiceSidList())
			}
		})
	}
}

func TestDomainPathResult_SetServiceSidList(t *testing.T) {
	tests := []struct {
		name             string
		pathRequest      PathRequest
		shortestPath     graph.Path
		ipv6SidAddresses []string
	}{
		{
			name:             "Set Service SID List",
			pathRequest:      NewMockPathRequest(gomock.NewController(t)),
			shortestPath:     graph.NewMockPath(gomock.NewController(t)),
			ipv6SidAddresses: []string{"2001:db8:0:1::1", "2001:db8:0:1::2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathResult, err := NewDomainPathResult(tt.pathRequest, tt.shortestPath, tt.ipv6SidAddresses)
			if err != nil {
				t.Error(err)
			}
			pathResult.SetServiceSidList(tt.ipv6SidAddresses)
			if !reflect.DeepEqual(tt.ipv6SidAddresses, pathResult.GetServiceSidList()) {
				t.Errorf("Expected %v, got %v", tt.ipv6SidAddresses, pathResult.serviceSidAddresses)
			}
		})
	}
}
