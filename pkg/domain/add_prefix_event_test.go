package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestNewAddPrefixEvent(t *testing.T) {
	tests := []struct {
		name         string
		key          *string
		igpRouterId  *string
		prefix       *string
		prefixLength *int32
	}{
		{
			name:         "Test NewAddPrefixEvent",
			key:          proto.String("2_0_2_0_0_fc00:0:c:129::_64_0000.0000.000c"),
			igpRouterId:  proto.String("0000.0000.000c"),
			prefix:       proto.String("fc00:0:c:129::"),
			prefixLength: proto.Int32(64),
		},
	}

	for _, tt := range tests {
		prefix, err := NewDomainPrefix(
			tt.key,
			tt.igpRouterId,
			tt.prefix,
			tt.prefixLength,
		)
		if err != nil {
			t.Errorf("Error creating DomainPrefix: %v", err)
		}
		event := NewAddPrefixEvent(prefix)
		assert.NotNil(t, event)
	}
}

func TestAddPrefixEvent_GetKey(t *testing.T) {
	tests := []struct {
		name         string
		key          *string
		igpRouterId  *string
		prefix       *string
		prefixLength *int32
	}{
		{
			name:         "Test AddPrefixEvent GetKey",
			key:          proto.String("2_0_2_0_0_fc00:0:c:129::_64_0000.0000.000c"),
			igpRouterId:  proto.String("0000.0000.000c"),
			prefix:       proto.String("fc00:0:c:129::"),
			prefixLength: proto.Int32(64),
		},
	}

	for _, tt := range tests {
		prefix, err := NewDomainPrefix(
			tt.key,
			tt.igpRouterId,
			tt.prefix,
			tt.prefixLength,
		)
		if err != nil {
			t.Errorf("Error creating DomainPrefix: %v", err)
		}
		event := NewAddPrefixEvent(prefix)
		assert.Equal(t, *tt.key, event.GetKey())
	}
}
