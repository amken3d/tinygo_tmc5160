package tinygo_tmc5160

import (
	math "github.com/orsinium-labs/tinymath"
	"log"
)

// Define constants
const _uStepCount uint32 = 256
const _motorAngle float32 = 1.8

var _ustepPerRev = uint32(360/_motorAngle) * _uStepCount
var _maxUstepPerTref = _ustepPerRev * _tRef

const _fclkMhz uint32 = 12 // Example clock frequency (12 MHz)
var _tRef = 16777216 / _fclkMhz

func roundFloat(val float32, precision uint) float32 {
	ratio := math.PowF(10, float32(precision))
	return math.Round(val*ratio) / ratio
}

func UstepPerSec(speedInternal uint32) uint32 {
	return speedInternal / _tRef
}

// speedToHz function using Stepper struct
func (stepper *Stepper) speedToHz(speedInternal uint32) float32 {
	// Check if the input speed is within the valid range, use VMax from the Stepper struct
	if speedInternal > uint32(715828) { // Limit speed based on your conditions
		return 0.0 // Return 0 for invalid speeds (e.g., negative or too large)
	}

	// Use the GearRatio and Velocity from the Stepper struct for calculation (you can modify this part as needed)
	log.Printf("Speed: %d, Gear Ratio: %f", speedInternal, stepper.GearRatio)
	return float32(speedInternal) * float32(_fclkMhz*1000000) / (2 * (1 << 23))
}

// Convert real-world frequency (Hz) to internal speed (v[5160A]) using the inverse formula
func speedFromHz(speedHz float32) int32 {
	if speedHz < 0 {
		return 0
	}

	// Applying the inverse formula
	return int32(speedHz*(2*float32(1<<23))/float32(_fclkMhz*1000000) + 1)
}

// Convert real-world acceleration (Hz/s) to internal acceleration
func accelFromHz(accelHz float32) int32 {
	if accelHz == 0 {
		return 0
	}

	// Apply proper scaling for acceleration conversion
	return int32(accelHz * float32(_fclkMhz*_fclkMhz) * 1000000 / (512.0 * 256.0) / float32(1<<24) * float32(_uStepCount))
}

// Convert threshold speed (Hz) to internal TSTEP value
func thrsSpeedToTstep(thrsSpeed uint32) uint32 {
	if thrsSpeed < 0 {
		return 0
	}
	r := (16777216 / thrsSpeed) * (_uStepCount / 256)
	// Correct scaling for threshold speed
	return constrain(r, 0, 1048575)
}

// Constrain function to limit values to a specific range (now works with int32 or uint32)
func constrain(value, min, max uint32) uint32 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

// Function returns the Max Steps per Seconds for a given value in rpm
func (stepper *Stepper) MaxStepsPerSecond(value float32) {
	//t := 16777216 / float32(stepper.Fclk)
	//vmax :=
	//rps := 16777216
	return
}
