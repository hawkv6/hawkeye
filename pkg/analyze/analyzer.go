package analyze

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/hawkv6/hawkeye/pkg/processor"
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
	log               *logrus.Entry
	processor         processor.Processor
	latencyMetrics    []float64
	jitterMetrics     []float64
	packetLossMetrics []float64
}

func NewMetricAnalyzer(processor processor.Processor) *MetricAnalyzer {
	return &MetricAnalyzer{
		log:               logging.DefaultLogger.WithField("subsystem", "analyze"),
		processor:         processor,
		latencyMetrics:    make([]float64, 0),
		jitterMetrics:     make([]float64, 0),
		packetLossMetrics: make([]float64, 0),
	}
}

func (analyzer *MetricAnalyzer) calculateMedian(metricName string, data stats.Float64Data) float64 {
	median, err := stats.Median(data)
	if err != nil {
		analyzer.log.Fatalf("Error calculating median %s", metricName)
	}
	analyzer.log.Debugf("Median %s: %f", metricName, median)
	return median
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

func (analyzer *MetricAnalyzer) calculateInterQuartileRange(metricName string, data stats.Float64Data) float64 {
	iqr, err := stats.InterQuartileRange(data)
	if err != nil {
		analyzer.log.Fatalf("Error calculating IQR %s", metricName)
	}
	analyzer.log.Debugf("IQR %s: %f", metricName, iqr)
	return iqr
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

func (analyzer *MetricAnalyzer) calculateStatisticalIndicators(metricName string, metrics []float64) {
	data := stats.LoadRawData(metrics)
	median := analyzer.calculateMedian(metricName, data)
	mean := analyzer.calculateMean(metricName, data)
	standardDeviation := analyzer.calculateStandardDeviation(metricName, data)
	interQuartileRange := analyzer.calculateInterQuartileRange(metricName, data)
	min := analyzer.calculateMin(metricName, data)
	max := analyzer.calculateMax(metricName, data)
	analyzer.log.Debugf("Stastical Indicators %s: Median: %f,  Mean: %f, Standard Deviation: %f, Interquartile Range: %f, min: %f, max: %f", metricName, median, mean, standardDeviation, interQuartileRange, min, max)

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

	binSize := 0.001

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

func (analyzer *MetricAnalyzer) analyzeOriginalMetric(metricName string, metrics []float64, color color.RGBA, shape int, dashes []vg.Length) (*plotter.Line, *plotter.Scatter) {
	analyzer.log.Debugf("Analyzing original %s metrics", metricName)
	analyzer.calculateStatisticalIndicators(metricName, metrics)
	analyzer.createHistogramWithBoxPlot(metricName, metrics, "original")
	return analyzer.createPlot(metricName, metrics, color, shape, dashes, "original")
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

func (analyzer *MetricAnalyzer) analyzeOriginalMetrics() {
	latencyLine, latencyPoints := analyzer.analyzeOriginalMetric("latency", analyzer.latencyMetrics, color.RGBA{R: 0, G: 126, B: 107, A: 255}, 0, nil)
	jitterLine, jitterPoints := analyzer.analyzeOriginalMetric("jitter", analyzer.jitterMetrics, color.RGBA{R: 140, G: 25, B: 95, A: 255}, 1, []vg.Length{vg.Points(5), vg.Points(5)})
	packetLossLine, packetLossPoints := analyzer.analyzeOriginalMetric("packet_loss", analyzer.packetLossMetrics, color.RGBA{R: 215, G: 40, B: 100, A: 255}, 2, []vg.Length{vg.Points(2), vg.Points(2)})

	analyzer.combinePlots(latencyLine, jitterLine, packetLossLine, latencyPoints, jitterPoints, packetLossPoints, "Original Network Metrics", "Index", "Values", "../png/original/combined_network_metrics.png")
}

func (analyzer *MetricAnalyzer) findOutliers(data stats.Float64Data) {
	outliers, err := data.QuartileOutliers()
	if err != nil {
		fmt.Println("Error in outliers calculation: ", err)
	}
	fmt.Println("Outliers: ", outliers)
}

func (analyzer *MetricAnalyzer) applyRobustNormalization(metricName string, metrics []float64) []float64 {
	data := stats.LoadRawData(metrics)
	median := analyzer.calculateMedian(metricName, data)
	interQuartileRange := analyzer.calculateInterQuartileRange(metricName, data)

	normalizedWithRobust := make([]float64, len(metrics))
	for i, value := range metrics {
		normalizedWithRobust[i] = (value - median) / interQuartileRange
	}

	analyzer.log.Debugf("%s normalized with robust: %v", metricName, normalizedWithRobust)

	analyzer.findOutliers(data)

	return normalizedWithRobust
}

func (analyzer *MetricAnalyzer) analyzeMetricRobustNormalized(metricName string, metrics []float64, color color.RGBA, shape int, dashes []vg.Length) (*plotter.Line, *plotter.Scatter) {
	analyzer.log.Debugf("Analyzing robust normalized %s metrics", metricName)
	robustNormalizedMetrics := analyzer.applyRobustNormalization(metricName, metrics)
	analyzer.calculateStatisticalIndicators(metricName, robustNormalizedMetrics)
	analyzer.createHistogramWithBoxPlot(metricName, robustNormalizedMetrics, "robust")
	return analyzer.createPlot(metricName, robustNormalizedMetrics, color, shape, dashes, "robust")
}

func (analyzer *MetricAnalyzer) analyzeRobustNormalization() {
	latencyLine, latencyPoints := analyzer.analyzeMetricRobustNormalized("latency", analyzer.latencyMetrics, color.RGBA{R: 0, G: 126, B: 107, A: 255}, 0, nil)
	jitterLine, jitterPoints := analyzer.analyzeMetricRobustNormalized("jitter", analyzer.jitterMetrics, color.RGBA{R: 140, G: 25, B: 95, A: 255}, 1, []vg.Length{vg.Points(5), vg.Points(5)})
	packetLossLine, packetLossPoints := analyzer.analyzeMetricRobustNormalized("packet_loss", analyzer.packetLossMetrics, color.RGBA{R: 215, G: 40, B: 100, A: 255}, 2, []vg.Length{vg.Points(2), vg.Points(2)})
	analyzer.combinePlots(latencyLine, jitterLine, packetLossLine, latencyPoints, jitterPoints, packetLossPoints, "Robust Normalization", "Index", "Values", "../png/robust/combined_network_metrics.png")
}

func (analyzer *MetricAnalyzer) applyStandardScale(metricName string, metrics []float64) []float64 {
	data := stats.LoadRawData(metrics)
	mean := analyzer.calculateMean(metricName, data)
	standardDeviation := analyzer.calculateStandardDeviation(metricName, data)

	normalizedWithStandard := make([]float64, len(metrics))
	for i, value := range metrics {
		normalizedWithStandard[i] = (value - mean) / standardDeviation
	}

	analyzer.log.Debugf("%s normalized with standard scaling: %v", metricName, normalizedWithStandard)

	return normalizedWithStandard
}

func (analyzer *MetricAnalyzer) analyzeMetricStandardScaled(metricName string, metrics []float64, color color.RGBA, shape int, dashes []vg.Length) (*plotter.Line, *plotter.Scatter) {
	analyzer.log.Debugf("Analyzing standard scaled %s metrics", metricName)
	standardScaledMetrics := analyzer.applyStandardScale(metricName, metrics)
	analyzer.calculateStatisticalIndicators(metricName, standardScaledMetrics)
	analyzer.createHistogramWithBoxPlot(metricName, standardScaledMetrics, "standard")
	return analyzer.createPlot(metricName, standardScaledMetrics, color, shape, dashes, "standard")
}

func (analyzer *MetricAnalyzer) analyzeStandardScaling() {
	latencyLine, latencyPoints := analyzer.analyzeMetricStandardScaled("latency", analyzer.latencyMetrics, color.RGBA{R: 0, G: 126, B: 107, A: 255}, 0, nil)
	jitterLine, jitterPoints := analyzer.analyzeMetricStandardScaled("jitter", analyzer.jitterMetrics, color.RGBA{R: 140, G: 25, B: 95, A: 255}, 1, []vg.Length{vg.Points(5), vg.Points(5)})
	packetLossLine, packetLossPoints := analyzer.analyzeMetricStandardScaled("packet_loss", analyzer.packetLossMetrics, color.RGBA{R: 215, G: 40, B: 100, A: 255}, 2, []vg.Length{vg.Points(2), vg.Points(2)})
	analyzer.combinePlots(latencyLine, jitterLine, packetLossLine, latencyPoints, jitterPoints, packetLossPoints, "Standard Scaling", "Index", "Values", "../png/standard/combined_network_metrics.png")
}

func (analyzer *MetricAnalyzer) applyMinMaxScaling(metricName string, metrics []float64) []float64 {
	min := analyzer.calculateMin(metricName, stats.LoadRawData(metrics))
	max := analyzer.calculateMax(metricName, stats.LoadRawData(metrics))

	normalizedWithMinMax := make([]float64, len(metrics))
	for i, value := range metrics {
		normalizedWithMinMax[i] = (value - min) / (max - min)
	}

	analyzer.log.Debugf("%s normalized with min max scaling: %v", metricName, normalizedWithMinMax)

	return normalizedWithMinMax
}

func (analyzer *MetricAnalyzer) analyzeMetricMinMaxScaled(metricName string, metrics []float64, color color.RGBA, shape int, dashes []vg.Length) (*plotter.Line, *plotter.Scatter) {
	analyzer.log.Debugf("Analyzing min max scaled %s metrics", metricName)
	minMaxScaledMetrics := analyzer.applyMinMaxScaling(metricName, metrics)
	analyzer.calculateStatisticalIndicators(metricName, minMaxScaledMetrics)
	analyzer.createHistogramWithBoxPlot(metricName, minMaxScaledMetrics, "minmax")
	return analyzer.createPlot(metricName, minMaxScaledMetrics, color, shape, dashes, "minmax")
}

func (analyzer *MetricAnalyzer) analyzeMinMaxScaling() {
	latencyLine, latencyPoints := analyzer.analyzeMetricMinMaxScaled("latency", analyzer.latencyMetrics, color.RGBA{R: 0, G: 126, B: 107, A: 255}, 0, nil)
	jitterLine, jitterPoints := analyzer.analyzeMetricMinMaxScaled("jitter", analyzer.jitterMetrics, color.RGBA{R: 140, G: 25, B: 95, A: 255}, 1, []vg.Length{vg.Points(5), vg.Points(5)})
	packetLossLine, packetLossPoints := analyzer.analyzeMetricMinMaxScaled("packet_loss", analyzer.packetLossMetrics, color.RGBA{R: 215, G: 40, B: 100, A: 255}, 2, []vg.Length{vg.Points(2), vg.Points(2)})
	analyzer.combinePlots(latencyLine, jitterLine, packetLossLine, latencyPoints, jitterPoints, packetLossPoints, "Min Max Scaling", "Index", "Values", "../png/minmax/combined_network_metrics.png")
}

func (analyzer *MetricAnalyzer) Visualize() {
	analyzer.latencyMetrics = analyzer.processor.GetLatencyMetrics()
	analyzer.jitterMetrics = analyzer.processor.GetJitterMetrics()
	analyzer.packetLossMetrics = analyzer.processor.GetPacketLossMetrics()

	analyzer.analyzeOriginalMetrics()
	analyzer.analyzeRobustNormalization()
	analyzer.analyzeStandardScaling()
	analyzer.analyzeMinMaxScaling()
}
