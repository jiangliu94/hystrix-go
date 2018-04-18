package calibrator

import (
	"time"
	"sync"
)

type DurationStream struct {
	// config stores the static configuration of this stream
	config Config
	mutex *sync.RWMutex
	// below stores dynamic data collected/calculated from the real time execution
	Buffer [][]time.Duration
	MeanBuffer []time.Duration
	DerivativeBuffer []float64
	CurrentDerivativeBufferIndex int
	LastUpdatedAt int64
	AccumulatedAverage time.Duration

}

func NewDurationStream(config Config) *DurationStream {
	buffer := make([][]time.Duration, config.AveragingWindowSize)
	for index := range buffer {
		buffer[index] = []time.Duration{}
	}
	meanBuffer := make([]time.Duration, config.AveragingWindowSize)
	derivativeBuffer := make([]float64, config.CalibrationWindowSize)

	return &DurationStream{
		mutex: 	&sync.RWMutex{},
		config: config,
		Buffer: buffer,
		MeanBuffer: meanBuffer,
		DerivativeBuffer: derivativeBuffer,
		CurrentDerivativeBufferIndex: 0,
		LastUpdatedAt: 0,
		AccumulatedAverage: 0,
	}

}

func (d *DurationStream) Append(value time.Duration) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	previousUpdate := d.LastUpdatedAt
	d.LastUpdatedAt = time.Now().Unix()
	index := d.getCurrentBufferIndex(d.LastUpdatedAt)

	if d.LastUpdatedAt != previousUpdate {
		d.Buffer[index] = []time.Duration{}
		d.updateAverage(d.LastUpdatedAt)
		d.updateDerivatives(d.LastUpdatedAt)
		// TODO: Implement the calibration logic
	}
	d.Buffer[index] = append(d.Buffer[index], value)
	d.updateMean(d.LastUpdatedAt, value)
}

func (d *DurationStream) getCurrentBufferIndex(now int64) int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return int(now % int64(len(d.Buffer)))
}

func (d *DurationStream) updateMean(now int64, value time.Duration) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	index := d.getCurrentBufferIndex(d.LastUpdatedAt)

	if len(d.Buffer[index]) <= 0 {
		d.MeanBuffer[index] = 0
		return
	}

	d.MeanBuffer[index] = (d.MeanBuffer[index] * time.Duration(len(d.Buffer[index]) - 1) + value) / time.Duration(len(d.Buffer[index]))
}

func (d *DurationStream) updateAverage(now int64) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	var sum time.Duration

	for _, value := range d.MeanBuffer {
		sum += value
	}

	size := len(d.Buffer) - 1
	if size <= 0 {
		size = 1
	}

	d.AccumulatedAverage = (sum - d.MeanBuffer[d.getCurrentBufferIndex(now)]) / time.Duration(size)

}


func (d *DurationStream) updateDerivatives(now int64) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	if len(d.Buffer) <= 2 || d.config.DerivativeCalculationInterval < 1 || d.config.DerivativeCalculationInterval % len(d.Buffer) == 0 {
		// nothing to calculate since the number of data point is too little
		// or it is always comparing the same data point
		return
	}
	if !(int(now) % d.config.CalibrationWindowSize == 0) {
		// not the right time to calculate the derivative
		return
	}

	lowerBound := d.getPreviousBufferIndex(now, d.config.DerivativeCalculationInterval + 1)
	upperBound := d.getPreviousBufferIndex(now, 1)
	d.DerivativeBuffer[d.getNewDerivativeBufferIndex()] = float64 (d.MeanBuffer[upperBound] - d.MeanBuffer[lowerBound]) / (float64 (d.config.DerivativeCalculationInterval) )
}

func (d *DurationStream) getPreviousBufferIndex(now int64, n int) int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	length := len(d.Buffer)
	return (int(now) + length - n % length) % length
}

func (d *DurationStream) getNewDerivativeBufferIndex() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return (d.CurrentDerivativeBufferIndex + 1) % len(d.DerivativeBuffer)
}

func (d *DurationStream) IsToScaleUp() bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for _, value := range d.DerivativeBuffer {
		if value < d.config.DerivativeThreshold.UpperBound {
			return false
		}
	}
	return float64(d.AccumulatedAverage) / d.config.PredefinedLimit.UpperBound >= d.config.UtilisationLimit.UpperBound
}

func (d *DurationStream) IsToScaleDown() bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for _, value := range d.DerivativeBuffer {
		if value > d.config.DerivativeThreshold.LowerBound {
			return false
		}
	}
	return float64(d.AccumulatedAverage) / d.config.PredefinedLimit.LowerBound <= d.config.UtilisationLimit.LowerBound
}

func (d *DurationStream) reset() {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for index := range d.Buffer {
		d.Buffer[index] = []time.Duration{}
		d.MeanBuffer[index] = 0
		d.DerivativeBuffer[index] = 0
	}
	d.CurrentDerivativeBufferIndex  = 0
	d.LastUpdatedAt = 0
	d.AccumulatedAverage = 0
}