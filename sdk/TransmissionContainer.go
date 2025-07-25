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

func (tc *TransmissionContainer) GetLength() (int, error) {
	payloadSize := tc.Packet.GetLength()
	size := int(unsafe.Sizeof(tc.TransmissionContainerPre)) + int(unsafe.Sizeof(tc.TransmissionContainerPost)) + payloadSize
	return size, nil
}

func (tc *TransmissionContainer) Encode() ([]byte, error) {
	data, err := tc.encodeToHexOrNot(true)
	return data, err
}

func (tc *TransmissionContainer) toHexString(data []byte) string {
	return strings.ToUpper(hex.EncodeToString(data))
}

func (tc *TransmissionContainer) encodeToHexOrNot(encodeIntoHex bool) ([]byte, error) {
	buffer := new(bytes.Buffer)

	bodyHex, err := tc.Packet.EncodeToHexString()
	if err != nil {
		return buffer.Bytes(), fmt.Errorf("%v", err)
	}

	tc.BodyLength = uint16(tc.Packet.GetLength()) // #nosec G115

	// Calculate checksum
	tc.Checksum, err = tc.getChecksum()
	if err != nil {
		return []byte{}, fmt.Errorf("failed to calculage packet checksum. %v", err)
	}

	tmpBuff := new(bytes.Buffer)
	err = binary.Write(tmpBuff, binary.BigEndian, tc.TransmissionContainerPre)
	if err != nil {
		return []byte{}, fmt.Errorf("%v", err)
	}
	if encodeIntoHex {
		buffer.Write([]byte(tc.toHexString(tmpBuff.Bytes())))
	} else {
		buffer.Write(tmpBuff.Bytes())
	}

	err = binary.Write(buffer, binary.BigEndian, []byte(bodyHex))
	if err != nil {
		return buffer.Bytes(), fmt.Errorf("%v", err)
	}

	tmpBuff = new(bytes.Buffer)
	err = binary.Write(tmpBuff, binary.BigEndian, tc.TransmissionContainerPost)
	if err != nil {
		return []byte{}, fmt.Errorf("%v", err)
	}
	if encodeIntoHex {
		buffer.Write([]byte(tc.toHexString(tmpBuff.Bytes())))
	} else {
		buffer.Write(tmpBuff.Bytes())
	}

	return buffer.Bytes(), nil
}

func (tc *TransmissionContainer) GetSrcMacToHexString() string {
	return strings.ToUpper(hex.EncodeToString(tc.SrcMac[:]))
}

func (tc *TransmissionContainer) GetDstMacToHexString() string {
	return strings.ToUpper(hex.EncodeToString(tc.DstMac[:]))
}

func DecodeTransmissionContainer(bufferHex *bytes.Buffer) (*TransmissionContainer, error) {
	buffer := new(bytes.Buffer)
	inputBytes, err := hex.DecodeString(bufferHex.String())
	if err != nil {
		return nil, err
	}
	_, err = buffer.Write(inputBytes)
	if err != nil {
		return nil, err
	}

	tc := TransmissionContainer{}

	err = binary.Read(buffer, binary.BigEndian, &tc.SrcMac)
	if err != nil {
		return nil, fmt.Errorf("failed to decode part of TransmissionContainerPre. %v", err)
	}

	err = binary.Read(buffer, binary.BigEndian, &tc.DstMac)
	if err != nil {
		return nil, fmt.Errorf("failed to decode part of TransmissionContainerPre. %v", err)
	}

	err = binary.Read(buffer, binary.BigEndian, &tc.BodyLength)
	if err != nil {
		return nil, fmt.Errorf("failed to decode part of  TransmissionContainerPre. %v", err)
	}

	packetLength := tc.BodyLength - 2 // ignore the last 2 checksum bytes; TODO This is really the reason to use -2 here???
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

	err = binary.Read(buffer, binary.BigEndian, &tc.Checksum)
	if err != nil {
		return nil, fmt.Errorf("failed to decode part of TransmissionContainerPost. %v", err)
	}
	// Validate checksum
	calculatedChecksum, err := tc.getChecksum()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate checksum of received Transmission Container. %v", err)
	}
	expectedChecksum := tc.Checksum
	if calculatedChecksum != expectedChecksum {
		return nil, fmt.Errorf("invalid transport container checksum. Received checksum: 0x%X, Calculated checksum: 0x%X", expectedChecksum, calculatedChecksum)
	}
	return &tc, nil
}

func (tc *TransmissionContainer) isResponse() bool {
	return (tc.Packet.CommandID & RESPONSE_MASK) == RESPONSE_MASK
}

func (tc *TransmissionContainer) isResponseFor(request *TransmissionContainer) bool {
	response := tc

	switch request.Packet.getCommandID() {
	case COMMANDID_SET_STATE:
		return response.Packet.TAG == request.Packet.TAG && response.Packet.getCommandID() == COMMANDID_HM_GET_TRANSITION && response.isResponse()
	default:
		return response.Packet.TAG == request.Packet.TAG && response.Packet.getCommandID() == request.Packet.getCommandID() && response.isResponse()
	}
}

func (tc *TransmissionContainer) getChecksum() (byte, error) {
	payloadHexBytes, err := tc.Packet.EncodeToHexString()
	if err != nil {
		return 0, err
	}

	fullTransmissionContainer := tc.GetSrcMacToHexString() + tc.GetDstMacToHexString() + tc.Packet.GetLengthHexString() + payloadHexBytes

	var b bytes.Buffer
	b.WriteString(fullTransmissionContainer)

	value := 0
	for _, v := range b.Bytes() {
		value = value + int(v)
	}

	checksum := byte(value & 255)
	return checksum, nil
}

func (tc *TransmissionContainer) Equal(o *TransmissionContainer) bool {
	if tc.TransmissionContainerPre != o.TransmissionContainerPre {
		return false
	}

	if !tc.Packet.Equal(&o.Packet) {
		return false
	}

	if tc.TransmissionContainerPost != o.TransmissionContainerPost {
		return false
	}

	return true
}

func (tc *TransmissionContainer) String() string {
	return fmt.Sprintf("SrcMAC=0x%X, DstMAC=0x%X, BodyLength=0x%X, packet=[%s], Checksum=0x%X, isResponse: %t", tc.SrcMac, tc.DstMac, tc.BodyLength, tc.Packet.String(), tc.Checksum, tc.isResponse())
}
