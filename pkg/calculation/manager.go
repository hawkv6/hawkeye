package calculation

import "github.com/hawkv6/hawkeye/pkg/domain"

const subsystem = "calculation"

type Manager interface {
	CalculateBestPath(domain.PathRequest) domain.PathResult
	CalculatePathUpdate(domain.StreamSession) *domain.PathResult
}
