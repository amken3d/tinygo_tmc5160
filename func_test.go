package tinygo_tmc5160

import "testing"

// Test speedToHz function
func TestSpeedToHz(t *testing.T) {
	tests := []struct {
		speedInternal int32
		expected      float32
	}{
		{10000, 7152.557617},  // Corrected expected result
		{50000, 35762.785156}, // Corrected expected result
		{-5000, 0},            // Corrected expected result
	}

	for _, test := range tests {
		t.Run("Testing speedToHz", func(t *testing.T) {
			result := speedToHz(uint32(test.speedInternal))
			if !(result == test.expected) {
				t.Errorf("speedToHz(%d) = %f; expected %f", test.speedInternal, result, test.expected)
			}
		})
	}
}

// Test speedFromHz function
func TestSpeedFromHz(t *testing.T) {
	tests := []struct {
		speedHz  float32
		expected int32
	}{
		{7152.55, 10000},  // Corrected expected result
		{35762.75, 50000}, // Corrected expected result
		{-3576.28, 0},     // Corrected expected result
	}

	for _, test := range tests {
		t.Run("Testing speedFromHz", func(t *testing.T) {
			result := speedFromHz(test.speedHz)
			if result != test.expected {
				t.Errorf("speedFromHz(%f) = %d; expected %d", test.speedHz, result, test.expected)
			}
		})
	}
}

// Test accelFromHz function
func TestAccelFromHz(t *testing.T) {
	tests := []struct {
		accelHz  float32
		expected int32
	}{
		{0.5, 73795},
		{1.0, 147590},
		{-0.25, -36897},
	}

	for _, test := range tests {
		t.Run("Testing accelFromHz", func(t *testing.T) {
			result := accelFromHz(test.accelHz)
			if result != test.expected {
				t.Errorf("accelFromHz(%f) = %d; expected %d", test.accelHz, result, test.expected)
			}
		})
	}
}

// Test thrsSpeedToTstep function
func TestThrsSpeedToTstep(t *testing.T) {
	tests := []struct {
		thrsSpeed uint32
		expected  uint32
	}{
		{1000, 16777},
		{10000, 1677},
		{1, 1048575},
	}

	for _, test := range tests {
		t.Run("Testing thrsSpeedToTstep", func(t *testing.T) {
			result := thrsSpeedToTstep(test.thrsSpeed)
			if result != test.expected {
				t.Errorf("thrsSpeedToTstep(%d) = %d; expected %d", test.thrsSpeed, result, test.expected)
			}
		})
	}
}
