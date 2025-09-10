package stride_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stride "github.com/gabrieleangeletti/stride"
)

func TestCalculateAverageHeartRate(t *testing.T) {
	timeseries := &stride.ActivityTimeseries{
		StartTime: time.Now(),
		Data: []stride.ActivityTimeseriesEntry{
			{Offset: 0, HeartRate: stride.Optional[uint8]{Value: 120, Valid: true}},
			{Offset: 10, HeartRate: stride.Optional[uint8]{Value: 130, Valid: true}},
			{Offset: 20, HeartRate: stride.Optional[uint8]{Value: 125, Valid: true}},
			{Offset: 30, HeartRate: stride.Optional[uint8]{Value: 140, Valid: true}},
			{Offset: 40, HeartRate: stride.Optional[uint8]{Value: 0, Valid: false}},
			{Offset: 50, HeartRate: stride.Optional[uint8]{Value: 135, Valid: true}},
			{Offset: 60, HeartRate: stride.Optional[uint8]{Value: 145, Valid: true}},
			{Offset: 70, HeartRate: stride.Optional[uint8]{Value: 150, Valid: true}},
			{Offset: 80, HeartRate: stride.Optional[uint8]{Value: 130, Valid: true}},
		},
	}

	config := stride.AvgHeartRateAnalysisConfig{
		Method:       stride.HeartRateMethodTimeWeighted,
		ExcludeZeros: true,
		MinValidRate: 40,
		MaxValidRate: 220,
		MaxHeartRate: 190,
	}

	t.Run("TimeWeighted", func(t *testing.T) {
		avgHR, err := stride.CalculateAverageHeartRate(timeseries, config)
		require.NoError(t, err)

		// Time-weighted: Only uses intervals between consecutive entries
		// (120*10 + 130*10 + 125*10 + 140*10 + 135*10 + 145*10 + 150*10) / 70
		// = (1200 + 1300 + 1250 + 1400 + 1350 + 1450 + 1500) / 70 = 9450 / 70 = 135.0
		// Note: Last entry (130) has no next entry, so it's not included in time-weighted calc
		expected := 135.0
		assert.InDelta(t, expected, avgHR, 0.01, "Time-weighted average should be close to expected value")
	})

	t.Run("Simple", func(t *testing.T) {
		config.Method = stride.HeartRateMethodSimple
		simpleAvgHR, err := stride.CalculateAverageHeartRate(timeseries, config)
		require.NoError(t, err)

		// Simple average: (120+130+125+140+135+145+150+130) / 8 = 1075 / 8 = 134.375
		expected := 134.375
		assert.InDelta(t, expected, simpleAvgHR, 0.01, "Simple average should be close to expected value")
	})

	t.Run("ZoneWeighted", func(t *testing.T) {
		config.Method = stride.HeartRateMethodZoneWeighted
		zoneAvgHR, err := stride.CalculateAverageHeartRate(timeseries, config)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, zoneAvgHR, 120.0, "Zone-weighted average should be at least 120")
		assert.LessOrEqual(t, zoneAvgHR, 160.0, "Zone-weighted average should be at most 160")
	})
}

func TestCalculateMaxHeartRate(t *testing.T) {
	timeseries := &stride.ActivityTimeseries{
		StartTime: time.Now(),
		Data: []stride.ActivityTimeseriesEntry{
			{Offset: 0, HeartRate: stride.Optional[uint8]{Value: 120, Valid: true}},
			{Offset: 10, HeartRate: stride.Optional[uint8]{Value: 130, Valid: true}},
			{Offset: 20, HeartRate: stride.Optional[uint8]{Value: 165, Valid: true}},
			{Offset: 30, HeartRate: stride.Optional[uint8]{Value: 140, Valid: true}},
			{Offset: 40, HeartRate: stride.Optional[uint8]{Value: 0, Valid: false}},
			{Offset: 50, HeartRate: stride.Optional[uint8]{Value: 175, Valid: true}}, // Peak
			{Offset: 60, HeartRate: stride.Optional[uint8]{Value: 145, Valid: true}},
			{Offset: 70, HeartRate: stride.Optional[uint8]{Value: 150, Valid: true}},
			{Offset: 80, HeartRate: stride.Optional[uint8]{Value: 130, Valid: true}},
		},
	}

	t.Run("Peak", func(t *testing.T) {
		config := stride.MaxHeartRateAnalysisConfig{
			Method:       stride.MaxHeartRateMethodPeak,
			ExcludeZeros: true,
			MinValidRate: 40,
			MaxValidRate: 220,
		}

		maxHR, err := stride.CalculateMaxHeartRate(timeseries, config)
		require.NoError(t, err)
		assert.Equal(t, 175, maxHR, "Peak max HR should be 175")
	})

	t.Run("RollingWindow", func(t *testing.T) {
		windowDuration := 30 * time.Second
		config := stride.MaxHeartRateAnalysisConfig{
			Method:         stride.MaxHeartRateMethodRollingWindow,
			ExcludeZeros:   true,
			MinValidRate:   40,
			MaxValidRate:   220,
			WindowDuration: &windowDuration,
		}

		maxHR, err := stride.CalculateMaxHeartRate(timeseries, config)
		require.NoError(t, err)
		assert.Equal(t, 175, maxHR, "Rolling window max HR should be 175")
	})

	t.Run("Percentile95", func(t *testing.T) {
		config := stride.MaxHeartRateAnalysisConfig{
			Method:          stride.MaxHeartRateMethodPercentile,
			ExcludeZeros:    true,
			MinValidRate:    40,
			MaxValidRate:    220,
			PercentileValue: 95.0,
		}

		maxHR, err := stride.CalculateMaxHeartRate(timeseries, config)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, maxHR, 165, "95th percentile max HR should be at least 165")
		assert.LessOrEqual(t, maxHR, 175, "95th percentile max HR should be at most 175")
	})

	t.Run("Percentile90", func(t *testing.T) {
		config := stride.MaxHeartRateAnalysisConfig{
			Method:          stride.MaxHeartRateMethodPercentile,
			ExcludeZeros:    true,
			MinValidRate:    40,
			MaxValidRate:    220,
			PercentileValue: 90.0,
		}

		maxHR, err := stride.CalculateMaxHeartRate(timeseries, config)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, maxHR, 150, "90th percentile max HR should be at least 150")
		assert.LessOrEqual(t, maxHR, 175, "90th percentile max HR should be at most 175")
	})
}

func TestCalculateAverageHeartRateEdgeCases(t *testing.T) {
	t.Run("EmptyData", func(t *testing.T) {
		timeseries := &stride.ActivityTimeseries{
			StartTime: time.Now(),
			Data:      []stride.ActivityTimeseriesEntry{},
		}
		config := stride.AvgHeartRateAnalysisConfig{
			Method: stride.HeartRateMethodSimple,
		}

		_, err := stride.CalculateAverageHeartRate(timeseries, config)
		assert.ErrorIs(t, err, stride.ErrEmptyTimeseriesData, "Should return ErrEmptyTimeseriesData for empty data")
	})

	t.Run("NoValidHeartRateData", func(t *testing.T) {
		timeseries := &stride.ActivityTimeseries{
			StartTime: time.Now(),
			Data: []stride.ActivityTimeseriesEntry{
				{Offset: 0, HeartRate: stride.Optional[uint8]{Value: 0, Valid: false}},
				{Offset: 10, HeartRate: stride.Optional[uint8]{Value: 0, Valid: false}},
			},
		}
		config := stride.AvgHeartRateAnalysisConfig{
			Method:       stride.HeartRateMethodSimple,
			ExcludeZeros: true,
		}

		_, err := stride.CalculateAverageHeartRate(timeseries, config)
		assert.ErrorIs(t, err, stride.ErrNoValidData, "Should return ErrNoValidData when no valid heart rate data")
	})

	t.Run("FilterByValidRate", func(t *testing.T) {
		timeseries := &stride.ActivityTimeseries{
			StartTime: time.Now(),
			Data: []stride.ActivityTimeseriesEntry{
				{Offset: 0, HeartRate: stride.Optional[uint8]{Value: 30, Valid: true}},   // Too low
				{Offset: 10, HeartRate: stride.Optional[uint8]{Value: 120, Valid: true}}, // Valid
				{Offset: 20, HeartRate: stride.Optional[uint8]{Value: 250, Valid: true}}, // Too high
			},
		}
		config := stride.AvgHeartRateAnalysisConfig{
			Method:       stride.HeartRateMethodSimple,
			MinValidRate: 40,
			MaxValidRate: 220,
		}

		avgHR, err := stride.CalculateAverageHeartRate(timeseries, config)
		require.NoError(t, err)
		assert.Equal(t, 120.0, avgHR, "Should only use the valid heart rate reading")
	})
}

func TestCalculateMaxHeartRateEdgeCases(t *testing.T) {
	t.Run("EmptyData", func(t *testing.T) {
		timeseries := &stride.ActivityTimeseries{
			StartTime: time.Now(),
			Data:      []stride.ActivityTimeseriesEntry{},
		}
		config := stride.MaxHeartRateAnalysisConfig{
			Method: stride.MaxHeartRateMethodPeak,
		}

		_, err := stride.CalculateMaxHeartRate(timeseries, config)
		assert.ErrorIs(t, err, stride.ErrEmptyTimeseriesData, "Should return ErrEmptyTimeseriesData for empty data")
	})

	t.Run("NoValidData", func(t *testing.T) {
		timeseries := &stride.ActivityTimeseries{
			StartTime: time.Now(),
			Data: []stride.ActivityTimeseriesEntry{
				{Offset: 0, HeartRate: stride.Optional[uint8]{Value: 0, Valid: false}},
				{Offset: 10, HeartRate: stride.Optional[uint8]{Value: 0, Valid: false}},
			},
		}
		config := stride.MaxHeartRateAnalysisConfig{
			Method:       stride.MaxHeartRateMethodPeak,
			ExcludeZeros: true,
		}

		_, err := stride.CalculateMaxHeartRate(timeseries, config)
		assert.ErrorIs(t, err, stride.ErrNoValidData, "Should return ErrNoValidData when no valid heart rate data")
	})

	t.Run("FilterByValidRate", func(t *testing.T) {
		timeseries := &stride.ActivityTimeseries{
			StartTime: time.Now(),
			Data: []stride.ActivityTimeseriesEntry{
				{Offset: 0, HeartRate: stride.Optional[uint8]{Value: 30, Valid: true}},   // Too low
				{Offset: 10, HeartRate: stride.Optional[uint8]{Value: 160, Valid: true}}, // Valid
				{Offset: 20, HeartRate: stride.Optional[uint8]{Value: 250, Valid: true}}, // Too high
			},
		}
		config := stride.MaxHeartRateAnalysisConfig{
			Method:       stride.MaxHeartRateMethodPeak,
			MinValidRate: 40,
			MaxValidRate: 220,
		}

		maxHR, err := stride.CalculateMaxHeartRate(timeseries, config)
		require.NoError(t, err)
		assert.Equal(t, 160, maxHR, "Should only use the valid heart rate reading")
	})
}
