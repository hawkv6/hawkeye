package calculation

import (
	"fmt"
	"os"
	reflect "reflect"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/stretchr/testify/assert"
)

func TestServiceFunctionChainCalculation_Execute(t *testing.T) {
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
		1:  graph.NewNetworkEdge("1", nodes[1], nodes[2], map[helper.WeightKey]float64{helper.IgpMetricKey: 10, helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1, helper.UtilizedBandwidthKey: 1000, helper.AvailableBandwidthKey: 999000}),  // 1->2 igp metric 10, latency 1ms, jitter 1us, loss 1%, utilization 1Mbit/s, available bandwidth 999Mbit/s
		2:  graph.NewNetworkEdge("2", nodes[1], nodes[3], map[helper.WeightKey]float64{helper.IgpMetricKey: 20, helper.LatencyKey: 2000, helper.JitterKey: 20, helper.PacketLossKey: 2, helper.UtilizedBandwidthKey: 9000, helper.AvailableBandwidthKey: 991000}),  // 1->3  igp metric 10, latency 2ms, jitter 2us, loss 2%, utilization 9Mbit/s, available bandwidth 991Mbit/s
		3:  graph.NewNetworkEdge("3", nodes[1], nodes[4], map[helper.WeightKey]float64{helper.IgpMetricKey: 20, helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1, helper.UtilizedBandwidthKey: 1000, helper.AvailableBandwidthKey: 999000}),  // 1->4  igp metric 10, latency 1ms, jitter 1us, loss 1%, utilization 1Mbit/s, available bandwidth 999Mbit/s
		4:  graph.NewNetworkEdge("4", nodes[2], nodes[5], map[helper.WeightKey]float64{helper.IgpMetricKey: 10, helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1, helper.UtilizedBandwidthKey: 10000, helper.AvailableBandwidthKey: 990000}), // 2->5, igp metric 10, latency 1ms, jitter 1us, loss 1%, utilization 10Mbit/s, available bandwidth 990Mbit/s
		5:  graph.NewNetworkEdge("5", nodes[3], nodes[5], map[helper.WeightKey]float64{helper.IgpMetricKey: 100, helper.LatencyKey: 3000, helper.JitterKey: 30, helper.PacketLossKey: 3, helper.UtilizedBandwidthKey: 5000, helper.AvailableBandwidthKey: 995000}), // 3->5  igp metric 10, latency 3ms, jitter 3us, loss 3%, utilization 5Mbit/s, available bandwidth 995Mbit/s
		6:  graph.NewNetworkEdge("6", nodes[3], nodes[6], map[helper.WeightKey]float64{helper.IgpMetricKey: 10, helper.LatencyKey: 4000, helper.JitterKey: 40, helper.PacketLossKey: 4, helper.UtilizedBandwidthKey: 5000, helper.AvailableBandwidthKey: 995000}),  // 3->6  igp metric 10, latency 4ms, jitter 4us, loss 4%, utilization 5Mbit/s, available bandwidth 995Mbit/s
		7:  graph.NewNetworkEdge("7", nodes[4], nodes[7], map[helper.WeightKey]float64{helper.IgpMetricKey: 10, helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1, helper.UtilizedBandwidthKey: 1000, helper.AvailableBandwidthKey: 999000}),  // 4->7  igp metric 10, latency 1ms, jitter 1us, loss 1%, utilization 1Mbit/s, available bandwidth 999Mbit/s
		8:  graph.NewNetworkEdge("8", nodes[5], nodes[8], map[helper.WeightKey]float64{helper.IgpMetricKey: 10, helper.LatencyKey: 6000, helper.JitterKey: 60, helper.PacketLossKey: 6, helper.UtilizedBandwidthKey: 1000, helper.AvailableBandwidthKey: 999000}),  // 5->8  igp metric 10, latency 6ms, jitter 6us, loss 6%, utilization 1Mbit/s, available bandwidth 999Mbit/s
		9:  graph.NewNetworkEdge("9", nodes[6], nodes[8], map[helper.WeightKey]float64{helper.IgpMetricKey: 20, helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1, helper.UtilizedBandwidthKey: 1000, helper.AvailableBandwidthKey: 999000}),  // 6->8  igp metric 10, latency 1ms, jitter 1us, loss 1%, utilization 1Mbit/s, available bandwidth 999Mbit/s
		10: graph.NewNetworkEdge("10", nodes[7], nodes[6], map[helper.WeightKey]float64{helper.IgpMetricKey: 30, helper.LatencyKey: 1000, helper.JitterKey: 10, helper.PacketLossKey: 1, helper.UtilizedBandwidthKey: 2000, helper.AvailableBandwidthKey: 998000}), // 7->6  igp metric 10, latency 1ms, jitter 1us, loss 1%, utilization 2Mbit/s, available bandwidth 998Mbit/s
		11: graph.NewNetworkEdge("11", nodes[7], nodes[8], map[helper.WeightKey]float64{helper.IgpMetricKey: 40, helper.LatencyKey: 5000, helper.JitterKey: 50, helper.PacketLossKey: 5, helper.UtilizedBandwidthKey: 5000, helper.AvailableBandwidthKey: 995000}), // 7->8  igp metric 20, latency 5ms, jitter 5us, loss 5%, utilization 5Mbit/s, available bandwidth 995Mbit/s
	}

	type Result struct {
		edgeNumbers []int
		totalCost   float64
	}
	type args struct {
		from                 graph.Node
		to                   graph.Node
		weightTypes          []helper.WeightKey
		calculationType      CalculationMode
		minConstraints       map[helper.WeightKey]float64
		maxConstraints       map[helper.WeightKey]float64
		serviceFunctionChain [][]string
		routerServiceMap     map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    Result
		wantErr bool
	}{
		{
			name: "Test sfc path igp metric",
			args: args{
				from:                 nodes[1],
				to:                   nodes[8],
				weightTypes:          []helper.WeightKey{helper.IgpMetricKey},
				calculationType:      CalculationModeSum,
				minConstraints:       make(map[helper.WeightKey]float64),
				maxConstraints:       make(map[helper.WeightKey]float64),
				serviceFunctionChain: [][]string{{"2", "5"}, {"2", "7"}, {"4", "7"}, {"4", "5"}},
				routerServiceMap:     map[string]string{"2": "2001:db8:f2::", "4": "2001:db8:f4::", "5": "2001:db8:f5::", "7": "2001:db8:f7::"},
			},
			want: Result{
				edgeNumbers: []int{1, 4, 8},
				totalCost:   edges[1].GetWeight(helper.IgpMetricKey) + edges[4].GetWeight(helper.IgpMetricKey) + edges[8].GetWeight(helper.IgpMetricKey),
			},
			wantErr: false,
		},
		{
			name: "Test sfc path igp metric non optimal order",
			args: args{
				from:                 nodes[1],
				to:                   nodes[8],
				weightTypes:          []helper.WeightKey{helper.IgpMetricKey},
				calculationType:      CalculationModeSum,
				minConstraints:       make(map[helper.WeightKey]float64),
				maxConstraints:       make(map[helper.WeightKey]float64),
				serviceFunctionChain: [][]string{{"4", "7"}, {"4", "5"}, {"2", "5"}, {"2", "7"}},
				routerServiceMap:     map[string]string{"2": "2001:db8:f2::", "4": "2001:db8:f4::", "5": "2001:db8:f5::", "7": "2001:db8:f7::"},
			},
			want: Result{
				edgeNumbers: []int{1, 4, 8},
				totalCost:   edges[1].GetWeight(helper.IgpMetricKey) + edges[4].GetWeight(helper.IgpMetricKey) + edges[8].GetWeight(helper.IgpMetricKey),
			},
			wantErr: false,
		},
		{
			name: "Test sfc path igp metric no sfc found",
			args: args{
				from:                 nodes[1],
				to:                   nodes[8],
				weightTypes:          []helper.WeightKey{helper.IgpMetricKey},
				calculationType:      CalculationModeSum,
				minConstraints:       make(map[helper.WeightKey]float64),
				maxConstraints:       make(map[helper.WeightKey]float64),
				serviceFunctionChain: [][]string{{"2", "7"}, {"4", "5"}},
				routerServiceMap:     map[string]string{"2": "2001:db8:f2::", "4": "2001:db8:f4::", "5": "2001:db8:f5::", "7": "2001:db8:f7::"},
			},
			want: Result{
				edgeNumbers: []int{},
				totalCost:   0,
			},
			wantErr: true,
		},
		{
			name: "Test sfc path low latency metric",
			args: args{
				from:                 nodes[1],
				to:                   nodes[8],
				weightTypes:          []helper.WeightKey{helper.LatencyKey},
				calculationType:      CalculationModeSum,
				minConstraints:       make(map[helper.WeightKey]float64),
				maxConstraints:       make(map[helper.WeightKey]float64),
				serviceFunctionChain: [][]string{{"2", "5"}, {"2", "7"}, {"4", "7"}, {"4", "5"}},
				routerServiceMap:     map[string]string{"2": "2001:db8:f2::", "4": "2001:db8:f4::", "5": "2001:db8:f5::", "7": "2001:db8:f7::"},
			},
			want: Result{
				edgeNumbers: []int{3, 7, 10, 9},
				totalCost:   edges[3].GetWeight(helper.LatencyKey) + edges[7].GetWeight(helper.LatencyKey) + edges[10].GetWeight(helper.LatencyKey) + edges[9].GetWeight(helper.LatencyKey),
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
			sfcCalculationOptions := &SfcCalculationOptions{tt.args.serviceFunctionChain, tt.args.routerServiceMap}
			calculationOptions := &CalculationOptions{networkGraph, tt.args.from, tt.args.to, tt.args.weightTypes, tt.args.calculationType, tt.args.maxConstraints, tt.args.minConstraints}
			calculation := NewServiceFunctionChainCalculation(calculationOptions, sfcCalculationOptions)
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
