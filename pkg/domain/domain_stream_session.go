package domain

import (
	"context"

	"github.com/go-playground/validator"
)

type StreamSession interface {
	GetContext() context.Context
	GetPathRequest() PathRequest
	GetPathResult() PathResult
	SetPathResult(PathResult)
}

type DomainStreamSession struct {
	pathRequest PathRequest `validate:"required"`
	pathResult  PathResult  `validate:"required"`
}

func NewDefaultStreamSession(pathRequest PathRequest, pathResponse PathResult) (*DomainStreamSession, error) {
	defaultStreamSession := &DomainStreamSession{
		pathRequest: pathRequest,
		pathResult:  pathResponse,
	}
	validator := validator.New()
	err := validator.Struct(defaultStreamSession)
	if err != nil {
		return nil, err
	}
	return defaultStreamSession, nil
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
