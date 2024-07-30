package domain

import (
	"context"
)

type StreamSession interface {
	GetContext() context.Context
	GetPathRequest() PathRequest
	GetPathResult() PathResult
	SetPathResult(PathResult)
}

type DomainStreamSession struct {
	pathRequest PathRequest
	pathResult  PathResult
}

func NewDomainStreamSession(pathRequest PathRequest, pathResponse PathResult) *DomainStreamSession {
	return &DomainStreamSession{
		pathRequest: pathRequest,
		pathResult:  pathResponse,
	}
}

func (streamSession *DomainStreamSession) GetContext() context.Context {
	return streamSession.pathRequest.GetContext()
}

func (streamSession *DomainStreamSession) GetPathRequest() PathRequest {
	return streamSession.pathRequest
}

func (streamSession *DomainStreamSession) GetPathResult() PathResult {
	return streamSession.pathResult
}

func (streamSession *DomainStreamSession) SetPathResult(pathResult PathResult) {
	streamSession.pathResult = pathResult
}
