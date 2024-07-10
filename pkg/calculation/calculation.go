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
	weightKeys      []helper.WeightKey
	calculationMode CalculationMode
	maxConstraints  map[helper.WeightKey]float64
	minConstraints  map[helper.WeightKey]float64
}

func NewBaseCalculation(options *CalculationOptions) *BaseCalculation {
	return &BaseCalculation{
		log:             logging.DefaultLogger.WithField("subsystem", subsystem),
		graph:           options.graph,
		source:          options.sourceNode,
		destination:     options.destinationNode,
		weightKeys:      options.weightKeys,
		calculationMode: options.calculationMode,
		maxConstraints:  options.maxConstraints,
		minConstraints:  options.minConstraints,
	}
}
