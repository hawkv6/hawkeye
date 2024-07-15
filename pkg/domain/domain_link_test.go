package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestNewDomainLink(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           *DomainLink
		wantErr                        bool
	}{
		{
			name:                           "Test NewDomainLink success",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want: &DomainLink{
				key:                            "2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6",
				igpRouterId:                    "0000.0000.000b",
				remoteIgpRouterId:              "0000.0000.0006",
				igpMetric:                      10,
				unidirLinkDelay:                2000,
				unidirDelayVariation:           100,
				maxLinkBWKbps:                  1000000,
				unidirAvailableBandwidth:       99766,
				unidirBandwidthUtilization:     234,
				unidirPacketLoss:               3.0059316283477027,
				normalizedUnidirLinkDelay:      0.05,
				normalizedUnidirDelayVariation: 0.016452169298129225,
				normalizedUnidirPacketLoss:     1e-10,
			},
			wantErr: false,
		},
		{
			name:                           "Test NewDomainLink validation Error:Field validation for 'MaxLinkBWKbps' failed on the 'min' tag",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(0),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           nil,
			wantErr:                        true,
		},
		{
			name:                           "Test NewDomainLink validation Error:Field validation for 'UnidirPacketLoss' failed on the 'min' tag",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(-0.5),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           nil,
			wantErr:                        true,
		},
		{
			name:                           "Test NewDomainLink validation Error:Field validation for 'UnidirPacketLoss' failed on the 'max' tag",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(100.5),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           nil,
			wantErr:                        true,
		},
		{
			name:                           "Test NewDomainLink validation Error:Field validation for 'NormalizedUnidirLinkDelay' failed on the 'min' tag",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(2.5),
			normalizedUnidirLinkDelay:      proto.Float64(-0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           nil,
			wantErr:                        true,
		},
		{
			name:                           "Test NewDomainLink validation Error:Field validation for 'NormalizedUnidirLinkDelay' failed on the 'max' tag",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(2.5),
			normalizedUnidirLinkDelay:      proto.Float64(1.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           nil,
			wantErr:                        true,
		},
		{
			name:                           "Test NewDomainLink validation Error:Field validation for 'NormalizedUnidirLinkDelayVariation' failed on the 'min' tag",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(2.5),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(-0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           nil,
			wantErr:                        true,
		},
		{
			name:                           "Test NewDomainLink validation Error:Field validation for 'NormalizedUnidirLinkDelayVariation' failed on the 'max' tag",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(2.5),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(1.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           nil,
			wantErr:                        true,
		},
		{
			name:                           "Test NewDomainLink validation Error:Field validation for 'NormalizedUnidirPacketLoss' failed on the 'min' tag",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(2.5),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(-0.5),
			want:                           nil,
			wantErr:                        true,
		},
		{
			name:                           "Test NewDomainLink validation Error:Field validation for 'NormalizedUnidirPacketLoss' failed on the 'max' tag",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(2.5),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1.5),
			want:                           nil,
			wantErr:                        true,
		},
	}

	for _, tt := range tests {
		link, err := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		if (err != nil) != tt.wantErr {
			t.Errorf("NewDomainLink() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if link == nil && tt.want != nil {
			t.Errorf("NewDomainLink() got nil, want %v", tt.want)
		}

	}
}

func TestDomainLink_GetKey(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           string
	}{
		{
			name:                           "Test DomainLink GetKey",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           "2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6",
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		assert.Equal(t, tt.want, link.GetKey())
	}
}

func TestDomainLink_GetIgpRouterId(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           string
	}{
		{
			name:                           "Test DomainLink GetIgpRouterId",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           "0000.0000.000b",
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		assert.Equal(t, tt.want, link.GetIgpRouterId())
	}
}

func TestDomainLink_GetRemoteIgpRouterId(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           string
	}{
		{
			name:                           "Test DomainLink GetRemoteIgpRouterId",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           "0000.0000.0006",
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		assert.Equal(t, tt.want, link.GetRemoteIgpRouterId())
	}
}

func TestDomainLink_GetIgpMetric(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           uint32
	}{
		{
			name:                           "Test DomainLink GetIgpMetric",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           10,
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		assert.Equal(t, tt.want, link.GetIgpMetric())
	}
}

func TestDomainLink_GetUndirLinkDelay(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           uint32
	}{
		{
			name:                           "Test DomainLink GetUnidirLinkDelay",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           2000,
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		assert.Equal(t, tt.want, link.GetUnidirLinkDelay())
	}
}

func TestDomainLink_GetUndirDelayVariation(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           uint32
	}{
		{
			name:                           "Test DomainLink GetUnidirDelayVariation",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           100,
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		assert.Equal(t, tt.want, link.GetUnidirDelayVariation())
	}
}

func TestDomainLink_GetMaxLinkBWKbps(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           uint64
	}{
		{
			name:                           "Test DomainLink GetMaxLinkBWKbps",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           1000000,
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		assert.Equal(t, tt.want, link.GetMaxLinkBWKbps())
	}
}

func TestDomainLink_GetUnidirAvailableBandwidth(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           uint32
	}{
		{
			name:                           "Test DomainLink GetUnidirAvailableBandwidth",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           99766,
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		assert.Equal(t, tt.want, link.GetUnidirAvailableBandwidth())
	}
}

func TestDomainLink_GetUnidirPacketLoss(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           float64
	}{
		{
			name:                           "Test DomainLink GetUnidirPacketLoss",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           3.0059316283477027,
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		assert.Equal(t, tt.want, link.GetUnidirPacketLoss())
	}
}

func TestDomainLink_GetUnidirBandwidthUtilization(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           uint32
	}{
		{
			name:                           "Test DomainLink GetUnidirBandwidthUtilization",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           234,
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		assert.Equal(t, tt.want, link.GetUnidirBandwidthUtilization())
	}
}

func TestDomainLink_GetNormalizedUnidirLinkDelay(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           float64
	}{
		{
			name:                           "Test DomainLink GetUnidirBandwidthUtilization",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           0.05,
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		assert.Equal(t, tt.want, link.GetNormalizedUnidirLinkDelay())
	}
}

func TestDomainLink_GetNormalizedUnidirDelayVariation(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           float64
	}{
		{
			name:                           "Test DomainLink GetUnidirBandwidthUtilization",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           0.016452169298129225,
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		assert.Equal(t, tt.want, link.GetNormalizedUnidirDelayVariation())
	}
}

func TestDomainLink_GetNormalizedUnidirPacketLoss(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           float64
	}{
		{
			name:                           "Test DomainLink GetUnidirBandwidthUtilization",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(0.01),
			want:                           0.01,
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		assert.Equal(t, tt.want, link.GetNormalizedUnidirPacketLoss())
	}
}

func TestDomainLink_SetNormalizedUnidirLinkDelay(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           float64
	}{
		{
			name:                           "Test DomainLink GetUnidirBandwidthUtilization",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           0.1,
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		link.SetNormalizedUnidirLinkDelay(tt.want)
		assert.Equal(t, tt.want, link.GetNormalizedUnidirLinkDelay())
	}
}

func TestDomainLink_SetNormalizedUnidirDelayVariation(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           float64
	}{
		{
			name:                           "Test DomainLink GetUnidirBandwidthUtilization",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(1e-10),
			want:                           0.1,
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		link.SetNormalizedUnidirDelayVariation(tt.want)
		assert.Equal(t, tt.want, link.GetNormalizedUnidirDelayVariation())
	}
}

func TestDomainLink_SetNormalizedUnidirPacketLoss(t *testing.T) {
	tests := []struct {
		name                           string
		key                            *string
		igpRouterId                    *string
		remoteIgpRouterId              *string
		igpMetric                      *uint32
		unidirLinkDelay                *uint32
		unidirDelayVariation           *uint32
		maxLinkBWKbps                  *uint64
		unidirAvailableBandwidth       *uint32
		unidirBandwidthUtilization     *uint32
		unidirPacketLoss               *float64
		normalizedUnidirLinkDelay      *float64
		normalizedUnidirDelayVariation *float64
		normalizedUnidirPacketLoss     *float64
		want                           float64
	}{
		{
			name:                           "Test DomainLink GetUnidirBandwidthUtilization",
			key:                            proto.String("2_0_2_0_0000.0000.000b_2001:db8:b6::b_0000.0000.0006_2001:db8:b6::6"),
			igpRouterId:                    proto.String("0000.0000.000b"),
			remoteIgpRouterId:              proto.String("0000.0000.0006"),
			igpMetric:                      proto.Uint32(10),
			unidirLinkDelay:                proto.Uint32(2000),
			unidirDelayVariation:           proto.Uint32(100),
			maxLinkBWKbps:                  proto.Uint64(1000000),
			unidirAvailableBandwidth:       proto.Uint32(99766),
			unidirBandwidthUtilization:     proto.Uint32(234),
			unidirPacketLoss:               proto.Float64(3.0059316283477027),
			normalizedUnidirLinkDelay:      proto.Float64(0.05),
			normalizedUnidirDelayVariation: proto.Float64(0.016452169298129225),
			normalizedUnidirPacketLoss:     proto.Float64(0.01),
			want:                           0.5,
		},
	}

	for _, tt := range tests {
		link, _ := NewDomainLink(
			tt.key,
			tt.igpRouterId,
			tt.remoteIgpRouterId,
			tt.igpMetric,
			tt.unidirLinkDelay,
			tt.unidirDelayVariation,
			tt.maxLinkBWKbps,
			tt.unidirAvailableBandwidth,
			tt.unidirBandwidthUtilization,
			tt.unidirPacketLoss,
			tt.normalizedUnidirLinkDelay,
			tt.normalizedUnidirDelayVariation,
			tt.normalizedUnidirPacketLoss,
		)
		link.SetNormalizedPacketLoss(tt.want)
		assert.Equal(t, tt.want, link.GetNormalizedUnidirPacketLoss())
	}
}
