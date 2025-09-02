package stride

import (
	"encoding/json"
	"time"
)

type TimeRFC3339 struct{ time.Time }

func (t *TimeRFC3339) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	tt, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}

	t.Time = tt
	return nil
}

func (t TimeRFC3339) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Format(time.RFC3339Nano))
}
