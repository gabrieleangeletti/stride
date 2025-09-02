package stride

import "errors"

var (
	ErrRateLimitExceeded             = errors.New("rate limit exceeded")
	ErrActivityIsNotOutdoorEndurance = errors.New("activity is not an outdoor endurance activity")
)
