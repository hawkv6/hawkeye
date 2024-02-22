package jagw

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hawkv6/hawkeye/pkg/adapter"
	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
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
	adapter              adapter.Adapter
	processor            processor.Processor
}

func NewDefaultJagwRequestService(config config.Config, adapter adapter.Adapter, processor processor.Processor) *DefaultJagwRequestService {
	return &DefaultJagwRequestService{
		log:               logging.DefaultLogger.WithField("subsystem", Subsystem),
		jagwRequestSocket: config.GetJagwServiceAddress() + ":" + strconv.FormatUint(uint64(config.GetJagwSubscriptionPort()), 10),
		adapter:           adapter,
		processor:         processor,
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

func (requestService *DefaultJagwRequestService) convertLsLinks(lsLinks []*jagw.LsLink) ([]domain.Link, error) {
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

func (requestService *DefaultJagwRequestService) GetLsLinks(network graph.Graph) error {
	client := jagw.NewRequestServiceClient(requestService.grpcClientConnection)
	request := &jagw.TopologyRequest{
		Keys:       []string{},
		Properties: []string{"Key", "IgpRouterId", "RemoteIgpRouterId", "UnidirLinkDelay"},
	}
	response, err := client.GetLsLinks(context.Background(), request)
	if err != nil {
		return err
	}
	links, err := requestService.convertLsLinks(response.LsLinks)
	if err != nil {
		return err
	}

	if err := requestService.processor.CreateNetworkGraph(links); err != nil {
		return err
	}
	return nil
}

func (requestService *DefaultJagwRequestService) convertLsPrefix(lsPrefixes []*jagw.LsPrefix) ([]domain.Prefix, error) {
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
	client := jagw.NewRequestServiceClient(requestService.grpcClientConnection)
	request := &jagw.TopologyRequest{
		Keys:       []string{},
		Properties: []string{"Key", "IgpRouterId", "Prefix", "PrefixLen"},
	}

	response, err := client.GetLsPrefixes(context.Background(), request)
	if err != nil {
		return err
	}

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
	client := jagw.NewRequestServiceClient(requestService.grpcClientConnection)
	request := &jagw.TopologyRequest{
		Keys:       []string{},
		Properties: []string{"Key", "IgpRouterId", "Srv6Sid"},
	}

	response, err := client.GetLsSrv6Sids(context.Background(), request)
	if err != nil {
		return err
	}

	sidList, err := requestService.convertLsSrv6Sids(response.LsSrv6Sids)
	if err != nil {
		return err
	}
	if err := requestService.processor.CreateSids(sidList); err != nil {
		return err
	}
	return nil
}

func (requestService *DefaultJagwRequestService) Close() {
	requestService.grpcClientConnection.Close()
}
