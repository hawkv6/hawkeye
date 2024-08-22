package processor

import "github.com/hawkv6/hawkeye/pkg/domain"

const Subsystem = "processor"

type Processor interface {
	ProcessNodes([]domain.Node)
	ProcessLinks([]domain.Link) error
	ProcessPrefixes([]domain.Prefix)
	ProcessSids([]domain.Sid)
	Start()
	Stop()
}

type NodeProcessor interface {
	ProcessNodes([]domain.Node)
}

type LinkProcessor interface {
	ProcessLinks([]domain.Link) error
}

type PrefixProcessor interface {
	ProcessPrefixes([]domain.Prefix)
}

type SidProcessor interface {
	ProcessSids([]domain.Sid)
}
