package cache

import (
	"github.com/hawkv6/hawkeye/pkg/domain"
)

const Subsystem = "cache"

type CacheService interface {
	Lock()
	Unlock()
	StoreClientNetwork(prefix domain.Prefix)
	StoreSids(sid domain.Sid)
	GetRouterIdFromNetworkAddress(string) (string, bool)
	GetSidFromRouterId(string) (string, bool)
	StoreNode(node domain.Node)
	GetNodeByKey(string) (domain.Node, bool)
	GetNodeByIgpRouterId(string) (domain.Node, bool)
	StoreLink(link domain.Link)
	GetLinkByKey(string) (domain.Link, bool)
}
