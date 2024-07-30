package calculation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMinimumPriorityQueue(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewMinimumPriorityQueue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, NewMinimumPriorityQueue())
		})
	}
}

func TestNewMaximumPriorityQueue(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewMaximumPriorityQueue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, NewMaximumPriorityQueue())
		})
	}
}

func TestPriorityQueue_Len(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestPriorityQueue_Len",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewMinimumPriorityQueue()
			assert.Equal(t, 0, pq.Len())
			pq.Push(&Item{})
			assert.Equal(t, 1, pq.Len())
		})
	}
}

func TestPriorityQueue_Less(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestPriorityQueue_Less",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewMinimumPriorityQueue()
			pq.Push(&Item{cost: 1})
			pq.Push(&Item{cost: 2})
			assert.True(t, pq.Less(0, 1))
			assert.False(t, pq.Less(1, 0))
		})
	}
}

func TestPriorityQueue_Swap(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestPriorityQueue_Swap",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewMinimumPriorityQueue()
			pq.Push(&Item{cost: 1})
			pq.Push(&Item{cost: 2})
			pq.Swap(0, 1)
			assert.Equal(t, float64(2), pq.items[0].cost)
			assert.Equal(t, float64(1), pq.items[1].cost)
		})
	}
}

func TestPriorityQueue_Push(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestPriorityQueue_Push",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewMinimumPriorityQueue()
			pq.Push(&Item{cost: 1})
			assert.Equal(t, 1, pq.Len())
			pq.Push(&Item{cost: 2})
			assert.Equal(t, 2, pq.Len())
		})
	}
}

func TestPriorityQueue_Pop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestPriorityQueue_Pop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewMinimumPriorityQueue()
			pq.Push(&Item{cost: 1})
			pq.Push(&Item{cost: 2})
			assert.Equal(t, 2, pq.Len())
			pq.Pop()
			assert.Equal(t, 1, pq.Len())
			assert.Equal(t, float64(1), pq.items[0].cost)
		})
	}
}

func TestPriorityQueue_Update(t *testing.T) {
	tests := []struct {
		name   string
		pqType string
	}{
		{
			name:   "TestPriorityQueue_Update MinPriorityQueue",
			pqType: "min",
		},
		{
			name:   "TestPriorityQueue_Update MaxPriorityQueue",
			pqType: "max",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pq *PriorityQueue
			if tt.pqType == "min" {
				pq = NewMinimumPriorityQueue()
			} else {
				pq = NewMaximumPriorityQueue()
			}
			item := &Item{cost: 1, nodeId: "1"}
			pq.Push(item)
			assert.Equal(t, float64(1), pq.items[0].cost)
			item2 := &Item{cost: 2, nodeId: "2"}
			pq.Push(item2)
			assert.Equal(t, float64(2), pq.items[1].cost)
			pq.Update(item, "1", 3)
			if tt.pqType == "min" {
				assert.Equal(t, float64(2), pq.items[0].cost)
				assert.Equal(t, float64(3), pq.items[1].cost)
			} else {
				assert.Equal(t, float64(3), pq.items[0].cost)
				assert.Equal(t, float64(2), pq.items[1].cost)
			}
		})
	}
}

func TestPriorityQueue_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestPriorityQueue_IsEmpty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewMinimumPriorityQueue()
			assert.True(t, pq.IsEmpty())
			pq.Push(&Item{})
			assert.False(t, pq.IsEmpty())
		})
	}
}

func TestPriorityQueue_Contains(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestPriorityQueue_Contains",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewMinimumPriorityQueue()
			item := &Item{cost: 1, nodeId: "1"}
			pq.Push(item)
			assert.True(t, pq.Contains("1"))
			assert.False(t, pq.Contains("2"))
		})
	}
}

func TestPriorityQueue_GetIndex(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestPriorityQueue_GetIndex",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewMinimumPriorityQueue()
			item := &Item{cost: 1, nodeId: "1"}
			pq.Push(item)
			got, exists := pq.GetIndex("1")
			assert.True(t, exists)
			assert.Equal(t, 0, got)
			got, exists = pq.GetIndex("2")
			assert.False(t, exists)
			assert.Equal(t, -1, got)
		})
	}
}
