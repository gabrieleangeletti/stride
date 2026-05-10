package stride

import "math"

type LLMRunSummary struct {
	Metadata       RunMetadata        `json:"runMetadata"`
	GlobalAverages GlobalAverages     `json:"globalAverages"`
	Distributions  Distributions      `json:"distributionsPctTime"`
	TerrainStats   TerrainPerformance `json:"performanceByTerrain"`
	TopoSplits     []TopographicSplit `json:"topographicSplits"`
	Decoupling     DecouplingDetail   `json:"decoupling"`
	ZoneDetails    []ZoneDetail       `json:"zoneDetails"`
	GAPBenchmarks  []GAPBenchmark     `json:"gapBenchmarks"`
	Thresholds     ThresholdMetrics   `json:"thresholds"`
	EndOfRun       EndOfRunMetrics    `json:"endOfRun"`
	Recovery       RecoveryMetrics    `json:"recovery"`
	Economy        EconomyMetrics     `json:"economy"`
	Athlete        AthleteBaseline    `json:"athlete"`
}

type RunMetadata struct {
	DistanceKm    float64 `json:"distanceKm"`
	MovingTimeMin int     `json:"movingTimeMin"`
	TotalAscentM  int     `json:"totalAscentM"`
	TotalDescentM int     `json:"totalDescentM"`
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
	HRZones           []int  `json:"hrZones1To5Pct"`
	GradeSteepDownPct int    `json:"gradeSteepDown"`
	GradeRunDownPct   int    `json:"gradeRunDown"`
	GradeFlatPct      int    `json:"gradeFlat"`
	GradeRunUpPct     int    `json:"gradeRunUp"`
	GradeHikeUpPct    int    `json:"gradeHikeUp"`
	Z2AvgPace         string `json:"z2AvgPace"`
	Z2AvgGAP          string `json:"z2AvgGap"`
}

type TerrainPerformance struct {
	UphillAvgHR               int           `json:"uphillAvgHr"`
	UphillAvgGAP              string        `json:"uphillAvgGap"`
	UphillHRStdDev            float64       `json:"uphillHrStdDev"`
	UphillVAM                 int           `json:"uphillVam"`
	DownhillAvgHR             int           `json:"downhillAvgHr"`
	DownhillAvgPace           string        `json:"downhillAvgPace"`
	DownhillEfficiency        float64       `json:"downhillEfficiency"`
	HikeRunTransitionGradePct float64       `json:"hikeRunTransitionGradePct"`
	VAMByGradient             []GradientVAM `json:"vamByGradient"`
	DownhillPaceByGrade       []GradePace   `json:"downhillPaceByGrade"`
}

type DecouplingDetail struct {
	FirstHalfAvgHR       float64 `json:"firstHalfAvgHr"`
	FirstHalfAvgGAP      string  `json:"firstHalfAvgGap"`
	SecondHalfAvgHR      float64 `json:"secondHalfAvgHr"`
	SecondHalfAvgGAP     string  `json:"secondHalfAvgGap"`
	AerobicDecouplingPct float64 `json:"aerobicDecouplingPct"`
	UphillDecouplingPct  float64 `json:"uphillDecouplingPct"`
}

type ZoneDetail struct {
	Zone    int    `json:"zone"`
	PctTime int    `json:"pctTime"`
	AvgGAP  string `json:"avgGap"`
	AvgHR   int    `json:"avgHr"`
}

type GAPBenchmark struct {
	Range string `json:"range"`
	AvgHR int    `json:"avgHr"`
	Count int    `json:"count"`
}

type EndOfRunMetrics struct {
	Last10PctAvgHR   int    `json:"last10PctAvgHr"`
	Last10PctAvgGAP  string `json:"last10PctAvgGap"`
	Last10PctAvgPace string `json:"last10PctAvgPace"`
}

type RecoveryMetrics struct {
	EndHRDrop60s           int `json:"endHrDrop60s"`
	PostMaxEffortHRDrop60s int `json:"postMaxEffortHrDrop60s"`
}

type EconomyMetrics struct {
	HRPerPaceFlat  float64 `json:"hrPerPaceFlat"`
	HRPerVAMUphill float64 `json:"hrPerVamUphill"`
}

type AthleteBaseline struct {
	MaxHR     int `json:"maxHr"`
	RestingHR int `json:"restingHr"`
	AeTHR     int `json:"aetHr"`
	AnTHR     int `json:"antHr"`
}

type ThresholdMetrics struct {
	Z4Z5AvgGAP        string `json:"z4z5AvgGap"`
	LongestZ4BlockSec int    `json:"longestZ4BlockSec"`
	ZoneThresholds    []int  `json:"zoneHrThresholds"`
}

type GradientVAM struct {
	Band    string `json:"band"`
	VAM     int    `json:"vam"`
	AvgHR   int    `json:"avgHr"`
	PctTime int    `json:"pctTime"`
}

type GradePace struct {
	Band    string `json:"band"`
	AvgPace string `json:"avgPace"`
	AvgHR   int    `json:"avgHr"`
}

type TopographicSplit struct {
	Type     string  `json:"type"`
	DistKm   float64 `json:"distKm"`
	GradePct float64 `json:"gradePct"`
	AvgGAP   string  `json:"avgGap"`
	AvgHR    int     `json:"avgHr"`

	sumGAPSpeed float64 `json:"-"`
	sumHR       int     `json:"-"`
	pointsCount int     `json:"-"`
}

type LLMSummaryConfig struct {
	GradeSmoothingWindow int             // Default 15.
	MinGradeDeltaM       float64         // Default 5.0.
	MinMovingSpeedMS     float64         // Default 0.5.
	PhaseThresholdM      float64         // Default 30.0.
	MinSegmentDistM      float64         // Default 100.
	GradeUpThreshold     float64         // Default 2.0.
	GradeDownThreshold   float64         // Default -2.0.
	GradeSteepDownBelow  float64         // Default -10.0.
	GradeHikeUpAbove     float64         // Default 8.0.
	ElevationHysteresisM float64         // Default 3.0.
	Athlete              AthleteBaseline // Optional: copied to output as-is.
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
	if config.ElevationHysteresisM == 0 {
		config.ElevationHysteresisM = 3.0
	}

	return config
}

// paceToSpeed thresholds for GAP benchmark bands (m/s).
var gapBenchmarkThresholds = []struct {
	label string
	speed float64
}{
	{"<4:00", 4.167},
	{"4:00-4:30", 3.704},
	{"4:30-5:00", 3.333},
	{"5:00-5:30", 3.030},
	{"5:30-6:00", 2.778},
	{"6:00-6:30", 2.564},
	{"6:30-7:00", 2.381},
	{"7:00-8:00", 2.083},
	{">8:00", 0},
}

type vamBand struct {
	name      string
	gradeLow  float64
	gradeHigh float64
}

var vamBands = []vamBand{
	{"2-5%", 2.0, 5.0},
	{"5-10%", 5.0, 10.0},
	{"10-15%", 10.0, 15.0},
	{"15-20%", 15.0, 20.0},
	{">20%", 20.0, 999.0},
}

type dhBand struct {
	name      string
	gradeLow  float64
	gradeHigh float64
}

var dhBands = []dhBand{
	{"-2 to -5%", -5.0, -2.0},
	{"-5 to -10%", -10.0, -5.0},
	{"<-10%", -999.0, -10.0},
}

// SummarizeForLLM processes augmented timeseries into the compressed LLM format.
func SummarizeForLLM(act *Activity, ts *ActivityTimeseries, config LLMSummaryConfig) (*LLMRunSummary, error) {
	config = config.ApplyDefaults()

	if act.Distance == 0 || !ts.Data[len(ts.Data)-1].Distance.Valid {
		AugmentGPXData(act, ts, AugmentConfig{ElevationHysteresisM: config.ElevationHysteresisM})
	}

	summary := &LLMRunSummary{}

	summary.Athlete = config.Athlete

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
	maxHR := float64(config.Athlete.MaxHR)
	if maxHR == 0 {
		maxHR = float64(summary.GlobalAverages.HRMax)
	}
	if maxHR == 0 {
		maxHR = 190
	}

	z1Max := float64(config.Athlete.AeTHR) * 0.9
	z2Max := float64(config.Athlete.AeTHR)
	if config.Athlete.AeTHR == 0 {
		z1Max = maxHR * 0.60
		z2Max = maxHR * 0.70
	}
	z3Max := float64(config.Athlete.AnTHR)
	if config.Athlete.AnTHR == 0 {
		z3Max = maxHR * 0.80
	}
	z4Max := maxHR * 0.95
	summary.Thresholds.ZoneThresholds = []int{int(z1Max), int(z2Max), int(z3Max), int(z4Max)}

	if act.AvgSpeed > 0 {
		summary.GlobalAverages.PaceAvg = formatPace(float64(act.AvgSpeed) / 1000.0)
	}
	if act.MovingTime > 0 && act.ElevationGain.Valid {
		summary.GlobalAverages.Vam = int(float64(act.ElevationGain.Value) / (float64(act.MovingTime) / 3600.0))
	}

	var hrZoneCounters [5]int
	var zoneGAPSum [5]float64
	var zoneHRSum [5]float64
	gradeCounters := map[string]int{"SteepDown": 0, "RunDown": 0, "Flat": 0, "RunUp": 0, "HikeUp": 0}

	currentSegment := TopographicSplit{Type: "Flat", DistKm: 0}
	var segmentStartElev float64
	var segmentStartDist float64
	if ts.Data[0].Altitude.Valid {
		segmentStartElev = ts.Data[0].Altitude.Value
	}
	if len(ts.Data) > 0 && ts.Data[0].Distance.Valid {
		segmentStartDist = float64(ts.Data[0].Distance.Value)
	}

	var upHrSum, upHrSqSum, upGapSum, upCount, upElevGain, upTime, upDist float64
	var downHrSum, downPaceSum, downCount float64

	type DecoupleState struct{ hrSum, speedSum, count float64 }
	halfIndex := float64(act.Distance) / 2
	firstHalf := DecoupleState{}
	secondHalf := DecoupleState{}

	uphillHalfIndex := float64(0) // filled after first pass
	var uphillFirstHalf, uphillSecondHalf DecoupleState

	var totalGAPSpeed float64
	var gapPoints int
	var totalCadence float64
	var cadencePoints int

	// GAP benchmark accumulators
	var gapBandHRSum [9]float64
	var gapBandCount [9]int

	// Z2-only accumulators
	var z2PaceSum, z2GAPSum float64
	var z2Count int

	// VAM by gradient band accumulators
	var vamGain [5]float64
	var vamTime [5]float64
	var vamHRSum [5]float64
	var vamCount [5]int

	// Downhill pace by grade band
	var dhPaceSum [3]float64
	var dhBandHRSum [3]float64
	var dhBandCount [3]int

	// Z4+ consecutive block tracking
	var currentZ4BlockSec int
	var longestZ4BlockSec int
	var z4z5GAPSum float64
	var z4z5Count int

	// Last 10% accumulators
	last10StartDist := float64(act.Distance) * 0.9
	var last10HrSum, last10GapSum, last10PaceSum float64
	var last10Count int

	// Recovery: HR at end and 60s before
	var hrAtEnd, hr60sBeforeEnd int
	var maxEffortHR int
	var maxEffortIdx int

	// Hike/run transition tracking
	var maxRunGrade float64
	var hikeTransitionFound bool

	// Flats economy accumulators
	var flatHrSum, flatPaceSum float64
	var flatCount int

	windowSize := config.GradeSmoothingWindow

	for i := windowSize; i < len(ts.Data); i++ {
		curr := ts.Data[i]
		prevWindow := ts.Data[i-windowSize]

		if !curr.Distance.Valid || !prevWindow.Distance.Valid || !curr.Altitude.Valid || !prevWindow.Altitude.Valid {
			continue
		}

		deltaDist := float64(curr.Distance.Value - prevWindow.Distance.Value)
		deltaElev := curr.Altitude.Value - prevWindow.Altitude.Value

		actualSpeed := 0.0
		timeDelta := 0.0
		if deltaDist > 0 {
			timeDelta = float64(curr.Offset - prevWindow.Offset)
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

		if curr.Cadence.Valid {
			totalCadence += float64(curr.Cadence.Value)
			cadencePoints++
		}

		if actualSpeed > config.MinMovingSpeedMS {
			totalGAPSpeed += gapSpeed
			gapPoints++

			if curr.HeartRate.Valid {
				hr := float64(curr.HeartRate.Value)

				if float64(curr.Distance.Value) < halfIndex {
					firstHalf.hrSum += hr
					firstHalf.speedSum += gapSpeed
					firstHalf.count++
				} else {
					secondHalf.hrSum += hr
					secondHalf.speedSum += gapSpeed
					secondHalf.count++
				}
			}

			// GAP benchmark accumulation
			matched := false
			for b := 0; b < len(gapBenchmarkThresholds)-1; b++ {
				if gapSpeed >= gapBenchmarkThresholds[b].speed {
					if curr.HeartRate.Valid {
						gapBandHRSum[b] += float64(curr.HeartRate.Value)
						gapBandCount[b]++
					}
					matched = true
					break
				}
			}
			if !matched && curr.HeartRate.Valid {
				last := len(gapBenchmarkThresholds) - 1
				gapBandHRSum[last] += float64(curr.HeartRate.Value)
				gapBandCount[last]++
			}

			// Z2-only
			if curr.HeartRate.Valid {
				hr := float64(curr.HeartRate.Value)
				if hr >= z1Max && hr < z2Max {
					z2PaceSum += actualSpeed
					z2GAPSum += gapSpeed
					z2Count++
				}
			}

			// Z4+Z5 threshold accumulation
			if curr.HeartRate.Valid {
				if float64(curr.HeartRate.Value) >= z3Max {
					currentZ4BlockSec += int(timeDelta)
					z4z5GAPSum += gapSpeed
					z4z5Count++
				} else {
					if currentZ4BlockSec > longestZ4BlockSec {
						longestZ4BlockSec = currentZ4BlockSec
					}
					currentZ4BlockSec = 0
				}
			}

			// Last 10%
			if float64(curr.Distance.Value) >= last10StartDist {
				if curr.HeartRate.Valid {
					last10HrSum += float64(curr.HeartRate.Value)
				}
				last10GapSum += gapSpeed
				last10PaceSum += actualSpeed
				last10Count++
			}
		}

		// Grade distribution
		if actualSpeed > config.MinMovingSpeedMS {
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
			zoneIdx := 4
			switch {
			case hr < z1Max:
				zoneIdx = 0
			case hr < z2Max:
				zoneIdx = 1
			case hr < z3Max:
				zoneIdx = 2
			case hr < z4Max:
				zoneIdx = 3
			}
			hrZoneCounters[zoneIdx]++
			if actualSpeed > config.MinMovingSpeedMS {
				zoneGAPSum[zoneIdx] += gapSpeed
				zoneHRSum[zoneIdx] += hr
			}

			// Max effort tracking
			if int(hr) > maxEffortHR {
				maxEffortHR = int(hr)
				maxEffortIdx = i
			}

			// Terrain specific performance
			if gradePct > config.GradeUpThreshold && actualSpeed > config.MinMovingSpeedMS {
				upHrSum += hr
				upHrSqSum += hr * hr
				upGapSum += gapSpeed
				upCount++
				upElevGain += deltaElev
				upDist += deltaDist
				upTime += timeDelta

				for b := range vamBands {
					if gradePct > vamBands[b].gradeLow && gradePct <= vamBands[b].gradeHigh {
						vamGain[b] += deltaElev
						vamTime[b] += timeDelta
						vamHRSum[b] += hr
						vamCount[b]++
						break
					}
				}

				// Hike/run transition: track max grade where still running (cadence-based)
				if curr.Cadence.Valid && !hikeTransitionFound {
					if float64(curr.Cadence.Value) >= 70 {
						if gradePct > maxRunGrade {
							maxRunGrade = gradePct
						}
					} else if gradePct > config.GradeUpThreshold {
						hikeTransitionFound = true
					}
				}

			} else if gradePct < config.GradeDownThreshold && actualSpeed > config.MinMovingSpeedMS {
				downHrSum += hr
				downPaceSum += actualSpeed
				downCount++

				for b := range dhBands {
					if gradePct >= dhBands[b].gradeLow && gradePct < dhBands[b].gradeHigh {
						dhPaceSum[b] += actualSpeed
						dhBandHRSum[b] += hr
						dhBandCount[b]++
						break
					}
				}
			}

			// Flats economy (|grade| < 2%)
			if gradePct >= config.GradeDownThreshold && gradePct <= config.GradeUpThreshold && actualSpeed > config.MinMovingSpeedMS {
				flatHrSum += hr
				flatPaceSum += actualSpeed
				flatCount++
			}
		}

		// Topographic Segmentation
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

		if actualSpeed > config.MinMovingSpeedMS {
			currentSegment.sumGAPSpeed += gapSpeed
			if curr.HeartRate.Valid {
				currentSegment.sumHR += int(curr.HeartRate.Value)
			}
			currentSegment.pointsCount++
		}

		// Capture HR at end and 60s before end
		if i == len(ts.Data)-1 && curr.HeartRate.Valid {
			hrAtEnd = int(curr.HeartRate.Value)
		}
		if i >= len(ts.Data)-61 && i <= len(ts.Data)-60 && curr.HeartRate.Valid {
			hr60sBeforeEnd = int(curr.HeartRate.Value)
		}
	}

	if len(ts.Data) > 0 {
		finalizeSegment(&summary.TopoSplits, &currentSegment, segmentStartElev, segmentStartDist, ts.Data[len(ts.Data)-1], config)
	}

	// Capture Z4+ block that continued to end
	if currentZ4BlockSec > longestZ4BlockSec {
		longestZ4BlockSec = currentZ4BlockSec
	}

	if upDist > 0 {
		uphillHalfIndex = upDist / 2
		var cumulativeUpDist float64
		uphillFirstHalf = DecoupleState{}
		uphillSecondHalf = DecoupleState{}
		for i := windowSize; i < len(ts.Data); i++ {
			curr := ts.Data[i]
			prevWindow := ts.Data[i-windowSize]
			if !curr.Distance.Valid || !prevWindow.Distance.Valid || !curr.Altitude.Valid || !prevWindow.Altitude.Valid {
				continue
			}
			deltaDist := float64(curr.Distance.Value - prevWindow.Distance.Value)
			deltaElev := curr.Altitude.Value - prevWindow.Altitude.Value
			if deltaDist <= config.MinGradeDeltaM {
				continue
			}
			gradePct := (deltaElev / deltaDist) * 100.0
			if gradePct <= config.GradeUpThreshold {
				continue
			}
			actualSpeed := 0.0
			timeDelta := float64(curr.Offset - prevWindow.Offset)
			if deltaDist > 0 && timeDelta > 0 {
				actualSpeed = deltaDist / timeDelta
			}
			if actualSpeed <= config.MinMovingSpeedMS || !curr.HeartRate.Valid {
				continue
			}
			gapSpeed := calculateGAP(actualSpeed, deltaElev/deltaDist)
			cumulativeUpDist += deltaDist
			if cumulativeUpDist <= uphillHalfIndex {
				uphillFirstHalf.hrSum += float64(curr.HeartRate.Value)
				uphillFirstHalf.speedSum += gapSpeed
				uphillFirstHalf.count++
			} else {
				uphillSecondHalf.hrSum += float64(curr.HeartRate.Value)
				uphillSecondHalf.speedSum += gapSpeed
				uphillSecondHalf.count++
			}
		}
	}

	// ===== Post-processing calculations =====

	if gapPoints > 0 {
		summary.GlobalAverages.GAPAvg = formatPace(totalGAPSpeed / float64(gapPoints))
	}
	if cadencePoints > 0 {
		summary.GlobalAverages.CadenceAvg = int(totalCadence/float64(cadencePoints) + 0.5)
	}

	// Distributions
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

	// Zone Details
	for z := range 5 {
		zd := ZoneDetail{Zone: z + 1, PctTime: summary.Distributions.HRZones[z]}
		if hrZoneCounters[z] > 0 {
			if gapPoints > 0 && zoneGAPSum[z] > 0 {
				zd.AvgGAP = formatPace(zoneGAPSum[z] / float64(hrZoneCounters[z]))
			}
			zd.AvgHR = int(zoneHRSum[z] / float64(hrZoneCounters[z]))
		}
		summary.ZoneDetails = append(summary.ZoneDetails, zd)
	}

	// Z2-only metrics
	if z2Count > 0 {
		summary.Distributions.Z2AvgPace = formatPace(z2PaceSum / float64(z2Count))
		summary.Distributions.Z2AvgGAP = formatPace(z2GAPSum / float64(z2Count))
	}

	// GAP Benchmarks
	for b := range gapBenchmarkThresholds {
		if gapBandCount[b] > 0 {
			summary.GAPBenchmarks = append(summary.GAPBenchmarks, GAPBenchmark{
				Range: gapBenchmarkThresholds[b].label,
				AvgHR: int(gapBandHRSum[b] / float64(gapBandCount[b])),
				Count: gapBandCount[b],
			})
		}
	}

	// Terrain Performance
	if upCount > 0 {
		summary.TerrainStats.UphillAvgHR = int(upHrSum / upCount)
		summary.TerrainStats.UphillAvgGAP = formatPace(upGapSum / upCount)
		meanHR := upHrSum / upCount
		summary.TerrainStats.UphillHRStdDev = math.Sqrt(upHrSqSum/upCount - meanHR*meanHR)
		if upTime > 0 {
			summary.TerrainStats.UphillVAM = int((upElevGain / upTime) * 3600.0)
		}
	}
	if downCount > 0 {
		summary.TerrainStats.DownhillAvgHR = int(downHrSum / downCount)
		summary.TerrainStats.DownhillAvgPace = formatPace(downPaceSum / downCount)
		avgDownPace := downPaceSum / downCount
		avgDownHR := downHrSum / downCount
		if avgDownHR > 0 {
			summary.TerrainStats.DownhillEfficiency = math.Round(avgDownPace/avgDownHR*1000) / 1000
		}
	}

	// VAM by gradient band
	totalUpTime := 0
	var vamOutputBandIdx []int
	for b := range vamBands {
		if vamTime[b] > 0 {
			vam := int((vamGain[b] / vamTime[b]) * 3600.0)
			avgHR := 0
			if vamCount[b] > 0 {
				avgHR = int(vamHRSum[b] / float64(vamCount[b]))
			}
			summary.TerrainStats.VAMByGradient = append(summary.TerrainStats.VAMByGradient, GradientVAM{
				Band:  vamBands[b].name,
				VAM:   vam,
				AvgHR: avgHR,
			})
			vamOutputBandIdx = append(vamOutputBandIdx, b)
			totalUpTime += int(vamTime[b])
		}
	}
	for i := range summary.TerrainStats.VAMByGradient {
		b := vamOutputBandIdx[i]
		if totalUpTime > 0 && vamTime[b] > 0 {
			summary.TerrainStats.VAMByGradient[i].PctTime = (int(vamTime[b]) * 100) / totalUpTime
		}
	}

	// Downhill pace by grade
	for b := range dhBands {
		if dhBandCount[b] > 0 {
			avgPace := formatPace(dhPaceSum[b] / float64(dhBandCount[b]))
			avgHR := int(dhBandHRSum[b] / float64(dhBandCount[b]))
			summary.TerrainStats.DownhillPaceByGrade = append(summary.TerrainStats.DownhillPaceByGrade, GradePace{
				Band:    dhBands[b].name,
				AvgPace: avgPace,
				AvgHR:   avgHR,
			})
		}
	}

	// Hike/run transition
	if hikeTransitionFound {
		summary.TerrainStats.HikeRunTransitionGradePct = math.Round(maxRunGrade*10) / 10
	}

	// Decoupling
	if firstHalf.count > 0 && secondHalf.count > 0 {
		ef1 := (firstHalf.speedSum / firstHalf.count) / (firstHalf.hrSum / firstHalf.count)
		ef2 := (secondHalf.speedSum / secondHalf.count) / (secondHalf.hrSum / secondHalf.count)
		if ef1 > 0 {
			decPct := math.Round(((ef1-ef2)/ef1)*100.0*10) / 10
			summary.Decoupling.AerobicDecouplingPct = decPct
		}
		summary.Decoupling.FirstHalfAvgHR = math.Round(firstHalf.hrSum/firstHalf.count*100) / 100
		summary.Decoupling.FirstHalfAvgGAP = formatPace(firstHalf.speedSum / firstHalf.count)
		summary.Decoupling.SecondHalfAvgHR = math.Round(secondHalf.hrSum/secondHalf.count*100) / 100
		summary.Decoupling.SecondHalfAvgGAP = formatPace(secondHalf.speedSum / secondHalf.count)
	}

	if uphillFirstHalf.count > 0 && uphillSecondHalf.count > 0 {
		uef1 := (uphillFirstHalf.speedSum / uphillFirstHalf.count) / (uphillFirstHalf.hrSum / uphillFirstHalf.count)
		uef2 := (uphillSecondHalf.speedSum / uphillSecondHalf.count) / (uphillSecondHalf.hrSum / uphillSecondHalf.count)
		if uef1 > 0 {
			summary.Decoupling.UphillDecouplingPct = math.Round(((uef1-uef2)/uef1)*100.0*10) / 10
		}
	}

	// Thresholds
	if z4z5Count > 0 {
		summary.Thresholds.Z4Z5AvgGAP = formatPace(z4z5GAPSum / float64(z4z5Count))
	}
	summary.Thresholds.LongestZ4BlockSec = longestZ4BlockSec

	// End of run
	if last10Count > 0 {
		if last10HrSum > 0 {
			summary.EndOfRun.Last10PctAvgHR = int(last10HrSum / float64(last10Count))
		}
		summary.EndOfRun.Last10PctAvgGAP = formatPace(last10GapSum / float64(last10Count))
		summary.EndOfRun.Last10PctAvgPace = formatPace(last10PaceSum / float64(last10Count))
	}

	// Recovery
	if hrAtEnd > 0 && hr60sBeforeEnd > 0 {
		summary.Recovery.EndHRDrop60s = hr60sBeforeEnd - hrAtEnd
	}
	if maxEffortIdx > 0 && maxEffortIdx+60 < len(ts.Data) {
		postHR := 0
		postCount := 0
		for j := maxEffortIdx + 1; j < len(ts.Data) && j <= maxEffortIdx+60; j++ {
			if ts.Data[j].HeartRate.Valid {
				postHR += int(ts.Data[j].HeartRate.Value)
				postCount++
			}
		}
		if postCount > 0 {
			summary.Recovery.PostMaxEffortHRDrop60s = maxEffortHR - (postHR / postCount)
		}
	}

	// Economy
	if flatHrSum > 0 && flatPaceSum > 0 && flatCount > 0 {
		avgFlatPace := flatPaceSum / float64(flatCount)
		avgFlatHR := flatHrSum / float64(flatCount)
		if avgFlatPace > 0 {
			summary.Economy.HRPerPaceFlat = math.Round(avgFlatHR/avgFlatPace*100) / 100
		}
	}
	if upCount > 0 && upTime > 0 && upHrSum > 0 {
		upVAM := (upElevGain / upTime) * 3600.0
		avgUpHR := upHrSum / upCount
		if upVAM > 0 {
			summary.Economy.HRPerVAMUphill = math.Round(avgUpHR/upVAM*100) / 100
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

	if seg.GradePct > config.GradeUpThreshold {
		seg.Type = "Uphill"
	} else if seg.GradePct < config.GradeDownThreshold {
		seg.Type = "Downhill"
	} else {
		seg.Type = "Flat"
	}

	if seg.pointsCount > 0 {
		seg.AvgGAP = formatPace(seg.sumGAPSpeed / float64(seg.pointsCount))
		seg.AvgHR = seg.sumHR / seg.pointsCount
	}

	*splits = append(*splits, *seg)
}
