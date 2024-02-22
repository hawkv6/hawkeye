package messaging

import "github.com/hawkv6/hawkeye/pkg/domain"

type DefaultMessagingChannels struct {
	pathRequestChan chan domain.PathRequest
	pathResultChan  chan domain.PathResult
}

func NewDefaultMessagingChannels() *DefaultMessagingChannels {
	return &DefaultMessagingChannels{
		pathRequestChan: make(chan domain.PathRequest),
		pathResultChan:  make(chan domain.PathResult),
	}
}

func (channels *DefaultMessagingChannels) GetPathRequestChan() chan domain.PathRequest {
	return channels.pathRequestChan
}

func (channels *DefaultMessagingChannels) GetPathResponseChan() chan domain.PathResult {
	return channels.pathResultChan
}
