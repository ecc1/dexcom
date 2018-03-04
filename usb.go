package dexcom

import (
	"log"

	"github.com/ecc1/serial"
)

const (
	// USB IDs for the Dexcom G4 receiver.
	dexcomVendor  = 0x22a3
	dexcomProduct = 0x0047
)

type usbConn serial.Port

// OpenUSB opens the USB serial device for a Dexcom G4 receiver.
func OpenUSB() (Connection, error) {
	device, err := serial.FindUSB(dexcomVendor, dexcomProduct)
	if err != nil {
		_, notFound := err.(serial.DeviceNotFoundError)
		if !notFound {
			log.Print(err)
		}
		return nil, err
	}
	port, err := serial.Open(device, 115200)
	return (*usbConn)(port), err
}

// Send writes data over the USB connection.
func (conn *usbConn) Send(data []byte) error {
	return (*serial.Port)(conn).Write(data)
}

// Receive reads data from the USB connection.
func (conn *usbConn) Receive(data []byte) error {
	return (*serial.Port)(conn).Read(data)
}

// Close closes the USB connection.
func (conn *usbConn) Close() {
	_ = (*serial.Port)(conn).Close()
}
