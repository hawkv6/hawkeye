package cache

import (
	"sync"

	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type InMemoryCache struct {
	log                         *logrus.Entry
	prefixStore                 map[string]domain.Prefix
	prefixToRouterIdMap         map[string]string
	sidStore                    map[string]domain.Sid
	igpRouterIdToSrAlgoToSidMap map[string]map[uint32]string
	nodeStore                   map[string]domain.Node
	igpRouterIdToRouterKeyMap   map[string]string
	serviceSidStore             map[string]map[string]bool
	mu                          sync.Mutex
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		log:                         logging.DefaultLogger.WithField("subsystem", "cache"),
		prefixStore:                 make(map[string]domain.Prefix),
		prefixToRouterIdMap:         make(map[string]string),
		sidStore:                    make(map[string]domain.Sid),
		igpRouterIdToSrAlgoToSidMap: make(map[string]map[uint32]string),
		nodeStore:                   make(map[string]domain.Node),
		igpRouterIdToRouterKeyMap:   make(map[string]string),
		serviceSidStore:             make(map[string]map[string]bool),
		mu:                          sync.Mutex{},
	}
}

func (cache *InMemoryCache) Lock() {
	cache.mu.Lock()
}

func (cache *InMemoryCache) Unlock() {
	cache.mu.Unlock()
}

func (cache *InMemoryCache) StoreClientNetwork(prefix domain.Prefix) {
	cache.prefixStore[prefix.GetKey()] = prefix
	networkAddress := prefix.GetPrefix()
	cache.prefixToRouterIdMap[networkAddress] = prefix.GetIgpRouterId()
}

func (cache *InMemoryCache) RemoveClientNetwork(prefix domain.Prefix) {
	networkAddress := prefix.GetPrefix()
	delete(cache.prefixStore, prefix.GetKey())
	delete(cache.prefixToRouterIdMap, networkAddress)
}

func (cache *InMemoryCache) GetClientNetworkByKey(key string) domain.Prefix {
	return cache.prefixStore[key]
}

func (cache *InMemoryCache) StoreSid(sid domain.Sid) {
	cache.sidStore[sid.GetKey()] = sid
	igpRouterId := sid.GetIgpRouterId()
	if _, ok := cache.igpRouterIdToSrAlgoToSidMap[sid.GetIgpRouterId()]; !ok {
		cache.igpRouterIdToSrAlgoToSidMap[igpRouterId] = make(map[uint32]string)
	}
	cache.igpRouterIdToSrAlgoToSidMap[igpRouterId][sid.GetAlgorithm()] = sid.GetKey()
}

func (cache *InMemoryCache) RemoveSid(sid domain.Sid) {
	delete(cache.sidStore, sid.GetKey())
	delete(cache.igpRouterIdToSrAlgoToSidMap, sid.GetIgpRouterId())
}

func (cache *InMemoryCache) GetSidByKey(key string) domain.Sid {
	return cache.sidStore[key]
}

func (cache *InMemoryCache) GetRouterIdFromNetworkAddress(networkAddress string) string {
	return cache.prefixToRouterIdMap[networkAddress]
}

func (cache *InMemoryCache) GetSrAlgorithmSid(igpRouterId string, srAlgorithm uint32) string {
	if _, ok := cache.igpRouterIdToSrAlgoToSidMap[igpRouterId]; !ok {
		return ""
	}
	if sidKey, ok := cache.igpRouterIdToSrAlgoToSidMap[igpRouterId][srAlgorithm]; !ok {
		return ""
	} else {
		return cache.sidStore[sidKey].GetSid()
	}
}

func (cache *InMemoryCache) StoreNode(node domain.Node) {
	cache.nodeStore[node.GetKey()] = node
	cache.igpRouterIdToRouterKeyMap[node.GetIgpRouterId()] = node.GetKey()
}

func (cache *InMemoryCache) RemoveNode(node domain.Node) {
	delete(cache.nodeStore, node.GetKey())
	delete(cache.igpRouterIdToRouterKeyMap, node.GetIgpRouterId())
}

func (cache *InMemoryCache) GetNodeByKey(key string) domain.Node {
	return cache.nodeStore[key]
}

func (cache *InMemoryCache) GetNodeByIgpRouterId(igpRouterId string) domain.Node {
	key, ok := cache.igpRouterIdToRouterKeyMap[igpRouterId]
	if !ok {
		return nil
	} else {
		return cache.nodeStore[key]
	}
}

func (cache *InMemoryCache) StoreServiceSid(serviceType, servicePrefixSid string) {
	if _, ok := cache.serviceSidStore[serviceType]; !ok {
		cache.serviceSidStore[serviceType] = make(map[string]bool)
	}
	cache.serviceSidStore[serviceType][servicePrefixSid] = true
}

func (cache *InMemoryCache) RemoveServiceSid(serviceType, servicePrefixSid string) {
	delete(cache.serviceSidStore[serviceType], servicePrefixSid)
	if len(cache.serviceSidStore[serviceType]) == 0 {
		delete(cache.serviceSidStore, serviceType)
	}
}

func (cache *InMemoryCache) GetServiceSids(serviceType string) []string {
	sids := make([]string, 0, len(cache.serviceSidStore[serviceType]))
	for sid := range cache.serviceSidStore[serviceType] {
		sids = append(sids, sid)
	}
	return sids
}

func (cache *InMemoryCache) DoesServiceSidExist(servicePrefixSid string) bool {
	for _, sids := range cache.serviceSidStore {
		if sids[servicePrefixSid] {
			return true
		}
	}
	return false
}
