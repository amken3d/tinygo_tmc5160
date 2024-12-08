package tinygo_tmc5160

import (
	"machine"
)

// CustomError is a lightweight error type used for TinyGo compatibility.
type CustomError string

func (e CustomError) Error() string {
	return string(e)
}

// SPIComm implements RegisterComm for SPI-based communication
type SPIComm struct {
	spi    machine.SPI
	csPins map[uint8]machine.Pin // Map to store CS pin for each TMC5160 by its address
}

// NewSPIComm creates a new SPIComm instance.
func NewSPIComm(spi machine.SPI, csPins map[uint8]machine.Pin) *SPIComm {
	return &SPIComm{
		spi:    spi,
		csPins: csPins,
	}
}

// Setup initializes the SPI communication with the TMC5160 and configures all CS pins.
func (comm *SPIComm) Setup() error {
	// Check if SPI is initialized
	if comm.spi == (machine.SPI{}) {
		return CustomError("SPI not initialized")
	}

	// Configure all CS pins (make them output and set them high)
	for _, csPin := range comm.csPins {
		csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		csPin.High() // Set all CS pins high initially
	}

	// Configure the SPI interface with the desired settings
	err := comm.spi.Configure(machine.SPIConfig{
		LSBFirst: false,
		Mode:     3,
	})
	if err != nil {
		return CustomError("Failed to configure SPI")
	}

	return nil
}

// WriteRegister sends a register write command to the TMC5160.
func (comm *SPIComm) WriteRegister(register uint8, value uint32, driverAddress uint8) error {
	// Assert the chip select pin (set CS low to start communication)
	csPin, exists := comm.csPins[driverAddress]
	if !exists {
		return CustomError("Invalid driver address")
	}
	csPin.Low()

	// Pass the register and value to the spiTransfer40 function to write to the device
	_, err := spiTransfer40(&comm.spi, register, value)

	// Deassert the chip select pin (set CS high to end communication)
	csPin.High()

	if err != nil {
		return CustomError("Failed to write register")
	}
	return nil
}

// ReadRegister sends a register read command to the TMC5160.
func (comm *SPIComm) ReadRegister(register uint8, driverAddress uint8) (uint32, error) {
	// Assert the chip select pin (set CS low to start communication)
	csPin, exists := comm.csPins[driverAddress]
	if !exists {
		return 0, CustomError("Invalid driver address")
	}
	csPin.Low()

	// Send the register read request and get the response
	response, err := spiTransfer40(&comm.spi, register, 0)

	// Deassert the chip select pin (set CS high to end communication)
	csPin.High()

	if err != nil {
		return 0, CustomError("Failed to read register")
	}
	return response, nil
}

func spiTransfer40(spi *machine.SPI, register uint8, txData uint32) (uint32, error) {
	// Prepare the 5-byte buffer for transmission (1 byte address + 4 bytes data)
	tx := []byte{
		register,           // Address byte
		byte(txData >> 24), // Upper 8 bits of data
		byte(txData >> 16), // Middle 8 bits of data
		byte(txData >> 8),  // Next 8 bits of data
		byte(txData),       // Lower 8 bits of data
	}

	// Prepare the receive buffer (5 bytes) for the response
	rx := make([]byte, 5)

	// Perform the SPI transaction
	err := spi.Tx(tx, rx)
	if err != nil {
		return 0, err
	}

	// Combine the received bytes into a 32-bit response, ignore the address byte
	rxData := uint32(rx[1])<<24 | uint32(rx[2])<<16 | uint32(rx[3])<<8 | uint32(rx[4])

	return rxData, nil
}
