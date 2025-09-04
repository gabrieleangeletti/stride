package stride

import (
	"bytes"
	"fmt"
	"math"
	"time"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/encoder"
	"github.com/muktihari/fit/profile/filedef"
	"github.com/muktihari/fit/profile/mesgdef"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/fit/proto"
)

type FITSport struct {
	Sport    typedef.Sport
	SubSport typedef.SubSport
}

func CreateFITFileInMemory(data *Activity, ts *ActivityTimeseries, sport Sport) ([]byte, error) {
	now := time.Now()

	activity := filedef.NewActivity()

	activity.FileId = *mesgdef.NewFileId(nil).
		SetType(typedef.FileActivity).
		SetTimeCreated(now).
		SetManufacturer(typedef.ManufacturerGarmin).
		SetProductName("Enduro 3")

	activity.Activity = mesgdef.NewActivity(nil).
		SetType(typedef.ActivityManual).
		SetTimestamp(data.StartTime).
		SetNumSessions(1)

	fitSport, err := sportToFitSport(sport)
	if err != nil {
		return nil, err
	}

	session := mesgdef.NewSession(nil).
		SetStartTime(data.StartTime).
		SetTotalElapsedTime(data.ElapsedTime).
		SetTotalMovingTime(data.MovingTime).
		SetTotalTimerTime(data.ElapsedTime).
		SetTotalDistance(data.Distance).
		SetSport(fitSport.Sport).
		SetSubSport(fitSport.SubSport).
		SetAvgSpeed(data.AvgSpeed)

	if data.AvgHR.Valid {
		session = session.SetAvgHeartRate(data.AvgHR.Value)
	}

	if data.MaxHR.Valid {
		session = session.SetMaxHeartRate(data.MaxHR.Value)
	}

	if data.ElevationGain.Valid {
		session = session.SetTotalAscent(data.ElevationGain.Value)
	}

	if data.ElevationLoss.Valid {
		session = session.SetTotalDescent(data.ElevationLoss.Value)
	}

	activity.Sessions = append(activity.Sessions, session)

	for _, d := range ts.Data {
		if d.IsEmpty() {
			continue
		}

		t := ts.StartTime.Add(time.Duration(d.Offset) * time.Second)

		record := mesgdef.NewRecord(nil).SetTimestamp(t)

		if d.Distance.Valid {
			record = record.SetDistance(d.Distance.Value)
		}

		if d.Velocity.Valid {
			record = record.SetSpeed(d.Velocity.Value)
		}

		if d.Altitude.Valid {
			record = record.SetAltitude(d.Altitude.Value)
		}

		if d.Cadence.Valid {
			record = record.SetCadence(d.Cadence.Value)
		}

		if d.HeartRate.Valid {
			record = record.SetHeartRate(d.HeartRate.Value)
		}

		if d.Latitude.Valid && d.Longitude.Valid {
			record = record.SetPositionLatDegrees(d.Latitude.Value)
			record = record.SetPositionLongDegrees(d.Longitude.Value)
		}

		activity.Records = append(activity.Records, record)
	}

	fit := activity.ToFIT(nil)

	buf := new(bytes.Buffer)

	enc := encoder.New(buf, encoder.WithProtocolVersion(proto.V2))
	if err := enc.Encode(&fit); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func FITFileToActivityTimeseries(data []byte) (*ActivityTimeseries, error) {
	dec := decoder.New(bytes.NewReader(data))

	fit, err := dec.Decode()
	if err != nil {
		return nil, err
	}

	activity := filedef.NewActivity(fit.Messages...)

	startTime := activity.Sessions[0].StartTime

	timeseries := ActivityTimeseries{
		StartTime: startTime,
		Data:      make([]ActivityTimeseriesEntry, 0),
	}

	for _, record := range activity.Records {
		entry := ActivityTimeseriesEntry{
			Offset:    int(record.Timestamp.Unix() - startTime.Unix()),
			HeartRate: Optional[uint8]{Value: record.HeartRate, Valid: record.HeartRate > 0},
			Cadence:   Optional[uint8]{Value: record.Cadence, Valid: record.Cadence > 0},
			Velocity:  Optional[uint16]{Value: uint16(record.SpeedScaled()), Valid: !math.IsNaN(record.SpeedScaled())},
			Altitude:  Optional[uint16]{Value: uint16(record.AltitudeScaled()), Valid: !math.IsNaN(record.AltitudeScaled())},
			Distance:  Optional[uint32]{Value: uint32(record.DistanceScaled()), Valid: !math.IsNaN(record.DistanceScaled())},
		}

		// Parse GPS coordinates if available
		if !math.IsNaN(record.PositionLatDegrees()) && !math.IsNaN(record.PositionLongDegrees()) {
			lat := record.PositionLatDegrees()
			lon := record.PositionLongDegrees()
			entry.Latitude = Optional[float64]{Value: lat, Valid: lat != 0}
			entry.Longitude = Optional[float64]{Value: lon, Valid: lon != 0}
		}

		timeseries.Data = append(timeseries.Data, entry)
	}

	return &timeseries, nil
}

func sportToFitSport(sport Sport) (FITSport, error) {
	switch sport {
	case SportCycling:
		return FITSport{Sport: typedef.SportCycling}, nil

	case SportElliptical:
		return FITSport{Sport: typedef.SportFitnessEquipment, SubSport: typedef.SubSportElliptical}, nil

	case SportHiking:
		return FITSport{Sport: typedef.SportHiking}, nil

	case SportInlineSkating:
		return FITSport{Sport: typedef.SportInlineSkating}, nil

	case SportKayaking:
		return FITSport{Sport: typedef.SportKayaking}, nil

	case SportRockClimbing:
		return FITSport{Sport: typedef.SportRockClimbing}, nil

	case SportRunning:
		return FITSport{Sport: typedef.SportRunning}, nil

	case SportStairStepper:
		return FITSport{Sport: typedef.SportFitnessEquipment, SubSport: typedef.SubSportStairClimbing}, nil

	case SportSurfing:
		return FITSport{Sport: typedef.SportSurfing}, nil

	case SportSwimming:
		return FITSport{Sport: typedef.SportSwimming}, nil

	case SportTrailRunning:
		return FITSport{Sport: typedef.SportRunning, SubSport: typedef.SubSportTrail}, nil

	default:
		return FITSport{}, fmt.Errorf("unknown sport: %s", sport)
	}
}
