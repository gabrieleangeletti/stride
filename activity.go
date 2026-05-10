package stride

import (
	"errors"
	"math"
	"time"
)

type ActivityConvertible interface {
	ToActivity() (*Activity, error)
}

type Activity struct {
	Provider      Provider
	Sport         Sport
	StartTime     time.Time        // UTC
	ElapsedTime   uint32           // seconds
	MovingTime    uint32           // seconds
	Distance      uint32           // meters
	AvgSpeed      uint16           // m / s
	AvgHR         Optional[uint8]  // beats / minute
	MaxHR         Optional[uint8]  // beats / minute
	ElevationGain Optional[uint16] // meters
	ElevationLoss Optional[uint16] // meters
}

type ActivityTimeseriesConvertible interface {
	ToTimeseries(startTime time.Time) (*ActivityTimeseries, error)
}

type ActivityTimeseries struct {
	StartTime time.Time
	Data      []ActivityTimeseriesEntry
}

func (ts ActivityTimeseries) EndTime() time.Time {
	return ts.StartTime.Add(time.Duration(ts.MaxOffset()) * time.Second)
}

func (ts ActivityTimeseries) ElapsedTime() time.Duration {
	return ts.EndTime().Sub(ts.StartTime)
}

func (ts ActivityTimeseries) MaxOffset() int {
	maxOffset := 0

	for _, entry := range ts.Data {
		if entry.Offset > maxOffset {
			maxOffset = entry.Offset
		}
	}

	return maxOffset
}

type hrMetrics struct {
	AvgHR int16
	MaxHR int16
}

func (ts *ActivityTimeseries) HRMetrics() (*hrMetrics, error) {
	avgHR, err := CalculateAverageHeartRate(ts, AvgHeartRateAnalysisConfig{
		Method:       HeartRateMethodTimeWeighted,
		ExcludeZeros: true,
		MinValidRate: 40,
		MaxValidRate: 220,
		MaxHeartRate: 193,
	})
	if err != nil {
		if !errors.Is(err, ErrNoValidData) {
			return nil, err
		}
	}

	var avgValue int16
	if avgHR > 0 {
		avgValue = int16(math.Round(avgHR))
	}

	thirtySec := 30 * time.Second

	maxHR, err := CalculateMaxHeartRate(ts, MaxHeartRateAnalysisConfig{
		Method:         MaxHeartRateMethodRollingWindow,
		WindowDuration: &thirtySec,
	})
	if err != nil {
		if !errors.Is(err, ErrNoValidData) {
			return nil, err
		}
	}

	var maxValue int16
	if maxHR > 0 {
		maxValue = int16(maxHR)
	}

	return &hrMetrics{
		AvgHR: avgValue,
		MaxHR: maxValue,
	}, nil
}

type ActivityTimeseriesEntry struct {
	Offset    int
	HeartRate Optional[uint8]
	Cadence   Optional[uint8]
	Distance  Optional[uint32]
	Altitude  Optional[float64]
	Velocity  Optional[uint16]
	Latitude  Optional[float64]
	Longitude Optional[float64]
}

func (a ActivityTimeseriesEntry) IsEmpty() bool {
	return !a.HeartRate.Valid &&
		!a.Cadence.Valid &&
		!a.Distance.Valid &&
		!a.Altitude.Valid &&
		!a.Velocity.Valid &&
		!a.Latitude.Valid &&
		!a.Longitude.Valid
}

func (a ActivityTimeseriesEntry) HasGPS() bool {
	return a.Latitude.Valid && a.Longitude.Valid
}

type AugmentConfig struct {
	ElevationHysteresisM float64 // Minimum cumulative elevation change (meters) to count as real gain/loss. Default 3.0.
}

func (c AugmentConfig) ApplyDefaults() AugmentConfig {
	config := c
	if config.ElevationHysteresisM == 0 {
		config.ElevationHysteresisM = 3.0
	}
	return config
}

// AugmentGPXData fills in missing Distance, Speed, Moving Time, and Ascent/Descent.
func AugmentGPXData(act *Activity, ts *ActivityTimeseries, config AugmentConfig) {
	config = config.ApplyDefaults()

	var totalDist float64
	var totalGain float64
	var totalLoss float64
	var runningDelta float64
	var movingSeconds uint32

	for i := 1; i < len(ts.Data); i++ {
		prev := ts.Data[i-1]
		curr := &ts.Data[i]

		if curr.HasGPS() && prev.HasGPS() {
			d := haversine(prev.Latitude.Value, prev.Longitude.Value, curr.Latitude.Value, curr.Longitude.Value)
			totalDist += d
			curr.Distance = Optional[uint32]{Value: uint32(totalDist), Valid: true}

			timeDelta := float64(curr.Offset - prev.Offset)
			if timeDelta > 0 {
				speed := d / timeDelta
				curr.Velocity = Optional[uint16]{Value: uint16(speed * 1000), Valid: true} // Storing as mm/s

				if speed > 0.5 {
					movingSeconds += uint32(timeDelta)
				}
			}
		}

		if curr.Altitude.Valid && prev.Altitude.Valid {
			deltaZ := curr.Altitude.Value - prev.Altitude.Value

			if (runningDelta > 0 && deltaZ < 0) || (runningDelta < 0 && deltaZ > 0) {
				runningDelta = deltaZ
			} else {
				runningDelta += deltaZ
			}

			if runningDelta >= config.ElevationHysteresisM {
				totalGain += runningDelta
				runningDelta = 0
			} else if runningDelta <= -config.ElevationHysteresisM {
				totalLoss -= runningDelta
				runningDelta = 0
			}
		}
	}

	act.Distance = uint32(totalDist)
	act.MovingTime = movingSeconds
	act.ElevationGain = Optional[uint16]{Value: uint16(totalGain), Valid: true}
	act.ElevationLoss = Optional[uint16]{Value: uint16(totalLoss), Valid: true}

	if movingSeconds > 0 {
		act.AvgSpeed = uint16((totalDist / float64(movingSeconds)) * 1000) // mm/s
	}
}
