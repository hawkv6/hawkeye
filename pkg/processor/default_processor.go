package processor

import (
	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type DefaultProcessor struct {
	log     *logrus.Entry
	network graph.Graph
	cache   cache.CacheService
}

func NewDefaultProcessor(graph graph.Graph, cache cache.CacheService) *DefaultProcessor {
	return &DefaultProcessor{
		log:     logging.DefaultLogger.WithField("subsystem", Subsystem),
		network: graph,
		cache:   cache,
	}
}

func (processor *DefaultProcessor) CreateNetworkGraph(links []domain.Link) error {
	for _, link := range links {
		if err := processor.addLinkToNetwork(link); err != nil {
			return err
		}
	}
	return nil
}

func (processor *DefaultProcessor) addLinkToNetwork(link domain.Link) error {
	from, err := processor.getOrCreateNode(link.GetIgpRouterId())
	if err != nil {
		return err
	}

	to, err := processor.getOrCreateNode(link.GetRemoteIgpRouterId())
	if err != nil {
		return err
	}

	edge := graph.NewDefaultEdge(link.GetKey(), from, to, map[string]float64{"delay": link.GetUnidirLinkDelay()})
	if err := processor.network.AddEdge(edge); err != nil {
		return err
	}

	return nil
}

func (processor *DefaultProcessor) getOrCreateNode(nodeId string) (graph.Node, error) {
	if processor.network.NodeExists(nodeId) {
		return processor.network.GetNode(nodeId)
	}
	node := graph.NewDefaultNode(nodeId)
	err := processor.network.AddNode(node)
	return node, err
}

func (processor *DefaultProcessor) getPrefixCounts(prefixes []domain.Prefix) (map[string]int, map[string]domain.Prefix) {
	prefixCounts := make(map[string]int)
	prefixMap := make(map[string]domain.Prefix)
	for _, prefix := range prefixes {
		// Ignore lo0 addressses
		if prefix.GetPrefixLength() == 128 {
			continue
		}
		network := prefix.GetPrefix()
		if _, ok := prefixCounts[network]; !ok {
			prefixCounts[network] = 1
		} else {
			prefixCounts[network]++
		}
		prefixMap[network] = prefix
	}
	return prefixCounts, prefixMap
}

func (processor *DefaultProcessor) getClientNetworks(prefixCounts map[string]int, prefixMap map[string]domain.Prefix) []domain.Prefix {
	// TODO: Find way to remove SRv6 locator prefixes from clientNetworks (with Segments)
	prefixes := make([]domain.Prefix, 0)
	for prefix, count := range prefixCounts {
		if count == 1 {
			prefixes = append(prefixes, prefixMap[prefix])
		}
	}
	return prefixes
}

func (processor *DefaultProcessor) CreateClientNetworks(prefixes []domain.Prefix) error {
	clientNetworks := processor.getClientNetworks(processor.getPrefixCounts(prefixes))
	for _, clientNetwork := range clientNetworks {
		processor.cache.StoreClientNetwork(clientNetwork)
		processor.log.Debugln("Client Network: ", clientNetwork)
	}
	return nil
}

func (processor *DefaultProcessor) CreateSids(sids []domain.Sid) error {
	for _, sid := range sids {
		processor.cache.StoreSids(sid)
	}
	return nil
}
