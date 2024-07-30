package calculation

import (
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
)

type CalculationTransformer interface {
	TransformResult(path graph.Path, pathRequest domain.PathRequest, algorithm uint32) domain.PathResult
}
