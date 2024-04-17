package messaging

import "github.com/hawkv6/hawkeye/pkg/domain"

type PathMessagingChannels struct {
	pathRequestChan chan domain.PathRequest
	pathResultChan  chan domain.PathResult
}

func NewPathMessagingChannels() *PathMessagingChannels {
	return &PathMessagingChannels{
		pathRequestChan: make(chan domain.PathRequest),
		pathResultChan:  make(chan domain.PathResult),
	}
}

func (channels *PathMessagingChannels) GetPathRequestChan() chan domain.PathRequest {
	return channels.pathRequestChan
}

func (channels *PathMessagingChannels) GetPathResponseChan() chan domain.PathResult {
	return channels.pathResultChan
}
