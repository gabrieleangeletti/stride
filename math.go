package stride

import "math"

// haversine calculates the distance between two lat/lon points in meters.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000.0 // Earth radius in meters

	rad := math.Pi / 180.0
	dLat := (lat2 - lat1) * rad
	dLon := (lon2 - lon1) * rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// calculateGAP computes Grade Adjusted Pace speed (m/s) using Minetti's energy cost formula.
func calculateGAP(actualSpeedMs float64, gradeFraction float64) float64 {
	if actualSpeedMs < 0.5 {
		return 0 // Not moving
	}

	// Limit grade to sensible running physics (-20% to +30%) to prevent polynomial explosion
	g := math.Max(-0.20, math.Min(0.30, gradeFraction))

	// Minetti's cost curve
	cost := 155.4*math.Pow(g, 5) - 30.4*math.Pow(g, 4) - 43.3*math.Pow(g, 3) + 46.3*math.Pow(g, 2) + 19.5*g + 3.6
	flatCost := 3.6

	// Ratio of energy cost relative to flat ground
	costRatio := cost / flatCost
	if costRatio < 0.3 {
		costRatio = 0.3 // Cap maximum downhill speed assist
	}

	return actualSpeedMs * costRatio
}
