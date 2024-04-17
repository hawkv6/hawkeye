package graph

import (
	"reflect"
	"testing"
)

func TestDefaultGraph_GetShortestPath(t *testing.T) {
	type fields struct {
		nodes map[int]Node
	}
	type args struct {
		from       Node
		to         Node
		weightKind string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "Test correct shortest path",
			fields: fields{
				nodes: map[int]Node{
					1: NewNetworkNode(1),
					2: NewNetworkNode(2),
					3: NewNetworkNode(3),
					4: NewNetworkNode(4),
					5: NewNetworkNode(5),
					6: NewNetworkNode(6),
					7: NewNetworkNode(7),
					8: NewNetworkNode(8),
				},
			},
			args: args{
				from:       NewNetworkNode(1),
				to:         NewNetworkNode(8),
				weightKind: "default",
			},
			want:    []int{3, 7, 10, 9},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edges := map[int]Edge{
				1:  NewNetworkEdge(1, tt.fields.nodes[1], tt.fields.nodes[2], map[string]float64{"default": 1}),
				2:  NewNetworkEdge(2, tt.fields.nodes[1], tt.fields.nodes[3], map[string]float64{"default": 2}),
				3:  NewNetworkEdge(3, tt.fields.nodes[1], tt.fields.nodes[4], map[string]float64{"default": 1}),
				4:  NewNetworkEdge(4, tt.fields.nodes[2], tt.fields.nodes[5], map[string]float64{"default": 1}),
				5:  NewNetworkEdge(5, tt.fields.nodes[3], tt.fields.nodes[5], map[string]float64{"default": 3}),
				6:  NewNetworkEdge(6, tt.fields.nodes[3], tt.fields.nodes[6], map[string]float64{"default": 4}),
				7:  NewNetworkEdge(7, tt.fields.nodes[4], tt.fields.nodes[7], map[string]float64{"default": 1}),
				8:  NewNetworkEdge(8, tt.fields.nodes[5], tt.fields.nodes[8], map[string]float64{"default": 6}),
				9:  NewNetworkEdge(9, tt.fields.nodes[6], tt.fields.nodes[8], map[string]float64{"default": 1}),
				10: NewNetworkEdge(10, tt.fields.nodes[7], tt.fields.nodes[6], map[string]float64{"default": 1}),
				11: NewNetworkEdge(11, tt.fields.nodes[7], tt.fields.nodes[8], map[string]float64{"default": 5}),
			}
			graph := NewNetworkGraph()
			for _, node := range tt.fields.nodes {
				graph.AddNode(node)
			}
			for _, edge := range edges {
				graph.AddEdge(edge)
			}
			got, err := graph.GetShortestPath(tt.args.from, tt.args.to, tt.args.weightKind)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultGraph.GetShortestPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			shortestPath := make([]Edge, len(tt.want))
			for index, node := range tt.want {
				shortestPath[index] = edges[node]
			}
			if !reflect.DeepEqual(got, shortestPath) {
				t.Errorf("DefaultGraph.GetShortestPath() = %v, want %v", got, shortestPath)
			}
		})
	}
}
