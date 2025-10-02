package stride

import "errors"

var (
	ErrRateLimitExceeded      = errors.New("rate limit exceeded")
	ErrActivityIsNotEndurance = errors.New("activity is not an endurance activity")
	ErrUnsupportedSportType   = errors.New("unsupported sport type")
)
