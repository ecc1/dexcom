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
	port *usbserial.Port
}

// Open locates the USB device for a Dexcom CGM receiver and opens it.
func OpenUSB() (Connection, error) {
	port, err := usbserial.Open(dexcomVendor, dexcomProduct)
	if err != nil {
		return nil, err
	}
	return &usbConn{port: port}, nil
}

func (conn usbConn) Send(data []byte) error {
	return conn.port.Write(data)
}

func (conn usbConn) Receive(data []byte) error {
	return conn.port.Read(data)
}

func (conn usbConn) Close() {
	conn.port.Close()
}
