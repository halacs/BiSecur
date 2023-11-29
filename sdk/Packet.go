package sdk

import (
	"bisecure/sdk/payload"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

type Packet struct {
	PacketPre
	payload payload.PayloadInterface // Command/response payload - variable size
	PacketPost
}

type PacketPre struct {
	TAG       byte   // ?
	Token     uint32 // Authentication token - 4 bytes (Session ID)
	CommandID byte   // Command ID - 1 byte
}

type PacketPost struct {
	Checksum byte // Checksum - 1 byte
}

func (p *PacketPre) getSize() int {
	return binary.Size(p)
}

func (p *PacketPost) getSize() int {
	return binary.Size(p)
}

func DecodePacket(packetLength uint16, buffer *bytes.Buffer) (*Packet, error) {
	p := Packet{}

	err := binary.Read(buffer, binary.BigEndian, &p.PacketPre)
	if err != nil {
		return nil, fmt.Errorf("failed to decode part of TransmissionContainerPre. %v", err)
	}

	payloadSize := int(packetLength) - (binary.Size(p.PacketPre) + binary.Size(p.PacketPost)) // don't read last checksum byte
	payloadBytes := make([]byte, payloadSize)
	_, err = buffer.Read(payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload part of Packet. %v", err)
	}

	// Let's decode packet payload according to Command ID
	switch p.CommandID {
	case COMMANDID_LOGIN:
		pl, err := payload.DecodeLoginPacket(payloadBytes)
		if err != nil {
			return nil, err
		}
		p.payload = pl
	default:
		p.payload = payload.EmptyPayload()
	}

	err = binary.Read(buffer, binary.BigEndian, &p.PacketPost)
	if err != nil {
		return nil, fmt.Errorf("failed to decode part of TransmissionContainerPre. %v", err)
	}

	return &p, nil
}

func (p *Packet) EncodeToHexString() (string, error) {
	byteData, err := p.encodeIntoHexOrNot(true)
	return strings.ToUpper(string(byteData)), err
}

func (p *Packet) Encode() ([]byte, error) {
	byteData, err := p.encodeIntoHexOrNot(false)
	return byteData, err
}

func (p *Packet) encodeIntoHexOrNot(encodeIntoHex bool) ([]byte, error) {
	var err error
	p.Checksum, err = p.getChecksum()
	if err != nil {
		return []byte{}, fmt.Errorf("failed to calculage packet checksum. %v", err)
	}

	buffer := new(bytes.Buffer)

	tmpBuff := new(bytes.Buffer)
	err = binary.Write(tmpBuff, binary.BigEndian, p.PacketPre)
	if err != nil {
		return []byte{}, fmt.Errorf("%v", err)
	}
	if encodeIntoHex {
		buffer.Write([]byte(hex.EncodeToString(tmpBuff.Bytes())))
	} else {
		buffer.Write(tmpBuff.Bytes())
	}

	payloadBytes := p.payload.Encode()
	err = binary.Write(buffer, binary.BigEndian, payloadBytes)
	if err != nil {
		return []byte{}, fmt.Errorf("%v", err)
	}

	tmpBuff = new(bytes.Buffer)
	err = binary.Write(tmpBuff, binary.BigEndian, p.PacketPost)
	if err != nil {
		return []byte{}, fmt.Errorf("%v", err)
	}
	if encodeIntoHex {
		buffer.Write([]byte(hex.EncodeToString(tmpBuff.Bytes())))
	} else {
		buffer.Write(tmpBuff.Bytes())
	}

	return buffer.Bytes(), nil
}

func (p *Packet) GetFrameLength() int {
	return p.PacketPre.getSize() + p.PacketPost.getSize()
}

func (p *Packet) GetLength() int {
	frameLength := p.GetFrameLength()
	dataLength := p.payload.Length()
	return frameLength + dataLength
}

func (p *Packet) GetLengthHexString() string {
	payloadLengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(payloadLengthBytes, uint16(p.GetLength()))
	return strings.ToUpper(hex.EncodeToString(payloadLengthBytes))
}

func (p *Packet) getChecksum() (byte, error) {
	value := byte(p.GetLength())

	// Summarize the head bytes
	headBuffer := new(bytes.Buffer)
	err := binary.Write(headBuffer, binary.BigEndian, p.PacketPre)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}
	for _, v := range headBuffer.Bytes() {
		value = value + v
	}

	// Summarize the payload bytes
	for _, v := range p.payload.ToByteArray() {
		value = value + v
	}

	value = value & 255 // grab only the lowest byte and drop the rest
	return value, nil
}

func (p *Packet) Equal(o *Packet) bool {
	if p.PacketPre != o.PacketPre {
		return false
	}

	if p.PacketPost != o.PacketPost {
		return false
	}

	if !bytes.Equal(p.payload.ToByteArray(), o.payload.ToByteArray()) {
		return false
	}

	return true
}
