package stride

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrMustUseRegistry = errors.New("direct JSON unmarshaling is unsupported, please use ActivityRegistry")
)

type Activity struct {
	ExternalID   string          `json:"externalId"` // The id given by the provider
	Provider     Provider        `json:"provider"`
	Sport        Sport           `json:"sport"`
	StartTime    time.Time       `json:"startTime"`
	EndTime      time.Time       `json:"endTime"`
	IanaTimezone string          `json:"ianaTimezone"`
	UTCOffset    int             `json:"utcOffset"` // seconds
	Metrics      ActivityMetrics `json:"metrics,omitempty"`
}

func (a Activity) MarshalJSON() ([]byte, error) {
	type Alias Activity

	if a.Metrics == nil {
		return json.Marshal(&struct {
			*Alias
			Metrics any `json:"metrics,omitempty"`
		}{Alias: (*Alias)(&a)})
	}

	return json.Marshal(&struct {
		*Alias
		Metrics ActivityMetrics `json:"metrics"`
	}{Alias: (*Alias)(&a), Metrics: a.Metrics})
}

func (a *Activity) UnmarshalJSON(b []byte) error {
	return ErrMustUseRegistry
}

type ActivityMetrics interface {
	GetKind() string
	Validate() error
}

type MetricsBase struct {
	Kind string `json:"kind"` // discriminator, e.g. "endurance", "strength"
}

func (m MetricsBase) GetKind() string { return m.Kind }
func (m MetricsBase) Validate() error { return nil }

// ---- Endurance (running/cycling/swimming/etc.)
type EnduranceMetrics struct {
	MetricsBase
	ElapsedTime int     `json:"elapsedTime"`        // seconds
	MovingTime  int     `json:"movingTime"`         // seconds
	Distance    float64 `json:"distance"`           // meters
	ElevGain    float64 `json:"elevGain,omitempty"` // meters
	ElevLoss    float64 `json:"elevLoss,omitempty"` // meters
	AvgSpeed    float64 `json:"avgSpeed,omitempty"` // meters / second
	AvgHR       *int16  `json:"avgHR,omitempty"`    // beats / minute
	MaxHR       *int16  `json:"maxHR,omitempty"`    // beats / minute
}

func NewEnduranceMetrics() *EnduranceMetrics {
	return &EnduranceMetrics{MetricsBase: MetricsBase{Kind: "endurance"}}
}

// ---- Strength (weightlifting)
type StrengthMetrics struct {
	MetricsBase
	Sets        int     `json:"sets"`
	Reps        int     `json:"reps"`
	AvgWeightKg float64 `json:"avgWeightKg,omitempty"`
}

func NewStrengthMetrics() *StrengthMetrics {
	return &StrengthMetrics{MetricsBase: MetricsBase{Kind: "strength"}}
}

// ---- Climbing
type ClimbingMetrics struct {
	MetricsBase
	Attempts int     `json:"attempts,omitempty"`
	Topouts  int     `json:"topouts,omitempty"`
	Grade    string  `json:"grade,omitempty"`
	HeightM  float64 `json:"heightM,omitempty"`
}

func NewClimbingMetrics() *ClimbingMetrics {
	return &ClimbingMetrics{MetricsBase: MetricsBase{Kind: "climbing"}}
}

// ---- Unknown (forward compatibility)
type UnknownMetrics struct {
	MetricsBase
	Raw json.RawMessage `json:"raw"`
}

func NewUnknown(kind string, raw json.RawMessage) *UnknownMetrics {
	return &UnknownMetrics{MetricsBase: MetricsBase{Kind: kind}, Raw: raw}
}
