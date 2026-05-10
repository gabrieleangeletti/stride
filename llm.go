package stride

import "math"

type LLMRunSummary struct {
	Metadata       RunMetadata        `json:"runMetadata"`
	GlobalAverages GlobalAverages     `json:"globalAverages"`
	Distributions  Distributions      `json:"distributionsPctTime"`
	TerrainStats   TerrainPerformance `json:"performanceByTerrain"`
	TopoSplits     []TopographicSplit `json:"topographicSplits"`
}

type RunMetadata struct {
	DistanceKm           float64 `json:"distanceKm"`
	MovingTimeMin        int     `json:"movingTimeMin"`
	TotalAscentM         int     `json:"totalAscentM"`
	TotalDescentM        int     `json:"totalDescentM"`
	AerobicDecouplingPct float64 `json:"aerobicDecouplingPct"`
}

type GlobalAverages struct {
	HRAvg      int    `json:"hrAvg"`
	HRMax      int    `json:"hrMax"`
	PaceAvg    string `json:"paceAvgMinKm"`
	GAPAvg     string `json:"gapAvgMinKm"`
	CadenceAvg int    `json:"cadenceAvg"`
	Vam        int    `json:"vamMHr"`
}

type Distributions struct {
	HRZones           []int `json:"hrZones1To5Pct"`
	GradeSteepDownPct int   `json:"gradeSteepDown"`
	GradeRunDownPct   int   `json:"gradeRunDown"`
	GradeFlatPct      int   `json:"gradeFlat"`
	GradeRunUpPct     int   `json:"gradeRunUp"`
	GradeHikeUpPct    int   `json:"gradeHikeUp"`
}

type TerrainPerformance struct {
	UphillAvgHR     int    `json:"uphillAvgHr"`
	UphillAvgGAP    string `json:"uphillAvgGap"`
	DownhillAvgHR   int    `json:"downhillAvgHr"`
	DownhillAvgPace string `json:"downhillAvgPace"`
}

type TopographicSplit struct {
	Type     string  `json:"type"` // "Uphill", "Downhill", "Flat"
	DistKm   float64 `json:"distKm"`
	GradePct float64 `json:"gradePct"`
	AvgGAP   string  `json:"avgGap"`
	AvgHR    int     `json:"avgHr"`

	// Internal accumulators for finalization
	sumGAPSpeed float64 `json:"-"`
	sumHR       int     `json:"-"`
	pointsCount int     `json:"-"`
}

type LLMSummaryConfig struct {
	GradeSmoothingWindow int     // Window size for smoothing grade (data points, roughly seconds). Default 15.
	MinGradeDeltaM       float64 // Minimum distance moved (meters) to calculate a stable grade. Default 5.0.
	MinMovingSpeedMS     float64 // Minimum speed (m/s) to consider as "moving" for stats. Default 0.5.
	PhaseThresholdM      float64 // Minimum continuous elevation change (meters) for climb/descent. Default 30.0.
	MinSegmentDistM      float64 // Minimum segment distance (meters); shorter segments are skipped. Default 100.

	DefaultMaxHR     int     // Fallback max HR when not available from data. Default 190.
	HRZone1Threshold float64 // Fraction of max HR for zone 1→2 boundary. Default 0.60.
	HRZone2Threshold float64 // Fraction of max HR for zone 2→3 boundary. Default 0.70.
	HRZone3Threshold float64 // Fraction of max HR for zone 3→4 boundary. Default 0.80.
	HRZone4Threshold float64 // Fraction of max HR for zone 4→5 boundary. Default 0.90.

	GradeUpThreshold    float64 // Grade % above which is considered uphill terrain. Default 2.0.
	GradeDownThreshold  float64 // Grade % below which is considered downhill terrain. Default -2.0.
	GradeSteepDownBelow float64 // Grade % below which is considered steep downhill. Default -10.0.
	GradeHikeUpAbove    float64 // Grade % above which is considered hike up. Default 8.0.
}

func (c LLMSummaryConfig) ApplyDefaults() LLMSummaryConfig {
	config := c

	if config.GradeSmoothingWindow == 0 {
		config.GradeSmoothingWindow = 15
	}
	if config.MinGradeDeltaM == 0 {
		config.MinGradeDeltaM = 5.0
	}
	if config.MinMovingSpeedMS == 0 {
		config.MinMovingSpeedMS = 0.5
	}
	if config.PhaseThresholdM == 0 {
		config.PhaseThresholdM = 30.0
	}
	if config.MinSegmentDistM == 0 {
		config.MinSegmentDistM = 100.0
	}
	if config.DefaultMaxHR == 0 {
		config.DefaultMaxHR = 190
	}
	if config.HRZone1Threshold == 0 {
		config.HRZone1Threshold = 0.60
	}
	if config.HRZone2Threshold == 0 {
		config.HRZone2Threshold = 0.70
	}
	if config.HRZone3Threshold == 0 {
		config.HRZone3Threshold = 0.80
	}
	if config.HRZone4Threshold == 0 {
		config.HRZone4Threshold = 0.90
	}
	if config.GradeUpThreshold == 0 {
		config.GradeUpThreshold = 2.0
	}
	if config.GradeDownThreshold == 0 {
		config.GradeDownThreshold = -2.0
	}
	if config.GradeSteepDownBelow == 0 {
		config.GradeSteepDownBelow = -10.0
	}
	if config.GradeHikeUpAbove == 0 {
		config.GradeHikeUpAbove = 8.0
	}

	return config
}

// SummarizeForLLM processes augmented timeseries into the compressed LLM format.
func SummarizeForLLM(act *Activity, ts *ActivityTimeseries, config LLMSummaryConfig) (*LLMRunSummary, error) {
	config = config.ApplyDefaults()

	// First, ensure we have distances and computed metrics
	if act.Distance == 0 || !ts.Data[len(ts.Data)-1].Distance.Valid {
		AugmentGPXData(act, ts)
	}

	summary := &LLMRunSummary{}

	// 1. Metadata & Global Averages
	summary.Metadata.DistanceKm = float64(act.Distance) / 1000.0
	summary.Metadata.MovingTimeMin = int(act.MovingTime) / 60
	if act.ElevationGain.Valid {
		summary.Metadata.TotalAscentM = int(act.ElevationGain.Value)
	}
	if act.ElevationLoss.Valid {
		summary.Metadata.TotalDescentM = int(act.ElevationLoss.Value)
	}

	hrMetrics, _ := ts.HRMetrics()
	if hrMetrics != nil {
		summary.GlobalAverages.HRAvg = int(hrMetrics.AvgHR)
		summary.GlobalAverages.HRMax = int(hrMetrics.MaxHR)
	}

	if act.AvgSpeed > 0 {
		summary.GlobalAverages.PaceAvg = formatPace(float64(act.AvgSpeed) / 1000.0)
	}
	if act.MovingTime > 0 && act.ElevationGain.Valid {
		summary.GlobalAverages.Vam = int(float64(act.ElevationGain.Value) / (float64(act.MovingTime) / 3600.0))
	}

	// State variables
	var hrZoneCounters [5]int
	gradeCounters := map[string]int{"SteepDown": 0, "RunDown": 0, "Flat": 0, "RunUp": 0, "HikeUp": 0}

	// Topographic Split State
	currentSegment := TopographicSplit{Type: "Flat", DistKm: 0}
	var segmentStartElev float64
	var segmentStartDist float64
	if ts.Data[0].Altitude.Valid {
		segmentStartElev = ts.Data[0].Altitude.Value
	}
	if len(ts.Data) > 0 && ts.Data[0].Distance.Valid {
		segmentStartDist = float64(ts.Data[0].Distance.Value)
	}

	// Terrain Performance Accumulators
	var upHrSum, upGapSum, upCount float64
	var downHrSum, downPaceSum, downCount float64

	// Aerobic Decoupling Accumulators
	type DecoupleState struct{ hrSum, speedSum, count float64 }
	halfIndex := act.Distance / 2 // Decouple based on distance
	firstHalf := DecoupleState{}
	secondHalf := DecoupleState{}

	var totalGAPSpeed float64
	var gapPoints int

	// Window for smoothing Grade (15 points / ~15 seconds)
	windowSize := config.GradeSmoothingWindow

	for i := windowSize; i < len(ts.Data); i++ {
		curr := ts.Data[i]
		prevWindow := ts.Data[i-windowSize]

		if !curr.Distance.Valid || !prevWindow.Distance.Valid || !curr.Altitude.Valid || !prevWindow.Altitude.Valid {
			continue
		}

		deltaDist := float64(curr.Distance.Value - prevWindow.Distance.Value)
		deltaElev := curr.Altitude.Value - prevWindow.Altitude.Value

		// Determine speed & GAP
		actualSpeed := 0.0
		if deltaDist > 0 {
			timeDelta := float64(curr.Offset - prevWindow.Offset)
			if timeDelta > 0 {
				actualSpeed = deltaDist / timeDelta
			}
		}

		gradeFraction := 0.0
		gradePct := 0.0
		if deltaDist > config.MinGradeDeltaM {
			gradeFraction = deltaElev / deltaDist
			gradePct = gradeFraction * 100.0
		}

		gapSpeed := calculateGAP(actualSpeed, gradeFraction)
		if actualSpeed > config.MinMovingSpeedMS { // Only track averages if actually moving
			totalGAPSpeed += gapSpeed
			gapPoints++

			// Aerobic Decoupling accumulation
			if curr.HeartRate.Valid {
				if curr.Distance.Value < halfIndex {
					firstHalf.hrSum += float64(curr.HeartRate.Value)
					firstHalf.speedSum += gapSpeed
					firstHalf.count++
				} else {
					secondHalf.hrSum += float64(curr.HeartRate.Value)
					secondHalf.speedSum += gapSpeed
					secondHalf.count++
				}
			}
		}

		// A. Distributions Tracking
		if actualSpeed > config.MinMovingSpeedMS { // Only bin grade if moving
			switch {
			case gradePct < config.GradeSteepDownBelow:
				gradeCounters["SteepDown"]++
			case gradePct >= config.GradeSteepDownBelow && gradePct < config.GradeDownThreshold:
				gradeCounters["RunDown"]++
			case gradePct >= config.GradeDownThreshold && gradePct <= config.GradeUpThreshold:
				gradeCounters["Flat"]++
			case gradePct > config.GradeUpThreshold && gradePct <= config.GradeHikeUpAbove:
				gradeCounters["RunUp"]++
			default:
				gradeCounters["HikeUp"]++
			}
		}

		if curr.HeartRate.Valid {
			hr := float64(curr.HeartRate.Value)
			maxHR := float64(summary.GlobalAverages.HRMax)
			if maxHR == 0 {
				maxHR = float64(config.DefaultMaxHR)
			}
			pctMax := hr / maxHR
			switch {
			case pctMax < config.HRZone1Threshold:
				hrZoneCounters[0]++
			case pctMax < config.HRZone2Threshold:
				hrZoneCounters[1]++
			case pctMax < config.HRZone3Threshold:
				hrZoneCounters[2]++
			case pctMax < config.HRZone4Threshold:
				hrZoneCounters[3]++
			default:
				hrZoneCounters[4]++
			}

			// Terrain Specific Performance
			if gradePct > config.GradeUpThreshold && actualSpeed > config.MinMovingSpeedMS {
				upHrSum += hr
				upGapSum += gapSpeed
				upCount++
			} else if gradePct < config.GradeDownThreshold && actualSpeed > config.MinMovingSpeedMS {
				downHrSum += hr
				downPaceSum += actualSpeed
				downCount++
			}
		}

		// B. Topographic Segmentation (ClimbPro logic)
		// Accumulate elevation from start of current phase
		segmentElevChange := curr.Altitude.Value - segmentStartElev

		isClimbing := segmentElevChange > config.PhaseThresholdM
		isDescending := segmentElevChange < -config.PhaseThresholdM

		if isClimbing && currentSegment.Type != "Uphill" {
			finalizeSegment(&summary.TopoSplits, &currentSegment, segmentStartElev, segmentStartDist, curr, config)
			currentSegment = TopographicSplit{Type: "Uphill"}
			segmentStartElev = curr.Altitude.Value
			segmentStartDist = float64(curr.Distance.Value)
		} else if isDescending && currentSegment.Type != "Downhill" {
			finalizeSegment(&summary.TopoSplits, &currentSegment, segmentStartElev, segmentStartDist, curr, config)
			currentSegment = TopographicSplit{Type: "Downhill"}
			segmentStartElev = curr.Altitude.Value
			segmentStartDist = float64(curr.Distance.Value)
		}

		// Accumulate segment stats
		if actualSpeed > config.MinMovingSpeedMS {
			currentSegment.sumGAPSpeed += gapSpeed
			if curr.HeartRate.Valid {
				currentSegment.sumHR += int(curr.HeartRate.Value)
			}
			currentSegment.pointsCount++
		}
	}

	// Finalize final segment
	if len(ts.Data) > 0 {
		finalizeSegment(&summary.TopoSplits, &currentSegment, segmentStartElev, segmentStartDist, ts.Data[len(ts.Data)-1], config)
	}

	// Final calculations
	if gapPoints > 0 {
		summary.GlobalAverages.GAPAvg = formatPace(totalGAPSpeed / float64(gapPoints))
	}

	// Calculate Distributions
	totalGradePoints := gradeCounters["SteepDown"] + gradeCounters["RunDown"] + gradeCounters["Flat"] + gradeCounters["RunUp"] + gradeCounters["HikeUp"]
	if totalGradePoints > 0 {
		summary.Distributions.GradeSteepDownPct = (gradeCounters["SteepDown"] * 100) / totalGradePoints
		summary.Distributions.GradeRunDownPct = (gradeCounters["RunDown"] * 100) / totalGradePoints
		summary.Distributions.GradeFlatPct = (gradeCounters["Flat"] * 100) / totalGradePoints
		summary.Distributions.GradeRunUpPct = (gradeCounters["RunUp"] * 100) / totalGradePoints
		summary.Distributions.GradeHikeUpPct = (gradeCounters["HikeUp"] * 100) / totalGradePoints
	}

	totalHRPoints := hrZoneCounters[0] + hrZoneCounters[1] + hrZoneCounters[2] + hrZoneCounters[3] + hrZoneCounters[4]
	if totalHRPoints > 0 {
		summary.Distributions.HRZones = []int{
			(hrZoneCounters[0] * 100) / totalHRPoints,
			(hrZoneCounters[1] * 100) / totalHRPoints,
			(hrZoneCounters[2] * 100) / totalHRPoints,
			(hrZoneCounters[3] * 100) / totalHRPoints,
			(hrZoneCounters[4] * 100) / totalHRPoints,
		}
	}

	if upCount > 0 {
		summary.TerrainStats.UphillAvgHR = int(upHrSum / upCount)
		summary.TerrainStats.UphillAvgGAP = formatPace(upGapSum / upCount)
	}
	if downCount > 0 {
		summary.TerrainStats.DownhillAvgHR = int(downHrSum / downCount)
		summary.TerrainStats.DownhillAvgPace = formatPace(downPaceSum / downCount) // We care about actual pace downhill, not GAP
	}

	// Aerobic Decoupling (Efficiency Factor = Speed / HR. Decoupling = % drop in efficiency)
	if firstHalf.count > 0 && secondHalf.count > 0 {
		ef1 := (firstHalf.speedSum / firstHalf.count) / (firstHalf.hrSum / firstHalf.count)
		ef2 := (secondHalf.speedSum / secondHalf.count) / (secondHalf.hrSum / secondHalf.count)
		if ef1 > 0 {
			// Positive decoupling means athlete became less efficient (fatigued)
			summary.Metadata.AerobicDecouplingPct = math.Round(((ef1-ef2)/ef1)*100.0*10) / 10 // rounded to 1 decimal
		}
	}

	return summary, nil
}

func finalizeSegment(splits *[]TopographicSplit, seg *TopographicSplit, startElev, startDist float64, curr ActivityTimeseriesEntry, config LLMSummaryConfig) {
	deltaDist := float64(curr.Distance.Value) - startDist
	if deltaDist < config.MinSegmentDistM {
		return
	}

	seg.DistKm = deltaDist / 1000.0

	deltaElev := curr.Altitude.Value - startElev
	seg.GradePct = math.Round(deltaElev/deltaDist*100*10) / 10

	if seg.pointsCount > 0 {
		seg.AvgGAP = formatPace(seg.sumGAPSpeed / float64(seg.pointsCount))
		seg.AvgHR = seg.sumHR / seg.pointsCount
	}

	*splits = append(*splits, *seg)
}
