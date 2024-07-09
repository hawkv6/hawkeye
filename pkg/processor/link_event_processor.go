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

type LinkEventProcessor struct {
	log   *logrus.Entry
	graph graph.Graph
	cache cache.Cache
}

func NewLinkEventProcessor(graph graph.Graph, cache cache.Cache) *LinkEventProcessor {
	return &LinkEventProcessor{
		log:   logging.DefaultLogger.WithField("subsystem", Subsystem),
		graph: graph,
		cache: cache,
	}
}

func (processor *LinkEventProcessor) getCurrentLinkWeights(link domain.Link) map[helper.WeightKey]float64 {
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

func (processor *LinkEventProcessor) deleteEdge(key string) {
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

func (processor *LinkEventProcessor) addEdgeToGraph(edge graph.Edge) error {
	processor.log.Debugf("Add edge with key %s to graph between %s and %s", edge.GetId(), edge.From().GetName(), edge.To().GetName())
	processor.log.Debugf("Weights: %v", edge.GetAllWeights())
	if err := processor.graph.AddEdge(edge); err != nil {
		return err
	}
	return nil
}

func (processor *LinkEventProcessor) getOrCreateNode(nodeIgpRouterId string) graph.Node {
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

func (processor *LinkEventProcessor) addLinkToGraph(link domain.Link) error {
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

func (processor *LinkEventProcessor) setEdgeWeight(edge graph.Edge, key helper.WeightKey, value float64) error {

	if value == 0 && key != helper.NormalizedLatencyKey && key != helper.NormalizedJitterKey && key != helper.NormalizedPacketLossKey {
		return fmt.Errorf("Value is 0, not setting %s", key)
	}
	currentValue := edge.GetWeight(key)

	if currentValue != value {
		edge.SetWeight(key, value)
	}
	return nil
}

func (processor *LinkEventProcessor) updateLinkInGraph(link domain.Link) error {
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

func (processor *LinkEventProcessor) ProcessLinks(links []domain.Link) error {
	for _, link := range links {
		if err := processor.addLinkToGraph(link); err != nil {
			return err
		}
	}
	processor.graph.UpdateSubGraphs()
	return nil
}

func (processor *LinkEventProcessor) HandleEvent(event domain.NetworkEvent) bool {
	switch eventType := event.(type) {
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
		return false
	case *domain.DeleteLinkEvent:
		processor.log.Debugln("Received DeleteLinkEvent: ", eventType.GetKey())
		processor.deleteEdge(eventType.GetKey())
	default:
		processor.log.Warnf("Unknown event type: %T", eventType)
		return false
	}
	return true
}
