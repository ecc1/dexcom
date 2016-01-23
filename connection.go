/*
Package dexcom provides functions to access a Dexcom CGM receiver
over a USB or BLE connection.

Based on the Python version at github.com/bewest/decoding-dexcom
*/
package dexcom

import (
	"log"
)

type Connection interface {
	Frame([]byte) []byte
	Send([]byte) error
	Receive([]byte) error
}

var conn Connection

// Open first attempts to open a USB connection;
// if that fails it tries a BLE connection.
func Open() error {
	var err error
	conn, err = OpenUSB()
	if err != nil {
		log.Println("USB:", err)
		conn, err = OpenBLE()
	}
	return err
}
