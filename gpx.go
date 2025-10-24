package stride

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

var (
	ErrNoTrackPoints = errors.New("no track points found")
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
			Timestamp: t,
		}

		if d.HasGPS() {
			point.Point = gpx.Point{
				Latitude:  d.Latitude.Value,
				Longitude: d.Longitude.Value,
			}
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

func ParseGPXFileFromMemory(data []byte) (*Activity, *ActivityTimeseries, Sport, error) {
	gpxFile, err := gpx.ParseBytes(data)
	if err != nil {
		return nil, nil, SportUnknown, fmt.Errorf("failed to parse GPX: %w", err)
	}

	if len(gpxFile.Tracks) == 0 || len(gpxFile.Tracks[0].Segments) == 0 {
		return nil, nil, SportUnknown, errors.New("no tracks or segments found in GPX")
	}

	track := gpxFile.Tracks[0]
	segment := track.Segments[0]
	points := segment.Points

	if len(points) == 0 {
		return nil, nil, SportUnknown, ErrNoTrackPoints
	}

	sport := gpxNameToSport(track.Type, track.Name)

	startTime := points[0].Timestamp
	endTime := points[len(points)-1].Timestamp

	ts := &ActivityTimeseries{
		StartTime: startTime,
	}

	for _, p := range points {
		offset := int(p.Timestamp.Sub(startTime).Seconds())

		entry := ActivityTimeseriesEntry{
			Offset:    offset,
			Latitude:  NewOptional(p.Point.Latitude, true),
			Longitude: NewOptional(p.Point.Longitude, true),
		}

		if !math.IsNaN(p.Point.Elevation.Value()) {
			entry.Altitude = NewOptional(uint16(p.Point.Elevation.Value()), true)
		}

		for _, ext := range p.Extensions.Nodes {
			if ext.XMLName.Local == "TrackPointExtension" {
				for _, sub := range ext.Nodes {
					switch sub.XMLName.Local {
					case "hr":
						var hr uint8
						fmt.Sscanf(sub.Data, "%d", &hr)
						entry.HeartRate = NewOptional(hr, true)
					case "cad":
						var cad uint8
						fmt.Sscanf(sub.Data, "%d", &cad)
						entry.Cadence = NewOptional(cad, true)
					}
				}
			}
		}

		ts.Data = append(ts.Data, entry)
	}

	elapsed := endTime.Sub(startTime)

	activity := &Activity{
		StartTime:   startTime,
		ElapsedTime: uint32(elapsed.Seconds()),
		// We can’t directly infer MovingTime, AvgSpeed, etc., from GPX alone,
		// so they’ll be left zeroed or estimated if desired.
	}

	return activity, ts, sport, nil
}

func gpxNameToSport(gpxType, gpxName string) Sport {
	switch gpxType {
	case "biking":
		return SportCycling
	case "fitness":
		// Could be elliptical or stair stepper — fallback to name
		if gpxName == "Elliptical" {
			return SportElliptical
		}
		if gpxName == "Stair Stepper" {
			return SportStairStepper
		}
		return SportUnknown
	case "hiking":
		return SportHiking
	case "skating":
		return SportInlineSkating
	case "paddling":
		return SportKayaking
	case "climbing":
		return SportRockClimbing
	case "running":
		return SportRunning
	case "trail_running":
		return SportTrailRunning
	case "water":
		return SportSurfing
	case "swimming":
		return SportSwimming
	default:
		return SportUnknown
	}
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
