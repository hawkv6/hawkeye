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

const tolerance = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestNetworkGraph_GetShortestPathSingleIntent(t *testing.T) {
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
		1:  graph.NewNetworkEdge("1", nodes[1], nodes[2], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency 1ms jitter 1us loss 1%
		2:  graph.NewNetworkEdge("2", nodes[1], nodes[3], map[helper.WeightKey]float64{helper.LatencyKey: 2000, helper.JitterKey: 20, helper.PacketLossKey: 2}),  // latency 2ms jitter 2us loss 2%
		3:  graph.NewNetworkEdge("3", nodes[1], nodes[4], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency 1ms jitter 1us loss 1%
		4:  graph.NewNetworkEdge("4", nodes[2], nodes[5], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency 1ms jitter 1us loss 1%
		5:  graph.NewNetworkEdge("5", nodes[3], nodes[5], map[helper.WeightKey]float64{helper.LatencyKey: 3000, helper.JitterKey: 30, helper.PacketLossKey: 3}),  // latency 3ms jitter 3us loss 3%
		6:  graph.NewNetworkEdge("6", nodes[3], nodes[6], map[helper.WeightKey]float64{helper.LatencyKey: 4000, helper.JitterKey: 40, helper.PacketLossKey: 4}),  // latency 4ms jitter 4us loss 4%
		7:  graph.NewNetworkEdge("7", nodes[4], nodes[7], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency ms jitter 1us loss 1%
		8:  graph.NewNetworkEdge("8", nodes[5], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 6000, helper.JitterKey: 60, helper.PacketLossKey: 6}),  // latency 6ms jitter 6us loss 6%
		9:  graph.NewNetworkEdge("9", nodes[6], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}),  // latency 1ms jitter 1us loss 1%
		10: graph.NewNetworkEdge("10", nodes[7], nodes[6], map[helper.WeightKey]float64{helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1}), // latency 1ms jitter 1us loss 1%
		11: graph.NewNetworkEdge("11", nodes[7], nodes[8], map[helper.WeightKey]float64{helper.LatencyKey: 5000, helper.JitterKey: 50, helper.PacketLossKey: 5}), // latency 5ms jitter 5us loss 5%
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
	}
	tests := []struct {
		name    string
		args    args
		want    Result
		wantErr bool
	}{
		{
			name: "Test correct shortest path with latency metric",
			args: args{
				from:            nodes[1],
				to:              nodes[8],
				weightTypes:     []helper.WeightKey{helper.LatencyKey},
				calculationType: CalculationModeSum,
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
				weightTypes:     []helper.WeightKey{helper.JitterKey},
				calculationType: CalculationModeSum,
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
				weightTypes:     []helper.WeightKey{helper.PacketLossKey},
				calculationType: CalculationModeSum,
			},
			want: Result{
				edgeNumbers: []int{3, 7, 10, 9},
				// totalCost:   1 - ((1 - 1/100.0) * (1 - 1/100.0) * (1 - 1/100.0) * (1 - 1/100.0)), // ~0.04%, packet loss of edge 3 * edge 7 * edge 10 * edge 9 -> 1% on each link gives in the end 0.04% loss
				totalCost: 4 * -math.Log(1-1/100.0), // we're using the formula -ln(1-p) to calculate the total loss
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
			calculationOptions := &CalculationOptions{networkGraph, tt.args.from, tt.args.to, tt.args.weightTypes, tt.args.calculationType, make(map[helper.WeightKey]float64), make(map[helper.WeightKey]float64)}
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
				diagram := generatePlantUMLDiagram(nodes, got.GetEdges(), tt.name, got.GetTotalCost())
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

func generatePlantUMLDiagram(nodes map[int]graph.Node, shortestPath []graph.Edge, title string, totalCost float64) string {
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
			color := ""
			if shortestPathEdges[edge.GetId()] {
				color = " #red"
			}
			weights := edge.GetAllWeights()
			from := edge.From().GetId()
			to := edge.To().GetId()
			if color == "" {
				builder.WriteString(fmt.Sprintf("n%s -- n%s : \"%s -> %s \\nLatency: %fms\\nJitter: %fus\\nLoss: %f\"\n",
					from, to, from, to, weights[helper.LatencyKey], weights[helper.JitterKey], weights[helper.PacketLossKey]))
			} else {
				builder.WriteString(fmt.Sprintf("n%s -- n%s %s : \"%s -> %s \\nLatency: %fms\\nJitter: %fus\\nLoss: %f\" SPF \n",
					from, to, color, from, to, weights[helper.LatencyKey], weights[helper.JitterKey], weights[helper.PacketLossKey]))
			}
		}
	}
	builder.WriteString("@enduml\n")

	return builder.String()
}