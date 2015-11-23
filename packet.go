package dexcom

import (
	"bytes"
	"fmt"

	"github.com/ecc1/usbserial"
)

const (
	startOfMessage       = 1
	minPacket = 6
	maxPacket = 1590
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

func sendPacket(port usbserial.Port, pkt []byte) error {
	return port.Write(pkt)
}

func receivePacket(port usbserial.Port) (cmd byte, data []byte, err error) {
	header := make([]byte, 4)
	if err = port.Read(header); err != nil {
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
		if err = port.Read(data); err != nil {
			return
		}
	}
	crc := make([]byte, 2)
	if err = port.Read(crc); err != nil {
		return
	}
	calc := crc16(append(header, data...))
	if !bytes.Equal(crc, calc) {
		err = fmt.Errorf("bad CRC (received %X, computed %X)", crc, calc)
		return
	}
	return
}
