package graph

import (
	"testing"

	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/stretchr/testify/assert"
)

func TestNewShortestPath(t *testing.T) {
	bottleNeckEdge := &NetworkEdge{
		id:   "2_0_2_0_0000.0000.0001_2001:db8:12::1_0000.0000.0002_2001:db8:12::2",
		from: NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
		to:   NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
		weights: map[helper.WeightKey]float64{
			helper.IgpMetricKey:            100,
			helper.LatencyKey:              20000,
			helper.JitterKey:               1000,
			helper.MaximumLinkBandwidthKey: 100000,
			helper.AvailableBandwidthKey:   99000,
			helper.UtilizedBandwidthKey:    1000,
			helper.PacketLossKey:           5.0,
			helper.NormalizedLatencyKey:    0.5,
			helper.NormalizedJitterKey:     0.5,
			helper.NormalizedPacketLossKey: 0.5,
		},
	}
	tests := []struct {
		testName        string
		edges           []Edge
		totalCost       float64
		delay           float64
		jitter          float64
		packetLoss      float64
		bottleNeckValue float64
		bottleneckEdge  Edge
	}{
		{
			testName: "TestNewShortestPath",
			edges: []Edge{
				bottleNeckEdge,
				&NetworkEdge{
					id:   "2_0_2_0_0000.0000.0002_2001:db8:23::2_0000.0000.0003_2001:db8:23::3",
					from: NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
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
			totalCost:       22000, // assuming intent is low latency
			delay:           22000,
			jitter:          1100,
			packetLoss:      (1 - ((1 - 0.05) * (1 - 0.001))) * 100,
			bottleNeckValue: 99000,
			bottleneckEdge:  bottleNeckEdge,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(tt.edges, tt.totalCost, tt.delay, tt.jitter, tt.packetLoss, tt.bottleNeckValue, tt.bottleneckEdge)
			assert.NotNil(t, shortestPath)
		})
	}
}

func TestShortestPath_GetEdges(t *testing.T) {
	tests := []struct {
		testName string
		edges    []Edge
	}{
		{
			testName: "TestShortestPath_GetEdges",
			edges: []Edge{
				&NetworkEdge{
					id:   "2_0_2_0_0000.0000.0001_2001:db8:12::1_0000.0000.0002_2001:db8:12::2",
					from: NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
					to:   NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
					weights: map[helper.WeightKey]float64{
						helper.IgpMetricKey:            100,
						helper.LatencyKey:              20000,
						helper.JitterKey:               1000,
						helper.MaximumLinkBandwidthKey: 100000,
						helper.AvailableBandwidthKey:   99000,
						helper.UtilizedBandwidthKey:    1000,
						helper.PacketLossKey:           5.0,
						helper.NormalizedLatencyKey:    0.5,
						helper.NormalizedJitterKey:     0.5,
						helper.NormalizedPacketLossKey: 0.5,
					},
				},
				&NetworkEdge{
					id:   "2_0_2_0_0000.0000.0002_2001:db8:23::2_0000.0000.0003_2001:db8:23::3",
					from: NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(tt.edges, 0, 0, 0, 0, 0, nil)
			assert.Equal(t, tt.edges, shortestPath.GetEdges())
		})
	}
}

func TestShortestPath_GetTotalCost(t *testing.T) {
	tests := []struct {
		testName  string
		totalCost float64
	}{
		{
			testName:  "TestShortestPath_GetTotalCost",
			totalCost: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(nil, tt.totalCost, 0, 0, 0, 0, nil)
			assert.Equal(t, tt.totalCost, shortestPath.GetTotalCost())
		})
	}
}

func TestShortestPath_SetTotalCost(t *testing.T) {
	tests := []struct {
		testName  string
		totalCost float64
		newCost   float64
	}{
		{
			testName:  "TestShortestPath_SetTotalCost",
			totalCost: 100,
			newCost:   200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(nil, tt.totalCost, 0, 0, 0, 0, nil)
			shortestPath.SetTotalCost(tt.newCost)
			assert.Equal(t, tt.newCost, shortestPath.totalCost)
		})
	}
}

func TestShortestPath_GetTotalDelay(t *testing.T) {
	tests := []struct {
		testName string
		delay    float64
	}{
		{
			testName: "TestShortestPath_GetTotalDelay",
			delay:    100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(nil, 0, tt.delay, 0, 0, 0, nil)
			assert.Equal(t, tt.delay, shortestPath.GetTotalDelay())
		})
	}
}

func TestShortestPath_SetTotalDelay(t *testing.T) {
	tests := []struct {
		testName string
		delay    float64
		newDelay float64
	}{
		{
			testName: "TestShortestPath_SetTotalDelay",
			delay:    100,
			newDelay: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(nil, 0, tt.delay, 0, 0, 0, nil)
			shortestPath.SetTotalDelay(tt.newDelay)
			assert.Equal(t, tt.newDelay, shortestPath.totalDelay)
		})
	}
}

func TestShortestPath_GetTotalJitter(t *testing.T) {
	tests := []struct {
		testName string
		jitter   float64
	}{
		{
			testName: "TestShortestPath_GetTotalJitter",
			jitter:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(nil, 0, 0, tt.jitter, 0, 0, nil)
			assert.Equal(t, tt.jitter, shortestPath.GetTotalJitter())
		})
	}
}

func TestShortestPath_SetTotalJitter(t *testing.T) {
	tests := []struct {
		testName  string
		jitter    float64
		newJitter float64
	}{
		{
			testName:  "TestShortestPath_SetTotalJitter",
			jitter:    100,
			newJitter: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(nil, 0, 0, tt.jitter, 0, 0, nil)
			shortestPath.SetTotalJitter(tt.newJitter)
			assert.Equal(t, tt.newJitter, shortestPath.totalJitter)
		})
	}
}

func TestShortestPath_GetTotalPacketLoss(t *testing.T) {
	tests := []struct {
		testName   string
		packetLoss float64
	}{
		{
			testName:   "TestShortestPath_GetTotalPacketLoss",
			packetLoss: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(nil, 0, 0, 0, tt.packetLoss, 0, nil)
			assert.Equal(t, tt.packetLoss, shortestPath.GetTotalPacketLoss())
		})
	}
}

func TestShortestPath_SetTotalPacketLoss(t *testing.T) {
	tests := []struct {
		testName   string
		packetLoss float64
		newLoss    float64
	}{
		{
			testName:   "TestShortestPath_SetTotalPacketLoss",
			packetLoss: 100,
			newLoss:    200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(nil, 0, 0, 0, tt.packetLoss, 0, nil)
			shortestPath.SetTotalPacketLoss(tt.newLoss)
			assert.Equal(t, tt.newLoss, shortestPath.totalPacketLoss)
		})
	}
}

func TestShortestPath_GetBottleneckEdge(t *testing.T) {
	bottleNeckEdge := &NetworkEdge{
		id:   "2_0_2_0_0000.0000.0001_2001:db8:12::1_0000.0000.0002_2001:db8:12::2",
		from: NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
		to:   NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
		weights: map[helper.WeightKey]float64{
			helper.IgpMetricKey:            100,
			helper.LatencyKey:              20000,
			helper.JitterKey:               1000,
			helper.MaximumLinkBandwidthKey: 100000,
			helper.AvailableBandwidthKey:   99000,
			helper.UtilizedBandwidthKey:    1000,
			helper.PacketLossKey:           5.0,
			helper.NormalizedLatencyKey:    0.5,
			helper.NormalizedJitterKey:     0.5,
			helper.NormalizedPacketLossKey: 0.5,
		},
	}
	tests := []struct {
		testName string
		edges    []Edge
	}{
		{
			testName: "TestShortestPath_GetBottleneckEdge",
			edges: []Edge{
				bottleNeckEdge,
				&NetworkEdge{
					id:   "2_0_2_0_0000.0000.0002_2001:db8:23::2_0000.0000.0003_2001:db8:23::3",
					from: NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(tt.edges, 0, 0, 0, 0, 0, bottleNeckEdge)
			assert.Equal(t, bottleNeckEdge, shortestPath.GetBottleneckEdge())
		})
	}
}

func TestShortestPath_GetBottleneckValue(t *testing.T) {
	tests := []struct {
		testName        string
		bottleNeckValue float64
	}{
		{
			testName:        "TestShortestPath_GetBottleneckValue",
			bottleNeckValue: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(nil, 0, 0, 0, 0, tt.bottleNeckValue, nil)
			assert.Equal(t, tt.bottleNeckValue, shortestPath.GetBottleneckValue())
		})
	}
}

func TestShortestPath_SetBottleneckEdge(t *testing.T) {
	bottleNeckEdge := &NetworkEdge{
		id:   "2_0_2_0_0000.0000.0001_2001:db8:12::1_0000.0000.0002_2001:db8:12::2",
		from: NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
		to:   NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
		weights: map[helper.WeightKey]float64{
			helper.IgpMetricKey:            100,
			helper.LatencyKey:              20000,
			helper.JitterKey:               1000,
			helper.MaximumLinkBandwidthKey: 100000,
			helper.AvailableBandwidthKey:   99000,
			helper.UtilizedBandwidthKey:    1000,
			helper.PacketLossKey:           5.0,
			helper.NormalizedLatencyKey:    0.5,
			helper.NormalizedJitterKey:     0.5,
			helper.NormalizedPacketLossKey: 0.5,
		},
	}
	tests := []struct {
		testName string
		edges    []Edge
	}{
		{
			testName: "TestShortestPath_SetBottleneckEdge",
			edges: []Edge{
				bottleNeckEdge,
				&NetworkEdge{
					id:   "2_0_2_0_0000.0000.0002_2001:db8:23::2_0000.0000.0003_2001:db8:23::3",
					from: NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(tt.edges, 0, 0, 0, 0, 0, nil)
			shortestPath.SetBottleneckEdge(bottleNeckEdge)
			assert.Equal(t, bottleNeckEdge, shortestPath.bottleneckEdge)
		})
	}
}

func TestShortestPath_SetBottleneckValue(t *testing.T) {
	tests := []struct {
		testName        string
		bottleNeckValue float64
		newValue        float64
	}{
		{
			testName:        "TestShortestPath_SetBottleneckValue",
			bottleNeckValue: 100,
			newValue:        200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(nil, 0, 0, 0, 0, tt.bottleNeckValue, nil)
			shortestPath.SetBottleneckValue(tt.newValue)
			assert.Equal(t, tt.newValue, shortestPath.bottleneckValue)
		})
	}
}

func TestShortestPath_SetRouterServiceMap(t *testing.T) {
	tests := []struct {
		testName         string
		routerServiceMap map[string]string
	}{
		{
			testName: "TestShortestPath_SetRouterServiceMap",
			routerServiceMap: map[string]string{
				"XR-1": "fw",
				"XR-2": "fw",
				"XR-6": "ids",
				"XR-7": "ids",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(nil, 0, 0, 0, 0, 0, nil)
			shortestPath.SetRouterServiceMap(tt.routerServiceMap)
			assert.Equal(t, tt.routerServiceMap, shortestPath.routerServiceMap)
		})
	}
}

func TestShortestPath_GetRouterServiceMap(t *testing.T) {
	tests := []struct {
		testName         string
		routerServiceMap map[string]string
	}{
		{
			testName: "TestShortestPath_GetRouterServiceMap",
			routerServiceMap: map[string]string{
				"XR-1": "fw",
				"XR-2": "fw",
				"XR-6": "ids",
				"XR-7": "ids",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			shortestPath := NewShortestPath(nil, 0, 0, 0, 0, 0, nil)
			shortestPath.routerServiceMap = tt.routerServiceMap
			assert.Equal(t, tt.routerServiceMap, shortestPath.GetRouterServiceMap())
		})
	}
}
