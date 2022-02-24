package cache

import (
	"github.com/golang/mock/gomock"
	. "github.com/lantern-db/lantern/graph/model"
	mock_model "github.com/lantern-db/lantern/graph/model/mock"
	"math"
	"testing"
	"time"
)

func TestEdgeCache_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	edge1 := mock_model.NewMockEdge(ctrl)
	edge1.EXPECT().Tail().Return(Key("tail")).AnyTimes()
	edge1.EXPECT().Head().Return(Key("head1")).AnyTimes()
	edge1.EXPECT().Weight().Return(Weight(1.0)).AnyTimes()
	edge1.EXPECT().Expiration().Return(NewExpiration(5 * time.Second)).AnyTimes()

	edge2 := mock_model.NewMockEdge(ctrl)
	edge2.EXPECT().Tail().Return(Key("tail")).AnyTimes()
	edge2.EXPECT().Head().Return(Key("head2")).AnyTimes()
	edge2.EXPECT().Weight().Return(Weight(1.0)).AnyTimes()
	edge2.EXPECT().Expiration().Return(NewExpiration(5 * time.Second)).AnyTimes()

	c := NewEdgeCache()
	t.Run("valid_case", func(t *testing.T) {
		c.Put(edge1)
		c.Put(edge2)
		if len(c.cache["tail"]) != 2 {
			t.Errorf("mismatch length")
		}
		c.Delete("tail", "head1")
		if len(c.cache["tail"]) != 1 {
			t.Errorf("mismatch_length")
		}
		c.Delete("tail", "head2")

		_, ok := c.cache["tail"]
		if ok {
			t.Errorf("keys tail shold be deleted")
		}
	})
}

func TestEdgeCache_GetAdjacent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	edge1 := mock_model.NewMockEdge(ctrl)
	edge1.EXPECT().Tail().Return(Key("tail")).AnyTimes()
	edge1.EXPECT().Head().Return(Key("head1")).AnyTimes()
	edge1.EXPECT().Weight().Return(Weight(1.0)).AnyTimes()
	edge1.EXPECT().Expiration().Return(NewExpiration(5 * time.Second)).AnyTimes()

	edge2 := mock_model.NewMockEdge(ctrl)
	edge2.EXPECT().Tail().Return(Key("tail")).AnyTimes()
	edge2.EXPECT().Head().Return(Key("head2")).AnyTimes()
	edge2.EXPECT().Weight().Return(Weight(1.0)).AnyTimes()
	edge2.EXPECT().Expiration().Return(NewExpiration(5 * time.Second)).AnyTimes()

	c := NewEdgeCache()

	t.Run("valid_case", func(t *testing.T) {
		c.Put(edge1)
		c.Put(edge2)
		if got, found := c.GetAdjacent("tail", math.MaxUint32); !found {
			t.Errorf("not found")
		} else {
			if len(got) != 2 {
				t.Errorf("mismatch length")
			}
		}
	})
}

func TestEdgeCache_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	edge := mock_model.NewMockEdge(ctrl)
	edge.EXPECT().Tail().Return(Key("tail")).AnyTimes()
	edge.EXPECT().Head().Return(Key("head")).AnyTimes()
	edge.EXPECT().Weight().Return(Weight(1.0)).AnyTimes()
	edge.EXPECT().Expiration().Return(NewExpiration(5 * time.Second)).AnyTimes()

	c := NewEdgeCache()
	t.Run("valid_case", func(t *testing.T) {
		c.Put(edge)
		got, found := c.Get("tail", "head")
		if !found {
			t.Errorf("not found")
		}
		if got.Weight() != 1.0 {
			t.Errorf("GetWeight() got = %v, want %v", got, 1.0)
		}
	})

}

func TestEdgeCache_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	edge := mock_model.NewMockEdge(ctrl)
	edge.EXPECT().Tail().Return(Key("tail")).AnyTimes()
	edge.EXPECT().Head().Return(Key("head")).AnyTimes()
	edge.EXPECT().Weight().Return(Weight(1.0)).AnyTimes()
	edge.EXPECT().Expiration().Return(NewExpiration(5 * time.Second)).AnyTimes()

	c := NewEdgeCache()
	t.Run("valid_case", func(t *testing.T) {
		c.Put(edge)
	})

	t.Run("not_expired", func(t *testing.T) {
		_, found := c.Get("tail", "head")
		if !found {
			t.Errorf("not found")
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("expired", func(t *testing.T) {
		_, found := c.Get("tail", "head")
		if found {
			t.Errorf("not expired")
		}
	})
}
