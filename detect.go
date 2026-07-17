package main

import "math"

// haversineDistance returns the great-circle distance in km between two coordinates.
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371.0
	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return r * c
}

func toRad(deg float64) float64 {
	return deg * math.Pi / 180
}

// NearestDetector finds the closest postal code to given coordinates.
type NearestDetector struct {
	data []PostalCode
}

// NewNearestDetector creates a detector from the given records.
func NewNearestDetector(data []PostalCode) *NearestDetector {
	return &NearestDetector{data: data}
}

// Detect finds the nearest postal code. Returns nil if data is empty.
func (d *NearestDetector) Detect(lat, lon float64) *PostalCodeResult {
	if len(d.data) == 0 {
		return nil
	}

	best := &d.data[0]
	bestDist := haversineDistance(lat, lon, best.Latitude, best.Longitude)

	for i := 1; i < len(d.data); i++ {
		dist := haversineDistance(lat, lon, d.data[i].Latitude, d.data[i].Longitude)
		if dist < bestDist {
			best = &d.data[i]
			bestDist = dist
		}
	}

	return &PostalCodeResult{
		Province:  best.Province,
		Regency:   best.Regency,
		District:  best.District,
		Village:   best.Village,
		Code:      best.Code,
		Latitude:  best.Latitude,
		Longitude: best.Longitude,
		Elevation: best.Elevation,
		Timezone:  best.Timezone,
		Distance:  bestDist,
	}
}
