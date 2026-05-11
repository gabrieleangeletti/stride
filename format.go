package stride

import (
	"fmt"
	"math"
)

// formatPace helper converts speed (m/s) into standard MM:SS/km pace
func formatPace(speedMs float64) string {
	if speedMs <= 0.2 || math.IsNaN(speedMs) { // Slower than 83 min/km
		return "0:00"
	}

	secPerKm := 1000.0 / speedMs
	if math.IsNaN(secPerKm) || math.IsInf(secPerKm, 0) || secPerKm > 3600 {
		return "0:00"
	}

	mins := int(secPerKm) / 60
	secs := int(secPerKm) % 60

	return fmt.Sprintf("%d:%02d", mins, secs)
}
