package processor

import "github.com/hawkv6/hawkeye/pkg/domain"

const Subsystem = "processor"

type Processor interface {
	CreateNetworkNodes([]domain.Node) error
	CreateNetworkEdges([]domain.Link) error
	CreateClientNetworks([]domain.Prefix) error
	CreateSids([]domain.Sid) error
}
