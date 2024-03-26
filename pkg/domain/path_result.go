package domain

import "github.com/hawkv6/hawkeye/pkg/graph"

type PathResult interface {
	PathRequest
	graph.PathResult
	GetIpv6SidAddresses() []string
}
