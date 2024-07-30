package graph

import (
	reflect "reflect"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/helper"
)

func TestNetworkNode_translateSrToFlexibleAlgorithm(t *testing.T) {
	tests := []struct {
		name         string
		srAlgorithms []uint32
		want         map[uint32]struct{}
	}{
		{
			name:         "Algo 0",
			srAlgorithms: []uint32{0},
			want:         map[uint32]struct{}{},
		},
		{
			name:         "Algo 1",
			srAlgorithms: []uint32{1},
			want:         map[uint32]struct{}{},
		},
		{
			name:         "Algo 0, 1",
			srAlgorithms: []uint32{0, 1},
			want:         map[uint32]struct{}{},
		},
		{
			name:         "Algo 0, 1, 128",
			srAlgorithms: []uint32{0, 1, 128},
			want:         map[uint32]struct{}{128: {}},
		},
		{
			name:         "Algo 0, 1, 128, 129",
			srAlgorithms: []uint32{0, 1, 128, 129},
			want:         map[uint32]struct{}{128: {}, 129: {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := translateSrToFlexibleAlgorithm(tt.srAlgorithms); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("translateSrToFlexibleAlgorithm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNetworkNode(t *testing.T) {
	tests := []struct {
		testName     string
		id           string
		name         string
		srAlgorithms []uint32
		want         *NetworkNode
	}{
		{
			testName:     "Add Node XR-1",
			id:           "2_0_0_0000.0000.0001",
			name:         "XR-1",
			srAlgorithms: []uint32{0, 1, 128, 129},
			want: &NetworkNode{
				id:                 "2_0_0_0000.0000.0001",
				name:               "XR-1",
				flexibleAlgorithms: map[uint32]struct{}{128: {}, 129: {}},
				edges:              map[string]Edge{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if got := NewNetworkNode(tt.id, tt.name, tt.srAlgorithms); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNetworkNode() testName %s = %v, want %v", tt.testName, got, tt.want)
			}
		})
	}
}

func TestNetworkNode_GetName(t *testing.T) {
	tests := []struct {
		testName     string
		id           string
		name         string
		srAlgorithms []uint32
		want         string
	}{
		{
			testName:     "Add Node XR-1",
			id:           "2_0_0_0000.0000.0001",
			name:         "XR-1",
			srAlgorithms: []uint32{0, 1, 128, 129},
			want:         "XR-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewNetworkNode(tt.id, tt.name, tt.srAlgorithms)
			if got := node.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworkNode_SetName(t *testing.T) {
	tests := []struct {
		testName     string
		id           string
		name         string
		srAlgorithms []uint32
		want         string
	}{
		{
			testName:     "Add Node XR-1",
			id:           "2_0_0_0000.0000.0001",
			name:         "XR-1",
			srAlgorithms: []uint32{0, 1, 128, 129},
			want:         "XR-2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewNetworkNode(tt.id, tt.name, tt.srAlgorithms)
			node.SetName("XR-2")
			if got := node.GetName(); got != tt.want {
				t.Errorf("SetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworkNode_GetId(t *testing.T) {
	tests := []struct {
		testName     string
		id           string
		name         string
		srAlgorithms []uint32
		want         string
	}{
		{
			testName:     "Add Node XR-1",
			id:           "2_0_0_0000.0000.0001",
			name:         "XR-1",
			srAlgorithms: []uint32{0, 1, 128, 129},
			want:         "2_0_0_0000.0000.0001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewNetworkNode(tt.id, tt.name, tt.srAlgorithms)
			if got := node.GetId(); got != tt.want {
				t.Errorf("GetId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworkNode_GetEdges(t *testing.T) {
	tests := []struct {
		testName     string
		id           string
		name         string
		srAlgorithms []uint32
		want         map[string]Edge
	}{
		{
			testName:     "Add Node XR-1",
			id:           "2_0_0_0000.0000.0001",
			name:         "XR-1",
			srAlgorithms: []uint32{0, 1, 128, 129},
			want: map[string]Edge{
				"2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1": &NetworkEdge{
					id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1",
					from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
					to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
					weights: map[helper.WeightKey]float64{
						helper.IgpMetricKey:            10,
						helper.LatencyKey:              2000,
						helper.JitterKey:               100,
						helper.MaximumLinkBandwidthKey: 1000000,
						helper.AvailableBandwidthKey:   999000,
						helper.UtilizedBandwidthKey:    1000,
						helper.PacketLossKey:           0.1,
						helper.NormalizedLatencyKey:    0.2,
						helper.NormalizedJitterKey:     0.2,
						helper.NormalizedPacketLossKey: 0.2,
					},
				},
				"2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0002_2001:db8:13::2": &NetworkEdge{
					id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0002_2001:db8:13::2",
					from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
					to:   NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
					weights: map[helper.WeightKey]float64{
						helper.IgpMetricKey:            10,
						helper.LatencyKey:              2000,
						helper.JitterKey:               100,
						helper.MaximumLinkBandwidthKey: 1000000,
						helper.AvailableBandwidthKey:   999000,
						helper.UtilizedBandwidthKey:    1000,
						helper.PacketLossKey:           0.1,
						helper.NormalizedLatencyKey:    0.2,
						helper.NormalizedJitterKey:     0.2,
						helper.NormalizedPacketLossKey: 0.2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewNetworkNode(tt.id, tt.name, tt.srAlgorithms)
			node.edges = tt.want
			if got := node.GetEdges(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEdges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworkNode_AddEdge(t *testing.T) {
	tests := []struct {
		testName     string
		id           string
		name         string
		srAlgorithms []uint32
		edge         Edge
		want         map[string]Edge
	}{
		{
			testName:     "Add Edge XR-3 to XR-1",
			id:           "2_0_0_0000.0000.0001",
			name:         "XR-1",
			srAlgorithms: []uint32{0, 1, 128, 129},
			edge: &NetworkEdge{
				id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1",
				from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
				to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
				weights: map[helper.WeightKey]float64{
					helper.IgpMetricKey:            10,
					helper.LatencyKey:              2000,
					helper.JitterKey:               100,
					helper.MaximumLinkBandwidthKey: 1000000,
					helper.AvailableBandwidthKey:   999000,
					helper.UtilizedBandwidthKey:    1000,
					helper.PacketLossKey:           0.1,
					helper.NormalizedLatencyKey:    0.2,
					helper.NormalizedJitterKey:     0.2,
					helper.NormalizedPacketLossKey: 0.2,
				},
			},
			want: map[string]Edge{
				"2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1": &NetworkEdge{
					id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1",
					from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
					to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
					weights: map[helper.WeightKey]float64{
						helper.IgpMetricKey:            10,
						helper.LatencyKey:              2000,
						helper.JitterKey:               100,
						helper.MaximumLinkBandwidthKey: 1000000,
						helper.AvailableBandwidthKey:   999000,
						helper.UtilizedBandwidthKey:    1000,
						helper.PacketLossKey:           0.1,
						helper.NormalizedLatencyKey:    0.2,
						helper.NormalizedJitterKey:     0.2,
						helper.NormalizedPacketLossKey: 0.2,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			node := NewNetworkNode(tt.id, tt.name, tt.srAlgorithms)
			node.AddEdge(tt.edge)
			if got := node.edges; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddEdge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworkNode_DeleteEdge(t *testing.T) {
	tests := []struct {
		testName     string
		id           string
		name         string
		srAlgorithms []uint32
		edges        map[string]Edge
		want         map[string]Edge
	}{
		{
			testName:     "Delete Edge from Node XR-1",
			id:           "2_0_0_0000.0000.0001",
			name:         "XR-1",
			srAlgorithms: []uint32{0, 1, 128, 129},
			edges: map[string]Edge{
				"2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1": &NetworkEdge{
					id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1",
					from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
					to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
					weights: map[helper.WeightKey]float64{
						helper.IgpMetricKey:            10,
						helper.LatencyKey:              2000,
						helper.JitterKey:               100,
						helper.MaximumLinkBandwidthKey: 1000000,
						helper.AvailableBandwidthKey:   999000,
						helper.UtilizedBandwidthKey:    1000,
						helper.PacketLossKey:           0.1,
						helper.NormalizedLatencyKey:    0.2,
						helper.NormalizedJitterKey:     0.2,
						helper.NormalizedPacketLossKey: 0.2,
					},
				},
				"2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0002_2001:db8:13::2": &NetworkEdge{
					id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0002_2001:db8:13::2",
					from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
					to:   NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
					weights: map[helper.WeightKey]float64{
						helper.IgpMetricKey:            10,
						helper.LatencyKey:              2000,
						helper.JitterKey:               100,
						helper.MaximumLinkBandwidthKey: 1000000,
						helper.AvailableBandwidthKey:   999000,
						helper.UtilizedBandwidthKey:    1000,
						helper.PacketLossKey:           0.1,
						helper.NormalizedLatencyKey:    0.2,
						helper.NormalizedJitterKey:     0.2,
						helper.NormalizedPacketLossKey: 0.2,
					},
				},
			},
			want: map[string]Edge{
				"2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1": &NetworkEdge{
					id:   "2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0001_2001:db8:13::1",
					from: NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
					to:   NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
					weights: map[helper.WeightKey]float64{
						helper.IgpMetricKey:            10,
						helper.LatencyKey:              2000,
						helper.JitterKey:               100,
						helper.MaximumLinkBandwidthKey: 1000000,
						helper.AvailableBandwidthKey:   999000,
						helper.UtilizedBandwidthKey:    1000,
						helper.PacketLossKey:           0.1,
						helper.NormalizedLatencyKey:    0.2,
						helper.NormalizedJitterKey:     0.2,
						helper.NormalizedPacketLossKey: 0.2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewNetworkNode(tt.id, tt.name, tt.srAlgorithms)
			node.edges = tt.edges
			node.DeleteEdge("2_0_2_0_0000.0000.0003_2001:db8:13::3_0000.0000.0002_2001:db8:13::2")
			if got := node.edges; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEdges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworkNode_SetFlexibleAlgorithms(t *testing.T) {
	tests := []struct {
		testName     string
		id           string
		name         string
		srAlgorithms []uint32
		edges        map[string]Edge
		want         map[uint32]struct{}
	}{
		{
			testName:     "Add Node XR-1",
			id:           "2_0_0_0000.0000.0001",
			name:         "XR-1",
			srAlgorithms: []uint32{0, 1, 128, 129},
			edges: map[string]Edge{
				"2_0_2_0_0000.0000.0001_2001:db8:12::1_0000.0000.0002_2001:db8:12::2": &NetworkEdge{
					id:   "2_0_2_0_0000.0000.0001_2001:db8:12::1_0000.0000.0002_2001:db8:12::2",
					from: NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128}),
					to:   NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128, 129}),
					weights: map[helper.WeightKey]float64{
						helper.IgpMetricKey:            10,
						helper.LatencyKey:              2000,
						helper.JitterKey:               100,
						helper.MaximumLinkBandwidthKey: 1000000,
						helper.AvailableBandwidthKey:   999000,
						helper.UtilizedBandwidthKey:    1000,
						helper.PacketLossKey:           0.1,
						helper.NormalizedLatencyKey:    0.2,
						helper.NormalizedJitterKey:     0.2,
						helper.NormalizedPacketLossKey: 0.2,
					},
				},
				"2_0_2_0_0000.0000.0001_2001:db8:13::1_0000.0000.0003_2001:db8:13::3": &NetworkEdge{
					id:   "2_0_2_0_0000.0000.0001_2001:db8:13::1_0000.0000.0003_2001:db8:13::3",
					from: NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128}),
					to:   NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128}),
					weights: map[helper.WeightKey]float64{
						helper.IgpMetricKey:            10,
						helper.LatencyKey:              2000,
						helper.JitterKey:               100,
						helper.MaximumLinkBandwidthKey: 1000000,
						helper.AvailableBandwidthKey:   999000,
						helper.UtilizedBandwidthKey:    1000,
						helper.PacketLossKey:           0.1,
						helper.NormalizedLatencyKey:    0.2,
						helper.NormalizedJitterKey:     0.2,
						helper.NormalizedPacketLossKey: 0.2,
					},
				},
			},
			want: map[uint32]struct{}{128: {}, 129: {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			node := NewNetworkNode(tt.id, tt.name, tt.srAlgorithms)
			node.edges = tt.edges
			node.SetFlexibleAlgorithms([]uint32{128, 129})
			if got := node.flexibleAlgorithms; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetFlexibleAlgorithms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworkNode_GetFlexibleAlgorithms(t *testing.T) {
	tests := []struct {
		testName     string
		id           string
		name         string
		srAlgorithms []uint32
		want         map[uint32]struct{}
	}{
		{
			testName:     "Add Node XR-1",
			id:           "2_0_0_0000.0000.0001",
			name:         "XR-1",
			srAlgorithms: []uint32{0, 1, 128, 129},
			want:         map[uint32]struct{}{128: {}, 129: {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			node := NewNetworkNode(tt.id, tt.name, tt.srAlgorithms)
			if got := node.GetFlexibleAlgorithms(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFlexibleAlgorithms() = %v, want %v", got, tt.want)
			}
		})
	}
}
