package main

import (
	"github.com/amken3d/tinygo_tmc5160/tmc5160"
	"machine"
	"time"
)

// Example usage to initialize a custom stepper with user-defined values
func main() {
	time.Sleep(8 * time.Second)
	println("Entering Main")

	spi := machine.SPI1
	csPin1 := machine.GPIO13
	enn1 := machine.GPIO18
	csPin1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	enn1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	err := spi.Configure(machine.SPIConfig{
		SCK:       machine.GPIO10,
		SDO:       machine.GPIO11,
		SDI:       machine.GPIO12,
		Frequency: 12000000,
		Mode:      3,
		LSBFirst:  false,
	})
	if err != nil {
		println("Failed to configure SPI")
		return
	}
	// Map the Driver addresses to their corresponding CS pins
	csPins := map[uint8]machine.Pin{
		0: csPin1, // Driver with address 0x01 uses csPin1
	}
	comm := tmc5160.NewSPIComm(*spi, csPins)
	println("Comms Registered ")

	stepper := tmc5160.NewDefaultStepper()
	driver1 := tmc5160.NewDriver(comm, 0, enn1, stepper)
	enn1.Low()
	err = driver1.Dump_TMC()
	if err != nil {
		return
	}
	//// Read the registers and log their values
	//GCONF := tmc5160.NewGCONF()
	//G, _ := driver1.ReadRegister(tmc5160.GCONF)
	//GCONF.Unpack(G)
	//println("GCONF{Recalibrate:", GCONF.Recalibrate, ", MultiStepFilt:", GCONF.MultistepFilt, "}")
	//IOIN := tmc5160.NewIOIN()
	//driver1.ReadRegister(tmc5160.IOIN)
	//i, _ := driver1.ReadRegister(tmc5160.IOIN)
	//IOIN.Unpack(i)
	//println("IOIN{SD Mode:", IOIN.SdMode, ", Version", tmc5160.ToHex(uint32(IOIN.Version)), "}")

}

//	func (d driver) GetStatus( driverAddress uint8) (uint32, error) {
//		// Assert the chip select pin (set CS low to start communication)
//		csPin, exists := comm.CsPins[driverAddress]
//		if !exists {
//			return 0, CustomError("Invalid driver address")
//		}
//		// Prepare the command to read the status register
//		tx := []byte{
//			GCONF, // Address of the GCONF register
//			0x00,  // Dummy byte to send
//			0x00,  // Dummy byte to send
//			0x00,  // Dummy byte to send
//			0x00,  // Dummy byte to send
//		}
//
//		// Prepare a buffer for the response (5 bytes)
//		rx := make([]byte, 5)
//
//		csPin.Low()
//
//		// Send the register read request and get the response
//		err := comm.spi.Tx(tx, rx)
//
//		// Deassert the chip select pin (set CS high to end communication)
//		csPin.High()
//
//		if err != nil {
//			return 0, CustomError("Failed to read register")
//		}
//		return uint32(rx[1]), nil
//	}
// String returns a formatted string displaying the unpacked values.
