package cache

import (
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/sirupsen/logrus"
)

type DefaultCacheService struct {
	log             *logrus.Entry
	prefixMap       map[string]domain.Prefix
	prefixRouterMap map[string]string
	sidStore        map[string]domain.Sid
	routerSidMap    map[string]string
}

func NewDefaultCacheService() *DefaultCacheService {
	return &DefaultCacheService{
		log:             logrus.WithField("subsystem", Subsystem),
		prefixMap:       make(map[string]domain.Prefix),
		prefixRouterMap: make(map[string]string),
		sidStore:        make(map[string]domain.Sid),
		routerSidMap:    make(map[string]string),
	}
}

func (cacheService *DefaultCacheService) StoreClientNetwork(prefix domain.Prefix) {
	cacheService.prefixMap[prefix.GetKey()] = prefix
	networkAddress := prefix.GetPrefix()
	cacheService.prefixRouterMap[networkAddress] = prefix.GetIgpRouterId()
}

func (cacheService *DefaultCacheService) StoreSids(sid domain.Sid) {
	cacheService.sidStore[sid.GetKey()] = sid
	cacheService.routerSidMap[sid.GetIgpRouterId()] = sid.GetKey()
}

func (cacheService *DefaultCacheService) GetRouterIdFromNetworkAddress(networkAddress string) (string, bool) {
	routerId, ok := cacheService.prefixRouterMap[networkAddress]
	return routerId, ok
}

func (cacheService *DefaultCacheService) GetSidFromRouterId(routerId string) (string, bool) {
	sidKey, ok := cacheService.routerSidMap[routerId]
	if !ok {
		return "", ok
	} else {
		sid, ok := cacheService.sidStore[sidKey]
		return sid.GetSid(), ok
	}
}
