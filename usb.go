package dexcom

import (
	"github.com/ecc1/usbserial"
)

const (
	// USB IDs for the Dexcom G4 receiver.
	dexcomVendor  = 0x22a3
	dexcomProduct = 0x0047
)

type usbConn struct {
	*usbserial.Port
}

// OpenUSB opens the USB serial device for a Dexcom G4 receiver.
func OpenUSB() (Connection, error) {
	port, err := usbserial.Open(dexcomVendor, dexcomProduct)
	return &usbConn{port}, err
}

// Send writes data over the USB connection.
func (conn *usbConn) Send(data []byte) error {
	return conn.Write(data)
}

// Receive reads data from the USB connection.
func (conn *usbConn) Receive(data []byte) error {
	return conn.Read(data)
}

// Close closes the USB connection.
func (conn *usbConn) Close() {
	_ = conn.Port.Close()
}
