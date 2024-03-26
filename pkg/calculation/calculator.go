package calculation

import "github.com/hawkv6/hawkeye/pkg/domain"

const subsystem = "calculation"

type Calculator interface {
	HandlePathRequest(domain.PathRequest) domain.PathResult
	UpdatePathSession(domain.StreamSession) *domain.PathResult
}
