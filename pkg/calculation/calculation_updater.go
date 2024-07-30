package calculation

import (
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/helper"
)

type CalculationUpdateOptions struct {
	currentPathResult     domain.PathResult
	currentAppliedSidList []string
	weightKeys            []helper.WeightKey
	calculationMode       CalculationMode
	newPathResult         domain.PathResult
	pathRequest           domain.PathRequest
	streamSession         domain.StreamSession
}

type CalculationUpdater interface {
	UpdateCalculation(*CalculationUpdateOptions) (domain.PathResult, error)
}
