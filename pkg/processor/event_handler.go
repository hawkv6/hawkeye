package processor

import "github.com/hawkv6/hawkeye/pkg/domain"

type EventHandler interface {
	HandleEvent(event domain.NetworkEvent) bool
}
