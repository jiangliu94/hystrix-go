package commandbuilder

import (
	"time"

	"github.com/myteksi/hystrix-go/hystrix"
)

// CommandBuilder builder for constructing new command

type CommandBuilder struct {
	commandName            string
	timeout                int
	commandGroup           string
	maxConcurrentRequests  int
	requestVolumeThreshold int
	sleepWindow            int
	errorPercentThreshold  int
	// for more details refer - https://github.com/Netflix/Hystrix/wiki/Configuration#maxqueuesize
	queueSizeRejectionThreshold *int
}

// New Create new command
func New(commandName string) *CommandBuilder {
	return &CommandBuilder{
		commandName:                 commandName,
		timeout:                     hystrix.DefaultTimeout,
		commandGroup:                commandName,
		maxConcurrentRequests:       hystrix.DefaultMaxConcurrent,
		requestVolumeThreshold:      hystrix.DefaultVolumeThreshold,
		queueSizeRejectionThreshold: &hystrix.DefaultQueueSizeRejectionThreshold,
		sleepWindow:                 hystrix.DefaultSleepWindow,
		errorPercentThreshold:       hystrix.DefaultErrorPercentThreshold,
	}
}

// WithTimeout modify timeout
func (cb *CommandBuilder) WithTimeout(timeoutInMs int) *CommandBuilder {
	cb.timeout = timeoutInMs
	return cb
}

// WithCommandGroup modify commandGroup
func (cb *CommandBuilder) WithCommandGroup(commandGroup string) *CommandBuilder {
	cb.commandGroup = commandGroup
	return cb
}

// WithMaxConcurrentRequests modify max concurrent requests
func (cb *CommandBuilder) WithMaxConcurrentRequests(maxConcurrentRequests int) *CommandBuilder {
	cb.maxConcurrentRequests = maxConcurrentRequests
	return cb
}

// WithRequestVolumeThreshold modify request volume threshold
func (cb *CommandBuilder) WithRequestVolumeThreshold(requestVolThreshold int) *CommandBuilder {
	cb.requestVolumeThreshold = requestVolThreshold
	return cb
}

// WithSleepWindow modify sleep window
func (cb *CommandBuilder) WithSleepWindow(sleepWindow int) *CommandBuilder {
	cb.sleepWindow = sleepWindow
	return cb
}

// WithErrorPercentageThreshold modify error percentage threshold
func (cb *CommandBuilder) WithErrorPercentageThreshold(errPercentThreshold int) *CommandBuilder {
	cb.errorPercentThreshold = errPercentThreshold
	return cb
}

// WithQueueSize modify queue size
func (cb *CommandBuilder) WithQueueSize(queueSize int) *CommandBuilder {
	if queueSize == 0 {
		zeroQueueSize := 0
		cb.queueSizeRejectionThreshold = &zeroQueueSize
	}
	cb.queueSizeRejectionThreshold = &queueSize
	return cb
}

// Build the command setting, Use hystrix.Initialize for setup
func (cb *CommandBuilder) Build() *hystrix.Settings {
	return &hystrix.Settings{
		CommandName:                 cb.commandName,
		QueueSizeRejectionThreshold: *cb.queueSizeRejectionThreshold,
		ErrorPercentThreshold:       cb.errorPercentThreshold,
		CommandGroup:                cb.commandGroup,
		Timeout:                     time.Duration(cb.timeout * 1000000),
		MaxConcurrentRequests:       cb.maxConcurrentRequests,
		RequestVolumeThreshold:      uint64(cb.requestVolumeThreshold),
		SleepWindow:                 time.Duration(cb.sleepWindow * 1000000),
	}
}
