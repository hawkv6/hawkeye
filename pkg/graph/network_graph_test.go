package graph

import (
	"math"
	"reflect"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/helper"
)

func setupGraph(nodes map[int]Node, edges map[int]Edge) (*NetworkGraph, error) {
	graph := NewNetworkGraph(helper.NewDefaultHelper())
	for _, node := range nodes {
		if _, err := graph.AddNode(node); err != nil {
			return nil, err
		}
	}
	for _, edge := range edges {
		if err := graph.AddEdge(edge); err != nil {
			return nil, err
		}
	}
	return graph, nil
}

const tolerance = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestNetworkGraph_GetShortestPathDefaultMetric(t *testing.T) {
	helper := helper.NewDefaultHelper()
	nodes := map[int]Node{
		1: NewNetworkNode(1),
		2: NewNetworkNode(2),
		3: NewNetworkNode(3),
		4: NewNetworkNode(4),
		5: NewNetworkNode(5),
		6: NewNetworkNode(6),
		7: NewNetworkNode(7),
		8: NewNetworkNode(8),
	}
	// 	     [1]
	//      / | \
	//    1/ 2|  \1
	//    /   |   \
	//  [2]  [3]  [4]
	//   |1   |4   \1
	//   |    |     \
	//  [5]  [6]-1-[7]
	//   \    |    /
	//    \6  |1  /5
	//     \  |  /
	//       [8]

	edges := map[int]Edge{
		1:  NewNetworkEdge(1, nodes[1], nodes[2], map[string]float64{helper.GetDefaultKey(): 1, helper.GetLatencyKey(): 1000, helper.GetJitterKey(): 10, helper.GetPacketLossKey(): 0.01}),  // latency 1ms jitter 1us loss 1%
		2:  NewNetworkEdge(2, nodes[1], nodes[3], map[string]float64{helper.GetDefaultKey(): 2, helper.GetLatencyKey(): 2000, helper.GetJitterKey(): 20, helper.GetPacketLossKey(): 0.02}),  // latency 2ms jitter 2us loss 2%
		3:  NewNetworkEdge(3, nodes[1], nodes[4], map[string]float64{helper.GetDefaultKey(): 1, helper.GetLatencyKey(): 1000, helper.GetJitterKey(): 10, helper.GetPacketLossKey(): 0.01}),  // latency 1ms jitter 1us loss 1%
		4:  NewNetworkEdge(4, nodes[2], nodes[5], map[string]float64{helper.GetDefaultKey(): 1, helper.GetLatencyKey(): 1000, helper.GetJitterKey(): 10, helper.GetPacketLossKey(): 0.01}),  // latency 1ms jitter 1us loss 1%
		5:  NewNetworkEdge(5, nodes[3], nodes[5], map[string]float64{helper.GetDefaultKey(): 3, helper.GetLatencyKey(): 3000, helper.GetJitterKey(): 30, helper.GetPacketLossKey(): 0.03}),  // latency 3ms jitter 3us loss 3%
		6:  NewNetworkEdge(6, nodes[3], nodes[6], map[string]float64{helper.GetDefaultKey(): 4, helper.GetLatencyKey(): 4000, helper.GetJitterKey(): 40, helper.GetPacketLossKey(): 0.04}),  // latency 4ms jitter 4us loss 4%
		7:  NewNetworkEdge(7, nodes[4], nodes[7], map[string]float64{helper.GetDefaultKey(): 1, helper.GetLatencyKey(): 1000, helper.GetJitterKey(): 10, helper.GetPacketLossKey(): 0.01}),  // latency ms jitter 1us loss 1%
		8:  NewNetworkEdge(8, nodes[5], nodes[8], map[string]float64{helper.GetDefaultKey(): 6, helper.GetLatencyKey(): 6000, helper.GetJitterKey(): 60, helper.GetPacketLossKey(): 0.06}),  // latency 6ms jitter 6us loss 6%
		9:  NewNetworkEdge(9, nodes[6], nodes[8], map[string]float64{helper.GetDefaultKey(): 1, helper.GetLatencyKey(): 1000, helper.GetJitterKey(): 10, helper.GetPacketLossKey(): 0.01}),  // latency 1ms jitter 1us loss 1%
		10: NewNetworkEdge(10, nodes[7], nodes[6], map[string]float64{helper.GetDefaultKey(): 1, helper.GetLatencyKey(): 1000, helper.GetJitterKey(): 10, helper.GetPacketLossKey(): 0.01}), // latency 1ms jitter 1us loss 1%
		11: NewNetworkEdge(11, nodes[7], nodes[8], map[string]float64{helper.GetDefaultKey(): 5, helper.GetLatencyKey(): 5000, helper.GetJitterKey(): 50, helper.GetPacketLossKey(): 0.05}), // latency 5ms jitter 5us loss 5%
	}
	type Result struct {
		edgeNumbers []int
		totalCost   float64
	}
	type args struct {
		from       Node
		to         Node
		weightType string
	}
	tests := []struct {
		name    string
		args    args
		want    Result
		wantErr bool
	}{
		// {
		// 	name: "Test correct shortest path with default metric",
		// 	args: args{
		// 		from:       nodes[1],
		// 		to:         nodes[8],
		// 		weightType: helper.GetDefaultKey(),
		// 	},
		// 	want: Result{
		// 		edgeNumbers: []int{3, 7, 10, 9},
		// 		totalCost:   4,
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "Test correct shortest path with latency metric",
		// 	args: args{
		// 		from:       nodes[1],
		// 		to:         nodes[8],
		// 		weightType: helper.GetLatencyKey(),
		// 	},
		// 	want: Result{
		// 		edgeNumbers: []int{3, 7, 10, 9},
		// 		totalCost:   1000 + 1000 + 1000 + 1000, // latency of edge 3 + edge 7 + edge 10 + edge 9
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "Test correct shortest path with jitter metric",
		// 	args: args{
		// 		from:       nodes[1],
		// 		to:         nodes[8],
		// 		weightType: helper.GetJitterKey(),
		// 	},
		// 	want: Result{
		// 		edgeNumbers: []int{3, 7, 10, 9},
		// 		totalCost:   10 + 10 + 10 + 10, // jitter of edge 3 + edge 7 + edge 10 + edge 9
		// 	},
		// 	wantErr: false,
		// },
		{
			name: "Test correct shortest path with packet loss metric",
			args: args{
				from:       nodes[1],
				to:         nodes[8],
				weightType: helper.GetPacketLossKey(),
			},
			want: Result{
				edgeNumbers: []int{3, 7, 10, 9},
				totalCost:   0.01 * 0.01 * 0.01 * 0.01, // packet loss of edge 3 * edge 7 * edge 10 * edge 9
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph, err := setupGraph(nodes, edges)
			if err != nil {
				t.Errorf("Error setting up graph")
			}
			got, err := graph.GetShortestPath(tt.args.from, tt.args.to, tt.args.weightType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get shortest path with default metric %s, error = %v, wantErr %v", helper.GetDefaultKey(), err, tt.wantErr)
				return
			}
			shortestPath := make([]Edge, len(tt.want.edgeNumbers))
			for index, node := range tt.want.edgeNumbers {
				shortestPath[index] = edges[node]
			}
			if !reflect.DeepEqual(got.GetEdges(), shortestPath) {
				t.Errorf("DefaultGraph.GetShortestPath() = %v, want %v", got.GetEdges(), shortestPath)
			}
			if !almostEqual(got.GetCost(), tt.want.totalCost) {
				t.Errorf("DefaultGraph.GetShortestPath() = %v, want %v", got.GetCost(), tt.want.totalCost)
			}
		})
	}
}
