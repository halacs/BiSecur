package sdk

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"unsafe"
)

type TransmissionContainer struct { // Transport Container
	TransmissionContainerPre
	Packet Packet // MCP packet payload - variable size
	TransmissionContainerPost
}
type TransmissionContainerPre struct {
	SrcMac     [6]byte // Source MAC address - 6 bytes
	DstMac     [6]byte // Destination MAC address - 6 bytes
	BodyLength uint16  // Length of the payload - 2 byte
}

type TransmissionContainerPost struct {
	Checksum byte // Checksum - 1 byte
}

func (p *TransmissionContainer) GetLength() (int, error) {
	payloadSize := p.Packet.GetLength()
	size := int(unsafe.Sizeof(p.TransmissionContainerPre)) + int(unsafe.Sizeof(p.TransmissionContainerPost)) + payloadSize
	return size, nil
}

func (p *TransmissionContainer) Encode() ([]byte, error) {
	data, err := p.encodeToHexOrNot(true)
	return data, err
}

func (p *TransmissionContainer) encodeToHexOrNot(encodeIntoHex bool) ([]byte, error) {
	buffer := new(bytes.Buffer)

	bodyHex, err := p.Packet.EncodeToHexString()
	if err != nil {
		return buffer.Bytes(), fmt.Errorf("%v", err)
	}

	p.BodyLength = uint16(p.Packet.GetLength())

	// Calculate checksum
	p.Checksum, err = p.getChecksum()
	if err != nil {
		return []byte{}, fmt.Errorf("failed to calculage packet checksum. %v", err)
	}

	tmpBuff := new(bytes.Buffer)
	err = binary.Write(tmpBuff, binary.BigEndian, p.TransmissionContainerPre)
	if err != nil {
		return []byte{}, fmt.Errorf("%v", err)
	}
	if encodeIntoHex {
		buffer.Write([]byte(hex.EncodeToString(tmpBuff.Bytes())))
	} else {
		buffer.Write(tmpBuff.Bytes())
	}

	err = binary.Write(buffer, binary.BigEndian, []byte(bodyHex))
	if err != nil {
		return buffer.Bytes(), fmt.Errorf("%v", err)
	}

	tmpBuff = new(bytes.Buffer)
	err = binary.Write(tmpBuff, binary.BigEndian, p.TransmissionContainerPost)
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

func (p *TransmissionContainer) GetSrcMacToHexString() string {
	return strings.ToUpper(hex.EncodeToString(p.SrcMac[:]))
}

func (p *TransmissionContainer) GetDstMacToHexString() string {
	return strings.ToUpper(hex.EncodeToString(p.DstMac[:]))
}

func DecodeTransmissionContainer(buffer *bytes.Buffer) (*TransmissionContainer, error) {
	tc := TransmissionContainer{}

	err := binary.Read(buffer, binary.BigEndian, &tc.TransmissionContainerPre.SrcMac)
	if err != nil {
		return nil, fmt.Errorf("failed to decode part of TransmissionContainerPre. %v", err)
	}

	err = binary.Read(buffer, binary.BigEndian, &tc.TransmissionContainerPre.DstMac)
	if err != nil {
		return nil, fmt.Errorf("failed to decode part of TransmissionContainerPre. %v", err)
	}

	err = binary.Read(buffer, binary.BigEndian, &tc.TransmissionContainerPre.BodyLength)
	if err != nil {
		return nil, fmt.Errorf("failed to decode part of  TransmissionContainerPre. %v", err)
	}

	packetLength := tc.TransmissionContainerPre.BodyLength - 2 // ignore the last 2 checksum bytes; TODO This is really the reason to use -2 here???
	packetBytes := make([]byte, packetLength)
	_, err = buffer.Read(packetBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read packet data. %v", err)
	}

	packetBuffer := new(bytes.Buffer)
	_, err = packetBuffer.Write(packetBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write packet data into buffer. %v", err)
	}

	packet, err := DecodePacket(packetLength, packetBuffer)
	if packet == nil || err != nil {
		return nil, fmt.Errorf("failed to decode packet data. %v", err)
	}
	tc.Packet = *packet

	err = binary.Read(buffer, binary.BigEndian, &tc.TransmissionContainerPost.Checksum)
	if err != nil {
		return nil, fmt.Errorf("failed to decode part of TransmissionContainerPost. %v", err)
	}

	// Validate checksum
	calculatedChecksum, err := tc.getChecksum()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate checksum of received Transmission Container. %v", err)
	}
	actualChecksum := tc.Checksum
	if calculatedChecksum != actualChecksum {
		return nil, fmt.Errorf("invalid checksum on transport container. Expected checksum value: %d, Actual checksum value: %d", calculatedChecksum, actualChecksum)
	}

	return &tc, nil
}

func (p *TransmissionContainer) getChecksum() (byte, error) {
	payloadHexBytes, err := p.Packet.EncodeToHexString()
	if err != nil {
		return 0, err
	}

	fullTransmissionContainer := p.GetSrcMacToHexString() + p.GetDstMacToHexString() + p.Packet.GetLengthHexString() + payloadHexBytes

	var b bytes.Buffer
	b.WriteString(fullTransmissionContainer)

	value := 0
	for _, v := range b.Bytes() {
		value = value + int(v)
	}

	checksum := byte(value & 255)
	return checksum, nil
}

func (p *TransmissionContainer) Equal(o *TransmissionContainer) bool {
	if p.TransmissionContainerPre != o.TransmissionContainerPre {
		return false
	}

	if !p.Packet.Equal(&o.Packet) {
		return false
	}

	if p.TransmissionContainerPost != o.TransmissionContainerPost {
		return false
	}

	return true
}
