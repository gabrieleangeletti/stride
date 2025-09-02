package stride

import (
	"time"
)

type EnduranceOutdoorActivity struct {
	ExternalID     string    `json:"externalId"`     // The id of the activity in the provider's system
	UserExternalID string    `json:"userExternalId"` // The id of the user in the provider's system
	Provider       Provider  `json:"provider"`
	Sport          Sport     `json:"sport"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
	IanaTimezone   string    `json:"ianaTimezone"`
	UTCOffset      int       `json:"utcOffset"`   // seconds
	ElapsedTime    int       `json:"elapsedTime"` // seconds
	MovingTime     int       `json:"movingTime"`  // seconds
	Distance       float64   `json:"distance"`    // meters
	ElevGain       *float64  `json:"elevGain"`    // meters
	ElevLoss       *float64  `json:"elevLoss"`    // meters
	AvgSpeed       float64   `json:"avgSpeed"`    // meters / second
	AvgHR          *int16    `json:"avgHR"`       // beats / minute
	MaxHR          *int16    `json:"maxHR"`       // beats / minute
}
