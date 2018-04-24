package calibrator

import "math"

type Limit struct {
	UpperBound float64
	LowerBound float64
}

func NewLimit() Limit {
	return Limit{
		UpperBound: math.MaxFloat64,
		LowerBound: -math.MaxFloat64,
	}
}

func (l Limit) Validate() Limit {
	if l.LowerBound > l.UpperBound {
		return NewLimit()
	}
	return l
}
