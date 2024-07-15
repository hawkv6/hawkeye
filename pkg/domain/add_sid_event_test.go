package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestNewAddSidEvent(t *testing.T) {
	tests := []struct {
		name        string
		key         *string
		igpRouterId *string
		sid         *string
		algorithm   *uint32
	}{
		{
			name:        "Test NewAddSidEvent",
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
		event := NewAddSidEvent(sid)
		assert.NotNil(t, event)
	}
}

func TestAddSidEvent_GetKey(t *testing.T) {
	tests := []struct {
		name        string
		key         *string
		igpRouterId *string
		sid         *string
		algorithm   *uint32
	}{
		{
			name:        "Test AddSidEvent GetKey",
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
		event := NewAddSidEvent(sid)
		assert.Equal(t, *tt.key, event.GetKey())
	}
}
