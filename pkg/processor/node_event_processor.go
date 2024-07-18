package processor

import (
	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type NodeEventProcessor struct {
	log   *logrus.Entry
	graph graph.Graph
	cache cache.Cache
}

func NewNodeEventProcessor(graph graph.Graph, cache cache.Cache) *NodeEventProcessor {
	return &NodeEventProcessor{
		log:   logging.DefaultLogger.WithField("subsystem", Subsystem),
		graph: graph,
		cache: cache,
	}
}

func (processor *NodeEventProcessor) updateNode(node graph.Node, name string, srAlgorithm []uint32) {
	processor.log.Debugf("Update node %s with igp router id %s in graph", name, node.GetId())
	node.SetName(name)
	node.SetFlexibleAlgorithms(srAlgorithm)
}

func (processor *NodeEventProcessor) updateNodeInGraph(igpRouterId, name string, srAlgorithm []uint32) {
	if processor.graph.NodeExists(igpRouterId) {
		processor.log.Debugf("Node with id %s already exists in graph", igpRouterId)
		processor.updateNode(processor.graph.GetNode(igpRouterId), name, srAlgorithm)
	} else {
		processor.graph.AddNode(graph.NewNetworkNode(igpRouterId, name, srAlgorithm))
	}
}

func (processor *NodeEventProcessor) addNodeToCache(node domain.Node) {
	processor.cache.StoreNode(node)
}

func (processor *NodeEventProcessor) removeNodeFromGraph(node domain.Node) {
	igpRouterId := node.GetIgpRouterId()
	if processor.graph.NodeExists(igpRouterId) {
		graphNode := processor.graph.GetNode(igpRouterId)
		processor.graph.DeleteNode(graphNode)
	} else {
		name := node.GetName()
		processor.log.Debugf("Node with igp router id %s and name %s does not exist in graph", igpRouterId, name)
	}
}

func (processor *NodeEventProcessor) deleteNode(key string) {
	processor.log.Debugf("Delete node with key %s", key)
	node := processor.cache.GetNodeByKey(key)
	if node != nil {
		processor.cache.RemoveNode(node)
		processor.removeNodeFromGraph(node)
	} else {
		processor.log.Debugf("Node with key %s does not exist in cache and thus also not in graph", key)
	}
}

func (processor *NodeEventProcessor) addOrUpdateNodeInGraphAndCache(node domain.Node) {
	igpRouterId := node.GetIgpRouterId()
	name := node.GetName()
	srAlgorithm := node.GetSrAlgorithm()
	processor.log.Debugf("Add node %s with igp router id %s and SR algorithm %v to graph and cache", name, igpRouterId, srAlgorithm)
	processor.updateNodeInGraph(igpRouterId, name, srAlgorithm)
	processor.addNodeToCache(node)

}

func (processor *NodeEventProcessor) ProcessNodes(nodes []domain.Node) {
	for _, node := range nodes {
		processor.addOrUpdateNodeInGraphAndCache(node)
	}
}

func (processor *NodeEventProcessor) HandleEvent(event domain.NetworkEvent) bool {
	switch eventType := event.(type) {
	case *domain.AddNodeEvent:
		processor.log.Debugln("Received AddNodeEvent: ", eventType.GetKey())
		processor.addOrUpdateNodeInGraphAndCache(eventType.Node)
	case *domain.UpdateNodeEvent:
		processor.log.Debugln("Received UpdateNodeEvent: ", eventType.GetKey())
		processor.addOrUpdateNodeInGraphAndCache(eventType.Node)
	case *domain.DeleteNodeEvent:
		processor.log.Debugln("Received DeleteNodeEvent: ", eventType.GetKey())
		processor.deleteNode(eventType.GetKey())
	default:
		processor.log.Warnf("Unknown event type: %T", eventType)
		return false
	}
	return true
}
