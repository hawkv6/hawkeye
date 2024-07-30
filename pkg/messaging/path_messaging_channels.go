package messaging

import "github.com/hawkv6/hawkeye/pkg/domain"

type PathMessagingChannels struct {
	pathRequestChan chan domain.PathRequest
	pathResultChan  chan domain.PathResult
	errorChan       chan error
}

func NewPathMessagingChannels() *PathMessagingChannels {
	return &PathMessagingChannels{
		pathRequestChan: make(chan domain.PathRequest),
		pathResultChan:  make(chan domain.PathResult),
		errorChan:       make(chan error),
	}
}

func (channels *PathMessagingChannels) GetPathRequestChan() chan domain.PathRequest {
	return channels.pathRequestChan
}

func (channels *PathMessagingChannels) GetPathResponseChan() chan domain.PathResult {
	return channels.pathResultChan
}

func (channels *PathMessagingChannels) GetErrorChan() chan error {
	return channels.errorChan
}
