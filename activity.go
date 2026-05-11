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

// WeightedAvg correctly computes time-weighted averages and standard deviation for uneven timeseries data.
type WeightedAvg struct {
	Sum    float64
	SumSq  float64 // For Standard Deviation
	Weight float64
	Count  int // Number of actual data points included
}

func (w *WeightedAvg) Add(val, weight float64) {
	w.Sum += val * weight
	w.SumSq += (val * val) * weight
	w.Weight += weight
	w.Count++
}

func (w *WeightedAvg) Avg() float64 {
	if w.Weight == 0 {
		return 0
	}
	return w.Sum / w.Weight
}

func (w *WeightedAvg) StdDev() float64 {
	if w.Weight == 0 {
		return 0
	}
	mean := w.Avg()
	variance := (w.SumSq / w.Weight) - (mean * mean)
	if variance < 0 {
		return 0 // prevent NaN from floating point precision issues
	}
	return math.Sqrt(variance)
}

type EnrichedPoint struct {
	Entry       *ActivityTimeseriesEntry
	TimeDelta   float64
	ActualSpeed float64
	GAPSpeed    float64
	GradePct    float64
	DistanceM   float64
	DeltaElevM  float64
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

// AugmentGPXData properly handles elevation hysteresis and initializes correctly
func AugmentGPXData(act *Activity, ts *ActivityTimeseries, config AugmentConfig) {
	config = config.ApplyDefaults()

	var totalDist float64
	var totalGain float64
	var totalLoss float64
	var movingSeconds uint32

	if len(ts.Data) == 0 {
		return
	}
	ts.Data[0].Distance = Optional[uint32]{Value: 0, Valid: true}

	var lastRecordedElev float64
	var elevInitialized bool

	for _, entry := range ts.Data {
		if entry.Altitude.Valid {
			lastRecordedElev = entry.Altitude.Value
			elevInitialized = true
			break
		}
	}

	for i := 1; i < len(ts.Data); i++ {
		prev := ts.Data[i-1]
		curr := &ts.Data[i]

		if curr.HasGPS() && prev.HasGPS() {
			d := haversine(prev.Latitude.Value, prev.Longitude.Value, curr.Latitude.Value, curr.Longitude.Value)
			totalDist += d
			timeDelta := float64(curr.Offset - prev.Offset)
			if timeDelta > 0 {
				speed := d / timeDelta
				curr.Velocity = Optional[uint16]{Value: uint16(speed * 1000), Valid: true}
				if speed > 0.5 { // Moving threshold
					movingSeconds += uint32(timeDelta)
				}
			}
		}
		curr.Distance = Optional[uint32]{Value: uint32(totalDist), Valid: true}

		if curr.Altitude.Valid && elevInitialized {
			deltaZ := curr.Altitude.Value - lastRecordedElev
			if deltaZ >= config.ElevationHysteresisM {
				totalGain += deltaZ
				lastRecordedElev = curr.Altitude.Value
			} else if deltaZ <= -config.ElevationHysteresisM {
				totalLoss -= deltaZ // deltaZ is negative, so subtract to add positive loss
				lastRecordedElev = curr.Altitude.Value
			}
		}
	}

	act.Distance = uint32(totalDist)
	act.MovingTime = movingSeconds
	act.ElevationGain = Optional[uint16]{Value: uint16(totalGain), Valid: true}
	act.ElevationLoss = Optional[uint16]{Value: uint16(totalLoss), Valid: true}
	if movingSeconds > 0 {
		act.AvgSpeed = uint16((totalDist / float64(movingSeconds)) * 1000)
	}
}

// DetectTopographicSplits uses a high/low watermark state machine to correctly isolate hills
func DetectTopographicSplits(enriched []EnrichedPoint, config LLMSummaryConfig) []TopographicSplit {
	if len(enriched) == 0 {
		return nil
	}

	var splits []TopographicSplit
	currentPhase := "Flat"

	startIdx := 0
	minIdx, maxIdx := 0, 0
	extremaIdx := 0

	minElev := enriched[0].Entry.Altitude.Value
	maxElev := enriched[0].Entry.Altitude.Value
	extremaElev := enriched[0].Entry.Altitude.Value

	addSplit := func(phase string, sIdx, eIdx int) {
		if eIdx <= sIdx {
			return
		}

		distKm := (enriched[eIdx].DistanceM - enriched[sIdx].DistanceM) / 1000.0
		startElev := enriched[sIdx].Entry.Altitude.Value
		endElev := enriched[eIdx].Entry.Altitude.Value

		gradePct := 0.0
		if distKm > 0 {
			gradePct = math.Round(((endElev-startElev)/(distKm*1000.0))*100*10) / 10
		}

		var gapAvg, hrAvg WeightedAvg

		for i := sIdx; i <= eIdx; i++ {
			pt := enriched[i]
			if pt.ActualSpeed > config.MinMovingSpeedMS {
				gapAvg.Add(pt.GAPSpeed, pt.TimeDelta)
				if pt.Entry.HeartRate.Valid {
					hrAvg.Add(float64(pt.Entry.HeartRate.Value), pt.TimeDelta)
				}
			}
		}

		splits = append(splits, TopographicSplit{
			Type:     phase,
			DistKm:   distKm,
			GradePct: gradePct,
			AvgGAP:   formatPace(gapAvg.Avg()),
			AvgHR:    int(hrAvg.Avg()),
		})
	}

	for i := 1; i < len(enriched); i++ {
		pt := enriched[i]
		if !pt.Entry.Altitude.Valid {
			continue
		}
		elev := pt.Entry.Altitude.Value

		switch currentPhase {
		case "Flat":
			if elev > maxElev {
				maxElev = elev
				maxIdx = i
			}
			if elev < minElev {
				minElev = elev
				minIdx = i
			}

			if elev-minElev > config.PhaseThresholdM {
				addSplit("Flat", startIdx, minIdx)
				currentPhase = "Uphill"
				startIdx = minIdx
				extremaElev = elev
				extremaIdx = i
			} else if maxElev-elev > config.PhaseThresholdM {
				addSplit("Flat", startIdx, maxIdx)
				currentPhase = "Downhill"
				startIdx = maxIdx
				extremaElev = elev
				extremaIdx = i
			}
		case "Uphill":
			if elev > extremaElev {
				extremaElev = elev
				extremaIdx = i
			} else if extremaElev-elev > config.PhaseThresholdM {
				addSplit("Uphill", startIdx, extremaIdx)
				currentPhase = "Downhill"
				startIdx = extremaIdx
				extremaElev = elev
				extremaIdx = i
			}
		case "Downhill":
			if elev < extremaElev {
				extremaElev = elev
				extremaIdx = i
			} else if elev-extremaElev > config.PhaseThresholdM {
				addSplit("Downhill", startIdx, extremaIdx)
				currentPhase = "Uphill"
				startIdx = extremaIdx
				extremaElev = elev
				extremaIdx = i
			}
		}
	}

	addSplit(currentPhase, startIdx, len(enriched)-1)

	return splits
}

// Helper to look back using actual timestamps, not array indices
func getHRAtOffset(ts *ActivityTimeseries, targetOffset int) int {
	bestDiff := math.MaxInt32
	bestHR := 0

	for _, entry := range ts.Data {
		if !entry.HeartRate.Valid {
			continue
		}
		diff := entry.Offset - targetOffset
		if diff < 0 {
			diff = -diff
		}
		if diff < bestDiff {
			bestDiff = diff
			bestHR = int(entry.HeartRate.Value)
		}
	}

	return bestHR
}

// Helper to get average HR within a specific time window (used for Post-Max Effort Recovery)
func getAvgHRInOffsetRange(ts *ActivityTimeseries, startOffset, endOffset int) int {
	var sum, count int
	for _, entry := range ts.Data {
		if entry.HeartRate.Valid && entry.Offset >= startOffset && entry.Offset <= endOffset {
			sum += int(entry.HeartRate.Value)
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / count
}
