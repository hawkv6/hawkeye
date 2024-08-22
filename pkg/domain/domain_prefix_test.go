package domain

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestNewDomainPrefix(t *testing.T) {
	tests := []struct {
		name         string
		key          *string
		igpRouterId  *string
		prefix       *string
		prefixLength *int32
		wantErr      bool
	}{
		{
			name:         "Test NewDomainPrefix success",
			key:          proto.String("2_0_2_0_0_fc00:0:c:129::_64_0000.0000.000c"),
			igpRouterId:  proto.String("0000.0000.000c"),
			prefix:       proto.String("fc00:0:c:129::"),
			prefixLength: proto.Int32(64),
			wantErr:      false,
		},
		{
			name:         "Test NewDomainPrefix failed with invalid prefix length",
			key:          proto.String("2_0_2_0_0_fc00:0:c:129::_64_0000.0000.000c"),
			igpRouterId:  proto.String("0000.0000.000c"),
			prefix:       proto.String("fc00:0:c:129::"),
			prefixLength: proto.Int32(129),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDomainPrefix(tt.key, tt.igpRouterId, tt.prefix, tt.prefixLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDomainPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDomainPrefix_GetKey(t *testing.T) {
	tests := []struct {
		name         string
		key          *string
		igpRouterId  *string
		prefix       *string
		prefixLength *int32
	}{
		{
			name:         "Test DomainPrefix GetKey",
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
		if got := prefix.GetKey(); got != *tt.key {
			t.Errorf("DomainPrefix.GetKey() = %v, want %v", got, *tt.key)
		}
	}
}

func TestDomainPrefix_GetIgpRouterId(t *testing.T) {
	tests := []struct {
		name         string
		key          *string
		igpRouterId  *string
		prefix       *string
		prefixLength *int32
	}{
		{
			name:         "Test DomainPrefix GetIgpRouterId",
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
		if got := prefix.GetIgpRouterId(); got != *tt.igpRouterId {
			t.Errorf("DomainPrefix.GetIgpRouterId() = %v, want %v", got, *tt.igpRouterId)
		}
	}
}

func TestDomainPrefix_GetPrefix(t *testing.T) {
	tests := []struct {
		name         string
		key          *string
		igpRouterId  *string
		prefix       *string
		prefixLength *int32
	}{
		{
			name:         "Test DomainPrefix GetPrefix",
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
		if got := prefix.GetPrefix(); got != *tt.prefix {
			t.Errorf("DomainPrefix.GetPrefix() = %v, want %v", got, *tt.prefix)
		}
	}
}

func TestDomainPrefix_GetPrefixLength(t *testing.T) {
	tests := []struct {
		name         string
		key          *string
		igpRouterId  *string
		prefix       *string
		prefixLength *int32
	}{
		{
			name:         "Test DomainPrefix GetPrefixLength",
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
		if got := prefix.GetPrefixLength(); got != uint8(*tt.prefixLength) {
			t.Errorf("DomainPrefix.GetPrefixLength() = %v, want %v", got, uint8(*tt.prefixLength))
		}
	}
}
