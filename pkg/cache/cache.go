package cache

import "github.com/hawkv6/hawkeye/pkg/domain"

const Subsystem = "cache"

type Cache interface {
	Lock()
	Unlock()
	StoreClientNetwork(domain.Prefix)
	RemoveClientNetwork(domain.Prefix)
	GetClientNetworkByKey(string) domain.Prefix
	StoreSid(domain.Sid)
	RemoveSid(domain.Sid)
	GetSidByKey(string) domain.Sid
	GetRouterIdFromNetworkAddress(string) string
	GetSrAlgorithmSid(string, uint32) string
	StoreNode(node domain.Node)
	RemoveNode(node domain.Node)
	GetNodeByKey(string) domain.Node
	GetNodeByIgpRouterId(string) domain.Node
	StoreServiceSid(string, string)
	RemoveServiceSid(string, string)
	GetServiceSids(string) []string
}
