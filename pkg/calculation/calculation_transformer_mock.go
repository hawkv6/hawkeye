// Code generated by MockGen. DO NOT EDIT.
// Source: calculation_transformer.go
//
// Generated by this command:
//
//	mockgen -source calculation_transformer.go -destination calculation_transformer_mock.go -package calculation
//

// Package calculation is a generated GoMock package.
package calculation

import (
	reflect "reflect"

	domain "github.com/hawkv6/hawkeye/pkg/domain"
	graph "github.com/hawkv6/hawkeye/pkg/graph"
	gomock "go.uber.org/mock/gomock"
)

// MockCalculationTransformer is a mock of CalculationTransformer interface.
type MockCalculationTransformer struct {
	ctrl     *gomock.Controller
	recorder *MockCalculationTransformerMockRecorder
}

// MockCalculationTransformerMockRecorder is the mock recorder for MockCalculationTransformer.
type MockCalculationTransformerMockRecorder struct {
	mock *MockCalculationTransformer
}

// NewMockCalculationTransformer creates a new mock instance.
func NewMockCalculationTransformer(ctrl *gomock.Controller) *MockCalculationTransformer {
	mock := &MockCalculationTransformer{ctrl: ctrl}
	mock.recorder = &MockCalculationTransformerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCalculationTransformer) EXPECT() *MockCalculationTransformerMockRecorder {
	return m.recorder
}

// TransformResult mocks base method.
func (m *MockCalculationTransformer) TransformResult(path graph.Path, pathRequest domain.PathRequest, algorithm uint32) domain.PathResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransformResult", path, pathRequest, algorithm)
	ret0, _ := ret[0].(domain.PathResult)
	return ret0
}

// TransformResult indicates an expected call of TransformResult.
func (mr *MockCalculationTransformerMockRecorder) TransformResult(path, pathRequest, algorithm any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransformResult", reflect.TypeOf((*MockCalculationTransformer)(nil).TransformResult), path, pathRequest, algorithm)
}