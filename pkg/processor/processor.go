package processor

import "github.com/hawkv6/hawkeye/pkg/domain"

const Subsystem = "processor"

type Processor interface {
	CreateGraphNodes([]domain.Node) error
	CreateGraphEdges([]domain.Link) error
	CreateClientNetworks([]domain.Prefix) error
	CreateSids([]domain.Sid) error
	Start()
	Stop()
}
