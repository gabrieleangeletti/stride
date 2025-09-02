package stride

func DefaultRegistry() *ActivityRegistry {
	r := NewActivityRegistry()

	r.RegisterKind("endurance", func() ActivityMetrics { return NewEnduranceMetrics() })
	r.RegisterKind("strength", func() ActivityMetrics { return NewStrengthMetrics() })
	r.RegisterKind("climbing", func() ActivityMetrics { return NewClimbingMetrics() })

	r.MapSportToKind(SportCycling, "endurance")
	r.MapSportToKind(SportElliptical, "endurance")
	r.MapSportToKind(SportHiking, "endurance")
	r.MapSportToKind(SportInlineSkating, "endurance")
	r.MapSportToKind(SportKayaking, "endurance")
	r.MapSportToKind(SportRockClimbing, "climbing")
	r.MapSportToKind(SportRunning, "endurance")
	r.MapSportToKind(SportStairStepper, "endurance")
	r.MapSportToKind(SportSurfing, "endurance")
	r.MapSportToKind(SportSwimming, "endurance")
	r.MapSportToKind(SportTrailRunning, "endurance")

	return r
}
