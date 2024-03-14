package processor

import (
	"fmt"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type DefaultProcessor struct {
	log       *logrus.Entry
	graph     graph.Graph
	cache     cache.CacheService
	eventChan chan domain.NetworkEvent
	quitChan  chan struct{}
	helper    helper.Helper
}

func NewDefaultProcessor(graph graph.Graph, cache cache.CacheService, eventChan chan domain.NetworkEvent, helper helper.Helper) *DefaultProcessor {
	return &DefaultProcessor{
		log:       logging.DefaultLogger.WithField("subsystem", Subsystem),
		graph:     graph,
		cache:     cache,
		eventChan: eventChan,
		helper:    helper,
		quitChan:  make(chan struct{}),
	}
}

func (processor *DefaultProcessor) setNodeName(node domain.Node, graphNode graph.Node) {
	nodeName := node.GetName()
	if graphNode.GetName() != nodeName {
		graphNode.SetName(nodeName)
	}
}

func (processor *DefaultProcessor) addNodeToGraph(node domain.Node) error {
	id := node.GetIgpRouterId()
	graphNode, exist := processor.graph.GetNode(id)
	if exist {
		processor.log.Debugf("Node with id %s already exists in graph", id)
	} else {
		processor.log.Debugf("Add node %s to graph with igp router id %s", node.GetName(), id)
		var err error
		graphNode, err = processor.graph.AddNode(graph.NewDefaultNode(id))
		if err != nil {
			return err
		}
	}
	processor.setNodeName(node, graphNode)
	return nil
}

func (processor *DefaultProcessor) addNodeToCache(node domain.Node) {
	processor.log.Debugf("Add node %s to cache with igp router id id %s", node.GetName(), node.GetIgpRouterId())
	processor.cache.StoreNode(node)
}

func (processor *DefaultProcessor) deleteNodeIfExists(key string) error {
	node, ok := processor.cache.GetNodeByKey(key)
	if !ok {
		return fmt.Errorf("Node with key %s does not exist in cache", key)
	}
	processor.log.Debugf("Delete node %s with igp router id %s from graph", node.GetName(), node.GetIgpRouterId())
	graphNode, exist := processor.graph.GetNode(node.GetIgpRouterId())
	if !exist {
		return fmt.Errorf("Node with igp router id %s does not exist in graph", node.GetIgpRouterId())
	}
	processor.graph.DeleteNode(graphNode)
	return nil
}

func (processor *DefaultProcessor) CreateGraphNodes(nodes []domain.Node) error {
	for _, node := range nodes {
		if err := processor.addNodeToGraph(node); err != nil {
			return err
		}
		processor.addNodeToCache(node)
	}
	return nil
}

func (processor *DefaultProcessor) CreateGraphEdges(links []domain.Link) error {
	for _, link := range links {
		if err := processor.addLinkToGraph(link); err != nil {
			return err
		}
	}
	return nil
}

func (processor *DefaultProcessor) getLinkWeights(link domain.Link) map[string]float64 {
	return map[string]float64{
		processor.helper.GetLatencyKey():            float64(link.GetUnidirLinkDelay()),
		processor.helper.GetJitterKey():             float64(link.GetUnidirDelayVariation()),
		processor.helper.GetAvailableBandwidthKey(): float64(link.GetUnidirAvailableBandwidth()),
		processor.helper.GetUtilizedBandwidthKey():  float64(link.GetUnidirBandwidthUtilization()),
		processor.helper.GetPacketLossKey():         float64(link.GetUnidirPacketLoss()),
	}
}

func (processor *DefaultProcessor) getOrCreateNode(nodeId string) (graph.Node, error) {
	// it's possible that link events are received before node events
	// so we need to ensure the node exists in the graph before adding the edge
	if processor.graph.NodeExists(nodeId) {
		node, _ := processor.graph.GetNode(nodeId)
		return node, nil
	}
	node, exists := processor.cache.GetNodeByIgpRouterId(nodeId)
	if !exists {
		processor.log.Errorf("Node with igp router id %s not in cache - create it only in graph", nodeId)
		return processor.graph.AddNode(graph.NewDefaultNode(nodeId))
	}
	if err := processor.addNodeToGraph(node); err != nil {
		return nil, err
	}
	graphNode, _ := processor.graph.GetNode(nodeId)
	return graphNode, nil
}

func (processor *DefaultProcessor) deleteEdgeIfExists(key string) error {
	if processor.graph.EdgeExists(key) {
		edge, exist := processor.graph.GetEdge(key)
		if !exist {
			return fmt.Errorf("Edge with key %s does not exist in graph", key)
		}
		processor.log.Debugf("Delete edge with key %s from graph between %s and %s", key, edge.From().GetName(), edge.To().GetName())
		processor.graph.DeleteEdge(edge)
	}
	return nil
}

func (processor *DefaultProcessor) addEdgeToGraph(edge graph.Edge) error {
	processor.log.Debugf("Add edge with key %s to graph between %s and %s", edge.GetId(), edge.From().GetName(), edge.To().GetName())
	if err := processor.graph.AddEdge(edge); err != nil {
		return err
	}
	return nil
}

func (processor *DefaultProcessor) addLinkToGraph(link domain.Link) error {
	from, err := processor.getOrCreateNode(link.GetIgpRouterId())
	if err != nil {
		return err
	}

	to, err := processor.getOrCreateNode(link.GetRemoteIgpRouterId())
	if err != nil {
		return err
	}

	key := link.GetKey()
	if !processor.graph.EdgeExists(key) {
		weights := processor.getLinkWeights(link)
		for _, weight := range weights {
			if weight == 0 {
				return fmt.Errorf("Link contains zero values, link %s is created during next update", key)
			}
		}
		return processor.addEdgeToGraph(graph.NewDefaultEdge(key, from, to, weights))
	}
	return nil
}

func (processor *DefaultProcessor) setEdgeWeight(edge graph.Edge, key string, value float64) error {
	if value == 0 {
		return fmt.Errorf("Value is 0, not setting %s", key)
	}
	currentValue, err := edge.GetWeight(key)
	if err != nil {
		return fmt.Errorf("Error getting %s weight: %v", key, err)
	}
	if currentValue != value {
		edge.SetWeight(key, value)
	}
	return nil
}

func (processor *DefaultProcessor) updateLinkInGraph(link domain.Link) error {
	key := link.GetKey()
	processor.log.Debugln("Updating link in graph with key: ", key)
	edge, exist := processor.graph.GetEdge(key)
	if !exist {
		processor.log.Debugf("Link with key %s does not exist in graph, create it", key)
		return processor.addLinkToGraph(link)
	}
	for weightKey, weightValue := range processor.getLinkWeights(link) {
		if err := processor.setEdgeWeight(edge, weightKey, weightValue); err != nil {
			return err
		}
	}
	processor.cache.StoreLink(link)
	return nil
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

func (processor *DefaultProcessor) Start() {
	processor.log.Infoln("Starting processor")
	for {
		select {
		case event := <-processor.eventChan:
			processor.processEvent(event)
		case <-processor.quitChan:
			processor.log.Infoln("Stopping processor")
			return
		}
	}
}

func (processor *DefaultProcessor) processEvent(event domain.NetworkEvent) {
	switch eventType := event.(type) {
	case *domain.AddNodeEvent:
		processor.log.Debugln("Received AddNodeEvent: ", eventType.GetKey())
		if err := processor.addNodeToGraph(eventType.Node); err != nil {
			processor.log.Errorln("Error creating network node: ", err)
		}
		processor.addNodeToCache(eventType.Node)
	case *domain.DeleteNodeEvent:
		processor.log.Debugln("Received DeleteNodeEvent: ", eventType.GetKey())
		if err := processor.deleteNodeIfExists(eventType.GetKey()); err != nil {
			processor.log.Errorln("Error deleting network node: ", err)
		}
	case *domain.AddLinkEvent:
		processor.log.Debugln("Received AddLinkEvent: ", eventType.GetKey())
		if err := processor.addLinkToGraph(eventType.Link); err != nil {
			processor.log.Errorln("Error adding link to graph: ", err)
		}
	case *domain.UpdateLinkEvent:
		processor.log.Debugln("Received UpdateLinkEvent: ", eventType.GetKey())
		if err := processor.updateLinkInGraph(eventType.Link); err != nil {
			processor.log.Errorln("Error updating link in graph: ", err)
		}
	case *domain.DeleteLinkEvent:
		processor.log.Debugln("Received DeleteLinkEvent: ", eventType.GetKey())
		if err := processor.deleteEdgeIfExists(eventType.GetKey()); err != nil {
			processor.log.Errorln("Error deleting edge: ", err)
		}
		// 	processor.addLinkToGraph(event.GetLink())
		// case domain.AddPrefixEvent:
		// 	processor.cache.StoreClientNetwork(event.GetPrefix())
		// case domain.AddSidEvent:
		// 	processor.cache.StoreSids(event.GetSid())
		// case domain.DeleteSidEvent:
		// 	processor.cache.DeleteSid(event.GetKey())
		// case domain.DeletePrefixEvent:
		// 	processor.cache.DeleteClientNetwork(event.GetKey())
	}
}

func (processor *DefaultProcessor) Stop() {
	processor.log.Infoln("Stopping processor")
	close(processor.quitChan)
}
