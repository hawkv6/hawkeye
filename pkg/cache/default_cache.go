package cache

import (
	"sync"

	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type DefaultCacheService struct {
	log             *logrus.Entry
	prefixMap       map[string]domain.Prefix
	prefixRouterMap map[string]string
	sidStore        map[string]domain.Sid
	routerSidMap    map[string]string
	nodeMap         map[string]domain.Node
	igpRouterIdMap  map[string]string
	linkMap         map[string]domain.Link
	mu              sync.Mutex
}

func NewDefaultCacheService() *DefaultCacheService {
	return &DefaultCacheService{
		log:             logging.DefaultLogger.WithField("subsystem", "cache"),
		prefixMap:       make(map[string]domain.Prefix),
		prefixRouterMap: make(map[string]string),
		sidStore:        make(map[string]domain.Sid),
		routerSidMap:    make(map[string]string),
		nodeMap:         make(map[string]domain.Node),
		igpRouterIdMap:  make(map[string]string),
		linkMap:         make(map[string]domain.Link),
		mu:              sync.Mutex{},
	}
}

func (cacheService *DefaultCacheService) Lock() {
	cacheService.mu.Lock()
}

func (cacheService *DefaultCacheService) Unlock() {
	cacheService.mu.Unlock()
}

func (cacheService *DefaultCacheService) StoreClientNetwork(prefix domain.Prefix) {
	cacheService.prefixMap[prefix.GetKey()] = prefix
	networkAddress := prefix.GetPrefix()
	cacheService.prefixRouterMap[networkAddress] = prefix.GetIgpRouterId()
}

func (cacheService *DefaultCacheService) RemoveClientNetwork(prefix domain.Prefix) {
	networkAddress := prefix.GetPrefix()
	delete(cacheService.prefixMap, prefix.GetKey())
	delete(cacheService.prefixRouterMap, networkAddress)
}

func (cacheService *DefaultCacheService) GetClientNetworkByKey(key string) (domain.Prefix, bool) {
	prefix, ok := cacheService.prefixMap[key]
	return prefix, ok
}

func (cacheService *DefaultCacheService) StoreSid(sid domain.Sid) {
	cacheService.sidStore[sid.GetKey()] = sid
	cacheService.routerSidMap[sid.GetIgpRouterId()] = sid.GetKey()
}

func (cacheService *DefaultCacheService) RemoveSid(sid domain.Sid) {
	delete(cacheService.sidStore, sid.GetKey())
	delete(cacheService.routerSidMap, sid.GetIgpRouterId())
}

func (cacheService *DefaultCacheService) GetSidByKey(key string) (domain.Sid, bool) {
	sid, ok := cacheService.sidStore[key]
	return sid, ok
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

func (cacheService *DefaultCacheService) StoreNode(node domain.Node) {
	cacheService.nodeMap[node.GetKey()] = node
	cacheService.igpRouterIdMap[node.GetIgpRouterId()] = node.GetKey()
}

func (cacheService *DefaultCacheService) GetNodeByKey(key string) (domain.Node, bool) {
	node, ok := cacheService.nodeMap[key]
	return node, ok
}

func (cacheService *DefaultCacheService) GetNodeByIgpRouterId(igpRouterId string) (domain.Node, bool) {
	key, ok := cacheService.igpRouterIdMap[igpRouterId]
	if !ok {
		return nil, ok
	} else {
		node, ok := cacheService.nodeMap[key]
		return node, ok
	}
}

func (cacheService *DefaultCacheService) StoreLink(link domain.Link) {
	cacheService.linkMap[link.GetKey()] = link
}

func (cacheService *DefaultCacheService) GetLinkByKey(key string) (domain.Link, bool) {
	link, ok := cacheService.linkMap[key]
	return link, ok
}
