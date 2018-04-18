package calibrator

import (
	"sync"
	"time"
	"github.com/myteksi/hystrix-go/hystrix/metric_collector"
	"gitlab.myteksi.net/gophers/go/commons/util/resilience/hystrix/metric_collector"
)

type Collector struct {
	mutex *sync.RWMutex
	CalibrationConfig Config
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

func NewCollector(config Config) metricCollector.MetricCollector {
	cfg := config.validate()
	collector := &Collector{
		CalibrationConfig: cfg,
		mutex: 	&sync.RWMutex{},
		numRequests: NewNumberStream(cfg),
		errors: NewNumberStream(cfg),
		successes: NewNumberStream(cfg),
		queueSize: NewNumberStream(cfg),
		failures: NewNumberStream(cfg),
		rejects: NewNumberStream(cfg),
		shortCircuits: NewNumberStream(cfg),
		timeouts: NewNumberStream(cfg),

		fallbackSuccesses: NewNumberStream(cfg),
		fallbackFailures: NewNumberStream(cfg),
		totalDuration: NewDurationStream(cfg),
		runDuration: NewDurationStream(cfg),
	}
	return collector
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
