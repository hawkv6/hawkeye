package calculation

import (
	"container/heap"
)

// Modified from https://pkg.go.dev/container/heap#example-package-PriorityQueue

type Item struct {
	nodeId interface{}
	cost   float64
	index  int
}

func (item *Item) GetNodeId() interface{} {
	return item.nodeId
}

type Comparator func(i, j *Item) bool

type PriorityQueue struct {
	items      []*Item
	comparator Comparator
}

func NewMinimumPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		items: make([]*Item, 0),
		comparator: func(i, j *Item) bool {
			return i.cost < j.cost
		},
	}
}

func NewMaximumPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		items: make([]*Item, 0),
		comparator: func(i, j *Item) bool {
			return i.cost > j.cost
		},
	}
}

func (pq PriorityQueue) Len() int { return len(pq.items) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq.comparator(pq.items[i], pq.items[j])
}

func (pq PriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(pq.items)
	item := x.(*Item)
	item.index = n
	pq.items = append(pq.items, item)
}

func (pq *PriorityQueue) Pop() any {
	old := pq.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	pq.items = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) Update(item *Item, nodeId int, distance float64) {
	item.nodeId = nodeId
	item.cost = distance
	heap.Fix(pq, item.index)
}

func (pq *PriorityQueue) IsEmpty() bool {
	return len(pq.items) == 0
}

func (pq *PriorityQueue) Contains(nodeId interface{}) bool {
	for _, item := range pq.items {
		if item.nodeId == nodeId {
			return true
		}
	}
	return false
}

func (pq *PriorityQueue) GetIndex(nodeId interface{}) (int, bool) {
	for i, item := range pq.items {
		if item.nodeId == nodeId {
			return i, true
		}
	}
	return -1, false
}
