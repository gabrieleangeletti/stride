package stride

import (
	"fmt"
	"strings"

	"github.com/twpayne/go-polyline"
)

func PolylineToWKT(poly string) (string, error) {
	coords, _, err := polyline.DecodeCoords([]byte(poly))
	if err != nil {
		return "", err
	}

	if len(coords) == 0 {
		return "", fmt.Errorf("empty polyline")
	}

	// WKT expects "lon lat", not "lat lon"
	points := make([]string, len(coords))
	for i, c := range coords {
		lat := c[0]
		lon := c[1]
		points[i] = fmt.Sprintf("%f %f", lon, lat)
	}

	wkt := fmt.Sprintf("LINESTRING(%s)", strings.Join(points, ", "))
	return wkt, nil
}
