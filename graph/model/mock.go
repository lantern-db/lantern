// Code generated by MockGen. DO NOT EDIT.
// Source: lanterne/model/graph.go

// Package model is a generated GoMock package.
package model

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockVertex is a mock of Vertex interface.
type MockVertex struct {
	ctrl     *gomock.Controller
	recorder *MockVertexMockRecorder
}

// MockVertexMockRecorder is the mock recorder for MockVertex.
type MockVertexMockRecorder struct {
	mock *MockVertex
}

// NewMockVertex creates a new mock instance.
func NewMockVertex(ctrl *gomock.Controller) *MockVertex {
	mock := &MockVertex{ctrl: ctrl}
	mock.recorder = &MockVertexMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVertex) EXPECT() *MockVertexMockRecorder {
	return m.recorder
}

// Digest mocks base method.
func (m *MockVertex) Digest() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Digest")
	ret0, _ := ret[0].(string)
	return ret0
}

// Digest indicates an expected call of Digest.
func (mr *MockVertexMockRecorder) Digest() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Digest", reflect.TypeOf((*MockVertex)(nil).Digest))
}