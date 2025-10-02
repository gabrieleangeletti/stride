package stride

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

func IsEnduranceActivity(sport Sport) bool {
	switch sport {
	case SportRunning, SportCycling, SportGravelCycling, SportHiking, SportTrailRunning, SportElliptical, SportSwimming, SportInlineSkating, SportStairStepper:
		return true
	default:
		return false
	}
}
