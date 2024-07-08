package calculation

import (
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type CalculationMode int

const (
	CalculationModeUndefined CalculationMode = iota
	CalculationModeSum
	CalculationModeMin
	CalculationModeMax
)

type Calculation interface {
	Execute() (graph.Path, error)
}

type BaseCalculation struct {
	log             *logrus.Entry
	graph           graph.Graph
	source          graph.Node
	destination     graph.Node
	weightTypes     []helper.WeightKey
	calculationMode CalculationMode
	maxConstraints  map[helper.WeightKey]float64
	minConstraints  map[helper.WeightKey]float64
}

func NewBaseCalculation(graph graph.Graph, source graph.Node, destination graph.Node, weightTypes []helper.WeightKey, calculationMode CalculationMode, maxConstraints, minConstraints map[helper.WeightKey]float64) *BaseCalculation {
	return &BaseCalculation{
		log:             logging.DefaultLogger.WithField("subsystem", subsystem),
		graph:           graph,
		source:          source,
		destination:     destination,
		weightTypes:     weightTypes,
		calculationMode: calculationMode,
		maxConstraints:  maxConstraints,
		minConstraints:  minConstraints,
	}
}
