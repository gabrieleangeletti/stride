package strava

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gabrieleangeletti/stride"
)

type TokenResponse struct {
	TokenType    string  `json:"token_type"`
	ExpiresAt    int     `json:"expires_at"`
	ExpiresIn    int     `json:"expires_in"`
	RefreshToken string  `json:"refresh_token"`
	AccessToken  string  `json:"access_token"`
	Athlete      Athlete `json:"athlete"`
}

// Resource state, indicates level of detail. Possible values: 2 -> "summary", 3 -> "detail"
type ResourceState int

const (
	ResourceStateSummary ResourceState = 2
	ResourceStateDetail  ResourceState = 3
)

type WebhookRegistrationResponse struct {
	ID int `json:"id"`
}

type WebhookSubscription struct {
	ID            int           `json:"id"`
	ApplicationID int           `json:"application_id"`
	CallbackURL   string        `json:"callback_url"`
	ResourceState ResourceState `json:"resource_state"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type WebhookObjectType string

const (
	WebhookActivity WebhookObjectType = "activity"
	WebhookAthlete  WebhookObjectType = "athlete"
)

type WebhookAspectType string

const (
	WebhookCreate WebhookAspectType = "create"
	WebhookUpdate WebhookAspectType = "update"
	WebhookDelete WebhookAspectType = "delete"
)

type WebhookEvent struct {
	ObjectType     WebhookObjectType `json:"object_type"`
	ObjectID       int               `json:"object_id"` // For activity events, the activity's ID. For athlete events, the athlete's ID.
	AspectType     WebhookAspectType `json:"aspect_type"`
	Updates        map[string]any    `json:"updates"`         // For activity update events, keys can contain "title," "type," and "private," which is always "true" (activity visibility set to Only You) or "false" (activity visibility set to Followers Only or Everyone). For app deauthorization events, there is always an "authorized" : "false" key-value pair.
	OwnerID        int               `json:"owner_id"`        // The athlete's ID.
	SubscriptionID int               `json:"subscription_id"` // The push subscription ID that is receiving this event.
	EventTime      int               `json:"event_time"`      // The time that the event occurred.
}

type Athlete struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	ResourceState int       `json:"resource_state"`
	Firstname     string    `json:"firstname"`
	Lastname      string    `json:"lastname"`
	Bio           string    `json:"bio"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	Sex           string    `json:"sex"`
	Premium       bool      `json:"premium"`
	Summit        bool      `json:"summit"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	BadgeTypeID   int       `json:"badge_type_id"`
	Weight        float64   `json:"weight"`
	ProfileMedium string    `json:"profile_medium"`
	Profile       string    `json:"profile"`
	Friend        string    `json:"friend"`
	Follower      string    `json:"follower"`
}

type ActivitySummary struct {
	ResourceState              int         `json:"resource_state"`
	Athlete                    Athlete     `json:"athlete"`
	Name                       string      `json:"name"`
	Distance                   float64     `json:"distance"`             // in meters
	MovingTime                 int         `json:"moving_time"`          // in seconds
	ElapsedTime                int         `json:"elapsed_time"`         // in seconds
	TotalElevationGain         float64     `json:"total_elevation_gain"` // in meters
	Type                       string      `json:"type"`
	SportType                  string      `json:"sport_type"`
	WorkoutType                int         `json:"workout_type"`
	ID                         int         `json:"id"`
	StartDate                  time.Time   `json:"start_date"`
	StartDateLocal             time.Time   `json:"start_date_local"`
	Timezone                   string      `json:"timezone"`
	UtcOffset                  float64     `json:"utc_offset"` // seconds
	LocationCity               string      `json:"location_city"`
	LocationState              string      `json:"location_state"`
	LocationCountry            string      `json:"location_country"`
	AchievementCount           int         `json:"achievement_count"`
	KudosCount                 int         `json:"kudos_count"`
	CommentCount               int         `json:"comment_count"`
	PhotoCount                 int         `json:"photo_count"`
	Map                        ActivityMap `json:"map"`
	Trainer                    bool        `json:"trainer"`
	Commute                    bool        `json:"commute"`
	Manual                     bool        `json:"manual"`
	Private                    bool        `json:"private"`
	Visibility                 string      `json:"visibility"`
	Flagged                    bool        `json:"flagged"`
	GearID                     string      `json:"gear_id"`
	StartLatLng                []float64   `json:"start_latlng"`
	EndLatLng                  []float64   `json:"end_latlng"`
	AverageSpeed               float64     `json:"average_speed"`   // in meters/second
	MaxSpeed                   float64     `json:"max_speed"`       // in meters/second
	AverageCadence             float64     `json:"average_cadence"` // in rpm
	AverageWatts               float64     `json:"average_watts"`
	MaxWatts                   float64     `json:"max_watts"`
	WeightedAverageWatts       float64     `json:"weighted_average_watts"`
	DeviceWatts                bool        `json:"device_watts"`
	Kilojoules                 float64     `json:"kilojoules"`
	HasHeartrate               bool        `json:"has_heartrate"`
	AverageHeartrate           float64     `json:"average_heartrate"`
	MaxHeartrate               float64     `json:"max_heartrate"`
	HeartrateOptOut            bool        `json:"heartrate_opt_out"`
	DisplayHideHeartrateOption bool        `json:"display_hide_heartrate_option"`
	ElevHigh                   float64     `json:"elev_high"` // in meters
	ElevLow                    float64     `json:"elev_low"`  // in meters
	UploadID                   int         `json:"upload_id"`
	UploadIDStr                string      `json:"upload_id_str"`
	ExternalID                 string      `json:"external_id"`
	FromAcceptedTag            bool        `json:"from_accepted_tag"`
	PrCount                    int         `json:"pr_count"`
	TotalPhotoCount            int         `json:"total_photo_count"`
	HasKudoed                  bool        `json:"has_kudoed"`
	SufferScore                float64     `json:"suffer_score"`
}

// IanaTimezone extracts the IANA timezone from the Strava timezone string.
// Strava stores the timezone as a string in the format "(GMT+00:00) Europe/Lisbon".
// This function returns the IANA timezone identifier, e.g. "Europe/Lisbon".
func (a ActivitySummary) IanaTimezone() string {
	if a.Timezone == "" {
		return ""
	}

	parts := strings.Split(a.Timezone, " ")

	return parts[len(parts)-1]
}

func (a ActivitySummary) ToEnduranceActivity() (*stride.EnduranceOutdoorActivity, error) {
	sport, err := a.mapSportType()
	if err != nil {
		return nil, err
	}

	if !stride.IsEnduranceOutdoorActivity(sport) {
		return nil, stride.ErrActivityIsNotOutdoorEndurance
	}

	activity := &stride.EnduranceOutdoorActivity{
		ExternalID:     strconv.Itoa(a.ID),
		UserExternalID: strconv.Itoa(a.Athlete.ID),
		Provider:       stride.ProviderStrava,
		Sport:          sport,
		StartTime:      a.StartDate,
		EndTime:        a.StartDate.Add(time.Duration(a.ElapsedTime) * time.Second),
		IanaTimezone:   a.IanaTimezone(),
		UTCOffset:      int(a.UtcOffset),
		ElapsedTime:    a.ElapsedTime,
		MovingTime:     a.MovingTime,
		Distance:       a.Distance,
		ElevGain:       &a.TotalElevationGain,
		AvgSpeed:       a.AverageSpeed,
	}

	if a.HasHeartrate && a.AverageHeartrate > 0 {
		avgHR := int16(a.AverageHeartrate)
		activity.AvgHR = &avgHR
	}

	if a.HasHeartrate && a.MaxHeartrate > 0 {
		maxHR := int16(a.MaxHeartrate)
		activity.MaxHR = &maxHR
	}

	return activity, nil
}

func (a ActivitySummary) SummaryPolyline() string {
	return a.Map.SummaryPolyline
}

func (a ActivitySummary) mapSportType() (stride.Sport, error) {
	switch a.SportType {
	case "Run":
		return stride.SportRunning, nil
	case "TrailRun":
		return stride.SportTrailRunning, nil
	case "Ride":
		return stride.SportCycling, nil
	case "Hike":
		return stride.SportHiking, nil
	default:
		return "", fmt.Errorf("%w: %s", stride.ErrUnsupportedSportType, a.SportType)
	}
}

type ActivityMap struct {
	ID              string `json:"id"`
	SummaryPolyline string `json:"summary_polyline"`
	ResourceState   int    `json:"resource_state"`
}

type ActivityStream struct {
	VelocitySmooth StreamSet[float64] `json:"velocity_smooth"`
	Cadence        StreamSet[int]     `json:"cadence"`
	Distance       StreamSet[float64] `json:"distance"`
	Altitude       StreamSet[float64] `json:"altitude"`
	Heartrate      StreamSet[int]     `json:"heartrate"`
	Time           StreamSet[int]     `json:"time"`
}

type StreamSet[T int | float64] struct {
	OriginalSize int    `json:"original_size"` // The number of data points in this stream.
	Resolution   string `json:"resolution"`    // The level of detail (sampling) in which this stream was returned. May take one of the following values: low, medium, high.
	SeriesType   string `json:"series_type"`   // The base series used in the case the stream was downsampled. May take one of the following values: distance, time.
	Data         []T    `json:"data"`          // The sequence of values for this stream.
}
