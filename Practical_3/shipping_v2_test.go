// shipping_v2_test.go
package shipping

import (
	"testing"
)

func TestCalculateShippingFee_V2(t *testing.T) {
	testCases := []struct {
		name        string
		weight      float64
		zone        string
		insured     bool
		expectedFee float64
		expectError bool
	}{
		// ===== PARTITION TESTS: WEIGHT =====
		// P1: Invalid weight - too small (≤ 0)
		{"P1: Weight zero", 0, "Domestic", false, 0, true},
		{"P1: Weight negative", -5, "International", false, 0, true},

		// P2: Valid Standard tier (0 < weight ≤ 10)
		{"P2: Standard package - Domestic", 5, "Domestic", false, 5.0, false},
		{"P2: Standard package - International", 8, "International", false, 20.0, false},
		{"P2: Standard package - Express", 10, "Express", false, 30.0, false},

		// P3: Valid Heavy tier (10 < weight ≤ 50)
		{"P3: Heavy package - Domestic", 25, "Domestic", false, 12.5, false},
		{"P3: Heavy package - International", 30, "International", false, 27.5, false},
		{"P3: Heavy package - Express", 40, "Express", false, 37.5, false},

		// P4: Invalid weight - too large (> 50)
		{"P4: Weight too large", 60, "Domestic", false, 0, true},
		{"P4: Weight way too large", 100, "International", false, 0, true},

		// ===== PARTITION TESTS: ZONE =====
		// P5: Valid zones (already covered above, but explicit tests here)
		{"P5: Valid zone - Domestic", 15, "Domestic", false, 12.5, false},
		{"P5: Valid zone - International", 15, "International", false, 27.5, false},
		{"P5: Valid zone - Express", 15, "Express", false, 37.5, false},

		// P6: Invalid zones
		{"P6: Invalid zone - Local", 10, "Local", false, 0, true},
		{"P6: Invalid zone - domestic (lowercase)", 10, "domestic", false, 0, true},
		{"P6: Invalid zone - empty string", 10, "", false, 0, true},

		// ===== PARTITION TESTS: INSURED =====
		// P7: Not insured (false) - covered by tests above
		{"P7: Not insured - Standard", 5, "Domestic", false, 5.0, false},
		{"P7: Not insured - Heavy", 25, "International", false, 27.5, false},

		// P8: Insured (true)
		// Insured Standard: 5 + (5 * 0.015) = 5 + 0.075 = 5.075
		{"P8: Insured Standard - Domestic", 5, "Domestic", true, 5.075, false},
		// Insured Standard: 20 + (20 * 0.015) = 20 + 0.3 = 20.3
		{"P8: Insured Standard - International", 8, "International", true, 20.3, false},
		// Insured Heavy: (5 + 7.5) + ((5 + 7.5) * 0.015) = 12.5 + 0.1875 = 12.6875
		{"P8: Insured Heavy - Domestic", 25, "Domestic", true, 12.6875, false},
		// Insured Heavy: (20 + 7.5) + ((20 + 7.5) * 0.015) = 27.5 + 0.4125 = 27.9125
		{"P8: Insured Heavy - International", 30, "International", true, 27.9125, false},

		// ===== BOUNDARY VALUE ANALYSIS: WEIGHT =====
		// Lower boundary around 0
		{"BVA: Weight at lower invalid (0)", 0, "Domestic", false, 0, true},
		{"BVA: Weight just above lower (0.1)", 0.1, "Domestic", false, 5.0, false},

		// Mid boundary around 10 (tier transition)
		{"BVA: Weight at tier boundary - Standard (10)", 10, "Domestic", false, 5.0, false},
		{"BVA: Weight at tier boundary - Heavy (10.1)", 10.1, "Domestic", false, 12.5, false},

		// Upper boundary around 50
		{"BVA: Weight at upper valid (50)", 50, "International", false, 27.5, false},
		{"BVA: Weight just above upper (50.1)", 50.1, "Express", false, 0, true},

		// ===== COMPREHENSIVE COMBINATION TESTS =====
		// All valid zones with Heavy package and insurance
		// Weight 35: (30 + 7.5) + ((30 + 7.5) * 0.015) = 37.5 + 0.5625 = 38.0625
		{"Combo: Heavy + Insured - Express", 35, "Express", true, 38.0625, false},
		// Edge case: minimum valid weight with maximum insured cost
		{"Combo: Min valid weight + Insured - Express", 0.1, "Express", true, 30.45, false},
		// Edge case: maximum valid weight with heavy surcharge and insurance
		// Weight 50: (20 + 7.5) + ((20 + 7.5) * 0.015) = 27.5 + 0.4125 = 27.9125
		{"Combo: Max valid weight + Insured - International", 50, "International", true, 27.9125, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fee, err := CalculateShippingFee(tc.weight, tc.zone, tc.insured)

			// Check error expectation
			if tc.expectError && err == nil {
				t.Errorf("Expected an error, but got nil. Fee: %f", fee)
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}

			// Check fee value only if no error was expected
			if !tc.expectError {
				const tolerance = 0.0001 // For floating-point comparison
				if fee < tc.expectedFee-tolerance || fee > tc.expectedFee+tolerance {
					t.Errorf("Expected fee %.4f, but got %.4f", tc.expectedFee, fee)
				}
			}
		})
	}
}