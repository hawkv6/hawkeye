package domain

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestNewDomainSid(t *testing.T) {
	tests := []struct {
		name        string
		key         *string
		igpRouterId *string
		sid         *string
		algorithm   *uint32
		wantErr     bool
	}{
		{
			name:        "Test NewDomainSid success",
			key:         proto.String("0_0000.0000.000b_fc00:0:b:0:1::"),
			igpRouterId: proto.String("0000.0000.000b"),
			sid:         proto.String("fc00:0:b:0:1::"),
			algorithm:   proto.Uint32(0),
			wantErr:     false,
		},
		{
			name:        "Test NewDomainSid failed with invalid algorithm",
			key:         proto.String("0_0000.0000.000b_fc00:0:b:0:1::"),
			igpRouterId: proto.String("0000.0000.000b"),
			sid:         proto.String("fc00:0:b:0:1::"),
			algorithm:   proto.Uint32(256),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDomainSid(tt.key, tt.igpRouterId, tt.sid, tt.algorithm)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDomainSid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDomainSid_GetKey(t *testing.T) {
	tests := []struct {
		name        string
		key         *string
		igpRouterId *string
		sid         *string
		algorithm   *uint32
	}{
		{
			name:        "Test DomainSid GetKey",
			key:         proto.String("0_0000.0000.000b_fc00:0:b:0:1::"),
			igpRouterId: proto.String("0000.0000.000b"),
			sid:         proto.String("fc00:0:b:0:1::"),
			algorithm:   proto.Uint32(0),
		},
	}

	for _, tt := range tests {
		sid, err := NewDomainSid(
			tt.key,
			tt.igpRouterId,
			tt.sid,
			tt.algorithm,
		)
		if err != nil {
			t.Errorf("Error creating DomainSid: %v", err)
		}
		if sid.GetKey() != *tt.key {
			t.Errorf("Expected %v, got %v", *tt.key, sid.GetKey())
		}
	}
}

func TestDomainSid_GetIgpRouterId(t *testing.T) {
	tests := []struct {
		name        string
		key         *string
		igpRouterId *string
		sid         *string
		algorithm   *uint32
	}{
		{
			name:        "Test DomainSid GetIgpRouterId",
			key:         proto.String("0_0000.0000.000b_fc00:0:b:0:1::"),
			igpRouterId: proto.String("0000.0000.000b"),
			sid:         proto.String("fc00:0:b:0:1::"),
			algorithm:   proto.Uint32(0),
		},
	}

	for _, tt := range tests {
		sid, err := NewDomainSid(
			tt.key,
			tt.igpRouterId,
			tt.sid,
			tt.algorithm,
		)
		if err != nil {
			t.Errorf("Error creating DomainSid: %v", err)
		}
		if sid.GetIgpRouterId() != *tt.igpRouterId {
			t.Errorf("Expected %v, got %v", *tt.igpRouterId, sid.GetIgpRouterId())
		}
	}
}

func TestDomainSid_GetSid(t *testing.T) {
	tests := []struct {
		name        string
		key         *string
		igpRouterId *string
		sid         *string
		algorithm   *uint32
	}{
		{
			name:        "Test DomainSid GetSid",
			key:         proto.String("0_0000.0000.000b_fc00:0:b:0:1::"),
			igpRouterId: proto.String("0000.0000.000b"),
			sid:         proto.String("fc00:0:b:0:1::"),
			algorithm:   proto.Uint32(0),
		},
	}

	for _, tt := range tests {
		sid, err := NewDomainSid(
			tt.key,
			tt.igpRouterId,
			tt.sid,
			tt.algorithm,
		)
		if err != nil {
			t.Errorf("Error creating DomainSid: %v", err)
		}
		if sid.GetSid() != *tt.sid {
			t.Errorf("Expected %v, got %v", *tt.sid, sid.GetSid())
		}
	}
}

func TestDomainSid_GetAlgorithm(t *testing.T) {
	tests := []struct {
		name        string
		key         *string
		igpRouterId *string
		sid         *string
		algorithm   *uint32
	}{
		{
			name:        "Test DomainSid GetAlgorithm",
			key:         proto.String("0_0000.0000.000b_fc00:0:b:0:1::"),
			igpRouterId: proto.String("0000.0000.000b"),
			sid:         proto.String("fc00:0:b:0:1::"),
			algorithm:   proto.Uint32(0),
		},
	}

	for _, tt := range tests {
		sid, err := NewDomainSid(
			tt.key,
			tt.igpRouterId,
			tt.sid,
			tt.algorithm,
		)
		if err != nil {
			t.Errorf("Error creating DomainSid: %v", err)
		}
		if sid.GetAlgorithm() != *tt.algorithm {
			t.Errorf("Expected %v, got %v", *tt.algorithm, sid.GetAlgorithm())
		}
	}
}
