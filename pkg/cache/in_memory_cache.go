package cache

import (
	"sync"

	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type InMemoryCache struct {
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

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
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

func (cache *InMemoryCache) Lock() {
	cache.mu.Lock()
}

func (cache *InMemoryCache) Unlock() {
	cache.mu.Unlock()
}

func (cache *InMemoryCache) StoreClientNetwork(prefix domain.Prefix) {
	cache.prefixMap[prefix.GetKey()] = prefix
	networkAddress := prefix.GetPrefix()
	cache.prefixRouterMap[networkAddress] = prefix.GetIgpRouterId()
}

func (cache *InMemoryCache) RemoveClientNetwork(prefix domain.Prefix) {
	networkAddress := prefix.GetPrefix()
	delete(cache.prefixMap, prefix.GetKey())
	delete(cache.prefixRouterMap, networkAddress)
}

func (cache *InMemoryCache) GetClientNetworkByKey(key string) (domain.Prefix, bool) {
	prefix, ok := cache.prefixMap[key]
	return prefix, ok
}

func (cache *InMemoryCache) StoreSid(sid domain.Sid) {
	cache.sidStore[sid.GetKey()] = sid
	cache.routerSidMap[sid.GetIgpRouterId()] = sid.GetKey()
}

func (cache *InMemoryCache) RemoveSid(sid domain.Sid) {
	delete(cache.sidStore, sid.GetKey())
	delete(cache.routerSidMap, sid.GetIgpRouterId())
}

func (cache *InMemoryCache) GetSidByKey(key string) (domain.Sid, bool) {
	sid, ok := cache.sidStore[key]
	return sid, ok
}

func (cache *InMemoryCache) GetRouterIdFromNetworkAddress(networkAddress string) (string, bool) {
	routerId, ok := cache.prefixRouterMap[networkAddress]
	return routerId, ok
}

func (cache *InMemoryCache) GetSidFromRouterId(routerId string) (string, bool) {
	sidKey, ok := cache.routerSidMap[routerId]
	if !ok {
		return "", ok
	} else {
		sid, ok := cache.sidStore[sidKey]
		return sid.GetSid(), ok
	}
}

func (cache *InMemoryCache) StoreNode(node domain.Node) {
	cache.nodeMap[node.GetKey()] = node
	cache.igpRouterIdMap[node.GetIgpRouterId()] = node.GetKey()
}

func (cache *InMemoryCache) GetNodeByKey(key string) (domain.Node, bool) {
	node, ok := cache.nodeMap[key]
	return node, ok
}

func (cache *InMemoryCache) GetNodeByIgpRouterId(igpRouterId string) (domain.Node, bool) {
	key, ok := cache.igpRouterIdMap[igpRouterId]
	if !ok {
		return nil, ok
	} else {
		node, ok := cache.nodeMap[key]
		return node, ok
	}
}

func (cache *InMemoryCache) StoreLink(link domain.Link) {
	cache.linkMap[link.GetKey()] = link
}

func (cache *InMemoryCache) GetLinkByKey(key string) (domain.Link, bool) {
	link, ok := cache.linkMap[key]
	return link, ok
}
