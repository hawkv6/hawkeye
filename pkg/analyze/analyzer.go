package analyze

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/hawkv6/hawkeye/pkg/normalizer"
	"github.com/montanaflynn/stats"
	"github.com/sirupsen/logrus"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Analyzer interface {
	Analyze()
}

type MetricAnalyzer struct {
	log                         *logrus.Entry
	normalizer                  normalizer.Normalizer
	currentLatencyMetrics       []float64
	currentJitterMetrics        []float64
	currentPacketLossMetrics    []float64
	normalizedLatencyMetrics    []float64
	normalizedJitterMetrics     []float64
	normalizedPacketLossMetrics []float64
	folderName                  string
	plotName                    string
}

func NewMetricAnalyzer(normalizer normalizer.Normalizer, folderName, plotName string) *MetricAnalyzer {
	return &MetricAnalyzer{
		log:                         logging.DefaultLogger.WithField("subsystem", "analyze"),
		normalizer:                  normalizer,
		folderName:                  folderName,
		plotName:                    plotName,
		currentLatencyMetrics:       normalizer.GetCurrentLatencyValues(),
		currentJitterMetrics:        normalizer.GetCurrentJitterValues(),
		currentPacketLossMetrics:    normalizer.GetCurrentPacketLossValues(),
		normalizedLatencyMetrics:    normalizer.GetNormalizedLatencyValues(),
		normalizedJitterMetrics:     normalizer.GetNormalizedJitterValues(),
		normalizedPacketLossMetrics: normalizer.GetNormalizedPacketLossValues(),
	}
}

func (analyzer *MetricAnalyzer) calculateMean(metricName string, data stats.Float64Data) float64 {
	mean, err := stats.Mean(data)
	if err != nil {
		analyzer.log.Fatalf("Error calculating mean %s", metricName)
	}
	analyzer.log.Debugf("Mean %s: %f", metricName, mean)
	return mean
}

func (analyzer *MetricAnalyzer) calculateStandardDeviation(metricName string, data stats.Float64Data) float64 {
	standardDeviation, err := stats.StandardDeviation(data)
	if err != nil {
		analyzer.log.Fatalf("Error calculating standard deviation %s", metricName)
	}
	analyzer.log.Debugf("Standard deviation %s: %f", metricName, standardDeviation)
	return standardDeviation
}

func (analyzer *MetricAnalyzer) calculateQuartiles(metricName string, data stats.Float64Data) (stats.Quartiles, float64) {
	quartiles, err := stats.Quartile(data)
	if err != nil {
		analyzer.log.Fatalf("Error calculating quartiles %s", metricName)
	}
	analyzer.log.Debugln("Q1: ", quartiles.Q1)
	analyzer.log.Debugln("Q2 / Median: ", quartiles.Q2)
	analyzer.log.Debugln("Q3", quartiles.Q3)
	interQuartileRange := quartiles.Q3 - quartiles.Q1
	analyzer.log.Debugln("Interquartile Range", interQuartileRange)
	return quartiles, interQuartileRange
}

func (analyzer *MetricAnalyzer) calculateMin(metricName string, data stats.Float64Data) float64 {
	min, err := stats.Min(data)
	if err != nil {
		analyzer.log.Fatalf("Error calculating min %s", metricName)
	}
	analyzer.log.Debugf("Min %s: %f", metricName, min)
	return min
}

func (analyzer *MetricAnalyzer) calculateMax(metricName string, data stats.Float64Data) float64 {
	max, err := stats.Max(data)
	if err != nil {
		analyzer.log.Fatalf("Error calculating max %s", metricName)
	}
	analyzer.log.Debugf("Max %s: %f", metricName, max)
	return max
}

func (analyzer *MetricAnalyzer) findOutliers(data stats.Float64Data) stats.Outliers {
	outliers, err := data.QuartileOutliers()
	if err != nil {
		fmt.Println("Error in outliers calculation: ", err)
	}
	return outliers
}

func (analyzer *MetricAnalyzer) calculateStatisticalIndicators(metricName string, metrics []float64) {
	data := stats.LoadRawData(metrics)
	mean := analyzer.calculateMean(metricName, data)
	standardDeviation := analyzer.calculateStandardDeviation(metricName, data)
	quartiles, interQuartileRange := analyzer.calculateQuartiles(metricName, data)
	median := quartiles.Q2
	min := analyzer.calculateMin(metricName, data)
	max := analyzer.calculateMax(metricName, data)
	outliers := analyzer.findOutliers(data)
	analyzer.log.Debugf("Stastical Indicators %s: Median: %f,  Mean: %f, Standard Deviation: %f, Q1: %f, Q3: %f Interquartile Range: %f, min: %f, max: %f, outliers: %+v", metricName, median, mean, standardDeviation, quartiles.Q1, quartiles.Q3, interQuartileRange, min, max, outliers)
}

func (analyzer *MetricAnalyzer) countOccurrences(values []float64, binSize float64) map[float64]int {
	counts := make(map[float64]int)
	for _, value := range values {
		// Round to the nearest bin.
		bin := math.Round(value/binSize) * binSize
		counts[bin]++
	}
	return counts
}

func (analyzer *MetricAnalyzer) createHistogramWithBoxPlot(metricName string, metrics []float64, folderName string) {
	p := plot.New()
	p.Title.Text = metricName + " histogram"
	p.X.Label.Text = "Value"
	p.Y.Label.Text = "Frequency"

	n := len(metrics)

	binSize := 0.0000001

	counts := analyzer.countOccurrences(metrics, binSize)

	metricPts := make(plotter.XYs, len(counts))
	index := 0
	for bin, count := range counts {
		metricPts[index] = plotter.XY{X: float64(bin), Y: float64(count)}
		index++
	}
	max := 0.0
	for _, count := range counts {
		if float64(count) > max {
			max = float64(count)
		}
	}

	h, err := plotter.NewHistogram(metricPts, len(counts))
	if err != nil {
		analyzer.log.Fatalf("Error creating histogram %s", err)
	}
	h.LogY = true
	p.Add(h)

	values := plotter.Values(metrics)
	boxPlot, err := plotter.NewBoxPlot(vg.Length(10), float64(n), values)
	if err != nil {
		analyzer.log.Fatalf("Error creating box plot %s", err)
	}
	boxPlot.Horizontal = true
	p.Add(boxPlot)

	if err := p.Save(8*vg.Inch, 8*vg.Inch, fmt.Sprintf("../png/%s/histogram_%s.png", folderName, metricName)); err != nil {
		analyzer.log.Fatalf("Error saving %s plot histogram to PNG", err)
	}
}

func (analyzer *MetricAnalyzer) createPlot(metricName string, metrics []float64, color color.RGBA, shape int, dashes []vg.Length, folderName string) (*plotter.Line, *plotter.Scatter) {
	p := plot.New()
	p.Title.Text = fmt.Sprintf("%s values", metricName)
	p.X.Label.Text = "Index"
	p.Y.Label.Text = metricName

	metricPts := make(plotter.XYs, len(metrics))
	for i, metric := range metrics {
		metricPts[i].X = float64(i)
		metricPts[i].Y = metric
	}
	metricLine, metricPoints, err := plotter.NewLinePoints(metricPts)
	if err != nil {
		analyzer.log.Fatalf("Error creating %s line plot %s", metricName, err)
	}

	metricLine.Color = color
	metricPoints.Color = color
	metricPoints.Shape = plotutil.Shape(shape)
	if dashes != nil {
		metricLine.Dashes = dashes
	}
	p.Add(metricLine, metricPoints)

	if err := p.Save(8*vg.Inch, 8*vg.Inch, fmt.Sprintf("../png/%s/%s_metrics.png", folderName, metricName)); err != nil {
		analyzer.log.Fatalf("Error saving %s plot to PNG %s", metricName, err)
	}
	return metricLine, metricPoints
}

func (anlayzer *MetricAnalyzer) combinePlots(latencyLine, jitterLine, packetLossLine *plotter.Line, latencyPoints, jitterPoints, packetLossPoints *plotter.Scatter, title, xLabel, yLabel string, filename string) {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	p.Add(latencyLine, latencyPoints, jitterLine, jitterPoints, packetLossLine, packetLossPoints)

	p.Legend.Add("latency", latencyLine, latencyPoints)
	p.Legend.Add("jitter", jitterLine, jitterPoints)
	p.Legend.Add("packet loss", packetLossLine, packetLossPoints)

	p.Legend.Top = true
	p.Legend.Left = true
	p.Legend.YOffs = -10

	if err := p.Save(8*vg.Inch, 8*vg.Inch, filename); err != nil {
		anlayzer.log.Errorf("Error saving plot to %s, error: %s", filename, err)
	}
}

func (analyzer *MetricAnalyzer) analyzeMetric(metricName, plotName, folderName string, metrics []float64, color color.RGBA, shape int, dashes []vg.Length) (*plotter.Line, *plotter.Scatter) {
	analyzer.log.Debugf("Analyzing %s %s ", plotName, metricName)
	analyzer.calculateStatisticalIndicators(metricName, metrics)
	analyzer.createHistogramWithBoxPlot(metricName, metrics, folderName)
	return analyzer.createPlot(metricName, metrics, color, shape, dashes, folderName)
}

func (analyzer *MetricAnalyzer) createAnalysis(latencyMetrics, jitterMetrics, packetLossMetrics []float64, plotName, folderName string) {
	latencyLine, latencyPoints := analyzer.analyzeMetric("latency", plotName, folderName, latencyMetrics, color.RGBA{R: 0, G: 126, B: 107, A: 255}, 0, nil)
	jitterLine, jitterPoints := analyzer.analyzeMetric("jitter", plotName, folderName, jitterMetrics, color.RGBA{R: 140, G: 25, B: 95, A: 255}, 1, []vg.Length{vg.Points(5), vg.Points(5)})
	packetLossLine, packetLossPoints := analyzer.analyzeMetric("packet_loss", plotName, folderName, packetLossMetrics, color.RGBA{R: 215, G: 40, B: 100, A: 255}, 2, []vg.Length{vg.Points(2), vg.Points(2)})
	analyzer.combinePlots(latencyLine, jitterLine, packetLossLine, latencyPoints, jitterPoints, packetLossPoints, plotName, "Index", "Values", fmt.Sprintf("../png/%s/combined_network_metrics.png", folderName))
}

// func (analyzer *MetricAnalyzer) applyStandardScale(metricName string, metrics []float64) []float64 {
// 	data := stats.LoadRawData(metrics)
// 	mean := analyzer.calculateMean(metricName, data)
// 	standardDeviation := analyzer.calculateStandardDeviation(metricName, data)

// 	normalizedWithStandard := make([]float64, len(metrics))
// 	for i, value := range metrics {
// 		normalizedWithStandard[i] = (value - mean) / standardDeviation
// 	}

// 	analyzer.log.Debugf("%s normalized with standard scaling: %v", metricName, normalizedWithStandard)

// 	return normalizedWithStandard
// }

// func (analyzer *MetricAnalyzer) analyzeMetricStandardScaled(metricName string, metrics []float64, color color.RGBA, shape int, dashes []vg.Length) (*plotter.Line, *plotter.Scatter) {
// 	analyzer.log.Debugf("Analyzing standard scaled %s metrics", metricName)
// 	standardScaledMetrics := analyzer.applyStandardScale(metricName, metrics)
// 	analyzer.calculateStatisticalIndicators(metricName, standardScaledMetrics)
// 	analyzer.createHistogramWithBoxPlot(metricName, standardScaledMetrics, "standard")
// 	return analyzer.createPlot(metricName, standardScaledMetrics, color, shape, dashes, "standard")
// }

// func (analyzer *MetricAnalyzer) analyzeStandardScaling() {
// 	latencyLine, latencyPoints := analyzer.analyzeMetricStandardScaled("latency", analyzer.currentLatencyMetrics, color.RGBA{R: 0, G: 126, B: 107, A: 255}, 0, nil)
// 	jitterLine, jitterPoints := analyzer.analyzeMetricStandardScaled("jitter", analyzer.currentJitterMetrics, color.RGBA{R: 140, G: 25, B: 95, A: 255}, 1, []vg.Length{vg.Points(5), vg.Points(5)})
// 	packetLossLine, packetLossPoints := analyzer.analyzeMetricStandardScaled("packet_loss", analyzer.currentPacketLossMetrics, color.RGBA{R: 215, G: 40, B: 100, A: 255}, 2, []vg.Length{vg.Points(2), vg.Points(2)})
// 	analyzer.combinePlots(latencyLine, jitterLine, packetLossLine, latencyPoints, jitterPoints, packetLossPoints, "Standard Scaling", "Index", "Values", "../png/standard/combined_network_metrics.png")
// }

// func (analyzer *MetricAnalyzer) applyMinMaxScaling(metricName string, metrics []float64, min, max float64) []float64 {
// 	normalizedWithMinMax := make([]float64, len(metrics))
// 	for i, value := range metrics {
// 		normalizedWithMinMax[i] = (value - min) / (max - min)
// 		if normalizedWithMinMax[i] > 1 {
// 			normalizedWithMinMax[i] = 1
// 		} else if normalizedWithMinMax[i] < 0 {
// 			normalizedWithMinMax[i] = 0
// 		}
// 	}
// 	analyzer.log.Debugf("%s normalized with min max scaling: %v", metricName, normalizedWithMinMax)
// 	return normalizedWithMinMax
// }

// func (analyzer *MetricAnalyzer) callMinMaxScaling(metricName string, metrics []float64) []float64 {
// 	min := analyzer.calculateMin(metricName, stats.LoadRawData(metrics))
// 	max := analyzer.calculateMax(metricName, stats.LoadRawData(metrics))
// 	return analyzer.applyMinMaxScaling(metricName, metrics, min, max)
// }

// func (analyzer *MetricAnalyzer) analyzeMetricMinMaxScaled(metricName string, metrics []float64, color color.RGBA, shape int, dashes []vg.Length) (*plotter.Line, *plotter.Scatter) {
// 	analyzer.log.Debugf("Analyzing min max scaled %s metrics", metricName)
// 	minMaxScaledMetrics := analyzer.callMinMaxScaling(metricName, metrics)
// 	analyzer.calculateStatisticalIndicators(metricName, minMaxScaledMetrics)
// 	analyzer.createHistogramWithBoxPlot(metricName, minMaxScaledMetrics, "minmax")
// 	return analyzer.createPlot(metricName, minMaxScaledMetrics, color, shape, dashes, "minmax")
// }

// func (analyzer *MetricAnalyzer) analyzeMinMaxScaling() {
// 	latencyLine, latencyPoints := analyzer.analyzeMetricMinMaxScaled("latency", analyzer.currentLatencyMetrics, color.RGBA{R: 0, G: 126, B: 107, A: 255}, 0, nil)
// 	jitterLine, jitterPoints := analyzer.analyzeMetricMinMaxScaled("jitter", analyzer.currentJitterMetrics, color.RGBA{R: 140, G: 25, B: 95, A: 255}, 1, []vg.Length{vg.Points(5), vg.Points(5)})
// 	packetLossLine, packetLossPoints := analyzer.analyzeMetricMinMaxScaled("packet_loss", analyzer.currentPacketLossMetrics, color.RGBA{R: 215, G: 40, B: 100, A: 255}, 2, []vg.Length{vg.Points(2), vg.Points(2)})
// 	analyzer.combinePlots(latencyLine, jitterLine, packetLossLine, latencyPoints, jitterPoints, packetLossPoints, "Min Max Scaling", "Index", "Values", "../png/minmax/combined_network_metrics.png")
// }

// func (analyzer *MetricAnalyzer) analyzeMetricIqrBasedMinMaxScaled(metricName string, metrics []float64, color color.RGBA, shape int, dashes []vg.Length) (*plotter.Line, *plotter.Scatter) {
// 	analyzer.log.Debugf("Analyzing IQR based min max scaled %s metrics", metricName)
// 	quartiles, interQuartileRange := analyzer.calculateQuartiles(metricName, metrics)
// 	min := math.Max(quartiles.Q1-1.5*interQuartileRange, 0)
// 	max := quartiles.Q3 + 1.5*interQuartileRange
// 	minMaxScaledMetrics := analyzer.applyMinMaxScaling(metricName, metrics, min, max)
// 	analyzer.calculateStatisticalIndicators(metricName, minMaxScaledMetrics)
// 	analyzer.createHistogramWithBoxPlot(metricName, minMaxScaledMetrics, "iqr-based-minmax")
// 	return analyzer.createPlot(metricName, minMaxScaledMetrics, color, shape, dashes, "iqr-based-minmax")
// }

// func (analyzer *MetricAnalyzer) analyzeIqrBasedMinMaxScaling() {
// 	latencyLine, latencyPoints := analyzer.analyzeMetricIqrBasedMinMaxScaled("latency", analyzer.currentLatencyMetrics, color.RGBA{R: 0, G: 126, B: 107, A: 255}, 0, nil)
// 	jitterLine, jitterPoints := analyzer.analyzeMetricIqrBasedMinMaxScaled("jitter", analyzer.currentJitterMetrics, color.RGBA{R: 140, G: 25, B: 95, A: 255}, 1, []vg.Length{vg.Points(5), vg.Points(5)})
// 	packetLossLine, packetLossPoints := analyzer.analyzeMetricIqrBasedMinMaxScaled("packet_loss", analyzer.currentPacketLossMetrics, color.RGBA{R: 215, G: 40, B: 100, A: 255}, 2, []vg.Length{vg.Points(2), vg.Points(2)})
// 	analyzer.combinePlots(latencyLine, jitterLine, packetLossLine, latencyPoints, jitterPoints, packetLossPoints, "Min Max Scaling IQR Based", "Index", "Values", "../png/iqr-based-minmax/combined_network_metrics.png")
// }

func (analyzer *MetricAnalyzer) Analyze() {
	analyzer.createAnalysis(analyzer.currentLatencyMetrics, analyzer.currentJitterMetrics, analyzer.currentPacketLossMetrics, "Original Metrics", "original")
	analyzer.createAnalysis(analyzer.normalizedLatencyMetrics, analyzer.normalizedJitterMetrics, analyzer.normalizedPacketLossMetrics, analyzer.plotName, analyzer.folderName)
}
