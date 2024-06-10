package normalizer

import (
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/montanaflynn/stats"
)

type ZScoreNormalizer struct {
	*BaseNormalizer
	meanLatency                 float64
	meanJitter                  float64
	meanPacketLoss              float64
	standardDeviationLatency    float64
	standardDeviationJitter     float64
	standardDeviationPacketLoss float64
}

func NewZScoreNormalizer() *ZScoreNormalizer {
	return &ZScoreNormalizer{
		BaseNormalizer: NewBaseNormalizer(),
	}
}

func (normalizer *ZScoreNormalizer) normalizeAndSetValue(value, mean, standardDeviation float64, setNormalizedValue func(float64)) float64 {
	zscore := (value - mean) / standardDeviation
	setNormalizedValue(zscore)
	return zscore
}

func (normalizer *ZScoreNormalizer) normalizeLinks(links []domain.Link) {
	for i := 0; i < len(links); i++ {
		link := links[i]
		delay := float64(link.GetUnidirLinkDelay())
		jitter := float64(link.GetUnidirDelayVariation())
		packetLoss := float64(link.GetUnidirPacketLoss())
		normalizer.normalizedLatencyValues = append(normalizer.normalizedLatencyValues, normalizer.normalizeAndSetValue(delay, normalizer.meanLatency, normalizer.standardDeviationLatency, link.SetNormalizedUnidirLinkDelay))
		normalizer.normalizedJitterValues = append(normalizer.normalizedJitterValues, normalizer.normalizeAndSetValue(jitter, normalizer.meanJitter, normalizer.standardDeviationJitter, link.SetNormalizedUnidirDelayVariation))
		normalizer.normalizedPacketLoss = append(normalizer.normalizedPacketLoss, normalizer.normalizeAndSetValue(packetLoss, normalizer.meanPacketLoss, normalizer.standardDeviationPacketLoss, link.SetNormalizedPacketLoss))
	}
}

func (normalizer *ZScoreNormalizer) setNormalizationIndicators(metrics []float64, mean, standardDeviation *float64) {
	data := stats.LoadRawData(metrics)

	if currentMean, err := stats.Mean(data); err != nil {
		normalizer.log.Fatalln("Error calculating mean", err)
	} else {
		*mean = currentMean
		normalizer.log.Debugln("Mean: ", *mean)
	}

	if currentStandardDeviation, err := stats.StandardDeviation(data); err != nil {
		normalizer.log.Fatalln("Error calculating standard deviation", err)
	} else {
		*standardDeviation = currentStandardDeviation
		normalizer.log.Debugln("Standard Deviation: ", *standardDeviation)
	}
}

func (normalizer *ZScoreNormalizer) Normalize(links []domain.Link) {
	normalizer.initializeNormalization(links)
	normalizer.setNormalizationIndicators(normalizer.currentLatencyValues, &normalizer.meanLatency, &normalizer.standardDeviationLatency)
	normalizer.setNormalizationIndicators(normalizer.currentJitterValues, &normalizer.meanJitter, &normalizer.standardDeviationJitter)
	normalizer.setNormalizationIndicators(normalizer.currentPacketLossValues, &normalizer.meanPacketLoss, &normalizer.standardDeviationPacketLoss)
	normalizer.normalizeLinks(links)
}
