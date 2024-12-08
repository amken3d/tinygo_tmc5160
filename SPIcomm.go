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
	spi machine.SPI
}

// NewSPIComm creates a new SPIComm instance.
func NewSPIComm(spi machine.SPI) *SPIComm {
	return &SPIComm{
		spi: spi,
	}
}

// Setup initializes the SPI communication with the TMC2209.
func (comm *SPIComm) Setup() error {
	// Check if SPI is initialized
	if comm.spi == (machine.SPI{}) {
		return CustomError("SPI not initialized")
	}

	// Configure the SPI interface with the desired settings
	err := comm.spi.Configure(machine.SPIConfig{
		LSBFirst: false,
		Mode:     3,
	})
	if err != nil {
		return CustomError("Failed to configure SPI")
	}

	// No built-in timeout in TinyGo, so timeout will be handled in the read/write methods
	return nil
}

// WriteRegister sends a register write command to the TMC5160.
func (comm *SPIComm) WriteRegister(register uint8, value uint32, driverIndex uint8) error {
	// Pass the register and value to the spiTransfer40 function to write to the device
	_, err := spiTransfer40(&comm.spi, register, value)
	if err != nil {
		return CustomError("Failed to write register")
	}
	return nil
}

// ReadRegister sends a register read command to the TMC5160.
func (comm *SPIComm) ReadRegister(register uint8, driverIndex uint8) (uint32, error) {
	// Send the register read request and get the response
	response, err := spiTransfer40(&comm.spi, register, 0)
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
