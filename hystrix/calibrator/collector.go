package calibrator

import (
	"sync"
	"time"
)

type Collector struct {
	mutex *sync.RWMutex
	// static configs
	CalibrationConfig Config

	// dynamic data collector
	numRequests *NumberStream
	errors      *NumberStream

	successes     *NumberStream
	queueSize     *NumberStream
	failures      *NumberStream
	rejects       *NumberStream
	shortCircuits *NumberStream
	timeouts      *NumberStream

	fallbackSuccesses *NumberStream
	fallbackFailures  *NumberStream
	totalDuration     *DurationStream
	runDuration       *DurationStream
}

func NewCollector(config Config) *Collector {
	cfg := config.validate()
	collector := &Collector{
		CalibrationConfig: cfg,
		mutex:             &sync.RWMutex{},
		numRequests:       NewNumberStream(cfg),
		errors:            NewNumberStream(cfg),
		successes:         NewNumberStream(cfg),
		queueSize:         NewNumberStream(cfg),
		failures:          NewNumberStream(cfg),
		rejects:           NewNumberStream(cfg),
		shortCircuits:     NewNumberStream(cfg),
		timeouts:          NewNumberStream(cfg),

		fallbackSuccesses: NewNumberStream(cfg),
		fallbackFailures:  NewNumberStream(cfg),
		totalDuration:     NewDurationStream(cfg),
		runDuration:       NewDurationStream(cfg),
	}
	return collector
}

// NumRequests returns the rolling number of requests
func (c *Collector) NumRequests() *NumberStream {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.numRequests
}

// QueueSize returns the rolling number of queue length
func (c *Collector) QueueSize() *NumberStream {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.queueSize
}

// Errors returns the rolling number of errors
func (c *Collector) Errors() *NumberStream {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.errors
}

// Successes returns the rolling number of successes
func (c *Collector) Successes() *NumberStream {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.successes
}

// Failures returns the rolling number of failures
func (c *Collector) Failures() *NumberStream {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.failures
}

// Rejects returns the rolling number of rejects
func (c *Collector) Rejects() *NumberStream {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.rejects
}

// ShortCircuits returns the rolling number of short circuits
func (c *Collector) ShortCircuits() *NumberStream {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.shortCircuits
}

// Timeouts returns the rolling number of timeouts
func (c *Collector) Timeouts() *NumberStream {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.timeouts
}

// FallbackSuccesses returns the rolling number of fallback successes
func (c *Collector) FallbackSuccesses() *NumberStream {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.fallbackSuccesses
}

// FallbackFailures returns the rolling number of fallback failures
func (c *Collector) FallbackFailures() *NumberStream {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.fallbackFailures
}

// TotalDuration returns the rolling total duration
func (c *Collector) TotalDuration() *DurationStream {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.totalDuration
}

// RunDuration returns the rolling run duration
func (c *Collector) RunDuration() *DurationStream {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.runDuration
}

func (c *Collector) IncrementAttempts() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.numRequests.Increment(1)
}

func (c *Collector) IncrementQueueSize() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.queueSize.Increment(1)
}

func (c *Collector) IncrementErrors() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.errors.Increment(1)
}

func (c *Collector) IncrementSuccesses() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.successes.Increment(1)
}

// IncrementFailures increments the number of failures seen in the latest time bucket.
func (c *Collector) IncrementFailures() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.failures.Increment(1)
}

// IncrementRejects increments the number of rejected requests seen in the latest time bucket.
func (c *Collector) IncrementRejects() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.rejects.Increment(1)
}

// IncrementShortCircuits increments the number of rejected requests seen in the latest time bucket.
func (c *Collector) IncrementShortCircuits() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.shortCircuits.Increment(1)
}

// IncrementTimeouts increments the number of requests that timec out in the latest time bucket.
func (c *Collector) IncrementTimeouts() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.timeouts.Increment(1)
}

// IncrementFallbackSuccesses increments the number of successful calls to the fallback function in the latest time bucket.
func (c *Collector) IncrementFallbackSuccesses() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.fallbackSuccesses.Increment(1)
}

// IncrementFallbackFailures increments the number of failed calls to the fallback function in the latest time bucket.
func (c *Collector) IncrementFallbackFailures() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.fallbackFailures.Increment(1)
}

// UpdateTotalDuration updates the total amount of time this circuit has been running.
func (c *Collector) UpdateTotalDuration(timeSinceStart time.Duration) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.totalDuration.Append(timeSinceStart)
}

// UpdateRunDuration updates the amount of time the latest request took to complete.
func (c *Collector) UpdateRunDuration(runDuration time.Duration) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.runDuration.Append(runDuration)
}

// Reset resets all metrics in this collector to 0.
func (c *Collector) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.numRequests.reset()
	c.errors.reset()
	c.successes.reset()
	c.rejects.reset()
	c.queueSize.reset()
	c.shortCircuits.reset()
	c.failures.reset()
	c.timeouts.reset()
	c.fallbackSuccesses.reset()
	c.fallbackFailures.reset()
	c.totalDuration.reset()
	c.runDuration.reset()
}
