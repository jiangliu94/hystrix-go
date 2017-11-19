package plugins

import (
	"sync/atomic"
	"testing"

	"github.com/myteksi/hystrix-go/plugins/mocks"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestNewDatadogCollector(t *testing.T) {
	mockDatadog := &mocks.DatadogClient{}
	attempt1 := int32(0)
	error1 := int32(0)
	mockDatadog.On("Count", "hystrix.attempts", int64(1), []string{"hystrixcircuit:commandName1", "commandGroup:commandGroup1"}, float64(1)).Run(func(args mock.Arguments) {
		atomic.AddInt32(&attempt1, 1)
	}).Return(nil)
	mockDatadog.On("Count", "hystrix.errors", int64(1), []string{"hystrixcircuit:commandName1", "commandGroup:commandGroup1"}, float64(1)).Run(func(args mock.Arguments) {
		atomic.AddInt32(&error1, 1)
	}).Return(nil)
	metricCollector := NewDatadogCollectorWithClient(mockDatadog)("commandName1", "commandGroup1")

	Convey("increment queue size", t, func() {
		metricCollector.IncrementAttempts()
		metricCollector.IncrementAttempts()
		metricCollector.IncrementErrors()

		So(2, ShouldEqual, attempt1)
		So(1, ShouldEqual, error1)
	})
}
