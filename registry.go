package stride

import (
	"encoding/json"
	"sync"
)

type ActivityMetricsFactory func() ActivityMetrics

type ActivityRegistry struct {
	mu     sync.RWMutex
	byKind map[string]ActivityMetricsFactory
	// optional: defaults mapping for each Sport
	defaultBySport map[Sport]string
}

func NewActivityRegistry() *ActivityRegistry {
	return &ActivityRegistry{
		byKind:         make(map[string]ActivityMetricsFactory),
		defaultBySport: make(map[Sport]string),
	}
}

func (r *ActivityRegistry) RegisterKind(kind string, f ActivityMetricsFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byKind[kind] = f
}

func (r *ActivityRegistry) MapSportToKind(s Sport, kind string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.defaultBySport[s] = kind
}

func (r *ActivityRegistry) factoryFor(kind string) ActivityMetricsFactory {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.byKind[kind]
}

func (r *ActivityRegistry) defaultKindFor(s Sport) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	k, ok := r.defaultBySport[s]
	return k, ok
}

// DecodeMetrics chooses the concrete metrics type by "kind" field.
// If kind is absent, it will try defaultBySport (if provided).
// If no default is found, it will return an UnknownMetrics instance.
func (r *ActivityRegistry) DecodeMetrics(raw json.RawMessage, sport Sport) (ActivityMetrics, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var peek struct {
		Kind string `json:"kind"`
	}
	_ = json.Unmarshal(raw, &peek)

	kind := peek.Kind
	if kind == "" {
		if dk, ok := r.defaultKindFor(sport); ok {
			kind = dk
		}
	}

	if kind == "" {
		return NewUnknown("unknown", raw), nil
	}

	if f := r.factoryFor(kind); f != nil {
		m := f()
		if err := json.Unmarshal(raw, m); err != nil {
			return nil, err
		}

		return m, nil
	}

	return NewUnknown(kind, raw), nil
}

func (r *ActivityRegistry) DecodeActivityJSON(b []byte) (*Activity, error) {
	var aux struct {
		ExternalID   string          `json:"externalId"` // The id given by the provider
		Provider     Provider        `json:"provider"`
		Sport        Sport           `json:"sport"`
		StartTime    TimeRFC3339     `json:"startTime"`
		EndTime      TimeRFC3339     `json:"endTime"`
		IanaTimezone string          `json:"ianaTimezone"`
		UTCOffset    int             `json:"utcOffset"` // seconds
		Metrics      json.RawMessage `json:"metrics"`
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return nil, err
	}

	m, err := r.DecodeMetrics(aux.Metrics, aux.Sport)
	if err != nil {
		return nil, err
	}

	return &Activity{
		ExternalID:   aux.ExternalID,
		Provider:     aux.Provider,
		Sport:        aux.Sport,
		StartTime:    aux.StartTime.Time,
		EndTime:      aux.EndTime.Time,
		IanaTimezone: aux.IanaTimezone,
		UTCOffset:    aux.UTCOffset,
		Metrics:      m,
	}, nil
}

func (r *ActivityRegistry) EncodeActivityJSON(a *Activity) ([]byte, error) {
	type alias Activity
	return json.Marshal(alias(*a))
}
