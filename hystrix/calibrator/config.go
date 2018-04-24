package calibrator

const (
	// DefaultDerivativeCalculationInterval is the constant controlling the period in seconds of stats derivative calculation
	DefaultDerivativeCalculationInterval = 1

	// DefaultCalibrationWindowSize is the constant controlling the number of consecutive events of triggering the circuit calculation
	DefaultCalibrationWindowSize = 60

	// DefaultAveragingWindowSize controls the window for averaging the data
	DefaultAveragingWindowSize = int64(3600)

	// DefaultAdjustmentThreshold controls the percentage of config calibration. e.g 5%
	DefaultAdjustmentThreshold = int64(5)
)

type Config struct {
	Name                          string
	DerivativeCalculationInterval int
	CalibrationWindowSize         int
	AveragingWindowSize           int64
	DerivativeThreshold           Limit
	PredefinedLimit               Limit
	UtilisationLimit              Limit
	AdjustmentThreshold           int64
}

func (c Config) validate() Config {
	dci := DefaultDerivativeCalculationInterval
	if c.DerivativeCalculationInterval >= dci {
		dci = c.DerivativeCalculationInterval
	}
	cws := DefaultCalibrationWindowSize
	if c.CalibrationWindowSize > 0 {
		cws = c.CalibrationWindowSize
	}
	aws := DefaultAveragingWindowSize
	if c.AveragingWindowSize > 0 {
		aws = c.AveragingWindowSize
	}
	at := DefaultAdjustmentThreshold
	if c.AdjustmentThreshold < 100 && c.AdjustmentThreshold >= 0 {
		at = c.AdjustmentThreshold
	}

	return Config{
		Name: c.Name,
		DerivativeCalculationInterval: dci,
		CalibrationWindowSize:         cws,
		AveragingWindowSize:           aws,
		DerivativeThreshold:           c.DerivativeThreshold.Validate(),
		PredefinedLimit:               c.PredefinedLimit.Validate(),
		UtilisationLimit:              c.UtilisationLimit.Validate(),
		AdjustmentThreshold:           at,
	}
}
