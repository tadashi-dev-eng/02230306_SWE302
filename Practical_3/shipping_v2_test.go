// shipping_v2_test.go
package shipping

import (
	"strings"
	"testing"
)

func TestCalculateShippingFee_V2(t *testing.T) {
	testCases := []struct {
		name         string
		weight       float64
		zone         string
		insured      bool
		expectedFee  float64
		expectError  bool
		errorContent string
	}{
		// Invalid Weight Tests (P1 & P4)
		{"Invalid: Weight too small", -5, "Domestic", false, 0, true, "invalid weight"},
		{"Invalid: Weight zero", 0, "International", true, 0, true, "invalid weight"},
		{"Invalid: Weight too large", 60, "Express", false, 0, true, "invalid weight"},

		// Invalid Zone Tests (P6)
		{"Invalid: Unknown zone", 10, "Local", false, 0, true, "invalid zone"},
		{"Invalid: Empty zone", 5, "", true, 0, true, "invalid zone"},
		{"Invalid: Case sensitive zone", 15, "domestic", false, 0, true, "invalid zone"},

		// Boundary Value Tests - Lower boundary (around 0)
		{"Boundary: Just above zero", 0.1, "Domestic", false, 5.0, false, ""},

		// Boundary Value Tests - Mid boundary (around 10)
		{"Boundary: Just below 10kg", 9.9, "Domestic", false, 5.0, false, ""},
		{"Boundary: Exactly 10kg", 10, "International", false, 20.0, false, ""},
		{"Boundary: Just above 10kg", 10.1, "Express", false, 37.5, false, ""},

		// Boundary Value Tests - Upper boundary (around 50)
		{"Boundary: Just below 50kg", 49.9, "Domestic", false, 12.5, false, ""},
		{"Boundary: Exactly 50kg", 50, "International", false, 27.5, false, ""},

		// Standard Package Tests (P2) - Weight 0 < w ≤ 10
		{"Standard: Domestic, No Insurance", 5, "Domestic", false, 5.0, false, ""},
		{"Standard: International, No Insurance", 8, "International", false, 20.0, false, ""},
		{"Standard: Express, No Insurance", 3, "Express", false, 30.0, false, ""},

		// Standard Package with Insurance (P2 + P7)
		{"Standard: Domestic with Insurance", 5, "Domestic", true, 5.075, false, ""},
		{"Standard: International with Insurance", 8, "International", true, 20.3, false, ""},
		{"Standard: Express with Insurance", 3, "Express", true, 30.45, false, ""},

		// Heavy Package Tests (P3) - Weight 10 < w ≤ 50
		{"Heavy: Domestic, No Insurance", 15, "Domestic", false, 12.5, false, ""},
		{"Heavy: International, No Insurance", 25, "International", false, 27.5, false, ""},
		{"Heavy: Express, No Insurance", 40, "Express", false, 37.5, false, ""},

		// Heavy Package with Insurance (P3 + P7)
		{"Heavy: Domestic with Insurance", 15, "Domestic", true, 12.6875, false, ""},
		{"Heavy: International with Insurance", 25, "International", true, 27.9125, false, ""},
		{"Heavy: Express with Insurance", 40, "Express", true, 38.0625, false, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fee, err := CalculateShippingFee(tc.weight, tc.zone, tc.insured)

			// Check error expectations
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error containing '%s', but got nil", tc.errorContent)
				} else if !strings.Contains(err.Error(), tc.errorContent) {
					t.Errorf("Expected error to contain '%s', but got: %v", tc.errorContent, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}

				// Check fee calculation (with small tolerance for floating point)
				tolerance := 0.001
				if fee < tc.expectedFee-tolerance || fee > tc.expectedFee+tolerance {
					t.Errorf("Expected fee %f, but got %f", tc.expectedFee, fee)
				}
			}
		})
	}
}
