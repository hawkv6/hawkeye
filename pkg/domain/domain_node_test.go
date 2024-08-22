package domain

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestNewDomainNode(t *testing.T) {
	tests := []struct {
		name        string
		key         *string
		igpRouterId *string
		nodeName    *string
		srAlgorithm []uint32
		wantErr     bool
		want        *DomainNode
	}{
		{
			name:        "Test NewDomainNode",
			key:         proto.String("2_0_0_0000.0000.0004"),
			igpRouterId: proto.String("0000.0000.0004"),
			nodeName:    proto.String("XR-4"),
			srAlgorithm: []uint32{0},
			wantErr:     false,
			want: &DomainNode{
				key:         "2_0_0_0000.0000.0004",
				igpRouterId: "0000.0000.0004",
				name:        "XR-4",
				srAlgorithm: []uint32{0},
			},
		},
		{
			name:        "Test NewDomainNode error no nodeName provided",
			key:         proto.String("2_0_0_0000.0000.0004"),
			igpRouterId: proto.String("0000.0000.0004"),
			nodeName:    nil,
			srAlgorithm: []uint32{0},
			wantErr:     true,
			want:        nil,
		},
	}

	for _, tt := range tests {
		node, err := NewDomainNode(
			tt.key,
			tt.igpRouterId,
			tt.nodeName,
			tt.srAlgorithm,
		)
		if (err != nil) != tt.wantErr {
			t.Errorf("Error creating DomainNode: %v", err)
			return
		}
		if !reflect.DeepEqual(node, tt.want) {
			t.Errorf("NewDomainNode() test '%s' = %v, want %v", tt.name, node, tt.want)
		}
	}
}

func TestDomainNode_GetKey(t *testing.T) {
	tests := []struct {
		name        string
		key         *string
		igpRouterId *string
		nodeName    *string
		srAlgorithm []uint32
	}{
		{
			name:        "Test GetKey",
			key:         proto.String("2_0_0_0000.0000.0004"),
			igpRouterId: proto.String("0000.0000.0004"),
			nodeName:    proto.String("XR-4"),
			srAlgorithm: []uint32{0},
		},
	}

	for _, tt := range tests {
		node, err := NewDomainNode(
			tt.key,
			tt.igpRouterId,
			tt.nodeName,
			tt.srAlgorithm,
		)
		if err != nil {
			t.Errorf("Error creating DomainNode: %v", err)
		}
		assert.Equal(t, *tt.key, node.GetKey())
	}
}

func TestDomainNode_GetIgpRouterId(t *testing.T) {
	tests := []struct {
		name        string
		key         *string
		igpRouterId *string
		nodeName    *string
		srAlgorithm []uint32
	}{
		{
			name:        "Test GetIgpRouterId",
			key:         proto.String("2_0_0_0000.0000.0004"),
			igpRouterId: proto.String("0000.0000.0004"),
			nodeName:    proto.String("XR-4"),
			srAlgorithm: []uint32{0},
		},
	}

	for _, tt := range tests {
		node, err := NewDomainNode(
			tt.key,
			tt.igpRouterId,
			tt.nodeName,
			tt.srAlgorithm,
		)
		if err != nil {
			t.Errorf("Error creating DomainNode: %v", err)
		}
		assert.Equal(t, *tt.igpRouterId, node.GetIgpRouterId())
	}
}

func TestDomainNode_GetName(t *testing.T) {
	tests := []struct {
		name        string
		key         *string
		igpRouterId *string
		nodeName    *string
		srAlgorithm []uint32
	}{
		{
			name:        "Test GetName",
			key:         proto.String("2_0_0_0000.0000.0004"),
			igpRouterId: proto.String("0000.0000.0004"),
			nodeName:    proto.String("XR-4"),
			srAlgorithm: []uint32{0},
		},
	}

	for _, tt := range tests {
		node, err := NewDomainNode(
			tt.key,
			tt.igpRouterId,
			tt.nodeName,
			tt.srAlgorithm,
		)
		if err != nil {
			t.Errorf("Error creating DomainNode: %v", err)
		}
		assert.Equal(t, *tt.nodeName, node.GetName())
	}
}

func TestDomainNode_GetSrAlgorithm(t *testing.T) {
	tests := []struct {
		name        string
		key         *string
		igpRouterId *string
		nodeName    *string
		srAlgorithm []uint32
	}{
		{
			name:        "Test GetSrAlgorithm",
			key:         proto.String("2_0_0_0000.0000.0004"),
			igpRouterId: proto.String("0000.0000.0004"),
			nodeName:    proto.String("XR-4"),
			srAlgorithm: []uint32{0},
		},
	}

	for _, tt := range tests {
		node, err := NewDomainNode(
			tt.key,
			tt.igpRouterId,
			tt.nodeName,
			tt.srAlgorithm,
		)
		if err != nil {
			t.Errorf("Error creating DomainNode: %v", err)
		}
		assert.Equal(t, tt.srAlgorithm, node.GetSrAlgorithm())
	}
}
