package domain

import "context"

type StreamSession interface {
	GetContext() context.Context
	GetPathRequest() PathRequest
	GetPathResult() PathResult
	SetPathResult(PathResult)
}
