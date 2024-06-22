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
