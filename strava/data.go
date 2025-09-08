package strava

import (
	"fmt"
	"strings"
	"time"

	"github.com/gabrieleangeletti/stride"
)

type LatLng = [2]float64

type SportType string

const (
	SportTypeAlpineSki                     = "AlpineSki"
	SportTypeBackcountrySki                = "BackcountrySki"
	SportTypeBadminton                     = "Badminton"
	SportTypeCanoeing                      = "Canoeing"
	SportTypeCrossfit                      = "Crossfit"
	SportTypeEBikeRide                     = "EBikeRide"
	SportTypeElliptical                    = "Elliptical"
	SportTypeEMountainBikeRide             = "EMountainBikeRide"
	SportTypeGolf                          = "Golf"
	SportTypeGravelRide                    = "GravelRide"
	SportTypeHandcycle                     = "Handcycle"
	SportTypeHighIntensityIntervalTraining = "HighIntensityIntervalTraining"
	SportTypeHike                          = "Hike"
	SportTypeIceSkate                      = "IceSkate"
	SportTypeInlineSkate                   = "InlineSkate"
	SportTypeKayaking                      = "Kayaking"
	SportTypeKitesurf                      = "Kitesurf"
	SportTypeMountainBikeRide              = "MountainBikeRide"
	SportTypeNordicSki                     = "NordicSki"
	SportTypePickleball                    = "Pickleball"
	SportTypePilates                       = "Pilates"
	SportTypeRacquetball                   = "Racquetball"
	SportTypeRide                          = "Ride"
	SportTypeRockClimbing                  = "RockClimbing"
	SportTypeRollerSki                     = "RollerSki"
	SportTypeRowing                        = "Rowing"
	SportTypeRun                           = "Run"
	SportTypeSail                          = "Sail"
	SportTypeSkateboard                    = "Skateboard"
	SportTypeSnowboard                     = "Snowboard"
	SportTypeSnowshoe                      = "Snowshoe"
	SportTypeSoccer                        = "Soccer"
	SportTypeSquash                        = "Squash"
	SportTypeStairStepper                  = "StairStepper"
	SportTypeStandUpPaddling               = "StandUpPaddling"
	SportTypeSurfing                       = "Surfing"
	SportTypeSwim                          = "Swim"
	SportTypeTableTennis                   = "TableTennis"
	SportTypeTennis                        = "Tennis"
	SportTypeTrailRun                      = "TrailRun"
	SportTypeVelomobile                    = "Velomobile"
	SportTypeVirtualRide                   = "VirtualRide"
	SportTypeVirtualRow                    = "VirtualRow"
	SportTypeVirtualRun                    = "VirtualRun"
	SportTypeWalk                          = "Walk"
	SportTypeWeightTraining                = "WeightTraining"
	SportTypeWheelchair                    = "Wheelchair"
	SportTypeWindsurf                      = "Windsurf"
	SportTypeWorkout                       = "Workout"
	SportTypeYoga                          = "Yoga"
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
	ObjectID       int64             `json:"object_id"` // For activity events, the activity's ID. For athlete events, the athlete's ID.
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

type MetaAthlete struct {
	ID int `json:"id"`
}

type ActivitySummary struct {
	ResourceState              int         `json:"resource_state"`
	Athlete                    MetaAthlete `json:"athlete"`
	Name                       string      `json:"name"`
	Distance                   float64     `json:"distance"`             // in meters
	MovingTime                 int         `json:"moving_time"`          // in seconds
	ElapsedTime                int         `json:"elapsed_time"`         // in seconds
	TotalElevationGain         float64     `json:"total_elevation_gain"` // in meters
	Type                       string      `json:"type"`
	SportType                  SportType   `json:"sport_type"`
	WorkoutType                int         `json:"workout_type"`
	ID                         int64       `json:"id"`
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
	StartLatLng                LatLng      `json:"start_latlng"`
	EndLatLng                  LatLng      `json:"end_latlng"`
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
	UploadID                   int64       `json:"upload_id"`
	UploadIDStr                string      `json:"upload_id_str"`
	ExternalID                 string      `json:"external_id"`
	FromAcceptedTag            bool        `json:"from_accepted_tag"`
	PrCount                    int         `json:"pr_count"`
	TotalPhotoCount            int         `json:"total_photo_count"`
	HasKudoed                  bool        `json:"has_kudoed"`
	SufferScore                float64     `json:"suffer_score"`
}

type ActivityDetailed struct {
	ResourceState      int         `json:"resource_state"`
	ID                 int64       `json:"id"`                     // The unique identifier of the activity
	ExternalID         string      `json:"external_id"`            // The identifier provided at upload time
	UploadID           int64       `json:"upload_id"`              // The identifier of the upload that resulted in this activity
	Athlete            MetaAthlete `json:"athlete"`                // The athlete who performed this activity
	Name               string      `json:"name"`                   // The name of the activity
	Distance           float64     `json:"distance"`               // The activity's distance, in meters
	MovingTime         int         `json:"moving_time"`            // The activity's moving time, in seconds
	ElapsedTime        int         `json:"elapsed_time"`           // The activity's elapsed time, in seconds
	TotalElevationGain float64     `json:"total_elevation_gain"`   // The activity's total elevation gain in meters
	ElevHigh           float64     `json:"elev_high"`              // The activity's highest elevation, in meters
	ElevLow            float64     `json:"elev_low"`               // The activity's lowest elevation, in meters
	Type               string      `json:"type"`                   // Deprecated. Prefer to use sport_type
	SportType          SportType   `json:"sport_type"`             // The sport type of this activity
	StartDate          time.Time   `json:"start_date"`             // The time at which the activity was started
	StartDateLocal     time.Time   `json:"start_date_local"`       // The time at which the activity was started in the local timezone
	Timezone           string      `json:"timezone"`               // The timezone of the activity
	StartLatLng        LatLng      `json:"start_latlng"`           // The start location of this activity
	EndLatLng          LatLng      `json:"end_latlng"`             // The end location of this activity
	AchievementCount   int         `json:"achievement_count"`      // The number of achievements gained during this activity
	KudosCount         int         `json:"kudos_count"`            // The number of kudos given for this activity
	CommentCount       int         `json:"comment_count"`          // The number of comments for this activity
	AthleteCount       int         `json:"athlete_count"`          // The number of athletes for taking part in a group activity
	PhotoCount         int         `json:"photo_count"`            // The number of Instagram photos for this activity
	TotalPhotoCount    int         `json:"total_photo_count"`      // The number of Instagram and Strava photos for this activity
	Map                ActivityMap `json:"map"`                    // The map data for this activity
	Trainer            bool        `json:"trainer"`                // Whether this activity was recorded on a training machine
	Commute            bool        `json:"commute"`                // Whether this activity is a commute
	Manual             bool        `json:"manual"`                 // Whether this activity was created manually
	Private            bool        `json:"private"`                // Whether this activity is private
	Flagged            bool        `json:"flagged"`                // Whether this activity is flagged
	WorkoutType        int         `json:"workout_type"`           // The activity's workout type
	UploadIDStr        string      `json:"upload_id_str"`          // The unique identifier of the upload in string format
	AverageSpeed       float64     `json:"average_speed"`          // The activity's average speed, in meters per second
	MaxSpeed           float64     `json:"max_speed"`              // The activity's max speed, in meters per second
	HasKudoed          bool        `json:"has_kudoed"`             // Whether the logged-in athlete has kudoed this activity
	HideFromHome       bool        `json:"hide_from_home"`         // Whether the activity is muted
	GearID             string      `json:"gear_id"`                // The id of the gear for the activity
	Kilojoules         float64     `json:"kilojoules"`             // The total work done in kilojoules during this activity. Rides only
	AverageWatts       float64     `json:"average_watts"`          // Average power output in watts during this activity. Rides only
	DeviceWatts        bool        `json:"device_watts"`           // Whether the watts are from a power meter, false if estimated
	MaxWatts           int         `json:"max_watts"`              // Maximum watts. Rides with power meter data only
	WeightedAvgWatts   int         `json:"weighted_average_watts"` // Similar to Normalized Power. Rides with power meter data only
	Description        string      `json:"description"`            // The description of the activity
	Photos             any         `json:"photos"`                 // The photos attached to this activity
	Gear               any         `json:"gear"`                   // The gear used in this activity
	Calories           float64     `json:"calories"`               // The number of kilocalories consumed during this activity
	SegmentEfforts     []any       `json:"segment_efforts"`        // The segments traversed in this activity
	DeviceName         string      `json:"device_name"`            // The name of the device used to record the activity
	EmbedToken         string      `json:"embed_token"`            // The token used to embed a Strava activity
	SplitsMetric       []any       `json:"splits_metric"`          // The splits of this activity in metric units (for runs)
	SplitsStandard     []any       `json:"splits_standard"`        // The splits of this activity in imperial units (for runs)
	Laps               []any       `json:"laps"`                   // The laps of this activity
	BestEfforts        []any       `json:"best_efforts"`           // The best efforts of this activity
}

// IanaTimezone extracts the IANA timezone from the Strava timezone string.
// Strava stores the timezone as a string in the format "(GMT+00:00) Europe/Lisbon".
// This function returns the IANA timezone identifier, e.g. "Europe/Lisbon".
func (a ActivityDetailed) IanaTimezone() string {
	if a.Timezone == "" {
		return ""
	}

	parts := strings.Split(a.Timezone, " ")

	return parts[len(parts)-1]
}

func (a ActivityDetailed) IsEnduranceOutdoorActivity() (bool, error) {
	sport, err := a.Sport()
	if err != nil {
		return false, err
	}

	return stride.IsEnduranceOutdoorActivity(sport), nil
}

func (a ActivityDetailed) SummaryPolyline() string {
	return a.Map.SummaryPolyline
}

func (a ActivityDetailed) Sport() (stride.Sport, error) {
	switch a.SportType {
	case SportTypeRun:
		return stride.SportRunning, nil

	case SportTypeTrailRun:
		return stride.SportTrailRunning, nil

	case SportTypeRide:
		return stride.SportCycling, nil

	case SportTypeGravelRide:
		return stride.SportGravelCycling, nil

	case SportTypeHike:
		return stride.SportHiking, nil
	default:
		return "", fmt.Errorf("%w: %s", stride.ErrUnsupportedSportType, a.SportType)
	}
}

func (a ActivityDetailed) ToActivity() (*stride.Activity, error) {
	return &stride.Activity{
		StartTime:     a.StartDate,
		ElapsedTime:   uint32(a.ElapsedTime),
		MovingTime:    uint32(a.MovingTime),
		Distance:      uint32(a.Distance),
		AvgSpeed:      uint16(a.AverageSpeed),
		ElevationGain: stride.Optional[uint16]{Valid: true, Value: uint16(a.TotalElevationGain)},
	}, nil
}

type ActivityMap struct {
	ID              string `json:"id"`
	SummaryPolyline string `json:"summary_polyline"`
	ResourceState   int    `json:"resource_state"`
}

type ActivityStream struct {
	Time           StreamSet[int]        `json:"time"`
	Distance       StreamSet[float64]    `json:"distance"`
	LatLng         StreamSet[[2]float64] `json:"latlng"`
	Altitude       StreamSet[float64]    `json:"altitude"`
	VelocitySmooth StreamSet[float64]    `json:"velocity_smooth"`
	Heartrate      StreamSet[int]        `json:"heartrate"`
	Cadence        StreamSet[int]        `json:"cadence"`
	Watts          StreamSet[int]        `json:"watts"`
	Temperature    StreamSet[int]        `json:"temp"`
	Moving         StreamSet[bool]       `json:"moving"`
	GradeSmooth    StreamSet[float64]    `json:"grade_smooth"`
}

func (s *ActivityStream) ToTimeseries(startTime time.Time) (*stride.ActivityTimeseries, error) {
	ts := stride.ActivityTimeseries{
		StartTime: startTime,
		Data:      []stride.ActivityTimeseriesEntry{},
	}

	for i := 0; i < len(s.Time.Data); i++ {
		data := stride.ActivityTimeseriesEntry{
			Offset: i,
		}

		if i < len(s.Heartrate.Data) {
			data.HeartRate = stride.Optional[uint8]{Value: uint8(s.Heartrate.Data[i]), Valid: s.Heartrate.Data[i] > 0}
		}

		if i < len(s.Cadence.Data) {
			data.Cadence = stride.Optional[uint8]{Value: uint8(s.Cadence.Data[i]), Valid: s.Cadence.Data[i] > 0}
		}

		if i < len(s.Distance.Data) {
			data.Distance = stride.Optional[uint32]{Value: uint32(s.Distance.Data[i]), Valid: s.Distance.Data[i] > 0}
		}

		if i < len(s.Altitude.Data) {
			data.Altitude = stride.Optional[uint16]{Value: uint16(s.Altitude.Data[i]), Valid: s.Altitude.Data[i] > 0}
		}

		if i < len(s.VelocitySmooth.Data) {
			data.Velocity = stride.Optional[uint16]{Value: uint16(s.VelocitySmooth.Data[i]), Valid: s.VelocitySmooth.Data[i] > 0}
		}

		if i < len(s.LatLng.Data) {
			latlng := s.LatLng.Data[i]
			data.Latitude = stride.Optional[float64]{Value: latlng[0], Valid: latlng[0] != 0 || latlng[1] != 0}
			data.Longitude = stride.Optional[float64]{Value: latlng[1], Valid: latlng[0] != 0 || latlng[1] != 0}
		}

		ts.Data = append(ts.Data, data)
	}

	return &ts, nil
}

type StreamSet[T ~bool | ~int | ~float64 | ~[2]float64] struct {
	OriginalSize int    `json:"original_size"` // The number of data points in this stream.
	Resolution   string `json:"resolution"`    // The level of detail (sampling) in which this stream was returned. May take one of the following values: low, medium, high.
	SeriesType   string `json:"series_type"`   // The base series used in the case the stream was downsampled. May take one of the following values: distance, time.
	Data         []T    `json:"data"`          // The sequence of values for this stream.
}
