package jagw

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hawkv6/hawkeye/pkg/adapter"
	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/hawkv6/hawkeye/pkg/processor"
	"github.com/jalapeno-api-gateway/jagw-go/jagw"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DefaultJagwRequestService struct {
	log                  *logrus.Entry
	jagwRequestSocket    string
	grpcClientConnection *grpc.ClientConn
	requestClient        jagw.RequestServiceClient
	adapter              adapter.Adapter
	processor            processor.Processor
}

func NewDefaultJagwRequestService(config config.Config, adapter adapter.Adapter, processor processor.Processor) *DefaultJagwRequestService {
	return &DefaultJagwRequestService{
		log:               logging.DefaultLogger.WithField("subsystem", Subsystem),
		jagwRequestSocket: config.GetJagwServiceAddress() + ":" + strconv.FormatUint(uint64(config.GetJagwRequestPort()), 10),
		adapter:           adapter,
		processor:         processor,
	}
}

func (requestService *DefaultJagwRequestService) Start() error {
	requestService.log.Debugln("Initializing JAGW Request Service")
	grpcClientConnection, err := grpc.Dial(requestService.jagwRequestSocket,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	requestService.grpcClientConnection = grpcClientConnection
	requestService.requestClient = jagw.NewRequestServiceClient(grpcClientConnection)
	return nil
}

func (requestService *DefaultJagwRequestService) convertLsNodes(lsNodes []*jagw.LsNode) ([]domain.Node, error) {
	requestService.log.Debugln("Converting LsNodes to internal structure")
	var nodes []domain.Node
	for _, lsNode := range lsNodes {
		node, err := requestService.adapter.ConvertNode(lsNode)
		if err != nil {
			return nil, fmt.Errorf("Error converting LsNode: %s", err.Error())
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (requestService *DefaultJagwRequestService) GetLsNodes() error {
	request := &jagw.TopologyRequest{
		Keys:       []string{},
		Properties: []string{"Key", "IgpRouterId", "Name"},
	}

	requestService.log.Debugln("Getting LsNodes from JAGW")
	response, err := requestService.requestClient.GetLsNodes(context.Background(), request)
	if err != nil {
		return err
	}

	requestService.log.Infof("Got %d LsNodes from JAGW", len(response.LsNodes))
	nodes, err := requestService.convertLsNodes(response.LsNodes)
	if err != nil {
		return err
	}
	if err := requestService.processor.CreateNetworkNodes(nodes); err != nil {
		return err
	}
	return nil
}

func (requestService *DefaultJagwRequestService) convertLsLinks(lsLinks []*jagw.LsLink) ([]domain.Link, error) {
	requestService.log.Debugln("Converting LsLinks to internal structure")
	var links []domain.Link
	for _, lsLink := range lsLinks {
		defaultLink, err := requestService.adapter.ConvertLink(lsLink)
		if err != nil {
			return nil, fmt.Errorf("Error converting LsLink: %s", err.Error())
		}
		links = append(links, defaultLink)
	}
	return links, nil
}

func (requestService *DefaultJagwRequestService) GetLsLinks() error {
	request := &jagw.TopologyRequest{
		Keys:       []string{},
		Properties: []string{"Key", "IgpRouterId", "RemoteIgpRouterId", "UnidirLinkDelay", "UnidirDelayVariation", "UnidirAvailableBW", "UnidirPacketLoss", "UnidirBWUtilization"},
	}

	requestService.log.Debugln("Getting LsLinks from JAGW")
	response, err := requestService.requestClient.GetLsLinks(context.Background(), request)
	if err != nil {
		return err
	}

	requestService.log.Infof("Got %d LsLinks from JAGW", len(response.LsLinks))
	links, err := requestService.convertLsLinks(response.LsLinks)
	if err != nil {
		return err
	}

	if err := requestService.processor.CreateNetworkEdges(links); err != nil {
		return err
	}
	return nil
}

func (requestService *DefaultJagwRequestService) convertLsPrefix(lsPrefixes []*jagw.LsPrefix) ([]domain.Prefix, error) {
	requestService.log.Debugf("Converting LsPrefixes to internal structure")
	var prefixes []domain.Prefix
	for _, lsPrefix := range lsPrefixes {
		prefix, err := requestService.adapter.ConvertPrefix(lsPrefix)
		if err != nil {
			return nil, fmt.Errorf("Error converting LsPrefix: %s", err.Error())
		}
		prefixes = append(prefixes, prefix)
	}
	return prefixes, nil
}

func (requestService *DefaultJagwRequestService) GetLsPrefixes() error {
	request := &jagw.TopologyRequest{
		Keys:       []string{},
		Properties: []string{"Key", "IgpRouterId", "Prefix", "PrefixLen"},
	}

	requestService.log.Debugln("Getting LsPrefixes from JAGW")
	response, err := requestService.requestClient.GetLsPrefixes(context.Background(), request)
	if err != nil {
		return err
	}

	requestService.log.Infof("Got %d LsPrefixes from JAGW", len(response.LsPrefixes))
	prefixes, err := requestService.convertLsPrefix(response.LsPrefixes)
	if err != nil {
		return err
	}

	if err := requestService.processor.CreateClientNetworks(prefixes); err != nil {
		return err
	}
	return nil
}

func (requestService *DefaultJagwRequestService) convertLsSrv6Sids(lsSrv6Sids []*jagw.LsSrv6Sid) ([]domain.Sid, error) {
	requestService.log.Debugln("Converting LsSrv6Sids to internal structure")
	var sidList []domain.Sid
	for _, lsSrv6Sid := range lsSrv6Sids {
		srv6Sid, err := requestService.adapter.ConvertSid(lsSrv6Sid)
		if err != nil {
			return nil, fmt.Errorf("Error converting LsSrv6Sid: %s", err.Error())
		}
		sidList = append(sidList, srv6Sid)
	}
	return sidList, nil
}

func (requestService *DefaultJagwRequestService) GetSrv6Sids() error {
	request := &jagw.TopologyRequest{
		Keys:       []string{},
		Properties: []string{"Key", "IgpRouterId", "Srv6Sid"},
	}

	requestService.log.Debugln("Getting SRv6 SIDs from JAGW")
	response, err := requestService.requestClient.GetLsSrv6Sids(context.Background(), request)
	if err != nil {
		return err
	}
	requestService.log.Infof("Got %d SRv6 SIDs from JAGW", len(response.LsSrv6Sids))
	sidList, err := requestService.convertLsSrv6Sids(response.LsSrv6Sids)
	if err != nil {
		return err
	}
	if err := requestService.processor.CreateSids(sidList); err != nil {
		return err
	}
	return nil
}

func (requestService *DefaultJagwRequestService) Stop() {
	requestService.log.Infoln("Closing JAGW Request Service")
	requestService.grpcClientConnection.Close()
}
