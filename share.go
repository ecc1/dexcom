package dexcom

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/ecc1/ble"
)

type bleConn struct {
	tx ble.Characteristic
	rx chan byte
}

// Dexcom G4 Share expects each BLE message to start with two 01 bytes.
func (conn *bleConn) Frame(data []byte) []byte {
	return append([]byte{1, 1}, data...)
}

const (
	// maximum size of writes to GATT characteristics
	gattMTU = 20
)

func (conn *bleConn) Send(data []byte) error {
	for {
		n := len(data)
		if n == 0 {
			return nil
		}
		if n > gattMTU {
			n = gattMTU
		}
		err := conn.tx.WriteValue(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
}

func (conn *bleConn) Receive(data []byte) error {
	for i := 0; i < len(data); i++ {
		b, ok := <-conn.rx
		if !ok {
			return fmt.Errorf("input channel closed")
		}
		data[i] = b
	}
	return nil
}

var (
	receiverService = dexcomUUID(0xa0b1)
)

func connect() error {
	device, err := ble.Discover(time.Minute, receiverService)
	if err != nil {
		return err
	}
	if !device.Connected() {
		err = device.Connect()
		if err != nil {
			return err
		}
		log.Printf("%s: connected\n", device.Name())
	} else {
		log.Printf("%s: already connected\n", device.Name())
	}
	if !device.Paired() {
		err = device.Pair()
		if err != nil {
			return err
		}
		log.Printf("%s: paired\n", device.Name())
	} else {
		log.Printf("%s: already paired\n", device.Name())
	}
	err = ble.Update()
	if err != nil {
		return err
	}
	return authenticate(device)
}

var (
	authentication = dexcomUUID(0xacac)
	authCode       = []byte(serialNumber + "000000")
)

func authenticate(device ble.Device) error {
	auth, err := ble.GetCharacteristic(authentication)
	if err != nil {
		return err
	}
	data, err := auth.ReadValue()
	if err != nil {
		return err
	}
	if bytes.Equal(data, authCode) {
		log.Printf("%s: already authenticated\n", device.Name())
		return nil
	}
	err = auth.WriteValue(authCode)
	if err != nil {
		return err
	}
	log.Printf("%s: authenticated\n", device.Name())
	return nil
}

var (
	heartbeat   = dexcomUUID(0x2b18)
	sendData    = dexcomUUID(0xb20a)
	receiveData = dexcomUUID(0xb20b)
)

func OpenBLE() (Connection, error) {
	err := connect()
	if err != nil {
		return nil, err
	}

	// We need to enable heartbeat notifications
	// or else we won't get any receiveData responses.
	err = ble.HandleNotify(heartbeat, func(data []byte) {})
	if err != nil {
		return nil, err
	}

	rx := make(chan byte, 1600)
	err = ble.HandleNotify(receiveData, func(data []byte) {
		for _, b := range data {
			rx <- b
		}
	})
	if err != nil {
		return nil, err
	}

	tx, err := ble.GetCharacteristic(sendData)
	if err != nil {
		return nil, err
	}

	return &bleConn{
		tx: tx,
		rx: rx,
	}, nil
}

func dexcomUUID(id uint16) string {
	return fmt.Sprintf("f0ac%04x-ebfa-f96f-28da-076c35a521db", id)
}
