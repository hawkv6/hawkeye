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
					1: NewDefaultNode(1),
					2: NewDefaultNode(2),
					3: NewDefaultNode(3),
					4: NewDefaultNode(4),
					5: NewDefaultNode(5),
					6: NewDefaultNode(6),
					7: NewDefaultNode(7),
					8: NewDefaultNode(8),
				},
			},
			args: args{
				from:       NewDefaultNode(1),
				to:         NewDefaultNode(8),
				weightKind: "default",
			},
			want:    []int{3, 7, 10, 9},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edges := map[int]Edge{
				1:  NewDefaultEdge(1, tt.fields.nodes[1], tt.fields.nodes[2], map[string]float64{"default": 1}),
				2:  NewDefaultEdge(2, tt.fields.nodes[1], tt.fields.nodes[3], map[string]float64{"default": 2}),
				3:  NewDefaultEdge(3, tt.fields.nodes[1], tt.fields.nodes[4], map[string]float64{"default": 1}),
				4:  NewDefaultEdge(4, tt.fields.nodes[2], tt.fields.nodes[5], map[string]float64{"default": 1}),
				5:  NewDefaultEdge(5, tt.fields.nodes[3], tt.fields.nodes[5], map[string]float64{"default": 3}),
				6:  NewDefaultEdge(6, tt.fields.nodes[3], tt.fields.nodes[6], map[string]float64{"default": 4}),
				7:  NewDefaultEdge(7, tt.fields.nodes[4], tt.fields.nodes[7], map[string]float64{"default": 1}),
				8:  NewDefaultEdge(8, tt.fields.nodes[5], tt.fields.nodes[8], map[string]float64{"default": 6}),
				9:  NewDefaultEdge(9, tt.fields.nodes[6], tt.fields.nodes[8], map[string]float64{"default": 1}),
				10: NewDefaultEdge(10, tt.fields.nodes[7], tt.fields.nodes[6], map[string]float64{"default": 1}),
				11: NewDefaultEdge(11, tt.fields.nodes[7], tt.fields.nodes[8], map[string]float64{"default": 5}),
			}
			graph := NewDefaultGraph()
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
