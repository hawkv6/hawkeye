package normalizer

import (
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/montanaflynn/stats"
)

type MinMaxNormalizer struct {
	*BaseNormalizer
	minLatency    float64
	minJitter     float64
	minPacketLoss float64
	maxLatency    float64
	maxJitter     float64
	maxPacketLoss float64
}

func NewMinMaxNormalizer() *MinMaxNormalizer {
	return &MinMaxNormalizer{
		BaseNormalizer: NewBaseNormalizer(),
	}
}

func (normalizer *MinMaxNormalizer) normalizeAndSetValue(value float64, min float64, max float64, setNormalizedValue func(float64)) float64 {
	normalizedValue := (value - min) / (max - min)
	if normalizedValue > 1 {
		normalizedValue = 1
	} else if normalizedValue < 0 {
		normalizedValue = 0
	}
	setNormalizedValue(normalizedValue)
	return normalizedValue
}

func (normalizer *MinMaxNormalizer) normalizeLinks(links []domain.Link) {
	normalizer.log.Debugln("Normalize links using min-max normalization")
	normalizer.log.Debugln("Min latency: ", normalizer.minLatency, " Max latency: ", normalizer.maxLatency)
	normalizer.log.Debugln("Min jitter: ", normalizer.minJitter, " Max jitter: ", normalizer.maxJitter)
	normalizer.log.Debugln("Min packet loss: ", normalizer.minPacketLoss, " Max packet loss: ", normalizer.maxPacketLoss)
	for i := 0; i < len(links); i++ {
		link := links[i]
		delay := float64(link.GetUnidirLinkDelay())
		jitter := float64(link.GetUnidirDelayVariation())
		packetLoss := float64(link.GetUnidirPacketLoss())
		normalizer.normalizedLatencyValues = append(normalizer.normalizedLatencyValues, normalizer.normalizeAndSetValue(delay, normalizer.minLatency, normalizer.maxLatency, link.SetNormalizedUnidirLinkDelay))
		normalizer.normalizedJitterValues = append(normalizer.normalizedJitterValues, normalizer.normalizeAndSetValue(jitter, normalizer.minJitter, normalizer.maxJitter, link.SetNormalizedUnidirDelayVariation))
		normalizer.normalizedPacketLoss = append(normalizer.normalizedPacketLoss, normalizer.normalizeAndSetValue(packetLoss, normalizer.minPacketLoss, normalizer.maxPacketLoss, link.SetNormalizedPacketLoss))

	}
}

func (normalizer *MinMaxNormalizer) calculateQuartiles(data stats.Float64Data) (float64, float64, float64) {
	if quartiles, err := stats.Quartile(data); err != nil {
		normalizer.log.Fatalf("Error calculating quartiles %s", err)
		return 0.0, 0.0, 0.0
	} else {
		interQuartileRange := quartiles.Q3 - quartiles.Q1
		normalizer.log.Debugln("Q1: ", quartiles.Q1)
		normalizer.log.Debugln("Q2 / Median: ", quartiles.Q2)
		normalizer.log.Debugln("Q3: ", quartiles.Q3)
		normalizer.log.Debugln("InterQuartileRange (IQR): ", interQuartileRange)
		return quartiles.Q1, quartiles.Q3, interQuartileRange
	}
}

func (normalizer *MinMaxNormalizer) calculateNormalizationIndicators(data stats.Float64Data, min *float64, max *float64) {
	normalizer.calculateMinMax(data, min, max)
	normalizer.calculateQuartiles(data)
}

func (normalizer *MinMaxNormalizer) calculateMinMax(data stats.Float64Data, min *float64, max *float64) {
	var err error
	*min, err = stats.Min(data)
	if err != nil {
		normalizer.log.Fatalf("Error calculating min: %s", err)
	}
	normalizer.log.Debugln("Min: ", *min)
	*max, err = stats.Max(data)
	if err != nil {
		normalizer.log.Fatalf("Error calculating max: %s", err)
	}
	normalizer.log.Debugln("Max: ", *max)
}

func (normalizer *MinMaxNormalizer) Normalize(links []domain.Link) {
	normalizer.initializeNormalization(links)
	normalizer.calculateNormalizationIndicators(stats.LoadRawData(normalizer.currentLatencyValues), &normalizer.minLatency, &normalizer.maxLatency)
	normalizer.calculateNormalizationIndicators(stats.LoadRawData(normalizer.currentJitterValues), &normalizer.minJitter, &normalizer.maxJitter)
	normalizer.calculateNormalizationIndicators(stats.LoadRawData(normalizer.currentPacketLossValues), &normalizer.minPacketLoss, &normalizer.maxPacketLoss)
	normalizer.normalizeLinks(links)
}
