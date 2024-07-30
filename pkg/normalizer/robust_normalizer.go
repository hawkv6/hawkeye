package normalizer

import (
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/montanaflynn/stats"
)

type RobustNormalizer struct {
	*BaseNormalizer
	interQuartileRangeLatency    float64
	interQuartileRangeJitter     float64
	interQuartileRangePacketLoss float64
	medianLatency                float64
	medianJitter                 float64
	medianPacketLoss             float64
}

func NewRobustNormalizer() *RobustNormalizer {
	return &RobustNormalizer{
		BaseNormalizer: NewBaseNormalizer(),
	}
}

func (normalizer *RobustNormalizer) setNormalizationIndicators(metrics []float64, median, interQuartileRange *float64) {
	data := stats.LoadRawData(metrics)
	if quartiles, err := stats.Quartile(data); err != nil {
		normalizer.log.Fatalln("Error calculating quartiles", err)
	} else {
		*median = quartiles.Q2
		*interQuartileRange = quartiles.Q3 - quartiles.Q1
		normalizer.log.Debugln("Q1: ", quartiles.Q1)
		normalizer.log.Debugln("Q2 / Median: ", quartiles.Q2)
		normalizer.log.Debugln("Q3: ", quartiles.Q3)
		normalizer.log.Debugln("InterQuartileRange (IQR): ", interQuartileRange)
	}
}

func (normalizer *RobustNormalizer) normalizeAndSetValue(value, median, interQuartileRange float64, setNormalizedValue func(float64)) float64 {
	normalizedValue := (value - median) / interQuartileRange
	setNormalizedValue(normalizedValue)
	return normalizedValue
}

func (normalizer *RobustNormalizer) normalizeLinks(links []domain.Link) {
	for i := 0; i < len(links); i++ {
		link := links[i]
		delay := float64(link.GetUnidirLinkDelay())
		jitter := float64(link.GetUnidirDelayVariation())
		packetLoss := float64(link.GetUnidirPacketLoss())
		normalizer.normalizedLatencyValues = append(normalizer.normalizedLatencyValues, normalizer.normalizeAndSetValue(delay, normalizer.medianLatency, normalizer.interQuartileRangeLatency, link.SetNormalizedUnidirLinkDelay))
		normalizer.normalizedJitterValues = append(normalizer.normalizedJitterValues, normalizer.normalizeAndSetValue(jitter, normalizer.medianJitter, normalizer.interQuartileRangeJitter, link.SetNormalizedUnidirDelayVariation))
		normalizer.normalizedPacketLoss = append(normalizer.normalizedPacketLoss, normalizer.normalizeAndSetValue(packetLoss, normalizer.medianPacketLoss, normalizer.interQuartileRangePacketLoss, link.SetNormalizedPacketLoss))

	}
}

func (normalizer *RobustNormalizer) Normalize(links []domain.Link) {
	normalizer.initializeNormalization(links)
	normalizer.setNormalizationIndicators(normalizer.currentLatencyValues, &normalizer.medianLatency, &normalizer.interQuartileRangeLatency)
	normalizer.setNormalizationIndicators(normalizer.currentJitterValues, &normalizer.medianJitter, &normalizer.interQuartileRangeJitter)
	normalizer.setNormalizationIndicators(normalizer.currentPacketLossValues, &normalizer.medianPacketLoss, &normalizer.interQuartileRangePacketLoss)
	normalizer.normalizeLinks(links)

}
