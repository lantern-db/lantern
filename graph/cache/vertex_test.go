package cache

import (
	"github.com/golang/mock/gomock"
	"github.com/piroyoung/lanterne/graph/model"
	"testing"
	"time"
)

func TestVertexCache_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := model.NewMockVertex(ctrl)
	m.EXPECT().Key().Return("mock").AnyTimes()
	m.EXPECT().Value().Return("mock").AnyTimes()

	c := NewVertexCache(10 * time.Second)
	t.Run("valid_case", func(t *testing.T) {
		c.Set("mock", m)
		c.Delete("mock")
		got, found := c.Get("mock")
		if found {
			t.Errorf("Get() got = %v, want %v", got, nil)
		}
	})
}

func TestVertexCache_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := model.NewMockVertex(ctrl)
	m.EXPECT().Key().Return("mock").AnyTimes()

	v := NewVertexCache(3 * time.Second)
	t.Run("valid_case", func(t *testing.T) {
		v.Set("mock", m)
		got, found := v.Get("mock")
		if !found {
			t.Errorf("Get() got = %v, want %v", got, m)
		}
	})

	t.Run("time_out", func(t *testing.T) {
		time.Sleep(4 * time.Second)
		got, found := v.Get("mock")
		if found {
			t.Errorf("Get() got = %v, want %v", got, nil)
		}
	})
}

func TestVertexCache_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := model.NewMockVertex(ctrl)
	m.EXPECT().Key().Return("mock").AnyTimes()

	v := NewVertexCache(3 * time.Second)
	t.Run("valid_case", func(t *testing.T) {
		v.Set("mock", m)
		got, found := v.Get("mock")
		if !found {
			t.Errorf("Get() got = %v, want %v", got, m)
		}
	})
}
