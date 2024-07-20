package domain

import (
	"github.com/go-playground/validator"
	"github.com/hawkv6/hawkeye/pkg/graph"
)

type PathResult interface {
	PathRequest
	graph.Path
	GetIpv6SidAddresses() []string
	GetServiceSidList() []string
	SetServiceSidList([]string)
}

type DomainPathResult struct {
	PathRequest
	graph.Path
	ipv6SidAddresses    []string
	serviceSidAddresses []string
}

type DomainPathResultInput struct {
	PathRequest      PathRequest `validate:"required"`
	Ipv6SidAddresses []string    `validate:"required,dive,ipv6"`
}

func NewDomainPathResult(pathRequest PathRequest, shortestPath graph.Path, ipv6SidAddresses []string) (*DomainPathResult, error) {
	domainPathResultInput := &DomainPathResultInput{
		Ipv6SidAddresses: ipv6SidAddresses,
		PathRequest:      pathRequest,
	}

	validator := validator.New()
	err := validator.Struct(domainPathResultInput)
	if err != nil {
		return nil, err
	}
	defaultPathResult := &DomainPathResult{
		PathRequest:      pathRequest,
		Path:             shortestPath,
		ipv6SidAddresses: domainPathResultInput.Ipv6SidAddresses,
	}
	return defaultPathResult, nil
}

func (pathResponse *DomainPathResult) GetIpv6SidAddresses() []string {
	return pathResponse.ipv6SidAddresses
}

func (pathResponse *DomainPathResult) GetServiceSidList() []string {
	return pathResponse.serviceSidAddresses
}

func (pathResponse *DomainPathResult) SetServiceSidList(serviceSidAddresses []string) {
	pathResponse.serviceSidAddresses = serviceSidAddresses
}
