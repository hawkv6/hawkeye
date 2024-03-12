package processor

import (
	"fmt"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type DefaultProcessor struct {
	log   *logrus.Entry
	graph graph.Graph
	cache cache.CacheService
}

func NewDefaultProcessor(graph graph.Graph, cache cache.CacheService) *DefaultProcessor {
	return &DefaultProcessor{
		log:   logging.DefaultLogger.WithField("subsystem", Subsystem),
		graph: graph,
		cache: cache,
	}
}

func (processor *DefaultProcessor) createNetworkNode(node domain.Node) error {
	id := node.GetIgpRouterId()
	if processor.graph.NodeExists(id) {
		return fmt.Errorf("Error creating network node - node already with id %s already exists", id)
	}
	graphNode, err := processor.graph.AddNode(graph.NewDefaultNode(id))
	if err != nil {
		return err
	}
	name := node.GetName()
	graphNode.SetName(name)
	processor.log.Debugf("Added node %s to graph with id %s", name, id)
	return nil
}

func (processor *DefaultProcessor) CreateNetworkNodes(nodes []domain.Node) error {
	for _, node := range nodes {
		if err := processor.createNetworkNode(node); err != nil {
			return err
		}
	}
	return nil
}

func (processor *DefaultProcessor) CreateNetworkEdges(links []domain.Link) error {
	for _, link := range links {
		if err := processor.addLinkToGraph(link); err != nil {
			return err
		}
	}
	return nil
}

func (processor *DefaultProcessor) getLinkWeights(link domain.Link) map[string]float64 {
	return map[string]float64{
		"latency":            float64(link.GetUnidirLinkDelay()),
		"jitter":             float64(link.GetUnidirDelayVariation()),
		"availableBandwidth": float64(link.GetUnidirAvailableBandwidth()),
		"utilizedBandwidth":  float64(link.GetUnidirBandwidthUtilization()),
		"loss":               float64(link.GetUnidirPacketLoss()),
	}
}
func (processor *DefaultProcessor) addLinkToGraph(link domain.Link) error {
	from, err := processor.getNode(link.GetIgpRouterId())
	if err != nil {
		return err
	}

	to, err := processor.getNode(link.GetRemoteIgpRouterId())
	if err != nil {
		return err
	}
	edge := graph.NewDefaultEdge(link.GetKey(), from, to, processor.getLinkWeights(link))
	if err := processor.graph.AddEdge(edge); err != nil {
		return err
	}
	processor.log.Debugf("Added edge to graph between %s and %s", from.GetName(), to.GetName())
	return nil
}

func (processor *DefaultProcessor) getNode(nodeId string) (graph.Node, error) {
	if !processor.graph.NodeExists(nodeId) {
		return nil, fmt.Errorf("Node with id %s does not exist in graph", nodeId)
	}
	return processor.graph.GetNode(nodeId)
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
		processor.log.Debugf("Added client network %s/%d to cache ", clientNetwork.GetPrefix(), clientNetwork.GetPrefixLength())
	}
	return nil
}

func (processor *DefaultProcessor) CreateSids(sids []domain.Sid) error {
	for _, sid := range sids {
		processor.cache.StoreSids(sid)
		processor.log.Debugf("Added SRv6 SID %s to cache", sid.GetSid())
	}
	return nil
}
