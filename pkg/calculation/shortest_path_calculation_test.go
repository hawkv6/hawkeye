package calculation

import (
	"math"
	"reflect"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
)

func setupGraph(nodes map[int]graph.Node, edges map[int]graph.Edge) (*graph.NetworkGraph, error) {
	graph := graph.NewNetworkGraph(helper.NewDefaultHelper())
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

func TestNetworkGraph_GetShortestPath(t *testing.T) {
	nodes := map[int]graph.Node{
		1: graph.NewNetworkNode(1),
		2: graph.NewNetworkNode(2),
		3: graph.NewNetworkNode(3),
		4: graph.NewNetworkNode(4),
		5: graph.NewNetworkNode(5),
		6: graph.NewNetworkNode(6),
		7: graph.NewNetworkNode(7),
		8: graph.NewNetworkNode(8),
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

	edges := map[int]graph.Edge{
		1:  graph.NewNetworkEdge(1, nodes[1], nodes[2], map[helper.WeightKey]float64{helper.DefaultKey: 1, helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 0.01}),  // latency 1ms jitter 1us loss 1%
		2:  graph.NewNetworkEdge(2, nodes[1], nodes[3], map[helper.WeightKey]float64{helper.DefaultKey: 2, helper.LatencyKey: 2000, helper.JitterKey: 20, helper.PacketLossKey: 0.02}),  // latency 2ms jitter 2us loss 2%
		3:  graph.NewNetworkEdge(3, nodes[1], nodes[4], map[helper.WeightKey]float64{helper.DefaultKey: 1, helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 0.01}),  // latency 1ms jitter 1us loss 1%
		4:  graph.NewNetworkEdge(4, nodes[2], nodes[5], map[helper.WeightKey]float64{helper.DefaultKey: 1, helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 0.01}),  // latency 1ms jitter 1us loss 1%
		5:  graph.NewNetworkEdge(5, nodes[3], nodes[5], map[helper.WeightKey]float64{helper.DefaultKey: 3, helper.LatencyKey: 3000, helper.JitterKey: 30, helper.PacketLossKey: 0.03}),  // latency 3ms jitter 3us loss 3%
		6:  graph.NewNetworkEdge(6, nodes[3], nodes[6], map[helper.WeightKey]float64{helper.DefaultKey: 4, helper.LatencyKey: 4000, helper.JitterKey: 40, helper.PacketLossKey: 0.04}),  // latency 4ms jitter 4us loss 4%
		7:  graph.NewNetworkEdge(7, nodes[4], nodes[7], map[helper.WeightKey]float64{helper.DefaultKey: 1, helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 0.01}),  // latency ms jitter 1us loss 1%
		8:  graph.NewNetworkEdge(8, nodes[5], nodes[8], map[helper.WeightKey]float64{helper.DefaultKey: 6, helper.LatencyKey: 6000, helper.JitterKey: 60, helper.PacketLossKey: 0.06}),  // latency 6ms jitter 6us loss 6%
		9:  graph.NewNetworkEdge(9, nodes[6], nodes[8], map[helper.WeightKey]float64{helper.DefaultKey: 1, helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 0.01}),  // latency 1ms jitter 1us loss 1%
		10: graph.NewNetworkEdge(10, nodes[7], nodes[6], map[helper.WeightKey]float64{helper.DefaultKey: 1, helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 0.01}), // latency 1ms jitter 1us loss 1%
		11: graph.NewNetworkEdge(11, nodes[7], nodes[8], map[helper.WeightKey]float64{helper.DefaultKey: 5, helper.LatencyKey: 5000, helper.JitterKey: 50, helper.PacketLossKey: 0.05}), // latency 5ms jitter 5us loss 5%
	}
	type Result struct {
		edgeNumbers []int
		totalCost   float64
	}
	type args struct {
		from            graph.Node
		to              graph.Node
		weightType      helper.WeightKey
		calculationType CalculationType
	}
	tests := []struct {
		name    string
		args    args
		want    Result
		wantErr bool
	}{
		{
			name: "Test correct shortest path with default metric",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightType:      helper.DefaultKey,
				calculationType: CalculationTypeSum,
			},
			want: Result{
				edgeNumbers: []int{3, 7, 10, 9},
				totalCost:   4,
			},
			wantErr: false,
		},
		{
			name: "Test correct shortest path with latency metric",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightType:      helper.LatencyKey,
				calculationType: CalculationTypeSum,
			},
			want: Result{
				edgeNumbers: []int{3, 7, 10, 9},
				totalCost:   1000 + 1000 + 1000 + 1000, // latency of edge 3 + edge 7 + edge 10 + edge 9
			},
			wantErr: false,
		},
		{
			name: "Test correct shortest path with jitter metric",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightType:      helper.JitterKey,
				calculationType: CalculationTypeSum,
			},
			want: Result{
				edgeNumbers: []int{3, 7, 10, 9},
				totalCost:   10 + 10 + 10 + 10, // jitter of edge 3 + edge 7 + edge 10 + edge 9
			},
			wantErr: false,
		},
		{
			name: "Test correct shortest path with packet loss metric",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightType:      helper.PacketLossKey,
				calculationType: CalculationTypeSum,
			},
			want: Result{
				edgeNumbers: []int{3, 7, 10, 9},
				totalCost:   1 - (1-0.01)*(1-0.01)*(1-0.01)*(1-0.01), // ~0.04%, packet loss of edge 3 * edge 7 * edge 10 * edge 9 -> 1% on each link gives in the end 0.04% loss
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			networkGraph, err := setupGraph(nodes, edges)
			if err != nil {
				t.Errorf("Error setting up graph")
			}
			calculation := NewShortestPathCalculation(networkGraph, tt.args.from, tt.args.to, helper.WeightKey(tt.args.weightType), CalculationTypeSum)
			got, err := calculation.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Get shortest path with default metric %s, error = %v, wantErr %v", tt.args.weightType, err, tt.wantErr)
				return
			}
			shortestPath := make([]graph.Edge, len(tt.want.edgeNumbers))
			for index, node := range tt.want.edgeNumbers {
				shortestPath[index] = edges[node]
			}
			if !reflect.DeepEqual(got.GetEdges(), shortestPath) {
				t.Errorf("DefaultGraph.GetShortestPath() = %v, want %v", got.GetEdges(), shortestPath)
			}
			if !almostEqual(got.GetTotalCost(), tt.want.totalCost) {
				t.Errorf("DefaultGraph.GetShortestPath() = %v, want %v", got.GetTotalCost(), tt.want.totalCost)
			}
		})
	}
}
