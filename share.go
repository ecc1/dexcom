package dexcom

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ecc1/ble"
)

type bleConn struct {
	*ble.Connection
	tx ble.Characteristic
	rx chan byte
}

const (
	addressEnvVar = "DEXCOM_G4_RECEIVER_ADDRESS"
	authEnvVar    = "DEXCOM_CGM_ID"

	discoveryTimeout = 10 * time.Second
	receiveTimeout   = 5 * time.Second
)

var (
	receiverAddress ble.Address

	receiverService = dexcomUUID(0xa0b1)
	authentication  = dexcomUUID(0xacac)
	heartbeat       = dexcomUUID(0x2b18)
	sendData        = dexcomUUID(0xb20a)
	receiveData     = dexcomUUID(0xb20b)

	authCode []byte
)

// Send writes data over the BLE connection.
func (conn *bleConn) Send(data []byte) error {
	// Dexcom G4 Share expects each BLE message to start with two 01 bytes.
	data = append([]byte{0x01, 0x01}, data...)
	for {
		n := len(data)
		if n == 0 {
			return nil
		}
		if n > ble.GATTMTU {
			n = ble.GATTMTU
		}
		err := conn.tx.WriteValue(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
}

// Receive reads data from the BLE connection.
func (conn *bleConn) Receive(data []byte) error {
	for i := 0; i < len(data); i++ {
		select {
		case b := <-conn.rx:
			data[i] = b
		case <-time.After(receiveTimeout):
			return fmt.Errorf("BLE receive timeout")
		}
	}
	return nil
}

func initAddress() error {
	if len(receiverAddress) != 0 {
		return nil
	}
	addr := os.Getenv(addressEnvVar)
	if len(addr) == 0 {
		return fmt.Errorf("%s environment variable is not set", addressEnvVar)
	}
	if !ble.ValidAddress(addr) {
		return fmt.Errorf("%s: %q is not a valid MAC address", addressEnvVar, addr)
	}
	receiverAddress = ble.Address(strings.ToUpper(addr))
	return nil
}

func connect(conn *ble.Connection) error {
	device, err := findDevice(conn)
	if err != nil {
		return err
	}
	reauth := false
	if !device.Connected() {
		reauth = true
		err = device.Connect()
		if err != nil {
			return err
		}
	}
	if !device.Paired() {
		reauth = true
		err = device.Pair()
		if err != nil {
			return err
		}
	}
	err = conn.Update()
	if err != nil {
		return err
	}
	return authenticate(device, reauth)
}

func findDevice(conn *ble.Connection) (ble.Device, error) {
	err := initAddress()
	if err != nil {
		return nil, err
	}
	device, err := conn.GetDeviceByAddress(receiverAddress)
	if err == nil && device.Connected() {
		return device, nil
	}
	// Remove device to avoid "Software caused connection abort" error.
	if err == nil {
		adapter, err := conn.GetAdapter()
		if err != nil {
			return nil, err
		}
		err = adapter.RemoveDevice(device)
		if err != nil {
			return nil, err
		}
	}
	return conn.Discover(discoveryTimeout, receiverAddress, receiverService)
}

func initAuthCode() error {
	if len(authCode) != 0 {
		return nil
	}
	id := os.Getenv(authEnvVar)
	switch len(id) {
	case 10:
		authCode = []byte(id + "000000")
		return nil
	case 0:
		return fmt.Errorf("%s environment variable is not set", authEnvVar)
	default:
		return fmt.Errorf("%s environment variable must be 2 letters followed by 8 digits", authEnvVar)
	}
}

func authenticate(device ble.Device, reauth bool) error {
	err := initAuthCode()
	if err != nil {
		return err
	}
	var auth ble.Characteristic
	auth, err = device.Conn().GetCharacteristic(authentication)
	if err != nil {
		return err
	}
	if !reauth {
		var data []byte
		data, err = auth.ReadValue()
		if err != nil {
			return err
		}
		if bytes.Equal(data, authCode) {
			return nil
		}
	}
	log.Printf("%s: authenticating", device.Name())
	err = auth.WriteValue(authCode)
	if err != nil {
		return err
	}
	return nil
}

// OpenBLE makes a BLE connection to a Dexcom G4 Share receiver.
func OpenBLE() (Connection, error) {
	conn, err := ble.Open()
	if err != nil {
		return nil, err
	}
	err = connect(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	// We need to enable heartbeat notifications
	// or else we won't get any receiveData responses.
	err = conn.HandleNotify(heartbeat, func(data []byte) {})
	if err != nil {
		conn.Close()
		return nil, err
	}
	rx := make(chan byte, 1600)
	err = conn.HandleNotify(receiveData, func(data []byte) {
		for _, b := range data {
			rx <- b
		}
	})
	if err != nil {
		conn.Close()
		return nil, err
	}
	tx, err := conn.GetCharacteristic(sendData)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &bleConn{Connection: conn, tx: tx, rx: rx}, nil
}

func dexcomUUID(id uint16) string {
	return fmt.Sprintf("f0ac%04x-ebfa-f96f-28da-076c35a521db", id)
}
