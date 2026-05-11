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
	AvgGAP  string `json:"avgGap,omitempty"`
	AvgHR   int    `json:"avgHr,omitempty"`
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
	Z4Z5AvgGAP        string `json:"z4z5AvgGap,omitempty"`
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

	if act.Distance == 0 || len(ts.Data) == 0 || !ts.Data[len(ts.Data)-1].Distance.Valid {
		AugmentGPXData(act, ts, AugmentConfig{ElevationHysteresisM: config.ElevationHysteresisM})
	}

	summary := &LLMRunSummary{Athlete: config.Athlete}

	// 1. Metadata & Global Averages
	summary.Metadata.DistanceKm = float64(act.Distance) / 1000.0
	summary.Metadata.MovingTimeMin = int(act.MovingTime) / 60
	if act.ElevationGain.Valid {
		summary.Metadata.TotalAscentM = int(act.ElevationGain.Value)
		if act.MovingTime > 0 {
			summary.GlobalAverages.Vam = int(float64(act.ElevationGain.Value) / (float64(act.MovingTime) / 3600.0))
		}
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

	// Thresholds setup
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

	// 2. Pre-process into EnrichedPoints
	var enriched []EnrichedPoint
	windowSize := config.GradeSmoothingWindow
	totalUpDist := 0.0

	for i := windowSize; i < len(ts.Data); i++ {
		curr := &ts.Data[i]
		prevWindow := ts.Data[i-windowSize]
		prev := ts.Data[i-1]

		if !curr.Distance.Valid || !prevWindow.Distance.Valid || !curr.Altitude.Valid || !prevWindow.Altitude.Valid {
			continue
		}

		timeDelta := float64(curr.Offset - prev.Offset)
		if timeDelta <= 0 {
			continue // Skip duplicate timestamps
		}

		pointDist := float64(curr.Distance.Value - prev.Distance.Value)
		actualSpeed := pointDist / timeDelta

		deltaDist := float64(curr.Distance.Value - prevWindow.Distance.Value)
		deltaElev := curr.Altitude.Value - prevWindow.Altitude.Value

		gradePct := 0.0
		gradeFraction := 0.0
		if deltaDist > config.MinGradeDeltaM {
			gradeFraction = deltaElev / deltaDist
			gradePct = gradeFraction * 100.0
		}

		gapSpeed := calculateGAP(actualSpeed, gradeFraction)

		enriched = append(enriched, EnrichedPoint{
			Entry:       curr,
			TimeDelta:   timeDelta,
			ActualSpeed: actualSpeed,
			GAPSpeed:    gapSpeed,
			GradePct:    gradePct,
			DistanceM:   float64(curr.Distance.Value),
			DeltaElevM:  curr.Altitude.Value - prev.Altitude.Value, // Point-to-point elevate for VAM
		})

		if gradePct > config.GradeUpThreshold && actualSpeed > config.MinMovingSpeedMS {
			totalUpDist += pointDist
		}
	}

	summary.TopoSplits = DetectTopographicSplits(enriched, config)

	// 3. Time-Weighted Aggregations
	var totalGAPSpeed, totalCadence WeightedAvg
	var hrZoneTimes [5]float64
	var zoneGAP [5]WeightedAvg
	var zoneHR [5]WeightedAvg
	gradeTimes := map[string]float64{"SteepDown": 0, "RunDown": 0, "Flat": 0, "RunUp": 0, "HikeUp": 0}
	totalMovingTime := 0.0

	// GAP Benchmark Accumulators
	var gapBandHR [9]WeightedAvg

	// VAM and DH Bands
	var vamGain [5]float64
	var vamTime [5]float64
	var vamHR [5]WeightedAvg
	var dhPace [3]WeightedAvg
	var dhHR [3]WeightedAvg

	// Decoupling & Economy
	var firstHalfSpeed, firstHalfHR, secondHalfSpeed, secondHalfHR WeightedAvg
	var up1Speed, up1HR, up2Speed, up2HR WeightedAvg
	var flatHR, flatPace WeightedAvg
	var upHR, upGAP, downHR, downPace WeightedAvg
	var z2Pace, z2GAP WeightedAvg

	var last10HR, last10GAP, last10Pace WeightedAvg
	var z4z5GAP WeightedAvg

	halfIndex := float64(act.Distance) / 2
	uphillHalfIndex := totalUpDist / 2
	last10StartDist := float64(act.Distance) * 0.9

	cumulativeUpDist := 0.0
	maxRunGrade := 0.0
	hikeTransitionFound := false

	var currentZ4BlockSec float64
	longestZ4BlockSec := 0.0
	maxEffortHR := 0
	maxEffortOffset := 0
	totalUpElevGain := 0.0
	totalUpTime := 0.0
	smoothedRunGrade := 0.0

	for _, pt := range enriched {
		if pt.ActualSpeed <= config.MinMovingSpeedMS {
			continue
		}

		totalMovingTime += pt.TimeDelta
		totalGAPSpeed.Add(pt.GAPSpeed, pt.TimeDelta)

		if pt.Entry.Cadence.Valid {
			totalCadence.Add(float64(pt.Entry.Cadence.Value), pt.TimeDelta)
		}

		// Grade Distributions
		switch {
		case pt.GradePct < config.GradeSteepDownBelow:
			gradeTimes["SteepDown"] += pt.TimeDelta
		case pt.GradePct >= config.GradeSteepDownBelow && pt.GradePct < config.GradeDownThreshold:
			gradeTimes["RunDown"] += pt.TimeDelta
		case pt.GradePct >= config.GradeDownThreshold && pt.GradePct <= config.GradeUpThreshold:
			gradeTimes["Flat"] += pt.TimeDelta
		case pt.GradePct > config.GradeUpThreshold && pt.GradePct <= config.GradeHikeUpAbove:
			gradeTimes["RunUp"] += pt.TimeDelta
		default:
			gradeTimes["HikeUp"] += pt.TimeDelta
		}

		// Terrain specific (Uphill / Downhill)
		if pt.GradePct > config.GradeUpThreshold {
			cumulativeUpDist += pt.ActualSpeed * pt.TimeDelta
			totalUpElevGain += pt.DeltaElevM
			totalUpTime += pt.TimeDelta

			// VAM Bands
			for b := range vamBands {
				if pt.GradePct > vamBands[b].gradeLow && pt.GradePct <= vamBands[b].gradeHigh {
					vamGain[b] += pt.DeltaElevM
					vamTime[b] += pt.TimeDelta
					if pt.Entry.HeartRate.Valid {
						vamHR[b].Add(float64(pt.Entry.HeartRate.Value), pt.TimeDelta)
					}
					break
				}
			}

			// Hike/Run Transition
			if pt.Entry.Cadence.Valid && !hikeTransitionFound {
				cad := float64(pt.Entry.Cadence.Value)
				// True if running RPM (70-100) OR running SPM (140-200+)
				isRunning := (cad >= 70 && cad <= 100) || (cad >= 140)

				if isRunning {
					smoothedRunGrade = (pt.GradePct * 0.1) + (smoothedRunGrade * 0.9)
					if smoothedRunGrade > maxRunGrade {
						maxRunGrade = smoothedRunGrade
					}
				} else {
					hikeTransitionFound = true
				}
			}
		} else if pt.GradePct < config.GradeDownThreshold {
			// Downhill Bands
			for b := range dhBands {
				if pt.GradePct >= dhBands[b].gradeLow && pt.GradePct < dhBands[b].gradeHigh {
					dhPace[b].Add(pt.ActualSpeed, pt.TimeDelta)
					if pt.Entry.HeartRate.Valid {
						dhHR[b].Add(float64(pt.Entry.HeartRate.Value), pt.TimeDelta)
					}
					break
				}
			}
		}

		// Distance-based Splits (Last 10%)
		if pt.DistanceM >= last10StartDist {
			last10GAP.Add(pt.GAPSpeed, pt.TimeDelta)
			last10Pace.Add(pt.ActualSpeed, pt.TimeDelta)
		}

		// Heart Rate Based Metrics
		if pt.Entry.HeartRate.Valid {
			hr := float64(pt.Entry.HeartRate.Value)

			// Recovery HR tracking
			if int(hr) > maxEffortHR {
				maxEffortHR = int(hr)
				maxEffortOffset = pt.Entry.Offset
			}

			// Zone distributions
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

			hrZoneTimes[zoneIdx] += pt.TimeDelta
			zoneGAP[zoneIdx].Add(pt.GAPSpeed, pt.TimeDelta)
			zoneHR[zoneIdx].Add(hr, pt.TimeDelta)

			if zoneIdx == 1 { // Z2 stats
				z2Pace.Add(pt.ActualSpeed, pt.TimeDelta)
				z2GAP.Add(pt.GAPSpeed, pt.TimeDelta)
			}

			// Z4+ tracking for Thresholds
			if hr >= z3Max {
				currentZ4BlockSec += pt.TimeDelta
				z4z5GAP.Add(pt.GAPSpeed, pt.TimeDelta)
			} else {
				if currentZ4BlockSec > longestZ4BlockSec {
					longestZ4BlockSec = currentZ4BlockSec
				}
				currentZ4BlockSec = 0
			}

			// Last 10% HR
			if pt.DistanceM >= last10StartDist {
				last10HR.Add(hr, pt.TimeDelta)
			}

			// Decoupling (Global)
			if pt.DistanceM < halfIndex {
				firstHalfSpeed.Add(pt.GAPSpeed, pt.TimeDelta)
				firstHalfHR.Add(hr, pt.TimeDelta)
			} else {
				secondHalfSpeed.Add(pt.GAPSpeed, pt.TimeDelta)
				secondHalfHR.Add(hr, pt.TimeDelta)
			}

			// Decoupling (Uphill) & Terrain Stats
			if pt.GradePct > config.GradeUpThreshold {
				upHR.Add(hr, pt.TimeDelta)
				upGAP.Add(pt.GAPSpeed, pt.TimeDelta)

				if cumulativeUpDist <= uphillHalfIndex {
					up1Speed.Add(pt.GAPSpeed, pt.TimeDelta)
					up1HR.Add(hr, pt.TimeDelta)
				} else {
					up2Speed.Add(pt.GAPSpeed, pt.TimeDelta)
					up2HR.Add(hr, pt.TimeDelta)
				}
			} else if pt.GradePct < config.GradeDownThreshold {
				downHR.Add(hr, pt.TimeDelta)
				downPace.Add(pt.ActualSpeed, pt.TimeDelta)
			} else {
				// Flats Economy (|grade| <= 2%)
				flatHR.Add(hr, pt.TimeDelta)
				flatPace.Add(pt.ActualSpeed, pt.TimeDelta)
			}

			// GAP Benchmarks
			matched := false
			for b := 0; b < len(gapBenchmarkThresholds)-1; b++ {
				if pt.GAPSpeed >= gapBenchmarkThresholds[b].speed {
					gapBandHR[b].Add(hr, pt.TimeDelta)
					matched = true
					break
				}
			}
			if !matched {
				gapBandHR[len(gapBenchmarkThresholds)-1].Add(hr, pt.TimeDelta)
			}
		}
	}

	// Finalize longest Z4 block
	if currentZ4BlockSec > longestZ4BlockSec {
		longestZ4BlockSec = currentZ4BlockSec
	}

	// 4. Map everything to the JSON Struct
	summary.GlobalAverages.GAPAvg = formatPace(totalGAPSpeed.Avg())
	summary.GlobalAverages.CadenceAvg = int(math.Round(totalCadence.Avg()))

	if totalMovingTime > 0 {
		summary.Distributions.GradeSteepDownPct = int(math.Round((gradeTimes["SteepDown"] / totalMovingTime) * 100))
		summary.Distributions.GradeRunDownPct = int(math.Round((gradeTimes["RunDown"] / totalMovingTime) * 100))
		summary.Distributions.GradeFlatPct = int(math.Round((gradeTimes["Flat"] / totalMovingTime) * 100))
		summary.Distributions.GradeRunUpPct = int(math.Round((gradeTimes["RunUp"] / totalMovingTime) * 100))
		summary.Distributions.GradeHikeUpPct = int(math.Round((gradeTimes["HikeUp"] / totalMovingTime) * 100))

		hrTimeTotal := hrZoneTimes[0] + hrZoneTimes[1] + hrZoneTimes[2] + hrZoneTimes[3] + hrZoneTimes[4]
		if hrTimeTotal > 0 {
			summary.Distributions.HRZones = []int{
				int(math.Round((hrZoneTimes[0] / hrTimeTotal) * 100)),
				int(math.Round((hrZoneTimes[1] / hrTimeTotal) * 100)),
				int(math.Round((hrZoneTimes[2] / hrTimeTotal) * 100)),
				int(math.Round((hrZoneTimes[3] / hrTimeTotal) * 100)),
				int(math.Round((hrZoneTimes[4] / hrTimeTotal) * 100)),
			}
		}
	}

	for z := 0; z < 5; z++ {
		pct := 0
		if len(summary.Distributions.HRZones) == 5 {
			pct = summary.Distributions.HRZones[z]
		}
		summary.ZoneDetails = append(summary.ZoneDetails, ZoneDetail{
			Zone:    z + 1,
			PctTime: pct,
			AvgGAP:  formatPace(zoneGAP[z].Avg()),
			AvgHR:   int(math.Round(zoneHR[z].Avg())),
		})
	}

	summary.Distributions.Z2AvgPace = formatPace(z2Pace.Avg())
	summary.Distributions.Z2AvgGAP = formatPace(z2GAP.Avg())

	// Terrain Performance
	summary.TerrainStats.UphillAvgHR = int(math.Round(upHR.Avg()))
	summary.TerrainStats.UphillAvgGAP = formatPace(upGAP.Avg())
	summary.TerrainStats.UphillHRStdDev = upHR.StdDev()
	if totalUpTime > 0 {
		summary.TerrainStats.UphillVAM = int(math.Round((totalUpElevGain / totalUpTime) * 3600.0))
	}

	summary.TerrainStats.DownhillAvgHR = int(math.Round(downHR.Avg()))
	summary.TerrainStats.DownhillAvgPace = formatPace(downPace.Avg())
	if summary.TerrainStats.DownhillAvgHR > 0 {
		summary.TerrainStats.DownhillEfficiency = math.Round((downPace.Avg()/float64(summary.TerrainStats.DownhillAvgHR))*1000) / 1000
	}
	if hikeTransitionFound {
		summary.TerrainStats.HikeRunTransitionGradePct = math.Round(maxRunGrade*10) / 10
	}

	// VAM & DH Bands output
	for b := range vamBands {
		if vamTime[b] > 0 {
			summary.TerrainStats.VAMByGradient = append(summary.TerrainStats.VAMByGradient, GradientVAM{
				Band:    vamBands[b].name,
				VAM:     int(math.Round((vamGain[b] / vamTime[b]) * 3600.0)),
				AvgHR:   int(math.Round(vamHR[b].Avg())),
				PctTime: int(math.Round((vamTime[b] / totalUpTime) * 100)),
			})
		}
	}
	for b := range dhBands {
		if dhPace[b].Count > 0 {
			summary.TerrainStats.DownhillPaceByGrade = append(summary.TerrainStats.DownhillPaceByGrade, GradePace{
				Band:    dhBands[b].name,
				AvgPace: formatPace(dhPace[b].Avg()),
				AvgHR:   int(math.Round(dhHR[b].Avg())),
			})
		}
	}

	// GAP Benchmarks
	for b := range gapBenchmarkThresholds {
		if gapBandHR[b].Count > 0 {
			summary.GAPBenchmarks = append(summary.GAPBenchmarks, GAPBenchmark{
				Range: gapBenchmarkThresholds[b].label,
				AvgHR: int(math.Round(gapBandHR[b].Avg())),
				Count: gapBandHR[b].Count,
			})
		}
	}

	// Decoupling
	if firstHalfHR.Avg() > 0 && secondHalfHR.Avg() > 0 {
		ef1 := firstHalfSpeed.Avg() / firstHalfHR.Avg()
		ef2 := secondHalfSpeed.Avg() / secondHalfHR.Avg()
		if ef1 > 0 {
			summary.Decoupling.AerobicDecouplingPct = math.Round(((ef1-ef2)/ef1)*100.0*10) / 10
		}
		summary.Decoupling.FirstHalfAvgHR = math.Round(firstHalfHR.Avg()*100) / 100
		summary.Decoupling.FirstHalfAvgGAP = formatPace(firstHalfSpeed.Avg())
		summary.Decoupling.SecondHalfAvgHR = math.Round(secondHalfHR.Avg()*100) / 100
		summary.Decoupling.SecondHalfAvgGAP = formatPace(secondHalfSpeed.Avg())
	}

	if up1HR.Avg() > 0 && up2HR.Avg() > 0 {
		uef1 := up1Speed.Avg() / up1HR.Avg()
		uef2 := up2Speed.Avg() / up2HR.Avg()
		if uef1 > 0 {
			summary.Decoupling.UphillDecouplingPct = math.Round(((uef1-uef2)/uef1)*100.0*10) / 10
		}
	}

	// Thresholds
	summary.Thresholds.Z4Z5AvgGAP = formatPace(z4z5GAP.Avg())
	summary.Thresholds.LongestZ4BlockSec = int(math.Round(longestZ4BlockSec))

	// End Of Run
	summary.EndOfRun.Last10PctAvgHR = int(math.Round(last10HR.Avg()))
	summary.EndOfRun.Last10PctAvgGAP = formatPace(last10GAP.Avg())
	summary.EndOfRun.Last10PctAvgPace = formatPace(last10Pace.Avg())

	// Recovery
	maxOffsetInFile := ts.MaxOffset()
	hrAtEnd := getHRAtOffset(ts, maxOffsetInFile)
	hr60sBeforeEnd := getHRAtOffset(ts, maxOffsetInFile-60)
	if hrAtEnd > 0 && hr60sBeforeEnd > 0 {
		summary.Recovery.EndHRDrop60s = hr60sBeforeEnd - hrAtEnd
	}
	if maxEffortHR > 0 {
		postHR := getAvgHRInOffsetRange(ts, maxEffortOffset+1, maxEffortOffset+60)
		if postHR > 0 {
			summary.Recovery.PostMaxEffortHRDrop60s = maxEffortHR - postHR
		}
	}

	// Economy
	if flatPace.Avg() > 0 {
		summary.Economy.HRPerPaceFlat = math.Round((flatHR.Avg()/flatPace.Avg())*100) / 100
	}
	if summary.TerrainStats.UphillVAM > 0 {
		summary.Economy.HRPerVAMUphill = math.Round((upHR.Avg()/float64(summary.TerrainStats.UphillVAM))*100) / 100
	}

	return summary, nil
}
