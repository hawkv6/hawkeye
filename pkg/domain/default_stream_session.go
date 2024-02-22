package domain

import (
	"github.com/go-playground/validator"
)

type DefaultStreamSession struct {
	PathRequest `validate:"required"`
	PathResult  `validate:"required"`
}

func NewDefaultStreamSession(pathRequest PathRequest, pathResponse PathResult) (*DefaultStreamSession, error) {
	defaultStreamSession := &DefaultStreamSession{
		PathRequest: pathRequest,
		PathResult:  pathResponse,
	}
	validator := validator.New()
	err := validator.Struct(defaultStreamSession)
	if err != nil {
		return nil, err
	}
	return defaultStreamSession, nil
}
