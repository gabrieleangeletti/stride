package stride

import "time"

type Activity struct {
	Sport         Sport
	StartTime     time.Time        // UTC
	ElapsedTime   uint32           // seconds
	MovingTime    uint32           // seconds
	Distance      uint32           // meters
	AvgSpeed      uint16           // m / s
	AvgHR         Optional[uint8]  // beats / minute
	MaxHR         Optional[uint8]  // beats / minute
	ElevationGain Optional[uint16] // meters
	ElevationLoss Optional[uint16] // meters
}

type ActivityTimeseries struct {
	StartTime time.Time
	Data      []ActivityTimeseriesEntry
}

func (ts ActivityTimeseries) EndTime() time.Time {
	return ts.StartTime.Add(time.Duration(ts.MaxOffset()) * time.Second)
}

func (ts ActivityTimeseries) ElapsedTime() time.Duration {
	return ts.EndTime().Sub(ts.StartTime)
}

func (ts ActivityTimeseries) MaxOffset() int {
	maxOffset := 0

	for _, entry := range ts.Data {
		if entry.Offset > maxOffset {
			maxOffset = entry.Offset
		}
	}

	return maxOffset
}

type ActivityTimeseriesEntry struct {
	Offset    int
	HeartRate Optional[uint8]
	Cadence   Optional[uint8]
	Distance  Optional[uint32]
	Altitude  Optional[uint16]
	Velocity  Optional[uint16]
	Latitude  Optional[float64]
	Longitude Optional[float64]
}

func (a ActivityTimeseriesEntry) IsEmpty() bool {
	return !a.HeartRate.Valid &&
		!a.Cadence.Valid &&
		!a.Distance.Valid &&
		!a.Altitude.Valid &&
		!a.Velocity.Valid &&
		!a.Latitude.Valid &&
		!a.Longitude.Valid
}

func (a ActivityTimeseriesEntry) HasGPS() bool {
	return a.Latitude.Valid && a.Longitude.Valid
}
