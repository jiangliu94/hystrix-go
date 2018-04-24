package calibrator

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewDurationStream(t *testing.T) {
	Convey("It should successfully setup a stream", t, func() {
		config := Config{
			Name: "Test Circuit Name",
			DerivativeCalculationInterval: 1,  // calculating derivative every second
			CalibrationWindowSize:         5,  //triggering calibration evaluation every 5 second
			AveragingWindowSize:           10, // keeping data in the last 10 seconds
			DerivativeThreshold: Limit{
				// scale down whenever it decreases
				LowerBound: 0,
				// scale up whenever it increases
				UpperBound: 0,
			},
			// no hard limit of upperbound and lowerbound
			PredefinedLimit: NewLimit(),
			UtilisationLimit: Limit{
				// only consider scale in when long term average is lower than 40% of hard limit
				LowerBound: 0.4,
				// only consider scale up when long term average is above 60% of hard limit
				UpperBound: 0.6,
			},
		}
		testStream := NewDurationStream(config)
		So(testStream.AccumulatedAverage, ShouldEqual, float64(0))
		So(len(testStream.Buffer), ShouldEqual, int64(10))
		So(len(testStream.MeanBuffer), ShouldEqual, 10)
		So(len(testStream.DerivativeBuffer), ShouldEqual, 5)
		So(testStream.CurrentDerivativeBufferIndex, ShouldEqual, float64(0))

	})
}

func TestAppend(t *testing.T) {

	Convey("when adding values to a number stream", t, func() {
		config := Config{
			Name: "Test Circuit Name",
			DerivativeCalculationInterval: 1,  // calculating derivative every second
			CalibrationWindowSize:         5,  //triggering calibration evaluation every 5 second
			AveragingWindowSize:           11, // keeping data in the last 10 seconds
			DerivativeThreshold: Limit{
				// scale down whenever it decreases
				LowerBound: 0,
				// scale up whenever it increases
				UpperBound: 0,
			},
			// no hard limit of upperbound and lowerbound
			PredefinedLimit: NewLimit(),
			UtilisationLimit: Limit{
				// only consider scale in when long term average is lower than 40% of hard limit
				LowerBound: 0.4,
				// only consider scale up when long term average is above 60% of hard limit
				UpperBound: 0.6,
			},
		}
		testStream := NewDurationStream(config)
		now := time.Now().Unix()
		for _, x := range []time.Duration{19, 11, 9, 10, 13, 7, 4, 9, 2, 15} {
			testStream.Append(x * time.Second)
			time.Sleep(1 * time.Second)
		}

		Convey("it should have the correct data and calculation", func() {
			So(testStream.LastUpdatedAt-now, ShouldBeGreaterThanOrEqualTo, 9)
			So(testStream.AccumulatedAverage, ShouldEqual, time.Duration(8*time.Second)+time.Duration(400*time.Millisecond))
			// current index is purely dependent on the actual test time, so not putting it into the test
			// the derivative buffer should be a ring of [(13-10), (7-13), (4-7), (9-4), (2-9)] with random starting index depends on the actual test run time
		})
	})
}

func TestDurationStream_IsToScaleUp(t *testing.T) {
	Convey("when adding values to a rolling number", t, func() {
		config := Config{
			Name: "Test Circuit Name",
			DerivativeCalculationInterval: 1, // calculating derivative every second
			CalibrationWindowSize:         3, //triggering calibration evaluation every 3 second
			AveragingWindowSize:           5, // keeping data in the last 5 seconds
			DerivativeThreshold: Limit{
				// scale down whenever it decreases
				LowerBound: 0,
				// scale up whenever it increases
				UpperBound: 0,
			},
			PredefinedLimit: Limit{
				LowerBound: 0,
				UpperBound: 5,
			},
			UtilisationLimit: Limit{
				// only consider scale in when long term average is lower than 40% of hard limit
				LowerBound: 0.4,
				// only consider scale up when long term average is above 60% of hard limit
				UpperBound: 0.6,
			},
		}
		testStream := NewDurationStream(config)
		for _, x := range []time.Duration{1, 2, 3, 4, 5, 4} {
			testStream.Append(x * time.Millisecond)
			time.Sleep(1 * time.Second)
		}

		Convey("it should scale up", func() {
			So(testStream.IsToScaleUp(), ShouldEqual, true)
		})
	})
}

func TestDurationStream_IsToScaleDown(t *testing.T) {
	Convey("when adding values to a rolling number", t, func() {
		config := Config{
			Name: "Test Circuit Name",
			DerivativeCalculationInterval: 1, // calculating derivative every second
			CalibrationWindowSize:         3, //triggering calibration evaluation every 3 second
			AveragingWindowSize:           5, // keeping data in the last 5 seconds
			DerivativeThreshold: Limit{
				// scale down whenever it decreases
				LowerBound: 0,
				// scale up whenever it increases
				UpperBound: 0,
			},
			PredefinedLimit: Limit{
				LowerBound: 20,
				UpperBound: 100,
			},
			UtilisationLimit: Limit{
				// only consider scale in when long term average is lower than 60% of hard limit
				LowerBound: 0.6,
				// only consider scale up when long term average is above 80% of hard limit
				UpperBound: 0.8,
			},
		}
		testStream := NewDurationStream(config)
		for _, x := range []time.Duration{60, 50, 40, 30, 20, 5} {
			testStream.Append(x * time.Millisecond)
			time.Sleep(1 * time.Second)
		}

		Convey("it should scale down", func() {
			So(testStream.IsToScaleDown(), ShouldEqual, true)
		})
	})
}
