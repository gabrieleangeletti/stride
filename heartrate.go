package stride

import (
	"errors"
	"fmt"
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

// HRThresholdAnalysisConfig defines the configuration for threshold analysis
type HRThresholdAnalysisConfig struct {
	LT1 uint8 // Aerobic threshold heart rate
	LT2 uint8 // Lactate threshold heart rate

	// Configurable parameters with sensible defaults
	BucketSizeSeconds          int     // Size of time buckets (default: 10)
	MinValidPointsPerBucket    int     // Minimum valid HR points needed per bucket (default: 5)
	ThresholdTolerancePercent  float64 // Tolerance around threshold (default: 0.05 for ±5%)
	LT1OverlapTolerancePercent float64 // If within LT1 + this %, choose LT1 over LT2 (default: 0.05)
	MinConsecutiveBuckets      int     // Minimum buckets for a valid sustained period (default: 5)
	ConsecutivePeriodThreshold float64 // Min percentage of buckets that must be LT1/LT2 in sequence (default: 0.90)
}

// ThresholdPeriod represents a single sustained period at a threshold
type ThresholdPeriod struct {
	Zone      thresholdType `json:"zone"`      // LT1 or LT2
	StartTime time.Time     `json:"startTime"` // Start time of the period
	EndTime   time.Time     `json:"endTime"`   // End time of the period
	Duration  int           `json:"duration"`  // Duration in seconds
	AvgHR     float64       `json:"avgHR"`     // Average heart rate during period
	MinHR     float64       `json:"minHR"`     // Minimum heart rate during period
	MaxHR     float64       `json:"maxHR"`     // Maximum heart rate during period
}

// HRThresholdAnalysisResult contains the analysis results
type HRThresholdAnalysisResult struct {
	TimeAtLT1Seconds int `json:"timeAtLT1Seconds"` // Total time spent at aerobic threshold
	TimeAtLT2Seconds int `json:"timeAtLT2Seconds"` // Total time spent at lactate threshold
	LT1Periods       int `json:"lt1Periods"`       // Number of sustained LT1 periods found
	LT2Periods       int `json:"lt2Periods"`       // Number of sustained LT2 periods found

	// Detailed period information
	Periods []ThresholdPeriod `json:"periods"` // All individual threshold periods

	// Workout quality metrics for LT1
	LT1AveragePeriodDuration int     `json:"lt1AveragePeriodDuration"` // Average duration of LT1 periods (seconds)
	LT1LongestPeriodDuration int     `json:"lt1LongestPeriodDuration"` // Longest LT1 period (seconds)
	LT1SmoothnessScore       float64 `json:"lt1SmoothnessScore"`       // 0-1, where 1 = one perfect sustained period

	// Workout quality metrics for LT2
	LT2AveragePeriodDuration int     `json:"lt2AveragePeriodDuration"` // Average duration of LT2 periods (seconds)
	LT2LongestPeriodDuration int     `json:"lt2LongestPeriodDuration"` // Longest LT2 period (seconds)
	LT2SmoothnessScore       float64 `json:"lt2SmoothnessScore"`       // 0-1, where 1 = one perfect sustained period
}

// thresholdType represents which threshold zone a bucket belongs to
type thresholdType int

const (
	ThresholdNone thresholdType = iota
	ThresholdLT1
	ThresholdLT2
)

// bucket represents a time bucket with its classification
type bucket struct {
	avgHeartRate float64
	threshold    thresholdType
	valid        bool
	index        int // Bucket index for timestamp calculation
}

// ApplyDefaults applies default values to unset configuration parameters
func (c *HRThresholdAnalysisConfig) ApplyDefaults() HRThresholdAnalysisConfig {
	config := *c

	if config.BucketSizeSeconds == 0 {
		config.BucketSizeSeconds = 10
	}

	if config.MinValidPointsPerBucket == 0 {
		config.MinValidPointsPerBucket = 5
	}

	if config.ThresholdTolerancePercent == 0 {
		config.ThresholdTolerancePercent = 0.05
	}

	if config.LT1OverlapTolerancePercent == 0 {
		config.LT1OverlapTolerancePercent = 0.05
	}

	if config.MinConsecutiveBuckets == 0 {
		config.MinConsecutiveBuckets = 5
	}

	if config.ConsecutivePeriodThreshold == 0 {
		config.ConsecutivePeriodThreshold = 0.90
	}

	return config
}

// AnalyzeHeartRateThresholds analyzes time spent at aerobic and lactate thresholds
func AnalyzeHeartRateThresholds(timeseries *ActivityTimeseries, config HRThresholdAnalysisConfig) (HRThresholdAnalysisResult, error) {
	if timeseries == nil {
		return HRThresholdAnalysisResult{}, fmt.Errorf("timeseries cannot be nil")
	}

	if len(timeseries.Data) == 0 {
		return HRThresholdAnalysisResult{}, fmt.Errorf("timeseries data is empty")
	}

	if config.LT1 == 0 || config.LT2 == 0 {
		return HRThresholdAnalysisResult{}, fmt.Errorf("LT1 and LT2 must be specified")
	}

	if config.LT2 <= config.LT1 {
		return HRThresholdAnalysisResult{}, fmt.Errorf("LT2 must be greater than LT1")
	}

	config = config.ApplyDefaults()

	buckets := createBuckets(timeseries, config)
	classifyBuckets(buckets, config)
	result := analyzePeriods(buckets, config, timeseries.StartTime)
	calculateSummaryMetrics(&result)

	return result, nil
}

// createBuckets divides the timeseries into time buckets and computes average HR
func createBuckets(timeseries *ActivityTimeseries, config HRThresholdAnalysisConfig) []bucket {
	if len(timeseries.Data) == 0 {
		return nil
	}

	numBuckets := (timeseries.MaxOffset() / config.BucketSizeSeconds) + 1
	buckets := make([]bucket, numBuckets)

	for i := range buckets {
		buckets[i].index = i
	}

	for _, entry := range timeseries.Data {
		if !entry.HeartRate.Valid {
			continue
		}

		bucketIdx := entry.Offset / config.BucketSizeSeconds
		if bucketIdx >= numBuckets {
			continue
		}

		buckets[bucketIdx].avgHeartRate += float64(entry.HeartRate.Value)
		buckets[bucketIdx].valid = true
	}

	bucketCounts := make([]int, numBuckets)
	for _, entry := range timeseries.Data {
		if !entry.HeartRate.Valid {
			continue
		}

		bucketIdx := entry.Offset / config.BucketSizeSeconds
		if bucketIdx < numBuckets {
			bucketCounts[bucketIdx]++
		}
	}

	for i := range buckets {
		if bucketCounts[i] >= config.MinValidPointsPerBucket {
			buckets[i].avgHeartRate = buckets[i].avgHeartRate / float64(bucketCounts[i])
			buckets[i].valid = true
		} else {
			buckets[i].valid = false
		}
	}

	return buckets
}

// classifyBuckets classifies each bucket as LT1, LT2, or neither
func classifyBuckets(buckets []bucket, config HRThresholdAnalysisConfig) {
	lt1Lower := float64(config.LT1) * (1 - config.ThresholdTolerancePercent)
	lt1Upper := float64(config.LT1) * (1 + config.ThresholdTolerancePercent)
	lt2Lower := float64(config.LT2) * (1 - config.ThresholdTolerancePercent)
	lt2Upper := float64(config.LT2) * (1 + config.ThresholdTolerancePercent)
	lt1OverlapUpper := float64(config.LT1) * (1 + config.LT1OverlapTolerancePercent)

	for i := range buckets {
		if !buckets[i].valid {
			buckets[i].threshold = ThresholdNone
			continue
		}

		hr := buckets[i].avgHeartRate
		inLT1 := hr >= lt1Lower && hr <= lt1Upper
		inLT2 := hr >= lt2Lower && hr <= lt2Upper

		if inLT1 && inLT2 {
			// Handle overlap: if within LT1 + overlap tolerance, choose LT1
			if hr <= lt1OverlapUpper {
				buckets[i].threshold = ThresholdLT1
			} else {
				buckets[i].threshold = ThresholdLT2
			}
		} else if inLT1 {
			buckets[i].threshold = ThresholdLT1
		} else if inLT2 {
			buckets[i].threshold = ThresholdLT2
		} else {
			buckets[i].threshold = ThresholdNone
		}
	}
}

// analyzePeriods finds sustained periods and counts time at each threshold
func analyzePeriods(buckets []bucket, config HRThresholdAnalysisConfig, startTime time.Time) HRThresholdAnalysisResult {
	result := HRThresholdAnalysisResult{
		Periods: make([]ThresholdPeriod, 0),
	}

	if len(buckets) == 0 {
		return result
	}

	// Find sequences allowing for gaps (thresholdNone) in between
	i := 0
	for i < len(buckets) {
		for i < len(buckets) && buckets[i].threshold == ThresholdNone {
			i++
		}
		if i >= len(buckets) {
			break
		}

		// Find end of sequence - continue through gaps until we hit
		// MinConsecutiveBuckets of consecutive ThresholdNone
		start := i
		consecutiveNone := 0
		lastThresholdIdx := i

		for i < len(buckets) {
			if buckets[i].threshold == ThresholdNone {
				consecutiveNone++
				if consecutiveNone >= config.MinConsecutiveBuckets {
					// Hit a significant gap, end the sequence at last threshold bucket
					break
				}
			} else {
				consecutiveNone = 0
				lastThresholdIdx = i
			}
			i++
		}

		// End is after the last threshold bucket (before the significant gap)
		end := lastThresholdIdx + 1

		// Analyze this sequence
		if end-start >= config.MinConsecutiveBuckets {
			analyzeSinglePeriod(buckets[start:end], config, &result, startTime)
		}

		// Continue from current position (after the gap if we hit one)
		if consecutiveNone >= config.MinConsecutiveBuckets {
			i = lastThresholdIdx + 1 + config.MinConsecutiveBuckets
		}
	}

	return result
}

// analyzeSinglePeriod analyzes a single consecutive sequence of buckets
func analyzeSinglePeriod(sequence []bucket, config HRThresholdAnalysisConfig, result *HRThresholdAnalysisResult, startTime time.Time) {
	lt1Count := 0
	lt2Count := 0
	for _, b := range sequence {
		switch b.threshold {
		case ThresholdLT1:
			lt1Count++

		case ThresholdLT2:
			lt2Count++
		}
	}

	totalThresholdBuckets := lt1Count + lt2Count
	if totalThresholdBuckets == 0 {
		return
	}

	if float64(lt1Count)/float64(len(sequence)) >= config.ConsecutivePeriodThreshold {
		addPeriod(sequence, ThresholdLT1, config, result, startTime)
		return
	}

	if float64(lt2Count)/float64(len(sequence)) >= config.ConsecutivePeriodThreshold {
		addPeriod(sequence, ThresholdLT2, config, result, startTime)
		return
	}

	analyzeSubPeriods(sequence, config, result, startTime)
}

// analyzeSubPeriods handles mixed sequences by finding homogeneous sub-sequences
func analyzeSubPeriods(sequence []bucket, config HRThresholdAnalysisConfig, result *HRThresholdAnalysisResult, startTime time.Time) {
	i := 0
	for i < len(sequence) {
		if sequence[i].threshold == ThresholdNone {
			i++
			continue
		}

		currentType := sequence[i].threshold
		start := i
		for i < len(sequence) && (sequence[i].threshold == currentType || sequence[i].threshold == ThresholdNone) {
			i++
		}

		subSeq := sequence[start:i]

		count := 0
		for _, b := range subSeq {
			if b.threshold == currentType {
				count++
			}
		}

		if len(subSeq) >= config.MinConsecutiveBuckets &&
			float64(count)/float64(len(subSeq)) >= config.ConsecutivePeriodThreshold {
			addPeriod(subSeq, currentType, config, result, startTime)
		}
	}
}

// addPeriod adds a threshold period to the result with full statistics
func addPeriod(sequence []bucket, zoneType thresholdType, config HRThresholdAnalysisConfig, result *HRThresholdAnalysisResult, startTime time.Time) {
	if len(sequence) == 0 {
		return
	}

	// Calculate period statistics
	var sumHR, minHR, maxHR float64
	count := 0
	minHR = 300 // Initialize to high value

	for _, b := range sequence {
		if b.valid && b.avgHeartRate > 0 {
			sumHR += b.avgHeartRate
			count++

			if b.avgHeartRate < minHR {
				minHR = b.avgHeartRate
			}

			if b.avgHeartRate > maxHR {
				maxHR = b.avgHeartRate
			}
		}
	}

	if count == 0 {
		return
	}

	avgHR := sumHR / float64(count)
	duration := len(sequence) * config.BucketSizeSeconds

	// Calculate start and end times based on bucket indices
	periodStartTime := startTime.Add(time.Duration(sequence[0].index*config.BucketSizeSeconds) * time.Second)
	periodEndTime := startTime.Add(time.Duration((sequence[len(sequence)-1].index+1)*config.BucketSizeSeconds) * time.Second)

	period := ThresholdPeriod{
		Zone:      zoneType,
		StartTime: periodStartTime,
		EndTime:   periodEndTime,
		Duration:  duration,
		AvgHR:     avgHR,
		MinHR:     minHR,
		MaxHR:     maxHR,
	}

	result.Periods = append(result.Periods, period)

	switch zoneType {
	case ThresholdLT1:
		result.TimeAtLT1Seconds += duration
		result.LT1Periods++
	case ThresholdLT2:
		result.TimeAtLT2Seconds += duration
		result.LT2Periods++
	}
}

// calculateSummaryMetrics computes average duration, longest period, and smoothness scores
func calculateSummaryMetrics(result *HRThresholdAnalysisResult) {
	var lt1Durations, lt2Durations []int

	for _, period := range result.Periods {
		switch period.Zone {
		case ThresholdLT1:
			lt1Durations = append(lt1Durations, period.Duration)
		case ThresholdLT2:
			lt2Durations = append(lt2Durations, period.Duration)
		}
	}

	// Calculate LT1 metrics
	if len(lt1Durations) > 0 {
		sum := 0
		longest := 0
		for _, d := range lt1Durations {
			sum += d
			if d > longest {
				longest = d
			}
		}
		result.LT1AveragePeriodDuration = sum / len(lt1Durations)
		result.LT1LongestPeriodDuration = longest

		// Smoothness score: ratio of longest period to total time
		// Further adjusted by penalizing multiple periods
		if result.TimeAtLT1Seconds > 0 {
			longestRatio := float64(longest) / float64(result.TimeAtLT1Seconds)
			periodPenalty := 1.0 / float64(len(lt1Durations))
			result.LT1SmoothnessScore = longestRatio * (0.5 + 0.5*periodPenalty)
		}
	}

	// Calculate LT2 metrics
	if len(lt2Durations) > 0 {
		sum := 0
		longest := 0
		for _, d := range lt2Durations {
			sum += d
			if d > longest {
				longest = d
			}
		}
		result.LT2AveragePeriodDuration = sum / len(lt2Durations)
		result.LT2LongestPeriodDuration = longest

		// Smoothness score: ratio of longest period to total time
		// Further adjusted by penalizing multiple periods
		if result.TimeAtLT2Seconds > 0 {
			longestRatio := float64(longest) / float64(result.TimeAtLT2Seconds)
			periodPenalty := 1.0 / float64(len(lt2Durations))
			result.LT2SmoothnessScore = longestRatio * (0.5 + 0.5*periodPenalty)
		}
	}
}
