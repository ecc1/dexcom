package dexcom

import (
	"fmt"

	"github.com/ecc1/usbserial"
)

const (
	// hex IDs that should match the idVendor and idProduct files
	// in /sys/bus/usb/devices
	dexcomVendor  = "22a3"
	dexcomProduct = "0047"
)

// A Device represents an open Dexcom CGM receiver.
type Device struct {
	port usbserial.Port
}

// Open locates the USB device for a Dexcom CGM receiver and opens it.
func Open() (Device, error) {
	port, err := usbserial.Open(dexcomVendor, dexcomProduct)
	return Device{port: port}, err
}

// Cmd creates a Dexcom packet with the given command and parameters,
// sends it to the device, and returns the response.
func (dev Device) Cmd(cmd Command, params ...[]byte) ([]byte, error) {
	pkt := marshalPacket(cmd, params...)
	err := sendPacket(dev.port, pkt)
	if err != nil {
		return nil, err
	}
	ack, response, err := receivePacket(dev.port)
	if err != nil {
		return nil, err
	}
	if ack != 1 {
		return nil, fmt.Errorf("unexpected ack (%X)", ack)
	}
	return response, nil
}
