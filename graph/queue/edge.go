package queue

import (
	"container/heap"
	"github.com/lantern-db/lantern/errors"
	m "github.com/lantern-db/lantern/graph/model"
)

type EdgeHeap []m.Edge

func (e EdgeHeap) Len() int {
	return len(e)
}

func (e EdgeHeap) Less(i, j int) bool {
	return e[i].Weight() > e[j].Weight()
}

func (e EdgeHeap) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e *EdgeHeap) Push(x interface{}) {
	*e = append(*e, x.(m.Edge))
}

func (e *EdgeHeap) Pop() interface{} {
	old := *e
	n := len(old)
	x := old[n-1]
	*e = old[0 : n-1]
	return x
}

type EdgePriorityQueue struct {
	edgeHeap *EdgeHeap
}

func NewEmptyPriorityQueue() EdgePriorityQueue {
	return EdgePriorityQueue{
		edgeHeap: &EdgeHeap{},
	}
}

func (q EdgePriorityQueue) Push(edge m.Edge) {
	heap.Push(q.edgeHeap, edge)
}

func (q EdgePriorityQueue) Pop() (m.Edge, error) {
	if len(*q.edgeHeap) > 0 {
		return heap.Pop(q.edgeHeap).(m.Edge), nil
	}
	return nil, errors.PriorityQueueEmptyError
}

func (q EdgePriorityQueue) Top(k uint32) []m.Edge {
	var edges []m.Edge
	for i := uint32(0); i < k; i++ {
		edge, err := q.Pop()
		if err != nil {
			break
		}
		edges = append(edges, edge)
	}
	return edges
}
