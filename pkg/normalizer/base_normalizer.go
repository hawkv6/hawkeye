package normalizer

import (
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type BaseNormalizer struct {
	log                     *logrus.Entry
	currentLatencyValues    []float64
	currentJitterValues     []float64
	currentPacketLossValues []float64
	normalizedLatencyValues []float64
	normalizedJitterValues  []float64
	normalizedPacketLoss    []float64
}

func NewBaseNormalizer() *BaseNormalizer {
	return &BaseNormalizer{
		log: logging.DefaultLogger.WithField("subsystem", subsystem),
	}
}

func (normalizer *BaseNormalizer) extractCurrentLinkMetrics(links []domain.Link) {
	currentJitterValues := make([]float64, 0, len(links))
	currentLatencyValues := make([]float64, 0, len(links))
	currentPacketLossValues := make([]float64, 0, len(links))
	for i := 0; i < len(links); i++ {
		link := links[i]
		currentLatencyValues = append(currentLatencyValues, float64(link.GetUnidirLinkDelay()))
		currentJitterValues = append(currentJitterValues, float64(link.GetUnidirDelayVariation()))
		currentPacketLossValues = append(currentPacketLossValues, float64(link.GetUnidirPacketLoss()))
	}
	normalizer.currentLatencyValues = currentLatencyValues
	normalizer.currentJitterValues = currentJitterValues
	normalizer.currentPacketLossValues = currentPacketLossValues
}

func (normalizer *BaseNormalizer) initializeNormalization(links []domain.Link) {
	normalizer.extractCurrentLinkMetrics(links)
	normalizer.normalizedLatencyValues = make([]float64, 0, len(normalizer.currentLatencyValues))
	normalizer.normalizedJitterValues = make([]float64, 0, len(normalizer.currentJitterValues))
	normalizer.normalizedPacketLoss = make([]float64, 0, len(normalizer.currentPacketLossValues))
}

func (normalizer *BaseNormalizer) GetCurrentLatencyValues() []float64 {
	return normalizer.currentLatencyValues
}

func (normalizer *BaseNormalizer) GetCurrentJitterValues() []float64 {
	return normalizer.currentJitterValues
}

func (normalizer *BaseNormalizer) GetCurrentPacketLossValues() []float64 {
	return normalizer.currentPacketLossValues
}

func (normalizer *BaseNormalizer) GetNormalizedLatencyValues() []float64 {
	return normalizer.normalizedLatencyValues
}

func (normalizer *BaseNormalizer) GetNormalizedJitterValues() []float64 {
	return normalizer.normalizedJitterValues
}

func (normalizer *BaseNormalizer) GetNormalizedPacketLossValues() []float64 {
	return normalizer.normalizedPacketLoss
}
