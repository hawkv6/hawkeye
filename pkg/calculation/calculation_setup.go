package calculation

import (
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
)

type CalculationOptions struct {
	graph           graph.Graph
	sourceNode      graph.Node
	destinationNode graph.Node
	weightKeys      []helper.WeightKey
	calculationMode CalculationMode
	maxConstraints  map[helper.WeightKey]float64
	minConstraints  map[helper.WeightKey]float64
}

type SfcCalculationOptions struct {
	serviceFunctionChain [][]string
	routerServiceMap     map[string]string
}

type CalculationSetup interface {
	PerformSetup(pathRequest domain.PathRequest) (*CalculationOptions, error)
	PerformServiceFunctionChainSetup(intent domain.Intent) (*SfcCalculationOptions, error)
	GetWeightKeysandCalculationMode(intents []domain.Intent) ([]helper.WeightKey, CalculationMode)
}
