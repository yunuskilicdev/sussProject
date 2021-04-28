package main

import "math"

const (
	earthRadiusMi = 3958
)

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

func Distance(fromLat float64, fromLon float64, toLat float64, toLon float64, fromRadius int64, toRadius int64) (mi float64) {
	lat1 := degreesToRadians(fromLat)
	lon1 := degreesToRadians(fromLon)
	lat2 := degreesToRadians(toLat)
	lon2 := degreesToRadians(toLon)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	mi = c * earthRadiusMi

	// Haversine formula calculate between coordinates, but we have uncertainty radius for each coordinate.
	// We should consider these radius. If sum of the radius are bigger than calculated difference, I subtract radius sum from difference
	fromRadiusMile := float64(fromRadius) / 1.609344
	toRadiusMile := float64(toRadius) / 1.609344

	totalUncertaintyMile := fromRadiusMile + toRadiusMile

	if mi > totalUncertaintyMile {
		mi -= totalUncertaintyMile
	}

	return mi
}
