package processor

import (
	"fmt"
	"time"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type NetworkProcessor struct {
	log                 *logrus.Entry
	graph               graph.Graph
	cache               cache.Cache
	eventChan           chan domain.NetworkEvent
	quitChan            chan struct{}
	prefixCounts        map[string]int
	updateChan          chan struct{}
	needsSubgraphUpdate bool
}

func NewNetworkProcessor(graph graph.Graph, cache cache.Cache, eventChan chan domain.NetworkEvent, helper helper.Helper, updateChan chan struct{}) *NetworkProcessor {
	return &NetworkProcessor{
		log:                 logging.DefaultLogger.WithField("subsystem", Subsystem),
		graph:               graph,
		cache:               cache,
		eventChan:           eventChan,
		quitChan:            make(chan struct{}),
		prefixCounts:        make(map[string]int),
		updateChan:          updateChan,
		needsSubgraphUpdate: false,
	}
}

func (processor *NetworkProcessor) updateNode(node graph.Node, name string, srAlgorithm []uint32) {
	processor.log.Debugf("Update node %s with igp router id %s in graph", name, node.GetId())
	node.SetName(name)
	node.SetFlexibleAlgorithms(srAlgorithm)
}

func (processor *NetworkProcessor) addNodeToGraph(igpRouterId, name string, srAlgorithm []uint32) {
	if processor.graph.NodeExists(igpRouterId) {
		processor.log.Debugf("Node with id %s already exists in graph", igpRouterId)
		processor.updateNode(processor.graph.GetNode(igpRouterId), name, srAlgorithm)
		return
	}
	processor.graph.AddNode(graph.NewNetworkNode(igpRouterId, name, srAlgorithm))
}

func (processor *NetworkProcessor) addNodeToCache(node domain.Node) {
	processor.cache.StoreNode(node)
}

func (processor *NetworkProcessor) removeNodeFromGraphIfExists(node domain.Node) {
	igpRouterId := node.GetIgpRouterId()
	name := node.GetName()
	if processor.graph.NodeExists(igpRouterId) {
		graphNode := processor.graph.GetNode(igpRouterId)
		processor.graph.DeleteNode(graphNode)
	} else {
		processor.log.Debugf("Node with igp router id %s and name %s does not exist in graph", igpRouterId, name)
	}
}

func (processor *NetworkProcessor) deleteNode(key string) {
	processor.log.Debugf("Delete node with key %s", key)
	node := processor.cache.GetNodeByKey(key)
	if node == nil {
		processor.log.Debugf("Node with key %s does not exist in cache and thus also not in graph", key)
		return
	}
	processor.cache.RemoveNode(node)
	processor.removeNodeFromGraphIfExists(node)
}

func (processor *NetworkProcessor) addOrUpdateNodeInGraphAndCache(node domain.Node) {
	igpRouterId := node.GetIgpRouterId()
	name := node.GetName()
	srAlgorithm := node.GetSrAlgorithm()
	processor.log.Debugf("Add node %s with igp router id %s and SR algorithm %v to graph and cache", name, igpRouterId, srAlgorithm)
	processor.addNodeToGraph(igpRouterId, name, srAlgorithm)
	processor.addNodeToCache(node)

}

func (processor *NetworkProcessor) ProcessNodes(nodes []domain.Node) {
	for _, node := range nodes {
		processor.addOrUpdateNodeInGraphAndCache(node)
	}
}

func (processor *NetworkProcessor) ProcessLinks(links []domain.Link) error {
	for _, link := range links {
		if err := processor.addLinkToGraph(link); err != nil {
			return err
		}
	}
	processor.graph.UpdateSubGraphs()
	return nil
}

func (processor *NetworkProcessor) getCurrentLinkWeights(link domain.Link) map[helper.WeightKey]float64 {
	return map[helper.WeightKey]float64{
		helper.IgpMetricKey:            float64(link.GetIgpMetric()),
		helper.LatencyKey:              float64(link.GetUnidirLinkDelay()),
		helper.JitterKey:               float64(link.GetUnidirDelayVariation()),
		helper.MaximumLinkBandwidth:    float64(link.GetMaxLinkBWKbps()),
		helper.AvailableBandwidthKey:   float64(link.GetUnidirAvailableBandwidth()),
		helper.UtilizedBandwidthKey:    float64(link.GetUnidirBandwidthUtilization()),
		helper.PacketLossKey:           float64(link.GetUnidirPacketLoss()),
		helper.NormalizedLatencyKey:    link.GetNormalizedUnidirLinkDelay(),
		helper.NormalizedJitterKey:     link.GetNormalizedUnidirDelayVariation(),
		helper.NormalizedPacketLossKey: link.GetNormalizedUnidirPacketLoss(),
	}
}

func (processor *NetworkProcessor) getOrCreateNode(nodeIgpRouterId string) graph.Node {
	// it's possible that link events are received before node events
	// so we need to ensure the node exists (maybe with temporary info) in the graph before adding the edge
	if processor.graph.NodeExists(nodeIgpRouterId) {
		return processor.graph.GetNode(nodeIgpRouterId)
	}
	node := processor.cache.GetNodeByIgpRouterId(nodeIgpRouterId)
	if node != nil {
		processor.log.Errorf("Node with igp router id %s not in cache - create it only in graph (will be added in cache by next AddLsNodeEvent)", nodeIgpRouterId)
		return processor.graph.AddNode(graph.NewNetworkNode(nodeIgpRouterId, node.GetName(), node.GetSrAlgorithm()))
	}
	return processor.graph.AddNode(graph.NewNetworkNode(nodeIgpRouterId, "", []uint32{}))
}

func (processor *NetworkProcessor) deleteEdge(key string) {
	processor.log.Debugln("Delete edge with key: ", key)
	if processor.graph.EdgeExists(key) {
		edge := processor.graph.GetEdge(key)
		if edge == nil {
			processor.log.Debugf("Edge with key %s does not exist in graph", key)
		}
		processor.log.Debugf("Delete edge with key %s from graph between %s and %s", key, edge.From().GetName(), edge.To().GetName())
		processor.graph.DeleteEdge(edge)
	} else {
		processor.log.Debugf("Edge with key %s does not exist in graph", key)
	}
}

func (processor *NetworkProcessor) addEdgeToGraph(edge graph.Edge) error {
	processor.log.Debugf("Add edge with key %s to graph between %s and %s", edge.GetId(), edge.From().GetName(), edge.To().GetName())
	processor.log.Debugf("Weights: %v", edge.GetAllWeights())
	if err := processor.graph.AddEdge(edge); err != nil {
		return err
	}
	return nil
}

func (processor *NetworkProcessor) addLinkToGraph(link domain.Link) error {
	key := link.GetKey()
	if !processor.graph.EdgeExists(key) {
		weights := processor.getCurrentLinkWeights(link)
		for weightKey, value := range weights {
			if value == 0 {
				return fmt.Errorf("Link contains zero values (%s), link %s is created during next update - ensure generic processor is running", weightKey, key)
			}
		}
		from := processor.getOrCreateNode(link.GetIgpRouterId())
		to := processor.getOrCreateNode(link.GetRemoteIgpRouterId())
		return processor.addEdgeToGraph(graph.NewNetworkEdge(key, from, to, weights))
	}
	return nil
}

func (processor *NetworkProcessor) setEdgeWeight(edge graph.Edge, key helper.WeightKey, value float64) error {

	if value == 0 && key != helper.NormalizedLatencyKey && key != helper.NormalizedJitterKey && key != helper.NormalizedPacketLossKey {
		return fmt.Errorf("Value is 0, not setting %s", key)
	}
	currentValue := edge.GetWeight(key)

	if currentValue != value {
		edge.SetWeight(key, value)
	}
	return nil
}

func (processor *NetworkProcessor) updateLinkInGraph(link domain.Link) error {
	key := link.GetKey()
	processor.log.Debugln("Updating link in graph with key: ", key)
	edge := processor.graph.GetEdge(key)
	if edge == nil {
		processor.log.Debugf("Link with key %s does not exist in graph, create it", key)
		return processor.addLinkToGraph(link)
	}
	for weightKey, weightValue := range processor.getCurrentLinkWeights(link) {
		if err := processor.setEdgeWeight(edge, weightKey, weightValue); err != nil {
			return err
		}
	}
	return nil
}

func (processor *NetworkProcessor) clearDuplicateAnnouncedPrefix(prefix domain.Prefix, networkAddress string, subnetLength uint8) {
	clientNetwork := processor.cache.GetClientNetworkByKey(prefix.GetKey())
	if clientNetwork == nil {
		processor.log.Debugf("Delete network %s/%d from cache since it's announced several times and thus not a client network: ", networkAddress, subnetLength)
		processor.prefixCounts[networkAddress]++
		processor.cache.RemoveClientNetwork(prefix)
	}
}

func (processor *NetworkProcessor) addNetworkToCache(prefix domain.Prefix, networkAddress string, subnetLength uint8) {
	processor.log.Debugf("Add network %s/%d to cache ", networkAddress, subnetLength)
	processor.prefixCounts[networkAddress] = 1
	processor.cache.StoreClientNetwork(prefix)
}

func (processor *NetworkProcessor) processPrefix(prefix domain.Prefix) {
	networkAddress := prefix.GetPrefix()
	subnetLength := prefix.GetPrefixLength()
	_, ok := processor.prefixCounts[networkAddress]
	if !ok {
		processor.addNetworkToCache(prefix, networkAddress, subnetLength)
	} else {
		processor.clearDuplicateAnnouncedPrefix(prefix, networkAddress, subnetLength)
	}
}

func (processor *NetworkProcessor) ProcessPrefixes(prefixes []domain.Prefix) {
	for _, prefix := range prefixes {
		processor.processPrefix(prefix)
	}
}

func (processor *NetworkProcessor) deleteClientNetwork(key string) error {
	prefix := processor.cache.GetClientNetworkByKey(key)
	if prefix == nil {
		return fmt.Errorf("Network with key %s does not exist in cache", key)
	}
	networkAddress := prefix.GetPrefix()
	subnetLength := prefix.GetPrefixLength()
	if _, ok := processor.prefixCounts[networkAddress]; !ok {
		return fmt.Errorf("Network %s/%d does not exist in prefix counts", networkAddress, subnetLength)
	}
	if processor.prefixCounts[networkAddress] > 1 {
		processor.log.Debugf("Decrement network %s/%d from announced prefix count", networkAddress, subnetLength)
		processor.prefixCounts[networkAddress]--
		return nil
	}
	processor.log.Debugf("Delete client network %s/%d from cache", networkAddress, subnetLength)
	delete(processor.prefixCounts, networkAddress)
	processor.cache.RemoveClientNetwork(prefix)
	return nil
}

func (processor *NetworkProcessor) addSidtoCache(sid domain.Sid) {
	processor.log.Debugf("Add SRv6 SID %s to cache", sid.GetSid())
	processor.cache.StoreSid(sid)
}

func (processor *NetworkProcessor) ProcessSids(sids []domain.Sid) {
	for _, sid := range sids {
		processor.addSidtoCache(sid)
	}
}

func (processor *NetworkProcessor) deleteSidFromCache(key string) {
	processor.log.Debugf("Delete SRv6 SID %s from cache", key)
	sid := processor.cache.GetSidByKey(key)
	if sid == nil {
		processor.log.Debugf("SID with key %s does not exist in cache", key)
	}
	processor.cache.RemoveSid(sid)
}

func (processor *NetworkProcessor) Start() {
	holdTime := time.Second * 3 // TODO make it configurable via env variable
	processor.log.Infof("Starting processing network updates with hold time %s", holdTime.String())

	timer := time.NewTimer(holdTime)
	defer timer.Stop()
	mutexesLocked := false

	for {
		select {
		case event := <-processor.eventChan:
			if !mutexesLocked {
				processor.log.Debugln("Locking cache and graph mutexes")
				processor.cache.Lock()
				processor.graph.Lock()
				mutexesLocked = true
			}
			processor.processEvent(event)
			timer.Reset(holdTime)
		case <-timer.C:
			if mutexesLocked {
				processor.log.Debugln("Unlocking cache and graph mutexes")
				processor.cache.Unlock()
				processor.graph.Unlock()
				mutexesLocked = false
			}
			if processor.needsSubgraphUpdate {
				processor.graph.UpdateSubGraphs()
				processor.needsSubgraphUpdate = false
			}
			processor.updateChan <- struct{}{}
		case <-processor.quitChan:
			if mutexesLocked {
				processor.cache.Unlock()
				processor.graph.Unlock()
			}
			return
		}
	}
}

func (processor *NetworkProcessor) processEvent(event domain.NetworkEvent) {
	switch eventType := event.(type) {
	case *domain.AddNodeEvent:
		processor.log.Debugln("Received AddNodeEvent: ", eventType.GetKey())
		processor.addOrUpdateNodeInGraphAndCache(eventType.Node)
		processor.needsSubgraphUpdate = true
	case *domain.UpdateNodeEvent:
		processor.log.Debugln("Received UpdateNodeEvent: ", eventType.GetKey())
		processor.addOrUpdateNodeInGraphAndCache(eventType.Node)
		processor.needsSubgraphUpdate = true
	case *domain.DeleteNodeEvent:
		processor.log.Debugln("Received DeleteNodeEvent: ", eventType.GetKey())
		processor.deleteNode(eventType.GetKey())
		processor.needsSubgraphUpdate = true
	case *domain.AddLinkEvent:
		processor.log.Debugln("Received AddLinkEvent: ", eventType.GetKey())
		if err := processor.addLinkToGraph(eventType.Link); err != nil {
			processor.log.Warnln("Error adding link to graph: ", err)
		}
		processor.needsSubgraphUpdate = true
	case *domain.UpdateLinkEvent:
		processor.log.Debugln("Received UpdateLinkEvent: ", eventType.GetKey())
		if err := processor.updateLinkInGraph(eventType.Link); err != nil {
			processor.log.Warnln("Error updating link in graph: ", err)
		}
	case *domain.DeleteLinkEvent:
		processor.log.Debugln("Received DeleteLinkEvent: ", eventType.GetKey())
		processor.deleteEdge(eventType.GetKey())
		processor.needsSubgraphUpdate = true
	case *domain.AddPrefixEvent:
		processor.log.Debugln("Received AddPrefixEvent: ", eventType.Prefix.GetKey())
		processor.processPrefix(eventType.Prefix)
	case *domain.DeletePrefixEvent:
		processor.log.Debugln("Received DeletePrefixEvent: ", eventType.GetKey())
		if err := processor.deleteClientNetwork(eventType.GetKey()); err != nil {
			processor.log.Warnln("Error deleting client network from cache: ", err)
		}
	case *domain.AddSidEvent:
		processor.log.Debugln("Received AddSidEvent: ", eventType.GetSid())
		processor.addSidtoCache(eventType.Sid)
	case *domain.DeleteSidEvent:
		processor.log.Debugln("Received DeleteSidEvent: ", eventType.GetKey())
		processor.deleteSidFromCache(eventType.GetKey())
	}
}

func (processor *NetworkProcessor) Stop() {
	processor.log.Infoln("Stopping processor")
	close(processor.quitChan)
}
