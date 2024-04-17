package cache

import (
	"github.com/hawkv6/hawkeye/pkg/domain"
)

const Subsystem = "cache"

type Cache interface {
	Lock()
	Unlock()
	StoreClientNetwork(prefix domain.Prefix)
	RemoveClientNetwork(prefix domain.Prefix)
	GetClientNetworkByKey(string) (domain.Prefix, bool)
	StoreSid(sid domain.Sid)
	RemoveSid(sid domain.Sid)
	GetSidByKey(string) (domain.Sid, bool)
	GetRouterIdFromNetworkAddress(string) (string, bool)
	GetSidFromRouterId(string) (string, bool)
	StoreNode(node domain.Node)
	GetNodeByKey(string) (domain.Node, bool)
	GetNodeByIgpRouterId(string) (domain.Node, bool)
	StoreLink(link domain.Link)
	GetLinkByKey(string) (domain.Link, bool)
}
