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
	ipv6SidAddresses    []string `validate:"required,dive,ipv6"`
	serviceSidAddresses []string
}

func NewDomainPathResult(pathRequest PathRequest, shortestPath graph.Path, ipv6SidAddresses []string) (*DomainPathResult, error) {
	defaultPathResult := &DomainPathResult{
		PathRequest:      pathRequest,
		Path:             shortestPath,
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

func (pathResponse *DomainPathResult) GetServiceSidList() []string {
	return pathResponse.serviceSidAddresses
}

func (pathResponse *DomainPathResult) SetServiceSidList(serviceSidAddresses []string) {
	pathResponse.serviceSidAddresses = serviceSidAddresses
}
