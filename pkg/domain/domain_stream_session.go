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
	pathRequest PathRequest
	pathResult  PathResult
}

type DomainStreamSessionInput struct {
	PathRequest  PathRequest `validate:"required"`
	PathResponse PathResult  `validate:"required"`
}

func NewDomainStreamSession(pathRequest PathRequest, pathResponse PathResult) (*DomainStreamSession, error) {
	domainStreamSessionInput := &DomainStreamSessionInput{
		PathRequest:  pathRequest,
		PathResponse: pathResponse,
	}
	validator := validator.New()
	err := validator.Struct(domainStreamSessionInput)
	if err != nil {
		return nil, err
	}
	defaultStreamSession := &DomainStreamSession{
		pathRequest: pathRequest,
		pathResult:  pathResponse,
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
