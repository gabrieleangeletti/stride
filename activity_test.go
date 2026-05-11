package stride_test

import (
	"math"
	"testing"

	"github.com/gabrieleangeletti/stride"
)

// Tests that the Elevation Hysteresis no longer drops valid climbs due to minor track dips
func TestAugmentGPXData_ElevationHysteresis(t *testing.T) {
	ts := &stride.ActivityTimeseries{
		Data: []stride.ActivityTimeseriesEntry{
			{Offset: 0, Altitude: stride.Optional[float64]{Value: 100, Valid: true}, Latitude: stride.Optional[float64]{Value: 0, Valid: true}, Longitude: stride.Optional[float64]{Value: 0, Valid: true}},
			{Offset: 10, Altitude: stride.Optional[float64]{Value: 102, Valid: true}, Latitude: stride.Optional[float64]{Value: 0.001, Valid: true}, Longitude: stride.Optional[float64]{Value: 0, Valid: true}}, // +2m
			{Offset: 20, Altitude: stride.Optional[float64]{Value: 101, Valid: true}, Latitude: stride.Optional[float64]{Value: 0.002, Valid: true}, Longitude: stride.Optional[float64]{Value: 0, Valid: true}}, // -1m (Dip shouldn't reset our +2)
			{Offset: 30, Altitude: stride.Optional[float64]{Value: 104, Valid: true}, Latitude: stride.Optional[float64]{Value: 0.003, Valid: true}, Longitude: stride.Optional[float64]{Value: 0, Valid: true}}, // +3m
		},
	}

	act := &stride.Activity{}
	config := stride.AugmentConfig{ElevationHysteresisM: 3.0}
	stride.AugmentGPXData(act, ts, config)

	// We started at 100m, ended at 104m. Because the dip wasn't > 3m,
	// we should capture the net 4m climb once it breaks the threshold!
	if act.ElevationGain.Value != 4 {
		t.Errorf("Expected 4m Elevation Gain, got %d", act.ElevationGain.Value)
	}
}

// Tests that topographic splits isolate hills based on extrema, not segment origins
func TestDetectTopographicSplits_Extrema(t *testing.T) {
	enriched := []stride.EnrichedPoint{
		{DistanceM: 0, Entry: &stride.ActivityTimeseriesEntry{Altitude: stride.Optional[float64]{Value: 100, Valid: true}}},
		{DistanceM: 100, Entry: &stride.ActivityTimeseriesEntry{Altitude: stride.Optional[float64]{Value: 140, Valid: true}}}, // Climb 40m
		{DistanceM: 200, Entry: &stride.ActivityTimeseriesEntry{Altitude: stride.Optional[float64]{Value: 135, Valid: true}}}, // Dip 5m (ignores)
		{DistanceM: 300, Entry: &stride.ActivityTimeseriesEntry{Altitude: stride.Optional[float64]{Value: 180, Valid: true}}}, // Peak!
		{DistanceM: 400, Entry: &stride.ActivityTimeseriesEntry{Altitude: stride.Optional[float64]{Value: 140, Valid: true}}}, // Drop 40m
	}

	config := stride.LLMSummaryConfig{PhaseThresholdM: 30.0}
	splits := stride.DetectTopographicSplits(enriched, config)

	if len(splits) != 2 {
		t.Fatalf("Expected exactly 2 splits (Uphill, Downhill), got %d", len(splits))
	}
	if splits[0].Type != "Uphill" || splits[0].DistKm != 0.3 {
		t.Errorf("Expected Uphill segment of 0.3km, got %s of %fkm", splits[0].Type, splits[0].DistKm)
	}
	if splits[1].Type != "Downhill" || splits[1].DistKm != 0.1 {
		t.Errorf("Expected Downhill segment of 0.1km, got %s of %fkm", splits[1].Type, splits[1].DistKm)
	}
}

// Demonstrates how the WeightedAvg abstraction fixes Smart Recording skews
func TestWeightedAvg(t *testing.T) {
	var hrAvg stride.WeightedAvg

	// Runner stays in Z1 (130 HR) for 60 seconds (long straightaway, 1 point)
	hrAvg.Add(130, 60)

	// Runner hits a short hill, HR spikes to 180, watch records densely (1 sec intervals x 5)
	for range 5 {
		hrAvg.Add(180, 1)
	}

	// Unweighted Average: (130 + 180 + 180 + 180 + 180 + 180) / 6 = 171 Avg HR (WRONG)
	// Time-Weighted Avg: ((130*60) + (180*5)) / 65 = 133 Avg HR (CORRECT)
	if math.Round(hrAvg.Avg()) != 134 { // 133.8
		t.Errorf("Time weighting failed, expected 134, got %f", hrAvg.Avg())
	}
}
