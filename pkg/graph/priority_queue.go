package graph

import (
	"container/heap"
)

// Modified from https://pkg.go.dev/container/heap#example-package-PriorityQueue

type Item struct {
	nodeId   interface{}
	distance float64
	index    int
}

func (item *Item) GetNodeId() interface{} {
	return item.nodeId
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) Update(item *Item, nodeId int, distance float64) {
	item.nodeId = nodeId
	item.distance = distance
	heap.Fix(pq, item.index)
}

func (pq *PriorityQueue) IsEmpty() bool {
	return len(*pq) == 0
}
