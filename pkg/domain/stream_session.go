package domain

import "context"

type StreamSession interface {
	GetContext() context.Context
	PathRequest
	PathResult
}
