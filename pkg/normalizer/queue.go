package normalizer

type Queue interface {
	EnqueueQ1(float64)
	EnqueueQ3(float64)
	EnqueueMin(float64)
	EnqueueMax(float64)
	GetAverageQ1() float64
	GetAverageQ3() float64
	GetAverageMax() float64
	GetAverageMin() float64
	GetQ1Elements() []float64
	GetQ3Elements() []float64
	GetMaxElements() []float64
	GetMinElements() []float64
}

type NormalizationQueue struct {
	rollingWindowSize int
	q1Elements        []float64
	q1Sum             float64
	q3Elements        []float64
	q3Sum             float64
	minElements       []float64
	minSum            float64
	maxElements       []float64
	maxSum            float64
}

func NewNormalizationQueue(rollingWindowSize int) *NormalizationQueue {
	return &NormalizationQueue{
		rollingWindowSize: rollingWindowSize,
		q1Elements:        make([]float64, 0, rollingWindowSize),
		q1Sum:             0,
		q3Elements:        make([]float64, 0, rollingWindowSize),
		q3Sum:             0,
		minElements:       make([]float64, 0, rollingWindowSize),
		minSum:            0,
		maxElements:       make([]float64, 0, rollingWindowSize),
		maxSum:            0,
	}
}

func (queue *NormalizationQueue) enqueue(elements *[]float64, sum *float64, value float64) {
	if len(*elements) == queue.rollingWindowSize {
		*sum -= (*elements)[0]
		*elements = (*elements)[1:]
	}
	*elements = append(*elements, value)
	*sum += value
}

func (queue *NormalizationQueue) EnqueueQ1(value float64) {
	queue.enqueue(&queue.q1Elements, &queue.q1Sum, value)
}

func (queue *NormalizationQueue) EnqueueQ3(value float64) {
	queue.enqueue(&queue.q3Elements, &queue.q3Sum, value)
}

func (queue *NormalizationQueue) EnqueueMin(value float64) {
	queue.enqueue(&queue.minElements, &queue.minSum, value)
}

func (queue *NormalizationQueue) EnqueueMax(value float64) {
	queue.enqueue(&queue.maxElements, &queue.maxSum, value)
}

func (queue *NormalizationQueue) getAverage(elements []float64, sum float64) float64 {
	if len(elements) == 0 {
		return 0
	}
	return sum / float64(len(elements))
}

func (queue *NormalizationQueue) GetAverageQ1() float64 {
	return queue.getAverage(queue.q1Elements, queue.q1Sum)
}

func (queue *NormalizationQueue) GetAverageQ3() float64 {
	return queue.getAverage(queue.q3Elements, queue.q3Sum)
}

func (queue *NormalizationQueue) GetAverageMin() float64 {
	return queue.getAverage(queue.minElements, queue.minSum)
}

func (queue *NormalizationQueue) GetAverageMax() float64 {
	return queue.getAverage(queue.maxElements, queue.maxSum)
}

func (queue *NormalizationQueue) GetQ1Elements() []float64 {
	return queue.q1Elements
}

func (queue *NormalizationQueue) GetQ3Elements() []float64 {
	return queue.q3Elements
}

func (queue *NormalizationQueue) GetMinElements() []float64 {
	return queue.minElements
}

func (queue *NormalizationQueue) GetMaxElements() []float64 {
	return queue.maxElements
}
