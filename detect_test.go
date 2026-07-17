package main

import (
	"math"
	"testing"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name                   string
		lat1, lon1, lat2, lon2 float64
		wantKm                 float64
		tolerance              float64
	}{
		{
			name: "same point",
			lat1: -6.547, lon1: 107.398, lat2: -6.547, lon2: 107.398,
			wantKm: 0, tolerance: 0.001,
		},
		{
			name: "known distance jakarta to bandung",
			lat1: -6.208, lon1: 106.845, lat2: -6.917, lon2: 107.619,
			wantKm: 120, tolerance: 10,
		},
		{
			name: "jakarta to surabaya",
			lat1: -6.208, lon1: 106.845, lat2: -7.257, lon2: 112.752,
			wantKm: 670, tolerance: 15,
		},
		{
			name: "antipodal points",
			lat1: 0, lon1: 0, lat2: 0, lon2: 180,
			wantKm: 20015, tolerance: 50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := haversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if math.Abs(got-tt.wantKm) > tt.tolerance {
				t.Errorf("haversineDistance(%v,%v,%v,%v) = %f, want ~%f (±%f)",
					tt.lat1, tt.lon1, tt.lat2, tt.lon2, got, tt.wantKm, tt.tolerance)
			}
		})
	}
}

func TestNearestDetector(t *testing.T) {
	data := []PostalCode{
		{Province: "Jawa Barat", Regency: "Purwakarta", District: "Jatiluhur", Village: "Kembangkuning", Code: 41152, Latitude: -6.549, Longitude: 107.412, Elevation: 112, Timezone: "WIB"},
		{Province: "Jawa Barat", Regency: "Bogor", District: "Bogor", Village: "Batu Tulis", Code: 16124, Latitude: -6.656, Longitude: 106.801, Elevation: 240, Timezone: "WIB"},
		{Province: "Jawa Tengah", Regency: "Semarang", District: "Semarang", Village: "Semarang", Code: 50111, Latitude: -6.966, Longitude: 110.419, Elevation: 3, Timezone: "WIB"},
	}
	d := NewNearestDetector(data)

	result := d.Detect(-6.547, 107.398)
	if result == nil {
		t.Fatal("expected a result, got nil")
	}
	if result.Code != 41152 {
		t.Errorf("expected nearest code 41152, got %d", result.Code)
	}
	if result.Distance < 0 {
		t.Errorf("distance should be >= 0, got %f", result.Distance)
	}
}

func TestNearestDetectorEmpty(t *testing.T) {
	d := NewNearestDetector(nil)

	result := d.Detect(-6.547, 107.398)
	if result != nil {
		t.Errorf("expected nil for empty data, got %+v", result)
	}
}

func TestNearestDetectorSingle(t *testing.T) {
	data := []PostalCode{
		{Province: "DKI Jakarta", Regency: "Jakarta", District: "Menteng", Village: "Menteng", Code: 10310, Latitude: -6.190, Longitude: 106.840, Elevation: 5, Timezone: "WIB"},
	}
	d := NewNearestDetector(data)

	result := d.Detect(-6.190, 106.840)
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Code != 10310 {
		t.Errorf("expected code 10310, got %d", result.Code)
	}
	if result.Distance > 0.001 {
		t.Errorf("expected ~0 distance, got %f", result.Distance)
	}
}
