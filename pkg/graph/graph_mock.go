// Code generated by MockGen. DO NOT EDIT.
// Source: graph.go
//
// Generated by this command:
//
//	mockgen -source graph.go -destination graph_mock.go -package graph
//

// Package graph is a generated GoMock package.
package graph

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockGraph is a mock of Graph interface.
type MockGraph struct {
	ctrl     *gomock.Controller
	recorder *MockGraphMockRecorder
}

// MockGraphMockRecorder is the mock recorder for MockGraph.
type MockGraphMockRecorder struct {
	mock *MockGraph
}

// NewMockGraph creates a new mock instance.
func NewMockGraph(ctrl *gomock.Controller) *MockGraph {
	mock := &MockGraph{ctrl: ctrl}
	mock.recorder = &MockGraphMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGraph) EXPECT() *MockGraphMockRecorder {
	return m.recorder
}

// AddEdge mocks base method.
func (m *MockGraph) AddEdge(arg0 Edge) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEdge", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddEdge indicates an expected call of AddEdge.
func (mr *MockGraphMockRecorder) AddEdge(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEdge", reflect.TypeOf((*MockGraph)(nil).AddEdge), arg0)
}

// AddNode mocks base method.
func (m *MockGraph) AddNode(arg0 Node) Node {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNode", arg0)
	ret0, _ := ret[0].(Node)
	return ret0
}

// AddNode indicates an expected call of AddNode.
func (mr *MockGraphMockRecorder) AddNode(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNode", reflect.TypeOf((*MockGraph)(nil).AddNode), arg0)
}

// DeleteEdge mocks base method.
func (m *MockGraph) DeleteEdge(arg0 Edge) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteEdge", arg0)
}

// DeleteEdge indicates an expected call of DeleteEdge.
func (mr *MockGraphMockRecorder) DeleteEdge(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEdge", reflect.TypeOf((*MockGraph)(nil).DeleteEdge), arg0)
}

// DeleteNode mocks base method.
func (m *MockGraph) DeleteNode(arg0 Node) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteNode", arg0)
}

// DeleteNode indicates an expected call of DeleteNode.
func (mr *MockGraphMockRecorder) DeleteNode(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNode", reflect.TypeOf((*MockGraph)(nil).DeleteNode), arg0)
}

// EdgeExists mocks base method.
func (m *MockGraph) EdgeExists(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EdgeExists", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// EdgeExists indicates an expected call of EdgeExists.
func (mr *MockGraphMockRecorder) EdgeExists(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EdgeExists", reflect.TypeOf((*MockGraph)(nil).EdgeExists), arg0)
}

// GetEdge mocks base method.
func (m *MockGraph) GetEdge(arg0 string) Edge {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEdge", arg0)
	ret0, _ := ret[0].(Edge)
	return ret0
}

// GetEdge indicates an expected call of GetEdge.
func (mr *MockGraphMockRecorder) GetEdge(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEdge", reflect.TypeOf((*MockGraph)(nil).GetEdge), arg0)
}

// GetEdges mocks base method.
func (m *MockGraph) GetEdges() map[string]Edge {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEdges")
	ret0, _ := ret[0].(map[string]Edge)
	return ret0
}

// GetEdges indicates an expected call of GetEdges.
func (mr *MockGraphMockRecorder) GetEdges() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEdges", reflect.TypeOf((*MockGraph)(nil).GetEdges))
}

// GetNode mocks base method.
func (m *MockGraph) GetNode(arg0 string) Node {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNode", arg0)
	ret0, _ := ret[0].(Node)
	return ret0
}

// GetNode indicates an expected call of GetNode.
func (mr *MockGraphMockRecorder) GetNode(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNode", reflect.TypeOf((*MockGraph)(nil).GetNode), arg0)
}

// GetNodes mocks base method.
func (m *MockGraph) GetNodes() map[string]Node {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodes")
	ret0, _ := ret[0].(map[string]Node)
	return ret0
}

// GetNodes indicates an expected call of GetNodes.
func (mr *MockGraphMockRecorder) GetNodes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodes", reflect.TypeOf((*MockGraph)(nil).GetNodes))
}

// GetSubGraph mocks base method.
func (m *MockGraph) GetSubGraph(arg0 uint32) Graph {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubGraph", arg0)
	ret0, _ := ret[0].(Graph)
	return ret0
}

// GetSubGraph indicates an expected call of GetSubGraph.
func (mr *MockGraphMockRecorder) GetSubGraph(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubGraph", reflect.TypeOf((*MockGraph)(nil).GetSubGraph), arg0)
}

// Lock mocks base method.
func (m *MockGraph) Lock() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Lock")
}

// Lock indicates an expected call of Lock.
func (mr *MockGraphMockRecorder) Lock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lock", reflect.TypeOf((*MockGraph)(nil).Lock))
}

// NodeExists mocks base method.
func (m *MockGraph) NodeExists(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NodeExists", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// NodeExists indicates an expected call of NodeExists.
func (mr *MockGraphMockRecorder) NodeExists(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NodeExists", reflect.TypeOf((*MockGraph)(nil).NodeExists), arg0)
}

// Unlock mocks base method.
func (m *MockGraph) Unlock() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Unlock")
}

// Unlock indicates an expected call of Unlock.
func (mr *MockGraphMockRecorder) Unlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unlock", reflect.TypeOf((*MockGraph)(nil).Unlock))
}

// UpdateSubGraphs mocks base method.
func (m *MockGraph) UpdateSubGraphs() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateSubGraphs")
}

// UpdateSubGraphs indicates an expected call of UpdateSubGraphs.
func (mr *MockGraphMockRecorder) UpdateSubGraphs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSubGraphs", reflect.TypeOf((*MockGraph)(nil).UpdateSubGraphs))
}
