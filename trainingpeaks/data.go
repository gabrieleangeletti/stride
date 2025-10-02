package trainingpeaks

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gabrieleangeletti/stride"
)

type TrainingPeaksDatetimeType time.Time

func (t TrainingPeaksDatetimeType) Time() time.Time {
	return time.Time(t)
}

func (t *TrainingPeaksDatetimeType) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	parsed, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	*t = TrainingPeaksDatetimeType(parsed)
	return nil
}

func (t TrainingPeaksDatetimeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format("2006-01-02T15:04:05"))
}

type TrainingPeaksLastModifiedType time.Time

func (t TrainingPeaksLastModifiedType) Time() time.Time {
	return time.Time(t)
}

func (t *TrainingPeaksLastModifiedType) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	parsed, err := time.Parse("2006-01-02T15:04:05.999", s)
	if err != nil {
		return err
	}
	*t = TrainingPeaksLastModifiedType(parsed)
	return nil
}

func (t TrainingPeaksLastModifiedType) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format("2006-01-02T15:04:05.999"))
}

type TrainingPeaksWorkoutSummary struct {
	WorkoutID                         int64                         `json:"workoutId"`
	AthleteID                         int                           `json:"athleteId"`
	Title                             string                        `json:"title"`
	WorkoutTypeValueID                int                           `json:"workoutTypeValueId"`
	Code                              *string                       `json:"code"`
	WorkoutDay                        TrainingPeaksDatetimeType     `json:"workoutDay"`
	StartTime                         TrainingPeaksDatetimeType     `json:"startTime"`
	StartTimePlanned                  *string                       `json:"startTimePlanned"`
	IsItAnOr                          bool                          `json:"isItAnOr"`
	IsHidden                          *bool                         `json:"isHidden"`
	Completed                         *bool                         `json:"completed"`
	Description                       *string                       `json:"description"`
	UserTags                          string                        `json:"userTags"`
	CoachComments                     *string                       `json:"coachComments"`
	WorkoutComments                   []any                         `json:"workoutComments"`
	NewComment                        *string                       `json:"newComment"`
	PublicSettingValue                int                           `json:"publicSettingValue"`
	SharedWorkoutInformationKey       string                        `json:"sharedWorkoutInformationKey"`
	SharedWorkoutInformationExpireKey string                        `json:"sharedWorkoutInformationExpireKey"`
	Distance                          float64                       `json:"distance"`
	DistancePlanned                   *float64                      `json:"distancePlanned"`
	DistanceCustomized                *float64                      `json:"distanceCustomized"`
	DistanceUnitsCustomized           *string                       `json:"distanceUnitsCustomized"`
	TotalTime                         float64                       `json:"totalTime"` // hours
	TotalTimePlanned                  *float64                      `json:"totalTimePlanned"`
	HeartRateMinimum                  int                           `json:"heartRateMinimum"`
	HeartRateMaximum                  int                           `json:"heartRateMaximum"`
	HeartRateAverage                  int                           `json:"heartRateAverage"`
	Calories                          int                           `json:"calories"`
	CaloriesPlanned                   *int                          `json:"caloriesPlanned"`
	TssActual                         float64                       `json:"tssActual"`
	TssPlanned                        *float64                      `json:"tssPlanned"`
	TssSource                         int                           `json:"tssSource"`
	If                                float64                       `json:"if"`
	IfPlanned                         *float64                      `json:"ifPlanned"`
	VelocityAverage                   float64                       `json:"velocityAverage"`
	VelocityPlanned                   *float64                      `json:"velocityPlanned"`
	VelocityMaximum                   float64                       `json:"velocityMaximum"`
	NormalizedSpeedActual             float64                       `json:"normalizedSpeedActual"`
	NormalizedPowerActual             *float64                      `json:"normalizedPowerActual"`
	PowerAverage                      *float64                      `json:"powerAverage"`
	PowerMaximum                      *float64                      `json:"powerMaximum"`
	Energy                            *float64                      `json:"energy"`
	EnergyPlanned                     *float64                      `json:"energyPlanned"`
	ElevationGain                     float64                       `json:"elevationGain"`
	ElevationGainPlanned              *float64                      `json:"elevationGainPlanned"`
	ElevationLoss                     float64                       `json:"elevationLoss"`
	ElevationMinimum                  float64                       `json:"elevationMinimum"`
	ElevationAverage                  float64                       `json:"elevationAverage"`
	ElevationMaximum                  float64                       `json:"elevationMaximum"`
	TorqueAverage                     *float64                      `json:"torqueAverage"`
	TorqueMaximum                     *float64                      `json:"torqueMaximum"`
	TempMin                           float64                       `json:"tempMin"`
	TempAvg                           float64                       `json:"tempAvg"`
	TempMax                           float64                       `json:"tempMax"`
	CadenceAverage                    int                           `json:"cadenceAverage"`
	CadenceMaximum                    int                           `json:"cadenceMaximum"`
	LastModifiedDate                  TrainingPeaksLastModifiedType `json:"lastModifiedDate"`
	EquipmentBikeID                   *int                          `json:"equipmentBikeId"`
	EquipmentShoeID                   int                           `json:"equipmentShoeId"`
	IsLocked                          *bool                         `json:"isLocked"`
	ComplianceDurationPercent         *float64                      `json:"complianceDurationPercent"`
	ComplianceDistancePercent         *float64                      `json:"complianceDistancePercent"`
	ComplianceTssPercent              *float64                      `json:"complianceTssPercent"`
	RPE                               *int                          `json:"rpe"`
	Feeling                           *int                          `json:"feeling"`
	Structure                         json.RawMessage               `json:"structure"`
	OrderOnDay                        *int                          `json:"orderOnDay"`
	PersonalRecordCount               int                           `json:"personalRecordCount"`
	SyncedTo                          *string                       `json:"syncedTo"`
	PoolLengthOptionID                *int                          `json:"poolLengthOptionId"`
}

func (s TrainingPeaksWorkoutSummary) GetID() string {
	return fmt.Sprintf("%d", s.WorkoutID)
}

func (s TrainingPeaksWorkoutSummary) GetStartTime() time.Time {
	return s.StartTime.Time()
}

func (s TrainingPeaksWorkoutSummary) GetEndTime() time.Time {
	return s.StartTime.Time().Add(time.Duration(s.TotalTime*3600) * time.Second)
}

func (s TrainingPeaksWorkoutSummary) ToActivity() (*stride.Activity, error) {
	return &stride.Activity{
		StartTime:     s.GetStartTime(),
		ElapsedTime:   uint32(s.TotalTime * 3600),
		Distance:      uint32(s.Distance),
		AvgSpeed:      uint16(s.VelocityAverage),
		ElevationGain: stride.Optional[uint16]{Valid: true, Value: uint16(s.ElevationGain)},
		ElevationLoss: stride.Optional[uint16]{Valid: true, Value: uint16(s.ElevationLoss)},
		AvgHR:         stride.Optional[uint8]{Valid: true, Value: uint8(s.HeartRateAverage)},
		MaxHR:         stride.Optional[uint8]{Valid: true, Value: uint8(s.HeartRateMaximum)},
	}, nil
}

type TrainingPeaksWorkoutDetail struct {
	WorkoutID            int64             `json:"workoutId"`
	TotalStats           json.RawMessage   `json:"totalStats"`
	LapStats             []json.RawMessage `json:"lapStats"`
	PeakCadences         []json.RawMessage `json:"peakCadences"`
	PeakHeartRates       []json.RawMessage `json:"peakHeartRates"`
	PeakPowers           []json.RawMessage `json:"peakPowers"`
	PeakSpeeds           []json.RawMessage `json:"peakSpeeds"`
	PeakSpeedsByDistance []json.RawMessage `json:"peakSpeedsByDistance"`
	SampleRate           int               `json:"sampleRate"`
	BoundingBox          json.RawMessage   `json:"boundingBox"`
	FivePointSignature   json.RawMessage   `json:"fivePointSignature"`
	WorkoutSampleList    json.RawMessage   `json:"workoutSampleList"`
	SwimLengthList       json.RawMessage   `json:"swimLengthList"`
}
