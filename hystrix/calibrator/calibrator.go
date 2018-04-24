package calibrator

import (
	"math"
	"sync"
	"time"

	"github.com/myteksi/hystrix-go/hystrix"
	"github.com/myteksi/hystrix-go/hystrix/metric_collector"
)

type Calibrator interface {
	Register(config Config)
	Calibrate()
}

type calibrator struct {
	mutex      *sync.RWMutex
	configs    map[string]*Config
	Collectors map[string]*Collector
	// TODO: provide datadog support for calibration event
}

func NewCalibrator() Calibrator {
	return &calibrator{
		mutex:      &sync.RWMutex{},
		configs:    map[string]*Config{},
		Collectors: map[string]*Collector{},
	}
}

func (c *calibrator) Register(config Config) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.configs[config.Name] = &config
	collector := NewCollector(config)
	c.Collectors[config.Name] = collector
	circuit, _, _ := hystrix.GetCircuit(config.Name)
	circuit.WithCollector(collector)
}

func (c *calibrator) Calibrate() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	settings := hystrix.GetCircuitSettings()

	for name, setting := range settings {
		go c.calibrate(name, setting)
	}
}

func (c *calibrator) calibrate(name string, setting *hystrix.Settings) {
	if setting == nil {
		return
	}
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	go func() {
		config, found := c.configs[name]
		collector, ok := c.Collectors[name]
		if !found || !ok {
			// not doing anything when there is no collector
			return
		}
		ticker := time.NewTicker(time.Duration(config.DerivativeCalculationInterval))

		for range ticker.C {
			changed := false
			if collector.numRequests.IsToScaleUp() {
				setting.MaxConcurrentRequests = scaleNumber(setting.MaxConcurrentRequests, true, config.AdjustmentThreshold)
				changed = true
			}
			if collector.numRequests.IsToScaleDown() {
				setting.MaxConcurrentRequests = scaleNumber(setting.MaxConcurrentRequests, false, config.AdjustmentThreshold)
				changed = true
			}
			if collector.runDuration.IsToScaleUp() {
				setting.Timeout = scaleTime(setting.Timeout, true, config.AdjustmentThreshold)
				changed = true
			}
			if collector.runDuration.IsToScaleDown() {
				setting.Timeout = scaleTime(setting.Timeout, false, config.AdjustmentThreshold)
				changed = true
			}
			if changed {
				hystrix.Initialize(setting)
			}
		}
	}()
}

func scaleTime(duration time.Duration, up bool, thresholdPercentage int64) time.Duration {
	if up {
		return time.Duration(math.Ceil(float64(duration.Nanoseconds()) * float64(100+thresholdPercentage) / (float64(100))))
	}
	return time.Duration(math.Ceil(float64(duration.Nanoseconds()) * float64(100-thresholdPercentage) / (float64(100))))
}

func scaleNumber(number int, up bool, thresholdPercentage int64) int {
	if up {
		return int(math.Ceil(float64(number) * float64(100+thresholdPercentage) / (float64(100))))
	}
	return int(math.Ceil(float64(number) * float64(100-thresholdPercentage) / (float64(100))))
}
