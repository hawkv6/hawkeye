package jagw

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hawkv6/hawkeye/pkg/adapter"
	"github.com/hawkv6/hawkeye/pkg/config"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/hawkv6/hawkeye/pkg/processor"
	"github.com/jalapeno-api-gateway/jagw-go/jagw"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type JagwRequestService struct {
	log                  *logrus.Entry
	jagwRequestSocket    string
	grpcClientConnection *grpc.ClientConn
	requestClient        jagw.RequestServiceClient
	adapter              adapter.Adapter
	processor            processor.Processor
	helper               helper.Helper
}

func NewJagwRequestService(config config.Config, adapter adapter.Adapter, processor processor.Processor, helper helper.Helper) *JagwRequestService {
	return &JagwRequestService{
		log:               logging.DefaultLogger.WithField("subsystem", Subsystem),
		jagwRequestSocket: config.GetJagwServiceAddress() + ":" + strconv.FormatUint(uint64(config.GetJagwRequestPort()), 10),
		adapter:           adapter,
		processor:         processor,
		helper:            helper,
	}
}

func (requestService *JagwRequestService) Init() error {
	requestService.log.Debugln("Initializing JAGW Request Service")
	grpcClientConnection, err := grpc.NewClient(requestService.jagwRequestSocket,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	requestService.grpcClientConnection = grpcClientConnection
	requestService.requestClient = jagw.NewRequestServiceClient(grpcClientConnection)
	return nil
}

func (requestService *JagwRequestService) Start() error {
	if err := requestService.getSrv6Sids(); err != nil {
		return fmt.Errorf("Error getting SRv6 SIDs from JAGW: %v", err)
	}
	if err := requestService.getLsPrefixes(); err != nil {
		return fmt.Errorf("Error getting LsPrefixes from JAGW: %v", err)
	}
	if err := requestService.getLsNodes(); err != nil {
		return fmt.Errorf("Error getting LsNodes from JAGW: %v", err)
	}
	if err := requestService.getLsLinks(); err != nil {
		return fmt.Errorf("Error getting LsLinks from JAGW: %v", err)
	}
	return nil
}

func (requestService *JagwRequestService) convertLsNodes(lsNodes []*jagw.LsNode) ([]domain.Node, error) {
	requestService.log.Debugln("Converting LsNodes to internal structure")
	nodes := make([]domain.Node, len(lsNodes))
	for i, lsNode := range lsNodes {
		node, err := requestService.adapter.ConvertNode(lsNode)
		if err != nil {
			return nil, fmt.Errorf("Error converting LsNode: %s", err.Error())
		}
		nodes[i] = node
	}
	return nodes, nil
}

func (requestService *JagwRequestService) getLsNodes() error {
	request := &jagw.TopologyRequest{
		Properties: requestService.helper.GetLsNodeProperties(),
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
	requestService.processor.ProcessNodes(nodes)
	return nil
}

func (requestService *JagwRequestService) convertLsLinks(lsLinks []*jagw.LsLink) ([]domain.Link, error) {
	requestService.log.Debugln("Converting LsLinks to internal structure")
	links := make([]domain.Link, len(lsLinks))
	for i, lsLink := range lsLinks {
		link, err := requestService.adapter.ConvertLink(lsLink)
		if err != nil {
			return nil, fmt.Errorf("Error converting LsLink: %s", err.Error())
		}
		links[i] = link
	}
	return links, nil
}

func (requestService *JagwRequestService) getLsLinks() error {
	request := &jagw.TopologyRequest{
		Properties: requestService.helper.GetLsLinkProperties(),
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

	if err := requestService.processor.ProcessLinks(links); err != nil {
		return err
	}
	return nil
}

func (requestService *JagwRequestService) shouldSkipPrefix(lsPrefix *jagw.LsPrefix) bool {
	if lsPrefix.Srv6Locator != nil {
		requestService.log.Debugf("Skip prefix %s because it has a SRv6 locator TLV", lsPrefix.GetKey())
		return true
	}
	if lsPrefix.PrefixLen != nil && *lsPrefix.PrefixLen == 128 {
		requestService.log.Debugf("Skip prefix %s because it has a prefix length of 128 (belongs to Loopback0)", lsPrefix.GetKey())
		return true
	}
	return false
}

func (requestService *JagwRequestService) convertLsPrefix(lsPrefixes []*jagw.LsPrefix) ([]domain.Prefix, error) {
	requestService.log.Debugf("Converting LsPrefixes to internal structure")
	var prefixes []domain.Prefix
	for _, lsPrefix := range lsPrefixes {
		if requestService.shouldSkipPrefix(lsPrefix) {
			continue
		}
		prefix, err := requestService.adapter.ConvertPrefix(lsPrefix)
		if err != nil {
			return nil, fmt.Errorf("Error converting LsPrefix: %s", err.Error())
		}
		prefixes = append(prefixes, prefix)
	}
	return prefixes, nil
}

func (requestService *JagwRequestService) getLsPrefixes() error {
	request := &jagw.TopologyRequest{
		Properties: requestService.helper.GetLsPrefixProperties(),
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

	requestService.processor.ProcessPrefixes(prefixes)
	return nil
}

func (requestService *JagwRequestService) convertLsSrv6Sids(lsSrv6Sids []*jagw.LsSrv6Sid) ([]domain.Sid, error) {
	requestService.log.Debugln("Converting LsSrv6Sids to internal structure")
	sidList := make([]domain.Sid, len(lsSrv6Sids))
	for i, lsSrv6Sid := range lsSrv6Sids {
		srv6Sid, err := requestService.adapter.ConvertSid(lsSrv6Sid)
		if err != nil {
			return nil, fmt.Errorf("Error converting LsSrv6Sid: %s", err.Error())
		}
		sidList[i] = srv6Sid
	}
	return sidList, nil
}

func (requestService *JagwRequestService) getSrv6Sids() error {
	request := &jagw.TopologyRequest{
		Properties: requestService.helper.GetLsSrv6SidsProperties(),
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
	requestService.processor.ProcessSids(sidList)
	return nil
}

func (requestService *JagwRequestService) Stop() {
	requestService.log.Infoln("Closing JAGW Request Service")
	requestService.grpcClientConnection.Close()
}
