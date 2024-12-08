package tinygo_tmc5160

import "machine"

const maxVMAX = 8388096

// PowerStageParameters represents the power stage parameters
type PowerStageParameters struct {
	drvStrength uint8
	bbmTime     uint8
	bbmClks     uint8
}

// MotorParameters represents the motor parameters
type MotorParameters struct {
	globalScaler   uint16
	ihold          uint8
	irun           uint8
	iholddelay     uint8
	pwmGradInitial uint16
	pwmOfsInitial  uint16
	freewheeling   uint8
}

// MotorDirection defines motor direction constants
type MotorDirection uint8

const (
	Clockwise MotorDirection = iota
	CounterClockwise
)

type StepperAngle float32
type Microstepping uint16

const (
	// Common stepper motor angles
	StepAngle_1_8  StepperAngle = 1.8
	StepAngle_0_9  StepperAngle = 0.9
	StepAngle_0_72 StepperAngle = 0.72
	StepAngle_1_2  StepperAngle = 1.2
	StepAngle_0_48 StepperAngle = 0.48

	// Common microstepping options
	Step_1   Microstepping = 1
	Step_2   Microstepping = 2
	Step_4   Microstepping = 4
	Step_8   Microstepping = 8
	Step_16  Microstepping = 16
	Step_32  Microstepping = 32
	Step_64  Microstepping = 64
	Step_128 Microstepping = 128
	Step_256 Microstepping = 256
)

const (
	DefaultAngle     StepperAngle = StepAngle_1_8
	DefaultGearRatio float32      = 1.0
	DefaultVSupply   float32      = 12.0
	DefaultRCoil     float32      = 1.2
	DefaultLCoil     float32      = 0.005
	DefaultIPeak     float32      = 2.0
	DefaultRSense    float32      = 0.1
	DefaultFclk      uint8        = 12
	DefaultStep_256               = 256
)

type Stepper struct {
	Angle       StepperAngle
	GearRatio   float32
	VelocitySPS float32 //  Velocity in Steps per sec
	VSupply     float32
	RCoil       float32
	LCoil       float32
	IPeak       float32
	RSense      float32
	MSteps      Microstepping
	Fclk        uint8 //Clock in Mhz

}

// NewStepper function initializes a Stepper with default values used for testing
func NewDefaultStepper() Stepper {
	return Stepper{
		Angle:       StepAngle_1_8, // Default to 1.8 degrees
		GearRatio:   1.0,           // Default to no reduction (1:1)
		VelocitySPS: 1000.0,        // Default velocity in steps per second
		VSupply:     12.0,          // Default 12V supply
		RCoil:       1.2,           // Default coil resistance (1.2 ohms)
		LCoil:       0.005,         // Default coil inductance (5 mH)
		IPeak:       2.0,           // Default peak current (2A)
		RSense:      0.1,           // Default sense resistance (0.1 ohms)
		MSteps:      Step_16,       // Default 16 Microsteps
		Fclk:        DefaultFclk,
	}
}

// NewStepper initializes a Stepper with user-defined values
func NewStepper(angle StepperAngle, gearRatio, velocitySPS, vSupply, rCoil, lCoil, iPeak, rSense float32, mSteps Microstepping, fclk uint8) Stepper {
	return Stepper{
		Angle:       angle,       // User-defined stepper angle (e.g., StepAngle_1_8)
		GearRatio:   gearRatio,   // User-defined gear ratio
		VelocitySPS: velocitySPS, // User-defined velocity in steps per second
		VSupply:     vSupply,     // User-defined supply voltage
		RCoil:       rCoil,       // User-defined coil resistance
		LCoil:       lCoil,       // User-defined coil inductance
		IPeak:       iPeak,       // User-defined peak current
		RSense:      rSense,      // User-defined sense resistance
		MSteps:      mSteps,      // User-defined microstepping setting
		Fclk:        fclk,        // User-defined clock frequency in MHz

	}
}

type TMC5160 struct {
	spi     machine.SPI
	csPin   machine.Pin
	stepper Stepper
}

func NewTMC5160(spi machine.SPI, csPin machine.Pin, stepper Stepper) *TMC5160 {
	return &TMC5160{
		spi:     spi,
		csPin:   csPin,
		stepper: stepper,
	}
}

// Enable Setup configures the TMC429 and SPI communication
func (t *TMC5160) Enable() {
	// Configure CS pin
	t.csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	t.csPin.Low()
}

// Disable Setup configures the TMC429 and SPI communication
func (t *TMC5160) Disable() {
	// Configure CS pin
	t.csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	t.csPin.High()
}

//
//// Begin initializes the TMC5160 driver with power and motor parameters
//func Begin(powerParams PowerStageParameters, motorParams MotorParameters, stepperDirection MotorDirection) bool {
//	// Clear the reset and charge pump undervoltage flags
//	gstat := NewGSTAT()
//	gstat.Reset = true
//	gstat.UvCp = true
//	writeRegister(GSTAT, gstat.Pack())
//
//	// Configure driver settings
//	drvConf := NewDRV_CONF()
//	drvConf.DrvStrength = constrain(powerParams.drvStrength, 0, 3)
//	drvConf.BBMTime = constrain(powerParams.bbmTime, 0, 24)
//	drvConf.BBMClks = constrain(powerParams.bbmClks, 0, 15)
//	writeRegister(DRV_CONF, drvConf.Pack())
//
//	// Set global scaler
//	writeRegister(GLOBAL_SCALER, uint32(constrain(uint8(motorParams.globalScaler), 32, 256)))
//
//	// Set initial currents and delay
//	iholdrun := NewIHOLD_IRUN()
//	iholdrun.Ihold = constrain(motorParams.ihold, 0, 31)
//	iholdrun.Ihold = constrain(motorParams.irun, 0, 31)
//	iholdrun.IholdDelay = 7
//	writeRegister(IHOLD_IRUN, iholdrun.Pack())
//
//	// Set PWM configuration values
//	pwmconf := NewPWMCONF()
//	writeRegister(PWMCONF, 0xC40C001E) // Reset default value pwm_ofs = 196,pwm_grad = 12,pwm_freq = 0, pwm_autoscale = false, pwm_autograd = false,freewheel = 3
//	pwmconf.PwmAutoscale = false       // Temporarily set to false for setting OFS and GRAD values
//	if _fclk > DEFAULT_F_CLK {
//		pwmconf.PwmFreq = 0
//	} else {
//		pwmconf.PwmFreq = 0b01 // Recommended: 35kHz with internal 12MHz clock
//	}
//	pwmconf.PwmGrad = uint8(motorParams.pwmGradInitial)
//	pwmconf.PwmOfs = uint8(motorParams.pwmOfsInitial)
//	pwmconf.Freewheel = motorParams.freewheeling
//	writeRegister(PWMCONF, pwmconf.Pack())
//
//	// Enable PWM auto-scaling and gradient adjustment
//	pwmconf.PwmAutoscale = true
//	pwmconf.PwmAutograd = true
//	writeRegister(PWMCONF, pwmconf.Pack())
//
//	// Recommended chop configuration settings
//	_chopConf := NewCHOPCONF()
//	_chopConf.Toff = 5
//	_chopConf.Tbl = 2
//	_chopConf.HstrtTfd = 4
//	_chopConf.HendOffset = 0
//	writeRegister(CHOPCONF, _chopConf.Pack())
//
//	// Use position mode
//	setRampMode(POSITIONING_MODE)
//
//	// Set StealthChop PWM mode and shaft direction
//	gconf := NewGCONF()
//	gconf.EnPwmMode = true // Enable stealthChop PWM mode
//	gconf.Shaft = stepperDirection == Clockwise
//	writeRegister(GCONF, gconf.Pack())
//
//	// Set default start, stop, threshold speeds
//	setRampSpeeds(0, 0.1, 0) // Start, stop, threshold speeds
//
//	// Set default D1 (must not be = 0 in positioning mode even with V1=0)
//	writeRegister(D_1, 100)
//
//	return false
//}
//
//// RampMode represents the mode of operation for the ramp generator
//type RampMode uint8
//
//const (
//	PositioningMode RampMode = iota
//	VelocityMode
//	HoldMode
//)
//
//// setRampMode sets the ramp mode for the TMC5160
//func setRampMode(mode RampMode) {
//	var rampModeValue uint32
//
//	switch mode {
//	case PositioningMode:
//		// Write the value for POSITIONING_MODE to RAMPMODE register
//		rampModeValue = POSITIONING_MODE
//	case VelocityMode:
//		// Write the value for VELOCITY_MODE_POS to RAMPMODE register
//		// First, set max speed to 0 as suggested in the original C++ code
//		setMaxSpeed(0)
//		rampModeValue = VELOCITY_MODE_POS
//	case HoldMode:
//		// Write the value for HOLD_MODE to RAMPMODE register
//		rampModeValue = HOLD_MODE
//	}
//
//	// Write the calculated value to the RAMPMODE register
//	writeRegister(RAMPMODE, rampModeValue)
//
//	// Optionally update the current ramp mode variable
//	_currentRampMode = mode
//}
//
//// setMaxSpeed sets the maximum speed to 0 (placeholder function)
//func setMaxSpeed(speed uint32) {
//	// This is a placeholder function that sets the speed
//	// Implement the actual logic to set the maximum speed register value
//	println("Max Speed set to:", speed)
//}
