package stride

import (
	"fmt"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

func CreateGPXFileInMemory(data *Activity, ts *ActivityTimeseries, sport Sport) ([]byte, error) {
	gpxFile := &gpx.GPX{
		Version: "1.1",
		Creator: "Stride",
		Name:    sportToGPXName(sport),
		Time:    &data.StartTime,
	}

	track := gpx.GPXTrack{
		Name: sportToGPXName(sport),
		Type: sportToGPXType(sport),
	}

	segment := gpx.GPXTrackSegment{}

	for _, d := range ts.Data {
		if d.IsEmpty() {
			continue
		}

		t := ts.StartTime.Add(time.Duration(d.Offset) * time.Second)

		point := gpx.GPXPoint{
			Point: gpx.Point{
				Latitude:  d.Latitude.Value,
				Longitude: d.Longitude.Value,
			},
			Timestamp: t,
		}

		// Skip points without GPS coordinates
		if !d.Latitude.Valid || !d.Longitude.Valid {
			continue
		}

		if d.Altitude.Valid {
			elevation := float64(d.Altitude.Value)
			point.Point.Elevation = *gpx.NewNullableFloat64(elevation)
		}

		if d.HeartRate.Valid {
			hrNode := point.Extensions.GetOrCreateNode("http://www.garmin.com/xmlschemas/TrackPointExtension/v1", "TrackPointExtension", "hr")
			hrNode.Data = fmt.Sprintf("%d", d.HeartRate.Value)
		}

		if d.Cadence.Valid {
			cadNode := point.Extensions.GetOrCreateNode("http://www.garmin.com/xmlschemas/TrackPointExtension/v1", "TrackPointExtension", "cad")
			cadNode.Data = fmt.Sprintf("%d", d.Cadence.Value)
		}

		segment.Points = append(segment.Points, point)
	}

	track.Segments = append(track.Segments, segment)
	gpxFile.Tracks = append(gpxFile.Tracks, track)

	xmlBytes, err := gpxFile.ToXml(gpx.ToXmlParams{
		Version: "1.1",
		Indent:  true,
	})
	if err != nil {
		return nil, err
	}

	return xmlBytes, nil
}

func sportToGPXName(sport Sport) string {
	switch sport {
	case SportCycling:
		return "Cycling"
	case SportElliptical:
		return "Elliptical"
	case SportHiking:
		return "Hiking"
	case SportInlineSkating:
		return "Inline Skating"
	case SportKayaking:
		return "Kayaking"
	case SportRockClimbing:
		return "Rock Climbing"
	case SportRunning:
		return "Running"
	case SportStairStepper:
		return "Stair Stepper"
	case SportSurfing:
		return "Surfing"
	case SportSwimming:
		return "Swimming"
	case SportTrailRunning:
		return "Trail Running"
	default:
		return "Activity"
	}
}

func sportToGPXType(sport Sport) string {
	switch sport {
	case SportCycling:
		return "biking"
	case SportElliptical:
		return "fitness"
	case SportHiking:
		return "hiking"
	case SportInlineSkating:
		return "skating"
	case SportKayaking:
		return "paddling"
	case SportRockClimbing:
		return "climbing"
	case SportRunning:
		return "running"
	case SportStairStepper:
		return "fitness"
	case SportSurfing:
		return "water"
	case SportSwimming:
		return "swimming"
	case SportTrailRunning:
		return "trail_running"
	default:
		return "activity"
	}
}
