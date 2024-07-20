package calculation

import (
	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type CalculationTransformerService struct {
	log   *logrus.Entry
	cache cache.Cache
}

func NewCalculationTransformerService(cache cache.Cache) *CalculationTransformerService {
	return &CalculationTransformerService{
		log:   logging.DefaultLogger.WithField("subsystem", subsystem),
		cache: cache,
	}
}

func (service *CalculationTransformerService) getNodesFromPath(path graph.Path) []string {
	nodeList := make([]string, 0)
	for _, edge := range path.GetEdges() {
		to := edge.To().GetId()
		nodeList = append(nodeList, to)
	}
	return nodeList
}

func (service *CalculationTransformerService) translatePathToSidList(path graph.Path, algorithm uint32) ([]string, []string) {
	nodeList := service.getNodesFromPath(path)
	serviceSidList := make([]string, 0)
	routerServiceMap := path.GetRouterServiceMap()
	service.log.Debugln("Node in Path: ", nodeList)
	var sidList []string
	for _, node := range nodeList {
		sid := service.cache.GetSrAlgorithmSid(node, algorithm)
		if sid == "" {
			service.log.Errorln("SID not found for router: ", node)
			continue
		}
		sidList = append(sidList, sid)
		if serviceSid, ok := routerServiceMap[node]; ok {
			sidList = append(sidList, serviceSid)
			serviceSidList = append(serviceSidList, serviceSid)
		}
	}
	service.log.Debugln("Translated SID List: ", sidList)
	return sidList, serviceSidList
}

func (service *CalculationTransformerService) TransformResult(path graph.Path, pathRequest domain.PathRequest, algorithm uint32) domain.PathResult {
	var sidList []string
	var serviceSidList []string
	if path == nil {
		service.log.Errorln("No path found, return destination IPv6 address as SID list")
		sidList = []string{pathRequest.GetIpv6DestinationAddress()}
	} else {
		sidList, serviceSidList = service.translatePathToSidList(path, algorithm)
	}
	pathResult, err := domain.NewDomainPathResult(pathRequest, path, sidList)
	if err != nil {
		service.log.Errorln("Error creating path result: ", err)
		return nil
	}
	pathResult.SetServiceSidList(serviceSidList)
	return pathResult
}
