package jagw

import (
	"context"
	"strconv"

	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/jalapeno-api-gateway/jagw-go/jagw"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DefaultJagwRequestService struct {
	log                  *logrus.Entry
	jagwRequestSocket    string
	grpcClientConnection *grpc.ClientConn
}

func NewDefaultJagwRequestService(config config.Config) *DefaultJagwRequestService {
	return &DefaultJagwRequestService{
		log:               logging.DefaultLogger.WithField("subsystem", Subsystem),
		jagwRequestSocket: config.GetJagwServiceAddress() + ":" + strconv.FormatUint(uint64(config.GetJagwSubscriptionPort()), 10),
	}
}

func (requestService *DefaultJagwRequestService) Init() error {
	requestService.log.Debugln("Initializing JAGW Request Service")
	grpcClientConnection, err := grpc.Dial(requestService.jagwRequestSocket,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	requestService.grpcClientConnection = grpcClientConnection
	return nil
}

func (requestService *DefaultJagwRequestService) GetLsLinks(network graph.Graph) error {
	client := jagw.NewRequestServiceClient(requestService.grpcClientConnection)
	defer requestService.grpcClientConnection.Close()
	request := &jagw.TopologyRequest{
		Keys:       []string{},
		Properties: []string{},
	}
	response, err := client.GetLsLinks(context.Background(), request)
	if err != nil {
		return err
	}
	for _, lsLink := range response.LsLinks {
		var from graph.Node
		if !network.NodeExists(*lsLink.IgpRouterId) {
			from = graph.NewDefaultNode(*lsLink.IgpRouterId)
			if err := network.AddNode(from); err != nil {
				return err
			}
		} else {
			if from, err = network.GetNode(*lsLink.IgpRouterId); err != nil {
				return err
			}
		}

		var to graph.Node
		if !network.NodeExists(*lsLink.RemoteIgpRouterId) {
			to = graph.NewDefaultNode(*lsLink.RemoteIgpRouterId)
			if err := network.AddNode(to); err != nil {
				return err
			}
		} else {
			if to, err = network.GetNode(*lsLink.RemoteIgpRouterId); err != nil {
				return err
			}
		}
		edge := graph.NewDefaultEdge(lsLink.Id, from, to, map[string]float64{"delay": float64(*lsLink.UnidirLinkDelay)})
		if err := network.AddEdge(edge); err != nil {
			return err
		}
	}
	return nil
}
