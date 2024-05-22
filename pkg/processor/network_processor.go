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
	log               *logrus.Entry
	graph             graph.Graph
	cache             cache.Cache
	eventChan         chan domain.NetworkEvent
	quitChan          chan struct{}
	prefixCounts      map[string]int
	updateChan        chan struct{}
	latencyMetrics    []float64
	jitterMetrics     []float64
	packetLossMetrics []float64
}

func NewNetworkProcessor(graph graph.Graph, cache cache.Cache, eventChan chan domain.NetworkEvent, helper helper.Helper, updateChan chan struct{}) *NetworkProcessor {
	return &NetworkProcessor{
		log:               logging.DefaultLogger.WithField("subsystem", Subsystem),
		graph:             graph,
		cache:             cache,
		eventChan:         eventChan,
		quitChan:          make(chan struct{}),
		prefixCounts:      make(map[string]int),
		updateChan:        updateChan,
		latencyMetrics:    []float64{},
		jitterMetrics:     []float64{},
		packetLossMetrics: []float64{},
	}
}

func (processor *NetworkProcessor) setNodeName(node domain.Node, graphNode graph.Node) {
	nodeName := node.GetName()
	if graphNode.GetName() != nodeName {
		graphNode.SetName(nodeName)
	}
}

func (processor *NetworkProcessor) addNodeToGraph(node domain.Node) error {
	id := node.GetIgpRouterId()
	graphNode, exist := processor.graph.GetNode(id)
	if exist {
		processor.log.Debugf("Node with id %s already exists in graph", id)
	} else {
		processor.log.Debugf("Add node %s to graph with igp router id %s", node.GetName(), id)
		var err error
		graphNode, err = processor.graph.AddNode(graph.NewNetworkNode(id))
		if err != nil {
			return err
		}
	}
	processor.setNodeName(node, graphNode)
	return nil
}

func (processor *NetworkProcessor) addNodeToCache(node domain.Node) {
	processor.log.Debugf("Add node %s to cache with igp router id id %s", node.GetName(), node.GetIgpRouterId())
	processor.cache.StoreNode(node)
}

func (processor *NetworkProcessor) deleteNode(key string) error {
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

func (processor *NetworkProcessor) CreateGraphNodes(nodes []domain.Node) error {
	for _, node := range nodes {
		if err := processor.addNodeToGraph(node); err != nil {
			return err
		}
		processor.addNodeToCache(node)
	}
	return nil
}

func (processor *NetworkProcessor) CreateGraphEdges(links []domain.Link) error {
	for _, link := range links {
		if err := processor.addLinkToGraph(link); err != nil {
			return err
		}
	}
	return nil
}

func (processor *NetworkProcessor) getLinkWeights(link domain.Link) map[helper.WeightKey]float64 {
	return map[helper.WeightKey]float64{
		helper.LatencyKey:            float64(link.GetUnidirLinkDelay()),
		helper.JitterKey:             float64(link.GetUnidirDelayVariation()),
		helper.MaximumLinkBandwidth:  float64(link.GetMaxLinkBWKbps()),
		helper.AvailableBandwidthKey: float64(link.GetUnidirAvailableBandwidth()),
		helper.UtilizedBandwidthKey:  float64(link.GetUnidirBandwidthUtilization()),
		helper.PacketLossKey:         float64(link.GetUnidirPacketLoss()),
	}
}

func (processor *NetworkProcessor) getOrCreateNode(nodeId string) (graph.Node, error) {
	// it's possible that link events are received before node events
	// so we need to ensure the node exists in the graph before adding the edge
	if processor.graph.NodeExists(nodeId) {
		node, _ := processor.graph.GetNode(nodeId)
		return node, nil
	}
	node, exists := processor.cache.GetNodeByIgpRouterId(nodeId)
	if !exists {
		processor.log.Errorf("Node with igp router id %s not in cache - create it only in graph", nodeId)
		return processor.graph.AddNode(graph.NewNetworkNode(nodeId))
	}
	if err := processor.addNodeToGraph(node); err != nil {
		return nil, err
	}
	graphNode, _ := processor.graph.GetNode(nodeId)
	return graphNode, nil
}

func (processor *NetworkProcessor) deleteEdge(key string) error {
	if processor.graph.EdgeExists(key) {
		edge, exist := processor.graph.GetEdge(key)
		if !exist {
			return fmt.Errorf("Edge with key %s does not exist in graph", key)
		}
		processor.log.Debugf("Delete edge with key %s from graph between %s and %s", key, edge.From().GetName(), edge.To().GetName())
		processor.graph.DeleteEdge(edge)
		return nil
	}
	return fmt.Errorf("Edge with key %s does not exist in graph", key)
}

func (processor *NetworkProcessor) addEdgeToGraph(edge graph.Edge) error {
	processor.log.Debugf("Add edge with key %s to graph between %s and %s", edge.GetId(), edge.From().GetName(), edge.To().GetName())
	if err := processor.graph.AddEdge(edge); err != nil {
		return err
	}
	return nil
}

func (processor *NetworkProcessor) GetLatencyMetrics() []float64 {
	return processor.latencyMetrics
}

func (processor *NetworkProcessor) GetJitterMetrics() []float64 {
	return processor.jitterMetrics
}

func (processor *NetworkProcessor) GetPacketLossMetrics() []float64 {
	return processor.packetLossMetrics
}

func (processor *NetworkProcessor) addWeightToMetrics(weightKey helper.WeightKey, value float64) {
	if weightKey == helper.LatencyKey {
		processor.latencyMetrics = append(processor.latencyMetrics, value)
	} else if weightKey == helper.JitterKey {
		processor.jitterMetrics = append(processor.jitterMetrics, value)
	} else if weightKey == helper.PacketLossKey {
		processor.packetLossMetrics = append(processor.packetLossMetrics, value)
	}
}

func (processor *NetworkProcessor) addLinkToGraph(link domain.Link) error {
	key := link.GetKey()
	if !processor.graph.EdgeExists(key) {
		weights := processor.getLinkWeights(link)
		for weightKey, value := range weights {
			if value == 0 {
				return fmt.Errorf("Link contains zero values (%s), link %s is created during next update", weightKey, key)
			}
			processor.addWeightToMetrics(weightKey, value)
		}
		from, err := processor.getOrCreateNode(link.GetIgpRouterId())
		if err != nil {
			return err
		}

		to, err := processor.getOrCreateNode(link.GetRemoteIgpRouterId())
		if err != nil {
			return err
		}

		return processor.addEdgeToGraph(graph.NewNetworkEdge(key, from, to, weights))
	}
	return nil
}

func (processor *NetworkProcessor) setEdgeWeight(edge graph.Edge, key helper.WeightKey, value float64) error {
	if value == 0 {
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
	return nil
}

func (processor *NetworkProcessor) addClientNetworkToCache(prefix domain.Prefix) {
	networkAddress := prefix.GetPrefix()
	subnetLength := prefix.GetPrefixLength()
	if subnetLength == 128 {
		processor.log.Debugf("Ignoring lo0 address %s/%d", networkAddress, subnetLength)
	}
	_, ok := processor.prefixCounts[networkAddress]
	if !ok {
		processor.log.Debugf("Add network %s/%d to cache ", networkAddress, subnetLength)
		processor.prefixCounts[networkAddress] = 1
		processor.cache.StoreClientNetwork(prefix)
	} else {
		if _, ok := processor.cache.GetClientNetworkByKey(prefix.GetKey()); ok {
			processor.log.Debugf("Delete non-client network %s/%d from cache ", networkAddress, subnetLength)
			processor.prefixCounts[networkAddress]++
			processor.cache.RemoveClientNetwork(prefix)
		}
	}
}

func (processor *NetworkProcessor) CreateClientNetworks(prefixes []domain.Prefix) {
	for _, prefix := range prefixes {
		processor.addClientNetworkToCache(prefix)
	}
}

func (processor *NetworkProcessor) deleteClientNetwork(key string) error {
	prefix, ok := processor.cache.GetClientNetworkByKey(key)
	if !ok {
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

func (processor *NetworkProcessor) CreateSids(sids []domain.Sid) {
	for _, sid := range sids {
		processor.addSidtoCache(sid)
	}
}

func (processor *NetworkProcessor) deleteSidFromCache(key string) {
	processor.log.Debugf("Delete SRv6 SID %s from cache", key)
	sid, ok := processor.cache.GetSidByKey(key)
	if !ok {
		processor.log.Debugf("SID with key %s does not exist in cache", key)
	}
	processor.cache.RemoveSid(sid)
}

func (processor *NetworkProcessor) Start() {
	holdTime := time.Second * 3
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
		if err := processor.addNodeToGraph(eventType.Node); err != nil {
			processor.log.Warnln("Error creating network node: ", err)
		}
		processor.addNodeToCache(eventType.Node)
	case *domain.DeleteNodeEvent:
		processor.log.Debugln("Received DeleteNodeEvent: ", eventType.GetKey())
		if err := processor.deleteNode(eventType.GetKey()); err != nil {
			processor.log.Warnln("Error deleting network node: ", err)
		}
	case *domain.AddLinkEvent:
		processor.log.Debugln("Received AddLinkEvent: ", eventType.GetKey())
		if err := processor.addLinkToGraph(eventType.Link); err != nil {
			processor.log.Warnln("Error adding link to graph: ", err)
		}
	case *domain.UpdateLinkEvent:
		processor.log.Debugln("Received UpdateLinkEvent: ", eventType.GetKey())
		if err := processor.updateLinkInGraph(eventType.Link); err != nil {
			processor.log.Warnln("Error updating link in graph: ", err)
		}
	case *domain.DeleteLinkEvent:
		processor.log.Debugln("Received DeleteLinkEvent: ", eventType.GetKey())
		if err := processor.deleteEdge(eventType.GetKey()); err != nil {
			processor.log.Warnln("Error deleting edge: ", err)
		}
	case *domain.AddPrefixEvent:
		processor.log.Debugln("Received AddPrefixEvent: ", eventType.Prefix.GetKey())
		processor.addClientNetworkToCache(eventType.Prefix)
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
