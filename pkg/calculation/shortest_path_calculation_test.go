package calculation

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/stretchr/testify/assert"
)

func setupGraph(nodes map[int]graph.Node, edges map[int]graph.Edge) (*graph.NetworkGraph, error) {
	graph := graph.NewNetworkGraph()
	for _, node := range nodes {
		graph.AddNode(node)
	}
	for _, edge := range edges {
		if err := graph.AddEdge(edge); err != nil {
			return nil, err
		}
	}
	return graph, nil
}

const tolerance = 1e-6

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestShortestPathCalculation_Execute_SingleIntentSum(t *testing.T) {
	srAlgorithm := []uint32{0}
	nodes := map[int]graph.Node{
		1: graph.NewNetworkNode("1", "1", srAlgorithm),
		2: graph.NewNetworkNode("2", "2", srAlgorithm),
		3: graph.NewNetworkNode("3", "3", srAlgorithm),
		4: graph.NewNetworkNode("4", "4", srAlgorithm),
		5: graph.NewNetworkNode("5", "5", srAlgorithm),
		6: graph.NewNetworkNode("6", "6", srAlgorithm),
		7: graph.NewNetworkNode("7", "7", srAlgorithm),
		8: graph.NewNetworkNode("8", "8", srAlgorithm),
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
		1:  graph.NewNetworkEdge("1", nodes[1], nodes[2], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1, helper.UtilizedBandwidthKey: 1000, helper.AvailableBandwidthKey: 999000}),  // 1->2 latency 1ms, jitter 1us, loss 1%, utilization 1Mbit/s, available bandwidth 999Mbit/s
		3:  graph.NewNetworkEdge("3", nodes[1], nodes[4], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1, helper.UtilizedBandwidthKey: 1000, helper.AvailableBandwidthKey: 999000}),  // 1->4 latency 1ms, jitter 1us, loss 1%, utilization 1Mbit/s, available bandwidth 999Mbit/s
		2:  graph.NewNetworkEdge("2", nodes[1], nodes[3], map[helper.WeightKey]float64{helper.LatencyKey: 2000, helper.JitterKey: 20, helper.PacketLossKey: 2, helper.UtilizedBandwidthKey: 9000, helper.AvailableBandwidthKey: 991000}),  // 1->3 latency 2ms, jitter 2us, loss 2%, utilization 9Mbit/s, available bandwidth 991Mbit/s
		4:  graph.NewNetworkEdge("4", nodes[2], nodes[5], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1, helper.UtilizedBandwidthKey: 10000, helper.AvailableBandwidthKey: 990000}), // 2->5,latency 1ms, jitter 1us, loss 1%, utilization 10Mbit/s, available bandwidth 990Mbit/s
		5:  graph.NewNetworkEdge("5", nodes[3], nodes[5], map[helper.WeightKey]float64{helper.LatencyKey: 3000, helper.JitterKey: 30, helper.PacketLossKey: 3, helper.UtilizedBandwidthKey: 5000, helper.AvailableBandwidthKey: 995000}),  // 3->5 latency 3ms, jitter 3us, loss 3%, utilization 5Mbit/s, available bandwidth 995Mbit/s
		6:  graph.NewNetworkEdge("6", nodes[3], nodes[6], map[helper.WeightKey]float64{helper.LatencyKey: 4000, helper.JitterKey: 40, helper.PacketLossKey: 4, helper.UtilizedBandwidthKey: 5000, helper.AvailableBandwidthKey: 995000}),  // 3->6 latency 4ms, jitter 4us, loss 4%, utilization 5Mbit/s, available bandwidth 995Mbit/s
		7:  graph.NewNetworkEdge("7", nodes[4], nodes[7], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1, helper.UtilizedBandwidthKey: 1000, helper.AvailableBandwidthKey: 999000}),  // 4->7 latency 1ms, jitter 1us, loss 1%, utilization 1Mbit/s, available bandwidth 999Mbit/s
		8:  graph.NewNetworkEdge("8", nodes[5], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 6000, helper.JitterKey: 60, helper.PacketLossKey: 6, helper.UtilizedBandwidthKey: 1000, helper.AvailableBandwidthKey: 999000}),  // 5->8 latency 6ms, jitter 6us, loss 6%, utilization 1Mbit/s, available bandwidth 999Mbit/s
		9:  graph.NewNetworkEdge("9", nodes[6], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1, helper.UtilizedBandwidthKey: 1000, helper.AvailableBandwidthKey: 999000}),  // 6->8 latency 1ms, jitter 1us, loss 1%, utilization 1Mbit/s, available bandwidth 999Mbit/s
		10: graph.NewNetworkEdge("10", nodes[7], nodes[6], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1, helper.UtilizedBandwidthKey: 2000, helper.AvailableBandwidthKey: 998000}), // 7->6 latency 1ms, jitter 1us, loss 1%, utilization 2Mbit/s, available bandwidth 998Mbit/s
		11: graph.NewNetworkEdge("11", nodes[7], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 5000, helper.JitterKey: 50, helper.PacketLossKey: 5, helper.UtilizedBandwidthKey: 5000, helper.AvailableBandwidthKey: 995000}), // 7->8 latency 5ms, jitter 5us, loss 5%, utilization 5Mbit/s, available bandwidth 995Mbit/s
	}

	type Result struct {
		edgeNumbers []int
		totalCost   float64
	}
	type args struct {
		from            graph.Node
		to              graph.Node
		weightTypes     []helper.WeightKey
		calculationType CalculationMode
		minConstraints  map[helper.WeightKey]float64
		maxConstraints  map[helper.WeightKey]float64
	}
	tests := []struct {
		name    string
		args    args
		want    Result
		wantErr bool
	}{
		{
			name: "Test shortest path intent type low latency",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightTypes:     []helper.WeightKey{helper.LatencyKey},
				calculationType: CalculationModeSum,
				minConstraints:  make(map[helper.WeightKey]float64),
				maxConstraints:  make(map[helper.WeightKey]float64),
			},
			want: Result{
				edgeNumbers: []int{3, 7, 10, 9},
				totalCost:   1000 + 1000 + 1000 + 1000, // latency of edge 3 + edge 7 + edge 10 + edge 9
			},
			wantErr: false,
		},
		{
			name: "Test shortest path intent type low jitter",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightTypes:     []helper.WeightKey{helper.JitterKey},
				calculationType: CalculationModeSum,
				minConstraints:  make(map[helper.WeightKey]float64),
				maxConstraints:  make(map[helper.WeightKey]float64),
			},
			want: Result{
				edgeNumbers: []int{3, 7, 10, 9},
				totalCost:   10 + 10 + 10 + 10, // jitter of edge 3 + edge 7 + edge 10 + edge 9
			},
			wantErr: false,
		},
		{
			name: "Test shortest path intent type low packet loss",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightTypes:     []helper.WeightKey{helper.PacketLossKey},
				calculationType: CalculationModeSum,
				minConstraints:  make(map[helper.WeightKey]float64),
				maxConstraints:  make(map[helper.WeightKey]float64),
			},
			want: Result{
				edgeNumbers: []int{3, 7, 10, 9},
				totalCost:   4 * -math.Log(1-1/100.0), // we're using the formula -ln(1-p) to calculate the total loss
			},
			wantErr: false,
		},
		{
			name: "Test shortest path intent type low utilization",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightTypes:     []helper.WeightKey{helper.UtilizedBandwidthKey},
				calculationType: CalculationModeSum,
				minConstraints:  make(map[helper.WeightKey]float64),
				maxConstraints:  make(map[helper.WeightKey]float64),
			},
			want: Result{
				edgeNumbers: []int{3, 7, 10, 9},
				totalCost:   5000, // utilization of edge 3 + edge 7 + edge 10 + edge 9
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
			calculationOptions := &CalculationOptions{networkGraph, tt.args.from, tt.args.to, tt.args.weightTypes, tt.args.calculationType, tt.args.maxConstraints, tt.args.minConstraints}
			calculation := NewShortestPathCalculation(calculationOptions)
			got, err := calculation.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Get shortest path with metric %s, error = %v, wantErr %v", tt.args.weightTypes, err, tt.wantErr)
				return
			}
			shortestPath := make([]graph.Edge, len(tt.want.edgeNumbers))
			for index, node := range tt.want.edgeNumbers {
				shortestPath[index] = edges[node]
			}
			if !reflect.DeepEqual(got.GetEdges(), shortestPath) {
				t.Errorf("DefaultGraph.GetShortestPath() = %v, want %v", got.GetEdges(), shortestPath)
			} else {
				diagram := generatePlantUMLDiagram(nodes, got.GetEdges(), tt.name, got.GetTotalCost(), tt.args.weightTypes, got.GetBottleneckEdge())
				fmt.Println(diagram)
				err := os.WriteFile(fmt.Sprintf("../../test/uml/%s.uml", tt.name), []byte(diagram), 0644)
				if err != nil {
					t.Errorf("Failed to save PlantUML diagram: %v", err)
				}
			}
			if !almostEqual(got.GetTotalCost(), tt.want.totalCost) {
				t.Errorf("DefaultGraph.GetShortestPath() = %v, want %v", got.GetTotalCost(), tt.want.totalCost)
			}
		})
	}
}
func TestShortestPathCalculation_Execute_SingleIntentMax(t *testing.T) {
	srAlgorithm := []uint32{0}
	nodes := map[int]graph.Node{
		1: graph.NewNetworkNode("1", "1", srAlgorithm),
		2: graph.NewNetworkNode("2", "2", srAlgorithm),
		3: graph.NewNetworkNode("3", "3", srAlgorithm),
		4: graph.NewNetworkNode("4", "4", srAlgorithm),
		5: graph.NewNetworkNode("5", "5", srAlgorithm),
		6: graph.NewNetworkNode("6", "6", srAlgorithm),
		7: graph.NewNetworkNode("7", "7", srAlgorithm),
		8: graph.NewNetworkNode("8", "8", srAlgorithm),
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
		1:  graph.NewNetworkEdge("1", nodes[1], nodes[2], map[helper.WeightKey]float64{helper.AvailableBandwidthKey: 100000}),   // 1->2  bandwidth 100Mbps
		2:  graph.NewNetworkEdge("2", nodes[1], nodes[3], map[helper.WeightKey]float64{helper.AvailableBandwidthKey: 900000}),   // 1->3  bandwidth 900Mbps
		3:  graph.NewNetworkEdge("3", nodes[1], nodes[4], map[helper.WeightKey]float64{helper.AvailableBandwidthKey: 1000000}),  // 1->4  bandwidth 1Gbps
		4:  graph.NewNetworkEdge("4", nodes[2], nodes[5], map[helper.WeightKey]float64{helper.AvailableBandwidthKey: 1000000}),  // 2->5, bandwidth 1Gbps
		5:  graph.NewNetworkEdge("5", nodes[3], nodes[5], map[helper.WeightKey]float64{helper.AvailableBandwidthKey: 500000}),   // 3->5  bandwidth 500Mbps
		6:  graph.NewNetworkEdge("6", nodes[3], nodes[6], map[helper.WeightKey]float64{helper.AvailableBandwidthKey: 500000}),   // 3->6  bandwidth 500Mbps
		7:  graph.NewNetworkEdge("7", nodes[4], nodes[7], map[helper.WeightKey]float64{helper.AvailableBandwidthKey: 1000000}),  // 4->7  bandwidth 1Gbps
		8:  graph.NewNetworkEdge("8", nodes[5], nodes[8], map[helper.WeightKey]float64{helper.AvailableBandwidthKey: 1000000}),  // 5->8  bandwidth 1Gbps
		9:  graph.NewNetworkEdge("9", nodes[6], nodes[8], map[helper.WeightKey]float64{helper.AvailableBandwidthKey: 900000}),   // 6->8  bandwidth 900Mbps
		10: graph.NewNetworkEdge("10", nodes[7], nodes[6], map[helper.WeightKey]float64{helper.AvailableBandwidthKey: 1000000}), // 7->6  bandwidth 1Gbps
		11: graph.NewNetworkEdge("11", nodes[7], nodes[8], map[helper.WeightKey]float64{helper.AvailableBandwidthKey: 500000}),  // 7->8  bandwidth 500Mbps
	}

	type Result struct {
		edgeNumbers []int
		totalCost   float64
	}
	type args struct {
		from            graph.Node
		to              graph.Node
		weightTypes     []helper.WeightKey
		calculationType CalculationMode
		minConstraints  map[helper.WeightKey]float64
		maxConstraints  map[helper.WeightKey]float64
	}
	tests := []struct {
		name    string
		args    args
		want    Result
		wantErr bool
	}{
		{
			name: "Test shortest path intent type high bandwidth",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightTypes:     []helper.WeightKey{helper.AvailableBandwidthKey},
				calculationType: CalculationModeMax,
				minConstraints:  make(map[helper.WeightKey]float64),
				maxConstraints:  make(map[helper.WeightKey]float64),
			},
			want: Result{
				edgeNumbers: []int{3, 7, 10, 9},
				totalCost:   900000, // bottleneck value is 900 Mpbs,
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
			calculationOptions := &CalculationOptions{networkGraph, tt.args.from, tt.args.to, tt.args.weightTypes, tt.args.calculationType, tt.args.maxConstraints, tt.args.minConstraints}
			calculation := NewShortestPathCalculation(calculationOptions)
			got, err := calculation.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Get shortest path with metric %s, error = %v, wantErr %v", tt.args.weightTypes, err, tt.wantErr)
				return
			}
			shortestPath := make([]graph.Edge, len(tt.want.edgeNumbers))
			for index, node := range tt.want.edgeNumbers {
				shortestPath[index] = edges[node]
			}
			if !reflect.DeepEqual(got.GetEdges(), shortestPath) {
				t.Errorf("DefaultGraph.GetShortestPath() = %v, want %v", got.GetEdges(), shortestPath)
			} else {
				diagram := generatePlantUMLDiagram(nodes, got.GetEdges(), tt.name, got.GetTotalCost(), tt.args.weightTypes, got.GetBottleneckEdge())
				fmt.Println(diagram)
				err := os.WriteFile(fmt.Sprintf("../../test/uml/%s.uml", tt.name), []byte(diagram), 0644)
				if err != nil {
					t.Errorf("Failed to save PlantUML diagram: %v", err)
				}
			}
			if !almostEqual(got.GetTotalCost(), tt.want.totalCost) {
				t.Errorf("DefaultGraph.GetShortestPath() = %v, want %v", got.GetTotalCost(), tt.want.totalCost)
			}
		})
	}
}

func TestShortestPathCalculation_Execute_SingleIntentMin(t *testing.T) {
	srAlgorithm := []uint32{0}
	nodes := map[int]graph.Node{
		1: graph.NewNetworkNode("1", "1", srAlgorithm),
		2: graph.NewNetworkNode("2", "2", srAlgorithm),
		3: graph.NewNetworkNode("3", "3", srAlgorithm),
		4: graph.NewNetworkNode("4", "4", srAlgorithm),
		5: graph.NewNetworkNode("5", "5", srAlgorithm),
		6: graph.NewNetworkNode("6", "6", srAlgorithm),
		7: graph.NewNetworkNode("7", "7", srAlgorithm),
		8: graph.NewNetworkNode("8", "8", srAlgorithm),
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
		1:  graph.NewNetworkEdge("1", nodes[1], nodes[2], map[helper.WeightKey]float64{helper.MaximumLinkBandwidthKey: 100000, helper.AvailableBandwidthKey: 100000}),    // 1->2  bandwidth 100Mbps
		2:  graph.NewNetworkEdge("2", nodes[1], nodes[3], map[helper.WeightKey]float64{helper.MaximumLinkBandwidthKey: 900000, helper.AvailableBandwidthKey: 900000}),    // 1->3  bandwidth 900Mbps
		3:  graph.NewNetworkEdge("3", nodes[1], nodes[4], map[helper.WeightKey]float64{helper.MaximumLinkBandwidthKey: 1000000, helper.AvailableBandwidthKey: 1000000}),  // 1->4  bandwidth 1Gbps
		4:  graph.NewNetworkEdge("4", nodes[2], nodes[5], map[helper.WeightKey]float64{helper.MaximumLinkBandwidthKey: 1000000, helper.AvailableBandwidthKey: 1000000}),  // 2->5, bandwidth 1Gbps
		5:  graph.NewNetworkEdge("5", nodes[3], nodes[5], map[helper.WeightKey]float64{helper.MaximumLinkBandwidthKey: 500000, helper.AvailableBandwidthKey: 500000}),    // 3->5  bandwidth 500Mbps
		6:  graph.NewNetworkEdge("6", nodes[3], nodes[6], map[helper.WeightKey]float64{helper.MaximumLinkBandwidthKey: 500000, helper.AvailableBandwidthKey: 500000}),    // 3->6  bandwidth 500Mbps
		7:  graph.NewNetworkEdge("7", nodes[4], nodes[7], map[helper.WeightKey]float64{helper.MaximumLinkBandwidthKey: 1000000, helper.AvailableBandwidthKey: 1000000}),  // 4->7  bandwidth 1Gbps
		8:  graph.NewNetworkEdge("8", nodes[5], nodes[8], map[helper.WeightKey]float64{helper.MaximumLinkBandwidthKey: 1000000, helper.AvailableBandwidthKey: 1000000}),  // 5->8  bandwidth 1Gbps
		9:  graph.NewNetworkEdge("9", nodes[6], nodes[8], map[helper.WeightKey]float64{helper.MaximumLinkBandwidthKey: 900000, helper.AvailableBandwidthKey: 900000}),    // 6->8  bandwidth 900Mbps
		10: graph.NewNetworkEdge("10", nodes[7], nodes[6], map[helper.WeightKey]float64{helper.MaximumLinkBandwidthKey: 1000000, helper.AvailableBandwidthKey: 1000000}), // 7->6  bandwidth 1Gbps
		11: graph.NewNetworkEdge("11", nodes[7], nodes[8], map[helper.WeightKey]float64{helper.MaximumLinkBandwidthKey: 500000, helper.AvailableBandwidthKey: 500000}),   // 7->8  bandwidth 500Mbps
	}

	type Result struct {
		edgeNumbers []int
		totalCost   float64
	}
	type args struct {
		from            graph.Node
		to              graph.Node
		weightTypes     []helper.WeightKey
		calculationType CalculationMode
		minConstraints  map[helper.WeightKey]float64
		maxConstraints  map[helper.WeightKey]float64
	}
	tests := []struct {
		name    string
		args    args
		want    Result
		wantErr bool
	}{
		{
			name: "Test shortest path intent type low bandwidth",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightTypes:     []helper.WeightKey{helper.MaximumLinkBandwidthKey},
				calculationType: CalculationModeMin,
				minConstraints:  make(map[helper.WeightKey]float64),
				maxConstraints:  make(map[helper.WeightKey]float64),
			},
			want: Result{
				edgeNumbers: []int{1, 4, 8},
				totalCost:   100000, // bottleneck value is 100 Mpbs,
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
			calculationOptions := &CalculationOptions{networkGraph, tt.args.from, tt.args.to, tt.args.weightTypes, tt.args.calculationType, tt.args.maxConstraints, tt.args.minConstraints}
			calculation := NewShortestPathCalculation(calculationOptions)
			got, err := calculation.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Get shortest path with metric %s, error = %v, wantErr %v", tt.args.weightTypes, err, tt.wantErr)
				return
			}
			shortestPath := make([]graph.Edge, len(tt.want.edgeNumbers))
			for index, node := range tt.want.edgeNumbers {
				shortestPath[index] = edges[node]
			}
			if !reflect.DeepEqual(got.GetEdges(), shortestPath) {
				t.Errorf("DefaultGraph.GetShortestPath() = %v, want %v", got.GetEdges(), shortestPath)
			} else {
				diagram := generatePlantUMLDiagram(nodes, got.GetEdges(), tt.name, got.GetTotalCost(), tt.args.weightTypes, got.GetBottleneckEdge())
				fmt.Println(diagram)
				err := os.WriteFile(fmt.Sprintf("../../test/uml/%s.uml", tt.name), []byte(diagram), 0644)
				if err != nil {
					t.Errorf("Failed to save PlantUML diagram: %v", err)
				}
			}
			if !almostEqual(got.GetTotalCost(), tt.want.totalCost) {
				t.Errorf("DefaultGraph.GetShortestPath() = %v, want %v", got.GetTotalCost(), tt.want.totalCost)
			}
		})
	}
}

func TestShortestPathCalculation_Execute_MultipleIntents(t *testing.T) {
	srAlgorithm := []uint32{0}
	nodes := map[int]graph.Node{
		1: graph.NewNetworkNode("1", "1", srAlgorithm),
		2: graph.NewNetworkNode("2", "2", srAlgorithm),
		3: graph.NewNetworkNode("3", "3", srAlgorithm),
		4: graph.NewNetworkNode("4", "4", srAlgorithm),
		5: graph.NewNetworkNode("5", "5", srAlgorithm),
		6: graph.NewNetworkNode("6", "6", srAlgorithm),
		7: graph.NewNetworkNode("7", "7", srAlgorithm),
		8: graph.NewNetworkNode("8", "8", srAlgorithm),
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
		1:  graph.NewNetworkEdge("1", nodes[1], nodes[2], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.NormalizedLatencyKey: 0.1, helper.JitterKey: 10, helper.NormalizedJitterKey: 0.1, helper.PacketLossKey: 2, helper.NormalizedPacketLossKey: 0.2, helper.AvailableBandwidthKey: 100000}),  // 1->2 latency 1ms, normalized latency 0.1, jitter 10us, normalized jitter 0.1, loss 2%, normalized loss 0.2, available bandwidth 100Mbit/s
		2:  graph.NewNetworkEdge("2", nodes[1], nodes[3], map[helper.WeightKey]float64{helper.LatencyKey: 2000, helper.NormalizedLatencyKey: 0.2, helper.JitterKey: 20, helper.NormalizedJitterKey: 0.2, helper.PacketLossKey: 1, helper.NormalizedPacketLossKey: 0.1, helper.AvailableBandwidthKey: 1000000}), // 1->3 latency 2ms, normalized latency 0.2, jitter 20us, normalized jitter 0.2, loss 1%, normalized loss 0.1, available bandwidth 1Gbit/s
		3:  graph.NewNetworkEdge("3", nodes[1], nodes[4], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.NormalizedLatencyKey: 0.1, helper.JitterKey: 10, helper.NormalizedJitterKey: 0.1, helper.PacketLossKey: 2, helper.NormalizedPacketLossKey: 0.2, helper.AvailableBandwidthKey: 1000000}), // 1->4 latency 1ms, normalized latency 0.1, jitter 10us, normalized jitter 0.1, loss 2%, normalized loss 0.2, available bandwidth 1Gbit/s
		4:  graph.NewNetworkEdge("4", nodes[2], nodes[5], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.NormalizedLatencyKey: 0.1, helper.JitterKey: 10, helper.NormalizedJitterKey: 0.1, helper.PacketLossKey: 1, helper.NormalizedPacketLossKey: 0.1, helper.AvailableBandwidthKey: 100000}),  // 2->5 latency 1ms, normalized latency 0.1, jitter 10us, normalized jitter 0.1, loss 1%, normalized loss 0.1, available bandwidth 100Mbit/s
		5:  graph.NewNetworkEdge("5", nodes[3], nodes[5], map[helper.WeightKey]float64{helper.LatencyKey: 3000, helper.NormalizedLatencyKey: 0.3, helper.JitterKey: 30, helper.NormalizedJitterKey: 0.3, helper.PacketLossKey: 6, helper.NormalizedPacketLossKey: 0.6, helper.AvailableBandwidthKey: 1000000}), // 3->5 latency 3ms, normalized latency 0.3, jitter 30us, normalized jitter 0.3, loss 6%, normalized loss 0.6, available bandwidth 1Gbit/s
		6:  graph.NewNetworkEdge("6", nodes[3], nodes[6], map[helper.WeightKey]float64{helper.LatencyKey: 4000, helper.NormalizedLatencyKey: 0.4, helper.JitterKey: 40, helper.NormalizedJitterKey: 0.4, helper.PacketLossKey: 5, helper.NormalizedPacketLossKey: 0.5, helper.AvailableBandwidthKey: 900000}),  // 3->6 latency 4ms, normalized latency 0.4, jitter 40us, normalized jitter 0.4, loss 5%, normalized loss 0.5, available bandwidth 900Mbit/s
		7:  graph.NewNetworkEdge("7", nodes[4], nodes[7], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.NormalizedLatencyKey: 0.1, helper.JitterKey: 10, helper.NormalizedJitterKey: 0.1, helper.PacketLossKey: 2, helper.NormalizedPacketLossKey: 0.2, helper.AvailableBandwidthKey: 900000}),  // 4->7 latency 1ms, normalized latency 0.1, jitter 10us, normalized jitter 0.1, loss 2%, normalized loss 0.2, available bandwidth 900Mbit/s
		8:  graph.NewNetworkEdge("8", nodes[5], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 6000, helper.NormalizedLatencyKey: 0.6, helper.JitterKey: 60, helper.NormalizedJitterKey: 0.6, helper.PacketLossKey: 1, helper.NormalizedPacketLossKey: 0.1, helper.AvailableBandwidthKey: 100000}),  // 5->8 latency 6ms, normalized latency 0.6, jitter 60us, normalized jitter 0.6, loss 1%, normalized loss 0.1, available bandwidth 100Mbit/s
		9:  graph.NewNetworkEdge("9", nodes[6], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.NormalizedLatencyKey: 0.1, helper.JitterKey: 10, helper.NormalizedJitterKey: 0.1, helper.PacketLossKey: 1, helper.NormalizedPacketLossKey: 0.1, helper.AvailableBandwidthKey: 900000}),  // 6->8 latency 1ms, normalized latency 0.1, jitter 10us, normalized jitter 0.1, loss 1%, normalized loss 0.1, available bandwidth 900Mbit/s
		10: graph.NewNetworkEdge("10", nodes[7], nodes[6], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.NormalizedLatencyKey: 0.1, helper.JitterKey: 10, helper.NormalizedJitterKey: 0.1, helper.PacketLossKey: 9, helper.NormalizedPacketLossKey: 0.9, helper.AvailableBandwidthKey: 10000}),  // 7->6 latency 1ms, normalized latency 0.1, jitter 10us, normalized jitter 0.1, loss 9%, normalized loss 0.9, available bandwidth 10Mbit/s
		11: graph.NewNetworkEdge("11", nodes[7], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 5000, helper.NormalizedLatencyKey: 0.5, helper.JitterKey: 50, helper.NormalizedJitterKey: 0.5, helper.PacketLossKey: 10, helper.NormalizedPacketLossKey: 1, helper.AvailableBandwidthKey: 995000}),  // 7->8 latency 5ms, normalized latency 0.5, jitter 50us, normalized jitter 0.5, loss 10%, normalized loss 1, available bandwidth 995Mbit/s
	}

	type Result struct {
		edgeNumbers []int
		totalCost   float64
	}
	type args struct {
		from            graph.Node
		to              graph.Node
		weightTypes     []helper.WeightKey
		calculationType CalculationMode
		minConstraints  map[helper.WeightKey]float64
		maxConstraints  map[helper.WeightKey]float64
	}
	tests := []struct {
		name    string
		args    args
		want    Result
		wantErr bool
	}{
		{
			name: "Test shortest path intent type low latency and low packet loss no constraints",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightTypes:     []helper.WeightKey{helper.NormalizedLatencyKey, helper.NormalizedPacketLossKey},
				calculationType: CalculationModeSum,
				minConstraints:  make(map[helper.WeightKey]float64),
				maxConstraints:  make(map[helper.WeightKey]float64),
			},
			want: Result{
				edgeNumbers: []int{1, 4, 8},
				totalCost:   float64(helper.TwoFactorWeights[0]*0.1 + helper.TwoFactorWeights[1]*0.2 + helper.TwoFactorWeights[0]*0.1 + helper.TwoFactorWeights[1]*0.1 + helper.TwoFactorWeights[0]*0.6 + helper.TwoFactorWeights[1]*0.1), // calculation of costs with normalized values and default helper weights
			},
			wantErr: false,
		},
		{
			name: "Test shortest path intent type low packet loss and low delay no constraints",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightTypes:     []helper.WeightKey{helper.NormalizedPacketLossKey, helper.NormalizedLatencyKey},
				calculationType: CalculationModeSum,
				minConstraints:  make(map[helper.WeightKey]float64),
				maxConstraints:  make(map[helper.WeightKey]float64),
			},
			want: Result{
				edgeNumbers: []int{1, 4, 8},
				totalCost:   float64(helper.TwoFactorWeights[0]*0.2 + helper.TwoFactorWeights[1]*0.1 + helper.TwoFactorWeights[0]*0.1 + helper.TwoFactorWeights[1]*0.1 + helper.TwoFactorWeights[0]*0.1 + helper.TwoFactorWeights[1]*0.6), // calculation of costs with normalized values and default helper weights
			},
			wantErr: false,
		},
		{
			name: "Test shortest path intent type low packet loss with delay constraint 5000ms",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightTypes:     []helper.WeightKey{helper.NormalizedPacketLossKey},
				calculationType: CalculationModeSum,
				minConstraints:  make(map[helper.WeightKey]float64),
				maxConstraints:  map[helper.WeightKey]float64{helper.NormalizedLatencyKey: 5000},
			},
			want: Result{
				edgeNumbers: []int{3, 7, 10, 9},
				totalCost: float64(edges[3].GetWeight(helper.NormalizedPacketLossKey)) +
					edges[7].GetWeight(helper.NormalizedPacketLossKey) +
					edges[10].GetWeight(helper.NormalizedPacketLossKey) +
					edges[9].GetWeight(helper.NormalizedPacketLossKey), // calculation of costs with normalized values and default helper weights
			},
			wantErr: false,
		},
		{
			name: "Test no shortest path intent type low packet loss and low delay packet loss constraint 2%",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightTypes:     []helper.WeightKey{helper.NormalizedPacketLossKey, helper.NormalizedLatencyKey},
				calculationType: CalculationModeSum,
				minConstraints:  make(map[helper.WeightKey]float64),
				maxConstraints:  map[helper.WeightKey]float64{helper.NormalizedPacketLossKey: 0.02},
			},
			want: Result{
				edgeNumbers: []int{},
				totalCost:   0,
			},
			wantErr: true,
		},
		{
			name: "Test no shortest path intent type low packet loss and min bandwidth constraint 10Gbps",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightTypes:     []helper.WeightKey{helper.NormalizedPacketLossKey, helper.NormalizedLatencyKey},
				calculationType: CalculationModeSum,
				minConstraints:  map[helper.WeightKey]float64{helper.AvailableBandwidthKey: 10000000}, // no 10Gbps link available
				maxConstraints:  make(map[helper.WeightKey]float64),
			},
			want: Result{
				edgeNumbers: []int{},
				totalCost:   0,
			},
			wantErr: true,
		},
		{
			name: "Test shortest path intent type low packet loss low latency and low jitter and min bandwidth constraint 900Mbps",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightTypes:     []helper.WeightKey{helper.NormalizedPacketLossKey, helper.NormalizedLatencyKey, helper.NormalizedJitterKey},
				calculationType: CalculationModeSum,
				minConstraints:  map[helper.WeightKey]float64{helper.AvailableBandwidthKey: 900000}, // 900Mbps link available
				maxConstraints:  make(map[helper.WeightKey]float64),
			},
			want: Result{
				edgeNumbers: []int{2, 6, 9},
				totalCost: float64(helper.ThreeFactorWeights[0]*edges[2].GetWeight(helper.NormalizedPacketLossKey) +
					helper.ThreeFactorWeights[1]*edges[2].GetWeight(helper.NormalizedLatencyKey) +
					helper.ThreeFactorWeights[2]*edges[2].GetWeight(helper.NormalizedJitterKey) +
					helper.ThreeFactorWeights[0]*edges[6].GetWeight(helper.NormalizedPacketLossKey) +
					helper.ThreeFactorWeights[1]*edges[6].GetWeight(helper.NormalizedLatencyKey) +
					helper.ThreeFactorWeights[2]*edges[6].GetWeight(helper.NormalizedJitterKey) +
					helper.ThreeFactorWeights[0]*edges[9].GetWeight(helper.NormalizedPacketLossKey) +
					helper.ThreeFactorWeights[1]*edges[9].GetWeight(helper.NormalizedLatencyKey) +
					helper.ThreeFactorWeights[2]*edges[9].GetWeight(helper.NormalizedJitterKey)), // calculation of costs with normalized values and default helper weights
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
			calculationOptions := &CalculationOptions{networkGraph, tt.args.from, tt.args.to, tt.args.weightTypes, tt.args.calculationType, tt.args.maxConstraints, tt.args.minConstraints}
			calculation := NewShortestPathCalculation(calculationOptions)
			got, err := calculation.Execute()
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
			}
			shortestPath := make([]graph.Edge, len(tt.want.edgeNumbers))
			for index, node := range tt.want.edgeNumbers {
				shortestPath[index] = edges[node]
			}
			if !reflect.DeepEqual(got.GetEdges(), shortestPath) {
				for _, edge := range got.GetEdges() {
					fmt.Println(edge.GetId())
				}
				t.Errorf("DefaultGraph.GetShortestPath() = %v, want %v", got.GetEdges(), shortestPath)
			} else {
				diagram := generatePlantUMLDiagram(nodes, got.GetEdges(), tt.name, got.GetTotalCost(), tt.args.weightTypes, got.GetBottleneckEdge())
				fmt.Println(diagram)
				err := os.WriteFile(fmt.Sprintf("../../test/uml/%s.uml", tt.name), []byte(diagram), 0644)
				if err != nil {
					t.Errorf("Failed to save PlantUML diagram: %v", err)
				}
			}

			if !almostEqual(got.GetTotalCost(), tt.want.totalCost) {
				t.Errorf("DefaultGraph.GetShortestPath() = %v, want %v", got.GetTotalCost(), tt.want.totalCost)
			}
		})
	}
}

func weightKeyExists(key helper.WeightKey, list []helper.WeightKey) bool {
	for _, item := range list {
		if item == key {
			return true
		}
	}
	return false
}

func generatePlantUMLDiagram(nodes map[int]graph.Node, shortestPath []graph.Edge, title string, totalCost float64, weightTypes []helper.WeightKey, bottleNeckEdge graph.Edge) string {
	var builder strings.Builder

	builder.WriteString("@startuml\n")

	builder.WriteString(fmt.Sprintf("title %s\n", title))
	builder.WriteString(fmt.Sprintf("caption Total cost: %f\n", totalCost))

	for _, node := range nodes {
		id := node.GetId()
		builder.WriteString(fmt.Sprintf("node \"%s\" as n%s\n", id, id))
	}

	shortestPathEdges := make(map[interface{}]bool)
	for _, edge := range shortestPath {
		shortestPathEdges[edge.GetId()] = true
	}

	for _, node := range nodes {
		for _, edge := range node.GetEdges() {
			weights := edge.GetAllWeights()
			from := edge.From().GetId()
			to := edge.To().GetId()

			var details []string

			details = append(details, fmt.Sprintf("Latency: %fms", weights[helper.LatencyKey]))
			details = append(details, fmt.Sprintf("Jitter: %fus", weights[helper.JitterKey]))
			details = append(details, fmt.Sprintf("Loss: %f%%", weights[helper.PacketLossKey]))
			details = append(details, fmt.Sprintf("Available Bandwidth: %fMbit/s", weights[helper.AvailableBandwidthKey]))

			if weightKeyExists(helper.UtilizedBandwidthKey, weightTypes) {
				details = append(details, fmt.Sprintf("Utilized Bandwidth: %fMbit/s", weights[helper.UtilizedBandwidthKey]))
			}

			detailsStr := strings.Join(details, "\\n")

			color := ""
			if shortestPathEdges[edge.GetId()] {
				color = " #green "
			}
			if edge == bottleNeckEdge {
				color = " #red "
			}

			builder.WriteString(fmt.Sprintf("n%s -- n%s%s : \"%s -> %s \\n%s\"\n", from, to, color, from, to, detailsStr))
		}
	}
	builder.WriteString("@enduml\n")

	return builder.String()
}
