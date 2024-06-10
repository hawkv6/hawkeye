package normalizer

import (
	"math"

	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/montanaflynn/stats"
)

type IQRMinMaxNormalizer struct {
	*MinMaxNormalizer
	latencyQueue    Queue
	jitterQueue     Queue
	packetLossQueue Queue
}

func NewIQRMinMaxNormalizer(latencyQueue, jitterQueue, packetLossQueue Queue) *IQRMinMaxNormalizer {
	return &IQRMinMaxNormalizer{
		MinMaxNormalizer: NewMinMaxNormalizer(),
		latencyQueue:     latencyQueue,
		jitterQueue:      jitterQueue,
		packetLossQueue:  packetLossQueue,
	}
}

func (normalizer *IQRMinMaxNormalizer) calculateNormalizationIndicators(data stats.Float64Data, queue Queue, lowerFence *float64, upperFence *float64) {
	min, max := 0.0, 0.0
	normalizer.calculateMinMax(data, &min, &max)
	q1, q3, _ := normalizer.calculateQuartiles(data)
	outliers, err := stats.QuartileOutliers(data)
	if err != nil {
		normalizer.log.Fatalf("Error calculating outliers %s", err)
	}
	normalizer.log.Debugf("Outliers: %+v ", outliers)

	queue.EnqueueMax(max)
	queue.EnqueueMin(min)
	queue.EnqueueQ1(q1)
	queue.EnqueueQ3(q3)

	normalizer.log.Debugln("Historical Q1 Values: ", queue.GetQ1Elements())
	q1Average := queue.GetAverageQ1()
	normalizer.log.Debugln("Rolling (Average) Q1: ", q1Average)

	normalizer.log.Debugln("Historical Q3 Values: ", queue.GetQ3Elements())
	q3Average := queue.GetAverageQ3()
	normalizer.log.Debugln("Rolling (Average) Q3: ", q3Average)
	iqrAverage := q3Average - q1Average
	normalizer.log.Debugln("Rolling IQR: ", iqrAverage)

	*upperFence = math.Min(q1Average+1.5*iqrAverage, queue.GetAverageMax())
	normalizer.log.Debugln("Upper fence calculated - rolling Q1 + 1.5 * rolling IQR: ", *upperFence)
	*lowerFence = math.Max(q1Average-1.5*iqrAverage, queue.GetAverageMin())
	normalizer.log.Debugln("Lower fence calculated - rolling Q1 - 1.5 * rolling IQR: ", *lowerFence)
}

func (normalizer *IQRMinMaxNormalizer) Normalize(links []domain.Link) {
	normalizer.initializeNormalization(links)
	normalizer.log.Debugln("Calculate normalization indicators for latency metric")
	normalizer.calculateNormalizationIndicators(stats.LoadRawData(normalizer.currentLatencyValues), normalizer.latencyQueue, &normalizer.minLatency, &normalizer.maxLatency)
	normalizer.log.Debugln("Calculate normalization indicators for jitter metric")
	normalizer.calculateNormalizationIndicators(stats.LoadRawData(normalizer.currentJitterValues), normalizer.jitterQueue, &normalizer.minJitter, &normalizer.maxJitter)
	normalizer.log.Debugln("Calculate normalization indicators for packet loss metric")
	normalizer.calculateNormalizationIndicators(stats.LoadRawData(normalizer.currentPacketLossValues), normalizer.packetLossQueue, &normalizer.minPacketLoss, &normalizer.maxPacketLoss)
	normalizer.normalizeLinks(links)
}
