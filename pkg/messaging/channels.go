package messaging

import (
	"github.com/hawkv6/hawkeye/pkg/domain"
)

type MessagingChannels interface {
	GetPathRequestChan() chan domain.PathRequest
	GetPathResponseChan() chan domain.PathResult
}
