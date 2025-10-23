package stride_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/gabrieleangeletti/stride"
)

func TestCalculateAverageHeartRate(t *testing.T) {
	timeseries := &ActivityTimeseries{
		StartTime: time.Now(),
		Data: []ActivityTimeseriesEntry{
			{Offset: 0, HeartRate: Optional[uint8]{Value: 120, Valid: true}},
			{Offset: 10, HeartRate: Optional[uint8]{Value: 130, Valid: true}},
			{Offset: 20, HeartRate: Optional[uint8]{Value: 125, Valid: true}},
			{Offset: 30, HeartRate: Optional[uint8]{Value: 140, Valid: true}},
			{Offset: 40, HeartRate: Optional[uint8]{Value: 0, Valid: false}},
			{Offset: 50, HeartRate: Optional[uint8]{Value: 135, Valid: true}},
			{Offset: 60, HeartRate: Optional[uint8]{Value: 145, Valid: true}},
			{Offset: 70, HeartRate: Optional[uint8]{Value: 150, Valid: true}},
			{Offset: 80, HeartRate: Optional[uint8]{Value: 130, Valid: true}},
		},
	}

	config := AvgHeartRateAnalysisConfig{
		Method:       HeartRateMethodTimeWeighted,
		ExcludeZeros: true,
		MinValidRate: 40,
		MaxValidRate: 220,
		MaxHeartRate: 190,
	}

	t.Run("TimeWeighted", func(t *testing.T) {
		avgHR, err := CalculateAverageHeartRate(timeseries, config)
		require.NoError(t, err)

		// Time-weighted: Only uses intervals between consecutive entries
		// (120*10 + 130*10 + 125*10 + 140*10 + 135*10 + 145*10 + 150*10) / 70
		// = (1200 + 1300 + 1250 + 1400 + 1350 + 1450 + 1500) / 70 = 9450 / 70 = 135.0
		// Note: Last entry (130) has no next entry, so it's not included in time-weighted calc
		expected := 135.0
		assert.InDelta(t, expected, avgHR, 0.01, "Time-weighted average should be close to expected value")
	})

	t.Run("Simple", func(t *testing.T) {
		config.Method = HeartRateMethodSimple
		simpleAvgHR, err := CalculateAverageHeartRate(timeseries, config)
		require.NoError(t, err)

		// Simple average: (120+130+125+140+135+145+150+130) / 8 = 1075 / 8 = 134.375
		expected := 134.375
		assert.InDelta(t, expected, simpleAvgHR, 0.01, "Simple average should be close to expected value")
	})

	t.Run("ZoneWeighted", func(t *testing.T) {
		config.Method = HeartRateMethodZoneWeighted
		zoneAvgHR, err := CalculateAverageHeartRate(timeseries, config)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, zoneAvgHR, 120.0, "Zone-weighted average should be at least 120")
		assert.LessOrEqual(t, zoneAvgHR, 160.0, "Zone-weighted average should be at most 160")
	})
}

func TestCalculateMaxHeartRate(t *testing.T) {
	timeseries := &ActivityTimeseries{
		StartTime: time.Now(),
		Data: []ActivityTimeseriesEntry{
			{Offset: 0, HeartRate: Optional[uint8]{Value: 120, Valid: true}},
			{Offset: 10, HeartRate: Optional[uint8]{Value: 130, Valid: true}},
			{Offset: 20, HeartRate: Optional[uint8]{Value: 165, Valid: true}},
			{Offset: 30, HeartRate: Optional[uint8]{Value: 140, Valid: true}},
			{Offset: 40, HeartRate: Optional[uint8]{Value: 0, Valid: false}},
			{Offset: 50, HeartRate: Optional[uint8]{Value: 175, Valid: true}}, // Peak
			{Offset: 60, HeartRate: Optional[uint8]{Value: 145, Valid: true}},
			{Offset: 70, HeartRate: Optional[uint8]{Value: 150, Valid: true}},
			{Offset: 80, HeartRate: Optional[uint8]{Value: 130, Valid: true}},
		},
	}

	t.Run("Peak", func(t *testing.T) {
		config := MaxHeartRateAnalysisConfig{
			Method:       MaxHeartRateMethodPeak,
			ExcludeZeros: true,
			MinValidRate: 40,
			MaxValidRate: 220,
		}

		maxHR, err := CalculateMaxHeartRate(timeseries, config)
		require.NoError(t, err)
		assert.Equal(t, 175, maxHR, "Peak max HR should be 175")
	})

	t.Run("RollingWindow", func(t *testing.T) {
		windowDuration := 30 * time.Second
		config := MaxHeartRateAnalysisConfig{
			Method:         MaxHeartRateMethodRollingWindow,
			ExcludeZeros:   true,
			MinValidRate:   40,
			MaxValidRate:   220,
			WindowDuration: &windowDuration,
		}

		maxHR, err := CalculateMaxHeartRate(timeseries, config)
		require.NoError(t, err)
		assert.Equal(t, 175, maxHR, "Rolling window max HR should be 175")
	})

	t.Run("Percentile95", func(t *testing.T) {
		config := MaxHeartRateAnalysisConfig{
			Method:          MaxHeartRateMethodPercentile,
			ExcludeZeros:    true,
			MinValidRate:    40,
			MaxValidRate:    220,
			PercentileValue: 95.0,
		}

		maxHR, err := CalculateMaxHeartRate(timeseries, config)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, maxHR, 165, "95th percentile max HR should be at least 165")
		assert.LessOrEqual(t, maxHR, 175, "95th percentile max HR should be at most 175")
	})

	t.Run("Percentile90", func(t *testing.T) {
		config := MaxHeartRateAnalysisConfig{
			Method:          MaxHeartRateMethodPercentile,
			ExcludeZeros:    true,
			MinValidRate:    40,
			MaxValidRate:    220,
			PercentileValue: 90.0,
		}

		maxHR, err := CalculateMaxHeartRate(timeseries, config)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, maxHR, 150, "90th percentile max HR should be at least 150")
		assert.LessOrEqual(t, maxHR, 175, "90th percentile max HR should be at most 175")
	})
}

func TestCalculateAverageHeartRateEdgeCases(t *testing.T) {
	t.Run("EmptyData", func(t *testing.T) {
		timeseries := &ActivityTimeseries{
			StartTime: time.Now(),
			Data:      []ActivityTimeseriesEntry{},
		}
		config := AvgHeartRateAnalysisConfig{
			Method: HeartRateMethodSimple,
		}

		_, err := CalculateAverageHeartRate(timeseries, config)
		assert.ErrorIs(t, err, ErrEmptyTimeseriesData, "Should return ErrEmptyTimeseriesData for empty data")
	})

	t.Run("NoValidHeartRateData", func(t *testing.T) {
		timeseries := &ActivityTimeseries{
			StartTime: time.Now(),
			Data: []ActivityTimeseriesEntry{
				{Offset: 0, HeartRate: Optional[uint8]{Value: 0, Valid: false}},
				{Offset: 10, HeartRate: Optional[uint8]{Value: 0, Valid: false}},
			},
		}
		config := AvgHeartRateAnalysisConfig{
			Method:       HeartRateMethodSimple,
			ExcludeZeros: true,
		}

		_, err := CalculateAverageHeartRate(timeseries, config)
		assert.ErrorIs(t, err, ErrNoValidData, "Should return ErrNoValidData when no valid heart rate data")
	})

	t.Run("FilterByValidRate", func(t *testing.T) {
		timeseries := &ActivityTimeseries{
			StartTime: time.Now(),
			Data: []ActivityTimeseriesEntry{
				{Offset: 0, HeartRate: Optional[uint8]{Value: 30, Valid: true}},   // Too low
				{Offset: 10, HeartRate: Optional[uint8]{Value: 120, Valid: true}}, // Valid
				{Offset: 20, HeartRate: Optional[uint8]{Value: 250, Valid: true}}, // Too high
			},
		}
		config := AvgHeartRateAnalysisConfig{
			Method:       HeartRateMethodSimple,
			MinValidRate: 40,
			MaxValidRate: 220,
		}

		avgHR, err := CalculateAverageHeartRate(timeseries, config)
		require.NoError(t, err)
		assert.Equal(t, 120.0, avgHR, "Should only use the valid heart rate reading")
	})
}

func TestCalculateMaxHeartRateEdgeCases(t *testing.T) {
	t.Run("EmptyData", func(t *testing.T) {
		timeseries := &ActivityTimeseries{
			StartTime: time.Now(),
			Data:      []ActivityTimeseriesEntry{},
		}
		config := MaxHeartRateAnalysisConfig{
			Method: MaxHeartRateMethodPeak,
		}

		_, err := CalculateMaxHeartRate(timeseries, config)
		assert.ErrorIs(t, err, ErrEmptyTimeseriesData, "Should return ErrEmptyTimeseriesData for empty data")
	})

	t.Run("NoValidData", func(t *testing.T) {
		timeseries := &ActivityTimeseries{
			StartTime: time.Now(),
			Data: []ActivityTimeseriesEntry{
				{Offset: 0, HeartRate: Optional[uint8]{Value: 0, Valid: false}},
				{Offset: 10, HeartRate: Optional[uint8]{Value: 0, Valid: false}},
			},
		}
		config := MaxHeartRateAnalysisConfig{
			Method:       MaxHeartRateMethodPeak,
			ExcludeZeros: true,
		}

		_, err := CalculateMaxHeartRate(timeseries, config)
		assert.ErrorIs(t, err, ErrNoValidData, "Should return ErrNoValidData when no valid heart rate data")
	})

	t.Run("FilterByValidRate", func(t *testing.T) {
		timeseries := &ActivityTimeseries{
			StartTime: time.Now(),
			Data: []ActivityTimeseriesEntry{
				{Offset: 0, HeartRate: Optional[uint8]{Value: 30, Valid: true}},   // Too low
				{Offset: 10, HeartRate: Optional[uint8]{Value: 160, Valid: true}}, // Valid
				{Offset: 20, HeartRate: Optional[uint8]{Value: 250, Valid: true}}, // Too high
			},
		}
		config := MaxHeartRateAnalysisConfig{
			Method:       MaxHeartRateMethodPeak,
			MinValidRate: 40,
			MaxValidRate: 220,
		}

		maxHR, err := CalculateMaxHeartRate(timeseries, config)
		require.NoError(t, err)
		assert.Equal(t, 160, maxHR, "Should only use the valid heart rate reading")
	})
}

// Helper function to create a timeseries with heart rate data
func createTimeseries(startTime time.Time, heartRates []uint8) *ActivityTimeseries {
	data := make([]ActivityTimeseriesEntry, len(heartRates))
	for i, hr := range heartRates {
		data[i] = ActivityTimeseriesEntry{
			Offset:    i,
			HeartRate: Optional[uint8]{Value: hr, Valid: true},
		}
	}
	return &ActivityTimeseries{
		StartTime: startTime,
		Data:      data,
	}
}

// Helper function to create a timeseries with optional heart rate data
func createTimeseriesWithGaps(startTime time.Time, heartRates []Optional[uint8]) *ActivityTimeseries {
	data := make([]ActivityTimeseriesEntry, len(heartRates))
	for i, hr := range heartRates {
		data[i] = ActivityTimeseriesEntry{
			Offset:    i,
			HeartRate: hr,
		}
	}
	return &ActivityTimeseries{
		StartTime: startTime,
		Data:      data,
	}
}

// Test basic configuration defaults
func TestConfigDefaults(t *testing.T) {
	config := HRThresholdAnalysisConfig{
		LT1: 150,
		LT2: 170,
	}

	config = config.ApplyDefaults()

	assert.Equal(t, 10, config.BucketSizeSeconds)
	assert.Equal(t, 5, config.MinValidPointsPerBucket)
	assert.Equal(t, 0.05, config.ThresholdTolerancePercent)
	assert.Equal(t, 0.05, config.LT1OverlapTolerancePercent)
	assert.Equal(t, 5, config.MinConsecutiveBuckets)
	assert.Equal(t, 0.90, config.ConsecutivePeriodThreshold)
}

// Test error handling for invalid inputs
func TestInvalidInputs(t *testing.T) {
	config := HRThresholdAnalysisConfig{LT1: 150, LT2: 170}

	// Nil timeseries
	_, err := AnalyzeHeartRateThresholds(nil, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")

	// Empty timeseries
	emptyTS := &ActivityTimeseries{StartTime: time.Now(), Data: []ActivityTimeseriesEntry{}}
	_, err = AnalyzeHeartRateThresholds(emptyTS, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty")

	// Missing LT1
	_, err = AnalyzeHeartRateThresholds(createTimeseries(time.Now(), []uint8{150}), HRThresholdAnalysisConfig{LT2: 170})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be specified")

	// LT2 <= LT1
	_, err = AnalyzeHeartRateThresholds(createTimeseries(time.Now(), []uint8{150}), HRThresholdAnalysisConfig{LT1: 170, LT2: 150})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be greater")
}

// Test single sustained LT1 period
func TestSingleLT1Period(t *testing.T) {
	// Create 60 seconds of heart rate data at LT1 (150 bpm)
	heartRates := make([]uint8, 60)
	for i := range heartRates {
		heartRates[i] = 150
	}

	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	ts := createTimeseries(startTime, heartRates)

	config := HRThresholdAnalysisConfig{
		LT1: 150,
		LT2: 170,
	}

	result, err := AnalyzeHeartRateThresholds(ts, config)
	assert.NoError(t, err)

	// Should have 60 seconds at LT1 (6 buckets of 10s each)
	assert.Equal(t, 60, result.TimeAtLT1Seconds)
	assert.Equal(t, 0, result.TimeAtLT2Seconds)
	assert.Equal(t, 1, result.LT1Periods)
	assert.Equal(t, 0, result.LT2Periods)

	// Check period details
	assert.Len(t, result.Periods, 1)
	period := result.Periods[0]
	assert.Equal(t, ThresholdLT1, period.Zone)
	assert.Equal(t, 60, period.Duration)
	assert.InDelta(t, 150.0, period.AvgHR, 0.1)
	assert.Equal(t, startTime, period.StartTime)
	assert.Equal(t, startTime.Add(60*time.Second), period.EndTime)

	// Check quality metrics
	assert.Equal(t, 60, result.LT1AveragePeriodDuration)
	assert.Equal(t, 60, result.LT1LongestPeriodDuration)
	assert.InDelta(t, 1.0, result.LT1SmoothnessScore, 0.01)
}

// Test single sustained LT2 period
func TestSingleLT2Period(t *testing.T) {
	// Create 100 seconds of heart rate data at LT2 (170 bpm)
	heartRates := make([]uint8, 100)
	for i := range heartRates {
		heartRates[i] = 170
	}

	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	ts := createTimeseries(startTime, heartRates)

	config := HRThresholdAnalysisConfig{
		LT1: 150,
		LT2: 170,
	}

	result, err := AnalyzeHeartRateThresholds(ts, config)
	assert.NoError(t, err)

	assert.Equal(t, 0, result.TimeAtLT1Seconds)
	assert.Equal(t, 100, result.TimeAtLT2Seconds)
	assert.Equal(t, 0, result.LT1Periods)
	assert.Equal(t, 1, result.LT2Periods)

	// Check quality metrics
	assert.Equal(t, 100, result.LT2AveragePeriodDuration)
	assert.Equal(t, 100, result.LT2LongestPeriodDuration)
	assert.InDelta(t, 1.0, result.LT2SmoothnessScore, 0.01)
}

// Test threshold tolerance zones
func TestThresholdTolerance(t *testing.T) {
	config := HRThresholdAnalysisConfig{
		LT1:                       150,
		LT2:                       170,
		ThresholdTolerancePercent: 0.05, // ±5%
	}

	// LT1 range: 142.5 - 157.5
	// Just below LT1 tolerance
	ts1 := createTimeseries(time.Now(), make([]uint8, 60))
	for i := range ts1.Data {
		ts1.Data[i].HeartRate = Optional[uint8]{Value: 142, Valid: true}
	}
	result1, _ := AnalyzeHeartRateThresholds(ts1, config)
	assert.Equal(t, 0, result1.TimeAtLT1Seconds, "142 bpm should be below LT1 range")

	// Just inside LT1 tolerance
	ts2 := createTimeseries(time.Now(), make([]uint8, 60))
	for i := range ts2.Data {
		ts2.Data[i].HeartRate = Optional[uint8]{Value: 143, Valid: true}
	}
	result2, _ := AnalyzeHeartRateThresholds(ts2, config)
	assert.Equal(t, 60, result2.TimeAtLT1Seconds, "143 bpm should be inside LT1 range")

	// Upper bound of LT1
	ts3 := createTimeseries(time.Now(), make([]uint8, 60))
	for i := range ts3.Data {
		ts3.Data[i].HeartRate = Optional[uint8]{Value: 157, Valid: true}
	}
	result3, _ := AnalyzeHeartRateThresholds(ts3, config)
	assert.Equal(t, 60, result3.TimeAtLT1Seconds, "157 bpm should be inside LT1 range")
}

// Test LT1/LT2 overlap resolution
func TestOverlapResolution(t *testing.T) {
	config := HRThresholdAnalysisConfig{
		LT1:                        150,
		LT2:                        160,
		ThresholdTolerancePercent:  0.05,
		LT1OverlapTolerancePercent: 0.05,
	}

	// LT1 range: 142.5 - 157.5
	// LT2 range: 152 - 168
	// Overlap: 152 - 157.5
	// LT1 + 5% = 157.5, so HR <= 157.5 should be LT1, HR > 157.5 should be LT2

	// HR = 155 (in overlap, should be LT1)
	ts1 := createTimeseries(time.Now(), make([]uint8, 60))
	for i := range ts1.Data {
		ts1.Data[i].HeartRate = Optional[uint8]{Value: 155, Valid: true}
	}
	result1, _ := AnalyzeHeartRateThresholds(ts1, config)
	assert.Equal(t, 60, result1.TimeAtLT1Seconds)
	assert.Equal(t, 0, result1.TimeAtLT2Seconds)

	// HR = 158 (in overlap but > LT1+5%, should be LT2)
	ts2 := createTimeseries(time.Now(), make([]uint8, 60))
	for i := range ts2.Data {
		ts2.Data[i].HeartRate = Optional[uint8]{Value: 158, Valid: true}
	}
	result2, _ := AnalyzeHeartRateThresholds(ts2, config)
	assert.Equal(t, 0, result2.TimeAtLT1Seconds)
	assert.Equal(t, 60, result2.TimeAtLT2Seconds)
}

// Test minimum consecutive buckets requirement
func TestMinConsecutiveBuckets(t *testing.T) {
	config := HRThresholdAnalysisConfig{
		LT1:                   150,
		LT2:                   170,
		MinConsecutiveBuckets: 5, // 50 seconds minimum
	}

	// Only 40 seconds at LT1 (should not count)
	heartRates := make([]uint8, 100)
	for i := range 40 {
		heartRates[i] = 150 // LT1
	}
	for i := 40; i < 100; i++ {
		heartRates[i] = 120 // Below threshold
	}

	ts := createTimeseries(time.Now(), heartRates)
	result, err := AnalyzeHeartRateThresholds(ts, config)
	assert.NoError(t, err)
	assert.Equal(t, 0, result.TimeAtLT1Seconds, "40s period should be too short")

	// Exactly 50 seconds at LT1 (should count)
	for i := range 50 {
		heartRates[i] = 150
	}
	ts2 := createTimeseries(time.Now(), heartRates)
	result2, _ := AnalyzeHeartRateThresholds(ts2, config)
	assert.Equal(t, 50, result2.TimeAtLT1Seconds, "50s period should be long enough")
}

// Test consecutive period threshold (90% rule)
func TestConsecutivePeriodThreshold(t *testing.T) {
	config := HRThresholdAnalysisConfig{
		LT1:                        150,
		LT2:                        170,
		MinConsecutiveBuckets:      4,
		ConsecutivePeriodThreshold: 0.90, // 90% must be in zone
	}

	// 10 buckets (100s): 9 at LT1 in middle, 1 gap in between (90% - should count all 100s)
	heartRates := make([]uint8, 100)
	for i := range 50 {
		heartRates[i] = 150 // LT1
	}
	for i := 50; i < 60; i++ {
		heartRates[i] = 120 // Gap in the middle
	}
	for i := 60; i < 100; i++ {
		heartRates[i] = 150 // LT1
	}

	startTime := time.Now()
	ts := createTimeseries(startTime, heartRates)
	result, err := AnalyzeHeartRateThresholds(ts, config)
	assert.NoError(t, err)
	assert.Equal(t, 100, result.TimeAtLT1Seconds, "Should count all 100s including gap in middle")
	assert.Equal(t, 1, result.LT1Periods)

	// Test gap at the end - should NOT count
	heartRates2 := make([]uint8, 100)
	for i := range 90 {
		heartRates2[i] = 150 // LT1
	}
	for i := 90; i < 100; i++ {
		heartRates2[i] = 120 // Gap at end
	}

	ts2 := createTimeseries(startTime, heartRates2)
	result2, _ := AnalyzeHeartRateThresholds(ts2, config)
	assert.Equal(t, 90, result2.TimeAtLT1Seconds, "Should NOT count gap at end")

	// Test gap at the start - should NOT count
	heartRates3 := make([]uint8, 100)
	for i := range 10 {
		heartRates3[i] = 120 // Gap at start
	}
	for i := 10; i < 100; i++ {
		heartRates3[i] = 150 // LT1
	}

	ts3 := createTimeseries(startTime, heartRates3)
	result3, _ := AnalyzeHeartRateThresholds(ts3, config)
	assert.Equal(t, 90, result3.TimeAtLT1Seconds, "Should NOT count gap at start")
}

// Test multiple separate periods
func TestMultipleSeparatePeriods(t *testing.T) {
	config := HRThresholdAnalysisConfig{
		LT1:                   150,
		LT2:                   170,
		MinConsecutiveBuckets: 2,
	}

	// Pattern: 60s LT1, 30s rest, 80s LT1, 30s rest, 50s LT2
	heartRates := make([]uint8, 250)

	// First LT1 period
	for i := range 60 {
		heartRates[i] = 150
	}
	// Rest
	for i := 60; i < 90; i++ {
		heartRates[i] = 120
	}
	// Second LT1 period
	for i := 90; i < 170; i++ {
		heartRates[i] = 150
	}
	// Rest
	for i := 170; i < 200; i++ {
		heartRates[i] = 120
	}
	// LT2 period
	for i := 200; i < 250; i++ {
		heartRates[i] = 170
	}

	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	ts := createTimeseries(startTime, heartRates)
	result, err := AnalyzeHeartRateThresholds(ts, config)
	assert.NoError(t, err)

	assert.Equal(t, 140, result.TimeAtLT1Seconds) // 60 + 80
	assert.Equal(t, 50, result.TimeAtLT2Seconds)
	assert.Equal(t, 2, result.LT1Periods)
	assert.Equal(t, 1, result.LT2Periods)

	// Check periods
	assert.Len(t, result.Periods, 3)

	// Check quality metrics
	assert.Equal(t, 70, result.LT1AveragePeriodDuration) // (60+80)/2
	assert.Equal(t, 80, result.LT1LongestPeriodDuration)
	assert.Less(t, result.LT1SmoothnessScore, 1.0, "Multiple periods should lower smoothness")

	assert.Equal(t, 50, result.LT2AveragePeriodDuration)
	assert.Equal(t, 50, result.LT2LongestPeriodDuration)
}

// Test spiky vs smooth workout
func TestSmoothnesScore(t *testing.T) {
	config := HRThresholdAnalysisConfig{
		LT1: 150,
		LT2: 170,
	}

	// Smooth workout: one 300s period
	smoothHR := make([]uint8, 300)
	for i := range smoothHR {
		smoothHR[i] = 170
	}
	smoothTS := createTimeseries(time.Now(), smoothHR)
	smoothResult, _ := AnalyzeHeartRateThresholds(smoothTS, config)

	// Spiky workout: three 100s periods with rests in between
	spikyHR := make([]uint8, 500)
	for i := range 100 {
		spikyHR[i] = 170
	}
	for i := 100; i < 200; i++ {
		spikyHR[i] = 120 // Rest
	}
	for i := 200; i < 300; i++ {
		spikyHR[i] = 170
	}
	for i := 300; i < 400; i++ {
		spikyHR[i] = 120 // Rest
	}
	for i := 400; i < 500; i++ {
		spikyHR[i] = 170
	}
	spikyTS := createTimeseries(time.Now(), spikyHR)
	spikyResult, _ := AnalyzeHeartRateThresholds(spikyTS, config)

	// Smooth should have higher score
	assert.Greater(t, smoothResult.LT2SmoothnessScore, spikyResult.LT2SmoothnessScore)
	assert.InDelta(t, 1.0, smoothResult.LT2SmoothnessScore, 0.01)
	assert.Less(t, spikyResult.LT2SmoothnessScore, 0.5)

	// Check period counts
	assert.Equal(t, 1, smoothResult.LT2Periods)
	assert.Equal(t, 3, spikyResult.LT2Periods)
}

// Test handling of missing heart rate data
func TestMissingHeartRateData(t *testing.T) {
	config := HRThresholdAnalysisConfig{
		LT1:                     150,
		LT2:                     170,
		MinValidPointsPerBucket: 5,
	}

	// Create 100 seconds with some missing data
	heartRates := make([]Optional[uint8], 100)
	for i := range heartRates {
		if i%10 < 7 { // 7 out of 10 points valid in each bucket
			heartRates[i] = Optional[uint8]{Value: 150, Valid: true}
		} else {
			heartRates[i] = Optional[uint8]{Valid: false}
		}
	}

	ts := createTimeseriesWithGaps(time.Now(), heartRates)
	result, err := AnalyzeHeartRateThresholds(ts, config)
	assert.NoError(t, err)

	// Should still detect LT1 (7 valid points per 10s bucket meets minimum)
	assert.Greater(t, result.TimeAtLT1Seconds, 0)
}

// Test bucket with insufficient data points
func TestInsufficientDataPoints(t *testing.T) {
	config := HRThresholdAnalysisConfig{
		LT1:                     150,
		LT2:                     170,
		MinValidPointsPerBucket: 5,
	}

	// Create data where first bucket has only 3 valid points
	heartRates := make([]Optional[uint8], 100)
	for i := range 3 {
		heartRates[i] = Optional[uint8]{Value: 150, Valid: true}
	}
	for i := 3; i < 10; i++ {
		heartRates[i] = Optional[uint8]{Valid: false}
	}
	// Rest of data is valid LT1
	for i := 10; i < 100; i++ {
		heartRates[i] = Optional[uint8]{Value: 150, Valid: true}
	}

	ts := createTimeseriesWithGaps(time.Now(), heartRates)
	result, err := AnalyzeHeartRateThresholds(ts, config)
	assert.NoError(t, err)

	// First bucket should be skipped, but remaining 90s should count
	assert.Equal(t, 90, result.TimeAtLT1Seconds)
}

// Test transition from LT1 to LT2
func TestLT1ToLT2Transition(t *testing.T) {
	config := HRThresholdAnalysisConfig{
		LT1: 150,
		LT2: 170,
	}

	// 60s at LT1, then 60s at LT2
	heartRates := make([]uint8, 120)
	for i := range 60 {
		heartRates[i] = 150
	}
	for i := 60; i < 120; i++ {
		heartRates[i] = 170
	}

	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	ts := createTimeseries(startTime, heartRates)
	result, err := AnalyzeHeartRateThresholds(ts, config)
	assert.NoError(t, err)

	assert.Equal(t, 60, result.TimeAtLT1Seconds)
	assert.Equal(t, 60, result.TimeAtLT2Seconds)
	assert.Equal(t, 1, result.LT1Periods)
	assert.Equal(t, 1, result.LT2Periods)

	// Verify periods are separate and correctly timed
	assert.Len(t, result.Periods, 2)
	assert.Equal(t, ThresholdLT1, result.Periods[0].Zone)
	assert.Equal(t, ThresholdLT2, result.Periods[1].Zone)
	assert.Equal(t, startTime, result.Periods[0].StartTime)
	assert.Equal(t, startTime.Add(60*time.Second), result.Periods[0].EndTime)
	assert.Equal(t, startTime.Add(60*time.Second), result.Periods[1].StartTime)
}

// Test heart rate statistics in periods
func TestPeriodHeartRateStatistics(t *testing.T) {
	config := HRThresholdAnalysisConfig{
		LT1: 150,
		LT2: 170,
	}

	// Create varying heart rates within LT1 zone
	heartRates := make([]uint8, 60)
	for i := range heartRates {
		// Vary between 145 and 155
		heartRates[i] = 145 + uint8(i%11)
	}

	ts := createTimeseries(time.Now(), heartRates)
	result, err := AnalyzeHeartRateThresholds(ts, config)
	assert.NoError(t, err)

	assert.Len(t, result.Periods, 1)
	period := result.Periods[0]

	assert.Greater(t, period.AvgHR, 145.0)
	assert.Less(t, period.AvgHR, 155.0)
	assert.GreaterOrEqual(t, period.MinHR, 145.0)
	assert.LessOrEqual(t, period.MaxHR, 155.0)
	assert.Less(t, period.MinHR, period.AvgHR)
	assert.Greater(t, period.MaxHR, period.AvgHR)
}

// Test custom configuration parameters
func TestCustomConfiguration(t *testing.T) {
	config := HRThresholdAnalysisConfig{
		LT1:                       150,
		LT2:                       170,
		BucketSizeSeconds:         5,    // 5s buckets instead of 10s
		ThresholdTolerancePercent: 0.10, // ±10% instead of ±5%
		MinConsecutiveBuckets:     10,   // 50s minimum (10 x 5s)
	}

	// 50 seconds at LT1
	heartRates := make([]uint8, 50)
	for i := range heartRates {
		heartRates[i] = 150
	}

	ts := createTimeseries(time.Now(), heartRates)
	result, err := AnalyzeHeartRateThresholds(ts, config)
	assert.NoError(t, err)

	assert.Equal(t, 50, result.TimeAtLT1Seconds)
	assert.Equal(t, 1, result.LT1Periods)
}

// Test empty result when no thresholds are met
func TestNoThresholdsReached(t *testing.T) {
	config := HRThresholdAnalysisConfig{
		LT1: 150,
		LT2: 170,
	}

	// Heart rate stays below both thresholds
	heartRates := make([]uint8, 100)
	for i := range heartRates {
		heartRates[i] = 120
	}

	ts := createTimeseries(time.Now(), heartRates)
	result, err := AnalyzeHeartRateThresholds(ts, config)
	assert.NoError(t, err)

	assert.Equal(t, 0, result.TimeAtLT1Seconds)
	assert.Equal(t, 0, result.TimeAtLT2Seconds)
	assert.Equal(t, 0, result.LT1Periods)
	assert.Equal(t, 0, result.LT2Periods)
	assert.Len(t, result.Periods, 0)
	assert.Equal(t, float64(0), result.LT1SmoothnessScore)
	assert.Equal(t, float64(0), result.LT2SmoothnessScore)
}
