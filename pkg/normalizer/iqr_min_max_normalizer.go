package normalizer

import (
	"math"

	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/montanaflynn/stats"
)

type IQRMinMaxNormalizer struct {
	*MinMaxNormalizer
}

func NewIQRMinMaxNormalizer() *IQRMinMaxNormalizer {
	return &IQRMinMaxNormalizer{
		MinMaxNormalizer: NewMinMaxNormalizer(),
	}
}

func (normalizer *IQRMinMaxNormalizer) calculateNormalizationIndicators(data stats.Float64Data, lowerFence *float64, upperFence *float64) {
	min, max := 0.0, 0.0
	normalizer.calculateMinMax(data, &min, &max)
	q1, q3, interQuartileRange := normalizer.calculateQuartiles(data)
	*upperFence = math.Min(q3+1.5*interQuartileRange, max)
	normalizer.log.Debugln("Upper fence (used min): ", *upperFence)
	*lowerFence = math.Max(q1-1.5*interQuartileRange, min)
	normalizer.log.Debugln("Lower fence (used max): ", *lowerFence)
}

func (normalizer *IQRMinMaxNormalizer) Normalize(links []domain.Link) {
	normalizer.initializeNormalization(links)
	normalizer.calculateNormalizationIndicators(stats.LoadRawData(normalizer.currentLatencyValues), &normalizer.minLatency, &normalizer.maxLatency)
	normalizer.calculateNormalizationIndicators(stats.LoadRawData(normalizer.currentJitterValues), &normalizer.minJitter, &normalizer.maxJitter)
	normalizer.calculateNormalizationIndicators(stats.LoadRawData(normalizer.currentPacketLossValues), &normalizer.minPacketLoss, &normalizer.maxPacketLoss)
	normalizer.normalizeLinks(links)
}
