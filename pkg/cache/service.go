package cache

import (
	"github.com/hawkv6/hawkeye/pkg/domain"
)

const Subsystem = "cache"

type CacheService interface {
	StoreClientNetwork(prefix domain.Prefix)
	StoreSids(sid domain.Sid)
	GetRouterIdFromNetworkAddress(string) (string, bool)
	GetSidFromRouterId(string) (string, bool)
}
