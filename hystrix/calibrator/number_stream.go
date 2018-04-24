package calibrator

import (
	"sync"
	"time"
)

type NumberStream struct {
	// config stores the static configuration of this stream
	config Config
	mutex  *sync.RWMutex
	// below stores dynamic data collected/calculated from the real time execution
	Buffer                       []float64
	DerivativeBuffer             []float64
	CurrentDerivativeBufferIndex int
	LastUpdatedAt                int64
	AccumulatedAverage           float64
}

func (s *NumberStream) Increment(value float64) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	previousUpdate := s.LastUpdatedAt
	s.LastUpdatedAt = time.Now().Unix()

	if s.LastUpdatedAt != previousUpdate {
		s.Buffer[s.getCurrentBufferIndex(s.LastUpdatedAt)] = 0
		s.updateAverage(s.LastUpdatedAt)
		s.updateDerivatives(s.LastUpdatedAt)
	}
	s.Buffer[s.getCurrentBufferIndex(s.LastUpdatedAt)] += value
}

func NewNumberStream(config Config) *NumberStream {
	buffer := make([]float64, config.AveragingWindowSize)
	derivativeBuffer := make([]float64, config.CalibrationWindowSize)
	return &NumberStream{
		mutex:                        &sync.RWMutex{},
		config:                       config,
		Buffer:                       buffer,
		DerivativeBuffer:             derivativeBuffer,
		CurrentDerivativeBufferIndex: 0,
		LastUpdatedAt:                0,
		AccumulatedAverage:           0,
	}

}

func (s *NumberStream) reset() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.Buffer = make([]float64, s.config.AveragingWindowSize)
	s.DerivativeBuffer = make([]float64, s.config.CalibrationWindowSize)
	s.CurrentDerivativeBufferIndex = 0
	s.LastUpdatedAt = 0
	s.AccumulatedAverage = 0

}

func (s *NumberStream) getCurrentBufferIndex(now int64) int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return int(now % int64(len(s.Buffer)))
}

func (s *NumberStream) getPreviousBufferIndex(now int64, n int) int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	length := len(s.Buffer)
	return (int(now) + length - n%length) % length
}

func (s *NumberStream) updateAverage(now int64) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var sum float64

	for _, value := range s.Buffer {
		sum += value
	}

	size := len(s.Buffer) - 1
	if size <= 0 {
		size = 1
	}

	s.AccumulatedAverage = (sum - s.Buffer[s.getCurrentBufferIndex(now)]) / float64(size)
}

func (s *NumberStream) updateDerivatives(now int64) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(s.Buffer) <= 2 || s.config.DerivativeCalculationInterval < 1 || s.config.DerivativeCalculationInterval%len(s.Buffer) == 0 {
		// nothing to calculate since the number of data point is too little
		// or it is always comparing the same data point
		return
	}
	if int(now)%s.config.DerivativeCalculationInterval != 0 {
		// not the right time to calculate the derivative
		return
	}

	lowerBound := s.getPreviousBufferIndex(now, s.config.DerivativeCalculationInterval+1)
	upperBound := s.getPreviousBufferIndex(now, 1)
	s.DerivativeBuffer[s.getNewDerivativeBufferIndex()] = (s.Buffer[upperBound] - s.Buffer[lowerBound]) / (float64(s.config.DerivativeCalculationInterval))

}

func (s *NumberStream) getNewDerivativeBufferIndex() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	s.CurrentDerivativeBufferIndex = (s.CurrentDerivativeBufferIndex + 1) % len(s.DerivativeBuffer)
	return s.CurrentDerivativeBufferIndex
}

func (s *NumberStream) IsToScaleUp() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, value := range s.DerivativeBuffer {
		if value < s.config.DerivativeThreshold.UpperBound {
			return false
		}
	}
	utilisation := s.AccumulatedAverage / s.config.PredefinedLimit.UpperBound
	// Only scale up when
	// 1. long term average above threshold
	// 2. long term average not reaching the hard upper limit
	// 3. the trend is increasing slowly and gently
	return utilisation >= s.config.UtilisationLimit.UpperBound && utilisation < 1
}

func (s *NumberStream) IsToScaleDown() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, value := range s.DerivativeBuffer {
		if value > s.config.DerivativeThreshold.LowerBound {
			return false
		}
	}

	// has to be upperbound here
	utilisation := s.AccumulatedAverage / s.config.PredefinedLimit.UpperBound
	// Only scale down when
	// 1. long term average is below the threshold
	// 2. long term average is not reaching the hard lower limit
	// 3. the trend is decreasing
	return utilisation <= s.config.UtilisationLimit.LowerBound && s.AccumulatedAverage > s.config.PredefinedLimit.LowerBound
}
