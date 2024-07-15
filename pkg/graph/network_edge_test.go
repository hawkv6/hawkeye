package graph

import (
	"testing"

	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/stretchr/testify/assert"
)

func TestNewNetworkEdge(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		from    Node
		to      Node
		weights map[helper.WeightKey]float64
	}{
		{
			name: "Test NewNetworkEdge successfully",
			id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1",
			from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
			to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
			weights: map[helper.WeightKey]float64{
				helper.IgpMetricKey:            10,
				helper.LatencyKey:              2000,
				helper.JitterKey:               100,
				helper.MaximumLinkBandwidth:    1000000,
				helper.AvailableBandwidthKey:   999000,
				helper.UtilizedBandwidthKey:    1000,
				helper.PacketLossKey:           0.1,
				helper.NormalizedLatencyKey:    0.2,
				helper.NormalizedJitterKey:     0.2,
				helper.NormalizedPacketLossKey: 0.2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edge := NewNetworkEdge(tt.id, tt.from, tt.to, tt.weights)
			assert.NotNil(t, edge)
		})
	}
}

func TestNetworkEdge_GetId(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		from    Node
		to      Node
		weights map[helper.WeightKey]float64
	}{
		{
			name: "Test NewNetworkEdge successfully",
			id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1",
			from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
			to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
			weights: map[helper.WeightKey]float64{
				helper.IgpMetricKey:            10,
				helper.LatencyKey:              2000,
				helper.JitterKey:               100,
				helper.MaximumLinkBandwidth:    1000000,
				helper.AvailableBandwidthKey:   999000,
				helper.UtilizedBandwidthKey:    1000,
				helper.PacketLossKey:           0.1,
				helper.NormalizedLatencyKey:    0.2,
				helper.NormalizedJitterKey:     0.2,
				helper.NormalizedPacketLossKey: 0.2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edge := NewNetworkEdge(tt.id, tt.from, tt.to, tt.weights)
			assert.Equal(t, tt.id, edge.GetId())
		})
	}
}

func TestNetworkEdge_From(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		from    Node
		to      Node
		weights map[helper.WeightKey]float64
	}{
		{
			name: "Test NewNetworkEdge successfully",
			id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1",
			from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
			to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
			weights: map[helper.WeightKey]float64{
				helper.IgpMetricKey:            10,
				helper.LatencyKey:              2000,
				helper.JitterKey:               100,
				helper.MaximumLinkBandwidth:    1000000,
				helper.AvailableBandwidthKey:   999000,
				helper.UtilizedBandwidthKey:    1000,
				helper.PacketLossKey:           0.1,
				helper.NormalizedLatencyKey:    0.2,
				helper.NormalizedJitterKey:     0.2,
				helper.NormalizedPacketLossKey: 0.2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edge := NewNetworkEdge(tt.id, tt.from, tt.to, tt.weights)
			assert.Equal(t, tt.from, edge.From())
		})
	}
}

func TestNetworkEdge_To(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		from    Node
		to      Node
		weights map[helper.WeightKey]float64
	}{
		{
			name: "Test NewNetworkEdge successfully",
			id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1",
			from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
			to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
			weights: map[helper.WeightKey]float64{
				helper.IgpMetricKey:            10,
				helper.LatencyKey:              2000,
				helper.JitterKey:               100,
				helper.MaximumLinkBandwidth:    1000000,
				helper.AvailableBandwidthKey:   999000,
				helper.UtilizedBandwidthKey:    1000,
				helper.PacketLossKey:           0.1,
				helper.NormalizedLatencyKey:    0.2,
				helper.NormalizedJitterKey:     0.2,
				helper.NormalizedPacketLossKey: 0.2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edge := NewNetworkEdge(tt.id, tt.from, tt.to, tt.weights)
			assert.Equal(t, tt.to, edge.To())
		})
	}
}

func TestNetworkEdge_GetAllWeights(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		from    Node
		to      Node
		weights map[helper.WeightKey]float64
	}{
		{
			name: "Test NewNetworkEdge successfully",
			id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1",
			from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
			to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
			weights: map[helper.WeightKey]float64{
				helper.IgpMetricKey:            10,
				helper.LatencyKey:              2000,
				helper.JitterKey:               100,
				helper.MaximumLinkBandwidth:    1000000,
				helper.AvailableBandwidthKey:   999000,
				helper.UtilizedBandwidthKey:    1000,
				helper.PacketLossKey:           0.1,
				helper.NormalizedLatencyKey:    0.2,
				helper.NormalizedJitterKey:     0.2,
				helper.NormalizedPacketLossKey: 0.2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edge := NewNetworkEdge(tt.id, tt.from, tt.to, tt.weights)
			assert.Equal(t, tt.weights, edge.GetAllWeights())
		})
	}
}

func TestNetworkEdge_GetWeight(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		from    Node
		to      Node
		weights map[helper.WeightKey]float64
	}{
		{
			name: "Test NewNetworkEdge successfully",
			id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1",
			from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
			to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
			weights: map[helper.WeightKey]float64{
				helper.IgpMetricKey:            10,
				helper.LatencyKey:              2000,
				helper.JitterKey:               100,
				helper.MaximumLinkBandwidth:    1000000,
				helper.AvailableBandwidthKey:   999000,
				helper.UtilizedBandwidthKey:    1000,
				helper.PacketLossKey:           0.1,
				helper.NormalizedLatencyKey:    0.2,
				helper.NormalizedJitterKey:     0.2,
				helper.NormalizedPacketLossKey: 0.2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edge := NewNetworkEdge(tt.id, tt.from, tt.to, tt.weights)
			for key, value := range tt.weights {
				assert.Equal(t, value, edge.GetWeight(key))
			}
		})
	}
}

func TestNetworkEdge_SetWeight(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		from    Node
		to      Node
		weights map[helper.WeightKey]float64
		newKey  helper.WeightKey
		newVal  float64
	}{
		{
			name: "Test NewNetworkEdge successfully",
			id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1",
			from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
			to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
			weights: map[helper.WeightKey]float64{
				helper.IgpMetricKey:            10,
				helper.LatencyKey:              2000,
				helper.JitterKey:               100,
				helper.MaximumLinkBandwidth:    1000000,
				helper.AvailableBandwidthKey:   999000,
				helper.UtilizedBandwidthKey:    1000,
				helper.PacketLossKey:           0.1,
				helper.NormalizedLatencyKey:    0.2,
				helper.NormalizedJitterKey:     0.2,
				helper.NormalizedPacketLossKey: 0.2,
			},
			newKey: helper.IgpMetricKey,
			newVal: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edge := NewNetworkEdge(tt.id, tt.from, tt.to, tt.weights)
			edge.SetWeight(tt.newKey, tt.newVal)
			assert.Equal(t, tt.newVal, edge.GetWeight(tt.newKey))
		})
	}
}

func TestNetworkEdge_GetFlexibleAlgorithms(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		from    Node
		to      Node
		weights map[helper.WeightKey]float64
		want    map[uint32]struct{}
	}{
		{
			name: "Test NewNetworkEdge successfully",
			id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1",
			from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
			to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
			weights: map[helper.WeightKey]float64{
				helper.IgpMetricKey:            10,
				helper.LatencyKey:              2000,
				helper.JitterKey:               100,
				helper.MaximumLinkBandwidth:    1000000,
				helper.AvailableBandwidthKey:   999000,
				helper.UtilizedBandwidthKey:    1000,
				helper.PacketLossKey:           0.1,
				helper.NormalizedLatencyKey:    0.2,
				helper.NormalizedJitterKey:     0.2,
				helper.NormalizedPacketLossKey: 0.2,
			},
			want: map[uint32]struct{}{128: {}},
		},
		{
			name: "Test NewNetworkEdge successfully",
			id:   "2_0_2_0_0000.0000.0003_2001:db8:12::2_0000.0000.0001_2001:db8:12::1",
			from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-2", []uint32{0, 1, 129}),
			to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
			weights: map[helper.WeightKey]float64{
				helper.IgpMetricKey:            10,
				helper.LatencyKey:              2000,
				helper.JitterKey:               100,
				helper.MaximumLinkBandwidth:    1000000,
				helper.AvailableBandwidthKey:   999000,
				helper.UtilizedBandwidthKey:    1000,
				helper.PacketLossKey:           0.1,
				helper.NormalizedLatencyKey:    0.2,
				helper.NormalizedJitterKey:     0.2,
				helper.NormalizedPacketLossKey: 0.2,
			},
			want: map[uint32]struct{}{129: {}},
		},
		{
			name: "Test NewNetworkEdge successfully",
			id:   "2_0_2_0_0000.0000.0003_2001:db8:24::2_0000.0000.0001_2001:db8:24::4",
			from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-2", []uint32{0, 1, 129}),
			to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-4", []uint32{0, 1}),
			weights: map[helper.WeightKey]float64{
				helper.IgpMetricKey:            10,
				helper.LatencyKey:              2000,
				helper.JitterKey:               100,
				helper.MaximumLinkBandwidth:    1000000,
				helper.AvailableBandwidthKey:   999000,
				helper.UtilizedBandwidthKey:    1000,
				helper.PacketLossKey:           0.1,
				helper.NormalizedLatencyKey:    0.2,
				helper.NormalizedJitterKey:     0.2,
				helper.NormalizedPacketLossKey: 0.2,
			},
			want: map[uint32]struct{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edge := NewNetworkEdge(tt.id, tt.from, tt.to, tt.weights)
			assert.Equal(t, tt.want, edge.GetFlexibleAlgorithms())
		})
	}
}

func TestNetworkEdge_UpdateFlexibleAlgorithms(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		from    Node
		to      Node
		weights map[helper.WeightKey]float64
		want    map[uint32]struct{}
	}{
		{
			name: "Test NewNetworkEdge successfully",
			id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1",
			from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
			to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
			weights: map[helper.WeightKey]float64{
				helper.IgpMetricKey:            10,
				helper.LatencyKey:              2000,
				helper.JitterKey:               100,
				helper.MaximumLinkBandwidth:    1000000,
				helper.AvailableBandwidthKey:   999000,
				helper.UtilizedBandwidthKey:    1000,
				helper.PacketLossKey:           0.1,
				helper.NormalizedLatencyKey:    0.2,
				helper.NormalizedJitterKey:     0.2,
				helper.NormalizedPacketLossKey: 0.2,
			},
			want: map[uint32]struct{}{128: {}, 129: {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edge := NewNetworkEdge(tt.id, tt.from, tt.to, tt.weights)
			tt.from.SetFlexibleAlgorithms([]uint32{128, 129})
			edge.UpdateFlexibleAlgorithms()
			assert.Equal(t, tt.want, edge.GetFlexibleAlgorithms())
		})
	}
}
