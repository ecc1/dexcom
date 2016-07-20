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

func marshalPacket(cmd Command, data []byte) []byte {
	buf := bytes.Buffer{}
	buf.WriteByte(startOfMessage)
	length := uint16(4 + len(data) + 2) // header, data, CRC-16
	buf.Write(MarshalUint16(length))
	buf.WriteByte(byte(cmd))
	buf.Write(data)
	body := buf.Bytes()
	buf.Write(MarshalUint16(crc16(body))) // append CRC-16
	return buf.Bytes()
}

func sendPacket(pkt []byte) error {
	return conn.Send(conn.Frame(pkt))
}

func receivePacket() (cmd byte, data []byte, err error) {
	header := make([]byte, 4)
	err = conn.Receive(header)
	if err != nil {
		return
	}
	if header[0] != startOfMessage {
		err = fmt.Errorf("unexpected message header (% X)", header)
		return
	}
	cmd = header[3]
	length := UnmarshalUint16(header[1:3])
	if length < minPacket || length > maxPacket {
		err = fmt.Errorf("invalid packet length (%d) in header (% X)", length, header)
		return
	}
	n := length - minPacket
	if n > 0 {
		data = make([]byte, n)
		err = conn.Receive(data)
		if err != nil {
			return
		}
	}
	crcBuf := make([]byte, 2)
	err = conn.Receive(crcBuf)
	if err != nil {
		return
	}
	crc := UnmarshalUint16(crcBuf)
	body := append(header, data...)
	calc := crc16(body)
	if crc != calc {
		err = CrcError{
			Kind:     "packet",
			Received: crc,
			Computed: calc,
			Context:  nil,
			Data:     body,
		}
		return
	}
	return
}

// Cmd creates a Dexcom packet with the given command and parameters,
// sends it to the device, and returns the response.
func Cmd(cmd Command, params ...byte) ([]byte, error) {
	pkt := marshalPacket(cmd, params)
	err := sendPacket(pkt)
	if err != nil {
		return nil, err
	}
	ack, response, err := receivePacket()
	if err != nil {
		return nil, err
	}
	if ack != 1 {
		return nil, fmt.Errorf("unexpected ack (%X) in response (% X)", ack, response)
	}
	return response, nil
}
