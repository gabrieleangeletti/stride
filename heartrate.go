package stride

import (
	"errors"
	"time"
)

type AvgHeartRateMethod int

const (
	HeartRateMethodSimple AvgHeartRateMethod = iota
	HeartRateMethodTimeWeighted
	HeartRateMethodZoneWeighted
)

type MaxHeartRateMethod int

const (
	MaxHeartRateMethodPeak MaxHeartRateMethod = iota
	MaxHeartRateMethodRollingWindow
	MaxHeartRateMethodPercentile
)

type AvgHeartRateAnalysisConfig struct {
	Method           AvgHeartRateMethod
	ExcludeZeros     bool
	MinValidRate     int
	MaxValidRate     int
	MaxHeartRate     int
	SamplingInterval *time.Duration
}

type MaxHeartRateAnalysisConfig struct {
	Method          MaxHeartRateMethod
	ExcludeZeros    bool
	MinValidRate    int
	MaxValidRate    int
	WindowDuration  *time.Duration
	PercentileValue float64
}

var (
	ErrEmptyTimeseriesData = errors.New("activity timeseries data is empty")
	ErrNoValidData         = errors.New("no valid heart rate data points found")
)

func CalculateAverageHeartRate(timeseries *ActivityTimeseries, config AvgHeartRateAnalysisConfig) (float64, error) {
	if len(timeseries.Data) == 0 {
		return 0, ErrEmptyTimeseriesData
	}

	switch config.Method {
	case HeartRateMethodSimple:
		return calculateSimpleAverageFromTimeseries(timeseries.Data, config)

	case HeartRateMethodTimeWeighted:
		return calculateTimeWeightedAverageFromTimeseries(timeseries.Data, config)

	case HeartRateMethodZoneWeighted:
		return calculateZoneWeightedAverageFromTimeseries(timeseries.Data, config)

	default:
		return calculateTimeWeightedAverageFromTimeseries(timeseries.Data, config)
	}
}

func CalculateMaxHeartRate(timeseries *ActivityTimeseries, config MaxHeartRateAnalysisConfig) (int, error) {
	if len(timeseries.Data) == 0 {
		return 0, ErrEmptyTimeseriesData
	}

	switch config.Method {
	case MaxHeartRateMethodPeak:
		return calculatePeakHeartRate(timeseries.Data, config)

	case MaxHeartRateMethodRollingWindow:
		return calculateRollingWindowMaxHeartRate(timeseries.Data, config)

	case MaxHeartRateMethodPercentile:
		return calculatePercentileMaxHeartRate(timeseries.Data, config)

	default:
		return calculatePeakHeartRate(timeseries.Data, config)
	}
}

func isValidMaxHeartRate(hr int, config MaxHeartRateAnalysisConfig) bool {
	if config.ExcludeZeros && hr == 0 {
		return false
	}

	if config.MinValidRate > 0 && hr < config.MinValidRate {
		return false
	}

	if config.MaxValidRate > 0 && hr > config.MaxValidRate {
		return false
	}

	return true
}

func calculatePeakHeartRate(data []ActivityTimeseriesEntry, config MaxHeartRateAnalysisConfig) (int, error) {
	maxHR := 0
	foundValid := false

	for _, entry := range data {
		if !entry.HeartRate.Valid {
			continue
		}

		hr := int(entry.HeartRate.Value)
		if !isValidMaxHeartRate(hr, config) {
			continue
		}

		if hr > maxHR {
			maxHR = hr
			foundValid = true
		}
	}

	if !foundValid {
		return 0, ErrNoValidData
	}

	return maxHR, nil
}

func calculateRollingWindowMaxHeartRate(data []ActivityTimeseriesEntry, config MaxHeartRateAnalysisConfig) (int, error) {
	windowDurationSeconds := 30
	if config.WindowDuration != nil {
		windowDurationSeconds = int(config.WindowDuration.Seconds())
	}

	if len(data) < 2 {
		return calculatePeakHeartRate(data, config)
	}

	overallMax := 0
	foundValid := false

	for i := range len(data) {
		windowStart := data[i].Offset
		windowEnd := windowStart + windowDurationSeconds
		windowMax := 0

		for j := i; j < len(data) && data[j].Offset < windowEnd; j++ {
			if !data[j].HeartRate.Valid {
				continue
			}

			hr := int(data[j].HeartRate.Value)
			if !isValidMaxHeartRate(hr, config) {
				continue
			}

			if hr > windowMax {
				windowMax = hr
			}
		}

		if windowMax > overallMax {
			overallMax = windowMax
			foundValid = true
		}
	}

	if !foundValid {
		return 0, ErrNoValidData
	}

	return overallMax, nil
}

func calculatePercentileMaxHeartRate(data []ActivityTimeseriesEntry, config MaxHeartRateAnalysisConfig) (int, error) {
	var validHeartRates []int

	for _, entry := range data {
		if !entry.HeartRate.Valid {
			continue
		}

		hr := int(entry.HeartRate.Value)
		if !isValidMaxHeartRate(hr, config) {
			continue
		}

		validHeartRates = append(validHeartRates, hr)
	}

	if len(validHeartRates) == 0 {
		return 0, ErrNoValidData
	}

	percentile := config.PercentileValue
	if percentile <= 0 || percentile > 100 {
		percentile = 95.0
	}

	sortHeartRates(validHeartRates)

	index := float64(len(validHeartRates)-1) * (percentile / 100.0)
	lowerIndex := int(index)
	upperIndex := lowerIndex + 1

	if upperIndex >= len(validHeartRates) {
		return validHeartRates[len(validHeartRates)-1], nil
	}

	weight := index - float64(lowerIndex)
	result := float64(validHeartRates[lowerIndex])*(1-weight) + float64(validHeartRates[upperIndex])*weight

	return int(result + 0.5), nil
}

func sortHeartRates(heartRates []int) {
	for i := 0; i < len(heartRates)-1; i++ {
		for j := i + 1; j < len(heartRates); j++ {
			if heartRates[i] > heartRates[j] {
				heartRates[i], heartRates[j] = heartRates[j], heartRates[i]
			}
		}
	}
}

func isValidHeartRate(hr int, config AvgHeartRateAnalysisConfig) bool {
	if config.ExcludeZeros && hr == 0 {
		return false
	}

	if config.MinValidRate > 0 && hr < config.MinValidRate {
		return false
	}

	if config.MaxValidRate > 0 && hr > config.MaxValidRate {
		return false
	}

	return true
}

func calculateSimpleAverageFromTimeseries(data []ActivityTimeseriesEntry, config AvgHeartRateAnalysisConfig) (float64, error) {
	var sum, count int

	for _, entry := range data {
		if !entry.HeartRate.Valid {
			continue
		}

		hr := int(entry.HeartRate.Value)
		if isValidHeartRate(hr, config) {
			sum += hr
			count++
		}
	}

	if count == 0 {
		return 0, ErrNoValidData
	}

	return float64(sum) / float64(count), nil
}

func calculateTimeWeightedAverageFromTimeseries(data []ActivityTimeseriesEntry, config AvgHeartRateAnalysisConfig) (float64, error) {
	if len(data) < 2 {
		return calculateSimpleAverageFromTimeseries(data, config)
	}

	var weightedSum float64
	var totalDuration int

	for i := 0; i < len(data)-1; i++ {
		if !data[i].HeartRate.Valid {
			continue
		}

		hr := int(data[i].HeartRate.Value)
		if !isValidHeartRate(hr, config) {
			continue
		}

		duration := data[i+1].Offset - data[i].Offset
		if duration <= 0 {
			continue
		}

		weightedSum += float64(hr) * float64(duration)
		totalDuration += duration
	}

	if totalDuration == 0 {
		return 0, ErrNoValidData
	}

	return weightedSum / float64(totalDuration), nil
}

func calculateZoneWeightedAverageFromTimeseries(data []ActivityTimeseriesEntry, config AvgHeartRateAnalysisConfig) (float64, error) {
	if len(data) < 2 {
		return calculateSimpleAverageFromTimeseries(data, config)
	}

	var weightedSum float64
	var totalWeightedDuration float64

	for i := 0; i < len(data)-1; i++ {
		if !data[i].HeartRate.Valid {
			continue
		}

		hr := int(data[i].HeartRate.Value)
		if !isValidHeartRate(hr, config) {
			continue
		}

		duration := data[i+1].Offset - data[i].Offset
		if duration <= 0 {
			continue
		}

		zoneWeight := getHeartRateZoneWeight(hr, config.MaxHeartRate)
		weightedSum += float64(hr) * float64(duration) * zoneWeight
		totalWeightedDuration += float64(duration) * zoneWeight
	}

	if totalWeightedDuration == 0 {
		return 0, ErrNoValidData
	}

	return weightedSum / totalWeightedDuration, nil
}

func getHeartRateZoneWeight(hr int, maxHeartRate int) float64 {
	if maxHeartRate <= 0 {
		maxHeartRate = 220
	}

	hrPercent := float64(hr) / float64(maxHeartRate)

	switch {
	case hrPercent <= 0.5:
		return 1.0
	case hrPercent <= 0.6:
		return 1.0
	case hrPercent <= 0.7:
		return 2.0
	case hrPercent <= 0.8:
		return 3.0
	case hrPercent <= 0.9:
		return 4.0
	default:
		return 5.0
	}
}
