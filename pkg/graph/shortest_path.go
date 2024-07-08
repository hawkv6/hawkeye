package graph

import (
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type ShortestPath struct {
	log              *logrus.Entry
	edges            []Edge
	totalCost        float64
	totalDelay       float64
	totalJitter      float64
	totalPacketLoss  float64
	bottleneckEdge   Edge
	bottleneckValue  float64
	routerServiceMap map[string]string
}

func NewShortestPath(edges []Edge, totalCost, delay, jitter, packetLoss, bottleNeckValue float64, bottleneckEdge Edge) *ShortestPath {
	return &ShortestPath{
		log:             logging.DefaultLogger.WithField("subsystem", Subsystem),
		edges:           edges,
		totalCost:       totalCost,
		totalDelay:      delay,
		totalJitter:     jitter,
		totalPacketLoss: packetLoss,
		bottleneckEdge:  bottleneckEdge,
		bottleneckValue: bottleNeckValue,
	}
}

func (path *ShortestPath) GetEdges() []Edge {
	return path.edges
}

func (path *ShortestPath) GetTotalCost() float64 {
	return path.totalCost
}

func (path *ShortestPath) SetTotalCost(cost float64) {
	path.totalCost = cost
}

func (path *ShortestPath) GetTotalDelay() float64 {
	return path.totalDelay
}

func (path *ShortestPath) SetTotalDelay(delay float64) {
	path.totalDelay = delay
}

func (path *ShortestPath) GetTotalJitter() float64 {
	return path.totalJitter
}

func (path *ShortestPath) SetTotalJitter(jitter float64) {
	path.totalJitter = jitter
}

func (path *ShortestPath) GetTotalPacketLoss() float64 {
	return path.totalPacketLoss
}

func (path *ShortestPath) SetTotalPacketLoss(packetLoss float64) {
	path.totalPacketLoss = packetLoss
}

func (path *ShortestPath) GetBottleneckEdge() Edge {
	return path.bottleneckEdge
}

func (path *ShortestPath) GetBottleneckValue() float64 {
	return path.bottleneckValue
}

func (path *ShortestPath) SetBottleneckEdge(edge Edge) {
	path.bottleneckEdge = edge
}

func (path *ShortestPath) SetBottleneckValue(value float64) {
	path.bottleneckValue = value
}

func (path *ShortestPath) SetRouterServiceMap(routerServiceMap map[string]string) {
	path.routerServiceMap = routerServiceMap
}

func (path *ShortestPath) GetRouterServiceMap() map[string]string {
	return path.routerServiceMap
}
