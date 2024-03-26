package domain

import (
	"github.com/go-playground/validator"
	"github.com/hawkv6/hawkeye/pkg/graph"
)

type DefaultPathResult struct {
	PathRequest
	graph.PathResult
	ipv6SidAddresses []string `validate:"required,dive,ipv6"`
}

func NewDefaultPathResult(pathRequest PathRequest, graphResult graph.PathResult, ipv6SidAddresses []string) (*DefaultPathResult, error) {
	defaultPathResult := &DefaultPathResult{
		PathRequest:      pathRequest,
		PathResult:       graphResult,
		ipv6SidAddresses: ipv6SidAddresses,
	}
	validator := validator.New()
	err := validator.Struct(defaultPathResult)
	if err != nil {
		return nil, err
	}
	return defaultPathResult, nil
}

func (defaultPathResponse *DefaultPathResult) GetIpv6SidAddresses() []string {
	return defaultPathResponse.ipv6SidAddresses
}
