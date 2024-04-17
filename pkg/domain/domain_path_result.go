package domain

import (
	"github.com/go-playground/validator"
	"github.com/hawkv6/hawkeye/pkg/graph"
)

type PathResult interface {
	PathRequest
	graph.PathResult
	GetIpv6SidAddresses() []string
}

type DomainPathResult struct {
	PathRequest
	graph.PathResult
	ipv6SidAddresses []string `validate:"required,dive,ipv6"`
}

func NewDomainPathResult(pathRequest PathRequest, graphResult graph.PathResult, ipv6SidAddresses []string) (*DomainPathResult, error) {
	defaultPathResult := &DomainPathResult{
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

func (pathResponse *DomainPathResult) GetIpv6SidAddresses() []string {
	return pathResponse.ipv6SidAddresses
}
