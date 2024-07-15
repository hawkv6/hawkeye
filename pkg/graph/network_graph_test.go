package graph

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/stretchr/testify/assert"
)

func TestNewNetworkGraph(t *testing.T) {
	tests := []struct {
		testName string
	}{
		{
			testName: "TestNewNetworkGraph",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			assert.NotNil(t, graph)
		})
	}
}

func TestNetworkGraph_Lock(t *testing.T) {
	tests := []struct {
		testName string
	}{
		{
			testName: "TestNetworkGraph Lock()",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			wg := &sync.WaitGroup{}
			for i := 1; i <= 8; i++ {
				wg.Add(1)
				go func(i int) {
					graph.Lock()
					defer graph.mu.Unlock()
					defer wg.Done()

					fromNodeID := fmt.Sprintf("Node%d", i)
					toNodeID := fmt.Sprintf("Node%d", i+1)
					edgeID := fmt.Sprintf("%s_to_%s", fromNodeID, toNodeID)

					fromNode := NewNetworkNode(fromNodeID, "XR-"+strconv.Itoa(i), []uint32{0, 1, 128})
					toNode := NewNetworkNode(toNodeID, "XR-"+strconv.Itoa(i+1), []uint32{0, 1, 128, 129})

					edge := &NetworkEdge{
						id:   edgeID,
						from: fromNode,
						to:   toNode,
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
					}
					graph.edges[edgeID] = edge
				}(i)
			}
			wg.Wait()
			assert.Equal(t, 8, len(graph.edges))
		})
	}
}

func TestNetworkGraph_Unlock(t *testing.T) {
	tests := []struct {
		testName string
	}{
		{
			testName: "TestNetworkGraph Lock()",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			wg := &sync.WaitGroup{}
			for i := 1; i <= 8; i++ {
				wg.Add(1)
				go func(i int) {
					graph.mu.Lock()
					defer graph.Unlock()
					defer wg.Done()

					fromNodeID := fmt.Sprintf("Node%d", i)
					toNodeID := fmt.Sprintf("Node%d", i+1)
					edgeID := fmt.Sprintf("%s_to_%s", fromNodeID, toNodeID)

					fromNode := NewNetworkNode(fromNodeID, "XR-"+strconv.Itoa(i), []uint32{0, 1, 128})
					toNode := NewNetworkNode(toNodeID, "XR-"+strconv.Itoa(i+1), []uint32{0, 1, 128, 129})

					edge := &NetworkEdge{
						id:   edgeID,
						from: fromNode,
						to:   toNode,
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
					}
					graph.edges[edgeID] = edge
				}(i)
			}

			wg.Wait()
			assert.Equal(t, 8, len(graph.edges))
		})
	}
}

func TestNetworkGraph_NodeExists(t *testing.T) {
	tests := []struct {
		testName     string
		id           string
		name         string
		srAlgorithms []uint32
		exists       bool
	}{
		{
			testName:     "TestNetworkGraph NodeExists() true",
			id:           "2_0_0_0000.0000.0001",
			name:         "XR-1",
			srAlgorithms: []uint32{0, 1, 128, 129},
			exists:       true,
		},
		{
			testName:     "TestNetworkGraph NodeExists() false",
			id:           "2_0_0_0000.0000.0002",
			name:         "XR-2",
			srAlgorithms: []uint32{0, 1, 128, 129},
			exists:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			if !tt.exists {
				assert.False(t, graph.NodeExists(tt.id))
				return
			}
			node := NewNetworkNode(tt.id, tt.name, tt.srAlgorithms)
			graph.nodes[tt.id] = node
			assert.True(t, graph.NodeExists(tt.id))
		})
	}
}

func TestNetworkGraph_GetNode(t *testing.T) {
	tests := []struct {
		testName     string
		id           string
		name         string
		srAlgorithms []uint32
		exists       bool
	}{
		{
			testName:     "TestNetworkGraph GetNode() exists",
			id:           "2_0_0_0000.0000.0001",
			name:         "XR-1",
			srAlgorithms: []uint32{0, 1, 128, 129},
			exists:       true,
		},
		{
			testName:     "TestNetworkGraph GetNode() does not exist",
			id:           "2_0_0_0000.0000.0002",
			name:         "XR-2",
			srAlgorithms: []uint32{0, 1, 128, 129},
			exists:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			if !tt.exists {
				assert.Nil(t, graph.GetNode(tt.id))
				return
			}
			node := NewNetworkNode(tt.id, tt.name, tt.srAlgorithms)
			graph.nodes[tt.id] = node
			assert.Equal(t, node, graph.GetNode(tt.id))
		})
	}
}

func TestNetworkGraph_GetNodes(t *testing.T) {
	tests := []struct {
		testName string
		nodes    map[string]Node
	}{
		{
			testName: "TestNetworkGraph GetNodes()",
			nodes: map[string]Node{
				"2_0_0_0000.0000.0001": &NetworkNode{
					id:   "2_0_0_0000.0000.0001",
					name: "XR-1",
				},
				"2_0_0_0000.0000.0002": &NetworkNode{
					id:   "2_0_0_0000.0000.0002",
					name: "XR-2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			graph.nodes = tt.nodes
			assert.Equal(t, tt.nodes, graph.GetNodes())
		})
	}
}

func TestNetworkGraph_AddNode(t *testing.T) {
	tests := []struct {
		testName     string
		id           string
		name         string
		srAlgorithms []uint32
	}{
		{
			testName:     "TestNetworkGraph AddNode()",
			id:           "2_0_0_0000.0000.0001",
			name:         "XR-1",
			srAlgorithms: []uint32{0, 1, 128, 129},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			assert.Nil(t, graph.nodes[tt.id])
			node := NewNetworkNode(tt.id, tt.name, tt.srAlgorithms)
			graph.AddNode(node)
			assert.Equal(t, node, graph.nodes[tt.id])
		})
	}
}

func TestNetworkGraph_DeleteNode(t *testing.T) {
	tests := []struct {
		testName     string
		id           string
		name         string
		srAlgorithms []uint32
		edges        map[string]Edge
	}{
		{
			testName:     "TestNetworkGraph DeleteNode()",
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
						helper.MaximumLinkBandwidth:    1000000,
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
						helper.MaximumLinkBandwidth:    1000000,
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
			graph := NewNetworkGraph()
			node := NewNetworkNode(tt.id, tt.name, tt.srAlgorithms)
			node.edges = tt.edges
			assert.Nil(t, graph.nodes[tt.id])
			graph.nodes[tt.id] = node
			graph.edges = tt.edges
			graph.DeleteNode(node)
			_, exists := graph.nodes[tt.id]
			assert.False(t, exists)
			for _, edge := range tt.edges {
				_, exists := graph.edges[edge.GetId()]
				assert.False(t, exists)
			}
		})
	}
}

func TestNetworkGraph_GetEdge(t *testing.T) {
	tests := []struct {
		testName string
		id       string
		from     Node
		to       Node
		weights  map[helper.WeightKey]float64
	}{
		{
			testName: "TestNetworkGraph GetEdge()",
			id:       "2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002",
			from: &NetworkNode{
				id:   "2_0_0_0000.0000.0001",
				name: "XR-1",
			},
			to: &NetworkNode{
				id:   "2_0_0_0000.0000.0002",
				name: "XR-2",
			},
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
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			edge := &NetworkEdge{
				id:      tt.id,
				from:    tt.from,
				to:      tt.to,
				weights: tt.weights,
			}
			assert.Nil(t, graph.GetEdge(tt.id))
			graph.edges[tt.id] = edge
			assert.Equal(t, edge, graph.GetEdge(tt.id))
		})
	}
}

func TestNetworkGraph_GetEdges(t *testing.T) {
	tests := []struct {
		testName string
		edges    map[string]Edge
	}{
		{
			testName: "TestNetworkGraph GetEdges()",
			edges: map[string]Edge{
				"2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002": &NetworkEdge{
					id: "2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002",
					from: &NetworkNode{
						id:   "2_0_0_0000.0000.0001",
						name: "XR-1",
					},
					to: &NetworkNode{
						id:   "2_0_0_0000.0000.0002",
						name: "XR-2",
					},
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
				"2_0_0_0000.0000.0002_to_2_0_0_0000.0000.0003": &NetworkEdge{
					id: "2_0_0_0000.0000.0002_to_2_0_0_0000.0000.0003",
					from: &NetworkNode{
						id:   "2_0_0_0000.0000.0002",
						name: "XR-2",
					},
					to: &NetworkNode{
						id:   "2_0_0_0000.0000.0003",
						name: "XR-3",
					},
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			assert.Equal(t, 0, len(graph.GetEdges()))
			graph.edges = tt.edges
			assert.Equal(t, tt.edges, graph.GetEdges())
		})
	}
}

func TestNetworkGraph_EdgeExists(t *testing.T) {
	tests := []struct {
		testName string
		id       string
		from     Node
		to       Node
		weights  map[helper.WeightKey]float64
	}{
		{
			testName: "TestNetworkGraph EdgeExists()",
			id:       "2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002",
			from: &NetworkNode{
				id:   "2_0_0_0000.0000.0001",
				name: "XR-1",
			},
			to: &NetworkNode{
				id:   "2_0_0_0000.0000.0002",
				name: "XR-2",
			},
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
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			edge := &NetworkEdge{
				id:      tt.id,
				from:    tt.from,
				to:      tt.to,
				weights: tt.weights,
			}
			assert.False(t, graph.EdgeExists(tt.id))
			graph.edges[tt.id] = edge
			assert.True(t, graph.EdgeExists(tt.id))
		})
	}
}
func TestNetworkGraph_AddEdge(t *testing.T) {
	from := NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128})
	to := NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128, 129})
	tests := []struct {
		testName  string
		to        Node
		from      Node
		edge      Edge
		wantErr   bool
		duplicate bool
	}{
		{
			testName: "TestNetworkGraph AddEdge - Valid Edge",
			from:     from,
			to:       to,
			edge: NewNetworkEdge("2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002", from, to, map[helper.WeightKey]float64{
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
			}),
			wantErr: false,
		},
		{
			testName: "TestNetworkGraph AddEdge - from node does not exist",
			from:     nil,
			to:       to,
			edge: NewNetworkEdge("2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002", from, to, map[helper.WeightKey]float64{
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
			}),
			wantErr: true,
		},
		{
			testName: "TestNetworkGraph AddEdge - to node does not exist",
			from:     from,
			to:       nil,
			edge: NewNetworkEdge("2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002", from, to, map[helper.WeightKey]float64{
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
			}),
			wantErr: true,
		},
		{
			testName:  "TestNetworkGraph AddEdge - edge already exists",
			from:      from,
			to:        to,
			duplicate: true,
			edge: NewNetworkEdge("2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002", from, to, map[helper.WeightKey]float64{
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
			}),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			if tt.from != nil {
				graph.nodes[tt.from.GetId()] = tt.from
			}
			if tt.to != nil {
				graph.nodes[tt.to.GetId()] = tt.to
			}
			if tt.duplicate {
				err := graph.AddEdge(tt.edge)
				assert.Nil(t, err)
			}
			err := graph.AddEdge(tt.edge)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddEdge() testname %s - error = %v, wantErr %v", tt.testName, err, tt.wantErr)
			}
		})
	}
}

func TestNetworkGraph_DeleteEdge(t *testing.T) {
	from := NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128})
	to := NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128, 129})
	tests := []struct {
		testName string
		edge     Edge
	}{
		{
			testName: "TestNetworkGraph DeleteEdge - Valid Edge",
			edge: NewNetworkEdge("2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002", from, to, map[helper.WeightKey]float64{
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
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			graph.edges[tt.edge.GetId()] = tt.edge
			from.edges[tt.edge.GetId()] = tt.edge
			to.edges[tt.edge.GetId()] = tt.edge
			assert.NotNil(t, graph.edges[tt.edge.GetId()])
			assert.NotNil(t, from.edges[tt.edge.GetId()])
			assert.NotNil(t, to.edges[tt.edge.GetId()])
			graph.DeleteEdge(tt.edge)
			_, exists := graph.edges[tt.edge.GetId()]
			assert.False(t, exists)
			assert.Nil(t, from.edges[tt.edge.GetId()])
			assert.Nil(t, to.edges[tt.edge.GetId()])
		})
	}
}

func TestNetworkGraph_addNodesToSubgraph(t *testing.T) {
	tests := []struct {
		testName string
		nodes    map[string]Node
	}{
		{
			testName: "TestNetworkGraph addNodesToSubgraph",
			nodes: map[string]Node{

				"2_0_0_0000.0000.0001": NewNetworkNode(
					"2_0_0_0000.0000.0001",
					"XR-1",
					[]uint32{0, 1, 128, 129},
				),
				"2_0_0_0000.0000.0002": NewNetworkNode(
					"2_0_0_0000.0000.0002",
					"XR-2",
					[]uint32{0, 1, 128},
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			graph.nodes = tt.nodes
			newSubGraphs := make(map[uint32]*NetworkGraph)
			graph.addNodesToSubgraph(newSubGraphs)
			assert.Equal(t, 2, len(newSubGraphs))
			assert.Equal(t, 2, len(newSubGraphs[128].nodes))
			assert.Equal(t, 1, len(newSubGraphs[129].nodes))
		})
	}
}

func TestNetworkGraph_addEdgesToSubgraph(t *testing.T) {
	tests := []struct {
		testName string
		edges    map[string]Edge
		wantErr  bool
	}{
		{
			testName: "TestNetworkGraph addEdgesToSubgraph",
			edges: map[string]Edge{
				"2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002": NewNetworkEdge(
					"2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002",
					NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
					NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
					map[helper.WeightKey]float64{
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
				),
				"2_0_0_0000.0000.0002_to_2_0_0_0000.0000.0003": NewNetworkEdge(
					"2_0_0_0000.0000.0002_to_2_0_0_0000.0000.0003",
					NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
					NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128, 129}),
					map[helper.WeightKey]float64{
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
				),
			},
			wantErr: false,
		},
		{
			testName: "TestNetworkGraph addEdgesToSubgraph failure node does not exist",
			edges: map[string]Edge{
				"2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002": NewNetworkEdge(
					"2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002",
					NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
					NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
					map[helper.WeightKey]float64{
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
				),
				"2_0_0_0000.0000.0002_to_2_0_0_0000.0000.0003": NewNetworkEdge(
					"2_0_0_0000.0000.0002_to_2_0_0_0000.0000.0003",
					NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
					NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128, 129}),
					map[helper.WeightKey]float64{
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
				),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			for _, edge := range tt.edges {
				graph.nodes[edge.From().GetId()] = edge.From()
				graph.nodes[edge.To().GetId()] = edge.To()
			}
			graph.edges = tt.edges
			newSubGraphs := make(map[uint32]*NetworkGraph)
			graph.addNodesToSubgraph(newSubGraphs)
			if tt.wantErr {
				newSubGraphs[128].nodes = nil // triggers error log message
				graph.addEdgesToSubgraph(newSubGraphs)
				return
			}
			graph.addEdgesToSubgraph(newSubGraphs)
			assert.Equal(t, 2, len(newSubGraphs[128].edges))
			assert.Equal(t, 0, len(newSubGraphs[129].edges))
		})
	}
}

func TestNetworkGraph_UpdateSubGraphs(t *testing.T) {
	tests := []struct {
		testName string
		edges    map[string]Edge
	}{
		{
			testName: "TestNetworkGraph UpdateSubGraphs",
			edges: map[string]Edge{
				"2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002": NewNetworkEdge(
					"2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002",
					NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
					NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
					map[helper.WeightKey]float64{
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
				),
				"2_0_0_0000.0000.0002_to_2_0_0_0000.0000.0003": NewNetworkEdge(
					"2_0_0_0000.0000.0002_to_2_0_0_0000.0000.0003",
					NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
					NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128, 129}),
					map[helper.WeightKey]float64{
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
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			for _, edge := range tt.edges {
				graph.nodes[edge.From().GetId()] = edge.From()
				graph.nodes[edge.To().GetId()] = edge.To()
			}
			graph.edges = tt.edges
			graph.UpdateSubGraphs()
			assert.Equal(t, 2, len(graph.subGraphs))
			assert.Equal(t, 3, len(graph.subGraphs[128].nodes))
			assert.Equal(t, 2, len(graph.subGraphs[129].nodes))
			assert.Equal(t, 2, len(graph.subGraphs[128].edges))
			assert.Equal(t, 0, len(graph.subGraphs[129].edges))
		})
	}
}

func TestNetworkGraph_GetSubGraphs(t *testing.T) {
	tests := []struct {
		testName string
		edges    map[string]Edge
	}{
		{
			testName: "TestNetworkGraph GetSubGraphs",
			edges: map[string]Edge{
				"2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002": NewNetworkEdge(
					"2_0_0_0000.0000.0001_to_2_0_0_0000.0000.0002",
					NewNetworkNode("2_0_0_0000.0000.0001", "XR-1", []uint32{0, 1, 128, 129}),
					NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
					map[helper.WeightKey]float64{
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
				),
				"2_0_0_0000.0000.0002_to_2_0_0_0000.0000.0003": NewNetworkEdge(
					"2_0_0_0000.0000.0002_to_2_0_0_0000.0000.0003",
					NewNetworkNode("2_0_0_0000.0000.0002", "XR-2", []uint32{0, 1, 128}),
					NewNetworkNode("2_0_0_0000.0000.0003", "XR-3", []uint32{0, 1, 128, 129}),
					map[helper.WeightKey]float64{
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
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			graph := NewNetworkGraph()
			for _, edge := range tt.edges {
				graph.nodes[edge.From().GetId()] = edge.From()
				graph.nodes[edge.To().GetId()] = edge.To()
			}
			graph.edges = tt.edges
			graph.UpdateSubGraphs()
			subGraph := graph.GetSubGraph(128)
			assert.Equal(t, 3, len(subGraph.nodes))
			assert.Equal(t, 2, len(subGraph.edges))
			subGraph = graph.GetSubGraph(129)
			assert.Equal(t, 2, len(subGraph.nodes))
			assert.Equal(t, 0, len(subGraph.edges))
		})
	}
}
