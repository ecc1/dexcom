package dexcom

import (
	"bytes"
	"fmt"
	"log"

	"github.com/ecc1/ble"
)

const (
	gattMTU = 20
)

var (
	// service
	receiverService = dexcomUUID(0xa0b1)

	// characteristics
	heartbeat      = dexcomUUID(0x2b18)
	authentication = dexcomUUID(0xacac)
	command        = dexcomUUID(0xb0cc)
	response       = dexcomUUID(0xb0cd)
	sendData       = dexcomUUID(0xb20a)
	receiveData    = dexcomUUID(0xb20b)

	authCode = []byte(serialNumber + "000000")
)

func dexcomUUID(id uint16) string {
	return "f0ac" + fmt.Sprintf("%04x", id) + "-ebfa-f96f-28da-076c35a521db"
}

type bleConn struct {
	txChar       ble.Characteristic
	inputChannel chan byte
}

func (conn *bleConn) Frame(data []byte) []byte {
	return append([]byte{1, 1}, data...)
}

func (conn *bleConn) Send(data []byte) error {
	for {
		n := len(data)
		if n == 0 {
			return nil
		}
		if n > gattMTU {
			n = gattMTU
		}
		err := conn.txChar.WriteValue(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
}

func (conn *bleConn) Receive(data []byte) error {
	for i := 0; i < len(data); i++ {
		b, ok := <-conn.inputChannel
		if !ok {
			return fmt.Errorf("input channel closed")
		}
		data[i] = b
	}
	return nil
}

func authenticate(device ble.Device, objects *ble.ObjectCache) error {
	auth, err := objects.GetCharacteristic(authentication)
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

func OpenBLE() (Connection, error) {
	objects, err := ble.ManagedObjects()
	if err != nil {
		log.Fatal(err)
	}

	device, err := objects.Discover(0, receiverService)
	if err != nil {
		log.Fatal(err)
	}

	if !device.Connected() {
		err = device.Connect()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s: connected\n", device.Name())
	} else {
		log.Printf("%s: already connected\n", device.Name())
	}

	if !device.Paired() {
		err = device.Pair()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s: paired\n", device.Name())
	} else {
		log.Printf("%s: already paired\n", device.Name())
	}

	err = objects.Update()
	if err != nil {
		log.Fatal(err)
	}

	err = authenticate(device, objects)
	if err != nil {
		log.Fatal(err)
	}

	// We need to enable heartbeat notifications
	// or else we won't get any receiveData responses.
	hbChar, err := objects.GetCharacteristic(heartbeat)
	if err != nil {
		log.Fatal(err)
	}
	err = hbChar.HandleNotify(func(data []byte) {})
	if err != nil {
		log.Fatal(err)
	}

	rxChar, err := objects.GetCharacteristic(receiveData)
	if err != nil {
		log.Fatal(err)
	}
	inputChannel := make(chan byte, 1600)
	err = rxChar.HandleNotify(func(data []byte) {
		for _, b := range data {
			inputChannel <- b
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	txChar, err := objects.GetCharacteristic(sendData)
	if err != nil {
		log.Fatal(err)
	}

	return &bleConn{
		txChar:       txChar,
		inputChannel: inputChannel,
	}, nil
}
