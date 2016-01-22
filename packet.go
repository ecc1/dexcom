package dexcom

import (
	"bytes"
	"fmt"
)

const (
	startOfMessage = 1
	minPacket      = 6
	maxPacket      = 1590
)

func marshalPacket(cmd Command, params ...[]byte) []byte {
	var buf bytes.Buffer
	buf.WriteByte(startOfMessage)
	data := []byte{}
	for _, p := range params {
		data = append(data, p...)
	}
	length := uint16(4 + len(data) + 2) // header, data, CRC-16
	buf.Write(MarshalUint16(length))
	buf.WriteByte(byte(cmd))
	buf.Write(data)
	buf.Write(crc16(buf.Bytes())) // append CRC-16
	return buf.Bytes()
}

func sendPacket(pkt []byte) error {
	return conn.Send(pkt)
}

func receivePacket() (cmd byte, data []byte, err error) {
	header := make([]byte, 4)
	if err = conn.Receive(header); err != nil {
		return
	}
	if header[0] != startOfMessage {
		err = fmt.Errorf("unexpected message header (%X)", header)
		return
	}
	cmd = header[3]
	length := UnmarshalUint16(header[1:3])
	if length < minPacket || length > maxPacket {
		err = fmt.Errorf("invalid packet length in header (%d)", length)
		return
	}
	n := length - minPacket
	if n > 0 {
		data = make([]byte, n)
		if err = conn.Receive(data); err != nil {
			return
		}
	}
	crc := make([]byte, 2)
	if err = conn.Receive(crc); err != nil {
		return
	}
	calc := crc16(append(header, data...))
	if !bytes.Equal(crc, calc) {
		err = fmt.Errorf("bad CRC (received %X, computed %X)", crc, calc)
		return
	}
	return
}

// Cmd creates a Dexcom packet with the given command and parameters,
// sends it to the device, and returns the response.
func Cmd(cmd Command, params ...[]byte) ([]byte, error) {
	pkt := marshalPacket(cmd, params...)
	err := sendPacket(pkt)
	if err != nil {
		return nil, err
	}
	ack, response, err := receivePacket()
	if err != nil {
		return nil, err
	}
	if ack != 1 {
		return nil, fmt.Errorf("unexpected ack (%X)", ack)
	}
	return response, nil
}
