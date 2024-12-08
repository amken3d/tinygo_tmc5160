package tinygo_tmc5160

type TMC5160 struct {
	comm    RegisterComm
	address uint8
}

func NewTMC5160(comm RegisterComm, address uint8) *TMC5160 {
	return &TMC5160{
		comm:    comm,
		address: address,
	}
}

// WriteRegister sends a register write command to the TMC5160.
func (driver *TMC5160) WriteRegister(reg uint8, value uint32) error {
	if driver.comm == nil {
		return CustomError("communication interface not set")
	}
	// Use the communication interface (RegisterComm) to write the register
	return driver.comm.WriteRegister(reg, value, driver.address)
}

// ReadRegister sends a register read command to the TMC5160 and returns the read value.
func (driver *TMC5160) ReadRegister(reg uint8) (uint32, error) {
	if driver.comm == nil {
		return 0, CustomError("communication interface not set")
	}
	// Use the communication interface (RegisterComm) to read the register
	return driver.comm.ReadRegister(reg, driver.address)
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
