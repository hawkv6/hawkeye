package domain

import (
	"context"

	"github.com/go-playground/validator"
)

type DefaultStreamSession struct {
	pathRequest PathRequest `validate:"required"`
	pathResult  PathResult  `validate:"required"`
}

func NewDefaultStreamSession(pathRequest PathRequest, pathResponse PathResult) (*DefaultStreamSession, error) {
	defaultStreamSession := &DefaultStreamSession{
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

func (streamSession *DefaultStreamSession) GetContext() context.Context {
	return streamSession.pathRequest.GetContext()
}

func (streamSession *DefaultStreamSession) GetPathRequest() PathRequest {
	return streamSession.pathRequest
}

func (streamSession *DefaultStreamSession) GetPathResult() PathResult {
	return streamSession.pathResult
}

func (streamSession *DefaultStreamSession) SetPathResult(pathResult PathResult) {
	streamSession.pathResult = pathResult
}
