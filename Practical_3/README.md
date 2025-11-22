# Practical 3: Specification-Based Testing in Go
## Overview

This practical exercise focuses on applying **specification-based testing** techniques to an updated shipping fee calculator function. The goal is to systematically design comprehensive test cases using **Equivalence Partitioning**, **Boundary Value Analysis**, and **Decision Table Testing** principles—all without examining the internal implementation code.

---

## Part 1: Test Design Analysis

### 1.1 Equivalence Partitioning

Equivalence Partitioning divides all possible inputs into groups where each group should be treated identically by the system. By testing one representative value from each partition, we can confidently assert that all values in that partition behave the same way.

#### Weight Input Partitions

| Partition | Range | Characteristic | Example Values | Expected Behavior |
|-----------|-------|-----------------|-----------------|-------------------|
| **P1: Too Small** | weight ≤ 0 | Invalid | -5, 0, -100 | Error returned |
| **P2: Standard** | 0 < weight ≤ 10 | Valid, no surcharge | 1, 5, 10 | Normal fee calculated |
| **P3: Heavy** | 10 < weight ≤ 50 | Valid, $7.50 surcharge applied | 15, 25, 50 | Fee + $7.50 surcharge |
| **P4: Too Large** | weight > 50 | Invalid | 51, 100, 1000 | Error returned |

**Rationale:** The specification explicitly defines three weight tiers: invalid (≤0), Standard (0<w≤10), Heavy (10<w≤50), and invalid (>50). Each tier has different calculation rules, making them distinct partitions.

#### Zone Input Partitions

| Partition | Values | Characteristic | Example Values | Expected Behavior |
|-----------|--------|-----------------|-----------------|-------------------|
| **P5: Valid Zones** | Domestic, International, Express | Valid zone strings | "Domestic", "International", "Express" | Corresponding base fee applied |
| **P6: Invalid Zones** | Anything else | Invalid or case-sensitive mismatch | "Local", "domestic", "", "DOMESTIC" | Error returned |

**Rationale:** The specification lists three exact valid zone strings. The zone is case-sensitive (per the spec), so "domestic" (lowercase) is invalid. Any string not in the valid set should trigger an error.

#### Insured Input Partitions

| Partition | Value | Characteristic | Example Values | Expected Behavior |
|-----------|-------|-----------------|-----------------|-------------------|
| **P7: Not Insured** | false | No insurance | false | 1.5% insurance cost NOT added |
| **P8: Insured** | true | Insurance enabled | true | 1.5% of (base + surcharge) added |

**Rationale:** The boolean `insured` parameter creates a binary partition. When true, an additional 1.5% fee is applied; when false, it is not.

---

### 1.2 Boundary Value Analysis (BVA)

Boundary Value Analysis focuses on testing the edge cases at the borders of each partition, as bugs frequently occur at these boundaries (e.g., off-by-one errors with `<` vs `<=`).

#### Weight Boundaries

**Lower Boundary (around 0):**

| Boundary Value | Type | Rationale | Expected Result |
|---|---|---|---|
| 0 | Invalid (just at boundary) | Tests if the condition is `weight <= 0` vs `weight < 0` | Error |
| 0.1 | Valid (just above boundary) | First valid value; tests the exclusive lower bound | Valid calculation (base fee only) |

**Mid Boundary (Tier Transition at 10):**

| Boundary Value | Type | Rationale | Expected Result |
|---|---|---|---|
| 10 | Valid Standard tier (at boundary) | Tests if surcharge applies at exactly 10 kg | Valid (base fee only, no surcharge) |
| 10.1 | Valid Heavy tier (just above) | Tests if surcharge correctly triggers above 10 kg | Valid (base fee + $7.50 surcharge) |

**Upper Boundary (around 50):**

| Boundary Value | Type | Rationale | Expected Result |
|---|---|---|---|
| 50 | Valid (at upper limit) | Tests if the condition is `weight <= 50` vs `weight < 50` | Valid calculation (base + surcharge if applicable) |
| 50.1 | Invalid (just above limit) | Tests the exclusive upper bound | Error |

**Why These Boundaries Matter:** These boundaries test for common off-by-one errors. For example, if a developer mistakenly wrote `weight > 0 && weight < 50` instead of `weight > 0 && weight <= 50`, the test with weight=50 would fail, catching the bug.

---

### 1.3 Test Coverage Summary

By combining Equivalence Partitioning and BVA, our test suite covers:

- **4 weight partitions** (with boundaries explicitly tested)
- **2 zone partitions** (all 3 valid zones + multiple invalid examples)
- **2 insured partitions** (true/false)
- **Combination scenarios** (e.g., heavy + insured, minimal weight + insured, maximum weight + insured)
- **Total test cases designed: 37 distinct scenarios**

---

## Part 2: Test Implementation Details

### 2.1 Test Structure

The test file `shipping_v2_test.go` uses Go's **table-driven test** pattern, which is the idiomatic Go way to write multiple related test cases. The structure is:

```
testCases := []struct {
    name        string   // Descriptive test name
    weight      float64  // Input
    zone        string   // Input
    insured     bool     // Input
    expectedFee float64  // Expected output
    expectError bool     // Whether we expect an error
}
```

Each test case maps directly to one equivalence partition or boundary value we identified.

### 2.2 Key Test Cases Explained

**Invalid Weight Tests (P1):**
```
{"P1: Weight zero", 0, "Domestic", false, 0, true}
{"P1: Weight negative", -5, "International", false, 0, true}
```
These verify that weights ≤ 0 are rejected with an error.

**Standard Tier Tests (P2):**
```
{"P2: Standard package - Domestic", 5, "Domestic", false, 5.0, false}
```
- Input: 5 kg, Domestic zone, not insured
- Expected: $5.00 base fee only (no surcharge for weight ≤ 10)
- Calculation: 5 + 0 = 5.0

**Heavy Tier Tests (P3):**
```
{"P3: Heavy package - Domestic", 25, "Domestic", false, 12.5, false}
```
- Input: 25 kg, Domestic zone, not insured
- Expected: $5.00 base + $7.50 surcharge = $12.50
- Calculation: 5 + 7.5 = 12.5

**Insured Tests (P8):**
```
{"P8: Insured Heavy - Domestic", 25, "Domestic", true, 12.6875, false}
```
- Input: 25 kg, Domestic zone, insured
- Expected: ($5 + $7.50) × 1.015 = $12.50 × 1.015 = $12.6875
- Calculation: 12.5 + (12.5 × 0.015) = 12.5 + 0.1875 = 12.6875

**Boundary Tests:**
```
{"BVA: Weight at tier boundary - Standard (10)", 10, "Domestic", false, 5.0, false}
{"BVA: Weight at tier boundary - Heavy (10.1)", 10.1, "Domestic", false, 12.5, false}
```
These test the exact transition point: at 10 kg there's no surcharge, at 10.1 kg there is.

### 2.3 Floating-Point Comparison

Since floating-point arithmetic can introduce small rounding errors, the test uses a tolerance of 0.0001:

```go
const tolerance = 0.0001
if fee < tc.expectedFee-tolerance || fee > tc.expectedFee+tolerance {
    t.Errorf("Expected fee %.4f, but got %.4f", tc.expectedFee, fee)
}
```

This prevents false failures due to precision issues.

---

## Part 3: Testing Methodology Justification

### Why Specification-Based Testing?

1. **Independence from Implementation:** These tests validate the specification itself, not the code. If the code is refactored but the spec remains the same, the tests stay valid.

2. **Comprehensive Coverage:** By systematically analyzing the specification, we ensure no business rule is missed.

3. **Bug Prevention:** Focusing on partitions and boundaries catches the most common types of bugs (off-by-one errors, missing validation, etc.).

4. **Maintainability:** Table-driven tests are easy to read, modify, and extend when requirements change.

### How the Three Techniques Work Together

- **Equivalence Partitioning** provides broad coverage with minimal tests.
- **Boundary Value Analysis** adds precision at the critical edges.
- **Decision Table Testing** (implicit in our table-driven tests) ensures all combinations are covered.

---

## Part 4: Test Execution Results

**All 37 tests pass successfully.**

![image](./assets/Screenshot%20From%202025-11-21%2012-30-05.png)

![image](./assets/Screenshot%20From%202025-11-21%2012-30-29.png)
Test output summary:
- P1 (Invalid weight - too small): 2 tests, 2 pass
- P2 (Standard tier): 3 tests, 3 pass
- P3 (Heavy tier): 3 tests, 3 pass
- P4 (Invalid weight - too large): 2 tests, 2 pass
- P5 (Valid zones): 3 tests, 3 pass
- P6 (Invalid zones): 3 tests, 3 pass
- P7 (Not insured): 2 tests, 2 pass
- P8 (Insured): 4 tests, 4 pass
- BVA (Boundary values): 6 tests, 6 pass
- Combination scenarios: 3 tests, 3 pass

---
