package queue

import (
	"github.com/golang/mock/gomock"
	"github.com/lantern-db/lantern/graph/model"
	mock_model "github.com/lantern-db/lantern/graph/model/mock"
	"math/rand"
	"testing"
)

func TestEdgeHeap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := NewEmptyPriorityQueue()

	for i := 1; i <= 10; i++ {
		edge := mock_model.NewMockEdge(ctrl)
		edge.EXPECT().Weight().Return(model.Weight(rand.Intn(1000))).AnyTimes()
		q.Push(edge)
	}

	isFirst := true
	var previous model.Edge
	for _, edge := range q.Top(20) {
		if isFirst {
			isFirst = false
			previous = edge
			continue
		}
		t.Run("valid_case", func(t *testing.T) {
			if previous.Weight() < edge.Weight() {
				t.Errorf("not sorted")
			}
			previous = edge

		})
	}
}
