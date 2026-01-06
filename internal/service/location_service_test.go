package service

import (
	"testing"
)

func TestCalculateDistance(t *testing.T) {
	tests := []struct {
		name      string
		lat1      float64
		lon1      float64
		lat2      float64
		lon2      float64
		expected  float64 // приблизительное расстояние в метрах
		tolerance float64 // допустимая погрешность
	}{
		{
			name:      "Москва - Санкт-Петербург",
			lat1:      55.7558,
			lon1:      37.6173,
			lat2:      59.9343,
			lon2:      30.3351,
			expected:  635000, // примерно 635 км
			tolerance: 10000,  // погрешность 10 км
		},
		{
			name:      "Одна и та же точка",
			lat1:      55.7558,
			lon1:      37.6173,
			lat2:      55.7558,
			lon2:      37.6173,
			expected:  0,
			tolerance: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance := CalculateDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			diff := distance - tt.expected
			if diff < 0 {
				diff = -diff
			}
			if diff > tt.tolerance {
				t.Errorf("CalculateDistance() = %v, expected around %v (tolerance %v), got difference %v",
					distance, tt.expected, tt.tolerance, diff)
			}
		})
	}
}
