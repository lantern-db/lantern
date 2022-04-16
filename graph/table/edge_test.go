package table

import (
	"github.com/golang/mock/gomock"
	pb "github.com/lantern-db/lantern-proto/go/lantern/v1"
	. "github.com/lantern-db/lantern/graph/model"
	mock_model "github.com/lantern-db/lantern/graph/model/mock"
	"testing"
	"time"
)

func TestEdgeTable_Append(t1 *testing.T) {
	ctrl := gomock.NewController(t1)
	defer ctrl.Finish()

	expiration := NewExpiration(30 * time.Second)

	e1 := mock_model.NewMockEdge(ctrl)
	e1.EXPECT().Tail().Return(Key("tail")).AnyTimes()
	e1.EXPECT().Head().Return(Key("head")).AnyTimes()
	e1.EXPECT().Expiration().Return(expiration).AnyTimes()
	e1.EXPECT().Weight().Return(Weight(1.0)).AnyTimes()
	e1.EXPECT().AsProto().Return(&pb.Edge{
		Tail:       "tail",
		Head:       "head",
		Weight:     1.0,
		Expiration: expiration.AsProtoTimestamp(),
	}).AnyTimes()

	type fields struct {
		edges []Edge
	}
	type args struct {
		edge Edge
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expect int
	}{
		{
			"valid_case",
			fields{[]Edge{}},
			args{e1},
			1,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &EdgeTable{
				edges: tt.fields.edges,
			}
			t.Append(tt.args.edge)
			if t.Len() != tt.expect {
				t1.Errorf("expect: %c, got: %c", tt.expect, t.Len())
			}
		})
	}
}

func TestEdgeTable_Flush(t1 *testing.T) {
	ctrl := gomock.NewController(t1)
	defer ctrl.Finish()

	expiration1 := NewExpiration(-1 * time.Second)
	e1 := mock_model.NewMockEdge(ctrl)
	e1.EXPECT().Tail().Return(Key("tail")).AnyTimes()
	e1.EXPECT().Head().Return(Key("head")).AnyTimes()
	e1.EXPECT().Expiration().Return(expiration1).AnyTimes()
	e1.EXPECT().Weight().Return(Weight(1.0)).AnyTimes()
	e1.EXPECT().AsProto().Return(&pb.Edge{
		Tail:       "tail",
		Head:       "head",
		Weight:     1.0,
		Expiration: expiration1.AsProtoTimestamp(),
	}).AnyTimes()

	expiration2 := NewExpiration(30 * time.Second)
	e2 := mock_model.NewMockEdge(ctrl)
	e2.EXPECT().Tail().Return(Key("tail")).AnyTimes()
	e2.EXPECT().Head().Return(Key("head")).AnyTimes()
	e2.EXPECT().Expiration().Return(expiration2).AnyTimes()
	e2.EXPECT().Weight().Return(Weight(2.0)).AnyTimes()
	e2.EXPECT().AsProto().Return(&pb.Edge{
		Tail:       "tail",
		Head:       "head",
		Weight:     2.0,
		Expiration: expiration2.AsProtoTimestamp(),
	}).AnyTimes()

	type fields struct {
		edges []Edge
	}
	tests := []struct {
		name   string
		fields fields
		expect int
	}{
		{
			"not_expired_case",
			fields{[]Edge{e2}},
			1,
		},
		{
			"expired_case",
			fields{[]Edge{e1}},
			0,
		},
		{
			"mix_case",
			fields{[]Edge{e1, e2}},
			1,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &EdgeTable{
				edges: tt.fields.edges,
			}
			t.sort()
			t.flush()
			if t.Len() != tt.expect {
				t1.Errorf("expect: %c, got: %c", tt.expect, t.Len())
			}
		})
	}
}

func TestEdgeTable_Sort(t1 *testing.T) {
	ctrl := gomock.NewController(t1)
	defer ctrl.Finish()

	expiration1 := NewExpiration(30 * time.Second)
	e1 := mock_model.NewMockEdge(ctrl)
	e1.EXPECT().Tail().Return(Key("tail")).AnyTimes()
	e1.EXPECT().Head().Return(Key("head")).AnyTimes()
	e1.EXPECT().Expiration().Return(expiration1).AnyTimes()
	e1.EXPECT().Weight().Return(Weight(1.0)).AnyTimes()
	e1.EXPECT().AsProto().Return(&pb.Edge{
		Tail:       "tail",
		Head:       "head",
		Weight:     1.0,
		Expiration: expiration1.AsProtoTimestamp(),
	}).AnyTimes()

	expiration2 := NewExpiration(31 * time.Second)
	e2 := mock_model.NewMockEdge(ctrl)
	e2.EXPECT().Tail().Return(Key("tail")).AnyTimes()
	e2.EXPECT().Head().Return(Key("head")).AnyTimes()
	e2.EXPECT().Expiration().Return(expiration2).AnyTimes()
	e2.EXPECT().Weight().Return(Weight(2.0)).AnyTimes()
	e2.EXPECT().AsProto().Return(&pb.Edge{
		Tail:       "tail",
		Head:       "head",
		Weight:     2.0,
		Expiration: expiration2.AsProtoTimestamp(),
	}).AnyTimes()

	type fields struct {
		edges []Edge
	}
	tests := []struct {
		name        string
		fields      fields
		expectedTop Edge
	}{
		{
			"valid_case",
			fields{[]Edge{e1, e2}},
			e2,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &EdgeTable{
				edges: tt.fields.edges,
			}
			t.sort()
			if t.edges[0] != tt.expectedTop {
				t1.Errorf("expect: %c, got: %c", tt.expectedTop, t.edges[0])
			}
		})
	}
}

func TestEdgeTable_Weight(t1 *testing.T) {
	ctrl := gomock.NewController(t1)
	defer ctrl.Finish()

	expiration1 := NewExpiration(-1 * time.Second)
	e1 := mock_model.NewMockEdge(ctrl)
	e1.EXPECT().Tail().Return(Key("tail")).AnyTimes()
	e1.EXPECT().Head().Return(Key("head")).AnyTimes()
	e1.EXPECT().Expiration().Return(expiration1).AnyTimes()
	e1.EXPECT().Weight().Return(Weight(1.0)).AnyTimes()
	e1.EXPECT().AsProto().Return(&pb.Edge{
		Tail:       "tail",
		Head:       "head",
		Weight:     1.0,
		Expiration: expiration1.AsProtoTimestamp(),
	}).AnyTimes()

	expiration2 := NewExpiration(30 * time.Second)
	e2 := mock_model.NewMockEdge(ctrl)
	e2.EXPECT().Tail().Return(Key("tail")).AnyTimes()
	e2.EXPECT().Head().Return(Key("head")).AnyTimes()
	e2.EXPECT().Expiration().Return(expiration2).AnyTimes()
	e2.EXPECT().Weight().Return(Weight(2.0)).AnyTimes()
	e2.EXPECT().AsProto().Return(&pb.Edge{
		Tail:       "tail",
		Head:       "head",
		Weight:     2.0,
		Expiration: expiration2.AsProtoTimestamp(),
	}).AnyTimes()

	type fields struct {
		edges []Edge
	}
	tests := []struct {
		name   string
		fields fields
		want   Weight
	}{
		{
			"valid_case",
			fields{[]Edge{e1, e2}},
			Weight(2.0),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &EdgeTable{
				edges: tt.fields.edges,
			}
			t.sort()
			if got := t.Weight(); got != tt.want {
				t1.Errorf("Weight() = %v, want %v", got, tt.want)
			}
		})
	}
}
