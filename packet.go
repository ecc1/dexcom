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

func (cgm *Cgm) sendPacket(pkt []byte) {
	err := cgm.conn.Send(pkt)
	cgm.SetError(err)
}

func (cgm *Cgm) receivePacket() []byte {
	header := make([]byte, 4)
	err := cgm.conn.Receive(header)
	if err != nil {
		cgm.SetError(err)
		return nil
	}
	if header[0] != startOfMessage {
		cgm.SetError(fmt.Errorf("unexpected message header % X", header))
		return nil
	}
	ack := Command(header[3])
	if ack != ACK {
		cgm.SetError(fmt.Errorf("unexpected response code %02X in header % X", ack, header))
		return nil
	}
	length := UnmarshalUint16(header[1:3])
	if length < minPacket || length > maxPacket {
		cgm.SetError(fmt.Errorf("invalid packet length %d in header % X", length, header))
		return nil
	}
	n := length - minPacket
	data := make([]byte, n+2)
	err = cgm.conn.Receive(data)
	if err != nil {
		cgm.SetError(err)
		return nil
	}
	crc := UnmarshalUint16(data[n:])
	data = data[:n]
	body := append(header, data...)
	calc := crc16(body)
	if crc != calc {
		cgm.SetError(CrcError{
			Kind:     "packet",
			Received: crc,
			Computed: calc,
			PageType: INVALID_PAGE,
			Data:     body,
		})
	}
	return data
}

// Cmd creates a Dexcom packet with the given command and parameters,
// sends it to the device, and returns the response.
func (cgm *Cgm) Cmd(cmd Command, params ...byte) []byte {
	if cgm.Error() != nil {
		return nil
	}
	pkt := marshalPacket(cmd, params)
	cgm.sendPacket(pkt)
	if cgm.Error() != nil {
		return nil
	}
	return cgm.receivePacket()
}
