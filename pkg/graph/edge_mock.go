// Code generated by MockGen. DO NOT EDIT.
// Source: edge.go
//
// Generated by this command:
//
//	mockgen -source edge.go -destination edge_mock.go -package graph
//

// Package graph is a generated GoMock package.
package graph

import (
	reflect "reflect"

	helper "github.com/hawkv6/hawkeye/pkg/helper"
	gomock "go.uber.org/mock/gomock"
)

// MockEdge is a mock of Edge interface.
type MockEdge struct {
	ctrl     *gomock.Controller
	recorder *MockEdgeMockRecorder
}

// MockEdgeMockRecorder is the mock recorder for MockEdge.
type MockEdgeMockRecorder struct {
	mock *MockEdge
}

// NewMockEdge creates a new mock instance.
func NewMockEdge(ctrl *gomock.Controller) *MockEdge {
	mock := &MockEdge{ctrl: ctrl}
	mock.recorder = &MockEdgeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEdge) EXPECT() *MockEdgeMockRecorder {
	return m.recorder
}

// From mocks base method.
func (m *MockEdge) From() Node {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "From")
	ret0, _ := ret[0].(Node)
	return ret0
}

// From indicates an expected call of From.
func (mr *MockEdgeMockRecorder) From() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "From", reflect.TypeOf((*MockEdge)(nil).From))
}

// GetAllWeights mocks base method.
func (m *MockEdge) GetAllWeights() map[helper.WeightKey]float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllWeights")
	ret0, _ := ret[0].(map[helper.WeightKey]float64)
	return ret0
}

// GetAllWeights indicates an expected call of GetAllWeights.
func (mr *MockEdgeMockRecorder) GetAllWeights() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllWeights", reflect.TypeOf((*MockEdge)(nil).GetAllWeights))
}

// GetFlexibleAlgorithms mocks base method.
func (m *MockEdge) GetFlexibleAlgorithms() map[uint32]struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlexibleAlgorithms")
	ret0, _ := ret[0].(map[uint32]struct{})
	return ret0
}

// GetFlexibleAlgorithms indicates an expected call of GetFlexibleAlgorithms.
func (mr *MockEdgeMockRecorder) GetFlexibleAlgorithms() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlexibleAlgorithms", reflect.TypeOf((*MockEdge)(nil).GetFlexibleAlgorithms))
}

// GetId mocks base method.
func (m *MockEdge) GetId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetId")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetId indicates an expected call of GetId.
func (mr *MockEdgeMockRecorder) GetId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetId", reflect.TypeOf((*MockEdge)(nil).GetId))
}

// GetWeight mocks base method.
func (m *MockEdge) GetWeight(kind helper.WeightKey) float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWeight", kind)
	ret0, _ := ret[0].(float64)
	return ret0
}

// GetWeight indicates an expected call of GetWeight.
func (mr *MockEdgeMockRecorder) GetWeight(kind any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWeight", reflect.TypeOf((*MockEdge)(nil).GetWeight), kind)
}

// SetWeight mocks base method.
func (m *MockEdge) SetWeight(kind helper.WeightKey, weight float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetWeight", kind, weight)
}

// SetWeight indicates an expected call of SetWeight.
func (mr *MockEdgeMockRecorder) SetWeight(kind, weight any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWeight", reflect.TypeOf((*MockEdge)(nil).SetWeight), kind, weight)
}

// To mocks base method.
func (m *MockEdge) To() Node {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "To")
	ret0, _ := ret[0].(Node)
	return ret0
}

// To indicates an expected call of To.
func (mr *MockEdgeMockRecorder) To() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "To", reflect.TypeOf((*MockEdge)(nil).To))
}

// UpdateFlexibleAlgorithms mocks base method.
func (m *MockEdge) UpdateFlexibleAlgorithms() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateFlexibleAlgorithms")
}

// UpdateFlexibleAlgorithms indicates an expected call of UpdateFlexibleAlgorithms.
func (mr *MockEdgeMockRecorder) UpdateFlexibleAlgorithms() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFlexibleAlgorithms", reflect.TypeOf((*MockEdge)(nil).UpdateFlexibleAlgorithms))
}