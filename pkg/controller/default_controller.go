package controller

import (
	"fmt"
	"net"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/hawkv6/hawkeye/pkg/messaging"
	"github.com/sirupsen/logrus"
)

type DefaultController struct {
	log             *logrus.Entry
	cache           cache.CacheService
	network         graph.Graph
	streamSessions  []domain.StreamSession
	pathRequestChan chan domain.PathRequest
	pathResultChan  chan domain.PathResult
}

func NewDefaultController(cache cache.CacheService, network graph.Graph, messagingChannels messaging.MessagingChannels) *DefaultController {
	return &DefaultController{
		log:             logging.DefaultLogger.WithField("subsystem", Subsystem),
		cache:           cache,
		network:         network,
		streamSessions:  make([]domain.StreamSession, 0),
		pathRequestChan: messagingChannels.GetPathRequestChan(),
		pathResultChan:  messagingChannels.GetPathResponseChan(),
	}
}

func getNetworkAddress(ip string) net.IP {
	ipv6Address := net.ParseIP(ip)
	mask := net.CIDRMask(64, 128)
	return ipv6Address.Mask(mask)
}

func (controller *DefaultController) getShortestPath(sourceIP, destinationIp, metric string) ([]string, error) {
	sourceIpv6 := getNetworkAddress(sourceIP)
	destinationIpv6 := getNetworkAddress(destinationIp)
	sourceRouterId, ok := controller.cache.GetRouterIdFromNetworkAddress(sourceIpv6.String())
	if !ok {
		return nil, fmt.Errorf("Router ID not found for source IP: %s", sourceRouterId)
	}
	destinationRouterId, ok := controller.cache.GetRouterIdFromNetworkAddress(destinationIpv6.String())
	if !ok {
		return nil, fmt.Errorf("Router ID not found for destination IP: %s", destinationRouterId)
	}
	controller.log.Debugln("Source Router ID: ", sourceRouterId)
	controller.log.Debugln("Destination Router ID: ", destinationRouterId)

	source, err := controller.network.GetNode(sourceRouterId)
	if err != nil {
		return nil, err
	}
	destination, err := controller.network.GetNode(destinationRouterId)
	if err != nil {
		return nil, err
	}
	result, err := controller.network.GetShortestPath(source, destination, metric)
	if err != nil {
		return nil, err
	}
	edgeList := make([]string, 0)
	for _, edge := range result {
		// from := edge.From().GetId().(string)
		to := edge.To().GetId().(string)
		// edgeList = append(edgeList, from, to)
		edgeList = append(edgeList, to)
	}
	return edgeList, nil
}

func (controller *DefaultController) Start() {
	controller.log.Infoln("Starting controller")
	for {
		pathRequest := <-controller.pathRequestChan
		controller.log.Debugln("Received path request: ", pathRequest)

		ipv6SourceAddress := pathRequest.GetIpv6SourceAddress()
		ipv6DestinationAddress := pathRequest.GetIpv6DestinationAddress()

		// TODO: change logic if there are multiple intents
		for _, intent := range pathRequest.GetIntents() {
			switch intentType := intent.GetIntentType(); intentType {
			case domain.IntentTypeLowLatency:
				nodeIds, err := controller.getShortestPath(ipv6SourceAddress, ipv6DestinationAddress, "delay")
				if err != nil {
					controller.log.Errorln("Error getting shortest path: ", err)
					continue
				}
				for i := 0; i < len(nodeIds); i = i + 2 {
					controller.log.Infof("Edge: %s -> %s", nodeIds[i], nodeIds[i+1])
				}
				sidList := make([]string, 0)
				for _, node := range nodeIds {
					sid, ok := controller.cache.GetSidFromRouterId(node)
					if !ok {
						controller.log.Errorln("SID not found for router: ", node)
						continue
					}
					sidList = append(sidList, sid)
				}
				controller.log.Debugln("SID List: ", sidList)
				pathResult, err := domain.NewDefaultPathResult(pathRequest, sidList)
				if err != nil {
					controller.log.Errorln("Error creating path result: ", err)
					continue
				}
				controller.pathResultChan <- pathResult
			default:
				controller.log.Errorln("Not (yet): ", intentType)
			}
		}
	}
}
