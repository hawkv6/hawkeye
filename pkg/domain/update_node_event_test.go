package domain

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestNewUpdateNodeEvent(t *testing.T) {
	tests := []struct {
		name        string
		key         *string
		igpRouterId *string
		nodeName    *string
		srAlgorithm []uint32
	}{
		{
			name:        "Add Node XR-4",
			key:         proto.String("2_0_0_0000.0000.0004"),
			igpRouterId: proto.String("0000.0000.0004"),
			nodeName:    proto.String("XR-4"),
			srAlgorithm: []uint32{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := NewDomainNode(
				tt.key,
				tt.igpRouterId,
				tt.nodeName,
				tt.srAlgorithm,
			)
			if err != nil {
				t.Error(err)
			}
			event := NewUpdateNodeEvent(node)
			assert.NotNil(t, event)
		})
	}
}
func TestUpdateNodeEvent_GetKey(t *testing.T) {
	tests := []struct {
		name        string
		key         *string
		igpRouterId *string
		nodeName    *string
		srAlgorihm  []uint32
	}{
		{
			name:        "Get Key of XR-4",
			key:         proto.String("2_0_0_0000.0000.0004"),
			igpRouterId: proto.String("0000.0000.0004"),
			nodeName:    proto.String("XR-4"),
			srAlgorihm:  []uint32{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := NewDomainNode(
				tt.key,
				tt.igpRouterId,
				tt.nodeName,
				tt.srAlgorihm,
			)
			if err != nil {
				t.Error(err)
			}
			event := &UpdateNodeEvent{
				Node: node,
			}
			actual := event.GetKey()
			reflect.DeepEqual(actual, *tt.key)
		})
	}
}
