package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestNewUpdateLinkEvent(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirLinkDelayVariation       *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
	}{
		{
			name:                           "Test NewUpdateLinkEvent",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirLinkDelayVariation:       proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link, err := NewDomainLink(
				tt.key,
				tt.igpRouterId,
				tt.remoteIgpRouterId,
				tt.igpMetric,
				tt.unidirLinkDelay,
				tt.unidirLinkDelayVariation,
				tt.maxLinkBWKbps,
				tt.unidirAvailableBandwidth,
				tt.unidirBandwidthUtilization,
				tt.unidirPacketLoss,
				tt.normalizedUnidirLinkDelay,
				tt.normalizedUnidirDelayVariation,
				tt.normalizedUnidirPacketLoss,
			)
			if err != nil {
				t.Error(err)
			}
			event := NewUpdateLinkEvent(link)
			assert.NotNil(t, event)
		})
	}
}

func TestUpdateLinkEvent_GetKey(t *testing.T) {
	tests := []struct {
		name string
		link Link
		want string
	}{
		{
			name: "Test UpdateLinkEvent GetKey",
			want: "2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link, err := NewDomainLink(
				proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
				proto.String("0000.0000.000b"),
				proto.String("0000.0000.0006"),
				proto.Uint32(10),
				proto.Uint32(2000),
				proto.Uint32(100),
				proto.Uint64(1000000),
				proto.Uint32(99766),
				proto.Uint32(234),
				proto.Float64(3.0059316283477027),
				proto.Float64(0.05),
				proto.Float64(0.016452169298129225),
				proto.Float64(1e-10),
			)
			if err != nil {
				t.Error(err)
			}
			event := NewUpdateLinkEvent(link)
			assert.Equal(t, tt.want, event.GetKey())
		})
	}
}
