package normalizer

import "github.com/hawkv6/hawkeye/pkg/domain"

const subsystem = "normalizer"

type Normalizer interface {
	Normalize([]domain.Link)
	GetCurrentLatencyValues() []float64
	GetCurrentJitterValues() []float64
	GetCurrentPacketLossValues() []float64
	GetNormalizedLatencyValues() []float64
	GetNormalizedJitterValues() []float64
	GetNormalizedPacketLossValues() []float64
}
