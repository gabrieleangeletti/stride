package stride

import "fmt"

type Provider string

const (
	ProviderStrava Provider = "strava"
)

type Sport string

const (
	SportCycling       Sport = "cycling"
	SportGravelCycling Sport = "gravel-cycling"
	SportElliptical    Sport = "elliptical"
	SportHiking        Sport = "hiking"
	SportInlineSkating Sport = "inline-skating"
	SportKayaking      Sport = "kayaking"
	SportRockClimbing  Sport = "rock-climbing"
	SportRunning       Sport = "running"
	SportStairStepper  Sport = "stair-stepper"
	SportSurfing       Sport = "surfing"
	SportSwimming      Sport = "swimming"
	SportTrailRunning  Sport = "trail-running"
)

var validSports = map[Sport]struct{}{
	SportCycling:       {},
	SportGravelCycling: {},
	SportElliptical:    {},
	SportHiking:        {},
	SportInlineSkating: {},
	SportKayaking:      {},
	SportRockClimbing:  {},
	SportRunning:       {},
	SportStairStepper:  {},
	SportSurfing:       {},
	SportSwimming:      {},
	SportTrailRunning:  {},
}

// ParseSport validates and converts a string to Sport
func ParseSport(s string) (Sport, error) {
	sp := Sport(s)

	if _, ok := validSports[sp]; !ok {
		return "", fmt.Errorf("%w: %q", ErrUnsupportedSportType, s)
	}

	return sp, nil
}

func IsEnduranceActivity(sport Sport) bool {
	switch sport {
	case SportRunning, SportCycling, SportGravelCycling, SportHiking, SportTrailRunning, SportElliptical, SportSwimming, SportInlineSkating, SportStairStepper:
		return true
	default:
		return false
	}
}
